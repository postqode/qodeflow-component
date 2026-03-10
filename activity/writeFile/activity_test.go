package writefile

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/postqode/qodeflow-core/activity"
	"github.com/postqode/qodeflow-core/support/test"
	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	ref := activity.GetRef(&Activity{})
	act := activity.Get(ref)

	assert.NotNil(t, act)
}

func TestWriteFileAndAppend(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "test.txt")

	act := &Activity{settings: &Settings{}}

	tc := test.NewActivityContext(activityMetadata)

	// Test WriteFile
	tc.SetInputObject(&Input{
		Method:   "WriteFile",
		FilePath: filePath,
		Content:  "Hello ",
	})

	done, err := act.Eval(tc)
	assert.True(t, done)
	assert.Nil(t, err)

	out := &Output{}
	err = tc.GetOutputObject(out)
	assert.Nil(t, err)
	assert.True(t, out.Success)

	data, err := os.ReadFile(filePath)
	assert.Nil(t, err)
	assert.Equal(t, "Hello ", string(data))

	// Test AppendFile
	tc.SetInputObject(&Input{
		Method:   "AppendFile",
		FilePath: filePath,
		Content:  "World!",
	})

	done, err = act.Eval(tc)
	assert.True(t, done)
	assert.Nil(t, err)

	err = tc.GetOutputObject(out)
	assert.Nil(t, err)
	assert.True(t, out.Success)

	data, err = os.ReadFile(filePath)
	assert.Nil(t, err)
	assert.Equal(t, "Hello World!", string(data))
}

func TestDeleteFile(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "delete_me.txt")

	err := os.WriteFile(filePath, []byte("delete this"), 0644)
	assert.Nil(t, err)

	act := &Activity{settings: &Settings{}}

	tc := test.NewActivityContext(activityMetadata)

	tc.SetInputObject(&Input{
		Method:   "DeleteFile",
		FilePath: filePath,
	})

	done, err := act.Eval(tc)
	assert.True(t, done)
	assert.Nil(t, err)

	out := &Output{}
	err = tc.GetOutputObject(out)
	assert.Nil(t, err)
	assert.True(t, out.Success)

	_, err = os.Stat(filePath)
	assert.True(t, os.IsNotExist(err))
}

func TestCreateDirectory(t *testing.T) {
	tempDir := t.TempDir()
	dirPath := filepath.Join(tempDir, "new_dir")

	act := &Activity{settings: &Settings{}}

	tc := test.NewActivityContext(activityMetadata)

	tc.SetInputObject(&Input{
		Method:   "CreateDirectory",
		FilePath: dirPath,
	})

	done, err := act.Eval(tc)
	assert.True(t, done)
	assert.Nil(t, err)

	out := &Output{}
	err = tc.GetOutputObject(out)
	assert.Nil(t, err)
	assert.True(t, out.Success)

	info, err := os.Stat(dirPath)
	assert.Nil(t, err)
	assert.True(t, info.IsDir())
}

// Integration tests — only run when env vars are provided.
//
// Run with a real file:
//
//	WRITE_FILE_PATH=/path/to/file.txt go test -v -run TestWriteFileActivity_Integration ./...
//
// Run with a real directory:
//
//	WRITE_DIR_PATH=/path/to/dir go test -v -run TestWriteFileActivity_Integration ./...
//
// Both at once:
//
//	WRITE_FILE_PATH=/path/to/file.txt WRITE_DIR_PATH=/path/to/dir go test -v -run TestWriteFileActivity_Integration ./...
func TestWriteFileActivity_Integration(t *testing.T) {
	filePath := os.Getenv("WRITE_FILE_PATH")
	dirPath := os.Getenv("WRITE_DIR_PATH")

	if filePath == "" && dirPath == "" {
		t.Skip("Skipping integration test: set WRITE_FILE_PATH and/or WRITE_DIR_PATH")
	}

	act := &Activity{settings: &Settings{}}

	// --- WriteFile ---
	if filePath != "" {
		t.Run("WriteFile", func(t *testing.T) {
			tc := test.NewActivityContext(activityMetadata)
			tc.SetInputObject(&Input{Method: "WriteFile", FilePath: filePath, Content: "Hello from integration test!\n"})

			done, err := act.Eval(tc)
			assert.True(t, done)
			assert.NoError(t, err)

			output := &Output{}
			tc.GetOutputObject(output)
			assert.True(t, output.Success, "WriteFile failed: %s", output.Error)

			result := output.Result.(map[string]any)
			fmt.Printf("[WriteFile] path=%s  size=%v bytes\n", result["filePath"], result["size"])
		})

		// --- AppendFile ---
		t.Run("AppendFile", func(t *testing.T) {
			tc := test.NewActivityContext(activityMetadata)
			tc.SetInputObject(&Input{Method: "AppendFile", FilePath: filePath, Content: "Appending more data!\n"})

			done, err := act.Eval(tc)
			assert.True(t, done)
			assert.NoError(t, err)

			output := &Output{}
			tc.GetOutputObject(output)
			assert.True(t, output.Success, "AppendFile failed: %s", output.Error)

			result := output.Result.(map[string]any)
			fmt.Printf("[AppendFile] path=%s  bytesWritten=%v\n", result["filePath"], result["bytesWritten"])
		})

		// --- DeleteFile ---
		t.Run("DeleteFile", func(t *testing.T) {
			tc := test.NewActivityContext(activityMetadata)
			tc.SetInputObject(&Input{Method: "DeleteFile", FilePath: filePath})

			done, err := act.Eval(tc)
			assert.True(t, done)
			assert.NoError(t, err)

			output := &Output{}
			tc.GetOutputObject(output)
			assert.True(t, output.Success, "DeleteFile failed: %s", output.Error)

			result := output.Result.(map[string]any)
			fmt.Printf("[DeleteFile] path=%s  deleted=%v\n", result["filePath"], result["deleted"])
		})
	}

	// --- CreateDirectory ---
	if dirPath != "" {
		t.Run("CreateDirectory", func(t *testing.T) {
			tc := test.NewActivityContext(activityMetadata)
			tc.SetInputObject(&Input{Method: "CreateDirectory", FilePath: dirPath})

			done, err := act.Eval(tc)
			assert.True(t, done)
			assert.NoError(t, err)

			output := &Output{}
			tc.GetOutputObject(output)
			assert.True(t, output.Success, "CreateDirectory failed: %s", output.Error)

			result := output.Result.(map[string]any)
			fmt.Printf("[CreateDirectory] path=%s  created=%v\n", result["dirPath"], result["created"])
		})
	}
}
