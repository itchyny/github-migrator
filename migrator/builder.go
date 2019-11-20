package migrator

import (
	"fmt"
	"html"
	"strings"
	"time"

	"github.com/itchyny/github-migrator/github"
	"github.com/itchyny/github-migrator/repo"
)

type builder struct {
	*migrator
	sourceCli      repo.Repo
	targetCli      repo.Repo
	source, target *github.Repo
	commentFilters commentFilters
	issue          *github.Issue
	pullReq        *github.PullReq
	comments       []*github.Comment
	events         []*github.Event
	commits        []*github.Commit
	commitDiff     string
	reviews        []*github.Review
	reviewComments []*github.ReviewComment
	projectByIDs   map[int]*github.Project
}

func buildImport(
	migrator *migrator,
	sourceCli, targetCli repo.Repo, sourceRepo, targetRepo *github.Repo, commentFilters commentFilters,
	issue *github.Issue, pullReq *github.PullReq,
	comments []*github.Comment, events []*github.Event,
	commits []*github.Commit, commitDiff string,
	reviews []*github.Review, reviewComments []*github.ReviewComment,
) (*github.Import, error) {
	return (&builder{
		migrator:       migrator,
		sourceCli:      sourceCli,
		targetCli:      targetCli,
		source:         sourceRepo,
		target:         targetRepo,
		commentFilters: commentFilters,
		issue:          issue,
		pullReq:        pullReq,
		comments:       comments,
		events:         events,
		commits:        commits,
		commitDiff:     commitDiff,
		reviews:        reviews,
		reviewComments: reviewComments,
	}).build()
}

func (b *builder) build() (*github.Import, error) {
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
		isMember, err := b.isTargetMember(target)
		if err != nil {
			return nil, err
		}
		if isMember {
			importIssue.Assignee = target
		}
	}
	comments, err := b.buildImportComments()
	if err != nil {
		return nil, err
	}
	return &github.Import{Issue: importIssue, Comments: comments}, nil
}

func (b *builder) buildImportBody() string {
	var suffix string
	if b.issue.Body != "" {
		suffix = "\n\n" + b.commentFilters.apply(b.issue.Body)
	}
	action := fmt.Sprintf("created the original %s<br>\n", b.issue.Type())
	if b.pullReq != nil {
		action += b.buildCompareLinkTag(b.target, b.pullReq.Base.SHA, b.pullReq.Head.SHA) +
			" " + b.buildPullRequestRefs() + "<br>\n"
	}
	action += "imported from " + buildIssueLinkTag(b.source, b.issue)
	tableRows := [][]string{
		[]string{
			b.buildImageTag(b.issue.User, 35),
			fmt.Sprintf("@%s %s", b.commentFilters.apply(b.issue.User.Login), action),
		},
	}
	if len(b.commitDiff) > 0 {
		tableRows = append(tableRows, []string{b.buildDiffDetails()})
	}
	if len(b.commits) > 0 {
		tableRows = append(tableRows, []string{b.buildCommitDetails()})
	}
	return b.buildTable(2, tableRows...) + suffix
}

func (b *builder) buildDiffDetails() string {
	summary := plural(b.pullReq.ChangedFiles, "file") + " changed"
	if b.pullReq.Additions > 0 {
		summary += ", " + plural(b.pullReq.Additions, "insertion") + "(+)"
	}
	if b.pullReq.Deletions > 0 {
		summary += ", " + plural(b.pullReq.Deletions, "deletion") + "(-)"
	}
	return b.buildDetails("  ", summary, "\n```diff\n"+
		escapeBackQuotes(truncateDiff(b.commitDiff))+
		"```\n")
}

func (b *builder) buildCommitDetails() string {
	summary := plural(b.pullReq.Commits, "commit")
	var commitRows [][]string
	for _, c := range b.commits {
		var dateString string
		committer := c.Committer
		if committer == nil {
			committer = c.Author
		}
		if committer == nil {
			committer = &github.User{Login: c.Commit.Committer.Name}
		}
		t, err := time.Parse(time.RFC3339, c.Commit.Committer.Date)
		if err == nil {
			dateString = t.Format(" on Mon 2, 2006")
		}
		commitRows = append(commitRows, []string{
			html.EscapeString(c.Commit.Message) + "<br>\n" +
				b.buildImageTag(committer, 16) +
				fmt.Sprintf(" @%s comitted%s", b.commentFilters.apply(committer.Login), dateString) +
				fmt.Sprintf(` <a href="%s">%s</a>`, b.commentFilters.apply(c.HTMLURL), c.SHA[:7]),
		})
	}
	return b.buildDetails("", summary, b.buildTable(1, commitRows...))
}

