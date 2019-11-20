package repo

import "github.com/itchyny/github-migrator/github"

// Get a user.
func (r *repo) GetUser(name string) (*github.User, error) {
	return r.cli.GetUser(name)
}
