FROM golang:alpine as builder

RUN adduser -D -g '' podcastMaker

WORKDIR /
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /podcastMaker

FROM alpine

RUN apk update && apk add --no-cache ffmpeg curl python2
RUN curl -L https://yt-dl.org/downloads/latest/youtube-dl -o /usr/local/bin/youtube-dl
RUN chmod a+rx /usr/local/bin/youtube-dl

VOLUME [ "/downloads" ]

WORKDIR /

COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /podcastMaker /podcastMaker

ENTRYPOINT ["/podcastMaker"]
