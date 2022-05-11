package download

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/floge77/mixcloudPodcaster/config"
)

var (
	DownloadDir  = "/Users/d068994/SAPDevelop/github.com/mixcloudPodacaster/testdownloads" //"/downloads"
	outputFormat = "%%(title)s%s%%(artist)s%s%%(uploader)s%s%%(upload_date)s.%%(ext)s"
)

const ()

type Downloader interface {
	Download() error
}

type YoutubeDlDownloader struct{ cfg *config.Config }

func NewDownloader(cfg *config.Config) Downloader {
	return &YoutubeDlDownloader{cfg: cfg}
}

func (d *YoutubeDlDownloader) Download() error {
	err := d.buildDestinationDirs()
	if err != nil {
		return err
	}
	for _, podcastConfig := range d.cfg.Podcasts {
		err := d.downloadChannel(podcastConfig)
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *YoutubeDlDownloader) buildDestinationDirs() error {
	for _, podcastConfig := range d.cfg.Podcasts {
		dir := filepath.Join(DownloadDir, podcastConfig.Provider, podcastConfig.Channel)
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return err
		}
		log.Printf("created downloadDir: %q", dir)
	}
	return nil
}

func (d *YoutubeDlDownloader) downloadChannel(podcastConfig *config.PodcastConfig) error {
	// ["youtube-dl", "-x", "-i", "--dateafter", "now-12months", "--audio-format", "mp3", "--embed-thumbnail", "--add-metadata", "--match-filter", f"duration>{lengthFilter}"]
	// ["--download-archive",	f"/downloads/{podcast['channelName']}/archive.txt", "-o", f"/downloads/{podcast['channelName']}/%(title)s__%(uploader)s__%(upload_date)s.%(ext)s", podcast['playlistToDownloadURL']]
	args := []string{"-x", "-i", "--dateafter", "now-3months", "--audio-format", "mp3",
		"--embed-thumbnail", "--add-metadata", "--match-filter", fmt.Sprintf("duration>%d", d.cfg.General.MinimalLengthMin*60), "--no-progress",
		"--download-archive", fmt.Sprintf("%s/%s/%s/archive.txt", DownloadDir, podcastConfig.Provider, podcastConfig.Channel), "-o",
		fmt.Sprintf("%s/%s/%s/%s", DownloadDir, podcastConfig.Provider, podcastConfig.Channel, fmt.Sprintf(outputFormat, d.cfg.General.Separator, d.cfg.General.Separator, d.cfg.General.Separator)), podcastConfig.DownloadUrl}
	cmd := exec.Command("youtube-dl", args...)

	log.Printf("downloading channel %q from %q", podcastConfig.Channel, podcastConfig.Provider)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}
