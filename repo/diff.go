package repo

// GetDiff gets the diff.
func (r *repo) GetDiff(sha string) (string, error) {
	return r.cli.GetDiff(r.path, sha)
}

// GetCompare gets the compare.
func (r *repo) GetCompare(base, head string) (string, error) {
	return r.cli.GetCompare(r.path, base, head)
}
