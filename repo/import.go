package repo

import "github.com/itchyny/github-migrator/github"

// Import an object.
func (r *Repo) Import(x *github.Import) (*github.ImportResult, error) {
	return r.cli.Import(r.path, x)
}

// GetImport gets the importing status.
func (r *Repo) GetImport(id int) (*github.ImportResult, error) {
	return r.cli.GetImport(r.path, id)
}
