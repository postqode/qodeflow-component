package email

// Settings for the trigger (IMAP connection settings)
type Settings struct {
	Host     string `md:"host,required"`     // The IMAP server host
	Port     int    `md:"port,required"`     // The IMAP server port
	Username string `md:"username,required"` // The email account username
	Password string `md:"password,required"` // The email account password
	UseTLS   bool   `md:"useTLS"`            // Use TLS for the IMAP connection
}

// HandlerSettings for each handler
type HandlerSettings struct {
	Folder       string `md:"folder"`       // The mailbox folder to monitor (default: INBOX)
	PollInterval string `md:"pollInterval"` // The polling interval (e.g., 30s, 1m, 5m)
}

type Attachment struct {
	Filename    string `md:"filename"`
	ContentType string `md:"contentType"`
	Content     []byte `md:"content"`
}

// Output represents the email data passed to the handler
type Output struct {
	From        string       `md:"from"`        // The sender email address
	To          string       `md:"to"`          // The recipient email address(es)
	Subject     string       `md:"subject"`     // The email subject
	Body        string       `md:"body"`        // The email body content
	Date        string       `md:"date"`        // The date the email was received
	Attachments []Attachment `md:"attachments"` // The attachment the email received
}

func (o *Output) ToMap() map[string]any {
	return map[string]any{
		"from":        o.From,
		"to":          o.To,
		"subject":     o.Subject,
		"body":        o.Body,
		"date":        o.Date,
		"attachments": o.Attachments,
	}
}

func (o *Output) FromMap(values map[string]any) error {
	if v, ok := values["from"]; ok {
		o.From, _ = v.(string)
	}
	if v, ok := values["to"]; ok {
		o.To, _ = v.(string)
	}
	if v, ok := values["subject"]; ok {
		o.Subject, _ = v.(string)
	}
	if v, ok := values["body"]; ok {
		o.Body, _ = v.(string)
	}
	if v, ok := values["date"]; ok {
		o.Date, _ = v.(string)
	}
	if v, ok := values["attachments"]; ok {
		if attachments, ok := v.([]Attachment); ok {
			o.Attachments = attachments
		}
	}
	return nil
}
