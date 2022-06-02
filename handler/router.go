package handler

import (
	"log"
	"net/http"
	"path/filepath"

	"github.com/floge77/PodcastMaker/config"
	"github.com/floge77/PodcastMaker/podcast"
	"github.com/gorilla/mux"
)

const podcastXmlName = "podcast.xml"

type Router interface {
	ServeAllPodcasts([]*podcast.PodcastFeed)
	ServeSinglePodcast(podcastFeed *podcast.PodcastFeed)
	Router() *mux.Router
}

type DefaultRouter struct {
	router *mux.Router
	config *config.Config
}

func NewRouter(config *config.Config) Router {
	router := mux.NewRouter()
	router.PathPrefix("/downloads/").Handler(http.StripPrefix("/downloads/", http.FileServer(http.Dir(config.General.DownloadDir+"/"))))
	router.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		_, _ = writer.Write([]byte("hello world"))
	})
	return &DefaultRouter{
		router: router,
		config: config,
	}
}

func (r *DefaultRouter) Router() *mux.Router {
	return r.router
}

func (r *DefaultRouter) ServeAllPodcasts(allPodcasts []*podcast.PodcastFeed) {
	for _, podcastFeed := range allPodcasts {
		r.ServeSinglePodcast(podcastFeed)
	}
}

func (r *DefaultRouter) ServeSinglePodcast(podcastFeed *podcast.PodcastFeed) {
	podcastXmlPath := filepath.Join(r.config.General.DownloadDir, podcastFeed.Path, podcastXmlName)
	handleFunc := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/xml")

		if !isAuthenticated() {
			w.Header().Set("WWW-Authenticate", "Basic")
			w.WriteHeader(http.StatusUnauthorized)
			http.Error(w, "you are not logged in", http.StatusUnauthorized)
			return
		}
		http.ServeFile(w, r, podcastXmlPath)
	}
	route := "/podcasts/" + podcastFeed.Feed.IAuthor
	log.Printf("Serving a podcast at %s", route)
	r.router.HandleFunc(route, handleFunc).Methods("GET")
}

func isAuthenticated() bool {
	return true
}
