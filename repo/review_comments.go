package repo

import "github.com/itchyny/github-migrator/github"

// ListReviewComments lists the review comments.
func (r *repo) ListReviewComments(pullNumber int) github.ReviewComments {
	return r.cli.ListReviewComments(r.path, pullNumber)
}
