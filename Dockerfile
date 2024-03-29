FROM golang:alpine as builder

RUN adduser -D -g '' podcastMaker

WORKDIR /
COPY downloader/ config/ handler/ podcast/ main.go go.mod go.sum .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /podcastMaker

FROM amd64/alpine

RUN apk update && apk add --no-cache ffmpeg curl python3 && ln -sf python3 /usr/bin/python
RUN curl -L https://yt-dl.org/downloads/latest/youtube-dl -o /usr/local/bin/youtube-dl
RUN chmod a+rx /usr/local/bin/youtube-dl

VOLUME [ "/downloads" ]

WORKDIR /

COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /podcastMaker /podcastMaker

ENTRYPOINT ["/podcastMaker"]
