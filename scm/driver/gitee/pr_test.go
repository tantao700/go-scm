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

func TestPullFind(t *testing.T) {
	defer gock.Off()

	gock.New("https://gitee.com/api/v5").
		Get("/repos/tantao700/hello-world/pulls/1347").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		File("testdata/pr.json")

	client := NewDefault()
	got, res, err := client.PullRequests.Find(context.Background(), "tantao700/hello-world", 1347)
	if err != nil {
		t.Error(err)
		return
	}

	want := new(scm.PullRequest)
	raw, _ := ioutil.ReadFile("testdata/pr.json.golden")
	json.Unmarshal(raw, want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
}

func TestPullList(t *testing.T) {
	defer gock.Off()

	gock.New("https://gitee.com/api/v5").
		Get("/repos/tantao700/hello-world/pulls").
		MatchParam("page", "1").
		MatchParam("per_page", "30").
		MatchParam("state", "all").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		SetHeaders(mockPageHeaders).
		File("testdata/pulls.json")

	client := NewDefault()
	got, res, err := client.PullRequests.List(context.Background(), "tantao700/hello-world", scm.PullRequestListOptions{Page: 1, Size: 30, Open: true, Closed: true})
	if err != nil {
		t.Error(err)
		return
	}

	want := []*scm.PullRequest{}
	raw, _ := ioutil.ReadFile("testdata/pulls.json.golden")
	json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
	t.Run("Page", testPage(res))
}

func TestPullListChanges(t *testing.T) {
	gock.New("https://gitee.com/api/v5").
		Get("/repos/tantao700/hello-world/pulls/1347/files").
		MatchParam("page", "1").
		MatchParam("per_page", "30").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		SetHeaders(mockPageHeaders).
		File("testdata/pr_files.json")

	client := NewDefault()
	got, res, err := client.PullRequests.ListChanges(context.Background(), "tantao700/hello-world", 1347, scm.ListOptions{Page: 1, Size: 30})
	if err != nil {
		t.Error(err)
		return
	}

	want := []*scm.Change{}
	raw, _ := ioutil.ReadFile("testdata/pr_files.json.golden")
	json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
	t.Run("Page", testPage(res))
}

func TestPullMerge(t *testing.T) {
	defer gock.Off()

	gock.New("https://gitee.com/api/v5").
		Put("/repos/tantao700/hello-world/pulls/1347/merge").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders)

	client := NewDefault()
	res, err := client.PullRequests.Merge(context.Background(), "tantao700/hello-world", 1347)
	if err != nil {
		t.Error(err)
		return
	}

	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
}

func TestPullClose(t *testing.T) {
	defer gock.Off()

	gock.New("https://gitee.com/api/v5").
		Patch("/repos/tantao700/hello-world/pulls/1347").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders)

	client := NewDefault()
	res, err := client.PullRequests.Close(context.Background(), "tantao700/hello-world", 1347)
	if err != nil {
		t.Error(err)
		return
	}

	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
}

func TestPullCreate(t *testing.T) {
	defer gock.Off()

	gock.New("https://gitee.com/api/v5").
		Post("/repos/tantao700/hello-world/pulls").
		Reply(201).
		Type("application/json").
		SetHeaders(mockHeaders).
		File("testdata/pr.json")

	input := scm.PullRequestInput{
		Title:  "new-feature",
		Body:   "Please pull these awesome changes",
		Source: "new-topic",
		Target: "master",
	}

	client := NewDefault()
	got, res, err := client.PullRequests.Create(context.Background(), "tantao700/hello-world", &input)
	if err != nil {
		t.Error(err)
		return
	}

	want := new(scm.PullRequest)
	raw, _ := ioutil.ReadFile("testdata/pr.json.golden")
	json.Unmarshal(raw, want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	t.Run("Request", testRequest(res))
	t.Run("rate", testRate(res))
}
