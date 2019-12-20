FROM golang:1.13.5-alpine AS build

ENV  CGO_ENABLED 0
WORKDIR /code
ADD  . ./
RUN  go install

FROM alpine:3.11.0
RUN apk add --no-cache ca-certificates
COPY --from=build /go/bin/hls-rip /usr/bin/hls-rip
ENTRYPOINT ["/usr/bin/hls-rip"]
