package migrator

import (
	"strings"

	"github.com/itchyny/github-migrator/github"
)

type commentFilter func(string) string

func applyCommentFilters(fs []commentFilter, src string) string {
	for _, f := range fs {
		src = f(src)
	}
	return src
}

func newRepoUrlFilter(sourceRepo, targetRepo *github.Repo) commentFilter {
	return commentFilter(func(src string) string {
		return strings.ReplaceAll(src, sourceRepo.HTMLURL, targetRepo.HTMLURL)
	})
}
