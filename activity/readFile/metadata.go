package readfile

import (
	"github.com/postqode/qodeflow-core/activity"
	"github.com/postqode/qodeflow-core/data/coerce"
)

var activityMetadata = activity.ToMetadata(&Settings{}, &Input{}, &Output{})

// Settings represents the configuration for the Read File activity
type Settings struct{}

// Input represents the input for the Read File activity
type Input struct {
	Method   string `md:"method"`
	FilePath string `md:"filePath"`
	Encoding string `md:"encoding"`
}

func (i *Input) ToMap() map[string]any {
	return map[string]any{
		"method":   i.Method,
		"filePath": i.FilePath,
		"encoding": i.Encoding,
	}
}

func (i *Input) FromMap(values map[string]any) error {
	var err error
	i.Method, err = coerce.ToString(values["method"])
	if err != nil {
		return err
	}
	i.FilePath, err = coerce.ToString(values["filePath"])
	if err != nil {
		return err
	}
	i.Encoding, err = coerce.ToString(values["encoding"])
	if err != nil {
		return err
	}
	return nil
}

// Output represents the output for the Read File activity
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
