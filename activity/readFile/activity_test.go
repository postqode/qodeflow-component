package readfile

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/postqode/qodeflow-core/support/test"
	"github.com/stretchr/testify/assert"
)

func TestReadFileActivity_ReadFile(t *testing.T) {
	// Create a temporary file
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.txt")
	content := "Hello, QodeFlow!"
	err := os.WriteFile(tmpFile, []byte(content), 0644)
	assert.NoError(t, err)

	act := &Activity{settings: &Settings{}}

	tc := test.NewActivityContext(activityMetadata)
	input := &Input{
		Method:   "ReadFile",
		FilePath: tmpFile,
	}
	tc.SetInputObject(input)

	done, err := act.Eval(tc)
	assert.True(t, done)
	assert.NoError(t, err)

	output := &Output{}
	tc.GetOutputObject(output)
	assert.True(t, output.Success)
	result := output.Result.(map[string]any)
	assert.Equal(t, content, result["content"])
	assert.Equal(t, tmpFile, result["filePath"])
	assert.Equal(t, len(content), result["size"])
}

func TestReadFileActivity_ReadFile_NonExistent(t *testing.T) {
	act := &Activity{settings: &Settings{}}

	tc := test.NewActivityContext(activityMetadata)
	input := &Input{
		Method:   "ReadFile",
		FilePath: "/nonexistent/path/file.txt",
	}
	tc.SetInputObject(input)

	done, err := act.Eval(tc)
	assert.True(t, done)
	assert.NoError(t, err)

	output := &Output{}
	tc.GetOutputObject(output)
	assert.False(t, output.Success)
	assert.NotEmpty(t, output.Error)
}

func TestReadFileActivity_FileExists_True(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "existing.txt")
	err := os.WriteFile(tmpFile, []byte("data"), 0644)
	assert.NoError(t, err)

	act := &Activity{settings: &Settings{}}

	tc := test.NewActivityContext(activityMetadata)
	input := &Input{
		Method:   "FileExists",
		FilePath: tmpFile,
	}
	tc.SetInputObject(input)

	done, err := act.Eval(tc)
	assert.True(t, done)
	assert.NoError(t, err)

	output := &Output{}
	tc.GetOutputObject(output)
	assert.True(t, output.Success)
	result := output.Result.(map[string]any)
	assert.Equal(t, true, result["exists"])
}

func TestReadFileActivity_FileExists_False(t *testing.T) {
	act := &Activity{settings: &Settings{}}

	tc := test.NewActivityContext(activityMetadata)
	input := &Input{
		Method:   "FileExists",
		FilePath: "/nonexistent/path/missing.txt",
	}
	tc.SetInputObject(input)

	done, err := act.Eval(tc)
	assert.True(t, done)
	assert.NoError(t, err)

	output := &Output{}
	tc.GetOutputObject(output)
	assert.True(t, output.Success)
	result := output.Result.(map[string]any)
	assert.Equal(t, false, result["exists"])
}

func TestReadFileActivity_ListDirectory(t *testing.T) {
	tmpDir := t.TempDir()
	// Create some files in temp dir
	err := os.WriteFile(filepath.Join(tmpDir, "file1.txt"), []byte("a"), 0644)
	assert.NoError(t, err)
	err = os.WriteFile(filepath.Join(tmpDir, "file2.txt"), []byte("b"), 0644)
	assert.NoError(t, err)
	err = os.Mkdir(filepath.Join(tmpDir, "subdir"), 0755)
	assert.NoError(t, err)

	act := &Activity{settings: &Settings{}}

	tc := test.NewActivityContext(activityMetadata)
	input := &Input{
		Method:   "ListDirectory",
		FilePath: tmpDir,
	}
	tc.SetInputObject(input)

	done, err := act.Eval(tc)
	assert.True(t, done)
	assert.NoError(t, err)

	output := &Output{}
	tc.GetOutputObject(output)
	assert.True(t, output.Success)
	result := output.Result.(map[string]any)
	assert.Equal(t, tmpDir, result["dirPath"])
	assert.Equal(t, 3, result["count"])
}

func TestReadFileActivity_GetFileInfo(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "info.txt")
	content := "file info test"
	err := os.WriteFile(tmpFile, []byte(content), 0644)
	assert.NoError(t, err)

	act := &Activity{settings: &Settings{}}

	tc := test.NewActivityContext(activityMetadata)
	input := &Input{
		Method:   "GetFileInfo",
		FilePath: tmpFile,
	}
	tc.SetInputObject(input)

	done, err := act.Eval(tc)
	assert.True(t, done)
	assert.NoError(t, err)

	output := &Output{}
	tc.GetOutputObject(output)
	assert.True(t, output.Success)
	result := output.Result.(map[string]any)
	assert.Equal(t, "info.txt", result["name"])
	assert.Equal(t, int64(len(content)), result["size"])
	assert.Equal(t, false, result["isDir"])
}

