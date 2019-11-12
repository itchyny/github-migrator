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
	pullReq        *github.PullReq
	comments       []*github.Comment
	reviews        []*github.Review
	reviewComments []*github.ReviewComment
	members        []*github.Member
}

func buildImport(
	sourceRepo, targetRepo *github.Repo, commentFilters commentFilters,
	issue *github.Issue, pullReq *github.PullReq,
	comments []*github.Comment, reviews []*github.Review, reviewComments []*github.ReviewComment,
	members []*github.Member,
) *github.Import {
	return (&builder{
		source:         sourceRepo,
		target:         targetRepo,
		commentFilters: commentFilters,
		issue:          issue,
		pullReq:        pullReq,
		comments:       comments,
		reviews:        reviews,
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
		Closed:    b.issue.State != github.IssueStateOpen,
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
	return b.buildUserActionBody(
		b.issue.User,
		fmt.Sprintf(
			"created the original %s<br>imported from %s",
			b.issue.Type(),
			b.buildIssueLinkTag(b.source, b.issue),
		),
		b.issue.Body,
	)
}

func (b *builder) buildImportComments() []*github.ImportComment {
	cs := append(
		append(
			b.buildImportIssueComments(),
			b.buildImportReviewComments()...,
		),
		b.buildImportReviews()...,
	)
	if c := b.buildClosedComment(); c != nil {
		cs = append(cs, c)
	}
	return cs
}

func (b *builder) buildImportIssueComments() []*github.ImportComment {
	xs := make([]*github.ImportComment, len(b.comments))
	for i, c := range b.comments {
		xs[i] = &github.ImportComment{
			Body:      b.buildUserActionBody(c.User, "commented", c.Body),
			CreatedAt: c.CreatedAt,
		}
	}
	return xs
}

func (b *builder) buildImportReviews() []*github.ImportComment {
	var xs []*github.ImportComment
	for _, c := range b.reviews {
		var action string
		if c.State == github.ReviewStateApproved {
			action = "approved"
		} else if c.State == github.ReviewStateChangesRequested {
			action = "requested changes"
		} else {
			continue
		}
		xs = append(xs, &github.ImportComment{
			Body:      b.buildUserActionBody(c.User, action, c.Body),
			CreatedAt: c.SubmittedAt,
		})
	}
	return xs
}

func (b *builder) buildImportReviewComments() []*github.ImportComment {
	var xs []*github.ImportComment
	indexByID := make(map[int]int)
	for _, c := range b.reviewComments {
		if i, ok := indexByID[c.InReplyToID]; ok {
			indexByID[c.ID] = i
			xs[i].Body += "\n\n" + b.buildUserActionBody(c.User, "commented", c.Body)
			continue
		}
		indexByID[c.ID] = len(xs)
		diffBody := strings.Join([]string{"```diff", "# " + c.Path, c.DiffHunk, "```"}, "\n")
		xs = append(xs, &github.ImportComment{
			Body:      diffBody + "\n\n" + b.buildUserActionBody(c.User, "commented", c.Body),
			CreatedAt: c.CreatedAt,
		})
	}
	return xs
}

func (b *builder) buildClosedComment() *github.ImportComment {
	if b.issue.State == github.IssueStateOpen {
		return nil
	}
	var user *github.User
	var action string
	var closedAt string
	if b.pullReq == nil {
		user = b.issue.ClosedBy
		action = "closed the issue"
		closedAt = b.issue.ClosedAt
	} else if b.pullReq.MergedBy != nil {
		user = b.pullReq.MergedBy
		action = "merged the pull request"
		closedAt = b.pullReq.MergedAt
	} else {
		user = b.issue.ClosedBy
		action = "closed the pull request without merging"
		closedAt = b.issue.ClosedAt
	}
	return &github.ImportComment{
		Body:      b.buildUserActionBody(user, action, ""),
		CreatedAt: closedAt,
	}
}

func (b *builder) buildUserActionBody(user *github.User, action, body string) string {
	var suffix string
	if body != "" {
		suffix = "\n\n" + b.commentFilters.apply(body)
	}
	return b.buildTable(
		b.buildImageTag(user),
		fmt.Sprintf("@%s %s", b.commentFilters.apply(user.Login), action),
	) + suffix
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
