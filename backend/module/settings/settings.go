package settings

import (
	"errors"
	"github.com/fsnotify/fsnotify"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"sync"
)

var (
	conf *AppConfig
	once sync.Once
)

func Instance() *AppConfig {
	once.Do(func() {
		conf = new(AppConfig)
	})
	return conf
}

func (receiver *AppConfig) LoadConfig() error {
	// 读取环境变量
	logMode := os.Getenv("LOG_MODE")
	// 读取配置
	viperConfig := viper.New()
	if logMode == "release" {
		viperConfig.SetConfigName("prod")
	} else {
		viperConfig.SetConfigName("dev")
	}
	viperConfig.AddConfigPath("config")
	viperConfig.SetConfigType("yaml")
	if e := viperConfig.ReadInConfig(); e != nil {
		return e
	}
	if e := viperConfig.Unmarshal(conf); e != nil {
		return e
	}
	// 配置热重载
	viper.OnConfigChange(func(event fsnotify.Event) {
		logrus.Infoln("config file changed:", event.Name)
		if e := viperConfig.Unmarshal(conf); e != nil {
			logrus.Infoln("config file update failed:", event.Name)
		}
	})
	// 监听配置文件
	viper.WatchConfig()
	return nil
}

func (receiver *AppConfig) GetSystemConf(name string) (*SystemConfig, error) {
	for _, v := range receiver.SystemConfigs {
		if v.Name == name {
			return &v, nil
		}
	}
	return nil, errors.New("system config not found")
}