func TestReadFileActivity_UnsupportedMethod(t *testing.T) {
	act := &Activity{settings: &Settings{}}

	tc := test.NewActivityContext(activityMetadata)
	input := &Input{
		Method:   "UnknownOp",
		FilePath: "/some/path",
	}
	tc.SetInputObject(input)

	done, err := act.Eval(tc)
	assert.True(t, done)
	assert.NoError(t, err)

	output := &Output{}
	tc.GetOutputObject(output)
	assert.False(t, output.Success)
	assert.Contains(t, output.Error, "unsupported method")
}

func TestReadFileActivity_MissingFilePath(t *testing.T) {
	act := &Activity{settings: &Settings{}}

	tc := test.NewActivityContext(activityMetadata)
	input := &Input{
		Method:   "ReadFile",
		FilePath: "",
	}
	tc.SetInputObject(input)

	done, err := act.Eval(tc)
	assert.True(t, done)
	assert.Error(t, err)
}

// Integration tests — only run when env vars are provided.
//
// Run with a real file:
//
//	READ_FILE_PATH=/path/to/file.txt go test -v -run TestReadFileActivity_Integration ./...
//
// Run with a real directory:
//
//	READ_DIR_PATH=/path/to/dir go test -v -run TestReadFileActivity_Integration ./...
//
// Both at once:
//
//	READ_FILE_PATH=/path/to/file.txt READ_DIR_PATH=/path/to/dir go test -v -run TestReadFileActivity_Integration ./...
func TestReadFileActivity_Integration(t *testing.T) {
	filePath := os.Getenv("READ_FILE_PATH")
	dirPath := os.Getenv("READ_DIR_PATH")

	if filePath == "" && dirPath == "" {
		t.Skip("Skipping integration test: set READ_FILE_PATH and/or READ_DIR_PATH")
	}

	act := &Activity{settings: &Settings{}}

	// --- ReadFile ---
	if filePath != "" {
		t.Run("ReadFile", func(t *testing.T) {
			tc := test.NewActivityContext(activityMetadata)
			tc.SetInputObject(&Input{Method: "ReadFile", FilePath: filePath})

			done, err := act.Eval(tc)
			assert.True(t, done)
			assert.NoError(t, err)

			output := &Output{}
			tc.GetOutputObject(output)
			assert.True(t, output.Success, "ReadFile failed: %s", output.Error)

			result := output.Result.(map[string]any)
			fmt.Printf("[ReadFile] path=%s  size=%v bytes\n", result["filePath"], result["size"])
			fmt.Printf("[ReadFile] content (first 200 chars):\n%.200s\n", result["content"])
		})

		// --- FileExists ---
		t.Run("FileExists", func(t *testing.T) {
			tc := test.NewActivityContext(activityMetadata)
			tc.SetInputObject(&Input{Method: "FileExists", FilePath: filePath})

			done, err := act.Eval(tc)
			assert.True(t, done)
			assert.NoError(t, err)

			output := &Output{}
			tc.GetOutputObject(output)
			assert.True(t, output.Success)
			result := output.Result.(map[string]any)
			assert.Equal(t, true, result["exists"])
			fmt.Printf("[FileExists] %s → exists=%v\n", result["filePath"], result["exists"])
		})

		// --- GetFileInfo ---
		t.Run("GetFileInfo", func(t *testing.T) {
			tc := test.NewActivityContext(activityMetadata)
			tc.SetInputObject(&Input{Method: "GetFileInfo", FilePath: filePath})

			done, err := act.Eval(tc)
			assert.True(t, done)
			assert.NoError(t, err)

			output := &Output{}
			tc.GetOutputObject(output)
			assert.True(t, output.Success, "GetFileInfo failed: %s", output.Error)

			result := output.Result.(map[string]any)
			fmt.Printf("[GetFileInfo] name=%s  size=%v  isDir=%v  modTime=%s\n",
				result["name"], result["size"], result["isDir"], result["modTime"])
		})
	}

	// --- ListDirectory ---
	if dirPath != "" {
		t.Run("ListDirectory", func(t *testing.T) {
			tc := test.NewActivityContext(activityMetadata)
			tc.SetInputObject(&Input{Method: "ListDirectory", FilePath: dirPath})

			done, err := act.Eval(tc)
			assert.True(t, done)
			assert.NoError(t, err)

			output := &Output{}
			tc.GetOutputObject(output)
			assert.True(t, output.Success, "ListDirectory failed: %s", output.Error)

			result := output.Result.(map[string]any)
			fmt.Printf("[ListDirectory] path=%s  count=%v\n", result["dirPath"], result["count"])
			if entries, ok := result["entries"].([]map[string]any); ok {
				for _, e := range entries {
					fmt.Printf("  - %s (isDir=%v, size=%v)\n", e["name"], e["isDir"], e["size"])
				}
			}
		})
	}
}
