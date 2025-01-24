package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AWS struct {
		Region          string
		AccessKeyID     string
		SecretAccessKey string
		SQSQueueURL     string
	}
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	cfg := &Config{}

	// AWS Configuration
	cfg.AWS.Region = os.Getenv("AWS_REGION")
	if cfg.AWS.Region == "" {
		cfg.AWS.Region = "us-east-1"
	}
	cfg.AWS.AccessKeyID = os.Getenv("AWS_ACCESS_KEY_ID")
	cfg.AWS.SecretAccessKey = os.Getenv("AWS_SECRET_ACCESS_KEY")
	cfg.AWS.SQSQueueURL = os.Getenv("SQS_QUEUE_URL")

	return cfg, nil
}
