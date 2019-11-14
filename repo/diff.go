package repo

// GetDiff gets the diff.
func (r *repo) GetDiff(sha string) (string, error) {
	return r.cli.GetDiff(r.path, sha)
}
