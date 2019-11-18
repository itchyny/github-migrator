package repo

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/itchyny/github-migrator/github"
)

func TestRepoListHooks(t *testing.T) {
	expected := []*github.Hook{
		&github.Hook{
			ID:   10,
			Name: "Test hook 1",
		},
		&github.Hook{
			ID:   10,
			Name: "Test hook 1",
		},
	}
	repo := New(github.NewMockClient(
		github.MockListHooks(func(path string) github.Hooks {
			assert.Contains(t, path, "/repos/example/test/hooks")
			assert.Contains(t, path, "per_page=100")
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
		github.MockGetHook(func(path string, hookID int) (*github.Hook, error) {
			assert.Contains(t, path, "/repos/example/test/hooks/1")
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
		github.MockCreateHook(func(path string, params *github.CreateHookParams) (*github.Hook, error) {
			assert.Equal(t, path, "/repos/example/test/hooks")
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
		github.MockUpdateHook(func(path string, hookID int, params *github.UpdateHookParams) (*github.Hook, error) {
			assert.Equal(t, path, "/repos/example/test/hooks/"+fmt.Sprint(hookID))
			return expected, nil
		}),
	), "example/test")
	got, err := repo.UpdateHook(1, &github.UpdateHookParams{
		Active: true,
	})
	assert.Nil(t, err)
	assert.Equal(t, got, expected)
}
