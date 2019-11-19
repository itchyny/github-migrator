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
				State:   github.IssueStateClosed,
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
				State:   github.IssueStateOpen,
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
		github.MockListPullReqs(func(string, *github.ListPullReqsParams) github.PullReqs {
			return github.PullReqsFromSlice(expected)
		}),
	), "example/test")
	got, err := github.PullReqsToSlice(repo.ListPullReqs())
	assert.Nil(t, err)
	assert.Equal(t, got, expected)
}

func TestRepoGetPullReq(t *testing.T) {
	expected := &github.PullReq{
		Issue: github.Issue{
			Number:  1,
			Title:   "Example title 1",
			State:   github.IssueStateClosed,
			Body:    "Example body 1",
			HTMLURL: "http://localhost/example/test/pull/1",
		},
		Merged: false,
		Draft:  true,
	}
	repo := New(github.NewMockClient(
		github.MockGetPullReq(func(string, int) (*github.PullReq, error) {
			return expected, nil
		}),
	), "example/test")
	got, err := repo.GetPullReq(1)
	assert.Nil(t, err)
	assert.Equal(t, got, expected)
}
