package repo

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/itchyny/github-migrator/github"
)

func TestRepoListPullReqs(t *testing.T) {
	expected := []*github.PullReq{
		&github.PullReq{
			Issue: github.Issue{
				Number:  1,
				Title:   "Example title 1",
				State:   "closed",
				Body:    "Example body 1",
				HTMLURL: "http://localhost/example/test/pull/1",
			},
			Merged: false,
			Draft:  true,
		},
		&github.PullReq{
			Issue: github.Issue{
				Number:  2,
				Title:   "Example title 2",
				State:   "open",
				Body:    "Example body 2",
				HTMLURL: "http://localhost/example/test/pull/2",
			},
			Merged: true,
			MergedBy: &github.User{
				Login: "sample-user-1",
			},
			Draft: false,
		},
	}
	repo := New(github.NewMockClient(
		github.MockListPullReqs(func(path string, _ *github.ListPullReqsParams) github.PullReqs {
			assert.Contains(t, path, "/repos/example/test/pulls")
			assert.Contains(t, path, "state=all")
			assert.Contains(t, path, "direction=asc")
			assert.Contains(t, path, "per_page=100")
			return github.PullReqsFromSlice(expected)
		}),
	), "example/test")
	got, err := github.PullReqsToSlice(repo.ListPullReqs())
	assert.Nil(t, err)
	assert.Equal(t, got, expected)
}
