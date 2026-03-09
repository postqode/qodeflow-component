package awssqs

import (
	"context"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/postqode/qodeflow-core/support/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockSQSClient struct {
	mock.Mock
}

func (m *MockSQSClient) SendMessage(ctx context.Context, params *sqs.SendMessageInput, optFns ...func(*sqs.Options)) (*sqs.SendMessageOutput, error) {
	args := m.Called(ctx, params, optFns)
	return args.Get(0).(*sqs.SendMessageOutput), args.Error(1)
}

func (m *MockSQSClient) ReceiveMessage(ctx context.Context, params *sqs.ReceiveMessageInput, optFns ...func(*sqs.Options)) (*sqs.ReceiveMessageOutput, error) {
	args := m.Called(ctx, params, optFns)
	return args.Get(0).(*sqs.ReceiveMessageOutput), args.Error(1)
}

func (m *MockSQSClient) DeleteMessage(ctx context.Context, params *sqs.DeleteMessageInput, optFns ...func(*sqs.Options)) (*sqs.DeleteMessageOutput, error) {
	args := m.Called(ctx, params, optFns)
	return args.Get(0).(*sqs.DeleteMessageOutput), args.Error(1)
}

func (m *MockSQSClient) ListQueues(ctx context.Context, params *sqs.ListQueuesInput, optFns ...func(*sqs.Options)) (*sqs.ListQueuesOutput, error) {
	args := m.Called(ctx, params, optFns)
	return args.Get(0).(*sqs.ListQueuesOutput), args.Error(1)
}

func TestSQSActivity_SendMessage(t *testing.T) {
	mockClient := new(MockSQSClient)
	act := &Activity{
		settings: &Settings{Region: "us-east-1"},
		client:   mockClient,
	}

	queueURL := "https://sqs.us-east-1.amazonaws.com/123456789012/MyQueue"
	messageBody := "Hello World"
	messageID := "msg-123"

	mockClient.On("SendMessage", mock.Anything, mock.MatchedBy(func(input *sqs.SendMessageInput) bool {
		return *input.QueueUrl == queueURL && *input.MessageBody == messageBody
	}), mock.Anything).Return(&sqs.SendMessageOutput{MessageId: &messageID}, nil)

	tc := test.NewActivityContext(activityMetadata)
	input := &Input{
		Method:      "SendMessage",
		QueueURL:    queueURL,
		MessageBody: messageBody,
	}
	tc.SetInputObject(input)

	done, err := act.Eval(tc)
	assert.True(t, done)
	assert.Nil(t, err)

	output := &Output{}
	tc.GetOutputObject(output)
	assert.True(t, output.Success)
	assert.Equal(t, messageID, output.Result.(map[string]any)["messageId"])
}

func TestSQSActivity_ReceiveMessage(t *testing.T) {
	mockClient := new(MockSQSClient)
	act := &Activity{
		settings: &Settings{Region: "us-east-1"},
		client:   mockClient,
	}

	queueURL := "https://sqs.us-east-1.amazonaws.com/123456789012/MyQueue"
	messageID := "msg-123"
	body := "Hello World"
	receiptHandle := "receipt-123"

	mockClient.On("ReceiveMessage", mock.Anything, mock.MatchedBy(func(input *sqs.ReceiveMessageInput) bool {
		return *input.QueueUrl == queueURL
	}), mock.Anything).Return(&sqs.ReceiveMessageOutput{
		Messages: []types.Message{
			{
				MessageId:     &messageID,
				Body:          &body,
				ReceiptHandle: &receiptHandle,
			},
		},
	}, nil)

	tc := test.NewActivityContext(activityMetadata)
	input := &Input{
		Method:   "ReceiveMessage",
		QueueURL: queueURL,
	}
	tc.SetInputObject(input)

	done, err := act.Eval(tc)
	assert.True(t, done)
	assert.Nil(t, err)

	output := &Output{}
	tc.GetOutputObject(output)
	assert.True(t, output.Success)
	results := output.Result.([]map[string]any)
	assert.Len(t, results, 1)
	assert.Equal(t, messageID, results[0]["messageId"])
}

func TestSQSActivity_DeleteMessage(t *testing.T) {
	mockClient := new(MockSQSClient)
	act := &Activity{
		settings: &Settings{Region: "us-east-1"},
		client:   mockClient,
	}

	queueURL := "https://sqs.us-east-1.amazonaws.com/123456789012/MyQueue"
	receiptHandle := "receipt-123"

	mockClient.On("DeleteMessage", mock.Anything, mock.MatchedBy(func(input *sqs.DeleteMessageInput) bool {
		return *input.QueueUrl == queueURL && *input.ReceiptHandle == receiptHandle
	}), mock.Anything).Return(&sqs.DeleteMessageOutput{}, nil)

	tc := test.NewActivityContext(activityMetadata)
	input := &Input{
		Method:        "DeleteMessage",
		QueueURL:      queueURL,
		ReceiptHandle: receiptHandle,
	}
	tc.SetInputObject(input)

	done, err := act.Eval(tc)
	assert.True(t, done)
	assert.Nil(t, err)

	output := &Output{}
	tc.GetOutputObject(output)
	assert.True(t, output.Success)
}

func TestSQSActivity_Integration(t *testing.T) {
	region := os.Getenv("AWS_REGION")
	accessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	secretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	queueURL := os.Getenv("AWS_SQS_QUEUE_URL")

	if region == "" || accessKey == "" || secretKey == "" || queueURL == "" {
		t.Skip("Skipping integration test: AWS credentials or SQS Queue URL not set")
	}

	settings := map[string]any{
		"region":          region,
		"accessKeyID":     accessKey,
		"secretAccessKey": secretKey,
	}

	iCtx := test.NewActivityInitContext(settings, nil)
	act, err := New(iCtx)
	assert.NoError(t, err)

	tc := test.NewActivityContext(act.Metadata())

	// Test SendMessage
	input := &Input{
		Method:      "SendMessage",
		QueueURL:    queueURL,
		MessageBody: "Integration test message",
	}
	tc.SetInputObject(input)

	done, err := act.Eval(tc)
	assert.NoError(t, err)
	assert.True(t, done)

	output := &Output{}
	tc.GetOutputObject(output)
	assert.True(t, output.Success)
	assert.NotNil(t, output.Result)
	messageID := output.Result.(map[string]any)["messageId"].(string)

	// Test ReceiveMessage
	input = &Input{
		Method:   "ReceiveMessage",
		QueueURL: queueURL,
	}
	tc.SetInputObject(input)
	done, err = act.Eval(tc)
	assert.NoError(t, err)

	tc.GetOutputObject(output)
	assert.True(t, output.Success)
	messages := output.Result.([]map[string]any)
	assert.NotEmpty(t, messages)

	// Find our message and delete it
	var receiptHandle string
	for _, msg := range messages {
		if msg["messageId"] == messageID {
			receiptHandle = msg["receiptHandle"].(string)
			break
		}
	}

	if receiptHandle != "" {
		input = &Input{
			Method:        "DeleteMessage",
			QueueURL:      queueURL,
			ReceiptHandle: receiptHandle,
		}
		tc.SetInputObject(input)
		done, err = act.Eval(tc)
		assert.NoError(t, err)

		tc.GetOutputObject(output)
		assert.True(t, output.Success)
	}
}
