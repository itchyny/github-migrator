package migrator

import (
	"regexp"
	"strings"
)

const (
	truncateLength      = 10000
	totalTruncateLength = 60000
)

// Since import fails on too large diff, truncate it.
// You may wonder building the diff from the api (without vnd.github.v3.diff header),
// but it's impossible to build the complete diff even if the changes are small.
func truncateDiff(diff string) string {
	s := new(strings.Builder)
	var i, j int
	for {
		i = strings.Index(diff, "\nindex ")
		if i < 0 {
			if len(diff) > truncateLength {
				s.WriteString("Too large diff\n")
				break
			}
			s.WriteString(diff)
			break
		}
		i++ // newline

		j = strings.Index(diff[i:], "\n")
		if j < 0 {
			if len(diff) > truncateLength {
				s.WriteString("Too large diff\n")
				break
			}
			s.WriteString(diff)
			break
		}
		j++ // newline
		s.WriteString(diff[:i+j])
		diff = diff[i+j:]

		i = strings.Index(diff, "\ndiff ")
		if i < 0 {
			if len(diff) > truncateLength {
				s.WriteString("Too large diff\n")
				break
			}
			s.WriteString(diff)
			break
		}
		i++ // newline
		if i > truncateLength {
			s.WriteString("Too large diff\n")
			diff = diff[i:]
			continue
		}
		s.WriteString(diff[:i])
		diff = diff[i:]
	}
	str := s.String()
	if len(str) > totalTruncateLength {
		str = str[:totalTruncateLength] + "\n\nToo large diff\n"
	}
	return str
}

var backquoteRe = regexp.MustCompile("((?:^|\n) *)```")

func escapeBackQuotes(src string) string {
	return backquoteRe.ReplaceAllString(src, "$1\u00a0```")
}
