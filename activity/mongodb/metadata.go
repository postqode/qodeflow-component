package mongodb

import (
	"github.com/postqode/qodeflow-core/activity"
	"github.com/postqode/qodeflow-core/data/coerce"
)

var activityMetadata = activity.ToMetadata(&Settings{}, &Input{}, &Output{})

// Settings represents the configuration for the MongoDB activity
type Settings struct {
	ConnectionURI string `md:"uri"`
	DbName        string `md:"dbName"`
	Collection    string `md:"collection"`
}

// Input represents the input for the MongoDB activity
type Input struct {
	Method   string      `md:"method"`
	KeyName  string      `md:"keyName"`
	KeyValue string      `md:"keyValue"`
	Data     interface{} `md:"data"`
}

func (i *Input) FromMap(values map[string]interface{}) error {
	var err error
	i.Method, err = coerce.ToString(values["method"])
	if err != nil {
		return err
	}
	i.KeyName, err = coerce.ToString(values["keyName"])
	if err != nil {
		return err
	}
	i.KeyValue, err = coerce.ToString(values["keyValue"])
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
		"method":   i.Method,
		"keyName":  i.KeyName,
		"keyValue": i.KeyValue,
		"data":     i.Data,
	}
}

// Output represents the output for the MongoDB activity
type Output struct {
	Output interface{} `md:"output"`
	Count  int64       `md:"count"`
}

func (o *Output) FromMap(values map[string]interface{}) error {
	var err error
	o.Output, err = coerce.ToAny(values["output"])
	if err != nil {
		return err
	}
	o.Count, err = coerce.ToInt64(values["count"])
	if err != nil {
		return err
	}
	return nil
}

func (o *Output) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"output": o.Output,
		"count":  o.Count,
	}
}
