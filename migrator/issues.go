package migrator

import (
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/itchyny/github-migrator/github"
)

var (
	beforeImportIssueDuration      = 1 * time.Second
	waitImportIssueInitialDuration = 3 * time.Second
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
	var lastIssueNumber int
	for {
		issue, err := sourceIssues.Next()
		if err != nil {
			if err != io.EOF {
				return err
			}
			break
		}
		for ; issue.Number > lastIssueNumber; lastIssueNumber++ {
			issue := issue
			var deleted bool
			if deleted = issue.Number > lastIssueNumber+1; deleted {
				issue = &github.Issue{
					Number:    lastIssueNumber + 1,
					HTMLURL:   fmt.Sprintf("%s/issues/%d", sourceRepo.HTMLURL, lastIssueNumber+1),
					CreatedAt: issue.CreatedAt,
					UpdatedAt: issue.CreatedAt,
					ClosedAt:  issue.CreatedAt,
				}
			}
			result, err := m.migrateIssue(sourceRepo, targetRepo, commentFilters, issue, targetIssuesBuffer, deleted)
			if err != nil {
				return err
			}
			if result != nil {
				if err := m.waitImportIssue(result.ID, issue); err != nil {
					return fmt.Errorf("importing %s failed: %w", issue.HTMLURL, err)
				}
			}
		}
	}
	return nil
}

func (m *migrator) migrateIssue(
	sourceRepo, targetRepo *github.Repo, commentFilters commentFilters,
	sourceIssue *github.Issue, targetIssuesBuffer *issuesBuffer, deleted bool,
) (*github.ImportResult, error) {
	fmt.Printf("[=>] migrating an issue: %s\n", sourceIssue.HTMLURL)
	targetIssue, err := targetIssuesBuffer.get(sourceIssue.Number)
	if err != nil {
		return nil, err
	}
	if targetIssue != nil {
		fmt.Printf("[--] skipping: %s (already exists)\n", targetIssue.HTMLURL)
		return nil, nil
	}
	time.Sleep(beforeImportIssueDuration)
	if deleted {
		fmt.Printf("[>>] creating a new issue: (original: %s is deleted)\n", sourceIssue.HTMLURL)
		return m.target.Import(&github.Import{
			Issue: &github.ImportIssue{
				Title: "[Deleted issue]",
				Body: fmt.Sprintf(`<table>
<tr>
  <td>This issue was imported from %s, which has already been deleted.</td>
</tr>
</table>
`, buildIssueLinkTag(sourceRepo, sourceIssue)),
				CreatedAt: sourceIssue.CreatedAt,
				UpdatedAt: sourceIssue.UpdatedAt,
				Closed:    true,
				ClosedAt:  sourceIssue.ClosedAt,
			},
			Comments: []*github.ImportComment{},
		})
	}
	comments, err := github.CommentsToSlice(m.source.ListComments(sourceIssue.Number))
	if err != nil {
		return nil, err
	}
	events, err := github.EventsToSlice(m.source.ListEvents(sourceIssue.Number))
	if err != nil {
		return nil, err
	}
	var sourcePullReq *github.PullReq
	var commits []*github.Commit
	var commitDiff string
	var reviews []*github.Review
	var reviewComments []*github.ReviewComment
	if sourceIssue.PullRequest != nil {
		sourcePullReq, err = m.source.GetPullReq(sourceIssue.Number)
		if err != nil {
			return nil, err
		}
		commits, err = github.CommitsToSlice(m.source.ListPullReqCommits(sourceIssue.Number))
		if err != nil {
			return nil, err
		}
		commitDiff, err = m.source.NewPath(sourcePullReq.Base.Repo.FullName).
			GetCompare(sourcePullReq.Base.SHA, sourcePullReq.Head.SHA)
		if err != nil {
			return nil, err
		}
		reviews, err = github.ReviewsToSlice(m.source.ListReviews(sourceIssue.Number))
		if err != nil {
			return nil, err
		}
		reviewComments, err = github.ReviewCommentsToSlice(m.source.ListReviewComments(sourceIssue.Number))
		if err != nil {
			return nil, err
		}
	}
	projects, err := m.listTargetProjects()
	if err != nil {
		return nil, err
	}
	imp, err := buildImport(
		m, m.source, m.target, sourceRepo, targetRepo, commentFilters,
		sourceIssue, sourcePullReq, comments, events,
		commits, commitDiff, reviews, reviewComments, projects,
	)
	if err != nil {
		return nil, err
	}
	fmt.Printf("[>>] creating a new issue: (original: %s)\n", sourceIssue.HTMLURL)
	return m.target.Import(imp)
}

func (m *migrator) waitImportIssue(id int, issue *github.Issue) error {
	var retry int
	duration := waitImportIssueInitialDuration
	for {
		time.Sleep(duration)
		if retry > 1 {
			duration *= 2
		}
		res, err := m.target.GetImport(id)
		if err != nil {
			return err
		}
		switch res.Status {
		case "imported":
			fmt.Printf("[<>] checking status: %s (importing %s)\n", res.Status, issue.HTMLURL)
			return nil
		case "failed":
			fmt.Printf("[!!] checking status: %s (importing %s)\n", res.Status, issue.HTMLURL)
			return errors.New("failed status")
		default:
			fmt.Printf("[??] checking status: %s (importing %s)\n", res.Status, issue.HTMLURL)
		}
		retry++
		if retry >= 5 {
			return errors.New("reached maximum retry count")
		}
	}
}
