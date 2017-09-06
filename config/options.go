package config

import (
	"blog/services"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Options struct {
	Database *services.DBOptions `yaml:"database"`
	ScoreRule *services.ScoreRule `yaml:"score"`
}

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

