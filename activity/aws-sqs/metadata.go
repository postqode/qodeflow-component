package awssqs

import (
	"github.com/postqode/qodeflow-core/activity"
	"github.com/postqode/qodeflow-core/data/coerce"
)

var activityMetadata = activity.ToMetadata(&Settings{}, &Input{}, &Output{})

// Settings represents the configuration for the AWS SQS activity
type Settings struct {
	Region          string `md:"region"`
	AccessKeyID     string `md:"accessKeyID"`
	SecretAccessKey string `md:"secretAccessKey"`
	SessionToken    string `md:"sessionToken"`
}

// Input represents the input for the AWS SQS activity
type Input struct {
	Method              string `md:"method"`
	QueueURL            string `md:"queueURL"`
	MessageBody         string `md:"messageBody"`
	ReceiptHandle       string `md:"receiptHandle"`
	MaxNumberOfMessages int    `md:"maxNumberOfMessages"`
	WaitTimeSeconds     int    `md:"waitTimeSeconds"`
}

func (i *Input) ToMap() map[string]any {
	return map[string]any{
		"method":              i.Method,
		"queueURL":            i.QueueURL,
		"messageBody":         i.MessageBody,
		"receiptHandle":       i.ReceiptHandle,
		"maxNumberOfMessages": i.MaxNumberOfMessages,
		"waitTimeSeconds":     i.WaitTimeSeconds,
	}
}

func (i *Input) FromMap(values map[string]any) error {
	var err error
	i.Method, err = coerce.ToString(values["method"])
	if err != nil {
		return err
	}
	i.QueueURL, err = coerce.ToString(values["queueURL"])
	if err != nil {
		return err
	}
	i.MessageBody, err = coerce.ToString(values["messageBody"])
	if err != nil {
		return err
	}
	i.ReceiptHandle, err = coerce.ToString(values["receiptHandle"])
	if err != nil {
		return err
	}
	i.MaxNumberOfMessages, err = coerce.ToInt(values["maxNumberOfMessages"])
	if err != nil {
		return err
	}
	i.WaitTimeSeconds, err = coerce.ToInt(values["waitTimeSeconds"])
	if err != nil {
		return err
	}
	return nil
}

// Output represents the output for the AWS SQS activity
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
