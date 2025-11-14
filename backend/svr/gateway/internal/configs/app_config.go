package configs

type AppConfig struct {
	RegistryCfg *Registry `yaml:"registry"`
	JwtCfg      *Jwt      `yaml:"jwt"`
}

type Jwt struct {
	SigningKey string `yaml:"key"`
}

type Registry struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}
