package repo

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/itchyny/github-migrator/github"
)

func TestRepoListMilestones(t *testing.T) {
	expected := []*github.Milestone{
		&github.Milestone{
			ID:    10,
			Title: "Test milestone 1",
		},
		&github.Milestone{
			ID:    10,
			Title: "Test milestone 1",
		},
	}
	repo := New(github.NewMockClient(
		github.MockListMilestones(func(string, *github.ListMilestonesParams) github.Milestones {
			return github.MilestonesFromSlice(expected)
		}),
	), "example/test")
	got, err := github.MilestonesToSlice(repo.ListMilestones(nil))
	assert.Nil(t, err)
	assert.Equal(t, got, expected)
}

func TestRepoGetMilestone(t *testing.T) {
	expected := &github.Milestone{
		ID:    1,
		Title: "Test milestone 1",
	}
	repo := New(github.NewMockClient(
		github.MockGetMilestone(func(string, int) (*github.Milestone, error) {
			return expected, nil
		}),
	), "example/test")
	got, err := repo.GetMilestone(1)
	assert.Nil(t, err)
	assert.Equal(t, got, expected)
}

func TestRepoCreateMilestone(t *testing.T) {
	expected := &github.Milestone{
		ID:    1,
		Title: "Test milestone 1",
	}
	repo := New(github.NewMockClient(
		github.MockCreateMilestone(func(string, *github.CreateMilestoneParams) (*github.Milestone, error) {
			return expected, nil
		}),
	), "example/test")
	got, err := repo.CreateMilestone(&github.CreateMilestoneParams{})
	assert.Nil(t, err)
	assert.Equal(t, got, expected)
}

func TestRepoUpdateMilestone(t *testing.T) {
	expected := &github.Milestone{
		ID:    1,
		Title: "Test milestone 1",
	}
	repo := New(github.NewMockClient(
		github.MockUpdateMilestone(func(string, int, *github.UpdateMilestoneParams) (*github.Milestone, error) {
			return expected, nil
		}),
	), "example/test")
	got, err := repo.UpdateMilestone(1, &github.UpdateMilestoneParams{})
	assert.Nil(t, err)
	assert.Equal(t, got, expected)
}

func TestRepoDeleteMilestone(t *testing.T) {
	repo := New(github.NewMockClient(
		github.MockDeleteMilestone(func(string, int) error {
			return nil
		}),
	), "example/test")
	err := repo.DeleteMilestone(1)
	assert.Nil(t, err)
}
