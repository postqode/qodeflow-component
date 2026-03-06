package awsses

import (
	"github.com/postqode/qodeflow-core/activity"
	"github.com/postqode/qodeflow-core/data/coerce"
)

var activityMetadata = activity.ToMetadata(&Settings{}, &Input{}, &Output{})

// Settings represents the configuration for the AWS SES activity
type Settings struct {
	Region          string `md:"region"`
	AccessKeyID     string `md:"accessKeyID"`
	SecretAccessKey string `md:"secretAccessKey"`
	SessionToken    string `md:"sessionToken"`
}

// Input represents the input for the AWS SES activity
type Input struct {
	From    string   `md:"from"`
	To      []string `md:"to"`
	Subject string   `md:"subject"`
	Body    string   `md:"body"`
	Files   any      `md:"files"`
}

func (i *Input) ToMap() map[string]any {
	return map[string]any{
		"from":    i.From,
		"to":      i.To,
		"subject": i.Subject,
		"body":    i.Body,
		"files":   i.Files,
	}
}

func (i *Input) FromMap(values map[string]any) error {
	var err error
	i.From, err = coerce.ToString(values["from"])
	if err != nil {
		return err
	}
	if values["to"] != nil {
		toArr, err := coerce.ToArray(values["to"])
		if err != nil {
			return err
		}
		i.To = make([]string, len(toArr))
		for idx, v := range toArr {
			i.To[idx], _ = coerce.ToString(v)
		}
	}
	i.Subject, err = coerce.ToString(values["subject"])
	if err != nil {
		return err
	}
	i.Body, err = coerce.ToString(values["body"])
	if err != nil {
		return err
	}
	i.Files, _ = coerce.ToAny(values["files"])
	return nil
}

// Output represents the output for the AWS SES activity
type Output struct {
	Success bool   `md:"success"`
	Error   string `md:"error"`
}

func (o *Output) ToMap() map[string]any {
	return map[string]any{
		"success": o.Success,
		"error":   o.Error,
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
	return nil
}
