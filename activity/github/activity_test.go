package github

import (
	"os"
	"testing"

	gh "github.com/google/go-github/v60/github"
	"github.com/postqode/qodeflow-core/support/test"
	"github.com/stretchr/testify/assert"
)

func getTestSettings() map[string]interface{} {
	token := os.Getenv("GITHUB_TEST_TOKEN")
	if token == "" {
		token = "test-token"
	}
	return map[string]interface{}{
		"token": token,
	}
}

func getTestInput(method string, data interface{}) *Input {
	owner := os.Getenv("GITHUB_OWNER")
	if owner == "" {
		owner = "postqode"
	}
	repo := os.Getenv("GITHUB_REPO")
	if repo == "" {
		repo = "qodeflow-component"
	}
	return &Input{
		Owner:  owner,
		Repo:   repo,
		Method: method,
		Data:   data,
	}
}

func skipIfNoToken(t *testing.T) {
	if os.Getenv("GITHUB_TEST_TOKEN") == "" {
		t.Skip("GITHUB_TEST_TOKEN not set, skipping functional test")
	}
}

func TestMetadata(t *testing.T) {
	act := &GitHubActivity{}
	assert.Equal(t, activityMetadata, act.Metadata())
}

func TestNew(t *testing.T) {
	settings := getTestSettings()
	initCtx := test.NewActivityInitContext(settings, nil)
	act, err := New(initCtx)
	assert.Nil(t, err)
	assert.NotNil(t, act)
}

func TestEval_NoToken(t *testing.T) {
	settings := map[string]interface{}{
		"token": "",
	}
	initCtx := test.NewActivityInitContext(settings, nil)
	act, _ := New(initCtx)

	tc := test.NewActivityContext(act.Metadata())
	input := getTestInput("LIST_REPOS", nil)
	tc.SetInputObject(input)

	done, err := act.Eval(tc)
	assert.False(t, done)
	assert.NotNil(t, err)
	assert.Equal(t, "github token is required", err.Error())
}

func TestEval_UnsupportedMethod(t *testing.T) {
	settings := getTestSettings()
	initCtx := test.NewActivityInitContext(settings, nil)
	act, _ := New(initCtx)

	tc := test.NewActivityContext(act.Metadata())
	input := getTestInput("INVALID_METHOD", nil)
	tc.SetInputObject(input)

	done, err := act.Eval(tc)
	assert.False(t, done)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "unsupported method")
}

// Functional Tests - require GITHUB_TEST_TOKEN

func TestEval_GetRepository(t *testing.T) {
	skipIfNoToken(t)
	settings := getTestSettings()
	initCtx := test.NewActivityInitContext(settings, nil)
	act, _ := New(initCtx)

	tc := test.NewActivityContext(act.Metadata())
	input := getTestInput("GET_REPOSITORY", nil)
	tc.SetInputObject(input)

	done, err := act.Eval(tc)
	assert.Nil(t, err)
	assert.True(t, done)

	output := &Output{}
	tc.GetOutputObject(output)
	assert.Empty(t, output.Error)
	assert.NotNil(t, output.Result)
}

func TestEval_ListRepos(t *testing.T) {
	skipIfNoToken(t)
	settings := getTestSettings()
	initCtx := test.NewActivityInitContext(settings, nil)
	act, _ := New(initCtx)

	tc := test.NewActivityContext(act.Metadata())
	input := getTestInput("LIST_REPOS", nil)
	tc.SetInputObject(input)

	done, err := act.Eval(tc)
	assert.Nil(t, err)
	assert.True(t, done)

	output := &Output{}
	tc.GetOutputObject(output)
	assert.Empty(t, output.Error)
	assert.NotNil(t, output.Result)
}

func TestEval_ListIssues(t *testing.T) {
	skipIfNoToken(t)
	settings := getTestSettings()
	initCtx := test.NewActivityInitContext(settings, nil)
	act, _ := New(initCtx)

	tc := test.NewActivityContext(act.Metadata())
	input := getTestInput("LIST_ISSUES", map[string]interface{}{})
	tc.SetInputObject(input)

	done, err := act.Eval(tc)
	assert.Nil(t, err)
	assert.True(t, done)

	output := &Output{}
	tc.GetOutputObject(output)
	assert.Empty(t, output.Error)
	assert.NotNil(t, output.Result)
}

