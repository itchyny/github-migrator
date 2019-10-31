package repo

import "github.com/itchyny/github-migrator/github"

// Get the repository.
func (r *repo) Get() (*github.Repo, error) {
	return r.cli.GetRepo(r.path)
}
