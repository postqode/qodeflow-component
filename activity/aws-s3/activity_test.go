package awss3

import (
	"bytes"
	"context"
	"io"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/postqode/qodeflow-core/support/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockS3Client is a mock of S3Client
type MockS3Client struct {
	mock.Mock
}

func (m *MockS3Client) PutObject(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error) {
	args := m.Called(ctx, params, optFns)
	return args.Get(0).(*s3.PutObjectOutput), args.Error(1)
}

func (m *MockS3Client) GetObject(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
	args := m.Called(ctx, params, optFns)
	return args.Get(0).(*s3.GetObjectOutput), args.Error(1)
}

func (m *MockS3Client) ListObjectsV2(ctx context.Context, params *s3.ListObjectsV2Input, optFns ...func(*s3.Options)) (*s3.ListObjectsV2Output, error) {
	args := m.Called(ctx, params, optFns)
	return args.Get(0).(*s3.ListObjectsV2Output), args.Error(1)
}

func (m *MockS3Client) DeleteObject(ctx context.Context, params *s3.DeleteObjectInput, optFns ...func(*s3.Options)) (*s3.DeleteObjectOutput, error) {
	args := m.Called(ctx, params, optFns)
	return args.Get(0).(*s3.DeleteObjectOutput), args.Error(1)
}

func TestEval_Upload(t *testing.T) {
	mockClient := new(MockS3Client)
	act := &Activity{
		client: mockClient,
	}

	tc := test.NewActivityContext(activityMetadata)
	input := &Input{
		Method: "Upload",
		Bucket: "test-bucket",
		Key:    "test-key",
		Data:   "hello world",
	}
	tc.SetInputObject(input)

	mockClient.On("PutObject", mock.Anything, mock.MatchedBy(func(p *s3.PutObjectInput) bool {
		body, _ := io.ReadAll(p.Body)
		return *p.Bucket == "test-bucket" && *p.Key == "test-key" && string(body) == "hello world"
	}), mock.Anything).Return(&s3.PutObjectOutput{}, nil)

	done, err := act.Eval(tc)
	assert.NoError(t, err)
	assert.True(t, done)

	output := &Output{}
	tc.GetOutputObject(output)
	assert.True(t, output.Success)
	mockClient.AssertExpectations(t)
}

func TestEval_Download(t *testing.T) {
	mockClient := new(MockS3Client)
	act := &Activity{
		client: mockClient,
	}

	tc := test.NewActivityContext(activityMetadata)
	input := &Input{
		Method: "Download",
		Bucket: "test-bucket",
		Key:    "test-key",
	}
	tc.SetInputObject(input)

	mockClient.On("GetObject", mock.Anything, mock.MatchedBy(func(p *s3.GetObjectInput) bool {
		return *p.Bucket == "test-bucket" && *p.Key == "test-key"
	}), mock.Anything).Return(&s3.GetObjectOutput{
		Body: io.NopCloser(bytes.NewReader([]byte("downloaded content"))),
	}, nil)

	done, err := act.Eval(tc)
	assert.NoError(t, err)
	assert.True(t, done)

	output := &Output{}
	tc.GetOutputObject(output)
	assert.True(t, output.Success)
	assert.Equal(t, "downloaded content", output.Result)
	mockClient.AssertExpectations(t)
}

func TestEval_List(t *testing.T) {
	mockClient := new(MockS3Client)
	act := &Activity{
		client: mockClient,
	}

	tc := test.NewActivityContext(activityMetadata)
	input := &Input{
		Method: "List",
		Bucket: "test-bucket",
		Key:    "folder/",
	}
	tc.SetInputObject(input)

	mockClient.On("ListObjectsV2", mock.Anything, mock.MatchedBy(func(p *s3.ListObjectsV2Input) bool {
		return *p.Bucket == "test-bucket" && *p.Prefix == "folder/"
	}), mock.Anything).Return(&s3.ListObjectsV2Output{
		Contents: []types.Object{
			{Key: aws.String("folder/file1.txt")},
			{Key: aws.String("folder/file2.txt")},
		},
	}, nil)

	done, err := act.Eval(tc)
	assert.NoError(t, err)
	assert.True(t, done)

	output := &Output{}
	tc.GetOutputObject(output)
	assert.True(t, output.Success)
	assert.Equal(t, []string{"folder/file1.txt", "folder/file2.txt"}, output.Result)
	mockClient.AssertExpectations(t)
}

func TestEval_Delete(t *testing.T) {
	mockClient := new(MockS3Client)
	act := &Activity{
		client: mockClient,
	}

	tc := test.NewActivityContext(activityMetadata)
	input := &Input{
		Method: "Delete",
		Bucket: "test-bucket",
		Key:    "test-key",
	}
	tc.SetInputObject(input)

	mockClient.On("DeleteObject", mock.Anything, mock.MatchedBy(func(p *s3.DeleteObjectInput) bool {
		return *p.Bucket == "test-bucket" && *p.Key == "test-key"
	}), mock.Anything).Return(&s3.DeleteObjectOutput{}, nil)

	done, err := act.Eval(tc)
	assert.NoError(t, err)
	assert.True(t, done)

	output := &Output{}
	tc.GetOutputObject(output)
	assert.True(t, output.Success)
	mockClient.AssertExpectations(t)
}

func TestEval_Integration(t *testing.T) {
	region := os.Getenv("AWS_REGION")
	accessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	secretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	bucket := os.Getenv("AWS_S3_BUCKET")

	if region == "" || accessKey == "" || secretKey == "" || bucket == "" {
		t.Skip("Skipping integration test: AWS credentials or S3 bucket not set")
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

	// Test Upload
	input := &Input{
		Method: "Upload",
		Bucket: bucket,
		Key:    "test-integration.txt",
		Data:   "Integration test content",
	}
	tc.SetInputObject(input)

	done, err := act.Eval(tc)
	assert.NoError(t, err)
	assert.True(t, done)

	output := &Output{}
	tc.GetOutputObject(output)
	assert.True(t, output.Success)

	// Test Delete (Cleanup)
	input.Method = "Delete"
	tc.SetInputObject(input)
	done, err = act.Eval(tc)
	assert.NoError(t, err)
	assert.True(t, done)

	tc.GetOutputObject(output)
	assert.True(t, output.Success)
}
