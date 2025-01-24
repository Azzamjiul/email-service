package main

import (
	"email-service/internal/mailer"
	"email-service/internal/sqs"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type EmailRequest struct {
	Recipient    string      `json:"recipient"`
	TemplateFile string      `json:"template_file"`
	Data         interface{} `json:"data"`
}

type SQSMessageRequest struct {
	Message interface{} `json:"message"`
}

var (
	mailerInstance *mailer.Mailer
	sqsPublisher   *sqs.Publisher
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		return
	}

	host := os.Getenv("HOST")
	sender := os.Getenv("MAILER_SENDER")
	username := os.Getenv("MAILER_USERNAME")
	password := os.Getenv("MAILER_PASSWORD")
	mailer_host := os.Getenv("MAILER_HOST")

	mailerInstance = mailer.New(sender, username, password, mailer_host)

	// Initialize SQS publisher
	queueURL := os.Getenv("SQS_QUEUE_URL")
	publisher, err := sqs.New(queueURL)
	if err != nil {
		fmt.Println("Error initializing SQS publisher:", err)
		return
	}
	sqsPublisher = publisher

	http.HandleFunc("/", handler)
	http.HandleFunc("/send-email", mailerHandler)
	http.HandleFunc("/publish-message", sqsHandler)

	fmt.Printf("Starting server on http://%s:8080\n", host)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!")
}

func mailerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var emailReq EmailRequest
	err := json.NewDecoder(r.Body).Decode(&emailReq)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = mailerInstance.Send(emailReq.Recipient, emailReq.TemplateFile, emailReq.Data)
	if err != nil {
		http.Error(w, "Failed to send email", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Mailer called successfully!")
}

func sqsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var messageReq SQSMessageRequest
	err := json.NewDecoder(r.Body).Decode(&messageReq)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = sqsPublisher.Publish(r.Context(), messageReq.Message)
	if err != nil {
		http.Error(w, "Failed to publish message", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Message published successfully!")
}
