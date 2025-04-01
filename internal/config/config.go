package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

const defaultConfigPath = "/app/config/config.yaml"

type AppConfig struct {
	Env        Env              `yaml:"env"`
	AppName    string           `yaml:"appName"`
	HTTPHealth HTTPHealthConfig `yaml:"httpHealth"`
	Metrics    MetricsConfig    `yaml:"metrics"`
	Grpc       GRPCBase         `yaml:"grpc"`
}

func (c *AppConfig) Validate() error {
	return nil
}

func New() (*AppConfig, error) {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = defaultConfigPath
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		log.Fatalf("fail while reading config file: %v", err)
	}

	var cfg AppConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		log.Fatalf("yaml parcing error: %v", err)
	}

	return &cfg, nil
}
