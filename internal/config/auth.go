package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type AuthConfig struct {
	Keys map[string]string `yaml:"keys"`
}

func initAuthConfig() (AuthConfig, error) {
	hmacKeysJSON := os.Getenv("AUTH_USER_KEYS")
	if hmacKeysJSON == "" {
		return AuthConfig{}, fmt.Errorf("HMAC_USER_KEYS environment variable not set")
	}

	var parsedKeys map[string]string
	if err := json.Unmarshal([]byte(hmacKeysJSON), &parsedKeys); err != nil {
		return AuthConfig{}, fmt.Errorf("failed to parse JSON from HMAC_USER_KEYS environment variable: %w", err)
	}

	return AuthConfig{Keys: parsedKeys}, nil
}
