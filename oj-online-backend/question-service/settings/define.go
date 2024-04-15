package settings

// 从viper反射到数据模型，需要设置`mapstructure`反射字段
type AppConfig struct {
	*SystemConfig   `mapstructure:"system"`
	*RedisConfig    `mapstructure:"redis"`
	*JwtConfig      `mapstructure:"jwt"`
	*SmsConfig      `mapstructure:"sms"`
	*LogConfig      `mapstructure:"log"`
	*RegistryConfig `mapstructure:"registry"`
	*MqConfig       `mapstructure:"rabbitmq"`
}

type SystemConfig struct {
	Name string `mapstructure:"name"`
	Port int    `mapstructure:"port"`
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
