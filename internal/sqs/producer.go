package sqs

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type Publisher struct {
	client   *sqs.Client
	queueURL string
	isFIFO   bool
}

func New(queueURL string) (*Publisher, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Printf("failed to load AWS config: %v", err)
		return nil, err
	}

	client := sqs.NewFromConfig(cfg)

	// Check if queue is FIFO by looking at the URL suffix
	isFIFO := false
	if len(queueURL) > 5 && queueURL[len(queueURL)-5:] == ".fifo" {
		isFIFO = true
	}

	return &Publisher{
		client:   client,
		queueURL: queueURL,
		isFIFO:   isFIFO,
	}, nil
}

func (p *Publisher) Publish(ctx context.Context, message interface{}) error {
	messageBody, err := json.Marshal(message)
	if err != nil {
		log.Printf("failed to marshal message: %v", err)
		return err
	}

	input := &sqs.SendMessageInput{
		QueueUrl:    &p.queueURL,
		MessageBody: String(string(messageBody)),
	}

	if p.isFIFO {
		// For FIFO queues, MessageGroupId is required
		input.MessageGroupId = String("default")

		// Generate a message deduplication ID based on the message content
		hash := fmt.Sprintf("%x", sha256.Sum256(messageBody))
		input.MessageDeduplicationId = String(hash)
	}

	_, err = p.client.SendMessage(ctx, input)
	if err != nil {
		log.Printf("failed to send message to SQS: %v", err)
		return err
	}

	return nil
}

// String returns a pointer to the string value passed in
func String(v string) *string {
	return &v
}
