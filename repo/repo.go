package repo

import (
	"fmt"

	"github.com/itchyny/github-migrator/github"
)

type Repo interface {
	Name() string
}

func New(cli github.Client, path string) *repo {
	return &repo{cli: cli, path: path}
}

type repo struct {
	cli  github.Client
	path string
}

func (r *repo) Name() string {
	return fmt.Sprintf("%s:%s", r.cli.Hostname(), r.path)
}
