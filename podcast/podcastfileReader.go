package podcast

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/floge77/mixcloudPodcaster/config"
)

type PodcastFileReader interface {
	GetPodcastItemsInformationForSingleDir(dir string) (*PodcastEntries, error)
}

type DefaultPodcastFileReader struct {
	cfg *config.Config
}

func NewPodcastFileReader(cfg *config.Config) PodcastFileReader {
	return &DefaultPodcastFileReader{cfg: cfg}
}

func (f *DefaultPodcastFileReader) GetPodcastItemsInformationForSingleDir(dir string) (*PodcastEntries, error) {
	fileNames, err := ReadDir(dir)
	if err != nil {
		return nil, err
	}
	podcastEntries := &PodcastEntries{}
	for _, name := range fileNames {
		// create Info structs for podcastsItems
		if strings.Contains(name, ".mp3") {
			item := f.getPodcastItemInfosFromFileName(dir, name)
			podcastEntries.Items = append(podcastEntries.Items, item)
		}
	}
	return podcastEntries, nil
}

func (f *DefaultPodcastFileReader) getPodcastItemInfosFromFileName(dir string, filename string) *PodcastItem {
	s := strings.Replace(filename, ".mp3", "", -1)
	fields := strings.Split(s, f.cfg.General.Separator)
	item := &PodcastItem{}
	item.Title = fields[0]
	item.Artist = fields[1]
	item.Channel = fields[2]
	item.ReleaseDate = f.getReleaseDateFromString(fields[3])
	fileSize, _ := f.extractFileSize(dir, filename)
	item.FileSize = fileSize
	item.FileName = filename
	return item
}

func (*DefaultPodcastFileReader) extractFileSize(dir string, filename string) (int64, error) {
	file, err := os.Stat(dir + "/" + filename)
	if err != nil {
		return 0, err
	}
	return file.Size(), nil
}

func (*DefaultPodcastFileReader) getReleaseDateFromString(date string) *time.Time {
	t, _ := time.Parse("20060102", date)
	return &t
}

func ReadDir(dirname string) ([]string, error) {
	file, err := os.Open(dirname)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	list, err := file.Readdirnames(0) // 0 to read all files and folders
	if err != nil {
		fmt.Printf("Could not read directory %v", dirname)
	}
	return list, nil
}
