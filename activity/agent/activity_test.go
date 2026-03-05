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

func TestEval_WithMappingResponseStructure(t *testing.T) {
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

func TestEval_WithStringResponseStructure(t *testing.T) {
	const wantResponse = `{"answer":"42"}`
	const jsonSchema = `{"type":"object","properties":{"answer":{"type":"string"}}}`

	act := &Activity{
		settings: &Settings{
			Provider:          "postqode",
			Model:             "gemini-2.0-flash",
			ResponseStructure: jsonSchema,
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

func TestEval_WithComplexMappingResponseStructure(t *testing.T) {
	const wantResponse = `{"person":{"firstName":"John","lastName":"Doe"},"address":{"street":"123 Main St","city":"Springfield","zipCode":"62701"}}`

	// This structure mirrors the real-world usage from the user's app JSON,
	// with address correctly typed as "object" (not "boolean").
	act := &Activity{
		settings: &Settings{
			Provider: "postqode",
			Model:    "gemini-2.0-flash",
			ResponseStructure: map[string]any{
				"mapping": map[string]any{
					"type": "object",
					"properties": map[string]any{
						"person": map[string]any{
							"type":        "object",
							"description": "Personal information",
							"properties": map[string]any{
								"firstName":  map[string]any{"type": "string", "description": "First name of the person"},
								"lastName":   map[string]any{"type": "string", "description": "Last name of the person"},
								"age":        map[string]any{"type": "number", "description": "Age in years"},
								"isEmployed": map[string]any{"type": "boolean", "description": "Whether the person is currently employed"},
							},
							"required": []any{"firstName", "lastName"},
						},
						"address": map[string]any{
							"type":        "object",
							"description": "Address information",
							"properties": map[string]any{
								"street":  map[string]any{"type": "string", "description": "Street address"},
								"city":    map[string]any{"type": "string", "description": "City name"},
								"zipCode": map[string]any{"type": "string", "description": "Postal/ZIP code"},
							},
						},
						"hobbies": map[string]any{
							"type":        "array",
							"description": "List of hobbies",
							"items": map[string]any{
								"type": "object",
								"properties": map[string]any{
									"name":            map[string]any{"type": "string", "description": "Name of the hobby"},
									"yearsExperience": map[string]any{"type": "number", "description": "Years of experience"},
								},
							},
						},
					},
					"required": []any{"person"},
				},
			},
		},
		generateFn: func(_ context.Context, msgs []*schema.Message) (*schema.Message, error) {
			return schema.AssistantMessage(wantResponse, nil), nil
		},
	}

	tc := test.NewActivityContext(act.Metadata())
	err := tc.SetInputObject(&Input{Message: "Give a sample person information"})
	assert.NoError(t, err)

	done, err := act.Eval(tc)
	assert.True(t, done)
	assert.NoError(t, err)

	output := &Output{}
	err = tc.GetOutputObject(output)
	assert.NoError(t, err)
	assert.Equal(t, wantResponse, output.Response)
}

func TestEval_WithSchemaKeyword(t *testing.T) {
	const wantResponse = `{"person":{"firstName":"John","lastName":"Doe"}}`

	// JSON Schema with $schema keyword — passed as a map (as it would arrive from a config file)
	act := &Activity{
		settings: &Settings{
			Provider: "postqode",
			Model:    "gemini-2.0-flash",
			ResponseStructure: map[string]any{
				"$schema": "https://json-schema.org/draft-07/schema",
				"type":    "object",
				"properties": map[string]any{
					"person": map[string]any{
						"type":        "object",
						"description": "Personal information",
						"properties": map[string]any{
							"firstName":  map[string]any{"type": "string", "description": "First name of the person"},
							"lastName":   map[string]any{"type": "string", "description": "Last name of the person"},
							"age":        map[string]any{"type": "number", "description": "Age in years"},
							"isEmployed": map[string]any{"type": "boolean", "description": "Whether the person is currently employed"},
						},
						"required": []any{"firstName", "lastName"},
					},
					"address": map[string]any{
						"type":        "object",
						"description": "Address information",
						"properties": map[string]any{
							"street":  map[string]any{"type": "string", "description": "Street address"},
							"city":    map[string]any{"type": "string", "description": "City name"},
							"zipCode": map[string]any{"type": "string", "description": "Postal/ZIP code"},
						},
					},
					"hobbies": map[string]any{
						"type":        "array",
						"description": "List of hobbies",
						"items": map[string]any{
							"type": "object",
							"properties": map[string]any{
								"name":            map[string]any{"type": "string", "description": "Name of the hobby"},
								"yearsExperience": map[string]any{"type": "number", "description": "Years of experience"},
							},
						},
					},
				},
				"required": []any{"person"},
			},
		},
		generateFn: func(_ context.Context, msgs []*schema.Message) (*schema.Message, error) {
			return schema.AssistantMessage(wantResponse, nil), nil
		},
	}

	tc := test.NewActivityContext(act.Metadata())
	err := tc.SetInputObject(&Input{Message: "Give a sample person information"})
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
