package awss3

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/postqode/qodeflow-core/activity"
	"github.com/postqode/qodeflow-core/data/metadata"
)

func init() {
	_ = activity.Register(&Activity{}, New)
}

// Activity represents the AWS S3 activity
type Activity struct {
	settings *Settings
	client   S3Client
}

// S3Client interface for mocking in tests
type S3Client interface {
	PutObject(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error)
	GetObject(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error)
	ListObjectsV2(ctx context.Context, params *s3.ListObjectsV2Input, optFns ...func(*s3.Options)) (*s3.ListObjectsV2Output, error)
	DeleteObject(ctx context.Context, params *s3.DeleteObjectInput, optFns ...func(*s3.Options)) (*s3.DeleteObjectOutput, error)
}

// New creates a new AWS S3 activity
func New(ctx activity.InitContext) (activity.Activity, error) {
	s := &Settings{}
	err := metadata.MapToStruct(ctx.Settings(), s, true)
	if err != nil {
		return nil, err
	}

	if s.Region == "" {
		return nil, fmt.Errorf("region is required for AWS S3 activity")
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

	client := s3.NewFromConfig(cfg)

	return &Activity{
		settings: s,
		client:   client,
	}, nil
}

// Metadata returns the metadata for the AWS S3 activity
func (act *Activity) Metadata() *activity.Metadata {
	return activityMetadata
}

// Eval executes the AWS S3 logic
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
	case "upload":
		err = act.upload(ctx, input, output)
	case "download":
		err = act.download(ctx, input, output)
	case "list":
		err = act.list(ctx, input, output)
	case "delete":
		err = act.delete(ctx, input, output)
	default:
		err = fmt.Errorf("unsupported method: %s", input.Method)
	}

	if err != nil {
		output.Success = false
		output.Error = err.Error()
		ctx.Logger().Errorf("AWS S3 operation %s failed: %v", input.Method, err)
	}

	err = ctx.SetOutputObject(output)
	if err != nil {
		return true, err
	}

	return true, nil
}

func (act *Activity) upload(ctx activity.Context, input *Input, output *Output) error {
	var body io.Reader
	switch d := input.Data.(type) {
	case string:
		body = strings.NewReader(d)
	case []byte:
		body = strings.NewReader(string(d))
	case io.Reader:
		body = d
	default:
		return fmt.Errorf("unsupported data type for upload: %T", input.Data)
	}

	_, err := act.client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(input.Bucket),
		Key:    aws.String(input.Key),
		Body:   body,
	})
	return err
}

func (act *Activity) download(ctx activity.Context, input *Input, output *Output) error {
	resp, err := act.client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(input.Bucket),
		Key:    aws.String(input.Key),
	})
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if input.LocalPath != "" {
		f, err := os.Create(input.LocalPath)
		if err != nil {
			return err
		}
		defer f.Close()
		_, err = io.Copy(f, resp.Body)
		return err
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	output.Result = string(data)
	return nil
}

func (act *Activity) list(ctx activity.Context, input *Input, output *Output) error {
	resp, err := act.client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: aws.String(input.Bucket),
		Prefix: aws.String(input.Key), // Use Key as Prefix for listing
	})
	if err != nil {
		return err
	}

	var keys []string
	for _, obj := range resp.Contents {
		keys = append(keys, *obj.Key)
	}
	output.Result = keys
	return nil
}

func (act *Activity) delete(ctx activity.Context, input *Input, output *Output) error {
	_, err := act.client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(input.Bucket),
		Key:    aws.String(input.Key),
	})
	return err
}
