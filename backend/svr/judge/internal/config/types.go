package config

type Config struct {
	Service  *Service   `yaml:"service"`
	Registry *Registry  `yaml:"registry"`
	MQ       *MQ        `yaml:"rabbitmq"`
	SandBox  []*SandBox `yaml:"sandbox"`
}

type Service struct {
	Name    string `yaml:"name"`
	NodeId  int    `yaml:"node_id"`
	Address string `yaml:"address"`
	Port    int    `yaml:"port"`
	Scheme  string `yaml:"scheme"` // http/grpc
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
