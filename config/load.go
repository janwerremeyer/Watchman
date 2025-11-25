package config

import (
	"log/slog"
	"os"
	"watchman/internal/gofy"

	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v3"
)

func LoadConfig() (*Config, error) {

	awsConfig, err := loadAwsConfig()
	if err != nil {
		return nil, err
	}

	wantedConfig, err := loadWantedConfig()

	return &Config{
		AWS:    *awsConfig,
		Wanted: wantedConfig,
	}, nil
}

func loadWantedConfig() ([]Wanted, error) {
	paths := [...]string{"./"}
	filename := "wanted.yaml"
	var cfg []Wanted

	for _, path := range paths {
		fullPath := path + filename
		if gofy.FileExists(fullPath) {
			data, err := os.ReadFile(filename)
			if err != nil {
				return nil, err
			}

			if err := yaml.Unmarshal(data, &cfg); err != nil {
				return nil, err
			}

			for _, c := range cfg {
				if c.Environment == nil {
					c.Environment = make(map[string]string)
				}
			}

			return cfg, nil
		}
	}

	slog.Error("AWS config file not found in any of the expected locations")

	return nil, os.ErrNotExist
}

func loadAwsConfig() (*AWSConfig, error) {
	paths := [...]string{"./"}
	filename := "aws.toml"
	var cfg AWSConfig

	for _, path := range paths {
		fullPath := path + filename
		if gofy.FileExists(fullPath) {
			if _, err := toml.DecodeFile(fullPath, &cfg); err != nil {
				return nil, err
			}
			return &cfg, nil
		}
	}

	slog.Error("AWS config file not found in any of the expected locations")

	return nil, os.ErrNotExist
}
