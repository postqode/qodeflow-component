package awsssm

import (
	"context"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"github.com/postqode/qodeflow-core/support/test"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockSSMClient struct {
	mock.Mock
}

func (m *MockSSMClient) GetParameter(ctx context.Context, params *ssm.GetParameterInput, optFns ...func(*ssm.Options)) (*ssm.GetParameterOutput, error) {
	args := m.Called(ctx, params, optFns)
	return args.Get(0).(*ssm.GetParameterOutput), args.Error(1)
}

func (m *MockSSMClient) PutParameter(ctx context.Context, params *ssm.PutParameterInput, optFns ...func(*ssm.Options)) (*ssm.PutParameterOutput, error) {
	args := m.Called(ctx, params, optFns)
	return args.Get(0).(*ssm.PutParameterOutput), args.Error(1)
}

func (m *MockSSMClient) DeleteParameter(ctx context.Context, params *ssm.DeleteParameterInput, optFns ...func(*ssm.Options)) (*ssm.DeleteParameterOutput, error) {
	args := m.Called(ctx, params, optFns)
	return args.Get(0).(*ssm.DeleteParameterOutput), args.Error(1)
}

func (m *MockSSMClient) GetParametersByPath(ctx context.Context, params *ssm.GetParametersByPathInput, optFns ...func(*ssm.Options)) (*ssm.GetParametersByPathOutput, error) {
	args := m.Called(ctx, params, optFns)
	return args.Get(0).(*ssm.GetParametersByPathOutput), args.Error(1)
}

func (m *MockSSMClient) DescribeParameters(ctx context.Context, params *ssm.DescribeParametersInput, optFns ...func(*ssm.Options)) (*ssm.DescribeParametersOutput, error) {
	args := m.Called(ctx, params, optFns)
	return args.Get(0).(*ssm.DescribeParametersOutput), args.Error(1)
}

func TestSSMActivity_GetParameter(t *testing.T) {
	mockClient := new(MockSSMClient)
	act := &Activity{
		settings: &Settings{Region: "us-east-1"},
		client:   mockClient,
	}

	paramName := "test-param"
	paramValue := "test-value"

	mockClient.On("GetParameter", mock.Anything, mock.MatchedBy(func(input *ssm.GetParameterInput) bool {
		return *input.Name == paramName
	}), mock.Anything).Return(&ssm.GetParameterOutput{
		Parameter: &types.Parameter{
			Name:  aws.String(paramName),
			Value: aws.String(paramValue),
			Type:  types.ParameterTypeString,
		},
	}, nil)

	tc := test.NewActivityContext(activityMetadata)
	input := &Input{
		Method:        "GetParameter",
		ParameterName: paramName,
	}
	tc.SetInputObject(input)

	done, err := act.Eval(tc)
	assert.True(t, done)
	assert.Nil(t, err)

	output := &Output{}
	tc.GetOutputObject(output)
	assert.True(t, output.Success)
	result := output.Result.(map[string]any)
	assert.Equal(t, paramName, result["name"])
	assert.Equal(t, paramValue, result["value"])
}

func TestSSMActivity_PutParameter(t *testing.T) {
	mockClient := new(MockSSMClient)
	act := &Activity{
		settings: &Settings{Region: "us-east-1"},
		client:   mockClient,
	}

	paramName := "test-param"
	paramValue := "test-value"

	mockClient.On("PutParameter", mock.Anything, mock.MatchedBy(func(input *ssm.PutParameterInput) bool {
		return *input.Name == paramName && *input.Value == paramValue
	}), mock.Anything).Return(&ssm.PutParameterOutput{Version: 1}, nil)

	tc := test.NewActivityContext(activityMetadata)
	input := &Input{
		Method:         "PutParameter",
		ParameterName:  paramName,
		ParameterValue: paramValue,
		ParameterType:  "String",
		Overwrite:      true,
	}
	tc.SetInputObject(input)

	done, err := act.Eval(tc)
	assert.True(t, done)
	assert.Nil(t, err)

	output := &Output{}
	tc.GetOutputObject(output)
	assert.True(t, output.Success)
	assert.Equal(t, int64(1), output.Result.(map[string]any)["version"])
}

func TestSSMActivity_DeleteParameter(t *testing.T) {
	mockClient := new(MockSSMClient)
	act := &Activity{
		settings: &Settings{Region: "us-east-1"},
		client:   mockClient,
	}

	paramName := "test-param"

	mockClient.On("DeleteParameter", mock.Anything, mock.MatchedBy(func(input *ssm.DeleteParameterInput) bool {
		return *input.Name == paramName
	}), mock.Anything).Return(&ssm.DeleteParameterOutput{}, nil)

	tc := test.NewActivityContext(activityMetadata)
	input := &Input{
		Method:        "DeleteParameter",
		ParameterName: paramName,
	}
	tc.SetInputObject(input)

	done, err := act.Eval(tc)
	assert.True(t, done)
	assert.Nil(t, err)

	output := &Output{}
	tc.GetOutputObject(output)
	assert.True(t, output.Success)
}

func TestSSMActivity_Integration(t *testing.T) {
	region := os.Getenv("AWS_REGION")
	accessKey := os.Getenv("AWS_ACCESS_KEY_ID")
	secretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")

	if region == "" || accessKey == "" || secretKey == "" {
		t.Skip("Skipping integration test: AWS credentials not set")
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

	paramName := "/qodeflow/test/param"
	paramValue := "integration-test-value"

	// Test PutParameter
	input := &Input{
		Method:         "PutParameter",
		ParameterName:  paramName,
		ParameterValue: paramValue,
		ParameterType:  "String",
		Overwrite:      true,
	}
	tc.SetInputObject(input)

	done, err := act.Eval(tc)
	assert.NoError(t, err)
	assert.True(t, done)

	output := &Output{}
	tc.GetOutputObject(output)
	assert.True(t, output.Success)

	// Test GetParameter
	input = &Input{
		Method:        "GetParameter",
		ParameterName: paramName,
	}
	tc.SetInputObject(input)
	done, err = act.Eval(tc)
	assert.NoError(t, err)

	tc.GetOutputObject(output)
	assert.True(t, output.Success)
	result := output.Result.(map[string]any)
	assert.Equal(t, paramValue, result["value"])

	// Test DeleteParameter
	input = &Input{
		Method:        "DeleteParameter",
		ParameterName: paramName,
	}
	tc.SetInputObject(input)
	done, err = act.Eval(tc)
	assert.NoError(t, err)

	tc.GetOutputObject(output)
	assert.True(t, output.Success)
}
