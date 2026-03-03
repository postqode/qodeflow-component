package email

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-message/mail"
	"github.com/postqode/qodeflow-core/data/metadata"
	"github.com/postqode/qodeflow-core/support/log"
	"github.com/postqode/qodeflow-core/trigger"
)

var triggerMd = trigger.NewMetadata(&Settings{}, &HandlerSettings{}, &Output{})

func init() {
	_ = trigger.Register(&Trigger{}, &Factory{})
}

// Factory is the email trigger factory
type Factory struct {
}

// Metadata implements trigger.Factory.Metadata
func (*Factory) Metadata() *trigger.Metadata {
	return triggerMd
}

// New implements trigger.Factory.New
func (*Factory) New(config *trigger.Config) (trigger.Trigger, error) {
	s := &Settings{}
	err := metadata.MapToStruct(config.Settings, s, true)
	if err != nil {
		return nil, err
	}

	// Default to TLS enabled
	if !s.UseTLS {
		// Check if the setting was explicitly provided; if port is 993, default to TLS
		if s.Port == 993 {
			s.UseTLS = true
		}
	}

	return &Trigger{settings: s}, nil
}

// Trigger is the email trigger
type Trigger struct {
	settings  *Settings
	handlers  []trigger.Handler
	logger    log.Logger
	stopChan  chan struct{}
	startTime time.Time
}

// Initialize implements trigger.Trigger.Initialize
func (t *Trigger) Initialize(ctx trigger.InitContext) error {
	t.handlers = ctx.GetHandlers()
	t.logger = ctx.Logger()
	return nil
}

// Start implements trigger.Trigger.Start
func (t *Trigger) Start() error {
	t.stopChan = make(chan struct{})
	t.startTime = time.Now()

	for _, handler := range t.handlers {
		hs := &HandlerSettings{}
		err := metadata.MapToStruct(handler.Settings(), hs, true)
		if err != nil {
			return err
		}

		// Set defaults
		if hs.Folder == "" {
			hs.Folder = "INBOX"
		}
		if hs.PollInterval == "" {
			hs.PollInterval = "1m"
		}

		pollDuration, err := time.ParseDuration(hs.PollInterval)
		if err != nil {
			return fmt.Errorf("unable to parse poll interval '%s': %s", hs.PollInterval, err.Error())
		}

		go t.pollMailbox(handler, hs, pollDuration)
	}

	return nil
}

// Stop implements trigger.Trigger.Stop
func (t *Trigger) Stop() error {
	if t.stopChan != nil {
		close(t.stopChan)
	}
	return nil
}

// pollMailbox periodically checks the mailbox for new (unseen) emails
func (t *Trigger) pollMailbox(handler trigger.Handler, hs *HandlerSettings, interval time.Duration) {
	t.logger.Infof("Starting email poller for folder '%s' with interval %s", hs.Folder, interval.String())

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	// Do an initial poll immediately
	t.checkMail(handler, hs)
	for {
		select {
		case <-ticker.C:
			t.checkMail(handler, hs)
		case <-t.stopChan:
			t.logger.Info("Stopping email poller")
			return
		}
	}
}

