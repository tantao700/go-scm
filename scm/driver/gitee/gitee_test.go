// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gitee

import (
	"net/url"
	"testing"

	"github.com/drone/go-scm/scm"
)

var mockHeaders = map[string]string{
	"X-GitHub-Request-Id":   "DD0E:6011:12F21A8:1926790:5A2064E2",
	"X-RateLimit-Limit":     "60",
	"X-RateLimit-Remaining": "59",
	"X-RateLimit-Reset":     "1512076018",
}

var mockPageHeaders = map[string]string{
	"Link": `<https://gitee.com/api/v5/resource?page=2>; rel="next",` +
		`<https://gitee.com/api/v5/resource?page=1>; rel="prev",` +
		`<https://gitee.com/api/v5/resource?page=1>; rel="first",` +
		`<https://gitee.com/api/v5/resource?page=5>; rel="last"`,
}

func TestClient(t *testing.T) {
	client, err := New("https://gitee.com/api/v5")
	if err != nil {
		t.Error(err)
	}
	if got, want := client.BaseURL.String(), "https://gitee.com/api/v5/"; got != want {
		t.Errorf("Want Client URL %q, got %q", want, got)
	}
}

func TestClient_Base(t *testing.T) {
	client, err := New("https://github.example.com/api/v3")
	if err != nil {
		t.Error(err)
	}
	if got, want := client.BaseURL.String(), "https://github.example.com/api/v3/"; got != want {
		t.Errorf("Want Client URL %q, got %q", want, got)
	}
}

func TestClient_Default(t *testing.T) {
	client := NewDefault()
	if got, want := client.BaseURL.String(), "https://gitee.com/api/v5/"; got != want {
		t.Errorf("Want Client URL %q, got %q", want, got)
	}
}

func TestClient_Error(t *testing.T) {
	_, err := New("http://a b.com/")
	if err == nil {
		t.Errorf("Expect error when invalid URL")
	}
}

func testRate(res *scm.Response) func(t *testing.T) {
	return func(t *testing.T) {
		if got, want := res.Rate.Limit, 60; got != want {
			t.Errorf("Want X-RateLimit-Limit %d, got %d", want, got)
		}
		if got, want := res.Rate.Remaining, 59; got != want {
			t.Errorf("Want X-RateLimit-Remaining %d, got %d", want, got)
		}
		if got, want := res.Rate.Reset, int64(1512076018); got != want {
			t.Errorf("Want X-RateLimit-Reset %d, got %d", want, got)
		}
	}
}

func testPage(res *scm.Response) func(t *testing.T) {
	return func(t *testing.T) {
		if got, want := res.Page.Next, 2; got != want {
			t.Errorf("Want next page %d, got %d", want, got)
		}
		if got, want := res.Page.Prev, 1; got != want {
			t.Errorf("Want prev page %d, got %d", want, got)
		}
		if got, want := res.Page.First, 1; got != want {
			t.Errorf("Want first page %d, got %d", want, got)
		}
		if got, want := res.Page.Last, 5; got != want {
			t.Errorf("Want last page %d, got %d", want, got)
		}
	}
}

func testRequest(res *scm.Response) func(t *testing.T) {
	return func(t *testing.T) {
		if got, want := res.ID, "DD0E:6011:12F21A8:1926790:5A2064E2"; got != want {
			t.Errorf("Want X-GitHub-Request-Id %q, got %q", want, got)
		}
	}
}

func TestWebsiteAddress(t *testing.T) {
	tests := []struct {
		api string
		web string
	}{
		{"https://gitee.com/api/v5/", "https://gitee.com/"},
		{"https://gitee.com/api/v5", "https://gitee.com/"},
		{"https://github.acme.com/api/v3", "https://github.acme.com/"},
		{"https://github.acme.com/api/v3/", "https://github.acme.com/"},
	}

	for _, test := range tests {
		parsed, _ := url.Parse(test.api)
		got, want := websiteAddress(parsed), test.web
		if got != want {
			t.Errorf("Want website address %q, got %q", want, got)
		}
	}
}
