package repo

import "github.com/itchyny/github-migrator/github"

// ListLabels lists the labels.
func (r *repo) ListLabels() github.Labels {
	return r.cli.ListLabels(r.path)
}
