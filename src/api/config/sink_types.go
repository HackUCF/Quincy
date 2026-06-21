package config

type Sinks struct {
	PGConfig PGConfig `yaml:"postgres" mapstructure:"postgres" json:"postgres"`
}

type PGConfig struct {
	Host     string `yaml:"host"     mapstructure:"host"     json:"host"     example:"localhost"`
	Port     uint16 `yaml:"port"     mapstructure:"port"     json:"port"     example:"5432"`
	Username string `yaml:"username" mapstructure:"username" json:"username" example:"postgres"`
	Password string `yaml:"password" mapstructure:"password" json:"password" example:"postgres"`
	Database string `yaml:"database" mapstructure:"database" json:"database" example:"quincy"`
	SSLMode  string `yaml:"ssl_mode" mapstructure:"ssl_mode" json:"ssl_mode" example:"prefer"`
}
