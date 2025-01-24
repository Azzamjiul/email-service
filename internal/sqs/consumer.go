package sqs

import (
	"context"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type Consumer struct {
	client   *sqs.Client
	queueURL string
}

func NewConsumer(queueURL string) (*Consumer, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Printf("failed to load AWS config: %v", err)
		return nil, err
	}

	client := sqs.NewFromConfig(cfg)

	return &Consumer{
		client:   client,
		queueURL: queueURL,
	}, nil
}

func (c *Consumer) Start(ctx context.Context) error {
	log.Printf("Starting SQS consumer for queue: %s", c.queueURL)

	for {
		select {
		case <-ctx.Done():
			log.Println("Shutting down SQS consumer...")
			return nil
		default:
			messages, err := c.receiveMessages(ctx)
			if err != nil {
				log.Printf("Couldn't get messages from queue %v. Here's why: %v\n", c.queueURL, err)
				time.Sleep(time.Second) // Wait before retry
				continue
			}

			for _, msg := range messages {
				if msg.Body != nil {
					log.Printf("Processing message: %s", *msg.Body)

					// Delete the message after processing
					if err := c.deleteMessage(ctx, msg.ReceiptHandle); err != nil {
						log.Printf("Failed to delete message: %v", err)
						continue
					}
					log.Printf("Message deleted successfully")
				}
			}
		}
	}
}

func (c *Consumer) receiveMessages(ctx context.Context) ([]types.Message, error) {
	output, err := c.client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(c.queueURL),
		MaxNumberOfMessages: 10,
		WaitTimeSeconds:     20,
	})
	if err != nil {
		return nil, err
	}

	return output.Messages, nil
}

func (c *Consumer) deleteMessage(ctx context.Context, receiptHandle *string) error {
	_, err := c.client.DeleteMessage(ctx, &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(c.queueURL),
		ReceiptHandle: receiptHandle,
	})
	return err
}
