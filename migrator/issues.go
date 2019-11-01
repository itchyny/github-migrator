package migrator

import (
	"fmt"
	"io"
	"strings"
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
	fs := newCommentFilters(
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
		if err := m.migrateIssue(sourceRepo, fs, issue, targetIssuesBuffer); err != nil {
			return err
		}
		time.Sleep(time.Second)
	}
	return nil
}

func (m *migrator) migrateIssue(
	sourceRepo *github.Repo, fs commentFilters,
	sourceIssue *github.Issue, targetIssuesBuffer *issuesBuffer,
) error {
	fmt.Printf("migrating: %s\n", sourceIssue.HTMLURL)
	targetIssue, err := targetIssuesBuffer.get(sourceIssue.Number)
	if err != nil {
		return err
	}
	if targetIssue != nil {
		fmt.Printf("skipping: %s (already exists)\n", targetIssue.HTMLURL)
		return nil
	}
	comments, err := github.CommentsToSlice(m.source.ListComments(sourceIssue.Number))
	if err != nil {
		return err
	}
	var reviewComments []*github.ReviewComment
	if sourceIssue.PullRequest != nil {
		reviewComments, err = github.ReviewCommentsToSlice(m.source.ListReviewComments(sourceIssue.Number))
		if err != nil {
			return err
		}
	}
	im, err := m.buildImport(sourceRepo, fs, sourceIssue, comments, reviewComments)
	if err != nil {
		return err
	}
	return m.target.Import(im)
}

func (m *migrator) buildImport(
	sourceRepo *github.Repo, fs commentFilters, issue *github.Issue,
	comments []*github.Comment, reviewComments []*github.ReviewComment,
) (*github.Import, error) {
	importIssue := &github.ImportIssue{
		Title:     issue.Title,
		Body:      buildImportBody(sourceRepo, fs, issue),
		CreatedAt: issue.CreatedAt,
		UpdatedAt: issue.UpdatedAt,
		Closed:    issue.State != "open",
		ClosedAt:  issue.ClosedAt,
		Labels:    buildImportLabels(issue),
	}
	if issue.Assignee != nil {
		target := fs.apply(issue.Assignee.Login)
		ok, err := m.isTargetMember(target)
		if err != nil {
			return nil, err
		}
		if ok {
			importIssue.Assignee = target
		}
	}
	return &github.Import{
		Issue:    importIssue,
		Comments: buildImportComments(fs, comments, reviewComments),
	}, nil
}

func buildImportBody(sourceRepo *github.Repo, fs commentFilters, issue *github.Issue) string {
	return buildTable(
		buildImageTag(fs, issue.User),
		fmt.Sprintf(
			"@%s created the original %s at %s<br>imported from %s",
			fs.apply(issue.User.Login), issue.Type(), formatTimestamp(issue.CreatedAt),
			buildIssueLinkTag(sourceRepo, issue),
		),
	) + "\n\n" + fs.apply(issue.Body)
}

func buildImportComments(
	fs commentFilters,
	comments []*github.Comment,
	reviewComments []*github.ReviewComment,
) []*github.ImportComment {
	xs := make([]*github.ImportComment, len(comments))
	for i, c := range comments {
		xs[i] = &github.ImportComment{
			Body: buildTable(
				buildImageTag(fs, c.User),
				fmt.Sprintf("@%s commented at %s", fs.apply(c.User.Login), formatTimestamp(c.CreatedAt)),
			) + "\n\n" + fs.apply(c.Body),
			CreatedAt: c.CreatedAt,
		}
	}
	reviewCommentsIDToIndex := make(map[int]int)
	for _, c := range reviewComments {
		if i, ok := reviewCommentsIDToIndex[c.InReplyToID]; ok {
			reviewCommentsIDToIndex[c.ID] = i
			xs[i].Body += "\n\n" + buildTable(
				buildImageTag(fs, c.User),
				fmt.Sprintf("@%s commented at %s", fs.apply(c.User.Login), formatTimestamp(c.CreatedAt)),
			) + "\n\n" + fs.apply(c.Body)
			continue
		}
		reviewCommentsIDToIndex[c.ID] = len(xs)
		xs = append(xs, &github.ImportComment{
			Body: strings.Join([]string{"```diff", fmt.Sprintf("# %s:%d", c.Path, c.Line), c.DiffHunk, "```\n\n"}, "\n") +
				buildTable(
					buildImageTag(fs, c.User),
					fmt.Sprintf("@%s commented at %s", fs.apply(c.User.Login), formatTimestamp(c.CreatedAt)),
				) + "\n\n" + fs.apply(c.Body),
			CreatedAt: c.CreatedAt,
		})
	}
	return xs
}

func formatTimestamp(src string) string {
	t, err := time.Parse(time.RFC3339, src)
	if err != nil {
		return ""
	}
	return t.Local().String()
}

func buildImageTag(fs commentFilters, user *github.User) string {
	return fmt.Sprintf(`<img src="https://github.com/%s.png" width="35">`, fs.apply(user.Login))
}

func buildTable(xs ...string) string {
	s := new(strings.Builder)
	s.WriteString("<table>\n")
	s.WriteString("  <tr>\n")
	for _, x := range xs {
		s.WriteString("    <td>\n")
		s.WriteString("      " + x + "\n")
		s.WriteString("    </td>\n")
	}
	s.WriteString("  </tr>\n")
	s.WriteString("</table>\n")
	return s.String()
}

func buildIssueLinkTag(repo *github.Repo, issue *github.Issue) string {
	return fmt.Sprintf(`<a href="%s">%s#%d</a>`, issue.HTMLURL, repo.FullName, issue.Number)
}

func buildImportLabels(issue *github.Issue) []string {
	xs := []string{}
	for _, l := range issue.Labels {
		xs = append(xs, l.Name)
	}
	return xs
}
