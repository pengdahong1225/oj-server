package configs

import (
	"gopkg.in/yaml.v3"
	"os"
)

var (
	ServerConf *ServerConfig
	AppConf    *AppConfig
)

func init() {
	ServerConf = new(ServerConfig)
	AppConf = new(AppConfig)
}

func LoadServerConfigFile(path string) error {
	// 读取文件
	fileData, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(fileData, ServerConf)
	if err != nil {
		return err
	}

	return nil
}

func LoadAppConfigFile(path string) error {
	// 读取文件
	fileData, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(fileData, AppConf)
	if err != nil {
		return err
	}

	return nil
}
