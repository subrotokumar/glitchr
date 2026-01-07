package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gitlab.com/subrotokumar/glitchr/pkg/queue"
)

func main() {
	log := slog.New(slog.NewTextHandler(os.Stdout, nil))
	log.Info("Raw uploaded video consumer started")

	q := queue.NewMessageQueue("ap-south-1", log)
	queueUrl := "https://sqs.ap-south-1.amazonaws.com/123456789012/your-queue"

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			log.Info("Shutting down consumer gracefully")
			return
		default:
			// Fetch messages
			messages, err := q.GetMessages(ctx, queueUrl, 10, 10)
			if err != nil {
				log.Error("Failed to get messages", "error", err)
				time.Sleep(2 * time.Second)
				continue
			}

			if len(messages) == 0 {
				time.Sleep(1 * time.Second)
				continue
			}

			for _, msg := range messages {
				log.Info("Processing message", "id", *msg.MessageId, "body", *msg.Body)

				// TODO: Add actual processing here
				// e.g., move file from raw → processed bucket

				// Delete message after successful processing
				err := q.DeleteMessage(ctx, queueUrl, *msg.ReceiptHandle)
				if err != nil {
					log.Error("Failed to delete message", "id", *msg.MessageId, "error", err)
				}
			}
		}
	}
}
