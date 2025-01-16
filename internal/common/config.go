package common

import "os"

type AsanaConfig struct {
	AccessToken string
	WorkspaceId string
}

type Config struct {
	Asana AsanaConfig
}

// Create new config with environment variables
func NewConfig() *Config {
	return &Config{
		Asana: AsanaConfig{
			AccessToken: getEnv("ASANA_ACCESS_TOKEN", ""),
			WorkspaceId: getEnv("ASANA_WORKSPACE_ID", ""),
		},
	}
}

// Get env value as string
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
