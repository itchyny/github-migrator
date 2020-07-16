package repo

// GetDiff gets the diff.
func (r *Repo) GetDiff(sha string) (string, error) {
	return r.cli.GetDiff(r.path, sha)
}

// GetCompare gets the compare.
func (r *Repo) GetCompare(base, head string) (string, error) {
	return r.cli.GetCompare(r.path, base, head)
}
