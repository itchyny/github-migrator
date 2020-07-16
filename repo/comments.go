package repo

import "github.com/itchyny/github-migrator/github"

// ListComments lists the comments.
func (r *Repo) ListComments(issueNumber int) github.Comments {
	return r.cli.ListComments(r.path, issueNumber)
}
