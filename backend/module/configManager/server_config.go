package configManager

type ServerConfig struct {
	NodeType string `yaml:"node_type"`
	NodeId   int    `yaml:"node_id"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Scheme   string `yaml:"scheme"`
}
