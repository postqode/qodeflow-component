package github

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/go-github/v60/github"
	"github.com/postqode/qodeflow-core/activity"
	"github.com/postqode/qodeflow-core/data/coerce"
	"github.com/postqode/qodeflow-core/data/metadata"
	"golang.org/x/oauth2"
)

func init() {
	_ = activity.Register(&GitHubActivity{}, New)
}

func toStringArray(val interface{}) []string {
	if val == nil {
		return nil
	}
	arr, err := coerce.ToArray(val)
	if err != nil {
		return nil
	}
	res := make([]string, len(arr))
	for i, v := range arr {
		res[i], _ = coerce.ToString(v)
	}
	return res
}

type GitHubActivity struct {
	settings *Settings
}

func New(ctx activity.InitContext) (activity.Activity, error) {
	s := &Settings{}
	err := metadata.MapToStruct(ctx.Settings(), s, true)
	if err != nil {
		return nil, err
	}

	return &GitHubActivity{settings: s}, nil
}

func (a *GitHubActivity) Metadata() *activity.Metadata {
	return activityMetadata
}

func (a *GitHubActivity) Eval(ctx activity.Context) (done bool, err error) {
	input := &Input{}
	err = ctx.GetInputObject(input)
	if err != nil {
		return true, err
	}

	if a.settings.Token == "" {
		return false, fmt.Errorf("github token is required")
		// can we introduce the logic to take the token from the global env if not provided.
		// Here also value can come from the properties as well if seen from qodeflow ui
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: a.settings.Token},
	)
	tc := oauth2.NewClient(context.Background(), ts)
	client := github.NewClient(tc)

	var result interface{}
	var opErr error

	switch strings.ToUpper(input.Method) {
	case "CREATE_ISSUE":
		data, ok := input.Data.(map[string]interface{})
		if !ok {
			opErr = fmt.Errorf("data must be a map for CREATE_ISSUE")
			break
		}
		title, _ := coerce.ToString(data["title"])
		body, _ := coerce.ToString(data["body"])

		state, _ := coerce.ToString(data["state"])
		stateReason, _ := coerce.ToString(data["stateReason"])
		labels := toStringArray(data["labels"])
		milestone, _ := coerce.ToInt(data["milestone"])
		assignees := toStringArray(data["assignees"])

		if title == "" {
			opErr = fmt.Errorf("title is required for CREATE_ISSUE")
			break
		}

		issueRequest := &github.IssueRequest{
			Title: &title,
			Body:  &body,
		}

		if state != "" {
			issueRequest.State = &state
		}
		if stateReason != "" {
			issueRequest.StateReason = &stateReason
		}
		if milestone != 0 {
			issueRequest.Milestone = &milestone
		}
		if len(labels) > 0 {
			issueRequest.Labels = &labels
		}
		if len(assignees) > 0 {
			issueRequest.Assignees = &assignees
		}

		issue, _, err := client.Issues.Create(context.Background(), input.Owner, input.Repo, issueRequest)
		result = issue
		opErr = err

	case "UPDATE_ISSUE":
		data, ok := input.Data.(map[string]interface{})
		if !ok {
			opErr = fmt.Errorf("data must be a map for UPDATE_ISSUE")
			break
		}
		issueNumber, _ := coerce.ToInt(data["issueNumber"])
		title, _ := coerce.ToString(data["title"])
		body, _ := coerce.ToString(data["body"])

		state, _ := coerce.ToString(data["state"])
		stateReason, _ := coerce.ToString(data["stateReason"])
		labels := toStringArray(data["labels"])
		milestone, _ := coerce.ToInt(data["milestone"])
		assignees := toStringArray(data["assignees"])

		if issueNumber == 0 {
			opErr = fmt.Errorf("issueNumber is required for UPDATE_ISSUE")
			break
		}

		issueRequest := &github.IssueRequest{
			Title: &title,
			Body:  &body,
		}

		if state != "" {
			issueRequest.State = &state
		}
		if stateReason != "" {
			issueRequest.StateReason = &stateReason
		}
		if milestone != 0 {
			issueRequest.Milestone = &milestone
		}
		if len(labels) > 0 {
			issueRequest.Labels = &labels
		}
		if len(assignees) > 0 {
			issueRequest.Assignees = &assignees
		}

		issue, _, err := client.Issues.Edit(context.Background(), input.Owner, input.Repo, issueNumber, issueRequest)
		result = issue
		opErr = err

	case "CREATE_COMMENT_ON_ISSUE":
		data, ok := input.Data.(map[string]interface{})
		if !ok {
			opErr = fmt.Errorf("data must be a map for COMMENT_ON_ISSUE")
			break
		}
		issueNumber, _ := coerce.ToInt(data["issueNumber"])
		body, _ := coerce.ToString(data["body"])

		if issueNumber == 0 || body == "" {
			opErr = fmt.Errorf("issueNumber and body are required for COMMENT_ON_ISSUE")
			break
		}

		comment, _, err := client.Issues.CreateComment(context.Background(), input.Owner, input.Repo, issueNumber, &github.IssueComment{
			Body: &body,
		})
		result = comment
		opErr = err

	case "UPDATE_COMMENT_ON_ISSUE":
		data, ok := input.Data.(map[string]interface{})
		if !ok {
			opErr = fmt.Errorf("data must be a map for UPDATE_COMMENT_ON_ISSUE")
			break
		}
		commentID, _ := coerce.ToInt64(data["commentID"])
		body, _ := coerce.ToString(data["body"])

		if commentID == 0 || body == "" {
			opErr = fmt.Errorf("commentID and body are required for UPDATE_COMMENT_ON_ISSUE")
			break
		}

		comment, _, err := client.Issues.EditComment(context.Background(), input.Owner, input.Repo, commentID, &github.IssueComment{
			Body: &body,
		})
		result = comment
		opErr = err

	case "LIST_ISSUES":
		data, ok := input.Data.(map[string]interface{})
		if !ok {
			opErr = fmt.Errorf("data must be a map for LIST_ISSUES")
			break
		}

		milestone, _ := coerce.ToString(data["milestone"])
		state, _ := coerce.ToString(data["state"])
		labels := toStringArray(data["labels"])
		sort, _ := coerce.ToString(data["sort"])
		direction, _ := coerce.ToString(data["direction"])
		since, _ := coerce.ToDateTime(data["since"])
		creator, _ := coerce.ToString(data["creator"])
		assignee, _ := coerce.ToString(data["assignee"])
		mentioned, _ := coerce.ToString(data["mentioned"])

		if input.Owner == "" || input.Repo == "" {
			opErr = fmt.Errorf("owner and repo are required for LIST_ISSUES")
			break
		}

		opt := &github.IssueListByRepoOptions{
			Milestone: milestone,
			State:     state,
			Labels:    labels,
			Sort:      sort,
			Direction: direction,
			Since:     since,
			Creator:   creator,
			Assignee:  assignee,
			Mentioned: mentioned,
		}

		issues, _, err := client.Issues.ListByRepo(context.Background(), input.Owner, input.Repo, opt)
		result = issues
		opErr = err

	case "CREATE_PULL_REQUEST":
		data, ok := input.Data.(map[string]interface{})
		if !ok {
			opErr = fmt.Errorf("data must be a map for CREATE_PULL_REQUEST")
			break
		}
		title, _ := coerce.ToString(data["title"])
		head, _ := coerce.ToString(data["head"])
		base, _ := coerce.ToString(data["base"])
		body, _ := coerce.ToString(data["body"])

		if title == "" || head == "" || base == "" {
			opErr = fmt.Errorf("title, head, and base are required for CREATE_PULL_REQUEST")
			break
		}

		newPR := &github.NewPullRequest{
			Title: &title,
			Head:  &head,
			Base:  &base,
			Body:  &body,
		}

		pr, _, err := client.PullRequests.Create(context.Background(), input.Owner, input.Repo, newPR)
		result = pr
		opErr = err

	case "GET_REPOSITORY":
		repo, _, err := client.Repositories.Get(context.Background(), input.Owner, input.Repo)
		result = repo
		opErr = err

	case "CREATE_REPOSITORY":
		data, ok := input.Data.(map[string]interface{})
		if !ok {
			opErr = fmt.Errorf("data must be a map for CREATE_REPOSITORY")
			break
		}
		name, _ := coerce.ToString(data["name"])
		if name == "" {
			opErr = fmt.Errorf("name is required for CREATE_REPOSITORY")
			break
		}

		newRepo := &github.Repository{
			Name: github.String(name),
		}

		if val, ok := data["description"]; ok {
			s, _ := coerce.ToString(val)
			newRepo.Description = github.String(s)
		}
		if val, ok := data["homepage"]; ok {
			s, _ := coerce.ToString(val)
			newRepo.Homepage = github.String(s)
		}
		if val, ok := data["private"]; ok {
			b, _ := coerce.ToBool(val)
			newRepo.Private = github.Bool(b)
		}
		if val, ok := data["visibility"]; ok {
			s, _ := coerce.ToString(val)
			newRepo.Visibility = github.String(s)
		}
		if val, ok := data["has_issues"]; ok {
			b, _ := coerce.ToBool(val)
			newRepo.HasIssues = github.Bool(b)
		}
		if val, ok := data["has_projects"]; ok {
			b, _ := coerce.ToBool(val)
			newRepo.HasProjects = github.Bool(b)
		}
		if val, ok := data["has_wiki"]; ok {
			b, _ := coerce.ToBool(val)
			newRepo.HasWiki = github.Bool(b)
		}
		if val, ok := data["has_discussions"]; ok {
			b, _ := coerce.ToBool(val)
			newRepo.HasDiscussions = github.Bool(b)
		}
		if val, ok := data["is_template"]; ok {
			b, _ := coerce.ToBool(val)
			newRepo.IsTemplate = github.Bool(b)
		}
		if val, ok := data["team_id"]; ok {
			i, _ := coerce.ToInt64(val)
			newRepo.TeamID = github.Int64(i)
		}
		if val, ok := data["auto_init"]; ok {
			b, _ := coerce.ToBool(val)
			newRepo.AutoInit = github.Bool(b)
		}
		if val, ok := data["gitignore_template"]; ok {
			s, _ := coerce.ToString(val)
			newRepo.GitignoreTemplate = github.String(s)
		}
		if val, ok := data["license_template"]; ok {
			s, _ := coerce.ToString(val)
			newRepo.LicenseTemplate = github.String(s)
		}
		if val, ok := data["allow_squash_merge"]; ok {
			b, _ := coerce.ToBool(val)
			newRepo.AllowSquashMerge = github.Bool(b)
		}
		if val, ok := data["allow_merge_commit"]; ok {
			b, _ := coerce.ToBool(val)
			newRepo.AllowMergeCommit = github.Bool(b)
		}
		if val, ok := data["allow_rebase_merge"]; ok {
			b, _ := coerce.ToBool(val)
			newRepo.AllowRebaseMerge = github.Bool(b)
		}
		if val, ok := data["allow_update_branch"]; ok {
			b, _ := coerce.ToBool(val)
			newRepo.AllowUpdateBranch = github.Bool(b)
		}
		if val, ok := data["allow_auto_merge"]; ok {
			b, _ := coerce.ToBool(val)
			newRepo.AllowAutoMerge = github.Bool(b)
		}
		if val, ok := data["allow_forking"]; ok {
			b, _ := coerce.ToBool(val)
			newRepo.AllowForking = github.Bool(b)
		}
		if val, ok := data["delete_branch_on_merge"]; ok {
			b, _ := coerce.ToBool(val)
			newRepo.DeleteBranchOnMerge = github.Bool(b)
		}
		if val, ok := data["use_squash_pr_title_as_default"]; ok {
			b, _ := coerce.ToBool(val)
			newRepo.UseSquashPRTitleAsDefault = github.Bool(b)
		}
		if val, ok := data["squash_merge_commit_title"]; ok {
			s, _ := coerce.ToString(val)
			newRepo.SquashMergeCommitTitle = github.String(s)
		}
		if val, ok := data["squash_merge_commit_message"]; ok {
			s, _ := coerce.ToString(val)
			newRepo.SquashMergeCommitMessage = github.String(s)
		}
		if val, ok := data["merge_commit_title"]; ok {
			s, _ := coerce.ToString(val)
			newRepo.MergeCommitTitle = github.String(s)
		}
		if val, ok := data["merge_commit_message"]; ok {
			s, _ := coerce.ToString(val)
			newRepo.MergeCommitMessage = github.String(s)
		}

		repo, _, err := client.Repositories.Create(context.Background(), input.Owner, newRepo)
		result = repo
		opErr = err

	case "DELETE_REPOSITORY":
		_, err := client.Repositories.Delete(context.Background(), input.Owner, input.Repo)
		opErr = err

	case "LIST_REPOS":
		opt := &github.RepositoryListOptions{
			ListOptions: github.ListOptions{PerPage: 10},
		}
		repos, _, err := client.Repositories.List(context.Background(), input.Owner, opt)
		result = repos
		opErr = err

	case "ASSIGN_ISSUE":
		data, ok := input.Data.(map[string]interface{})
		if !ok {
			opErr = fmt.Errorf("data must be a map for ASSIGN_ISSUE")
			break
		}
		issueNumber, _ := coerce.ToInt(data["issueNumber"])
		assignees := toStringArray(data["assignees"])

		if issueNumber == 0 || len(assignees) == 0 {
			opErr = fmt.Errorf("issueNumber and assignees are required for ASSIGN_ISSUE")
			break
		}

		issue, _, err := client.Issues.AddAssignees(context.Background(), input.Owner, input.Repo, issueNumber, assignees)
		result = issue
		opErr = err

	case "REMOVE_ASSIGNEE":
		data, ok := input.Data.(map[string]interface{})
		if !ok {
			opErr = fmt.Errorf("data must be a map for REMOVE_ASSIGNEE")
			break
		}
		issueNumber, _ := coerce.ToInt(data["issueNumber"])
		assignees := toStringArray(data["assignees"])

		if issueNumber == 0 || len(assignees) == 0 {
			opErr = fmt.Errorf("issueNumber and assignees are required for REMOVE_ASSIGNEE")
			break
		}

		issue, _, err := client.Issues.RemoveAssignees(context.Background(), input.Owner, input.Repo, issueNumber, assignees)
		result = issue
		opErr = err

	case "CREATE_FILE":
		data, ok := input.Data.(map[string]interface{})
		if !ok {
			opErr = fmt.Errorf("data must be a map for CREATE_FILE")
			break
		}
		path, _ := coerce.ToString(data["path"])
		content, _ := coerce.ToString(data["content"])
		message, _ := coerce.ToString(data["message"])
		branch, _ := coerce.ToString(data["branch"])

		if path == "" || content == "" || message == "" {
			opErr = fmt.Errorf("path, content, and message are required for CREATE_FILE")
			break
		}

		opts := &github.RepositoryContentFileOptions{
			Message: github.String(message),
			Content: []byte(content),
		}
		if branch != "" {
			opts.Branch = github.String(branch)
		}

		resp, _, err := client.Repositories.CreateFile(context.Background(), input.Owner, input.Repo, path, opts)
		result = resp
		opErr = err

	case "UPDATE_FILE":
		data, ok := input.Data.(map[string]interface{})
		if !ok {
			opErr = fmt.Errorf("data must be a map for UPDATE_FILE")
			break
		}
		path, _ := coerce.ToString(data["path"])
		content, _ := coerce.ToString(data["content"])
		message, _ := coerce.ToString(data["message"])
		branch, _ := coerce.ToString(data["branch"])
		sha, _ := coerce.ToString(data["sha"])

		if path == "" || content == "" || message == "" || sha == "" {
			opErr = fmt.Errorf("path, content, message, and sha are required for UPDATE_FILE")
			break
		}

		opts := &github.RepositoryContentFileOptions{
			Message: github.String(message),
			Content: []byte(content),
			SHA:     github.String(sha),
		}
		if branch != "" {
			opts.Branch = github.String(branch)
		}

		resp, _, err := client.Repositories.UpdateFile(context.Background(), input.Owner, input.Repo, path, opts)
		result = resp
		opErr = err

	case "DELETE_FILE":
		data, ok := input.Data.(map[string]interface{})
		if !ok {
			opErr = fmt.Errorf("data must be a map for DELETE_FILE")
			break
		}
		path, _ := coerce.ToString(data["path"])
		message, _ := coerce.ToString(data["message"])
		sha, _ := coerce.ToString(data["sha"])
		branch, _ := coerce.ToString(data["branch"])

		if path == "" || message == "" || sha == "" {
			opErr = fmt.Errorf("path, message, and sha are required for DELETE_FILE")
			break
		}

		opts := &github.RepositoryContentFileOptions{
			Message: github.String(message),
			SHA:     github.String(sha),
		}
		if branch != "" {
			opts.Branch = github.String(branch)
		}

		resp, _, err := client.Repositories.DeleteFile(context.Background(), input.Owner, input.Repo, path, opts)
		result = resp
		opErr = err

	case "ADD_COLLABORATOR":
		data, ok := input.Data.(map[string]interface{})
		if !ok {
			opErr = fmt.Errorf("data must be a map for ADD_COLLABORATOR")
			break
		}
		user, _ := coerce.ToString(data["username"])
		permission, _ := coerce.ToString(data["permission"])

		if user == "" {
			opErr = fmt.Errorf("username is required for ADD_COLLABORATOR")
			break
		}

		opts := &github.RepositoryAddCollaboratorOptions{}
		if permission != "" {
			opts.Permission = permission
		}

		resp, _, err := client.Repositories.AddCollaborator(context.Background(), input.Owner, input.Repo, user, opts)
		result = resp
		opErr = err

	case "REMOVE_COLLABORATOR":
		data, ok := input.Data.(map[string]interface{})
		if !ok {
			opErr = fmt.Errorf("data must be a map for REMOVE_COLLABORATOR")
			break
		}
		user, _ := coerce.ToString(data["username"])

		if user == "" {
			opErr = fmt.Errorf("username is required for REMOVE_COLLABORATOR")
			break
		}

		_, err := client.Repositories.RemoveCollaborator(context.Background(), input.Owner, input.Repo, user)
		opErr = err

	case "GET_PULL_REQUEST":
		data, ok := input.Data.(map[string]interface{})
		if !ok {
			opErr = fmt.Errorf("data must be a map for GET_PULL_REQUEST")
			break
		}
		number, _ := coerce.ToInt(data["pullRequestNumber"])
		if number == 0 {
			opErr = fmt.Errorf("pullRequestNumber is required for GET_PULL_REQUEST")
			break
		}

		pr, _, err := client.PullRequests.Get(context.Background(), input.Owner, input.Repo, number)
		result = pr
		opErr = err

	case "MERGE_PULL_REQUEST":
		data, ok := input.Data.(map[string]interface{})
		if !ok {
			opErr = fmt.Errorf("data must be a map for MERGE_PULL_REQUEST")
			break
		}
		number, _ := coerce.ToInt(data["pullRequestNumber"])
		commitMsg, _ := coerce.ToString(data["commitMessage"])

		if number == 0 {
			opErr = fmt.Errorf("pullRequestNumber is required for MERGE_PULL_REQUEST")
			break
		}

		res, _, err := client.PullRequests.Merge(context.Background(), input.Owner, input.Repo, number, commitMsg, nil)
		result = res
		opErr = err

	default:
		return false, fmt.Errorf("unsupported method: %s", input.Method)
	}

	output := &Output{}
	if opErr != nil {
		output.Error = opErr.Error()
	} else {
		output.Result = result
	}

	err = ctx.SetOutputObject(output)
	if err != nil {
		return true, err
	}

	return true, nil
}
