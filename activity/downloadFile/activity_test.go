package downloadfile

import (
	"net/http"
	"net/http/httptest"
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

func TestDownloadFile_Success(t *testing.T) {
	// Create a mock HTTP server
	mockContent := "hello world"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "/testfile.txt", r.URL.Path)
		assert.Equal(t, "custom-value", r.Header.Get("X-Custom-Header"))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mockContent))
	}))
	defer server.Close()

	// Temporary directory for the download
	tempDir, err := os.MkdirTemp("", "downloadfile-test")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	destPath := filepath.Join(tempDir, "downloaded.txt")

	act := &Activity{
		settings: &Settings{},
	}
	tc := test.NewActivityContext(act.Metadata())
	input := &Input{
		URL:         server.URL + "/testfile.txt",
		Destination: destPath,
		Method:      "GET",
		Headers:     map[string]any{"X-Custom-Header": "custom-value"},
		Timeout:     5,
	}

	err = tc.SetInputObject(input)
	assert.NoError(t, err)

	done, err := act.Eval(tc)
	assert.True(t, done)
	assert.NoError(t, err)

	output := &Output{}
	err = tc.GetOutputObject(output)
	assert.NoError(t, err)

	assert.True(t, output.Success)
	assert.Empty(t, output.Error)

	resultMap, ok := output.Result.(map[string]any)
	assert.True(t, ok)
	assert.Equal(t, destPath, resultMap["filePath"])
	assert.Equal(t, int64(len(mockContent)), resultMap["size"])

	// Verify the file was created and contains the correct data
	content, err := os.ReadFile(destPath)
	assert.NoError(t, err)
	assert.Equal(t, mockContent, string(content))
}

func TestDownloadFile_ToDirectory(t *testing.T) {
	mockContent := "hello directory"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(mockContent))
	}))
	defer server.Close()

	tempDir, err := os.MkdirTemp("", "downloadfile-test-dir")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	act := &Activity{settings: &Settings{}}
	tc := test.NewActivityContext(act.Metadata())
	input := &Input{
		URL:         server.URL + "/myfile.txt",
		Destination: tempDir, // pointing to directory
	}

	tc.SetInputObject(input)
	_, err = act.Eval(tc)
	assert.NoError(t, err)

	output := &Output{}
	tc.GetOutputObject(output)
	assert.True(t, output.Success)

	expectedPath := filepath.Join(tempDir, "myfile.txt")
	resultMap := output.Result.(map[string]any)
	assert.Equal(t, expectedPath, resultMap["filePath"])

	content, err := os.ReadFile(expectedPath)
	assert.NoError(t, err)
	assert.Equal(t, mockContent, string(content))
}

func TestDownloadFile_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	tempDir, err := os.MkdirTemp("", "downloadfile-test-notfound")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)

	destPath := filepath.Join(tempDir, "notfound.txt")

	act := &Activity{settings: &Settings{}}
	tc := test.NewActivityContext(act.Metadata())
	input := &Input{
		URL:         server.URL + "/missing",
		Destination: destPath,
	}

	tc.SetInputObject(input)
	_, err = act.Eval(tc)
	assert.NoError(t, err)

	output := &Output{}
	tc.GetOutputObject(output)
	assert.False(t, output.Success)
	assert.Contains(t, output.Error, "received non-success status code: 404")

	// Ensure file was not created
	_, err = os.Stat(destPath)
	assert.True(t, os.IsNotExist(err))
}
