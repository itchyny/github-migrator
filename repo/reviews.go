package repo

import "github.com/itchyny/github-migrator/github"

// ListReviews lists the reviews.
func (r *repo) ListReviews(pullNumber int) github.Reviews {
	return r.cli.ListReviews(r.path, pullNumber)
}
