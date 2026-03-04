package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	einoopenrouter "github.com/cloudwego/eino-ext/components/model/openrouter"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/flow/agent/react"
	"github.com/cloudwego/eino/schema"
	"github.com/eino-contrib/jsonschema"

	"github.com/postqode/qodeflow-core/activity"
	"github.com/postqode/qodeflow-core/data/metadata"
)

func init() {
	_ = activity.Register(&Activity{}, New)
}

// generateFunc wraps the agent's generate call so it can be swapped in tests.
type generateFunc func(ctx context.Context, msgs []*schema.Message) (*schema.Message, error)

// Activity is the LLM Agent activity powered by cloudwego/eino.
type Activity struct {
	settings   *Settings
	generateFn generateFunc
}

// New creates a per-configuration Activity instance. It resolves the model
// provider via eino-ext and builds an eino ReAct agent internally.
func New(ctx activity.InitContext) (activity.Activity, error) {
	s := &Settings{}
	if err := metadata.MapToStruct(ctx.Settings(), s, true); err != nil {
		return nil, err
	}

	if s.Provider == "" {
		return nil, fmt.Errorf("agent activity: 'provider' setting is required")
	}
	if s.Model == "" {
		return nil, fmt.Errorf("agent activity: 'model' setting is required")
	}

	bgCtx := context.Background()

	chatModel, err := newChatModel(bgCtx, s)
	if err != nil {
		return nil, fmt.Errorf("agent activity: failed to initialise model for provider %q: %w", s.Provider, err)
	}

	agentCfg := &react.AgentConfig{
		ToolCallingModel: chatModel,
	}

	// Prepend the system prompt on every turn when configured.
	if s.SystemPrompt != "" {
		prompt := s.SystemPrompt
		agentCfg.MessageModifier = func(_ context.Context, msgs []*schema.Message) []*schema.Message {
			if len(msgs) > 0 && msgs[0].Role == schema.System {
				return msgs // already has a system message; leave untouched
			}
			return append([]*schema.Message{schema.SystemMessage(prompt)}, msgs...)
		}
	}

	llmAgent, err := react.NewAgent(bgCtx, agentCfg)
	if err != nil {
		return nil, fmt.Errorf("agent activity: failed to create react agent: %w", err)
	}

	return &Activity{
		settings: s,
		generateFn: func(ctx context.Context, msgs []*schema.Message) (*schema.Message, error) {
			return llmAgent.Generate(ctx, msgs)
		},
	}, nil
}

// newChatModel constructs the appropriate eino-ext model for the given provider.
func newChatModel(ctx context.Context, s *Settings) (model.ToolCallingChatModel, error) {

	openrouterConfig := &einoopenrouter.Config{
		Model:  s.Model,
		APIKey: s.APIKey,
	}

	if s.ResponseStructure != "" {
		var jsSchema *jsonschema.Schema
		if err := json.Unmarshal([]byte(s.ResponseStructure), &jsSchema); err != nil {
			return nil, fmt.Errorf("agent activity: invalid response structure JSON: %w", err)
		}

		openrouterConfig.ResponseFormat = &einoopenrouter.ChatCompletionResponseFormat{
			Type: einoopenrouter.ChatCompletionResponseFormatTypeJSONSchema,
			JSONSchema: &einoopenrouter.ChatCompletionResponseFormatJSONSchema{
				Name:        "ResponseFormat",
				Description: "Response output structure",
				JSONSchema:  jsSchema,
				Strict:      true,
			},
		}
	}

	switch strings.ToLower(s.Provider) {
	case "openrouter":
		return einoopenrouter.NewChatModel(ctx, openrouterConfig)

	case "postqode":
		openrouterConfig.BaseURL = "https://api.postqode.ai/gateway/v1"
		return einoopenrouter.NewChatModel(ctx, openrouterConfig)

	default:
		return nil, fmt.Errorf("unsupported provider %q (supported: openai, openrouter, anthropic)", s.Provider)
	}
}

// Metadata returns the activity's metadata.
func (a *Activity) Metadata() *activity.Metadata {
	return activityMd
}

// Eval implements activity.Activity – sends the user message to the LLM agent
// and writes the text response to the output.
func (a *Activity) Eval(ctx activity.Context) (done bool, err error) {
	input := &Input{}
	if err = ctx.GetInputObject(input); err != nil {
		return true, err
	}

	if input.Message == "" {
		return true, fmt.Errorf("agent activity: 'message' input cannot be empty")
	}

	msgs := []*schema.Message{
		schema.UserMessage(input.Message),
	}

	resp, err := a.generateFn(context.Background(), msgs)
	if err != nil {
		return true, fmt.Errorf("agent activity: generation failed: %w", err)
	}

	output := &Output{Response: resp.Content}
	if err = ctx.SetOutputObject(output); err != nil {
		return true, err
	}

	return true, nil
}
