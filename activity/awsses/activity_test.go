package awsses

import (
	"context"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/postqode/qodeflow-core/support/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockSESClient is a mock of the SESClient interface
type MockSESClient struct {
	mock.Mock
}

func (m *MockSESClient) SendEmail(ctx context.Context, params *sesv2.SendEmailInput, optFns ...func(*sesv2.Options)) (*sesv2.SendEmailOutput, error) {
	args := m.Called(ctx, params, optFns)
	return args.Get(0).(*sesv2.SendEmailOutput), args.Error(1)
}

func TestActivity_Eval(t *testing.T) {
	mockClient := new(MockSESClient)
	mockClient.On("SendEmail", mock.Anything, mock.Anything, mock.Anything).Return(&sesv2.SendEmailOutput{}, nil)

	act := &Activity{
		settings: &Settings{
			Region: "us-east-1",
		},
		client: mockClient,
	}
	tc := test.NewActivityContext(act.Metadata())

	input := &Input{
		From:    "sender@example.com",
		To:      []string{"test@example.com"},
		Subject: "Test Subject",
		Body:    "Test Body",
	}
	tc.SetInputObject(input)

	done, err := act.Eval(tc)
	assert.Nil(t, err)
	assert.True(t, done)
	mockClient.AssertExpectations(t)

	output := &Output{}
	err = tc.GetOutputObject(output)
	assert.Nil(t, err)
	assert.True(t, output.Success)
	assert.Empty(t, output.Error)
}

func TestActivity_Metadata(t *testing.T) {
	act := &Activity{}
	md := act.Metadata()
	assert.NotNil(t, md)
	assert.NotNil(t, md.Settings["region"])
	assert.NotNil(t, md.Input["to"])
}

func TestActivity_Eval_Integration(t *testing.T) {
	region := os.Getenv("AWS_REGION")
	accessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	secretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	from := os.Getenv("AWS_SES_FROM")
	to := os.Getenv("AWS_SES_TO")

	if region == "" || accessKey == "" || secretKey == "" || from == "" || to == "" {
		t.Skip("Skipping integration test: AWS credentials or test emails not set")
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
	input := &Input{
		From:    from,
		To:      []string{to},
		Subject: "Integration Test",
		Body:    "<h1>Integration Test Success</h1>",
	}
	tc.SetInputObject(input)

	done, err := act.Eval(tc)
	assert.NoError(t, err)
	assert.True(t, done)

	output := &Output{}
	err = tc.GetOutputObject(output)
	assert.NoError(t, err)
	assert.True(t, output.Success)
}
