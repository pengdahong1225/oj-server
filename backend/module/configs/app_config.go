package configs

type AppConfig struct {
	SandBoxCfg  []*SandBox `yaml:"sandbox"`
	MysqlCfg    *Mysql     `yaml:"mysql"`
	RedisCfg    *Redis     `yaml:"redis"`
	RegistryCfg *Registry  `yaml:"registry"`
	MQCfg       *MQ        `yaml:"rabbitmq"`
	JwtCfg      *Jwt       `yaml:"jwt"`
	SmsCfg      *Sms       `yaml:"sms"`
}

type SandBox struct {
	Type string `yaml:"type"`
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
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

type Jwt struct {
	SigningKey string `yaml:"key"`
}

type Sms struct {
	AccessKeyId     string `yaml:"accessKeyId"`
	AccessKeySecret string `yaml:"accessKeySecret"`
	Endpoint        string `yaml:"endpoint"`
	SignName        string `yaml:"signName"`
	TemplateCode    string `yaml:"templateCode"`
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
