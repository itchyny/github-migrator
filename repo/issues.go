package repo

import "github.com/itchyny/github-migrator/github"

// ListIssues lists the issues.
func (r *repo) ListIssues() ([]*github.Issue, error) {
	return r.cli.ListIssues(r.path, &github.ListIssuesParams{
		Filter:    github.ListIssuesParamFilterAll,
		State:     github.ListIssuesParamStateAll,
		Direction: github.ListIssuesParamDirectionAsc,
	})
}
