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
	return r.cli.GetProject(r.path, projectID)
}
