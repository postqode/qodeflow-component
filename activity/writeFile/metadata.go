package writefile

import (
	"github.com/postqode/qodeflow-core/activity"
	"github.com/postqode/qodeflow-core/data/coerce"
)

var activityMetadata = activity.ToMetadata(&Settings{}, &Input{}, &Output{})

// Settings represents the configuration for the Write File activity
type Settings struct{}

// Input represents the input for the Write File activity
type Input struct {
	Method   string `md:"method"`
	FilePath string `md:"filePath"`
	Content  any    `md:"content"`
}

func (i *Input) ToMap() map[string]any {
	return map[string]any{
		"method":   i.Method,
		"filePath": i.FilePath,
		"content":  i.Content,
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
	i.Content, _ = coerce.ToAny(values["content"])
	return nil
}

// Output represents the output for the Write File activity
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
