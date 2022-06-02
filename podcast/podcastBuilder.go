package podcast

import (
	"fmt"
	"net/url"
	"os"
	"path"
	"strings"

	podcastFeed "github.com/eduncan911/podcast"
	"github.com/floge77/PodcastMaker/config"
)

const podcastXmlName = "podcast.xml"

type PodcastBuilder interface {
	WriteAllPodcastFeedsXml() ([]*PodcastFeed, error)
	CreateSinglePodcastFeed(dirName string) (*PodcastFeed, error)
}

type DefaultPodcastBuilder struct {
	infoReader PodcastFileReader
	config     *config.Config
}

func NewPodcastBuilder(config *config.Config) PodcastBuilder {
	return &DefaultPodcastBuilder{infoReader: &DefaultPodcastFileReader{cfg: config}, config: config}
}

func (b *DefaultPodcastBuilder) CreateSinglePodcastFeed(dirName string) (*PodcastFeed, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	defer os.Chdir(pwd)
	err = os.Chdir(b.config.General.DownloadDir)
	if err != nil {
		return nil, err
	}
	return b.createPodcastFeed(dirName)
}

func (b *DefaultPodcastBuilder) WriteAllPodcastFeedsXml() ([]*PodcastFeed, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	defer os.Chdir(pwd)
	err = os.Chdir(b.config.General.DownloadDir)
	if err != nil {
		return nil, err
	}
	providerDirectories, err := ReadDir(b.config.General.DownloadDir)
	if err != nil {
		return nil, err
	}
	allFeeds := []*PodcastFeed{}
	for _, provider := range providerDirectories { // e.g. youtube, mixcloud
		uploaderDirectories, err := ReadDir(path.Join(b.config.General.DownloadDir, provider))
		if err != nil {
			return nil, err
		}
		for _, uploader := range uploaderDirectories { // different channels
			feed, err := b.CreateSinglePodcastFeed(path.Join(provider, uploader))
			if err != nil {
				return nil, err
			}
			allFeeds = append(allFeeds, feed)
			file, err := os.Create(path.Join(provider, uploader, podcastXmlName))
			if err != nil {
				return nil, err
			}
			err = feed.Feed.Encode(file)
			if err != nil {
				return nil, err
			}
		}
	}
	return allFeeds, nil
}

func (b *DefaultPodcastBuilder) createPodcastFeed(dirName string) (*PodcastFeed, error) {
	entries, err := b.infoReader.GetPodcastItemsInformationForSingleDir(dirName)
	if err != nil {
		return nil, err
	}
	// dirName is e.g. mixcloud.com/q-dance
	provider := strings.Split(dirName, "/")[0]
	channel := strings.Split(dirName, "/")[1]
	feed := &PodcastFeed{Feed: podcastFeed.New(
		fmt.Sprintf("%s - %s - Podcast", channel, provider),
		fmt.Sprintf("https://%s/%s", provider, channel),
		"",
		nil,
		nil,
	), Path: dirName}
	feed.Feed.IAuthor = dirName
	for _, entry := range entries.Items {
		item := podcastFeed.Item{
			Title:       entry.Title,
			Description: entry.Channel,
			Author:      &podcastFeed.Author{Name: entry.Artist},
			PubDate:     entry.ReleaseDate,
		}
		downloadURL := &url.URL{
			Scheme: "http",
			Host:   b.config.General.HostName, // + ":80",
			Path:   path.Join("/downloads", dirName, entry.FileName),
		}
		item.AddEnclosure(downloadURL.String(), podcastFeed.MP3, entry.FileSize)
		_, err = feed.Feed.AddItem(item)
		if err != nil {
			return nil, err
		}
	}
	return feed, nil
}
