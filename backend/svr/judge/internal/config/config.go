package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

var (
	MyConfig *Config
)

func init() {
	MyConfig = new(Config)
}

func (c *Config) ReadConfigFile(path string) error {
	fileData, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(fileData, c)
	if err != nil {
		return err
	}

	return nil
}
