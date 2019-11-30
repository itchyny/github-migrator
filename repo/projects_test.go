package repo

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/itchyny/github-migrator/github"
)

func TestRepoListProjects(t *testing.T) {
	expected := []*github.Project{
		&github.Project{
			ID:   10,
			Name: "Test project 1",
		},
		&github.Project{
			ID:   10,
			Name: "Test project 1",
		},
	}
	repo := New(github.NewMockClient(
		github.MockListProjects(func(string, *github.ListProjectsParams) github.Projects {
			return github.ProjectsFromSlice(expected)
		}),
	), "example/test")
	got, err := github.ProjectsToSlice(repo.ListProjects())
	assert.Nil(t, err)
	assert.Equal(t, got, expected)
}

func TestRepoGetProject(t *testing.T) {
	expected := &github.Project{
		ID:   1,
		Name: "Test project 1",
	}
	repo := New(github.NewMockClient(
		github.MockGetProject(func(projectID int) (*github.Project, error) {
			return expected, nil
		}),
	), "example/test")
	got, err := repo.GetProject(1)
	assert.Nil(t, err)
	assert.Equal(t, got, expected)
}

func TestRepoCreateProject(t *testing.T) {
	expected := &github.Project{
		ID:    1,
		Name:  "Test project 1",
		Body:  "Test body",
		State: github.ProjectStateClosed,
	}
	repo := New(github.NewMockClient(
		github.MockCreateProject(func(string, *github.CreateProjectParams) (*github.Project, error) {
			return expected, nil
		}),
	), "example/test")
	got, err := repo.CreateProject(&github.CreateProjectParams{
		Name: "Test project 1",
		Body: "Test body",
	})
	assert.Nil(t, err)
	assert.Equal(t, got, expected)
}

func TestRepoUpdateProject(t *testing.T) {
	expected := &github.Project{
		ID:    1,
		Name:  "Test project 1",
		Body:  "Test body",
		State: github.ProjectStateClosed,
	}
	repo := New(github.NewMockClient(
		github.MockUpdateProject(func(projectID int, params *github.UpdateProjectParams) (*github.Project, error) {
			return expected, nil
		}),
	), "example/test")
	got, err := repo.UpdateProject(1, &github.UpdateProjectParams{
		Name:  "Test project 1",
		Body:  "Test body",
		State: github.ProjectStateClosed,
	})
	assert.Nil(t, err)
	assert.Equal(t, got, expected)
}

func TestRepoDeleteProject(t *testing.T) {
	repo := New(github.NewMockClient(
		github.MockDeleteProject(func(int) error {
			return nil
		}),
	), "example/test")
	err := repo.DeleteProject(1)
	assert.Nil(t, err)
}
