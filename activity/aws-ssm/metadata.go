package awsssm

import (
	"github.com/postqode/qodeflow-core/activity"
	"github.com/postqode/qodeflow-core/data/coerce"
)

var activityMetadata = activity.ToMetadata(&Settings{}, &Input{}, &Output{})

// Settings represents the configuration for the AWS SSM activity
type Settings struct {
	Region          string `md:"region"`
	AccessKeyID     string `md:"accessKeyID"`
	SecretAccessKey string `md:"secretAccessKey"`
	SessionToken    string `md:"sessionToken"`
}

// Input represents the input for the AWS SSM activity
type Input struct {
	Method         string `md:"method"`
	ParameterName  string `md:"parameterName"`
	ParameterValue string `md:"parameterValue"`
	ParameterType  string `md:"parameterType"`
	Path           string `md:"path"`
	Recursive      bool   `md:"recursive"`
	WithDecryption bool   `md:"withDecryption"`
	Overwrite      bool   `md:"overwrite"`
}

func (i *Input) ToMap() map[string]any {
	return map[string]any{
		"method":         i.Method,
		"parameterName":  i.ParameterName,
		"parameterValue": i.ParameterValue,
		"parameterType":  i.ParameterType,
		"path":           i.Path,
		"recursive":      i.Recursive,
		"withDecryption": i.WithDecryption,
		"overwrite":      i.Overwrite,
	}
}

func (i *Input) FromMap(values map[string]any) error {
	var err error
	i.Method, err = coerce.ToString(values["method"])
	if err != nil {
		return err
	}
	i.ParameterName, err = coerce.ToString(values["parameterName"])
	if err != nil {
		return err
	}
	i.ParameterValue, err = coerce.ToString(values["parameterValue"])
	if err != nil {
		return err
	}
	i.ParameterType, err = coerce.ToString(values["parameterType"])
	if err != nil {
		return err
	}
	i.Path, err = coerce.ToString(values["path"])
	if err != nil {
		return err
	}
	i.Recursive, err = coerce.ToBool(values["recursive"])
	if err != nil {
		return err
	}
	i.WithDecryption, err = coerce.ToBool(values["withDecryption"])
	if err != nil {
		return err
	}
	i.Overwrite, err = coerce.ToBool(values["overwrite"])
	if err != nil {
		return err
	}
	return nil
}

// Output represents the output for the AWS SSM activity
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
