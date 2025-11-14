package configs

type AppConfig struct {
	MysqlCfg    *Mysql    `yaml:"mysql"`
	RedisCfg    *Redis    `yaml:"redis"`
	RegistryCfg *Registry `yaml:"registry"`
	MQCfg       *MQ       `yaml:"rabbitmq"`
}

type Mysql struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
	Db   string `yaml:"db"`
	User string `yaml:"user"`
	Pwd  string `yaml:"password"`
}

type Redis struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type Registry struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type MQ struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	PassWord string `yaml:"password"`
	VHost    string `yaml:"vhost"`
}
