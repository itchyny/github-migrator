package migrator

import (
	"net/url"
	"regexp"
	"strings"

	"github.com/itchyny/github-migrator/github"
)

type commentFilter func(string) string

func newRepoURLFilter(sourceRepo, targetRepo *github.Repo) commentFilter {
	sourceURL, _ := url.Parse(sourceRepo.HTMLURL)
	targetURL, _ := url.Parse(targetRepo.HTMLURL)
	replaceImageLinks := sourceURL.Scheme != targetURL.Scheme || sourceURL.Host != targetURL.Host
	var imageMarkdownPattern, imageHTMLPattern *regexp.Regexp
	if replaceImageLinks {
		urlPatten := sourceURL.Scheme + `://` + regexp.QuoteMeta(sourceURL.Host) + `[^"<>()]+`
		imageMarkdownPattern = regexp.MustCompile(`(?i)!\[[^]]*\]\((` + urlPatten + `)\)`)
		imageHTMLPattern = regexp.MustCompile(`(?i)<img [^<>]*\bsrc="(` + urlPatten + `)"[^<>]*>`)
	}
	return commentFilter(func(src string) string {
		src = strings.ReplaceAll(src, sourceRepo.HTMLURL, targetRepo.HTMLURL)
		if replaceImageLinks {
			src = imageMarkdownPattern.ReplaceAllString(src, `<a href="$1">$0</a>`)
			src = imageHTMLPattern.ReplaceAllString(src, `<a href="$1">$0</a>`)
		}
		return src
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
