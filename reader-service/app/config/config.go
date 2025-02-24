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
	InfluxDBURL          string
	InfluxDBToken        string
	InfluxDBOrg          string
	InfluxDBBucket       string
	LogLevel             string
	WorkerCount          int
	TaskQueueCap         int // Buffered task queue
	TaskQueueSize        int
	MaxProcessedEvents   int
	MinCriticalityEvents int
	PrometheusPort       string
	Env                  string
	ServerAddress        string
}

// LoadConfig reads configuration from environment variables
func LoadConfig() *Config {
	cfg := &Config{
		NATSURL:              getEnv("NATS_URL", "nats://localhost:4222"),
		NATSUsername:         getEnv("NATS_USERNAME", "username"),
		NATSPassword:         getEnv("NATS_PASSWORD", "password"),
		InfluxDBURL:          getEnv("INFLUXDB_URL", "http://localhost:8086"),
		InfluxDBToken:        getEnv("INFLUXDB_TOKEN", "0kDO0SGFORhP0YbJLFNh8WRZ4T-iuY7uVR279NDUtHscRX8rJct1QTuAxMeYl3Rp_Kvx-4oZEYZDsuHjMNILeQ=="),
		InfluxDBOrg:          getEnv("INFLUXDB_ORG", "event-org"),
		InfluxDBBucket:       getEnv("INFLUXDB_BUCKET", "event-bucket"),
		LogLevel:             getEnv("LOG_LEVEL", "info"),
		WorkerCount:          getEnvAsInt("WORKER_COUNT", 5),
		TaskQueueCap:         getEnvAsInt("TASK_QUEUE_CAP", 100),
		TaskQueueSize:        getEnvAsInt("TASK_QUEUE_SIZE", 100),
		MaxProcessedEvents:   getEnvAsInt("MAX_PROCESSED_EVENTS", 100),
		MinCriticalityEvents: getEnvAsInt("MIN_CRITICALITY_EVENTS", 5),
		PrometheusPort:       getEnv("PROMETHEUS_PORT", ":9090"),
		Env:                  getEnv("ENV", "dev"),
		ServerAddress:        getEnv("SERVER_ADDRESS", "localhost:8080"),
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
