package repo

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/itchyny/github-migrator/github"
)

func TestRepoUpdate(t *testing.T) {
	expected := &github.Repo{
		Name:        "test",
		FullName:    "example/test",
		Description: "New description",
		HTMLURL:     "http://localhost/example/test",
		Homepage:    "http://localhost/new",
		Private:     true,
	}
	repo := New(github.NewMockClient(
		github.MockGetRepo(func(path string) (*github.Repo, error) {
			return &github.Repo{
				Name:        "test",
				FullName:    "example/test",
				Description: "Test repository.",
				HTMLURL:     "http://localhost/example/test",
				Homepage:    "http://localhost/",
				Private:     false,
			}, nil
		}),
		github.MockUpdateRepo(func(path string, params *github.UpdateRepoParams) (*github.Repo, error) {
			assert.Equal(t, path, "/repos/example/test")
			assert.Equal(t, params.Name, "test")
			assert.Equal(t, params.Description, "New description")
			assert.Equal(t, params.Homepage, "http://localhost/new")
			assert.Equal(t, params.Private, true)
			return expected, nil
		}),
	), "example/test")
	got, err := repo.Update(&github.UpdateRepoParams{
		Name:        "test",
		Description: "New description",
		Homepage:    "http://localhost/new",
		Private:     true,
	})
	assert.Nil(t, err)
	assert.Equal(t, got, expected)
}
