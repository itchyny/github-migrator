package repo

import "github.com/itchyny/github-migrator/github"

// Import an object.
func (r *repo) Import(x *github.Import) error {
	return r.cli.Import(r.path, x)
}
