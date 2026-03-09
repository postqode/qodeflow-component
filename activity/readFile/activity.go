package readfile

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/postqode/qodeflow-core/activity"
	"github.com/postqode/qodeflow-core/data/metadata"
)

func init() {
	_ = activity.Register(&Activity{}, New)
}

// Activity represents the Read File activity
type Activity struct {
	settings *Settings
}

// New creates a new Read File activity
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

// Metadata returns the metadata for the Read File activity
func (act *Activity) Metadata() *activity.Metadata {
	return activityMetadata
}

// Eval executes the Read File logic
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
	case "readfile":
		err = act.readFile(ctx, input, output)
	case "fileexists":
		err = act.fileExists(ctx, input, output)
	case "listdirectory":
		err = act.listDirectory(ctx, input, output)
	case "getfileinfo":
		err = act.getFileInfo(ctx, input, output)
	default:
		err = fmt.Errorf("unsupported method: %s", input.Method)
	}

	if err != nil {
		output.Success = false
		output.Error = err.Error()
		ctx.Logger().Errorf("Read File operation %s failed: %v", input.Method, err)
	}

	err = ctx.SetOutputObject(output)
	if err != nil {
		return true, err
	}

	return true, nil
}

func (act *Activity) readFile(ctx activity.Context, input *Input, output *Output) error {
	data, err := os.ReadFile(input.FilePath)
	if err != nil {
		return fmt.Errorf("failed to read file %q: %w", input.FilePath, err)
	}

	output.Result = map[string]any{
		"content":  string(data),
		"filePath": input.FilePath,
		"size":     len(data),
	}
	return nil
}

func (act *Activity) fileExists(ctx activity.Context, input *Input, output *Output) error {
	_, err := os.Stat(input.FilePath)
	if err != nil {
		if os.IsNotExist(err) {
			output.Result = map[string]any{
				"exists":   false,
				"filePath": input.FilePath,
			}
			return nil
		}
		return fmt.Errorf("failed to check file %q: %w", input.FilePath, err)
	}

	output.Result = map[string]any{
		"exists":   true,
		"filePath": input.FilePath,
	}
	return nil
}

func (act *Activity) listDirectory(ctx activity.Context, input *Input, output *Output) error {
	entries, err := os.ReadDir(input.FilePath)
	if err != nil {
		return fmt.Errorf("failed to list directory %q: %w", input.FilePath, err)
	}

	var items []map[string]any
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			continue
		}
		items = append(items, map[string]any{
			"name":    entry.Name(),
			"isDir":   entry.IsDir(),
			"size":    info.Size(),
			"modTime": info.ModTime().Format(time.RFC3339),
		})
	}

	output.Result = map[string]any{
		"dirPath": input.FilePath,
		"entries": items,
		"count":   len(items),
	}
	return nil
}

func (act *Activity) getFileInfo(ctx activity.Context, input *Input, output *Output) error {
	info, err := os.Stat(input.FilePath)
	if err != nil {
		return fmt.Errorf("failed to get file info for %q: %w", input.FilePath, err)
	}

	output.Result = map[string]any{
		"name":    info.Name(),
		"size":    info.Size(),
		"isDir":   info.IsDir(),
		"modTime": info.ModTime().Format(time.RFC3339),
		"mode":    info.Mode().String(),
	}
	return nil
}
