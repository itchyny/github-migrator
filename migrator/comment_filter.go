package migrator

import (
	"regexp"
	"strings"

	"github.com/itchyny/github-migrator/github"
)

type commentFilter func(string) string

func newRepoURLFilter(sourceRepo, targetRepo *github.Repo) commentFilter {
	return commentFilter(func(src string) string {
		return strings.ReplaceAll(src, sourceRepo.HTMLURL, targetRepo.HTMLURL)
	})
}

func newUserMappingFilter(userMapping map[string]string) commentFilter {
	if len(userMapping) == 0 {
		return commentFilter(func(src string) string {
			return src
		})
	}
	return commentFilter(func(src string) string {
		for k, v := range userMapping {
			r, err := regexp.Compile(`\b` + k + `\b`)
			if err != nil {
				continue
			}
			src = r.ReplaceAllString(src, v)
		}
		return src
	})
}

type commentFilters []commentFilter

func newCommentFilters(fs ...commentFilter) commentFilters {
	return commentFilters(fs)
}

func (fs commentFilters) apply(src string) string {
	for _, f := range fs {
		src = f(src)
	}
	return src
}
