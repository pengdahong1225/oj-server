package configs

type AppConfig struct {
	MysqlCfg    *Mysql    `yaml:"mysql"`
	RedisCfg    *Redis    `yaml:"redis"`
	RegistryCfg *Registry `yaml:"registry"`
	SmsCfg      *Sms      `yaml:"sms"`
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
