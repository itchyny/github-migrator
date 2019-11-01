package migrator

import (
	"fmt"
	"io"
	"time"

	"github.com/itchyny/github-migrator/github"
)

func (m *migrator) migrateIssues() error {
	sourceRepo, err := m.getSourceRepo()
	if err != nil {
		return err
	}
	targetRepo, err := m.getTargetRepo()
	if err != nil {
		return err
	}
	sourceIssues := m.source.ListIssues()
	targetIssuesBuffer := newIssuesBuffer(m.target.ListIssues())
	commentFilters := newCommentFilters(
		newRepoURLFilter(sourceRepo, targetRepo),
		newUserMappingFilter(m.userMapping),
	)
	for {
		issue, err := sourceIssues.Next()
		if err != nil {
			if err != io.EOF {
				return err
			}
			break
		}
		result, err := m.migrateIssue(sourceRepo, targetRepo, commentFilters, issue, targetIssuesBuffer)
		if err != nil {
			return err
		}
		if result != nil {
			fmt.Printf("%#v\n", result)
		}
	}
	return nil
}

func (m *migrator) migrateIssue(
	sourceRepo, targetRepo *github.Repo, commentFilters commentFilters,
	sourceIssue *github.Issue, targetIssuesBuffer *issuesBuffer,
) (*github.ImportResult, error) {
	fmt.Printf("migrating: %s\n", sourceIssue.HTMLURL)
	targetIssue, err := targetIssuesBuffer.get(sourceIssue.Number)
	if err != nil {
		return nil, err
	}
	if targetIssue != nil {
		fmt.Printf("skipping: %s (already exists)\n", targetIssue.HTMLURL)
		return nil, nil
	}
	comments, err := github.CommentsToSlice(m.source.ListComments(sourceIssue.Number))
	if err != nil {
		return nil, err
	}
	var reviewComments []*github.ReviewComment
	if sourceIssue.PullRequest != nil {
		reviewComments, err = github.ReviewCommentsToSlice(m.source.ListReviewComments(sourceIssue.Number))
		if err != nil {
			return nil, err
		}
	}
	members, err := m.listTargetMembers()
	if err != nil {
		return nil, err
	}
	time.Sleep(time.Second)
	return m.target.Import(
		buildImport(
			sourceRepo, targetRepo, commentFilters,
			sourceIssue, comments, reviewComments, members,
		),
	)
}
