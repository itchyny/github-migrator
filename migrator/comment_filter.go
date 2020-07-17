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

func newUserMappingFilter(userMapping map[string]string, targetRepo *github.Repo) commentFilter {
	if len(userMapping) == 0 {
		return commentFilter(func(src string) string {
			return src
		})
	}
	froms := make([]string, 0, len(userMapping))
	tos := make([]string, 0, len(userMapping))
	userMappingRev := make(map[string]string, len(userMapping))
	for k, v := range userMapping {
		froms = append(froms, k)
		tos = append(tos, v)
		userMappingRev[v] = k
	}
	re1 := regexp.MustCompile(buildPattern(froms))
	re2 := regexp.MustCompile(buildPattern(tos))
	re3 := regexp.MustCompile(`https?://[-.a-zA-Z0-9/_%]*` + buildPattern(tos))
	targetURL, _ := url.Parse(targetRepo.HTMLURL)
	return commentFilter(func(src string) string {
		src = re1.ReplaceAllStringFunc(src, func(from string) string {
			return userMapping[from]
		})
		src = re3.ReplaceAllStringFunc(src, func(url string) string {
			if strings.Contains(url, "://"+targetURL.Host+"/") {
				return url
			}
			return re2.ReplaceAllStringFunc(url, func(to string) string {
				return userMappingRev[to]
			})
		})
		return src
	})
}

func buildPattern(xs []string) string {
	var pattern strings.Builder
	pattern.WriteString(`\b(`)
	for i, x := range xs {
		if i > 0 {
			pattern.WriteByte('|')
		}
		pattern.WriteString(regexp.QuoteMeta(x))
		i++
	}
	pattern.WriteString(`)\b`)
	return pattern.String()
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
