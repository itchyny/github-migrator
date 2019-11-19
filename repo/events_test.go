package repo

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/itchyny/github-migrator/github"
)

func TestRepoListEvents(t *testing.T) {
	expected := []*github.Event{
		{
			Actor: &github.User{Login: "user-1"},
			Event: "labeled",
			Label: &github.EventLabel{Name: "label-1"},
		},
		{
			Actor: &github.User{Login: "user-2"},
			Event: "labeled",
			Label: &github.EventLabel{Name: "label-2"},
		},
	}
	repo := New(github.NewMockClient(
		github.MockListEvents(func(string, int) github.Events {
			return github.EventsFromSlice(expected)
		}),
	), "example/test")
	got, err := github.EventsToSlice(repo.ListEvents(1))
	assert.Nil(t, err)
	assert.Equal(t, got, expected)
}
