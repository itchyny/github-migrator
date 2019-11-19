package repo

import "github.com/itchyny/github-migrator/github"

// ListProjects lists the projects.
func (r *repo) ListProjects() github.Projects {
	return r.cli.ListProjects(r.path, &github.ListProjectsParams{
		State: github.ListProjectsParamStateAll,
	})
}

// GetProject gets the project.
func (r *repo) GetProject(projectID int) (*github.Project, error) {
	return r.cli.GetProject(projectID)
}

// CreateProject creates a project.
func (r *repo) CreateProject(params *github.CreateProjectParams) (*github.Project, error) {
	return r.cli.CreateProject(r.path, params)
}

// UpdateProject updates the project.
func (r *repo) UpdateProject(projectID int, params *github.UpdateProjectParams) (*github.Project, error) {
	return r.cli.UpdateProject(projectID, params)
}
