package awss3

import (
	"github.com/postqode/qodeflow-core/activity"
	"github.com/postqode/qodeflow-core/data/coerce"
)

var activityMetadata = activity.ToMetadata(&Settings{}, &Input{}, &Output{})

// Settings represents the configuration for the AWS S3 activity
type Settings struct {
	Region          string `md:"region"`
	AccessKeyID     string `md:"accessKeyID"`
	SecretAccessKey string `md:"secretAccessKey"`
	SessionToken    string `md:"sessionToken"`
}

// Input represents the input for the AWS S3 activity
type Input struct {
	Method    string `md:"method"`
	Bucket    string `md:"bucket"`
	Key       string `md:"key"`
	Data      any    `md:"data"`
	LocalPath string `md:"localPath"`
}

func (i *Input) ToMap() map[string]any {
	return map[string]any{
		"method":    i.Method,
		"bucket":    i.Bucket,
		"key":       i.Key,
		"data":      i.Data,
		"localPath": i.LocalPath,
	}
}

func (i *Input) FromMap(values map[string]any) error {
	var err error
	i.Method, err = coerce.ToString(values["method"])
	if err != nil {
		return err
	}
	i.Bucket, err = coerce.ToString(values["bucket"])
	if err != nil {
		return err
	}
	i.Key, err = coerce.ToString(values["key"])
	if err != nil {
		return err
	}
	i.Data, _ = coerce.ToAny(values["data"])
	i.LocalPath, err = coerce.ToString(values["localPath"])
	if err != nil {
		return err
	}
	return nil
}

// Output represents the output for the AWS S3 activity
type Output struct {
	Success bool   `md:"success"`
	Error   string `md:"error"`
	Result  any    `md:"result"`
}

func (o *Output) ToMap() map[string]any {
	return map[string]any{
		"success": o.Success,
		"error":   o.Error,
		"result":  o.Result,
	}
}

func (o *Output) FromMap(values map[string]any) error {
	var err error
	o.Success, err = coerce.ToBool(values["success"])
	if err != nil {
		return err
	}
	o.Error, err = coerce.ToString(values["error"])
	if err != nil {
		return err
	}
	o.Result, _ = coerce.ToAny(values["result"])
	return nil
}
