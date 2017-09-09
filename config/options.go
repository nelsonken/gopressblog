package config

import (
	"blog/services"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// Options services's options
type Options struct {
	Database  *services.DBOptions     `yaml:"database"`
	ScoreRule *services.ScoreRule     `yaml:"score"`
	Elastic   *services.ElasticOption `yaml:"elastic"`
}

// GetConfig getconfig from file
func GetConfig(configFile string, opts *Options) {
	options, err := ioutil.ReadFile(configFile)
	if err != nil {
		panic(err)
	}

	err = yaml.Unmarshal(options, opts)
	if err != nil {
		panic(err)
	}
}
