package migrator

import (
	"fmt"
	"strings"

	"github.com/itchyny/github-migrator/github"
)

type builder struct {
	source, target *github.Repo
	commentFilters commentFilters
	issue          *github.Issue
	comments       []*github.Comment
	reviewComments []*github.ReviewComment
	members        []*github.Member
}

func buildImport(
	sourceRepo, targetRepo *github.Repo, commentFilters commentFilters,
	issue *github.Issue,
	comments []*github.Comment, reviewComments []*github.ReviewComment,
	members []*github.Member,
) *github.Import {
	return (&builder{
		source:         sourceRepo,
		target:         targetRepo,
		commentFilters: commentFilters,
		issue:          issue,
		comments:       comments,
		reviewComments: reviewComments,
		members:        members,
	}).build()
}

func (b *builder) build() *github.Import {
	importIssue := &github.ImportIssue{
		Title:     b.issue.Title,
		Body:      b.buildImportBody(),
		CreatedAt: b.issue.CreatedAt,
		UpdatedAt: b.issue.UpdatedAt,
		Closed:    b.issue.State != "open",
		ClosedAt:  b.issue.ClosedAt,
		Labels:    b.buildImportLabels(b.issue),
	}
	if b.issue.Assignee != nil {
		target := b.commentFilters.apply(b.issue.Assignee.Login)
		if b.isTargetMember(target) {
			importIssue.Assignee = target
		}
	}
	return &github.Import{
		Issue:    importIssue,
		Comments: b.buildImportComments(),
	}
}

func (b *builder) buildImportBody() string {
	return b.buildTable(
		b.buildImageTag(b.issue.User),
		fmt.Sprintf(
			"@%s created the original %s<br>imported from %s",
			b.commentFilters.apply(b.issue.User.Login),
			b.issue.Type(),
			b.buildIssueLinkTag(b.source, b.issue),
		),
	) + "\n\n" + b.commentFilters.apply(b.issue.Body)
}

func (b *builder) buildImportComments() []*github.ImportComment {
	xs := append(
		b.buildImportIssueComments(),
		b.buildImportReviewComments()...,
	)
	for _, x := range xs {
		x.Body = b.commentFilters.apply(x.Body)
	}
	return xs
}

func (b *builder) buildImportIssueComments() []*github.ImportComment {
	xs := make([]*github.ImportComment, len(b.comments))
	for i, c := range b.comments {
		xs[i] = &github.ImportComment{
			Body:      b.buildCommentedTable(c.User, c.Body),
			CreatedAt: c.CreatedAt,
		}
	}
	return xs
}

func (b *builder) buildImportReviewComments() []*github.ImportComment {
	var xs []*github.ImportComment
	indexByID := make(map[int]int)
	for _, c := range b.reviewComments {
		if i, ok := indexByID[c.InReplyToID]; ok {
			indexByID[c.ID] = i
			xs[i].Body += "\n\n" + b.buildCommentedTable(c.User, c.Body)
			continue
		}
		indexByID[c.ID] = len(xs)
		xs = append(xs, &github.ImportComment{
			Body: strings.Join([]string{"```diff", "# " + c.Path, c.DiffHunk, "```\n\n"}, "\n") +
				b.buildCommentedTable(c.User, c.Body),
			CreatedAt: c.CreatedAt,
		})
	}
	return xs
}

func (b *builder) buildCommentedTable(user *github.User, body string) string {
	return b.buildTable(
		b.buildImageTag(user), fmt.Sprintf("@%s commented", user.Login),
	) + "\n\n" + body
}

func (b *builder) buildImageTag(user *github.User) string {
	target := b.commentFilters.apply(user.Login)
	if !b.isTargetMember(target) {
		target = "github"
	}
	return fmt.Sprintf(`<img src="https://github.com/%s.png" width="35">`, target)
}

func (b *builder) buildTable(xs ...string) string {
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

func (b *builder) buildIssueLinkTag(repo *github.Repo, issue *github.Issue) string {
	return fmt.Sprintf(`<a href="%s">%s#%d</a>`, issue.HTMLURL, repo.FullName, issue.Number)
}

func (b *builder) buildImportLabels(issue *github.Issue) []string {
	xs := []string{}
	for _, l := range issue.Labels {
		xs = append(xs, l.Name)
	}
	return xs
}

func (b *builder) isTargetMember(name string) bool {
	if !b.target.Private {
		return true
	}
	if strings.HasPrefix(b.target.FullName, name+"/") {
		return true
	}
	for _, m := range b.members {
		if m.Login == name {
			return true
		}
	}
	return false
}
