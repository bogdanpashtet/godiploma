package config

type HTTPHealthConfig struct {
	// Port to start the server on.
	// Required.
	Port int `yaml:"port" env:"PORT"`

	// ReadyEndpoint is HTTP path for ready endpoint.
	// If empty, ready endpoint is disabled.
	ReadyEndpoint string `yaml:"readyEndpoint" env:"READY_ENDPOINT"`

	// LiveEndpoint is HTTP path for live endpoint.
	// If empty, live endpoint is disabled.
	LiveEndpoint string `yaml:"liveEndpoint" env:"LIVE_ENDPOINT"`
}
