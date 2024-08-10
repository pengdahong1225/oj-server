package settings

import (
	"github.com/pengdahong1225/Oj-Online-Server/config"
)

var Conf *config.AppConfig

func init() {
	Conf = new(config.AppConfig)
}

func Init() error {
	return config.LoadConfig(Conf)
}
