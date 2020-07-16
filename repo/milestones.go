package repo

import "github.com/itchyny/github-migrator/github"

// ListMilestones lists the milestones.
func (r *Repo) ListMilestones(params *github.ListMilestonesParams) github.Milestones {
	return r.cli.ListMilestones(r.path, params)
}

// GetMilestone gets the milestone.
func (r *Repo) GetMilestone(milestoneNumber int) (*github.Milestone, error) {
	return r.cli.GetMilestone(r.path, milestoneNumber)
}

// CreateMilestone creates a milestone.
func (r *Repo) CreateMilestone(params *github.CreateMilestoneParams) (*github.Milestone, error) {
	return r.cli.CreateMilestone(r.path, params)
}

// UpdateMilestone updates the milestone.
func (r *Repo) UpdateMilestone(milestoneNumber int, params *github.UpdateMilestoneParams) (*github.Milestone, error) {
	return r.cli.UpdateMilestone(r.path, milestoneNumber, params)
}

// DeleteMilestone deletes the milestone.
func (r *Repo) DeleteMilestone(milestoneNumber int) error {
	return r.cli.DeleteMilestone(r.path, milestoneNumber)
}
