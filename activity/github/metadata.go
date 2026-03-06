package github

import (
	"github.com/postqode/qodeflow-core/activity"
	"github.com/postqode/qodeflow-core/data/coerce"
)

var activityMetadata = activity.ToMetadata(&Settings{}, &Input{}, &Output{})

// Settings represents the configuration for the GitHub activity
type Settings struct {
	Token string `md:"token"`
}

// Input represents the input for the GitHub activity
type Input struct {
	Owner  string      `md:"owner"`
	Repo   string      `md:"repo"`
	Method string      `md:"method"`
	Data   interface{} `md:"data"`
}

func (i *Input) FromMap(values map[string]interface{}) error {
	var err error
	i.Owner, err = coerce.ToString(values["owner"])
	if err != nil {
		return err
	}
	i.Repo, err = coerce.ToString(values["repo"])
	if err != nil {
		return err
	}
	i.Method, err = coerce.ToString(values["method"])
	if err != nil {
		return err
	}
	i.Data, err = coerce.ToAny(values["data"])
	if err != nil {
		return err
	}
	return nil
}

func (i *Input) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"owner":  i.Owner,
		"repo":   i.Repo,
		"method": i.Method,
		"data":   i.Data,
	}
}

// Output represents the output for the GitHub activity
type Output struct {
	Result interface{} `md:"result"`
	Error  string      `md:"error"`
}

func (o *Output) FromMap(values map[string]interface{}) error {
	var err error
	o.Result, err = coerce.ToAny(values["result"])
	if err != nil {
		return err
	}
	o.Error, err = coerce.ToString(values["error"])
	if err != nil {
		return err
	}
	return nil
}

func (o *Output) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"result": o.Result,
		"error":  o.Error,
	}
}
