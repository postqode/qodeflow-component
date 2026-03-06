package awsses

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"
	"github.com/postqode/qodeflow-core/activity"
	"github.com/postqode/qodeflow-core/data/metadata"
)

func init() {
	_ = activity.Register(&Activity{}, New)
}

// Activity represents the AWS SES activity
type Activity struct {
	settings *Settings
	client   SESClient
}

// SESClient interface for mocking in tests
type SESClient interface {
	SendEmail(ctx context.Context, params *sesv2.SendEmailInput, optFns ...func(*sesv2.Options)) (*sesv2.SendEmailOutput, error)
}

// New creates a new AWS SES activity
func New(ctx activity.InitContext) (activity.Activity, error) {
	s := &Settings{}
	err := metadata.MapToStruct(ctx.Settings(), s, true)
	if err != nil {
		return nil, err
	}

	if s.Region == "" {
		return nil, fmt.Errorf("region is required for AWS SES activity")
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

	client := sesv2.NewFromConfig(cfg)

	return &Activity{
		settings: s,
		client:   client,
	}, nil
}

// Metadata returns the metadata for the AWS SES activity
func (act *Activity) Metadata() *activity.Metadata {
	return activityMetadata
}

// Eval executes the AWS SES sending logic
func (act *Activity) Eval(ctx activity.Context) (done bool, err error) {
	input := &Input{}
	err = ctx.GetInputObject(input)
	if err != nil {
		return true, err
	}

	// For now, we'll implement a simple email sending.
	// Attachments (files) require RawMessage which is more complex.
	// Implementing simple message first.

	params := &sesv2.SendEmailInput{
		FromEmailAddress: aws.String(input.From),
		Destination: &types.Destination{
			ToAddresses: input.To,
		},
		Content: &types.EmailContent{
			Simple: &types.Message{
				Subject: &types.Content{
					Data: aws.String(input.Subject),
				},
				Body: &types.Body{
					Html: &types.Content{
						Data: aws.String(input.Body),
					},
				},
			},
		},
	}

	_, err = act.client.SendEmail(context.TODO(), params)

	output := &Output{
		Success: err == nil,
	}

	if err != nil {
		output.Error = err.Error()
		ctx.Logger().Errorf("Failed to send AWS SES email: %v", err)
	}

	err = ctx.SetOutputObject(output)
	if err != nil {
		return true, err
	}

	return true, nil
}
