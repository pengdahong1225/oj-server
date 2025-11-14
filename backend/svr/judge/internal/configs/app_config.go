package configs

type AppConfig struct {
	SandBoxCfg  []*SandBox `yaml:"sandbox"`
	RegistryCfg *Registry  `yaml:"registry"`
	MQCfg       *MQ        `yaml:"rabbitmq"`
}

type SandBox struct {
	Type string `yaml:"type"`
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
