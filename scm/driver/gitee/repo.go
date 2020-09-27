// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gitee

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/drone/go-scm/scm"
)

type repository struct {
	ID    int `json:"id"`
	Owner struct {
		ID        int    `json:"id"`
		Login     string `json:"login"`
		AvatarURL string `json:"avatar_url"`
	} `json:"owner"`
	Name          string    `json:"name"`
	FullName      string    `json:"full_name"`
	Private       bool      `json:"private"`
	Fork          bool      `json:"fork"`
	HTMLURL       string    `json:"html_url"`
	SSHURL        string    `json:"ssh_url"`
	CloneURL      string    `json:"clone_url"`
	DefaultBranch string    `json:"default_branch"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	Permissions   struct {
		Admin bool `json:"admin"`
		Push  bool `json:"push"`
		Pull  bool `json:"pull"`
	} `json:"permission"`
}

type hookCreate struct {
	ID                  int    `json:"id,omitempty"`
	AccessToken         string `json:"access_token"`
	URL                 string `json:"url"`
	Password            string `json:"password"`                        //optinal
	PushEvents          bool   `json:"push_events" default:"true"`      //optinal
	TagPushEvents       bool   `json:"tag_push_events,omitempty"`       //optinal
	IssuesEvents        bool   `json:"issues_events,omitempty"`         //optinal
	NoteEvents          bool   `json:"note_events,omitempty"`           //optinal
	MergeRequestsEvents bool   `json:"merge_requests_events,omitempty"` //optinal
	CreatedAt           string `json:"created_at,omitempty"`
	ProjectID           int    `json:"project_id,omitempty"`
	Result              string `json:"result,omitempty"`
	ResultCode          int    `json:"result_code,omitempty"`
}

// RepositoryService implements the repository service for
// the GitHub driver.
type RepositoryService struct {
	client *wrapper
}

// Find returns the repository by name.
func (s *RepositoryService) Find(ctx context.Context, repo string) (*scm.Repository, *scm.Response, error) {
	path := fmt.Sprintf("repos/%s", repo)
	out := new(repository)
	res, err := s.client.do(ctx, "GET", api(path), nil, out)
	return convertRepository(out), res, err
}

// FindHook returns a repository hook.
func (s *RepositoryService) FindHook(ctx context.Context, repo string, id string) (*scm.Hook, *scm.Response, error) {
	path := fmt.Sprintf("repos/%s/hooks/%s", repo, id)
	out := new(hookCreate)
	res, err := s.client.do(ctx, "GET", api(path), nil, out)
	return convertHook(out), res, err
}

// FindPerms returns the repository permissions.
func (s *RepositoryService) FindPerms(ctx context.Context, repo string) (*scm.Perm, *scm.Response, error) {
	path := fmt.Sprintf("repos/%s", repo)
	out := new(repository)
	res, err := s.client.do(ctx, "GET", api(path), nil, out)
	return convertRepository(out).Perm, res, err
}

// List returns the user repository list.
func (s *RepositoryService) List(ctx context.Context, opts scm.ListOptions) ([]*scm.Repository, *scm.Response, error) {
	path := fmt.Sprintf("user/repos?%s", encodeListOptions(opts))
	out := []*repository{}
	res, err := s.client.do(ctx, "GET", api(path), nil, &out)
	return convertRepositoryList(out), res, err
}

// ListHooks returns a list or repository hooks.
func (s *RepositoryService) ListHooks(ctx context.Context, repo string, opts scm.ListOptions) ([]*scm.Hook, *scm.Response, error) {
	path := fmt.Sprintf("repos/%s/hooks?%s", repo, encodeListOptions(opts))
	out := []*hookCreate{}
	res, err := s.client.do(ctx, "GET", api(path), nil, &out)
	return convertHookList(out), res, err
}

// ListStatus returns a list of commit statuses.
func (s *RepositoryService) ListStatus(ctx context.Context, repo, ref string, opts scm.ListOptions) ([]*scm.Status, *scm.Response, error) {
	path := fmt.Sprintf("repos/%s/statuses/%s?%s", repo, ref, encodeListOptions(opts))
	out := []*status{}
	res, err := s.client.do(ctx, "GET", api(path), nil, &out)
	return convertStatusList(out), res, err
}

// CreateHook creates a new repository webhook.
func (s *RepositoryService) CreateHook(ctx context.Context, repo string, input *scm.HookInput) (*scm.Hook, *scm.Response, error) {
	path := fmt.Sprintf("repos/%s/hooks", repo)
	in := new(hookCreate)
	in = convertHookInput(input)
	token, _ := ctx.Value(scm.TokenKey{}).(*scm.Token)
	in.AccessToken = token.Token

	out := new(hookCreate)
	res, err := s.client.do(ctx, "POST", api(path), in, out)
	return convertHook(out), res, err
}

// CreateStatus creates a new commit status.
func (s *RepositoryService) CreateStatus(ctx context.Context, repo, ref string, input *scm.StatusInput) (*scm.Status, *scm.Response, error) {
	path := fmt.Sprintf("repos/%s/statuses/%s", repo, ref)
	in := &status{
		State:       convertFromState(input.State),
		Context:     input.Label,
		Description: input.Desc,
		TargetURL:   input.Target,
	}
	out := new(status)
	res, err := s.client.do(ctx, "POST", api(path), in, out)
	return convertStatus(out), res, err
}

// CreateDeployStatus creates a new deployment status.
func (s *RepositoryService) CreateDeployStatus(ctx context.Context, repo string, input *scm.DeployStatus) (*scm.DeployStatus, *scm.Response, error) {
	path := fmt.Sprintf("repos/%s/deployments/%d/statuses", repo, input.Number)
	in := &deployStatus{
		State:          convertFromState(input.State),
		Environment:    input.Environment,
		EnvironmentURL: input.EnvironmentURL,
		Description:    input.Desc,
		TargetURL:      input.Target,
	}
	out := new(deployStatus)
	res, err := s.client.do(ctx, "POST", api(path), in, out)
	return convertDeployStatus(out), res, err
}

// UpdateHook updates a repository webhook.
func (s *RepositoryService) UpdateHook(ctx context.Context, repo, id string, input *scm.HookInput) (*scm.Hook, *scm.Response, error) {
	path := fmt.Sprintf("repos/%s/hooks/%s", repo, id)
	in := new(hookCreate)
	in.ID, _ = strconv.Atoi(id)
	in = convertHookInput(input)

	out := new(hookCreate)
	res, err := s.client.do(ctx, "PATCH", api(path), in, out)
	return convertHook(out), res, err
}

// DeleteHook deletes a repository webhook.
func (s *RepositoryService) DeleteHook(ctx context.Context, repo, id string) (*scm.Response, error) {
	path := fmt.Sprintf("repos/%s/hooks/%s", repo, id)
	return s.client.do(ctx, "DELETE", api(path), nil, nil)
}

// helper function to convert from the gogs repository list to
// the common repository structure.
func convertRepositoryList(from []*repository) []*scm.Repository {
	to := []*scm.Repository{}
	for _, v := range from {
		to = append(to, convertRepository(v))
	}
	return to
}

// helper function to convert from the gogs repository structure
// to the common repository structure.
func convertRepository(from *repository) *scm.Repository {
	return &scm.Repository{
		ID:        strconv.Itoa(from.ID),
		Name:      from.Name,
		Namespace: from.Owner.Login,
		Perm: &scm.Perm{
			Push:  from.Permissions.Push,
			Pull:  from.Permissions.Pull,
			Admin: from.Permissions.Admin,
		},
		Link:     from.HTMLURL,
		Branch:   from.DefaultBranch,
		Private:  from.Private,
		Clone:    from.CloneURL,
		CloneSSH: from.SSHURL,
		Created:  from.CreatedAt,
		Updated:  from.UpdatedAt,
	}
}

func convertHookList(from []*hookCreate) []*scm.Hook {
	to := []*scm.Hook{}
	for _, v := range from {
		to = append(to, convertHook(v))
	}
	return to
}

func convertHookInput(from *scm.HookInput) *hookCreate {
	in := new(hookCreate)
	in.Password = from.Secret
	in.URL = from.Target

	in.PushEvents = from.Events.Push
	in.TagPushEvents = from.Events.Tag
	in.IssuesEvents = from.Events.Issue
	in.MergeRequestsEvents = from.Events.PullRequest // pull request is merge event ?

	return in
}

func convertHook(from *hookCreate) *scm.Hook {

	events := []string{"create", "delete", "deployment"}

	if from.PushEvents {
		events = append(events, "push")
	}

	if from.MergeRequestsEvents {
		events = append(events, "pull_request")
		events = append(events, "pull_request_review_comment")
	}

	if from.IssuesEvents {
		events = append(events, "issues")
		events = append(events, "issue_comment")
	}

	fromID := func() string {
		if from.ID == 0 {
			return ""
		}
		return strconv.Itoa(from.ID)
	}()
	return &scm.Hook{
		ID:         fromID,
		Active:     true,
		Target:     from.URL,
		Events:     events,
		SkipVerify: true,
	}
}

func convertFromHookEvents(from scm.HookEvents) []string {
	var events []string
	if from.Push {
		events = append(events, "push")
	}
	if from.PullRequest {
		events = append(events, "pull_request")
	}
	if from.PullRequestComment {
		events = append(events, "pull_request_review_comment")
	}
	if from.Issue {
		events = append(events, "issues")
	}
	if from.IssueComment || from.PullRequestComment {
		events = append(events, "issue_comment")
	}
	if from.Branch || from.Tag {
		events = append(events, "create")
		events = append(events, "delete")
	}
	if from.Deployment {
		events = append(events, "deployment")
	}
	return events
}

type status struct {
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	State       string    `json:"state"`
	TargetURL   string    `json:"target_url"`
	Description string    `json:"description"`
	Context     string    `json:"context"`
}

type deployStatus struct {
	ID             int64  `json:"id"`
	Environment    string `json:"environment"`
	EnvironmentURL string `json:"environment_url"`
	State          string `json:"state"`
	TargetURL      string `json:"log_url"`
	Description    string `json:"description"`
}

func convertStatusList(from []*status) []*scm.Status {
	to := []*scm.Status{}
	for _, v := range from {
		to = append(to, convertStatus(v))
	}
	return to
}

func convertStatus(from *status) *scm.Status {
	return &scm.Status{
		State:  convertState(from.State),
		Label:  from.Context,
		Desc:   from.Description,
		Target: from.TargetURL,
	}
}

func convertDeployStatus(from *deployStatus) *scm.DeployStatus {
	return &scm.DeployStatus{
		Number:         from.ID,
		State:          convertState(from.State),
		Desc:           from.Description,
		Target:         from.TargetURL,
		Environment:    from.Environment,
		EnvironmentURL: from.EnvironmentURL,
	}
}

func convertState(from string) scm.State {
	switch from {
	case "error":
		return scm.StateError
	case "failure":
		return scm.StateFailure
	case "pending":
		return scm.StatePending
	case "success":
		return scm.StateSuccess
	default:
		return scm.StateUnknown
	}
}

func convertFromState(from scm.State) string {
	switch from {
	case scm.StatePending, scm.StateRunning:
		return "pending"
	case scm.StateSuccess:
		return "success"
	case scm.StateFailure:
		return "failure"
	default:
		return "error"
	}
}