// checkMail connects to the IMAP server and fetches messages
func (t *Trigger) checkMail(handler trigger.Handler, hs *HandlerSettings) {
	t.logger.Debug("Checking for new emails...")

	c, err := t.connect()
	if err != nil {
		t.logger.Errorf("Failed to connect to IMAP server: %s", err.Error())
		return
	}
	defer func() {
		_ = c.Logout()
	}()

	// Login
	if err := c.Login(t.settings.Username, t.settings.Password); err != nil {
		t.logger.Errorf("Failed to login to IMAP server: %s", err.Error())
		return
	}

	// Select the mailbox folder
	_, err = c.Select(hs.Folder, false)
	if err != nil {
		t.logger.Errorf("Failed to select folder '%s': %s", hs.Folder, err.Error())
		return
	}

	// Search for unseen messages
	criteria := imap.NewSearchCriteria()
	criteria.WithoutFlags = []string{imap.SeenFlag}
	criteria.Since = t.startTime

	uids, err := c.Search(criteria)
	if err != nil {
		t.logger.Errorf("Failed to search for unseen messages: %s", err.Error())
		return
	}

	if len(uids) == 0 {
		t.logger.Debug("No new emails found")
		return
	}

	t.logger.Infof("Found %d new email(s)", len(uids))

	seqSet := new(imap.SeqSet)
	seqSet.AddNum(uids...)

	// Fetch messages with envelope and body
	section := &imap.BodySectionName{}
	items := []imap.FetchItem{imap.FetchEnvelope, section.FetchItem()}

	messages := make(chan *imap.Message, len(uids))
	done := make(chan error, 1)
	go func() {
		done <- c.Fetch(seqSet, items, messages)
	}()

	for msg := range messages {
		t.processMessage(handler, msg, section)
	}

	if err := <-done; err != nil {
		t.logger.Errorf("Error fetching messages: %s", err.Error())
	}
}

// connect establishes a connection to the IMAP server
func (t *Trigger) connect() (*client.Client, error) {
	addr := fmt.Sprintf("%s:%d", t.settings.Host, t.settings.Port)

	if t.settings.UseTLS {
		tlsConfig := &tls.Config{
			ServerName: t.settings.Host,
		}
		return client.DialTLS(addr, tlsConfig)
	}

	return client.Dial(addr)
}

// processMessage extracts email data and triggers the handler
func (t *Trigger) processMessage(handler trigger.Handler, msg *imap.Message, section *imap.BodySectionName) {
	if msg == nil || msg.Envelope == nil {
		return
	}

	envelope := msg.Envelope

	// Extract sender
	from := ""
	if len(envelope.From) > 0 {
		from = formatAddress(envelope.From[0])
	}

	// Extract recipients
	toAddrs := make([]string, 0, len(envelope.To))
	for _, addr := range envelope.To {
		toAddrs = append(toAddrs, formatAddress(addr))
	}
	to := strings.Join(toAddrs, ", ")

	// Extract body
	body := ""
	var attachments []Attachment
	if msg.Body != nil {
		r := msg.GetBody(section)
		if r != nil {
			mr, err := mail.CreateReader(r)
			if err == nil {
				for {
					p, err := mr.NextPart()
					if err != nil {
						break
					}
					switch h := p.Header.(type) {
					case *mail.InlineHeader:
						b, err := io.ReadAll(p.Body)
						if err == nil {
							body = string(b)
						}
					case *mail.AttachmentHeader:
						filename, _ := h.Filename()
						contentType, _, _ := h.ContentType()
						b, err := io.ReadAll(p.Body)
						if err != nil {
							t.logger.Errorf("Failed to read attachment '%s': %s", filename, err.Error())
							continue
						}
						attachments = append(attachments, Attachment{
							Filename:    filename,
							ContentType: contentType,
							Content:     b,
						})
					}
				}
			}
		}
	}

	output := &Output{
		From:        from,
		To:          to,
		Subject:     envelope.Subject,
		Body:        body,
		Date:        envelope.Date.Format(time.RFC3339),
		Attachments: attachments,
	}

	t.logger.Debugf("Processing email from '%s' with subject '%s'", from, envelope.Subject)

	_, err := handler.Handle(context.Background(), output.ToMap())
	if err != nil {
		t.logger.Errorf("Error handling email trigger: %s", err.Error())
	}
}

// formatAddress formats an IMAP address to a string
func formatAddress(addr *imap.Address) string {
	if addr == nil {
		return ""
	}
	if addr.PersonalName != "" {
		return fmt.Sprintf("%s <%s@%s>", addr.PersonalName, addr.MailboxName, addr.HostName)
	}
	return fmt.Sprintf("%s@%s", addr.MailboxName, addr.HostName)
}
