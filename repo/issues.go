package repo

import "github.com/itchyny/github-migrator/github"

// ListIssues lists the issues.
func (r *repo) ListIssues() github.Issues {
	return r.cli.ListIssues(r.path, &github.ListIssuesParams{
		Filter:    github.ListIssuesParamFilterAll,
		State:     github.ListIssuesParamStateAll,
		Direction: github.ListIssuesParamDirectionAsc,
	})
}

// GetIssue gets the issue.
func (r *repo) GetIssue(issueNumber int) (*github.Issue, error) {
	return r.cli.GetIssue(r.path, issueNumber)
}

// AddAssignees assigns users to the issue.
func (r *repo) AddAssignees(issueNumber int, assignees []string) error {
	return r.cli.AddAssignees(r.path, issueNumber, assignees)
}
