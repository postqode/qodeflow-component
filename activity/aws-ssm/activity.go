package awsssm

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"github.com/postqode/qodeflow-core/activity"
	"github.com/postqode/qodeflow-core/data/metadata"
)

func init() {
	_ = activity.Register(&Activity{}, New)
}

// Activity represents the AWS SSM activity
type Activity struct {
	settings *Settings
	client   SSMClient
}

// SSMClient interface for mocking in tests
type SSMClient interface {
	GetParameter(ctx context.Context, params *ssm.GetParameterInput, optFns ...func(*ssm.Options)) (*ssm.GetParameterOutput, error)
	PutParameter(ctx context.Context, params *ssm.PutParameterInput, optFns ...func(*ssm.Options)) (*ssm.PutParameterOutput, error)
	DeleteParameter(ctx context.Context, params *ssm.DeleteParameterInput, optFns ...func(*ssm.Options)) (*ssm.DeleteParameterOutput, error)
	GetParametersByPath(ctx context.Context, params *ssm.GetParametersByPathInput, optFns ...func(*ssm.Options)) (*ssm.GetParametersByPathOutput, error)
	DescribeParameters(ctx context.Context, params *ssm.DescribeParametersInput, optFns ...func(*ssm.Options)) (*ssm.DescribeParametersOutput, error)
}

// New creates a new AWS SSM activity
func New(ctx activity.InitContext) (activity.Activity, error) {
	s := &Settings{}
	err := metadata.MapToStruct(ctx.Settings(), s, true)
	if err != nil {
		return nil, err
	}

	if s.Region == "" {
		return nil, fmt.Errorf("region is required for AWS SSM activity")
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

	client := ssm.NewFromConfig(cfg)

	return &Activity{
		settings: s,
		client:   client,
	}, nil
}

// Metadata returns the metadata for the AWS SSM activity
func (act *Activity) Metadata() *activity.Metadata {
	return activityMetadata
}

// Eval executes the AWS SSM logic
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
	case "getparameter":
		err = act.getParameter(ctx, input, output)
	case "putparameter":
		err = act.putParameter(ctx, input, output)
	case "deleteparameter":
		err = act.deleteParameter(ctx, input, output)
	case "getparametersbypath":
		err = act.getParametersByPath(ctx, input, output)
	case "describeparameters":
		err = act.describeParameters(ctx, input, output)
	default:
		err = fmt.Errorf("unsupported method: %s", input.Method)
	}

	if err != nil {
		output.Success = false
		output.Error = err.Error()
		ctx.Logger().Errorf("AWS SSM operation %s failed: %v", input.Method, err)
	}

	err = ctx.SetOutputObject(output)
	if err != nil {
		return true, err
	}

	return true, nil
}

func (act *Activity) getParameter(ctx activity.Context, input *Input, output *Output) error {
	if input.ParameterName == "" {
		return fmt.Errorf("parameterName is required for GetParameter")
	}

	resp, err := act.client.GetParameter(context.TODO(), &ssm.GetParameterInput{
		Name:           aws.String(input.ParameterName),
		WithDecryption: aws.Bool(input.WithDecryption),
	})
	if err != nil {
		return err
	}

	output.Result = map[string]any{
		"name":  aws.ToString(resp.Parameter.Name),
		"value": aws.ToString(resp.Parameter.Value),
		"type":  string(resp.Parameter.Type),
	}
	return nil
}

func (act *Activity) putParameter(ctx activity.Context, input *Input, output *Output) error {
	if input.ParameterName == "" {
		return fmt.Errorf("parameterName is required for PutParameter")
	}
	if input.ParameterValue == "" {
		return fmt.Errorf("parameterValue is required for PutParameter")
	}

	paramType := types.ParameterTypeString
	if input.ParameterType != "" {
		paramType = types.ParameterType(input.ParameterType)
	}

	resp, err := act.client.PutParameter(context.TODO(), &ssm.PutParameterInput{
		Name:      aws.String(input.ParameterName),
		Value:     aws.String(input.ParameterValue),
		Type:      paramType,
		Overwrite: aws.Bool(input.Overwrite),
	})
	if err != nil {
		return err
	}

	output.Result = map[string]any{
		"version": resp.Version,
	}
	return nil
}

func (act *Activity) deleteParameter(ctx activity.Context, input *Input, output *Output) error {
	if input.ParameterName == "" {
		return fmt.Errorf("parameterName is required for DeleteParameter")
	}

	_, err := act.client.DeleteParameter(context.TODO(), &ssm.DeleteParameterInput{
		Name: aws.String(input.ParameterName),
	})
	return err
}

func (act *Activity) getParametersByPath(ctx activity.Context, input *Input, output *Output) error {
	if input.Path == "" {
		return fmt.Errorf("path is required for GetParametersByPath")
	}

	resp, err := act.client.GetParametersByPath(context.TODO(), &ssm.GetParametersByPathInput{
		Path:           aws.String(input.Path),
		Recursive:      aws.Bool(input.Recursive),
		WithDecryption: aws.Bool(input.WithDecryption),
	})
	if err != nil {
		return err
	}

	var params []map[string]any
	for _, p := range resp.Parameters {
		params = append(params, map[string]any{
			"name":  aws.ToString(p.Name),
			"value": aws.ToString(p.Value),
			"type":  string(p.Type),
		})
	}

	output.Result = params
	return nil
}

func (act *Activity) describeParameters(ctx activity.Context, input *Input, output *Output) error {
	resp, err := act.client.DescribeParameters(context.TODO(), &ssm.DescribeParametersInput{})
	if err != nil {
		return err
	}

	var params []map[string]any
	for _, p := range resp.Parameters {
		params = append(params, map[string]any{
			"name":        aws.ToString(p.Name),
			"type":        string(p.Type),
			"description": aws.ToString(p.Description),
		})
	}

	output.Result = params
	return nil
}
