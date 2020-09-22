// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gitee

import (
	"context"
	"fmt"
	"time"

	"github.com/drone/go-scm/scm"
)

type gitService struct {
	client *wrapper
}

func (s *gitService) FindBranch(ctx context.Context, repo, name string) (*scm.Reference, *scm.Response, error) {
	path := fmt.Sprintf("api/v5/repos/%s/branches/%s", repo, name)
	out := new(branch)
	res, err := s.client.do(ctx, "GET", path, nil, out)
	return convertBranch(out), res, err
}

func (s *gitService) FindCommit(ctx context.Context, repo, ref string) (*scm.Commit, *scm.Response, error) {
	path := fmt.Sprintf("api/v5/repos/%s/commits/%s", repo, scm.TrimRef(ref))
	out := new(commit)
	res, err := s.client.do(ctx, "GET", path, nil, out)
	return convertCommit(out), res, err
}

func (s *gitService) FindTag(ctx context.Context, repo, name string) (*scm.Reference, *scm.Response, error) {
	path := fmt.Sprintf("api/v5/repos/%s/tags/%s", repo, name)
	out := new(tag)
	res, err := s.client.do(ctx, "GET", path, nil, out)
	return convertTag(out), res, err
}

func (s *gitService) ListBranches(ctx context.Context, repo string, opts scm.ListOptions) ([]*scm.Reference, *scm.Response, error) {
	path := fmt.Sprintf("api/v5/repos/%s/branches?%s", repo, encodeListOptions(opts))
	out := []*branch{}
	res, err := s.client.do(ctx, "GET", path, nil, &out)
	return convertBranchList(out), res, err
}

func (s *gitService) ListCommits(ctx context.Context, repo string, opts scm.CommitListOptions) ([]*scm.Commit, *scm.Response, error) {
	path := fmt.Sprintf("api/v5/repos/%s/commits?%s", repo, encodeCommitListOptions(opts))
	out := []*commit{}
	res, err := s.client.do(ctx, "GET", path, nil, &out)
	return convertCommitList(out), res, err
}

func (s *gitService) ListTags(ctx context.Context, repo string, opts scm.ListOptions) ([]*scm.Reference, *scm.Response, error) {
	path := fmt.Sprintf("api/v5/repos/%s/tags?%s", repo, encodeListOptions(opts))
	out := []*tag{}
	res, err := s.client.do(ctx, "GET", path, nil, &out)
	return convertTagList(out), res, err
}

func (s *gitService) ListChanges(ctx context.Context, repo, ref string, opts scm.ListOptions) ([]*scm.Change, *scm.Response, error) {
	path := fmt.Sprintf("api/v5/repos/%s/commits/%s/diff", repo, ref)
	out := []*change{}
	res, err := s.client.do(ctx, "GET", path, nil, &out)
	return convertChangeList(out), res, err
}

func (s *gitService) CompareChanges(ctx context.Context, repo, source, target string, _ scm.ListOptions) ([]*scm.Change, *scm.Response, error) {
	path := fmt.Sprintf("api/v5/repos/%s/repository/compare?from=%s&to=%s", repo, source, target)
	out := new(compare)
	res, err := s.client.do(ctx, "GET", path, nil, &out)
	return convertChangeList(out.Diffs), res, err
}

type branch struct {
	Name   string `json:"name"`
	Commit struct {
		ID  string `json:"id"`
		Sha string `json:"sha"`
	}
}

type tag struct {
	Name    string `json:"name"`
	Message string `json:"message"`
	Commit  struct {
		Sha  string `json:"sha"`
		Date string `json:"date"`
	}
}

type commit struct {
	ID             string    `json:"id"`
	Title          string    `json:"title"`
	Message        string    `json:"message"`
	AuthorName     string    `json:"author_name"`
	AuthorEmail    string    `json:"author_email"`
	AuthorDate     time.Time `json:"authored_date"`
	CommittedDate  time.Time `json:"committed_date"`
	CommitterName  string    `json:"committer_name"`
	CommitterEmail string    `json:"committer_email"`
	Created        time.Time `json:"created_at"`
}

type compare struct {
	Diffs []*change `json:"diffs"`
}

func convertCommitList(from []*commit) []*scm.Commit {
	to := []*scm.Commit{}
	for _, v := range from {
		to = append(to, convertCommit(v))
	}
	return to
}

func convertCommit(from *commit) *scm.Commit {
	return &scm.Commit{
		Message: from.Message,
		Sha:     from.ID,
		Author: scm.Signature{
			Login: from.AuthorName,
			Name:  from.AuthorName,
			Email: from.AuthorEmail,
			Date:  from.AuthorDate,
		},
		Committer: scm.Signature{
			Login: from.CommitterName,
			Name:  from.CommitterName,
			Email: from.CommitterEmail,
			Date:  from.CommittedDate,
		},
	}
}

func convertBranchList(from []*branch) []*scm.Reference {
	to := []*scm.Reference{}
	for _, v := range from {
		to = append(to, convertBranch(v))
	}
	return to
}

func convertBranch(from *branch) *scm.Reference {
	return &scm.Reference{
		Name: scm.TrimRef(from.Name),
		Path: scm.ExpandRef(from.Name, "refs/heads/"),
		Sha:  from.Commit.ID,
	}
}

func convertTagList(from []*tag) []*scm.Reference {
	to := []*scm.Reference{}
	for _, v := range from {
		to = append(to, convertTag(v))
	}
	return to
}

func convertTag(from *tag) *scm.Reference {
	return &scm.Reference{
		Name: scm.TrimRef(from.Name),
		Path: scm.ExpandRef(from.Name, "refs/tags/"),
		Sha:  from.Commit.Sha,
	}
}
