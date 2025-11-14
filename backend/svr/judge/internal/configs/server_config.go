package configs

type ServerConfig struct {
	Name    string `yaml:"name"`
	NodeId  int    `yaml:"node_id"`
	Address string `yaml:"address"`
	Port    int    `yaml:"port"`
	Scheme  string `yaml:"scheme"` // http/grpc
}
