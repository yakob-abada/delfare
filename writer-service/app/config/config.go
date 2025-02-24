package config

import (
	"os"
	"strconv"
)

// Config holds service configurations
type Config struct {
	Env            string
	NATSURL        string
	NATSUsername   string
	NATSPassword   string
	InfluxDBURL    string
	InfluxDBOrg    string
	InfluxDBBucket string
	InfluxDBToken  string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	return &Config{
		Env:            getEnv("ENV", "dev"),
		NATSURL:        getEnv("NATS_URL", "nats://localhost:4222"),
		NATSUsername:   getEnv("NATS_USERNAME", "username"),
		NATSPassword:   getEnv("NATS_PASSWORD", "password"),
		InfluxDBURL:    getEnv("INFLUXDB_URL", "http://localhost:8086"),
		InfluxDBToken:  getEnv("INFLUXDB_TOKEN", "token"),
		InfluxDBOrg:    getEnv("INFLUXDB_ORG", "event-org"),
		InfluxDBBucket: getEnv("INFLUXDB_BUCKET", "event-bucket"),
	}
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
