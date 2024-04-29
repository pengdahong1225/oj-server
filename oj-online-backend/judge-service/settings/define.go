package settings

// 从viper反射到数据模型，需要设置`mapstructure`反射字段
type AppConfig struct {
	*SystemConfig   `mapstructure:"system"`
	*RedisConfig    `mapstructure:"redis"`
	*RegistryConfig `mapstructure:"registry"`
	*LogConfig      `mapstructure:"log"`
	*MqConfig       `mapstructure:"mq"`
}

type SystemConfig struct {
	Name string `mapstructure:"name"`
	Port int    `mapstructure:"port"`
}

type RegistryConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type RedisConfig struct {
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
