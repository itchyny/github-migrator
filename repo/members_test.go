package repo

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/itchyny/github-migrator/github"
)

func TestListMembers(t *testing.T) {
	expected := []*github.Member{
		{
			Login: "user1",
		},
		{
			Login: "user2",
		},
	}
	repo := New(github.NewMockClient(
		github.MockListMembers(func(path string) github.Members {
			assert.Equal(t, path, "/orgs/example/members")
			return github.MembersFromSlice(expected)
		}),
	), "example")
	got, err := github.MembersToSlice(repo.ListMembers())
	assert.Nil(t, err)
	assert.Equal(t, got, expected)
}
