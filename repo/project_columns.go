package repo

import "github.com/itchyny/github-migrator/github"

// ListProjectColumns lists the project columns.
func (r *repo) ListProjectColumns(projectID int) github.ProjectColumns {
	return r.cli.ListProjectColumns(projectID)
}

// GetProjectColumn gets the project column.
func (r *repo) GetProjectColumn(projectColumnID int) (*github.ProjectColumn, error) {
	return r.cli.GetProjectColumn(projectColumnID)
}

// CreateProjectColumn creates a project column.
func (r *repo) CreateProjectColumn(projectID int, name string) (*github.ProjectColumn, error) {
	return r.cli.CreateProjectColumn(projectID, name)
}

// UpdateProjectColumn updates the project column..
func (r *repo) UpdateProjectColumn(projectColumnID int, name string) (*github.ProjectColumn, error) {
	return r.cli.UpdateProjectColumn(projectColumnID, name)
}