func TestEval_CreateIssue(t *testing.T) {
	skipIfNoToken(t)
	settings := getTestSettings()
	initCtx := test.NewActivityInitContext(settings, nil)
	act, _ := New(initCtx)

	tc := test.NewActivityContext(act.Metadata())
	input := getTestInput("CREATE_ISSUE", map[string]interface{}{
		"title": "Test Issue from Unit Test",
		"body":  "This issue was created by a unit test.",
	})
	tc.SetInputObject(input)

	done, err := act.Eval(tc)
	assert.Nil(t, err)
	assert.True(t, done)

	output := &Output{}
	tc.GetOutputObject(output)
	assert.Empty(t, output.Error)
	assert.NotNil(t, output.Result)
}

func TestEval_AssignIssue(t *testing.T) {
	skipIfNoToken(t)
	// This test requires a valid issue number. We use a placeholder or skip if not provided.
	issueNum := os.Getenv("GITHUB_ISSUE_NUMBER")
	if issueNum == "" {
		t.Skip("GITHUB_ISSUE_NUMBER not set, skipping ASSIGN_ISSUE test")
	}

	settings := getTestSettings()
	initCtx := test.NewActivityInitContext(settings, nil)
	act, _ := New(initCtx)

	tc := test.NewActivityContext(act.Metadata())
	input := getTestInput("ASSIGN_ISSUE", map[string]interface{}{
		"issueNumber": issueNum,
		"assignees":   []string{os.Getenv("GITHUB_OWNER")},
	})
	tc.SetInputObject(input)

	done, err := act.Eval(tc)
	assert.Nil(t, err)
	assert.True(t, done)
}

func TestEval_RemoveAssignee(t *testing.T) {
	skipIfNoToken(t)
	issueNum := os.Getenv("GITHUB_ISSUE_NUMBER")
	if issueNum == "" {
		t.Skip("GITHUB_ISSUE_NUMBER not set, skipping REMOVE_ASSIGNEE test")
	}

	settings := getTestSettings()
	initCtx := test.NewActivityInitContext(settings, nil)
	act, _ := New(initCtx)

	tc := test.NewActivityContext(act.Metadata())
	input := getTestInput("REMOVE_ASSIGNEE", map[string]interface{}{
		"issueNumber": issueNum,
		"assignees":   []string{os.Getenv("GITHUB_OWNER")},
	})
	tc.SetInputObject(input)

	done, err := act.Eval(tc)
	assert.Nil(t, err)
	assert.True(t, done)
}

func TestEval_Files(t *testing.T) {
	skipIfNoToken(t)
	settings := getTestSettings()
	initCtx := test.NewActivityInitContext(settings, nil)
	act, _ := New(initCtx)

	// Create File
	tc := test.NewActivityContext(act.Metadata())
	path := "test_file_new.txt"
	input := getTestInput("CREATE_FILE", map[string]interface{}{
		"path":    path,
		"content": "Hello GitHub! New file.",
		"message": "Test commit create",
	})
	tc.SetInputObject(input)

	done, err := act.Eval(tc)
	assert.Nil(t, err)
	assert.True(t, done)

	output := &Output{}
	tc.GetOutputObject(output)
	assert.Empty(t, output.Error)

	// Update File if SHA is provided
	sha := os.Getenv("GITHUB_FILE_SHA")
	if sha != "" {
		input = getTestInput("UPDATE_FILE", map[string]interface{}{
			"path":    path,
			"content": "Hello GitHub! Updated file.",
			"message": "Test commit update",
			"sha":     sha,
		})
		tc.SetInputObject(input)
		act.Eval(tc)
		tc.GetOutputObject(output)
		assert.Empty(t, output.Error)
	}

	// Delete File if SHA is provided
	if sha != "" {
		input = getTestInput("DELETE_FILE", map[string]interface{}{
			"path":    path,
			"message": "Test commit delete",
			"sha":     sha,
		})
		tc.SetInputObject(input)
		act.Eval(tc)
		tc.GetOutputObject(output)
		assert.Empty(t, output.Error)
	}
}

