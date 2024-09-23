package settings

import "errors"

// 从viper反射到数据模型，需要设置`mapstructure`反射字段
type AppConfig struct {
	SystemConfigs   []SystemConfig `mapstructure:"system"`
	*MysqlConfig    `mapstructure:"mysql"`
	*RedisConfig    `mapstructure:"redis"`
	*LogConfig      `mapstructure:"log"`
	*RegistryConfig `mapstructure:"registry"`
	*MqConfig       `mapstructure:"rabbitmq"`
	*JwtConfig      `mapstructure:"jwt"`
	*SmsConfig      `mapstructure:"sms"`
}

type SystemConfig struct {
	Name string `mapstructure:"name"`
	Port int    `mapstructure:"port"`
	Host string `mapstructure:"host"`
}

type MysqlConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
	Db   string `mapstructure:"db"`
	User string `mapstructure:"user"`
	Pwd  string `mapstructure:"password"`
}

type RedisConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type JwtConfig struct {
	SigningKey string `mapstructure:"key"`
}

type SmsConfig struct {
	AccessKeyId     string `mapstructure:"accessKeyId"`
	AccessKeySecret string `mapstructure:"accessKeySecret"`
	Endpoint        string `mapstructure:"endpoint"`
	SignName        string `mapstructure:"signName"`
	TemplateCode    string `mapstructure:"templateCode"`
}

type RegistryConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type LogConfig struct {
	Path  string `mapstructure:"path"`
	Level string `mapstructure:"level"`
}

type MqConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	PassWord string `mapstructure:"password"`
	VHost    string `mapstructure:"vhost"`
}

func (receiver *AppConfig) GetSystemConf(name string) (*SystemConfig, error) {
	for _, v := range receiver.SystemConfigs {
		if v.Name == name {
			return &v, nil
		}
	}
	return nil, errors.New("system config not found")
}
