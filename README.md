# hls-rip

Tool for ripping m3u8 playlists. It downloads playlists sequentially and
segments in parallel. The number of parallel workers can be controlled with the
``-w`` flag.

```
% hls-rip -w 64 https://website.example.com/videos/2017/interesting-video.m3u8
```

It requires variation playlists and segments to be declared using relative
paths.

## Installing

Installing this requires [latest version of Go properly installed and
configured](https://golang.org/doc/install), then just ``go get`` should be
enough:

```
% go get github.com/fsouza/hls-rip
```
