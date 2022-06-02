package main

import (
	"log"
	"net/http"
	"time"

	"github.com/floge77/PodcastMaker/config"
	"github.com/floge77/PodcastMaker/downloader"
	"github.com/floge77/PodcastMaker/handler"
	"github.com/floge77/PodcastMaker/podcast"
)

func main() {
	var err error
	var allPodcasts []*podcast.PodcastFeed

	cfgReader := config.NewConfigReader()
	cfg, err := cfgReader.Read("config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	// download and create routes on start - then every 10h
	downloader := downloader.NewDownloader(cfg)
	err = downloader.Download()
	if err != nil {
		log.Fatal(err)
	}
	builder := podcast.NewPodcastBuilder(cfg)
	allPodcasts, err = builder.WriteAllPodcastFeedsXml()
	if err != nil {
		log.Fatal(err)
	}
	router := handler.NewRouter(cfg)
	router.ServeAllPodcasts(allPodcasts)

	// download podcasts and update routes every 10h
	go func() {
		for range time.Tick(10 * time.Hour) {
			err = downloader.Download()
			if err != nil {
				log.Fatal(err)
			}
			allPodcasts, err = builder.WriteAllPodcastFeedsXml()
			if err != nil {
				log.Fatal(err)
			}
			router.ServeAllPodcasts(allPodcasts)
		}
	}()

	server := &http.Server{
		Handler:      router.Router(),
		Addr:         ":" + "80",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Println("Router running at Port 80")
	log.Fatal(server.ListenAndServe())
}
