package email

import (
	"net/smtp"
	"os"
	"testing"

	"github.com/postqode/qodeflow-core/support/test"
	"github.com/stretchr/testify/assert"
)

func TestActivity_Eval(t *testing.T) {
	mailerCalled := false
	mockMailer := func(addr string, a smtp.Auth, from string, to []string, msg []byte) error {
		mailerCalled = true
		assert.Equal(t, "smtp.example.com:587", addr)
		assert.Equal(t, "user", from)
		assert.Equal(t, []string{"test@example.com"}, to)
		assert.Contains(t, string(msg), "Subject: Test Subject")
		assert.Contains(t, string(msg), "Test Body")
		return nil
	}

	act := &Activity{
		settings: &Settings{
			Host:     "smtp.example.com",
			Port:     "587",
			Username: "user",
			Password: "pass",
		},
		mailer: mockMailer,
	}
	tc := test.NewActivityContext(act.Metadata())

	input := &Input{
		To:      []string{"test@example.com"},
		Subject: "Test Subject",
		Body:    "Test Body",
	}
	tc.SetInputObject(input)

	done, err := act.Eval(tc)
	assert.Nil(t, err)
	assert.True(t, done)
	assert.True(t, mailerCalled)

	output := &Output{}
	err = tc.GetOutputObject(output)
	assert.Nil(t, err)
	assert.True(t, output.Success)
	assert.Empty(t, output.Error)
}

func TestActivity_Metadata(t *testing.T) {
	act := &Activity{}
	md := act.Metadata()
	assert.NotNil(t, md)
	assert.NotNil(t, md.Settings["host"])
	assert.NotNil(t, md.Input["to"])
}

func TestActivity_New_Validation(t *testing.T) {
	// 1. Missing all settings
	settings := map[string]any{}
	iCtx := test.NewActivityInitContext(settings, nil)
	act, err := New(iCtx)
	assert.Error(t, err)
	assert.Nil(t, act)
	assert.Contains(t, err.Error(), "host is required")

	// 2. Missing port
	settings = map[string]any{
		"host": "smtp.example.com",
	}
	iCtx = test.NewActivityInitContext(settings, nil)
	act, err = New(iCtx)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "port is required")

	// 3. Missing username/password
	settings = map[string]any{
		"host": "smtp.example.com",
		"port": "587",
	}
	iCtx = test.NewActivityInitContext(settings, nil)
	act, err = New(iCtx)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "username is required")
}

func TestActivity_Eval_Integration(t *testing.T) {
	host := os.Getenv("EMAIL_HOST")
	port := os.Getenv("EMAIL_PORT")
	user := os.Getenv("EMAIL_USERNAME")
	pass := os.Getenv("EMAIL_PASSWORD")

	if host == "" || port == "" || user == "" || pass == "" {
		t.Skip("Skipping integration test: EMAIL_HOST, EMAIL_PORT, EMAIL_USERNAME, or EMAIL_PASSWORD not set")
	}

	settings := map[string]any{
		"host":     host,
		"port":     port,
		"username": user,
		"password": pass,
	}

	iCtx := test.NewActivityInitContext(settings, nil)
	act, err := New(iCtx)
	assert.NoError(t, err)

	tc := test.NewActivityContext(act.Metadata())
	input := &Input{
		To:      []string{"vedant.khatri@postqode.com"}, // the person who receives the email. Replace this
		Subject: "Integration Test",
		Body:    "<h1>Integration Test Success</h1>",
	}
	tc.SetInputObject(input)

	done, err := act.Eval(tc)
	assert.NoError(t, err)
	assert.True(t, done)

	output := &Output{}
	err = tc.GetOutputObject(output)
	assert.NoError(t, err)
	assert.True(t, output.Success)
}

func TestActivity_Eval_WithAttachments(t *testing.T) {
	mailerCalled := false
	mockMailer := func(addr string, a smtp.Auth, from string, to []string, msg []byte) error {
		mailerCalled = true
		msgStr := string(msg)
		assert.Contains(t, msgStr, "Content-Type: multipart/mixed;")
		assert.Contains(t, msgStr, "<p>please find the attachments with this mail. Thanks.</p>")
		assert.Contains(t, msgStr, "Content-Disposition: attachment; filename=\"test-file.txt\"")
		assert.Contains(t, msgStr, "Content-Disposition: attachment; filename=\"attachment-2.dat\"")
		return nil
	}

	act := &Activity{
		settings: &Settings{
			Host:     "smtp.example.com",
			Port:     "587",
			Username: "user",
			Password: "pass",
		},
		mailer: mockMailer,
	}
	tc := test.NewActivityContext(act.Metadata())

	input := &Input{
		To:      []string{"test@example.com"}, // provide the emails whom email will be sent
		Subject: "Test Attachments",
		Body:    "<p>please find the attachments with this mail. Thanks.</p>",
		Files: []any{
			map[string]any{
				"filename": "test-file.txt",
				"data":     "plain text content",
				"mimeType": "text/plain",
			},
			[]byte("binary data"),
		},
	}
	tc.SetInputObject(input)

	done, err := act.Eval(tc)
	assert.NoError(t, err)
	assert.True(t, done)
	assert.True(t, mailerCalled)
}

func TestActivity_Eval_WithAttachments_Integration(t *testing.T) {
	host := os.Getenv("EMAIL_HOST")
	port := os.Getenv("EMAIL_PORT")
	user := os.Getenv("EMAIL_USERNAME")
	pass := os.Getenv("EMAIL_PASSWORD")

	if host == "" || port == "" || user == "" || pass == "" {
		t.Skip("Skipping integration test: EMAIL_HOST, EMAIL_PORT, EMAIL_USERNAME, or EMAIL_PASSWORD not set")
	}

	settings := map[string]any{
		"host":     host,
		"port":     port,
		"username": user,
		"password": pass,
	}

	iCtx := test.NewActivityInitContext(settings, nil)
	act, err := New(iCtx)
	assert.NoError(t, err)

	tc := test.NewActivityContext(act.Metadata())
	input := &Input{
		To:      []string{"test@example.com"}, // provide the emails whom email will be sent
		Subject: "Integration Test with Attachments",
		Body:    "<p>please find the attachments with this mail. Thanks.</p>",
		Files: []any{
			map[string]any{
				"filename": "integration-test.txt",
				"data":     "This is an attachment from integration test.",
				"mimeType": "text/plain",
			},
		},
	}
	tc.SetInputObject(input)

	done, err := act.Eval(tc)
	assert.NoError(t, err)
	assert.True(t, done)

	output := &Output{}
	err = tc.GetOutputObject(output)
	assert.NoError(t, err)
	assert.True(t, output.Success)
}
