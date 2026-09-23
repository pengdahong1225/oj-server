package config

type Config struct {
	Service  *Service  `yaml:"service"`
	Mysql    *Mysql    `yaml:"mysql"`
	Redis    *Redis    `yaml:"redis"`
	Registry *Registry `yaml:"registry"`
	Sms      *Sms      `yaml:"sms"`
}

type Service struct {
	Name    string `yaml:"name"`
	NodeId  int    `yaml:"node_id"`
	Address string `yaml:"address"`
	Port    int    `yaml:"port"`
	Scheme  string `yaml:"scheme"` // http/grpc
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

type Sms struct {
	AccessKeyId     string `yaml:"accessKeyId"`
	AccessKeySecret string `yaml:"accessKeySecret"`
	Endpoint        string `yaml:"endpoint"`
	SignName        string `yaml:"signName"`
	TemplateCode    string `yaml:"templateCode"`
}
