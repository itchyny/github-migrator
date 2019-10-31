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
		if err := m.migrateIssue(sourceRepo, issue, targetIssuesBuffer); err != nil {
			return err
		}
	}
	return nil
}

func (m *migrator) migrateIssue(sourceRepo *github.Repo, sourceIssue *github.Issue, targetIssuesBuffer *issuesBuffer) error {
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
	return m.target.Import(buildImport(sourceRepo, sourceIssue, comments))
}

func buildImport(repo *github.Repo, issue *github.Issue, comments []*github.Comment) *github.Import {
	importIssue := &github.ImportIssue{
		Title:     issue.Title,
		Body:      buildImportBody(repo, issue),
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
		Comments: buildImportComments(comments),
	}
}

func buildImportBody(repo *github.Repo, issue *github.Issue) string {
	return buildTable(
		buildImageTag(issue.User),
		fmt.Sprintf(
			"Original %s by @%s - imported from %s",
			issue.Type(),
			issue.User.Login,
			buildIssueLinkTag(repo, issue),
		),
	) + "\n\n" + issue.Body
}

func buildImportComments(comments []*github.Comment) []*github.ImportComment {
	xs := make([]*github.ImportComment, len(comments))
	for i, c := range comments {
		xs[i] = &github.ImportComment{
			Body: buildTable(
				buildImageTag(c.User),
				fmt.Sprintf("@%s commented", c.User.Login),
			) + "\n\n" + c.Body,
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
