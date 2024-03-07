package global

// 从viper反射到数据模型，需要设置`mapstructure`反射字段
type config struct {
	System_   systemConfig   `mapstructure:"system"`
	Registry_ registryConfig `mapstructure:"registry"`
	Sql_      mysqlConfig    `mapstructure:"mysql"`
	Redis_    redisConfig    `mapstructure:"redis"`
	Log_      logConfig      `mapstructure:"log"`
}

type systemConfig struct {
	Name string `mapstructure:"name"`
	Port int    `mapstructure:"port"`
}

type registryConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type mysqlConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
	Db   string `mapstructure:"db"`
	User string `mapstructure:"user"`
	Pwd  string `mapstructure:"password"`
}

type redisConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type logConfig struct {
	Path  string `mapstructure:"path"`
	Level string `mapstructure:"level"`
}
