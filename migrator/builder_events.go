package migrator

import (
	"fmt"
	"html"
	"math"
	"strings"
	"time"

	"github.com/itchyny/github-migrator/github"
)

func (b *builder) buildImportEventComments() ([]*github.ImportComment, error) {
	xs := make([]*github.ImportComment, 0, len(b.events))
	egs := groupEventsByCreated(b.events)
	for _, eg := range egs {
		body, err := b.buildImportEventGroupBody(eg)
		if err != nil {
			return nil, err
		}
		if body != "" {
			xs = append(xs, &github.ImportComment{
				Body:      b.buildUserActionBody(getEventUser(eg[0]), body, ""),
				CreatedAt: eg[0].CreatedAt,
			})
		}
	}
	return xs, nil
}

func getEventUser(e *github.Event) *github.User {
	switch e.Event {
	case "assigned", "unassigned":
		return e.Assigner
	default:
		return e.Actor
	}
}

func groupEventsByCreated(xs []*github.Event) [][]*github.Event {
	ess := make([][]*github.Event, 0, len(xs))
	eventGroupTypes := map[string]int{
		"closed":                   1,
		"merged":                   1,
		"reopened":                 1,
		"labeled":                  2,
		"unlabeled":                2,
		"rename":                   3,
		"head_ref_deleted":         4,
		"head_ref_restored":        4,
		"head_ref_force_pushed":    5,
		"base_ref_force_pushed":    5,
		"locked":                   6,
		"unlocked":                 6,
		"pinned":                   7,
		"unpinned":                 7,
		"assigned":                 8,
		"unassigned":               8,
		"review_requested":         9,
		"review_request_removed":   9,
		"review_dismissed":         9,
		"ready_for_review":         9,
		"converted_note_to_issue":  10,
		"added_to_project":         10,
		"moved_columns_in_project": 10,
		"removed_from_project":     10,
		"milestoned":               11,
		"demilestoned":             11,
		"deployed":                 12,
	}
	for _, x := range xs {
		var appended bool
		for i, es := range ess {
			if getEventUser(es[0]).Login == getEventUser(x).Login &&
				nearTime(es[0].CreatedAt, x.CreatedAt) &&
				eventGroupTypes[es[0].Event] == eventGroupTypes[x.Event] {
				ess[i] = append(ess[i], x)
				appended = true
				break
			}
		}
		if appended {
			continue
		}
		ess = append(ess, []*github.Event{x})
	}
	return ess
}

func nearTime(s1, s2 string) bool {
	t1, err := time.Parse(time.RFC3339, s1)
	if err != nil {
		panic(err)
	}
	t2, err := time.Parse(time.RFC3339, s2)
	if err != nil {
		panic(err)
	}
	diff := t1.Sub(t2)
	return math.Abs(float64(diff)) < float64(10*time.Second)
}

const (
	actionClosed = 1 << iota
	actionMerged
	actionReopened
)

