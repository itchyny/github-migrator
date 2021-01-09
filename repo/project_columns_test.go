package repo

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/itchyny/github-migrator/github"
)

func TestRepoListProjectColumns(t *testing.T) {
	expected := []*github.ProjectColumn{
		{
			ID:   1,
			Name: "Test project column 1",
		},
		{
			ID:   2,
			Name: "Test project column 2",
		},
	}
	repo := New(github.NewMockClient(
		github.MockListProjectColumns(func(int) github.ProjectColumns {
			return github.ProjectColumnsFromSlice(expected)
		}),
	), "example/test")
	got, err := github.ProjectColumnsToSlice(repo.ListProjectColumns(1))
	assert.Nil(t, err)
	assert.Equal(t, got, expected)
}

func TestRepoGetProjectColumn(t *testing.T) {
	expected := &github.ProjectColumn{
		ID:   1,
		Name: "Test project column 1",
	}
	repo := New(github.NewMockClient(
		github.MockGetProjectColumn(func(int) (*github.ProjectColumn, error) {
			return expected, nil
		}),
	), "example/test")
	got, err := repo.GetProjectColumn(1)
	assert.Nil(t, err)
	assert.Equal(t, got, expected)
}

func TestRepoCreateProjectColumn(t *testing.T) {
	expected := &github.ProjectColumn{
		ID:   1,
		Name: "Test project column 1",
	}
	repo := New(github.NewMockClient(
		github.MockCreateProjectColumn(func(int, string) (*github.ProjectColumn, error) {
			return expected, nil
		}),
	), "example/test")
	got, err := repo.CreateProjectColumn(10, "Test project column 1")
	assert.Nil(t, err)
	assert.Equal(t, got, expected)
}

func TestRepoUpdateProjectColumn(t *testing.T) {
	expected := &github.ProjectColumn{
		ID:   1,
		Name: "Test project column 1",
	}
	repo := New(github.NewMockClient(
		github.MockUpdateProjectColumn(func(int, string) (*github.ProjectColumn, error) {
			return expected, nil
		}),
	), "example/test")
	got, err := repo.UpdateProjectColumn(1, "Test project column 1")
	assert.Nil(t, err)
	assert.Equal(t, got, expected)
}
