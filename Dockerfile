FROM golang:1.14.1-alpine AS build

ENV  CGO_ENABLED 0
WORKDIR /code
ADD  . ./
RUN  go install

FROM alpine:3.11.5
RUN apk add --no-cache ca-certificates
COPY --from=build /go/bin/hls-rip /usr/bin/hls-rip
ENTRYPOINT ["/usr/bin/hls-rip"]
