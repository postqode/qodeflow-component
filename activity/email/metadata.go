package email

import (
	"github.com/postqode/qodeflow-core/activity"
	"github.com/postqode/qodeflow-core/data/coerce"
)

var activityMetadata = activity.ToMetadata(&Settings{}, &Input{}, &Output{})

// Settings represents the configuration for the email activity
type Settings struct {
	Host     string `md:"host"`
	Port     string `md:"port"`
	Username string `md:"username"`
	Password string `md:"password"`
}

// Input represents the input for the email activity
type Input struct {
	To      []string    `md:"to"`
	Subject string      `md:"subject"`
	Body    string      `md:"body"`
	Files   interface{} `md:"files"`
}

func (i *Input) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"to":      i.To,
		"subject": i.Subject,
		"body":    i.Body,
		"files":   i.Files,
	}
}

func (i *Input) FromMap(values map[string]interface{}) error {
	var err error
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

// Output represents the output for the email activity
type Output struct {
	Success bool   `md:"success"`
	Error   string `md:"error"`
}

func (o *Output) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"success": o.Success,
		"error":   o.Error,
	}
}

func (o *Output) FromMap(values map[string]interface{}) error {
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
