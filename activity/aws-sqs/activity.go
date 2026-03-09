package awssqs

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/postqode/qodeflow-core/activity"
	"github.com/postqode/qodeflow-core/data/metadata"
)

func init() {
	_ = activity.Register(&Activity{}, New)
}

// Activity represents the AWS SQS activity
type Activity struct {
	settings *Settings
	client   SQSClient
}

// SQSClient interface for mocking in tests
type SQSClient interface {
	SendMessage(ctx context.Context, params *sqs.SendMessageInput, optFns ...func(*sqs.Options)) (*sqs.SendMessageOutput, error)
	ReceiveMessage(ctx context.Context, params *sqs.ReceiveMessageInput, optFns ...func(*sqs.Options)) (*sqs.ReceiveMessageOutput, error)
	DeleteMessage(ctx context.Context, params *sqs.DeleteMessageInput, optFns ...func(*sqs.Options)) (*sqs.DeleteMessageOutput, error)
	ListQueues(ctx context.Context, params *sqs.ListQueuesInput, optFns ...func(*sqs.Options)) (*sqs.ListQueuesOutput, error)
}

// New creates a new AWS SQS activity
func New(ctx activity.InitContext) (activity.Activity, error) {
	s := &Settings{}
	err := metadata.MapToStruct(ctx.Settings(), s, true)
	if err != nil {
		return nil, err
	}

	if s.Region == "" {
		return nil, fmt.Errorf("region is required for AWS SQS activity")
	}

	var cfg aws.Config
	if s.AccessKeyID != "" && s.SecretAccessKey != "" {
		cfg, err = config.LoadDefaultConfig(context.TODO(),
			config.WithRegion(s.Region),
			config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(s.AccessKeyID, s.SecretAccessKey, s.SessionToken)),
		)
	} else {
		cfg, err = config.LoadDefaultConfig(context.TODO(),
			config.WithRegion(s.Region),
		)
	}

	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config, %v", err)
	}

	client := sqs.NewFromConfig(cfg)

	return &Activity{
		settings: s,
		client:   client,
	}, nil
}

// Metadata returns the metadata for the AWS SQS activity
func (act *Activity) Metadata() *activity.Metadata {
	return activityMetadata
}

// Eval executes the AWS SQS logic
func (act *Activity) Eval(ctx activity.Context) (done bool, err error) {
	input := &Input{}
	err = ctx.GetInputObject(input)
	if err != nil {
		return true, err
	}

	output := &Output{
		Success: true,
	}

	switch strings.ToLower(input.Method) {
	case "sendmessage":
		err = act.sendMessage(ctx, input, output)
	case "receivemessage":
		err = act.receiveMessage(ctx, input, output)
	case "deletemessage":
		err = act.deleteMessage(ctx, input, output)
	case "listqueues":
		err = act.listQueues(ctx, input, output)
	default:
		err = fmt.Errorf("unsupported method: %s", input.Method)
	}

	if err != nil {
		output.Success = false
		output.Error = err.Error()
		ctx.Logger().Errorf("AWS SQS operation %s failed: %v", input.Method, err)
	}

	err = ctx.SetOutputObject(output)
	if err != nil {
		return true, err
	}

	return true, nil
}

func (act *Activity) sendMessage(ctx activity.Context, input *Input, output *Output) error {
	if input.QueueURL == "" {
		return fmt.Errorf("queueURL is required for SendMessage")
	}
	if input.MessageBody == "" {
		return fmt.Errorf("messageBody is required for SendMessage")
	}

	resp, err := act.client.SendMessage(context.TODO(), &sqs.SendMessageInput{
		QueueUrl:    aws.String(input.QueueURL),
		MessageBody: aws.String(input.MessageBody),
	})
	if err != nil {
		return err
	}

	output.Result = map[string]any{
		"messageId": *resp.MessageId,
	}
	return nil
}

func (act *Activity) receiveMessage(ctx activity.Context, input *Input, output *Output) error {
	if input.QueueURL == "" {
		return fmt.Errorf("queueURL is required for ReceiveMessage")
	}

	maxMessages := int32(input.MaxNumberOfMessages)
	if maxMessages <= 0 {
		maxMessages = 1
	}
	if maxMessages > 10 {
		maxMessages = 10
	}

	waitTime := int32(input.WaitTimeSeconds)

	resp, err := act.client.ReceiveMessage(context.TODO(), &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(input.QueueURL),
		MaxNumberOfMessages: maxMessages,
		WaitTimeSeconds:     waitTime,
	})
	if err != nil {
		return err
	}

	var messages []map[string]any
	for _, msg := range resp.Messages {
		m := map[string]any{
			"messageId":     aws.ToString(msg.MessageId),
			"body":          aws.ToString(msg.Body),
			"receiptHandle": aws.ToString(msg.ReceiptHandle),
		}
		messages = append(messages, m)
	}

	output.Result = messages
	return nil
}

func (act *Activity) deleteMessage(ctx activity.Context, input *Input, output *Output) error {
	if input.QueueURL == "" {
		return fmt.Errorf("queueURL is required for DeleteMessage")
	}
	if input.ReceiptHandle == "" {
		return fmt.Errorf("receiptHandle is required for DeleteMessage")
	}

	_, err := act.client.DeleteMessage(context.TODO(), &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(input.QueueURL),
		ReceiptHandle: aws.String(input.ReceiptHandle),
	})
	return err
}

func (act *Activity) listQueues(ctx activity.Context, input *Input, output *Output) error {
	resp, err := act.client.ListQueues(context.TODO(), &sqs.ListQueuesInput{})
	if err != nil {
		return err
	}

	output.Result = resp.QueueUrls
	return nil
}
