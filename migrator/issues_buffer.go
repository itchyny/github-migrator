package migrator

import (
	"io"

	"github.com/itchyny/github-migrator/github"
)

type issuesBuffer struct {
	src    github.Issues
	issues []*github.Issue
}

func newIssuesBuffer(is github.Issues) *issuesBuffer {
	return &issuesBuffer{src: is}
}

func (ib *issuesBuffer) get(num int) (*github.Issue, error) {
	for _, issue := range ib.issues {
		if issue.Number == num {
			return issue, nil
		}
	}
	for {
		issue, err := ib.src.Next()
		if err != nil {
			if err != io.EOF {
				return nil, err
			}
			break
		}
		ib.issues = append(ib.issues, issue)
		if issue.Number == num {
			return issue, nil
		} else if issue.Number > num {
			return nil, nil
		}
	}
	return nil, nil
}
