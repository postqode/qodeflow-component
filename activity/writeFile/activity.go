package writefile

import (
	"fmt"
	"os"
	"strings"

	"github.com/postqode/qodeflow-core/activity"
	"github.com/postqode/qodeflow-core/data/coerce"
	"github.com/postqode/qodeflow-core/data/metadata"
)

func init() {
	_ = activity.Register(&Activity{}, New)
}

// Activity represents the Write File activity
type Activity struct {
	settings *Settings
}

// New creates a new Write File activity
func New(ctx activity.InitContext) (activity.Activity, error) {
	s := &Settings{}
	err := metadata.MapToStruct(ctx.Settings(), s, true)
	if err != nil {
		return nil, err
	}

	return &Activity{
		settings: s,
	}, nil
}

// Metadata returns the metadata for the Write File activity
func (act *Activity) Metadata() *activity.Metadata {
	return activityMetadata
}

// Eval executes the Write File logic
func (act *Activity) Eval(ctx activity.Context) (done bool, err error) {
	input := &Input{}
	err = ctx.GetInputObject(input)
	if err != nil {
		return true, err
	}

	if input.FilePath == "" {
		return true, fmt.Errorf("filePath is required")
	}

	output := &Output{
		Success: true,
	}

	switch strings.ToLower(input.Method) {
	case "writefile":
		err = act.writeFile(ctx, input, output)
	case "appendfile":
		err = act.appendFile(ctx, input, output)
	case "deletefile":
		err = act.deleteFile(ctx, input, output)
	case "createdirectory":
		err = act.createDirectory(ctx, input, output)
	default:
		err = fmt.Errorf("unsupported method: %s", input.Method)
	}

	if err != nil {
		output.Success = false
		output.Error = err.Error()
		ctx.Logger().Errorf("Write File operation %s failed: %v", input.Method, err)
	}

	err = ctx.SetOutputObject(output)
	if err != nil {
		return true, err
	}

	return true, nil
}

func (act *Activity) writeFile(ctx activity.Context, input *Input, output *Output) error {
	contentStr, _ := coerce.ToString(input.Content)
	err := os.WriteFile(input.FilePath, []byte(contentStr), 0644)
	if err != nil {
		return fmt.Errorf("failed to write file %q: %w", input.FilePath, err)
	}

	output.Result = map[string]any{
		"filePath": input.FilePath,
		"size":     len(contentStr),
	}
	return nil
}

func (act *Activity) appendFile(ctx activity.Context, input *Input, output *Output) error {
	contentStr, _ := coerce.ToString(input.Content)

	f, err := os.OpenFile(input.FilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open file %q for appending: %w", input.FilePath, err)
	}
	defer f.Close()

	n, err := f.WriteString(contentStr)
	if err != nil {
		return fmt.Errorf("failed to append to file %q: %w", input.FilePath, err)
	}

	output.Result = map[string]any{
		"filePath":     input.FilePath,
		"bytesWritten": n,
	}
	return nil
}

func (act *Activity) deleteFile(ctx activity.Context, input *Input, output *Output) error {
	err := os.Remove(input.FilePath)
	if err != nil {
		return fmt.Errorf("failed to delete file %q: %w", input.FilePath, err)
	}

	output.Result = map[string]any{
		"filePath": input.FilePath,
		"deleted":  true,
	}
	return nil
}

func (act *Activity) createDirectory(ctx activity.Context, input *Input, output *Output) error {
	err := os.MkdirAll(input.FilePath, 0755)
	if err != nil {
		return fmt.Errorf("failed to create directory %q: %w", input.FilePath, err)
	}

	output.Result = map[string]any{
		"dirPath": input.FilePath,
		"created": true,
	}
	return nil
}
