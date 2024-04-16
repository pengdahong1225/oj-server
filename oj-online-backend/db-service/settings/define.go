package settings

// 从viper反射到数据模型，需要设置`mapstructure`反射字段
type AppConfig struct {
	*SystemConfig   `mapstructure:"system"`
	*MysqlConfig    `mapstructure:"mysql"`
	*RedisConfig    `mapstructure:"redis"`
	*RegistryConfig `mapstructure:"registry"`
	*LogConfig      `mapstructure:"log"`
}

type SystemConfig struct {
	Name string `mapstructure:"name"`
	Port int    `mapstructure:"port"`
}

type RegistryConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
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

type LogConfig struct {
	Path  string `mapstructure:"path"`
	Level string `mapstructure:"level"`
}
