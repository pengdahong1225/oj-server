package global

// 从viper反射到数据模型，需要设置`mapstructure`反射字段
type config struct {
	System_    systemConfig   `mapstructure:"system"`
	Redis_     redisConfig    `mapstructure:"redis"`
	JWTConfig_ jwtConfig      `mapstructure:"jwt"`
	SMS_       sms            `mapstructure:"sms"`
	Registry_  registryConfig `mapstructure:"registry"`
}

type systemConfig struct {
	Name   string `mapstructure:"name"`
	Port   int    `mapstructure:"port"`
	Prefix string `mapstructure:"prefix"`
}

type redisConfig struct {
	Ip   string `mapstructure:"ip"`
	Port int    `mapstructure:"port"`
}

type jwtConfig struct {
	SigningKey string `mapstructure:"key"`
}

type sms struct {
	AccessKeyId     string `mapstructure:"accessKeyId"`
	AccessKeySecret string `mapstructure:"accessKeySecret"`
	Endpoint        string `mapstructure:"endpoint"`
	SignName        string `mapstructure:"signName"`
	TemplateCode    string `mapstructure:"templateCode"`
}

type registryConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}