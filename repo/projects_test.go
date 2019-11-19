package repo

import (
	"fmt"
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
		github.MockListProjects(func(path string, _ *github.ListProjectsParams) github.Projects {
			assert.Contains(t, path, "/repos/example/test/projects")
			assert.Contains(t, path, "state=all")
			assert.Contains(t, path, "per_page=100")
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
		github.MockCreateProject(func(path string, params *github.CreateProjectParams) (*github.Project, error) {
			assert.Equal(t, path, "/repos/example/test/projects")
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
		github.MockUpdateProject(func(path string, projectID int, params *github.UpdateProjectParams) (*github.Project, error) {
			assert.Equal(t, path, "/projects/"+fmt.Sprint(projectID))
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