func (b *builder) buildImportComments() ([]*github.ImportComment, error) {
	issueComments := b.buildImportIssueComments()
	eventComments, err := b.buildImportEventComments()
	if err != nil {
		return nil, err
	}
	reviewComments := b.buildImportReviewComments()
	importReviews := b.buildImportReviews()
	return append(
		append(
			append(
				issueComments,
				eventComments...,
			),
			reviewComments...,
		),
		importReviews...,
	), nil
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

func (b *builder) buildPullRequestRefs() string {
	return fmt.Sprintf(
		"into <code>%s</code> from <code>%s</code>",
		html.EscapeString(b.pullReq.Base.Ref),
		html.EscapeString(b.pullReq.Head.Ref),
	)
}

func (b *builder) buildUserActionBody(user *github.User, action, body string) string {
	var suffix string
	if body != "" {
		suffix = "\n\n" + b.commentFilters.apply(body)
	}
	return b.buildTable(2, []string{
		b.buildImageTag(user, 35),
		fmt.Sprintf("@%s %s", b.commentFilters.apply(user.Login), action),
	}) + suffix
}

func (b *builder) buildImageTag(user *github.User, width int) string {
	target := b.commentFilters.apply(user.Login)
	if !b.isAvailableUser(target) {
		target = "github"
	}
	return fmt.Sprintf(`<img src="https://github.com/%s.png" width="%d">`, target, width)
}

func (b *builder) buildTable(width int, xss ...[]string) string {
	s := new(strings.Builder)
	s.WriteString("<table>\n")
	for i, xs := range xss {
		if i > 0 {
			s.WriteString("<tr></tr>\n")
		}
		s.WriteString("<tr>\n")
		for i, x := range xs {
			if i == len(xs)-1 && len(xs) < width {
				s.WriteString(fmt.Sprintf("  <td colspan=\"%d\">\n", width-i))
			} else if i == 0 && len(xs) == 2 && strings.HasPrefix(x, `<img src="`) && !strings.Contains(x, "\n") {
				s.WriteString("  <td width=\"60\">\n")
			} else {
				s.WriteString("  <td>\n")
			}
			x := makeIndent("    ", x)
			if !strings.HasSuffix(x, "\n") {
				x += "\n"
			}
			s.WriteString(x)
			s.WriteString("  </td>\n")
		}
		s.WriteString("</tr>\n")
	}
	s.WriteString("</table>\n")
	return s.String()
}

func (b *builder) buildDetails(indent, summary, details string) string {
	s := new(strings.Builder)
	s.WriteString(indent + "<details>\n")
	s.WriteString(fmt.Sprintf(indent+"  <summary>%s</summary>\n", summary))
	s.WriteString(makeIndent(indent+"  ", details))
	s.WriteString(indent + "</details>\n")
	return s.String()
}

func makeIndent(indent, str string) string {
	if strings.Contains(str, "```") {
		return str
	}
	xs := strings.Split(str, "\n")
	for i, x := range xs {
		if x == "" {
			break // avoid indented code block
		}
		xs[i] = indent + x
	}
	return strings.Join(xs, "\n")
}

func buildIssueLinkTag(repo *github.Repo, issue *github.Issue) string {
	return fmt.Sprintf(`<a href="%s">%s#%d</a>`, issue.HTMLURL, repo.FullName, issue.Number)
}

func (b *builder) buildCommitLinkTag(repo *github.Repo, sha string) string {
	return fmt.Sprintf(`<a href="%s/commit/%s">%s</a>`, repo.HTMLURL, sha, sha[:7])
}

func (b *builder) buildCompareLinkTag(repo *github.Repo, base, head string) string {
	return fmt.Sprintf(`<a href="%s/compare/%s...%s">%s...%s</a>`, repo.HTMLURL, base, head, base[:7], head[:7])
}

func (b *builder) buildImportLabels(issue *github.Issue) []string {
	xs := []string{}
	for _, l := range issue.Labels {
		xs = append(xs, l.Name)
	}
	return xs
}

func (b *builder) isAvailableUser(name string) bool {
	u, _ := b.lookupUser(name)
	return u != nil
}
