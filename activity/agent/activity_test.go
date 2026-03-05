package agent

import (
	"context"
	"testing"

	"github.com/cloudwego/eino/schema"
	"github.com/postqode/qodeflow-core/activity"
	"github.com/postqode/qodeflow-core/support/test"
	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	ref := activity.GetRef(&Activity{})
	act := activity.Get(ref)
	assert.NotNil(t, act)
}

func TestEval(t *testing.T) {
	const wantResponse = "Hello from the mock agent"

	// Build an Activity with a mock generate function – no real API call needed.
	act := &Activity{
		settings: &Settings{
			Provider:     "openrouter",
			Model:        "anthropic/claude-3-5-sonnet",
			SystemPrompt: "You are a helpful assistant.",
		},
		generateFn: func(_ context.Context, msgs []*schema.Message) (*schema.Message, error) {
			return schema.AssistantMessage(wantResponse, nil), nil
		},
	}

	tc := test.NewActivityContext(act.Metadata())

	err := tc.SetInputObject(&Input{Message: "Hi"})
	assert.NoError(t, err)

	done, err := act.Eval(tc)
	assert.True(t, done)
	assert.NoError(t, err)

	output := &Output{}
	err = tc.GetOutputObject(output)
	assert.NoError(t, err)
	assert.Equal(t, wantResponse, output.Response)
}

func TestEval_WithResponseStructure(t *testing.T) {
	const wantResponse = `{"answer":"42"}`

	act := &Activity{
		settings: &Settings{
			Provider:     "postqode",
			Model:        "gemini-2.0-flash",
			SystemPrompt: "You are a helpful assistant.",
			ResponseStructure: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"answer": map[string]any{
						"type": "string",
					},
				},
			},
		},
		generateFn: func(_ context.Context, msgs []*schema.Message) (*schema.Message, error) {
			return schema.AssistantMessage(wantResponse, nil), nil
		},
	}

	tc := test.NewActivityContext(act.Metadata())

	err := tc.SetInputObject(&Input{Message: "What is the answer?"})
	assert.NoError(t, err)

	done, err := act.Eval(tc)
	assert.True(t, done)
	assert.NoError(t, err)

	output := &Output{}
	err = tc.GetOutputObject(output)
	assert.NoError(t, err)
	assert.Equal(t, wantResponse, output.Response)
}

func TestEval_PostqodeProvider(t *testing.T) {
	const wantResponse = "Response via postqode provider"

	act := &Activity{
		settings: &Settings{
			Provider: "postqode",
			Model:    "gemini-2.0-flash",
		},
		generateFn: func(_ context.Context, msgs []*schema.Message) (*schema.Message, error) {
			return schema.AssistantMessage(wantResponse, nil), nil
		},
	}

	tc := test.NewActivityContext(act.Metadata())

	err := tc.SetInputObject(&Input{Message: "Hello"})
	assert.NoError(t, err)

	done, err := act.Eval(tc)
	assert.True(t, done)
	assert.NoError(t, err)

	output := &Output{}
	err = tc.GetOutputObject(output)
	assert.NoError(t, err)
	assert.Equal(t, wantResponse, output.Response)
}

func TestEval_EmptyMessage(t *testing.T) {
	act := &Activity{
		settings: &Settings{Provider: "openrouter", Model: "anthropic/claude-3-5-sonnet"},
		generateFn: func(_ context.Context, _ []*schema.Message) (*schema.Message, error) {
			return schema.AssistantMessage("should not reach here", nil), nil
		},
	}

	tc := test.NewActivityContext(act.Metadata())
	err := tc.SetInputObject(&Input{Message: ""})
	assert.NoError(t, err)

	done, err := act.Eval(tc)
	assert.True(t, done)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "message")
}
