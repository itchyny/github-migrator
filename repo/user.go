package repo

import "github.com/itchyny/github-migrator/github"

// GetUser gets a user.
func (r *Repo) GetUser(name string) (*github.User, error) {
	return r.cli.GetUser(name)
}
