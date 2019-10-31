package repo

import (
	"fmt"

	"github.com/itchyny/github-migrator/github"
)

// Repo represents a GitHub repository.
type Repo interface {
	Name() string
	ListIssues() ([]*github.Issue, error)
}

// New creates a new Repo.
func New(cli github.Client, path string) Repo {
	return &repo{cli: cli, path: path}
}

type repo struct {
	cli  github.Client
	path string
}

// Name ...
func (r *repo) Name() string {
	return fmt.Sprintf("%s:%s", r.cli.Hostname(), r.path)
}
