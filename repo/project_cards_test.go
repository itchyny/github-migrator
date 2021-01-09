package repo

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/itchyny/github-migrator/github"
)

func TestRepoListProjectCards(t *testing.T) {
	expected := []*github.ProjectCard{
		{
			ID:   1,
			Note: "Test project card 1",
		},
		{
			ID:   2,
			Note: "Test project card 2",
		},
	}
	repo := New(github.NewMockClient(
		github.MockListProjectCards(func(int) github.ProjectCards {
			return github.ProjectCardsFromSlice(expected)
		}),
	), "example/test")
	got, err := github.ProjectCardsToSlice(repo.ListProjectCards(1))
	assert.Nil(t, err)
	assert.Equal(t, got, expected)
}

func TestRepoGetProjectCard(t *testing.T) {
	expected := &github.ProjectCard{
		ID:   1,
		Note: "Test project card 1",
	}
	repo := New(github.NewMockClient(
		github.MockGetProjectCard(func(int) (*github.ProjectCard, error) {
			return expected, nil
		}),
	), "example/test")
	got, err := repo.GetProjectCard(1)
	assert.Nil(t, err)
	assert.Equal(t, got, expected)
}

func TestRepoCreateProjectCard(t *testing.T) {
	expected := &github.ProjectCard{
		ID:   1,
		Note: "Test project card 1",
	}
	repo := New(github.NewMockClient(
		github.MockCreateProjectCard(func(int, *github.CreateProjectCardParams) (*github.ProjectCard, error) {
			return expected, nil
		}),
	), "example/test")
	got, err := repo.CreateProjectCard(10, nil)
	assert.Nil(t, err)
	assert.Equal(t, got, expected)
}

func TestRepoUpdateProjectCard(t *testing.T) {
	expected := &github.ProjectCard{
		ID:   1,
		Note: "Test project card 1",
	}
	repo := New(github.NewMockClient(
		github.MockUpdateProjectCard(func(int, *github.UpdateProjectCardParams) (*github.ProjectCard, error) {
			return expected, nil
		}),
	), "example/test")
	got, err := repo.UpdateProjectCard(1, nil)
	assert.Nil(t, err)
	assert.Equal(t, got, expected)
}

func TestRepoMoveProjectCard(t *testing.T) {
	expected := &github.ProjectCard{
		ID:   1,
		Note: "Test project card 1",
	}
	repo := New(github.NewMockClient(
		github.MockMoveProjectCard(func(int, *github.MoveProjectCardParams) (*github.ProjectCard, error) {
			return expected, nil
		}),
	), "example/test")
	got, err := repo.MoveProjectCard(1, nil)
	assert.Nil(t, err)
	assert.Equal(t, got, expected)
}