func TestEval_Collaborators(t *testing.T) {
	skipIfNoToken(t)
	username := os.Getenv("GITHUB_COLLABORATOR")
	if username == "" {
		t.Skip("GITHUB_COLLABORATOR not set, skipping collaborator tests")
	}

	settings := getTestSettings()
	initCtx := test.NewActivityInitContext(settings, nil)
	act, _ := New(initCtx)

	tc := test.NewActivityContext(act.Metadata())
	input := getTestInput("ADD_COLLABORATOR", map[string]interface{}{
		"username": username,
	})
	tc.SetInputObject(input)

	done, err := act.Eval(tc)
	assert.Nil(t, err)
	assert.True(t, done)
}

func TestEval_RemoveCollaborator(t *testing.T) {
	skipIfNoToken(t)
	username := os.Getenv("GITHUB_COLLABORATOR")
	if username == "" {
		t.Skip("GITHUB_COLLABORATOR not set, skipping REMOVE_COLLABORATOR test")
	}

	settings := getTestSettings()
	initCtx := test.NewActivityInitContext(settings, nil)
	act, _ := New(initCtx)

	tc := test.NewActivityContext(act.Metadata())
	input := getTestInput("REMOVE_COLLABORATOR", map[string]interface{}{
		"username": username,
	})
	tc.SetInputObject(input)

	done, err := act.Eval(tc)
	assert.Nil(t, err)
	assert.True(t, done)
}

func TestEval_PullRequests(t *testing.T) {
	skipIfNoToken(t)
	prNum := os.Getenv("GITHUB_PR_NUMBER")
	if prNum == "" {
		t.Skip("GITHUB_PR_NUMBER not set, skipping pull request tests")
	}

	settings := getTestSettings()
	initCtx := test.NewActivityInitContext(settings, nil)
	act, _ := New(initCtx)

	tc := test.NewActivityContext(act.Metadata())
	input := getTestInput("GET_PULL_REQUEST", map[string]interface{}{
		"pullRequestNumber": prNum,
	})
	tc.SetInputObject(input)

	done, err := act.Eval(tc)
	assert.Nil(t, err)
	assert.True(t, done)
}

func TestEval_UpdateIssue(t *testing.T) {
	skipIfNoToken(t)
	issueNum := os.Getenv("GITHUB_ISSUE_NUMBER")
	if issueNum == "" {
		t.Skip("GITHUB_ISSUE_NUMBER not set, skipping UPDATE_ISSUE test")
	}

	settings := getTestSettings()
	initCtx := test.NewActivityInitContext(settings, nil)
	act, _ := New(initCtx)

	tc := test.NewActivityContext(act.Metadata())
	input := getTestInput("UPDATE_ISSUE", map[string]interface{}{
		"issueNumber": issueNum,
		"title":       "Updated Issue Title from Unit Test",
		"body":        "This issue body was updated by a unit test.",
	})
	tc.SetInputObject(input)

	done, err := act.Eval(tc)
	assert.Nil(t, err)
	assert.True(t, done)

	output := &Output{}
	tc.GetOutputObject(output)
	assert.Empty(t, output.Error)
}

func TestEval_IssueComments(t *testing.T) {
	skipIfNoToken(t)
	issueNum := os.Getenv("GITHUB_ISSUE_NUMBER")
	if issueNum == "" {
		t.Skip("GITHUB_ISSUE_NUMBER not set, skipping issue comment tests")
	}

	settings := getTestSettings()
	initCtx := test.NewActivityInitContext(settings, nil)
	act, _ := New(initCtx)

	// Create Comment
	tc := test.NewActivityContext(act.Metadata())
	input := getTestInput("CREATE_COMMENT_ON_ISSUE", map[string]interface{}{
		"issueNumber": issueNum,
		"body":        "Test comment from unit test",
	})
	tc.SetInputObject(input)

	done, err := act.Eval(tc)
	assert.Nil(t, err)
	assert.True(t, done)

	output := &Output{}
	tc.GetOutputObject(output)
	assert.Empty(t, output.Error)

	// Update Comment if CommentID is provided
	commentID := os.Getenv("GITHUB_COMMENT_ID")
	if commentID != "" {
		input = getTestInput("UPDATE_COMMENT_ON_ISSUE", map[string]interface{}{
			"commentID": commentID,
			"body":      "Updated test comment from unit test",
		})
		tc.SetInputObject(input)
		act.Eval(tc)
	}
}

