package repo

import (
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
		github.MockListPullReqCommits(func(string, int) github.Commits {
			return github.CommitsFromSlice(expected)
		}),
	), "example/test")
	got, err := github.CommitsToSlice(repo.ListPullReqCommits(10))
	assert.Nil(t, err)
	assert.Equal(t, got, expected)
}
