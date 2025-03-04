package config

import (
	"context"
	"log"
	"os"
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
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
	MaxProcessedEvents   int
	MinCriticalityEvents int
	Env                  string
}

// LoadConfig reads configuration from environment variables
func LoadConfig(ctx context.Context) *Config {
	cfg := &Config{
		NATSURL:              getEnv("NATS_URL", "nats://localhost:4222"),
		NATSUsername:         getSecret(ctx, "NATS_USERNAME", "username"),
		NATSPassword:         getSecret(ctx, "NATS_PASSWORD", "password"),
		InfluxDBURL:          getEnv("INFLUXDB_URL", "http://localhost:8086"),
		InfluxDBToken:        getSecret(ctx, "INFLUXDB_TOKEN", "0kDO0SGFORhP0YbJLFNh8WRZ4T-iuY7uVR279NDUtHscRX8rJct1QTuAxMeYl3Rp_Kvx-4oZEYZDsuHjMNILeQ=="),
		InfluxDBOrg:          getEnv("INFLUXDB_ORG", "event-org"),
		InfluxDBBucket:       getEnv("INFLUXDB_BUCKET", "event-bucket"),
		LogLevel:             getEnv("LOG_LEVEL", "info"),
		WorkerCount:          getEnvAsInt("WORKER_COUNT", 5),
		MaxProcessedEvents:   getEnvAsInt("MAX_PROCESSED_EVENTS", 100),
		MinCriticalityEvents: getEnvAsInt("MIN_CRITICALITY_EVENTS", 5),
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

func getSecret(ctx context.Context, key string, defaultValue string) string {
	env := getEnv("ENV", "dev")

	if env == "dev" {
		return getEnv(key, defaultValue)
	}

	return getSecretFromAWS(ctx, key)
}

func getSecretFromAWS(ctx context.Context, secretName string) string {
	// Load AWS SDK configuration
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion("us-east-1")) // Change region accordingly
	if err != nil {
		log.Printf("Error loading AWS config: %v", err)
		return ""
	}

	// Create AWS Secrets Manager client
	client := secretsmanager.NewFromConfig(cfg)

	// Get secret value
	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretName),
	}

	result, err := client.GetSecretValue(ctx, input)
	if err != nil {
		log.Printf("Error retrieving secret %s from AWS Secrets Manager: %v", secretName, err)
		return ""
	}

	log.Println("Successfully retrieved secret from AWS Secrets Manager")
	return aws.ToString(result.SecretString)
}
