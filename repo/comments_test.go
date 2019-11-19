package repo

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/itchyny/github-migrator/github"
)

func TestRepoListComments(t *testing.T) {
	expected := []*github.Comment{
		{
			Body:    "Example body 1",
			HTMLURL: "http://localhost/example/test/issues/1#issuecomment-1",
		},
		{
			Body:    "Example body 2",
			HTMLURL: "http://localhost/example/test/issues/1#issuecomment-2",
		},
	}
	repo := New(github.NewMockClient(
		github.MockListComments(func(string, int) github.Comments {
			return github.CommentsFromSlice(expected)
		}),
	), "example/test")
	got, err := github.CommentsToSlice(repo.ListComments(1))
	assert.Nil(t, err)
	assert.Equal(t, got, expected)
}
