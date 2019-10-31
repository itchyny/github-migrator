package repo

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/itchyny/github-migrator/github"
)

func TestRepoListReviewComments(t *testing.T) {
	expected := []*github.ReviewComment{
		{
			Body: "Example body 1",
		},
		{
			Body: "Example body 2",
		},
	}
	repo := New(github.NewMockClient(
		github.MockListReviewComments(func(path string, pullNumber int) github.ReviewComments {
			assert.Contains(t, path, "/repos/example/test/pulls/1/comments")
			assert.Equal(t, pullNumber, 1)
			return github.ReviewCommentsFromSlice(expected)
		}),
	), "example/test")
	got, err := github.ReviewCommentsToSlice(repo.ListReviewComments(1))
	assert.Nil(t, err)
	assert.Equal(t, got, expected)
}
