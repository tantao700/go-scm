// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gitee

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/drone/go-scm/scm"

	"github.com/google/go-cmp/cmp"
	"github.com/h2non/gock"
)

// TODO(bradrydzewski) missing commit link
// TODO(bradrydzewski) missing commit author avatar
// TODO(bradrydzewski) missing commit committer avatar

func TestGitFindCommit(t *testing.T) {
	defer gock.Off()

	gock.New("https://gitee.com").
		Get("/api/v5/repos/mirrors/diaspora/commits/7fd1a60b01f91b314f59955a4e4d4e80d8edf11d").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		File("testdata/commit.json")

	client := NewDefault()
	got, res, err := client.Git.FindCommit(context.Background(), "mirrors/diaspora", "7fd1a60b01f91b314f59955a4e4d4e80d8edf11d")
	if err != nil {
		t.Error(err)
		return
	}

	want := new(scm.Commit)
	raw, _ := ioutil.ReadFile("testdata/commit.json.golden")
	json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
}

func TestGitFindBranch(t *testing.T) {
	defer gock.Off()

	gock.New("https://gitee.com").
		Get("/api/v5/repos/mirrors/diaspora/branches/master").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		File("testdata/branch.json")

	client := NewDefault()
	got, res, err := client.Git.FindBranch(context.Background(), "mirrors/diaspora", "master")
	if err != nil {
		t.Error(err)
		return
	}

	want := new(scm.Reference)
	raw, _ := ioutil.ReadFile("testdata/branch.json.golden")
	json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
}

func TestGitFindTag(t *testing.T) {
	defer gock.Off()

	gock.New("https://gitee.com").
		Get("/api/v5/repos/mirrors/diaspora/tags/v0.1.0.0").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		File("testdata/tag.json")

	client := NewDefault()
	got, res, err := client.Git.FindTag(context.Background(), "mirrors/diaspora", "v0.1.0.0")
	if err != nil {
		t.Error(err)
		return
	}

	want := new(scm.Reference)
	raw, _ := ioutil.ReadFile("testdata/tag.json.golden")
	json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
}

func TestGitListCommits(t *testing.T) {
	defer gock.Off()

	gock.New("https://gitee.com").
		Get("api/v5/repos/mirrors/diaspora/commits").
		MatchParam("page", "1").
		MatchParam("per_page", "30").
		MatchParam("ref_name", "master").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		SetHeaders(mockPageHeaders).
		File("testdata/commits.json")

	client := NewDefault()
	got, res, err := client.Git.ListCommits(context.Background(), "mirrors/diaspora", scm.CommitListOptions{Ref: "master", Page: 1, Size: 30})
	if err != nil {
		t.Error(err)
		return
	}

	want := []*scm.Commit{}
	raw, _ := ioutil.ReadFile("testdata/commits.json.golden")
	json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
	t.Run("Page", testPage(res))
}

func TestGitListBranches(t *testing.T) {
	defer gock.Off()

	gock.New("https://gitee.com").
		Get("/api/v5/repos/mirrors/diaspora/branches").
		MatchParam("page", "1").
		MatchParam("per_page", "30").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		SetHeaders(mockPageHeaders).
		File("testdata/branches.json")

	client := NewDefault()
	got, res, err := client.Git.ListBranches(context.Background(), "mirrors/diaspora", scm.ListOptions{Page: 1, Size: 30})
	if err != nil {
		t.Error(err)
		return
	}

	want := []*scm.Reference{}
	raw, _ := ioutil.ReadFile("testdata/branches.json.golden")
	json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
	t.Run("Page", testPage(res))
}

func TestGitListTags(t *testing.T) {
	defer gock.Off()

	gock.New("https://gitee.com").
		Get("/api/v5/repos/mirrors/diaspora/repository/tags").
		MatchParam("page", "1").
		MatchParam("per_page", "30").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		SetHeaders(mockPageHeaders).
		File("testdata/tags.json")

	client := NewDefault()
	got, res, err := client.Git.ListTags(context.Background(), "mirrors/diaspora", scm.ListOptions{Page: 1, Size: 30})
	if err != nil {
		t.Error(err)
		return
	}

	want := []*scm.Reference{}
	raw, _ := ioutil.ReadFile("testdata/tags.json.golden")
	json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
	t.Run("Page", testPage(res))
}

func TestGitListChanges(t *testing.T) {
	defer gock.Off()

	gock.New("https://gitee.com").
		Get("/api/v5/repos/mirrors/diaspora/commits/6104942438c14ec7bd21c6cd5bd995272b3faff6/diff").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		File("testdata/commit_diff.json")

	client := NewDefault()
	got, res, err := client.Git.ListChanges(context.Background(), "mirrors/diaspora", "6104942438c14ec7bd21c6cd5bd995272b3faff6", scm.ListOptions{})
	if err != nil {
		t.Error(err)
		return
	}

	want := []*scm.Change{}
	raw, _ := ioutil.ReadFile("testdata/commit_diff.json.golden")
	json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
}

func TestGitCompareChanges(t *testing.T) {
	defer gock.Off()

	gock.New("https://gitee.com").
		Get("/api/v5/repos/mirrors/diaspora/repository/compare").
		MatchParam("from", "ae1d9fb46aa2b07ee9836d49862ec4e2c46fbbba").
		MatchParam("to", "6104942438c14ec7bd21c6cd5bd995272b3faff6").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		File("testdata/compare.json")

	client := NewDefault()
	got, res, err := client.Git.CompareChanges(context.Background(), "mirrors/diaspora", "ae1d9fb46aa2b07ee9836d49862ec4e2c46fbbba", "6104942438c14ec7bd21c6cd5bd995272b3faff6", scm.ListOptions{})
	if err != nil {
		t.Error(err)
		return
	}

	want := []*scm.Change{}
	raw, _ := ioutil.ReadFile("testdata/compare.json.golden")
	json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
}
