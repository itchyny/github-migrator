package repo

import "github.com/itchyny/github-migrator/github"

// ListEvents lists the events.
func (r *repo) ListEvents(issueNumber int) github.Events {
	return r.cli.ListEvents(r.path, issueNumber)
}
