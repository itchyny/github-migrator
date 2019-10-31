package migrator

import (
	"fmt"
	"io"

	"github.com/itchyny/github-migrator/github"
)

func (m *migrator) migrateIssues() error {
	sourceIssues := m.source.ListIssues()
	targetIssuesBuffer := newIssuesBuffer(m.target.ListIssues())
	for {
		issue, err := sourceIssues.Next()
		if err != nil {
			if err != io.EOF {
				return err
			}
			break
		}
		if err := m.migrateIssue(issue, targetIssuesBuffer); err != nil {
			return err
		}
	}
	return nil
}

func (m *migrator) migrateIssue(sourceIssue *github.Issue, targetIssuesBuffer *issuesBuffer) error {
	fmt.Printf("importing: %s\n", sourceIssue.HTMLURL)
	targetIssue, err := targetIssuesBuffer.get(sourceIssue.Number)
	if err != nil {
		return err
	}
	if targetIssue != nil {
		fmt.Printf("skipping: %s (already exists)\n", targetIssue.HTMLURL)
		return nil
	}
	return nil
}
