package repo

import "github.com/itchyny/github-migrator/github"

// Update the repository.
func (r *Repo) Update(params *github.UpdateRepoParams) (*github.Repo, error) {
	return r.cli.UpdateRepo(r.path, params)
}
