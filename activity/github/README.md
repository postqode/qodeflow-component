<!--
title: GitHub Activity
weight: 4630
-->

# GitHub Activity
This activity allows you to perform operations on GitHub, such as creating issues, pull requests, and listing repositories.

## Installation

### Qodeflow CLI
```bash
qodeflow install github.com/postqode/qodeflow-component/activity/github
```

## Configuration

### Settings:
| Name     | Type   | Description
|:---      | :---   | :---
| token     | string | GitHub Personal Access Token

### Input:
| Name    | Type   | Description
|:---     | :---   | :---
| owner   | string | The owner of the repository (e.g., "postqode")
| repo    | string | The name of the repository (e.g., "qodeflow-component")
| method  | string | The operation to perform (CREATE_ISSUE, LIST_ISSUES, CREATE_PULL_REQUEST, GET_REPOSITORY, LIST_REPOS)
| data    | any    | Data for the operation (e.g., mapping for issue title/body)

### Output:
| Name    | Type   | Description
|:---     | :---   | :---
| result  | any    | The result of the GitHub operation
| error   | string | Error message if the operation failed

## Supported Methods

### CREATE_ISSUE
Requires `data` to contain `title` (required), `body` (optional), `state` (optional), `stateReason` (optional), `labels` (array of strings, optional), `milestone` (int, optional), and `assignees` (array of strings, optional).

### UPDATE_ISSUE
Requires `data` to contain `issueNumber` (int, required), and optional fields: `title`, `body`, `state`, `stateReason`, `labels`, `milestone`, `assignees`.

### CREATE_COMMENT_ON_ISSUE
Requires `data` to contain `issueNumber` (int) and `body` (string).

### UPDATE_COMMENT_ON_ISSUE
Requires `data` to contain `commentID` (int64) and `body` (string).

### CREATE_PULL_REQUEST
Requires `data` to contain `title`, `head`, and `base` (all required), and `body` (optional).

### GET_PULL_REQUEST
Requires `data` to contain `pullRequestNumber` (int).

### MERGE_PULL_REQUEST
Requires `data` to contain `pullRequestNumber` (int) and `commitMessage` (string, optional).

### LIST_ISSUES
Lists issues in the specified repository. Supports filtering via `data`:
- `milestone`: string (milestone number or `*`, `none`)
- `state`: string (`open`, `closed`, `all`)
- `labels`: array of strings
- `sort`: string (`created`, `updated`, `comments`)
- `direction`: string (`asc`, `desc`)
- `since`: string (ISO 8601 timestamp)
- `creator`: string (username)
- `assignee`: string (username, `*`, `none`)
- `mentioned`: string (username)

### CREATE_REPOSITORY
Requires `data` to contain `name` (string). Optional configuration fields in `data`:
- `description`: string
- `homepage`: string
- `private`: boolean
- `visibility`: string (`public`, `private`, `internal`)
- `has_issues`: boolean
- `has_projects`: boolean
- `has_wiki`: boolean
- `has_discussions`: boolean
- `is_template`: boolean
- `team_id`: int64 (Required for organization repositories if not an owner)
- `auto_init`: boolean
- `gitignore_template`: string (e.g., `Go`, `Node`)
- `license_template`: string (e.g., `mit`, `apache-2.0`)
- `allow_squash_merge`: boolean
- `allow_merge_commit`: boolean
- `allow_rebase_merge`: boolean
- `allow_update_branch`: boolean
- `allow_auto_merge`: boolean
- `allow_forking`: boolean
- `delete_branch_on_merge`: boolean
- `use_squash_pr_title_as_default`: boolean
- `squash_merge_commit_title`: string (`PR_TITLE`, `COMMIT_OR_PR_TITLE`)
- `squash_merge_commit_message`: string (`PR_BODY`, `COMMIT_MESSAGES`, `BLANK`)
- `merge_commit_title`: string (`PR_TITLE`, `MERGE_MESSAGE`)
- `merge_commit_message`: string (`PR_BODY`, `PR_TITLE`, `BLANK`)

### ASSIGN_ISSUE
Requires `data` to contain `issueNumber` (int) and `assignees` (array of strings).

### REMOVE_ASSIGNEE
Requires `data` to contain `issueNumber` (int) and `assignees` (array of strings).

### CREATE_FILE
Requires `data` to contain `path` (string), `content` (string), and `message` (string). Optionally `branch` (string).

### UPDATE_FILE
Requires `data` to contain `path` (string), `content` (string), `message` (string), and `sha` (string). Optionally `branch` (string).

### DELETE_FILE
Requires `data` to contain `path` (string), `message` (string), `sha` (string), and optionally `branch` (string).

### ADD_COLLABORATOR
Requires `data` to contain `username` (string) and optionally `permission` (string: pull, push, maintain, admin, triage).

### REMOVE_COLLABORATOR
Requires `data` to contain `username` (string).

### GET_REPOSITORY
Gets repository information.

### LIST_REPOS
Lists repositories for the authenticated user/org.

## Examples

### Create GitHub Issue
```json
{
  "id": "create_issue",
  "name": "Create GitHub Issue",
  "activity": {
    "ref": "github.com/postqode/qodeflow-component/activity/github",
    "input": {
      "owner": "postqode",
      "repo": "qodeflow-component",
      "method": "CREATE_ISSUE",
      "data": {
        "title": "Bug: Something is broken",
        "body": "Detailed description of the bug...",
        "labels": ["bug", "critical"]
      }
    }
  }
}
```

### Update Issue
```json
{
  "id": "update_issue",
  "name": "Update Issue",
  "activity": {
    "ref": "github.com/postqode/qodeflow-component/activity/github",
    "input": {
      "owner": "postqode",
      "repo": "qodeflow-component",
      "method": "UPDATE_ISSUE",
      "data": {
        "issueNumber": 1,
        "state": "closed"
      }
    }
  }
}
```

### Assign Issue
```json
{
  "id": "assign_issue",
  "name": "Assign Issue",
  "activity": {
    "ref": "github.com/postqode/qodeflow-component/activity/github",
    "input": {
      "owner": "postqode",
      "repo": "qodeflow-component",
      "method": "ASSIGN_ISSUE",
      "data": {
        "issueNumber": 1,
        "assignees": ["octocat"]
      }
    }
  }
}
```
