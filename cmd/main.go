package main

import (
	"context"
	"email-service/internal/config"
	"email-service/internal/sqs"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	consumer, err := sqs.NewConsumer(cfg.AWS.SQSQueueURL)
	if err != nil {
		log.Fatalf("Failed to create SQS consumer: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		cancel()
	}()

	// Start consuming messages
	if err := consumer.Start(ctx); err != nil {
		log.Fatalf("Consumer error: %v", err)
	}
}
