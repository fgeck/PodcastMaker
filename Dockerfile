FROM golang:alpine as builder

RUN adduser -D -g '' mixcloudPodcaster

WORKDIR /
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build  -o /mixcloudPodcaster

FROM alpine

RUN apk update && apk add --no-cache ffmpeg curl python3 py3-pip
RUN curl -L https://yt-dl.org/downloads/latest/youtube-dl -o /usr/local/bin/youtube-dl
RUN chmod a+rx /usr/local/bin/youtube-dl

VOLUME [ "/downloads" ]

WORKDIR /

COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /mixcloudPodcaster /mixcloudPodcaster

ENV VERBOSE=false
ENV PARALLEL=false

ENTRYPOINT ["/mixcloudPodcaster"]