package repo

import "github.com/itchyny/github-migrator/github"

// ListLabels lists the labels.
func (r *Repo) ListLabels() github.Labels {
	return r.cli.ListLabels(r.path)
}

// CreateLabel creates a new label.
func (r *Repo) CreateLabel(params *github.CreateLabelParams) (*github.Label, error) {
	return r.cli.CreateLabel(r.path, params)
}

// UpdateLabel creates a new label.
func (r *Repo) UpdateLabel(name string, params *github.UpdateLabelParams) (*github.Label, error) {
	return r.cli.UpdateLabel(r.path, name, params)
}
