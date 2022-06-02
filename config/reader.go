package config

import (
	"io/ioutil"
	"net/url"
	"strings"

	"gopkg.in/yaml.v2"
)

type Config struct {
	General struct {
		MinimalLengthMin      int    `yaml:"minimalLengthMin"`
		DownloadIntervalHours int    `yaml:"downloadIntervalHours"`
		DownloadDir           string `yaml:"downloadDir"`
		HostName              string `yaml:"hostName"`
		Separator             string `yaml:"separator"`
	} `yaml:"general"`
	Podcasts []*PodcastConfig
}

type PodcastConfig struct {
	Provider    string
	Channel     string
	DownloadUrl string
}

type configHelper struct {
	General struct {
		MinimalLengthMin      int    `yaml:"minimalLengthMin"`
		DownloadIntervalHours int    `yaml:"downloadIntervalHours"`
		DownloadDir           string `yaml:"downloadDir"`
		HostName              string `yaml:"hostName"`
		Separator             string `yaml:"separator"`
	} `yaml:"general"`
	Podcasts []string `yaml:"podcasts"`
}

type ConfigReader interface {
	Read(string) (*Config, error)
}
type DefaultConfigReader struct{}

func NewConfigReader() ConfigReader {
	return &DefaultConfigReader{}
}

func (r *DefaultConfigReader) Read(path string) (*Config, error) {
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var configHelper *configHelper
	err = yaml.Unmarshal(yamlFile, &configHelper)
	if err != nil {
		return nil, err
	}
	config := &Config{General: configHelper.General, Podcasts: []*PodcastConfig{}}

	for _, podcast := range configHelper.Podcasts {
		downloadUrl, err := url.Parse(podcast)
		if err != nil {
			return nil, err
		}
		//regex := regexp.MustCompile(`\.?([^.]*.com)`)
		//provider := regex.FindStringSubmatch(podcast)[1]
		name := strings.Replace(downloadUrl.Path, "/", "", -1)
		config.Podcasts = append(config.Podcasts, &PodcastConfig{
			Provider:    downloadUrl.Host,
			Channel:     name,
			DownloadUrl: podcast})
	}
	return config, nil
}
