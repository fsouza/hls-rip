FROM golang:1.16.3-alpine AS build

ENV  CGO_ENABLED 0
WORKDIR /code
ADD  . ./
RUN  go install

FROM alpine:3.13.4
RUN apk add --no-cache ca-certificates
COPY --from=build /go/bin/hls-rip /usr/bin/hls-rip
ENTRYPOINT ["/usr/bin/hls-rip"]
