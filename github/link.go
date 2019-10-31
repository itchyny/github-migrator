package github

import (
	"net/http"
	"net/url"
	"strings"
)

type link struct {
	next, last string
}

func parseLink(src []string) *link {
	l := &link{}
	for _, s := range src {
		for _, x := range strings.Split(s, ", ") {
			ys := strings.Split(x, "; ")
			if len(ys) != 2 {
				continue
			}
			switch ys[1] {
			case `rel="next"`:
				l.next = ys[0][1:][:len(ys[0])-2]
			case `rel="last"`:
				l.last = ys[0][1:][:len(ys[0])-2]
			}
		}
	}
	return l
}

func getNext(header http.Header) string {
	next := parseLink(header["Link"]).next
	if next == "" {
		return ""
	}

	if _, err := url.Parse(next); err != nil {
		return ""
	}

	return next
}
