package queue

import (
	"context"
	"log/slog"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type Queue struct {
	SqsClient *sqs.Client
	log       *slog.Logger
}

func NewMessageQueue(region string, log *slog.Logger) *Queue {
	ctx := context.Background()
	sdkConfig, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		log.Info("Couldn't load default configuration. Have you set up your AWS account?")
		log.Error(err.Error())
	}
	sqsClient := sqs.NewFromConfig(sdkConfig)
	return &Queue{
		SqsClient: sqsClient,
	}
}
