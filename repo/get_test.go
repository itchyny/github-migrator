package repo

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/itchyny/github-migrator/github"
)

func TestRepoGet(t *testing.T) {
	expected := &github.Repo{
		Name:        "test",
		Description: "Test repository.",
		HTMLURL:     "http://localhost/example/test",
	}
	repo := New(github.NewMockClient(
		github.MockGetRepo(func(path string) (*github.Repo, error) {
			assert.Contains(t, path, "/repos/example/test")
			return expected, nil
		}),
	), "example/test")
	got, err := repo.Get()
	assert.Nil(t, err)
	assert.Equal(t, got, expected)
}
