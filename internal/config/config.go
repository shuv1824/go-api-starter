package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type ModeType string

const (
	ModeTypeDebug  ModeType = "debug"
	ModeTypeTest   ModeType = "test"
	ModeTypeReleae ModeType = "release"
)

type DatabaseConfig struct {
	Type     string `yaml:"type"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	DbName   string `yaml:"dbname"`
	SSLMode  string `yaml:"sslmode"`
}

type Config struct {
	Mode     ModeType       `yaml:"mode"`
	Port     int            `yaml:"port"`
	Database DatabaseConfig `yaml:"database"`
}

func InitConfig(filePath string) (*Config, error) {
	cfg := Config{
		Mode: ModeTypeDebug,
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