func (b *builder) buildImportEventGroupBody(eg []*github.Event) (string, error) {
	var actions []string
	var merged bool
	var addedLabels []string
	var removedLabels []string

	for _, e := range eg {
		switch e.Event {
		case "closed":
			if !merged {
				if b.pullReq == nil {
					actions = append(actions, "closed the issue")
				} else {
					actions = append(actions, "closed the pull request without merging")
				}
			}
		case "merged":
			merged = true
			actions = append(actions,
				fmt.Sprintf(
					"merged the pull request<br>\ncommit %s ",
					b.buildCommitLinkTag(b.targetRepo, e.CommitID),
				)+b.buildPullRequestRefs(),
			)
		case "reopened":
			actions = append(actions, fmt.Sprintf("reopened the %s", b.issue.Type()))
		case "labeled":
			addedLabels = append(addedLabels, e.Label.Name)
		case "unlabeled":
			removedLabels = append(removedLabels, e.Label.Name)
		case "renamed":
			actions = append(actions,
				fmt.Sprintf(
					"changed the title <b><s>%s</s></b> <b>%s</b>",
					html.EscapeString(e.Rename.From), html.EscapeString(e.Rename.To),
				),
			)
		case "head_ref_deleted":
			actions = append(actions,
				fmt.Sprintf(
					"deleted the <code>%s</code> branch",
					html.EscapeString(b.pullReq.Head.Ref),
				),
			)
		case "head_ref_restored":
			actions = append(actions,
				fmt.Sprintf(
					"restored the <code>%s</code> branch",
					html.EscapeString(b.pullReq.Head.Ref),
				),
			)
		case "head_ref_force_pushed", "base_ref_force_pushed":
			ref := b.pullReq.Head.Ref
			if e.Event == "base_ref_force_pushed" {
				ref = b.pullReq.Base.Ref
			}
			actions = append(actions,
				fmt.Sprintf(
					"force-pushed the <code>%s</code> branch",
					html.EscapeString(ref),
				),
			)
		case "locked":
			actions = append(actions,
				fmt.Sprintf(
					"locked as <b>%s</b> and limited conversation to collaborators",
					html.EscapeString(strings.ReplaceAll(e.LockReason, "-", " ")),
				),
			)
		case "unlocked":
			actions = append(actions, "unlocked this conversation")
		case "pinned", "unpinned":
			actions = append(actions, e.Event+` this issue`)
		case "assigned", "unassigned":
			if len(eg) == 1 && len(e.Assignees) <= 1 && e.Assigner.Login == e.Assignee.Login {
				if e.Event == "assigned" {
					return "self-assigned this", nil
				}
				return "removed their assignment", nil
			}
			var targets []*github.User
			if len(e.Assignees) > 0 {
				for _, u := range e.Assignees {
					targets = append(targets, u)
				}
			} else {
				targets = append(targets, e.Assignee)
			}
			actions = append(actions, e.Event+" "+b.mentionAll(targets))
		case "review_requested", "review_request_removed":
			var actionStr string
			if e.Event == "review_requested" {
				actionStr = "requested a review"
			} else {
				actionStr = "removed the request for review"
			}
			if e.RequestedTeam != nil {
				actions = append(actions,
					fmt.Sprintf(
						`%s from <b>%s</b>`,
						actionStr,
						b.commentFilters.apply(e.RequestedTeam.Name),
					),
				)
				break
			}
			if len(eg) == 1 && len(e.Reviewers) <= 1 && e.Actor.Login == e.Reviewer.Login {
				if e.Event == "review_requested" {
					return "self-requested a review", nil
				}
				return "removed their request for review", nil
			}
			var targets []*github.User
			if len(e.Reviewers) > 0 {
				for _, u := range e.Reviewers {
					targets = append(targets, u)
				}
			} else {
				targets = append(targets, e.Reviewer)
			}
			actions = append(actions, actionStr+" from "+b.mentionAll(targets))
		case "review_dismissed":
			var target *github.User
			for _, r := range b.reviews {
				if r.ID == e.DismissedReview.ReviewID {
					target = r.User
					break
				}
			}
			if target != nil {
				actions = append(actions,
					fmt.Sprintf(
						`dismissed @%s's review<br>%s`,
						b.commentFilters.apply(target.Login),
						html.EscapeString(e.DismissedReview.DismissalMessage),
					),
				)
			} else {
				actions = append(actions,
					fmt.Sprintf(
						`dismissed a review<br>%s`,
						html.EscapeString(e.DismissedReview.DismissalMessage),
					),
				)
			}
		case "ready_for_review":
			actions = append(actions, "marked this pull request as ready for review")
		case "converted_note_to_issue":
			p, err := b.getProject(e.ProjectCard.ProjectID)
			if err != nil {
				return "", err
			}
			actions = append(actions,
				fmt.Sprintf(
					`created this issue from a note in <b><a href="%s">%s</a></b> (<code>%s</code>)`,
					b.lookupMigratedProject(p).HTMLURL, html.EscapeString(p.Name),
					html.EscapeString(e.ProjectCard.ColumnName),
				),
			)
		case "added_to_project":
			p, err := b.getProject(e.ProjectCard.ProjectID)
			if err != nil {
				return "", err
			}
			actions = append(actions,
				fmt.Sprintf(
					`added this to <code>%s</code> in <b><a href="%s">%s</a></b>`,
					html.EscapeString(e.ProjectCard.ColumnName),
					b.lookupMigratedProject(p).HTMLURL, html.EscapeString(p.Name),
				),
			)
		case "moved_columns_in_project":
			p, err := b.getProject(e.ProjectCard.ProjectID)
			if err != nil {
				return "", err
			}
			actions = append(actions,
				fmt.Sprintf(
					`moved this from <code>%s</code> to <code>%s</code> in <b><a href="%s">%s</a></b>`,
					html.EscapeString(e.ProjectCard.PreviousColumnName),
					html.EscapeString(e.ProjectCard.ColumnName),
					b.lookupMigratedProject(p).HTMLURL, html.EscapeString(p.Name),
				),
			)
		case "removed_from_project":
			p, err := b.getProject(e.ProjectCard.ProjectID)
			if err != nil {
				return "", err
			}
			actions = append(actions,
				fmt.Sprintf(
					`removed this from <code>%s</code> in <b><a href="%s">%s</a></b>`,
					html.EscapeString(e.ProjectCard.ColumnName),
					b.lookupMigratedProject(p).HTMLURL, html.EscapeString(p.Name),
				),
			)
		case "milestoned", "demilestoned":
			var actionStr string
			if e.Event == "milestoned" {
				actionStr = "added this to"
			} else {
				actionStr = "removed this from"
			}
			actions = append(actions,
				fmt.Sprintf(
					`%s the <b><a href="%s">%s</a></b> milestone`,
					actionStr,
					b.milestoneByTitle[e.Milestone.Title].HTMLURL,
					html.EscapeString(e.Milestone.Title),
				),
			)
		case "deployed":
			actions = append(actions, `deployed this`)
		case "referenced", "mentioned", "comment_deleted",
			"subscribed", "unsubscribed", "base_ref_changed":
		default:
			fmt.Printf("%#v\n", e)
			panic(e.Event)
		}
	}

	var action string
	if len(actions) > 0 {
		for i, a := range actions {
			if i > 0 {
				if i == len(actions)-1 {
					action += " and "
				} else {
					action += ", "
				}
			}
			action += a
		}
		return action, nil
	}

	if len(addedLabels) > 0 {
		action += "added " + quoteLabels(addedLabels)
	}
	if len(removedLabels) > 0 {
		if action != "" {
			action += " and "
		}
		action += "removed " + quoteLabels(removedLabels)
	}
	if len(addedLabels) > 0 || len(removedLabels) > 0 {
		action += pluralUnit(len(addedLabels)+len(removedLabels), " label")
	}
	return action, nil
}

func (b *builder) mentionAll(users []*github.User) string {
	var s string
	for i, u := range users {
		if i > 0 {
			s += " "
		}
		s += "@" + b.commentFilters.apply(u.Login)
	}
	return s
}

func quoteLabels(xs []string) string {
	ys := make([]string, len(xs))
	for i, x := range xs {
		ys[i] = "<b><code>" + html.EscapeString(x) + "</code></b>"
	}
	return strings.Join(ys, " ")
}

func (b *builder) lookupMigratedProject(orig *github.Project) *github.Project {
	if !strings.HasPrefix(orig.HTMLURL, b.sourceRepo.HTMLURL+"/projects/") {
		return orig
	}
	found := lookupProject(b.targetProjects, orig)
	if found != nil {
		return found
	}
	return orig
}
