package repo

import "github.com/itchyny/github-migrator/github"

// ListMilestones lists the milestones.
func (r *repo) ListMilestones(params *github.ListMilestonesParams) github.Milestones {
	return r.cli.ListMilestones(r.path, params)
}

// GetMilestone gets the milestone.
func (r *repo) GetMilestone(milestoneNumber int) (*github.Milestone, error) {
	return r.cli.GetMilestone(r.path, milestoneNumber)
}

// CreateMilestone creates a milestone.
func (r *repo) CreateMilestone(params *github.CreateMilestoneParams) (*github.Milestone, error) {
	return r.cli.CreateMilestone(r.path, params)
}

// UpdateMilestone updates the milestone.
func (r *repo) UpdateMilestone(milestoneNumber int, params *github.UpdateMilestoneParams) (*github.Milestone, error) {
	return r.cli.UpdateMilestone(r.path, milestoneNumber, params)
}

// DeleteMilestone deletes the milestone.
func (r *repo) DeleteMilestone(milestoneNumber int) error {
	return r.cli.DeleteMilestone(r.path, milestoneNumber)
}
