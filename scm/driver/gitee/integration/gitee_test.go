// Copyright 2017 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package integration

import (
	"net/http"
	"os"
	"testing"

	"github.com/drone/go-scm/scm/driver/gitee"
	"github.com/drone/go-scm/scm/transport"
)

func Testgitee(t *testing.T) {
	if os.Getenv("gitee_TOKEN") == "" {
		t.Skipf("missing gitee_TOKEN environment variable")
		return
	}

	client, _ := gitee.New("https://gitee.com/")
	client.Client = &http.Client{
		Transport: &transport.PrivateToken{
			Token: os.Getenv("gitee_TOKEN"),
		},
	}

	t.Run("Contents", testContents(client))
	t.Run("Git", testGit(client))
	t.Run("Issues", testIssues(client))
	t.Run("Organizations", testOrgs(client))
	t.Run("PullRequests", testPullRequests(client))
	t.Run("Repositories", testRepos(client))
	t.Run("Users", testUsers(client))
}
