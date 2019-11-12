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
			State:   github.IssueStateClosed,
			Body:    "Example body 1",
			HTMLURL: "http://localhost/example/test/issues/1",
		},
		{
			Number:  2,
			Title:   "Example title 2",
			State:   github.IssueStateOpen,
			Body:    "Example body 2",
			HTMLURL: "http://localhost/example/test/issues/2",
		},
	}
	repo := New(github.NewMockClient(
		github.MockListIssues(func(path string, _ *github.ListIssuesParams) github.Issues {
			assert.Contains(t, path, "/repos/example/test/issues")
			assert.Contains(t, path, "filter=all")
			assert.Contains(t, path, "state=all")
			assert.Contains(t, path, "direction=asc")
			assert.Contains(t, path, "per_page=100")
			return github.IssuesFromSlice(expected)
		}),
	), "example/test")
	got, err := github.IssuesToSlice(repo.ListIssues())
	assert.Nil(t, err)
	assert.Equal(t, got, expected)
}

func TestRepoGetIssue(t *testing.T) {
	expected := &github.Issue{
		Number:  1,
		Title:   "Example title 1",
		State:   github.IssueStateClosed,
		Body:    "Example body 1",
		HTMLURL: "http://localhost/example/test/issue/1",
		ClosedBy: &github.User{
			Login: "test-user",
		},
	}
	repo := New(github.NewMockClient(
		github.MockGetIssue(func(path string, issueNumber int) (*github.Issue, error) {
			assert.Contains(t, path, "/repos/example/test/issues/1")
			return expected, nil
		}),
	), "example/test")
	got, err := repo.GetIssue(1)
	assert.Nil(t, err)
	assert.Equal(t, got, expected)
}
