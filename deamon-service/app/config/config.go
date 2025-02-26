package config

import (
	"log"
	"os"
	"strconv"
)

// Config Struct
type Config struct {
	NATSURL              string
	NATSUsername         string
	NATSPassword         string
	LogLevel             string
	CriticalityThreshold int
	EventProcessLimit    int
	Env                  string
}

// LoadConfig reads configuration from environment variables
func LoadConfig() *Config {
	cfg := &Config{
		NATSURL:              getEnv("NATS_URL", "nats://localhost:4222"),
		NATSUsername:         getEnv("NATS_USERNAME", "username"),
		NATSPassword:         getEnv("NATS_PASSWORD", "password"),
		LogLevel:             getEnv("LOG_LEVEL", "info"),
		CriticalityThreshold: getEnvAsInt("CRITICALITY_THRESHOLD", 5),
		EventProcessLimit:    getEnvAsInt("EVENT_PROCESS_LIMIT", 5),
		Env:                  getEnv("ENV", "dev"),
	}

	log.Println("Config Loaded Successfully")
	return cfg
}

// Helper function to read environment variables with default values
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// Helper function to read integer environment variables
func getEnvAsInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
