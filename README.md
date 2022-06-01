# PodcastMaker

## Configure

## Run

### Docker Run

```bash
    make docker-build
    docker run -d -p 8888:80 -v $(pwd)/config.yaml:/config.yaml -v $(pwd)/downloads:/downloads --name podcastmaker floge77/podcastmaker
```

### Docker-Compose

```bash
    make docker-build
    docker-compose up -d 
```
