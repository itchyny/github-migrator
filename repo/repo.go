package repo

import "github.com/itchyny/github-migrator/github"

// Repo represents a GitHub repository.
type Repo struct {
	cli  github.Client
	path string
}

// New creates a new Repo.
func New(cli github.Client, path string) *Repo {
	return &Repo{cli: cli, path: path}
}

// NewPath creates a new Repo with the same client.
func (r *Repo) NewPath(path string) *Repo {
	return New(r.cli, path)
}
