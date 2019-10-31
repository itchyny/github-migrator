package repo

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/itchyny/github-migrator/github"
)

func TestRepoListLabels(t *testing.T) {
	expected := []*github.Label{
		{
			ID:          1,
			Name:        "bug",
			Description: "This is a bug.",
			Color:       "fc2929",
			Default:     false,
		},
		{
			ID:          2,
			Name:        "design",
			Description: "This is a design issue.",
			Color:       "bfdadc",
			Default:     false,
		},
	}
	repo := New(github.NewMockClient(
		github.MockListLabels(func(path string) github.Labels {
			assert.Equal(t, path, "/repos/example/test/labels")
			return github.LabelsFromSlice(expected)
		}),
	), "example/test")
	got, err := github.LabelsToSlice(repo.ListLabels())
	assert.Nil(t, err)
	assert.Equal(t, got, expected)
}
