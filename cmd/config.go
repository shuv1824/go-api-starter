package cmd

import (
	"os"

	"gopkg.in/yaml.v3"
)

type modeType string

const (
	modeTypeDebug  modeType = "debug"
	modeTypeTest   modeType = "test"
	modeTypeReleae modeType = "release"
)

type Config struct {
	Mode modeType `yaml:"mode"`
}

func InitConfig(filePath string) (*Config, error) {
	cfg := Config{
		Mode: modeTypeDebug,
	}

	fileBytes, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(fileBytes, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
