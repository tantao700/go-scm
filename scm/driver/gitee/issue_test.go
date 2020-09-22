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

func TestIssueFind(t *testing.T) {
	defer gock.Off()

	gock.New("https://gitee.com").
		Get("/api/v5/repos/mirrors/diaspora/issues/1").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		File("testdata/issue.json")

	client := NewDefault()
	got, res, err := client.Issues.Find(context.Background(), "mirrors/diaspora", 1)
	if err != nil {
		t.Error(err)
		return
	}

	want := new(scm.Issue)
	raw, _ := ioutil.ReadFile("testdata/issue.json.golden")
	json.Unmarshal(raw, want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
}

func TestIssueCommentFind(t *testing.T) {
	defer gock.Off()

	gock.New("https://gitee.com").
		Get("/api/v5/repos/mirrors/diaspora/issues/2/notes/1").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		File("testdata/issue_note.json")

	client := NewDefault()
	got, res, err := client.Issues.FindComment(context.Background(), "mirrors/diaspora", 2, 1)
	if err != nil {
		t.Error(err)
		return
	}

	want := new(scm.Comment)
	raw, _ := ioutil.ReadFile("testdata/issue_note.json.golden")
	json.Unmarshal(raw, want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
}

func TestIssueList(t *testing.T) {
	defer gock.Off()

	gock.New("https://gitee.com").
		Get("/api/v5/repos/mirrors/diaspora/issues").
		MatchParam("page", "1").
		MatchParam("per_page", "30").
		MatchParam("state", "all").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		SetHeaders(mockPageHeaders).
		File("testdata/issues.json")

	client := NewDefault()
	got, res, err := client.Issues.List(context.Background(), "mirrors/diaspora", scm.IssueListOptions{Page: 1, Size: 30, Open: true, Closed: true})
	if err != nil {
		t.Error(err)
		return
	}

	want := []*scm.Issue{}
	raw, _ := ioutil.ReadFile("testdata/issues.json.golden")
	json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
	t.Run("Page", testPage(res))
}

func TestIssueListComments(t *testing.T) {
	defer gock.Off()

	gock.New("https://gitee.com").
		Get("/api/v5/repos/mirrors/diaspora/issues/1/notes").
		MatchParam("page", "1").
		MatchParam("per_page", "30").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		SetHeaders(mockPageHeaders).
		File("testdata/issue_notes.json")

	client := NewDefault()
	got, res, err := client.Issues.ListComments(context.Background(), "mirrors/diaspora", 1, scm.ListOptions{Size: 30, Page: 1})
	if err != nil {
		t.Error(err)
		return
	}

	want := []*scm.Comment{}
	raw, _ := ioutil.ReadFile("testdata/issue_notes.json.golden")
	json.Unmarshal(raw, &want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
	t.Run("Page", testPage(res))
}

func TestIssueCreate(t *testing.T) {
	defer gock.Off()

	gock.New("https://gitee.com").
		Post("/api/v5/repos/mirrors/diaspora/issues").
		Reply(200).
		Type("application/json").
		SetHeaders(mockHeaders).
		File("testdata/issue.json")

	input := scm.IssueInput{
		Title: "Found a bug",
		Body:  "I'm having a problem with this.",
	}

	client := NewDefault()
	got, res, err := client.Issues.Create(context.Background(), "mirrors/diaspora", &input)
	if err != nil {
		t.Error(err)
		return
	}

	want := new(scm.Issue)
	raw, _ := ioutil.ReadFile("testdata/issue.json.golden")
	json.Unmarshal(raw, want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
}

func TestIssueCreateComment(t *testing.T) {
	defer gock.Off()

	gock.New("https://gitee.com").
		Post("/api/v5/repos/mirrors/diaspora/issues/1/notes").
		MatchParam("body", "lgtm").
		Reply(201).
		Type("application/json").
		SetHeaders(mockHeaders).
		File("testdata/issue_note.json")

	input := &scm.CommentInput{
		Body: "lgtm",
	}

	client := NewDefault()
	got, res, err := client.Issues.CreateComment(context.Background(), "mirrors/diaspora", 1, input)
	if err != nil {
		t.Error(err)
		return
	}

	want := new(scm.Comment)
	raw, _ := ioutil.ReadFile("testdata/issue_note.json.golden")
	json.Unmarshal(raw, want)

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Unexpected Results")
		t.Log(diff)
	}

	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
}

func TestIssueCommentDelete(t *testing.T) {
	defer gock.Off()

	gock.New("https://gitee.com").
		Delete("/api/v5/repos/mirrors/diaspora/issues/2/notes/1").
		Reply(204).
		Type("application/json").
		SetHeaders(mockHeaders)

	client := NewDefault()
	res, err := client.Issues.DeleteComment(context.Background(), "mirrors/diaspora", 2, 1)
	if err != nil {
		t.Error(err)
		return
	}

	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
}

func TestIssueClose(t *testing.T) {
	defer gock.Off()

	gock.New("https://gitee.com").
		Put("/api/v5/repos/mirrors/diaspora/issues/1").
		MatchParam("state_event", "close").
		Reply(204).
		Type("application/json").
		SetHeaders(mockHeaders)

	client := NewDefault()
	res, err := client.Issues.Close(context.Background(), "mirrors/diaspora", 1)
	if err != nil {
		t.Error(err)
		return
	}

	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
}

func TestIssueLock(t *testing.T) {
	defer gock.Off()

	gock.New("https://gitee.com").
		Put("/api/v5/repos/mirrors/diaspora/issues/1").
		MatchParam("discussion_locked", "true").
		Reply(204).
		Type("application/json").
		SetHeaders(mockHeaders)

	client := NewDefault()
	res, err := client.Issues.Lock(context.Background(), "mirrors/diaspora", 1)
	if err != nil {
		t.Error(err)
		return
	}

	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
}

func TestIssueUnlock(t *testing.T) {
	defer gock.Off()

	gock.New("https://gitee.com").
		Put("/api/v5/repos/mirrors/diaspora/issues/1").
		MatchParam("discussion_locked", "false").
		Reply(204).
		Type("application/json").
		SetHeaders(mockHeaders)

	client := NewDefault()
	res, err := client.Issues.Unlock(context.Background(), "mirrors/diaspora", 1)
	if err != nil {
		t.Error(err)
		return
	}

	t.Run("Request", testRequest(res))
	t.Run("Rate", testRate(res))
}
