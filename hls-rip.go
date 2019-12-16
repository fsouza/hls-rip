// Copyright 2017 Francisco Souza. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"github.com/grafov/m3u8"
)

var (
	nworkers   int
	workingDir string
)

func main() {
	flag.IntVar(&nworkers, "w", 16, "Number of workers for parallel download of segments")
	flag.StringVar(&workingDir, "wd", "", "Working dir. Defaults to the current directory")
	flag.Parse()

	for _, arg := range flag.Args() {
		err := rip(arg)
		if err != nil {
			log.Printf("failed to rip %q: %s", arg, err)
		}
	}
}

func rip(playlistURL string) error {
	folderName, err := getFolderName(playlistURL)
	if err != nil {
		return err
	}
	return ripPlaylist(playlistURL, folderName)
}

func ripPlaylist(url, folder string) error {
	baseURL, fileName := path.Split(url)
	filePath := filepath.Join(folder, fileName)
	err := download(url, filePath)
	if err != nil {
		return err
	}

	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	playlist, listType, err := m3u8.DecodeFrom(f, true)
	if err != nil {
		return err
	}
	switch listType {
	case m3u8.MEDIA:
		return ripSegments(playlist.(*m3u8.MediaPlaylist), baseURL, folder)
	case m3u8.MASTER:
		master := playlist.(*m3u8.MasterPlaylist)
		for _, variant := range master.Variants {
			variantURL := baseURL + "/" + variant.URI
			err = ripPlaylist(variantURL, folder)
			if err != nil {
				return fmt.Errorf("failed to rip variant %q: %s", variant.URI, err)
			}
		}
	}
	return nil
}

func ripSegments(p *m3u8.MediaPlaylist, baseURL, folder string) error {
	var wg sync.WaitGroup

	ch := make(chan m3u8.MediaSegment, nworkers)
	errs := make(chan error, nworkers)

	wg.Add(nworkers)
	for i := 0; i < nworkers; i++ {
		go func() {
			defer wg.Done()
			for seg := range ch {
				filePath := filepath.Join(folder, seg.URI)
				url := baseURL + "/" + seg.URI
				err := download(url, filePath)
				if err != nil {
					errs <- err
					return
				}
			}
		}()
	}

	for _, seg := range p.Segments {
		if seg != nil {
			ch <- *seg
		}
	}

	close(ch)
	wg.Wait()
	close(errs)
	return <-errs
}

func download(url, path string) error {
	dir, _ := filepath.Split(path)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}

	//nolint:gosec
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, resp.Body)

	return err
}

func getFolderName(playlistURL string) (string, error) {
	parts := strings.Split(playlistURL, "/")
	return filepath.Abs(filepath.Join(workingDir, parts[len(parts)-2]))
}
