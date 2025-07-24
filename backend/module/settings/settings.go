package settings

import (
	"github.com/fsnotify/fsnotify"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"path"
)

var (
	AppConf *AppConfig
)

func init() {
	AppConf = new(AppConfig)
}

type AppConfig struct {
	SandBoxCfg  []*SandBox `mapstructure:"sandbox"`
	MysqlCfg    *Mysql     `mapstructure:"mysql"`
	RedisCfg    *Redis     `mapstructure:"redis"`
	RegistryCfg *Registry  `mapstructure:"registry"`
	MQCfg       *MQ        `mapstructure:"rabbitmq"`
	JwtCfg      *Jwt       `mapstructure:"jwt"`
	SmsCfg      *Sms       `mapstructure:"sms"`
}

func (receiver *AppConfig) LoadConfig(config string) error {
	filePath := path.Dir(config)
	fileName := path.Base(config)

	// 读取配置
	v := viper.New()
	v.AddConfigPath(filePath)
	v.SetConfigName(fileName)
	v.SetConfigType("yaml")
	if e := v.ReadInConfig(); e != nil {
		return e
	}
	if e := v.Unmarshal(AppConf); e != nil {
		return e
	}
	// 配置热重载
	viper.OnConfigChange(func(event fsnotify.Event) {
		logrus.Errorf("config file changed: %s", event.Name)
		if e := v.Unmarshal(AppConf); e != nil {
			logrus.Errorf("config file update failed: %s", event.Name)
		}
	})
	// 监听配置文件
	viper.WatchConfig()
	return nil
}
