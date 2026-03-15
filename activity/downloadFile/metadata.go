package downloadfile

import (
	"github.com/postqode/qodeflow-core/activity"
	"github.com/postqode/qodeflow-core/data/coerce"
)

var activityMetadata = activity.ToMetadata(&Settings{}, &Input{}, &Output{})

// Settings represents the configuration for the Download File activity
type Settings struct{}

// Input represents the input for the Download File activity
type Input struct {
	URL         string         `md:"url,required"`
	Destination string         `md:"destination,required"`
	Method      string         `md:"method"`
	Headers     map[string]any `md:"headers"`
	Timeout     int            `md:"timeout"`
	Append      bool           `md:"append"`
}

func (i *Input) ToMap() map[string]any {
	return map[string]any{
		"url":         i.URL,
		"destination": i.Destination,
		"method":      i.Method,
		"headers":     i.Headers,
		"timeout":     i.Timeout,
		"append":      i.Append,
	}
}

func (i *Input) FromMap(values map[string]any) error {
	var err error
	i.URL, err = coerce.ToString(values["url"])
	if err != nil {
		return err
	}
	i.Destination, err = coerce.ToString(values["destination"])
	if err != nil {
		return err
	}
	i.Method, _ = coerce.ToString(values["method"])
	if i.Method == "" {
		i.Method = "GET"
	}

	if values["headers"] != nil {
		i.Headers, err = coerce.ToObject(values["headers"])
		if err != nil {
			return err
		}
	}

	i.Timeout, _ = coerce.ToInt(values["timeout"])
	i.Append, _ = coerce.ToBool(values["append"])
	return nil
}

// Output represents the output for the Download File activity
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
	o.Error, _ = coerce.ToString(values["error"])
	o.Result, _ = coerce.ToAny(values["result"])
	return nil
}
