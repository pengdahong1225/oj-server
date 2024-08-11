package setting

import (
	"github.com/pengdahong1225/Oj-Online-Server/pkg/settings"
	"sync"
)

var (
	conf *settings.AppConfig
	once sync.Once
)

func Instance() *settings.AppConfig {
	once.Do(func() {
		conf = new(settings.AppConfig)
		err := settings.LoadConfig(conf)
		if err != nil {
			panic(err)
		}
	})
	return conf
}
