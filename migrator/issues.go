package migrator

import (
	"fmt"
	"io"
	"strings"

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
	for {
		issue, err := sourceIssues.Next()
		if err != nil {
			if err != io.EOF {
				return err
			}
			break
		}
		if err := m.migrateIssue(sourceRepo, targetRepo, issue, targetIssuesBuffer); err != nil {
			return err
		}
	}
	return nil
}

func (m *migrator) migrateIssue(sourceRepo, targetRepo *github.Repo, sourceIssue *github.Issue, targetIssuesBuffer *issuesBuffer) error {
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
	return m.target.Import(buildImport(sourceRepo, targetRepo, sourceIssue, comments))
}

func buildImport(sourceRepo, targetRepo *github.Repo, issue *github.Issue, comments []*github.Comment) *github.Import {
	importIssue := &github.ImportIssue{
		Title:     issue.Title,
		Body:      buildImportBody(sourceRepo, targetRepo, issue),
		CreatedAt: issue.CreatedAt,
		UpdatedAt: issue.UpdatedAt,
		Closed:    issue.State != "open",
		ClosedAt:  issue.ClosedAt,
		Labels:    buildImportLabels(issue),
	}
	if issue.Assignee != nil {
		importIssue.Assignee = issue.Assignee.Login
	}
	return &github.Import{
		Issue:    importIssue,
		Comments: buildImportComments(sourceRepo, targetRepo, comments),
	}
}

func buildImportBody(sourceRepo, targetRepo *github.Repo, issue *github.Issue) string {
	return buildTable(
		buildImageTag(issue.User),
		fmt.Sprintf(
			"Original %s by @%s - imported from %s",
			issue.Type(),
			issue.User.Login,
			buildIssueLinkTag(sourceRepo, issue),
		),
	) + "\n\n" + strings.ReplaceAll(issue.Body, sourceRepo.HTMLURL, targetRepo.HTMLURL)
}

func buildImportComments(sourceRepo, targetRepo *github.Repo, comments []*github.Comment) []*github.ImportComment {
	xs := make([]*github.ImportComment, len(comments))
	for i, c := range comments {
		xs[i] = &github.ImportComment{
			Body: buildTable(
				buildImageTag(c.User),
				fmt.Sprintf("@%s commented", c.User.Login),
			) + "\n\n" + strings.ReplaceAll(c.Body, sourceRepo.HTMLURL, targetRepo.HTMLURL),
			CreatedAt: c.CreatedAt,
		}
	}
	return xs
}

func buildImageTag(user *github.User) string {
	return fmt.Sprintf(`<img src="https://github.com/%s.png" width="35">`, user.Login)
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
