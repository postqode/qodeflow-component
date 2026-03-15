package downloadfile

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/postqode/qodeflow-core/activity"
	"github.com/postqode/qodeflow-core/data/metadata"
)

func init() {
	_ = activity.Register(&Activity{}, New)
}

// Activity represents the Download File activity
type Activity struct {
	settings *Settings
}

// New creates a new Download File activity
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

// Metadata returns the metadata for the Download File activity
func (act *Activity) Metadata() *activity.Metadata {
	return activityMetadata
}

// Eval executes the Download File logic
func (act *Activity) Eval(ctx activity.Context) (done bool, err error) {
	input := &Input{}
	err = ctx.GetInputObject(input)
	if err != nil {
		return true, err
	}

	if input.URL == "" {
		return true, fmt.Errorf("url is required")
	}
	if input.Destination == "" {
		return true, fmt.Errorf("destination is required")
	}

	output := &Output{
		Success: true,
	}

	err = act.downloadFile(ctx, input, output)
	if err != nil {
		output.Success = false
		output.Error = err.Error()
		ctx.Logger().Errorf("Download File operation from %s failed: %v", input.URL, err)
	}

	err = ctx.SetOutputObject(output)
	if err != nil {
		return true, err
	}

	return true, nil
}

func (act *Activity) downloadFile(ctx activity.Context, input *Input, output *Output) error {
	method := "GET"
	if input.Method != "" {
		method = strings.ToUpper(input.Method)
	}

	req, err := http.NewRequest(method, input.URL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	if input.Headers != nil {
		for k, v := range input.Headers {
			req.Header.Set(k, fmt.Sprintf("%v", v))
		}
	}

	client := &http.Client{}
	if input.Timeout > 0 {
		client.Timeout = time.Duration(input.Timeout) * time.Second
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("received non-success status code: %d", resp.StatusCode)
	}

	destPath := input.Destination
	info, err := os.Stat(destPath)
	if err == nil && info.IsDir() {
		// If destination is a directory, append base name of URL
		parsedURL, err := url.Parse(input.URL)
		if err != nil {
			return fmt.Errorf("failed to parse URL to determine filename: %w", err)
		}
		filename := filepath.Base(parsedURL.Path)
		if filename == "." || filename == "/" {
			filename = "downloaded_file"
		}
		destPath = filepath.Join(destPath, filename)
	} else if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to check destination path: %w", err)
	}

	// Ensure the directory for the destination file exists
	dir := filepath.Dir(destPath)
	if dir != "" && dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create destination directory %s: %w", dir, err)
		}
	}

	var out *os.File
	fileExists := false
	if _, err := os.Stat(destPath); err == nil {
		fileExists = true
	}

	if input.Append && fileExists {
		// Append mode: open existing file for appending
		out, err = os.OpenFile(destPath, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return fmt.Errorf("failed to open file for append %q: %w", destPath, err)
		}
	} else {
		// Overwrite mode (default): create/truncate file
		// If file exists and append is false, this will overwrite (delete and recreate)
		out, err = os.Create(destPath)
		if err != nil {
			return fmt.Errorf("failed to create file %q: %w", destPath, err)
		}
	}
	defer out.Close()

	n, err := io.Copy(out, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to write to file %q: %w", destPath, err)
	}

	output.Result = map[string]any{
		"filePath": destPath,
		"size":     n,
	}

	return nil
}
