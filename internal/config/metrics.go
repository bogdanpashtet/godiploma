package config

type MetricsConfig struct {
	// Port to start the server on.
	// Required.
	Port int `yaml:"port" env:"PORT"`

	// Endpoint is HTTP path for metrics endpoint.
	// Defaults to "/metrics".
	Endpoint string `yaml:"endpoint" env:"ENDPOINT"`
}
