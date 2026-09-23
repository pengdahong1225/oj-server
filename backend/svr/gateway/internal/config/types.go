package config

type Config struct {
	Service  *Service  `yaml:"service"`
	Registry *Registry `yaml:"registry"`
	Jwt      *Jwt      `yaml:"jwt"`
}

type Service struct {
	Name    string `yaml:"name"`
	NodeId  int    `yaml:"node_id"`
	Address string `yaml:"address"`
	Port    int    `yaml:"port"`
	Scheme  string `yaml:"scheme"` // http/grpc
}

type Jwt struct {
	SigningKey string `yaml:"key"`
}

type Registry struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}
