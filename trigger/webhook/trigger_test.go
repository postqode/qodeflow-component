package webhook

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/postqode/qodeflow-core/support/log"
	"github.com/postqode/qodeflow-core/trigger"
	"github.com/stretchr/testify/assert"
)

func TestWebhookTrigger(t *testing.T) {
	fixture := NewWebhookFixture(t, 8081)

	// Test case 1: POST request with JSON body
	t.Run("PostWithJSON", func(t *testing.T) {
		fixture.SetHandler("/test", "POST", func(ctx context.Context, triggerData interface{}) (map[string]interface{}, error) {
			return map[string]interface{}{
				"code": 200,
				"data": map[string]interface{}{"status": "ok"},
			}, nil
		})
	})

	// Test case 2: Path parameters
	t.Run("PathParams", func(t *testing.T) {
		fixture.SetHandler("/user/{id}", "GET", func(ctx context.Context, triggerData interface{}) (map[string]interface{}, error) {
			out := triggerData.(*Output)
			return map[string]interface{}{
				"code": 200,
				"data": map[string]interface{}{"userId": out.PathParams["id"]},
			}, nil
		})
	})

	err := fixture.Start()
	assert.Nil(t, err)
	defer fixture.Stop()

	// Wait for server to start
	time.Sleep(100 * time.Millisecond)

	t.Run("VerifyPostWithJSON", func(t *testing.T) {
		body := map[string]interface{}{"foo": "bar"}
		resp, err := fixture.Post("http://localhost:8081/test", body)
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)
		assert.Equal(t, "ok", result["status"])
	})

	t.Run("VerifyPathParams", func(t *testing.T) {
		resp, err := http.Get("http://localhost:8081/user/123")
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var result map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&result)
		assert.Equal(t, "123", result["userId"])
	})
}

// WebhookFixture is a test fixture for the Webhook trigger
type WebhookFixture struct {
	t        *testing.T
	port     int
	trig     trigger.Trigger
	handlers []*mockHandler
}

func NewWebhookFixture(t *testing.T, port int) *WebhookFixture {
	factory := &Factory{}
	config := &trigger.Config{
		Id: "test-webhook",
		Settings: map[string]interface{}{
			"port": port,
		},
	}
	trig, _ := factory.New(config)

	return &WebhookFixture{
		t:    t,
		port: port,
		trig: trig,
	}
}

func (f *WebhookFixture) Start() error {
	ctx := &mockInitContext{fixture: f}
	if err := f.trig.Initialize(ctx); err != nil {
		return err
	}
	return f.trig.Start()
}

func (f *WebhookFixture) Stop() error {
	return f.trig.Stop()
}

func (f *WebhookFixture) SetHandler(path, method string, handleFunc func(context.Context, interface{}) (map[string]interface{}, error)) {
	f.handlers = append(f.handlers, &mockHandler{
		path:       path,
		method:     method,
		handleFunc: handleFunc,
	})
}

func (f *WebhookFixture) Post(url string, body interface{}) (*http.Response, error) {
	jsonBody, _ := json.Marshal(body)
	return http.Post(url, "application/json", bytes.NewBuffer(jsonBody))
}

type mockInitContext struct {
	trigger.InitContext
	fixture *WebhookFixture
}

func (m *mockInitContext) Logger() log.Logger {
	return log.RootLogger()
}

func (m *mockInitContext) GetHandlers() []trigger.Handler {
	var handlers []trigger.Handler
	for _, h := range m.fixture.handlers {
		handlers = append(handlers, h)
	}
	return handlers
}

type mockHandler struct {
	trigger.Handler
	path       string
	method     string
	handleFunc func(context.Context, interface{}) (map[string]interface{}, error)
}

func (m *mockHandler) Settings() map[string]interface{} {
	return map[string]interface{}{
		"method": m.method,
		"path":   m.path,
	}
}

func (m *mockHandler) Handle(ctx context.Context, triggerData interface{}) (map[string]interface{}, error) {
	if m.handleFunc != nil {
		return m.handleFunc(ctx, triggerData)
	}
	return nil, nil
}
