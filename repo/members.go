package repo

import (
	"strings"

	"github.com/itchyny/github-migrator/github"
)

// ListMembers lists the members.
func (r *Repo) ListMembers() github.Members {
	return r.cli.ListMembers(strings.Split(r.path, "/")[0])
}
