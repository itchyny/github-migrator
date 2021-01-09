package repo

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/itchyny/github-migrator/github"
)

func TestRepoListHooks(t *testing.T) {
	expected := []*github.Hook{
		{
			ID:   10,
			Name: "Test hook 1",
		},
		{
			ID:   10,
			Name: "Test hook 1",
		},
	}
	repo := New(github.NewMockClient(
		github.MockListHooks(func(string) github.Hooks {
			return github.HooksFromSlice(expected)
		}),
	), "example/test")
	got, err := github.HooksToSlice(repo.ListHooks())
	assert.Nil(t, err)
	assert.Equal(t, got, expected)
}

func TestRepoGetHook(t *testing.T) {
	expected := &github.Hook{
		ID:   1,
		Name: "Test hook 1",
	}
	repo := New(github.NewMockClient(
		github.MockGetHook(func(string, int) (*github.Hook, error) {
			return expected, nil
		}),
	), "example/test")
	got, err := repo.GetHook(1)
	assert.Nil(t, err)
	assert.Equal(t, got, expected)
}

func TestRepoCreateHook(t *testing.T) {
	expected := &github.Hook{
		ID:     1,
		Name:   "Test hook 1",
		Active: true,
	}
	repo := New(github.NewMockClient(
		github.MockCreateHook(func(string, *github.CreateHookParams) (*github.Hook, error) {
			return expected, nil
		}),
	), "example/test")
	got, err := repo.CreateHook(&github.CreateHookParams{
		Active: true,
	})
	assert.Nil(t, err)
	assert.Equal(t, got, expected)
}

func TestRepoUpdateHook(t *testing.T) {
	expected := &github.Hook{
		ID:     1,
		Name:   "Test hook 1",
		Active: true,
	}
	repo := New(github.NewMockClient(
		github.MockUpdateHook(func(string, int, *github.UpdateHookParams) (*github.Hook, error) {
			return expected, nil
		}),
	), "example/test")
	got, err := repo.UpdateHook(1, &github.UpdateHookParams{
		Active: true,
	})
	assert.Nil(t, err)
	assert.Equal(t, got, expected)
}
