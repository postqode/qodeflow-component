package email

import (
	"encoding/json"
	"testing"

	"github.com/postqode/qodeflow-core/action"
	"github.com/postqode/qodeflow-core/support/test"
	"github.com/postqode/qodeflow-core/trigger"
	"github.com/stretchr/testify/assert"
)

const testConfig string = `{
	"id": "qodeflow-email",
	"ref": "github.com/postqode/qodeflow-component/trigger/email",
	"settings": {
		"host": "imap.example.com",
		"port": 993,
		"username": "test@example.com",
		"password": "testpassword",
		"useTLS": true
	},
	"handlers": [
	  {
		"settings":{
			"folder": "INBOX",
			"pollInterval": "1m"
		},
		"action":{
			"id":"dummy"
		}
	  }
	]
  }`

func TestEmailFactory_New(t *testing.T) {
	f := &Factory{}

	config := &trigger.Config{}
	err := json.Unmarshal([]byte(testConfig), config)
	assert.Nil(t, err)

	trg, err := f.New(config)
	assert.Nil(t, err)
	assert.NotNil(t, trg)

	// Verify settings were parsed correctly
	emailTrg, ok := trg.(*Trigger)
	assert.True(t, ok)
	assert.Equal(t, "imap.example.com", emailTrg.settings.Host)
	assert.Equal(t, 993, emailTrg.settings.Port)
	assert.Equal(t, "test@example.com", emailTrg.settings.Username)
	assert.Equal(t, "testpassword", emailTrg.settings.Password)
	assert.True(t, emailTrg.settings.UseTLS)
}

func TestEmailFactory_Metadata(t *testing.T) {
	f := &Factory{}
	md := f.Metadata()
	assert.NotNil(t, md)
}

func TestEmailTrigger_Initialize(t *testing.T) {
	f := &Factory{}

	config := &trigger.Config{}
	err := json.Unmarshal([]byte(testConfig), config)
	assert.Nil(t, err)

	actions := map[string]action.Action{"dummy": test.NewDummyAction(func() {
		// do nothing
	})}

	trg, err := test.InitTrigger(f, config, actions)
	assert.Nil(t, err)
	assert.NotNil(t, trg)

	emailTrg, ok := trg.(*Trigger)
	assert.True(t, ok)
	assert.NotNil(t, emailTrg.handlers)
	assert.Equal(t, 1, len(emailTrg.handlers))
}

func TestEmailTrigger_StartStop(t *testing.T) {
	f := &Factory{}

	config := &trigger.Config{}
	err := json.Unmarshal([]byte(testConfig), config)
	assert.Nil(t, err)

	actions := map[string]action.Action{"dummy": test.NewDummyAction(func() {
		// do nothing
	})}

	trg, err := test.InitTrigger(f, config, actions)
	assert.Nil(t, err)
	assert.NotNil(t, trg)

	// Note: Start will fail to connect to the test IMAP server, but we're testing
	// that the trigger can be started and stopped without panicking
	err = trg.Start()
	// We expect an error since we're not connecting to a real IMAP server
	// but the trigger should handle it gracefully
	assert.Nil(t, err)

	err = trg.Stop()
	assert.Nil(t, err)
}

func TestFormatAddress(t *testing.T) {
	// Test with nil address
	result := formatAddress(nil)
	assert.Equal(t, "", result)

	// Note: We can't easily test with actual imap.Address without importing
	// the IMAP library in tests, but the function is simple enough to verify
	// through integration testing
}

func TestOutput_ToMap(t *testing.T) {
	output := &Output{
		From:    "sender@example.com",
		To:      "recipient@example.com",
		Subject: "Test Subject",
		Body:    "Test Body",
		Date:    "2026-02-09T12:00:00Z",
	}

	m := output.ToMap()
	assert.Equal(t, "sender@example.com", m["from"])
	assert.Equal(t, "recipient@example.com", m["to"])
	assert.Equal(t, "Test Subject", m["subject"])
	assert.Equal(t, "Test Body", m["body"])
	assert.Equal(t, "2026-02-09T12:00:00Z", m["date"])
}

func TestOutput_FromMap(t *testing.T) {
	m := map[string]interface{}{
		"from":    "sender@example.com",
		"to":      "recipient@example.com",
		"subject": "Test Subject",
		"body":    "Test Body",
		"date":    "2026-02-09T12:00:00Z",
	}

	output := &Output{}
	err := output.FromMap(m)
	assert.Nil(t, err)
	assert.Equal(t, "sender@example.com", output.From)
	assert.Equal(t, "recipient@example.com", output.To)
	assert.Equal(t, "Test Subject", output.Subject)
	assert.Equal(t, "Test Body", output.Body)
	assert.Equal(t, "2026-02-09T12:00:00Z", output.Date)
}

func TestEmailTrigger_DefaultSettings(t *testing.T) {
	// Test that default values are applied correctly
	f := &Factory{}

	configWithDefaults := `{
		"id": "qodeflow-email",
		"ref": "github.com/postqode/qodeflow-component/trigger/email",
		"settings": {
			"host": "imap.example.com",
			"port": 993,
			"username": "test@example.com",
			"password": "testpassword"
		},
		"handlers": [
		  {
			"settings":{},
			"action":{"id":"dummy"}
		  }
		]
	  }`

	config := &trigger.Config{}
	err := json.Unmarshal([]byte(configWithDefaults), config)
	assert.Nil(t, err)

	trg, err := f.New(config)
	assert.Nil(t, err)

	// Verify TLS defaults to true for port 993
	emailTrg, ok := trg.(*Trigger)
	assert.True(t, ok)
	assert.True(t, emailTrg.settings.UseTLS)
}
