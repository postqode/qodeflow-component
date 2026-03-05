package email

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"mime/multipart"
	"net/smtp"
	"net/textproto"
	"os"
	"path/filepath"
	"strings"

	"github.com/postqode/qodeflow-core/activity"
	"github.com/postqode/qodeflow-core/data/coerce"
	"github.com/postqode/qodeflow-core/data/metadata"
)

func init() {
	_ = activity.Register(&Activity{}, New)
}

// MailerFunc represents a function that sends an email
type MailerFunc func(addr string, a smtp.Auth, from string, to []string, msg []byte) error

// Activity represents the email activity
type Activity struct {
	settings *Settings
	mailer   MailerFunc
}

// New creates a new email activity
func New(ctx activity.InitContext) (activity.Activity, error) {
	s := &Settings{}
	err := metadata.MapToStruct(ctx.Settings(), s, true)
	if err != nil {
		return nil, err
	}

	// Fallback to environment variables if settings are empty
	if s.Host == "" {
		s.Host = os.Getenv("EMAIL_HOST")
	}
	if s.Port == "" {
		s.Port = os.Getenv("EMAIL_PORT")
	}
	if s.Username == "" {
		s.Username = os.Getenv("EMAIL_USERNAME")
	}
	if s.Password == "" {
		s.Password = os.Getenv("EMAIL_PASSWORD")
	}

	// Validation
	if s.Host == "" {
		return nil, fmt.Errorf("host is required for email activity")
	}
	if s.Port == "" {
		return nil, fmt.Errorf("port is required for email activity")
	}
	if s.Username == "" {
		return nil, fmt.Errorf("username is required for email activity")
	}
	if s.Password == "" {
		return nil, fmt.Errorf("password is required for email activity")
	}

	return &Activity{
		settings: s,
		mailer:   smtp.SendMail,
	}, nil
}

// Metadata returns the metadata for the email activity
func (act *Activity) Metadata() *activity.Metadata {
	return activityMetadata
}

// Eval executes the email sending logic
func (act *Activity) Eval(ctx activity.Context) (done bool, err error) {
	input := &Input{}
	err = ctx.GetInputObject(input)
	if err != nil {
		return true, err
	}

	// 1. Setup Authentication
	auth := smtp.PlainAuth("", act.settings.Username, act.settings.Password, act.settings.Host)
	addr := fmt.Sprintf("%s:%s", act.settings.Host, act.settings.Port)

	var msg []byte

	// 2. Format the Message
	if input.Files == nil {
		// Simple HTML message (Backward compatibility)
		headerTo := strings.Join(input.To, ",")
		mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
		subjectHeader := fmt.Sprintf("Subject: %s\n", input.Subject)
		toHeader := fmt.Sprintf("To: %s\n", headerTo)
		msg = []byte(toHeader + subjectHeader + mime + input.Body)
	} else {
		// Multipart message for attachments
		buf := new(bytes.Buffer)
		writer := multipart.NewWriter(buf)

		// Main Headers
		buf.WriteString(fmt.Sprintf("To: %s\r\n", strings.Join(input.To, ",")))
		buf.WriteString(fmt.Sprintf("Subject: %s\r\n", input.Subject))
		buf.WriteString("MIME-Version: 1.0\r\n")
		buf.WriteString(fmt.Sprintf("Content-Type: multipart/mixed; boundary=%s\r\n", writer.Boundary()))
		buf.WriteString("\r\n")

		// Body Part (HTML)
		bodyHeader := make(textproto.MIMEHeader)
		bodyHeader.Set("Content-Type", "text/html; charset=\"UTF-8\"")
		bodyHeader.Set("Content-Transfer-Encoding", "quoted-printable")
		bodyPart, err := writer.CreatePart(bodyHeader)
		if err != nil {
			return true, err
		}
		_, _ = bodyPart.Write([]byte(input.Body))

		// Process files
		var files []any
		if fArr, ok := input.Files.([]any); ok {
			files = fArr
		} else {
			files = []any{input.Files}
		}

		for i, file := range files {
			var fileData []byte
			var fileName string
			var mimeType string

			switch v := file.(type) {
			case []byte:
				fileData = v
				fileName = fmt.Sprintf("attachment-%d.dat", i+1)
				mimeType = "application/octet-stream"
			case string:
				fileData = []byte(v)
				fileName = fmt.Sprintf("attachment-%d.txt", i+1)
				mimeType = "text/plain"
			case map[string]any:
				if data, ok := v["data"]; ok {
					fileData, _ = coerce.ToBytes(data)
				}
				if name, ok := v["filename"].(string); ok {
					fileName = name
				} else if name, ok := v["name"].(string); ok {
					fileName = name
				}
				if mt, ok := v["mimeType"].(string); ok {
					mimeType = mt
				}
			}

			if len(fileData) > 0 {
				if fileName == "" {
					fileName = fmt.Sprintf("attachment-%d.dat", i+1)
				}
				if mimeType == "" {
					mimeType = "application/octet-stream"
				}

				partHeader := make(textproto.MIMEHeader)
				partHeader.Set("Content-Type", mimeType)
				partHeader.Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filepath.Base(fileName)))
				partHeader.Set("Content-Transfer-Encoding", "base64")

				part, err := writer.CreatePart(partHeader)
				if err != nil {
					return true, err
				}

				encoded := make([]byte, base64.StdEncoding.EncodedLen(len(fileData)))
				base64.StdEncoding.Encode(encoded, fileData)
				_, _ = part.Write(encoded)
			}
		}

		writer.Close()
		msg = buf.Bytes()
	}

	// 3. Send the email
	err = act.mailer(addr, auth, act.settings.Username, input.To, msg)

	output := &Output{
		Success: err == nil,
	}

	if err != nil {
		output.Error = err.Error()
		ctx.Logger().Errorf("Failed to send email: %v", err)
	}

	err = ctx.SetOutputObject(output)
	if err != nil {
		return true, err
	}

	return true, nil
}
