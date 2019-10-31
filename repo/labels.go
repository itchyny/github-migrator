package repo

import "github.com/itchyny/github-migrator/github"

// ListLabels lists the labels.
func (r *repo) ListLabels() github.Labels {
	return r.cli.ListLabels(r.path)
}

// CreateLabel creates a new label.
func (r *repo) CreateLabel(param *github.CreateLabelParams) (*github.Label, error) {
	return r.cli.CreateLabel(r.path, param)
}
