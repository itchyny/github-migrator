package repo

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/itchyny/github-migrator/github"
)

func TestRepoListPullReqCommits(t *testing.T) {
	expected := []*github.Commit{
		&github.Commit{
			HTMLURL: "http://localhost/example/test/commit/xxx",
			SHA:     "xxx",
		},
		&github.Commit{
			HTMLURL: "http://localhost/example/test/commit/yyy",
			SHA:     "yyy",
		},
	}
	repo := New(github.NewMockClient(
		github.MockListPullReqCommits(func(path string, pullNumber int) github.Commits {
			assert.Contains(t, path, fmt.Sprintf("/repos/example/test/pulls/%d/commits", pullNumber))
			return github.CommitsFromSlice(expected)
		}),
	), "example/test")
	got, err := github.CommitsToSlice(repo.ListPullReqCommits(10))
	assert.Nil(t, err)
	assert.Equal(t, got, expected)
}
