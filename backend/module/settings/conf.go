package settings

// 从viper反射到数据模型，需要设置`mapstructure`反射字段

type SandBox struct {
	Type string `mapstructure:"type"`
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type Mysql struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
	Db   string `mapstructure:"db"`
	User string `mapstructure:"user"`
	Pwd  string `mapstructure:"password"`
}

type Redis struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type Jwt struct {
	SigningKey string `mapstructure:"key"`
}

type Sms struct {
	AccessKeyId     string `mapstructure:"accessKeyId"`
	AccessKeySecret string `mapstructure:"accessKeySecret"`
	Endpoint        string `mapstructure:"endpoint"`
	SignName        string `mapstructure:"signName"`
	TemplateCode    string `mapstructure:"templateCode"`
}

type Registry struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

type MQ struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	PassWord string `mapstructure:"password"`
	VHost    string `mapstructure:"vhost"`
}
