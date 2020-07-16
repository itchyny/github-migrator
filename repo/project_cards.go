package repo

import "github.com/itchyny/github-migrator/github"

// ListProjectCards lists the project cards.
func (r *Repo) ListProjectCards(columnID int) github.ProjectCards {
	return r.cli.ListProjectCards(columnID)
}

// GetProjectCard gets the project card.
func (r *Repo) GetProjectCard(projectCardID int) (*github.ProjectCard, error) {
	return r.cli.GetProjectCard(projectCardID)
}

// CreateProjectCard creates a project card.
func (r *Repo) CreateProjectCard(columnID int, params *github.CreateProjectCardParams) (*github.ProjectCard, error) {
	return r.cli.CreateProjectCard(columnID, params)
}

// UpdateProjectCard updates the project card..
func (r *Repo) UpdateProjectCard(projectCardID int, params *github.UpdateProjectCardParams) (*github.ProjectCard, error) {
	return r.cli.UpdateProjectCard(projectCardID, params)
}

// MoveProjectCard moves the project card..
func (r *Repo) MoveProjectCard(projectCardID int, params *github.MoveProjectCardParams) (*github.ProjectCard, error) {
	return r.cli.MoveProjectCard(projectCardID, params)
}
