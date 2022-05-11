package handler

import (
	"log"
	"net/http"

	"github.com/floge77/mixcloudPodcaster/config"
	"github.com/floge77/mixcloudPodcaster/podcast"
	"github.com/gorilla/mux"
)

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
	handleFunc := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/xml")

		if err := podcastFeed.Feed.Encode(w); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
	route := "/podcasts/" + podcastFeed.Feed.IAuthor
	log.Printf("Serving a podcast at %s", route)
	r.router.HandleFunc(route, handleFunc).Methods("GET")
}
