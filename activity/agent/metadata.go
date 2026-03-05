package agent

import (
	"github.com/postqode/qodeflow-core/activity"
	"github.com/postqode/qodeflow-core/data/coerce"
)

var activityMd = activity.ToMetadata(&Settings{}, &Input{}, &Output{})

// Settings holds the configuration for the LLM Agent activity.
type Settings struct {
	// Provider is the LLM provider name (e.g. openai, openrouter, anthropic)
	Provider string `md:"provider,required"`
	// Model is the model identifier (e.g. gpt-4o, claude-3-5-sonnet-20241022)
	Model string `md:"model,required"`
	// APIKey is the provider API key
	APIKey string `md:"apiKey"`
	// SystemPrompt is an optional system-level instruction for the agent
	SystemPrompt string `md:"systemPrompt"`
	// ResponseStructure is an optional jsonSchema, which specifies the response structure which agent should respond with
	ResponseStructure any `md:"responseStructure"`
}

// Input holds the runtime input for the LLM Agent activity.
type Input struct {
	// Message is the user message sent to the agent
	Message string `md:"message,required"`
}

func (r *Input) FromMap(values map[string]any) error {
	strVal, _ := coerce.ToString(values["message"])
	r.Message = strVal
	return nil
}

func (r *Input) ToMap() map[string]any {
	return map[string]any{
		"message": r.Message,
	}
}

// Output holds the result produced by the LLM Agent activity.
type Output struct {
	// Response is the text response from the agent
	Response string `md:"response"`
}

func (o *Output) FromMap(values map[string]any) error {
	strVal, _ := coerce.ToString(values["response"])
	o.Response = strVal
	return nil
}

func (o *Output) ToMap() map[string]any {
	return map[string]any{
		"response": o.Response,
	}
}
