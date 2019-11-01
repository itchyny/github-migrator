package repo

import (
	"strings"

	"github.com/itchyny/github-migrator/github"
)

// List the members.
func (r *repo) ListMembers() github.Members {
	return r.cli.ListMembers(strings.Split(r.path, "/")[0])
}
