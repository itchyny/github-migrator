package repo

import "github.com/itchyny/github-migrator/github"

// ListProjects lists the projects.
func (r *Repo) ListProjects() github.Projects {
	return r.cli.ListProjects(r.path, &github.ListProjectsParams{
		State: github.ListProjectsParamStateAll,
	})
}

// GetProject gets the project.
func (r *Repo) GetProject(projectID int) (*github.Project, error) {
	return r.cli.GetProject(projectID)
}

// CreateProject creates a project.
func (r *Repo) CreateProject(params *github.CreateProjectParams) (*github.Project, error) {
	return r.cli.CreateProject(r.path, params)
}

// UpdateProject updates the project.
func (r *Repo) UpdateProject(projectID int, params *github.UpdateProjectParams) (*github.Project, error) {
	return r.cli.UpdateProject(projectID, params)
}

// DeleteProject deletes the project.
func (r *Repo) DeleteProject(projectID int) error {
	return r.cli.DeleteProject(projectID)
}
