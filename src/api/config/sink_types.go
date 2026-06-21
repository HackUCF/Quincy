package config

type Sinks struct {
	PGConfig   PGConfig   `yaml:"postgres" mapstructure:"postgres" json:"postgres"`
	OTelConfig OTelConfig `yaml:"otel"     mapstructure:"otel"     json:"otel"`
}

type OTelConfig struct {
	Endpoint  string `yaml:"endpoint"    mapstructure:"endpoint"    json:"endpoint"    example:"http://localhost:4318"`
	Username  string `yaml:"username"    mapstructure:"username"    json:"username"    example:"admin"`
	Password  string `yaml:"password"    mapstructure:"password"    json:"password"    example:"secret"`
	BasicAuth string `yaml:"basic_auth"  mapstructure:"basic_auth"  json:"basic_auth"  example:"dXNlcjpwYXNz"`

	// OpenObserve specific header
	StreamName string `yaml:"stream_name" mapstructure:"stream_name" json:"stream_name" example:"quincy"`

	// Batching settings
	Batching OTelBatching `yaml:"batching"    mapstructure:"batching"    json:"batching"`
}

// OTelBatching controls how log records are batched before export.
// Zero values use defaults (BatchSize: 20, ExportInterval: 5s, MaxQueueSize: 200).
type OTelBatching struct {
	BatchSize      int `yaml:"batch_size"      mapstructure:"batch_size"      json:"batch_size"      example:"20"`
	ExportInterval int `yaml:"export_interval" mapstructure:"export_interval" json:"export_interval" example:"5"`
	MaxQueueSize   int `yaml:"max_queue_size"  mapstructure:"max_queue_size"  json:"max_queue_size"  example:"200"`
}

type PGConfig struct {
	Host     string `yaml:"host"     mapstructure:"host"     json:"host"     example:"localhost"`
	Port     uint16 `yaml:"port"     mapstructure:"port"     json:"port"     example:"5432"`
	Username string `yaml:"username" mapstructure:"username" json:"username" example:"postgres"`
	Password string `yaml:"password" mapstructure:"password" json:"password" example:"postgres"`
	Database string `yaml:"database" mapstructure:"database" json:"database" example:"quincy"`
	SSLMode  string `yaml:"ssl_mode" mapstructure:"ssl_mode" json:"ssl_mode" example:"prefer"`
}
