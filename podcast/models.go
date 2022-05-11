package podcast

import (
	"time"

	podcastFeed "github.com/eduncan911/podcast"
)

type PodcastEntries struct {
	Items []*PodcastItem
}

type PodcastFeed struct {
	Feed podcastFeed.Podcast
	Path string
}

type PodcastItem struct {
	Title       string
	Artist      string
	Channel     string
	ReleaseDate *time.Time
	FileSize    int64
	FileName    string
}
