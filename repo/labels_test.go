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
		github.MockListLabels(func(string) github.Labels {
			return github.LabelsFromSlice(expected)
		}),
	), "example/test")
	got, err := github.LabelsToSlice(repo.ListLabels())
	assert.Nil(t, err)
	assert.Equal(t, got, expected)
}

func TestRepoCreateLabel(t *testing.T) {
	expected := &github.Label{
		ID:          1,
		Name:        "bug",
		Description: "This is a bug.",
		Color:       "fc2929",
		Default:     false,
	}
	repo := New(github.NewMockClient(
		github.MockCreateLabel(func(string, *github.CreateLabelParams) (*github.Label, error) {
			return expected, nil
		}),
	), "example/test")
	got, err := repo.CreateLabel(&github.CreateLabelParams{
		Name:        "bug",
		Description: "This is a bug.",
		Color:       "fc2929",
	})
	assert.Nil(t, err)
	assert.Equal(t, got, expected)
}

func TestRepoUpdateLabel(t *testing.T) {
	expected := &github.Label{
		ID:          1,
		Name:        "warn",
		Description: "This is a warning.",
		Color:       "fcfc29",
		Default:     false,
	}
	repo := New(github.NewMockClient(
		github.MockUpdateLabel(func(string, string, *github.UpdateLabelParams) (*github.Label, error) {
			return expected, nil
		}),
	), "example/test")
	got, err := repo.UpdateLabel("bug", &github.UpdateLabelParams{
		Name:        "warn",
		Description: "This is a warning.",
		Color:       "fcfc29",
	})
	assert.Nil(t, err)
	assert.Equal(t, got, expected)
}
