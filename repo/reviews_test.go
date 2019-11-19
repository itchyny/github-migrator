package repo

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/itchyny/github-migrator/github"
)

func TestRepoListReviews(t *testing.T) {
	expected := []*github.Review{
		{
			ID:    1,
			State: github.ReviewStateApproved,
			Body:  "Example body 1",
		},
		{
			ID:    2,
			State: github.ReviewStateChangesRequested,
			Body:  "Example body 2",
		},
	}
	repo := New(github.NewMockClient(
		github.MockListReviews(func(_ string, pullNumber int) github.Reviews {
			assert.Equal(t, pullNumber, 1)
			return github.ReviewsFromSlice(expected)
		}),
	), "example/test")
	got, err := github.ReviewsToSlice(repo.ListReviews(1))
	assert.Nil(t, err)
	assert.Equal(t, got, expected)
}
