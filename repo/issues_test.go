package repo

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/itchyny/github-migrator/github"
)

func TestRepoListIssues(t *testing.T) {
	expected := []*github.Issue{
		{
			Number:  1,
			Title:   "Example title 1",
			State:   "closed",
			Body:    "Example body 1",
			HTMLURL: "https://github.com/example/test/issues/1",
		},
		{
			Number:  2,
			Title:   "Example title 2",
			State:   "open",
			Body:    "Example body 2",
			HTMLURL: "https://github.com/example/test/issues/2",
		},
	}
	repo := New(github.NewMockClient(
		github.MockListIssues(func(path string, _ *github.ListIssuesParams) ([]*github.Issue, error) {
			assert.Contains(t, path, "/repos/example/test/issues")
			assert.Contains(t, path, "filter=all")
			assert.Contains(t, path, "state=all")
			assert.Contains(t, path, "direction=asc")
			assert.Contains(t, path, "per_page=100")
			return expected, nil
		}),
	), "example/test")
	got, err := repo.ListIssues()
	assert.Nil(t, err)
	assert.Equal(t, got, expected)
}