func TestEval_RepositoryOperations(t *testing.T) {
	skipIfNoToken(t)
	repoName := os.Getenv("GITHUB_NEW_REPO_NAME")
	if repoName == "" {
		t.Skip("GITHUB_NEW_REPO_NAME not set, skipping repository lifecycle tests")
	}

	settings := getTestSettings()
	initCtx := test.NewActivityInitContext(settings, nil)
	act, _ := New(initCtx)

	// Create Repository
	tc := test.NewActivityContext(act.Metadata())
	input := getTestInput("CREATE_REPOSITORY", map[string]interface{}{
		"name": repoName,
	})
	tc.SetInputObject(input)

	done, err := act.Eval(tc)
	assert.Nil(t, err)
	assert.True(t, done)

	output := &Output{}
	tc.GetOutputObject(output)
	assert.Empty(t, output.Error)

	// Delete Repository (Optional - usually best to keep it for manual check or use a temp name)
	// We'll just verify the method exists and can be called if needed.
}

func TestWorkflow_Lifecycle(t *testing.T) {
	skipIfNoToken(t)
	repoName := "qodeflow-test"
	settings := getTestSettings()
	initCtx := test.NewActivityInitContext(settings, nil)
	act, _ := New(initCtx)

	// 1. Create Repository
	t.Run("CreateRepository", func(t *testing.T) {
		tc := test.NewActivityContext(act.Metadata())
		input := getTestInput("CREATE_REPOSITORY", map[string]interface{}{
			"name": repoName,
		})
		tc.SetInputObject(input)
		done, err := act.Eval(tc)
		assert.Nil(t, err)
		assert.True(t, done)
		output := &Output{}
		tc.GetOutputObject(output)
		assert.Empty(t, output.Error)
	})

	// 2. Create File
	t.Run("CreateFile", func(t *testing.T) {
		tc := test.NewActivityContext(act.Metadata())
		input := getTestInput("CREATE_FILE", map[string]interface{}{
			"path":    "README.md",
			"content": "# Qodeflow Test\nThis is a test repo.",
			"message": "Initial commit",
		})
		input.Repo = repoName // Use the newly created repo
		tc.SetInputObject(input)
		done, err := act.Eval(tc)
		assert.Nil(t, err)
		assert.True(t, done)
		output := &Output{}
		tc.GetOutputObject(output)
		assert.Empty(t, output.Error)
	})

	// 3. Create Issue
	var issueNum int
	t.Run("CreateIssue", func(t *testing.T) {
		tc := test.NewActivityContext(act.Metadata())
		input := getTestInput("CREATE_ISSUE", map[string]interface{}{
			"title": "testing issue",
			"body":  "This is a test issue for workflow verification.",
		})
		input.Repo = repoName
		tc.SetInputObject(input)
		done, err := act.Eval(tc)
		assert.Nil(t, err)
		assert.True(t, done)
		output := &Output{}
		tc.GetOutputObject(output)
		assert.Empty(t, output.Error)

		// Result is gh.Issue
		if res, ok := output.Result.(*gh.Issue); ok {
			issueNum = res.GetNumber()
		}
	})

	// 4. Assign Issue
	t.Run("AssignIssue", func(t *testing.T) {
		if issueNum == 0 {
			t.Skip("Issue number not captured, skipping assign")
		}
		tc := test.NewActivityContext(act.Metadata())
		owner := os.Getenv("GITHUB_OWNER")
		if owner == "" {
			owner = "postqode"
		}
		input := getTestInput("ASSIGN_ISSUE", map[string]interface{}{
			"issueNumber": issueNum,
			"assignees":   []string{owner},
		})
		input.Repo = repoName
		tc.SetInputObject(input)
		done, err := act.Eval(tc)
		assert.Nil(t, err)
		assert.True(t, done)
		output := &Output{}
		tc.GetOutputObject(output)
		assert.Empty(t, output.Error)
	})
}
