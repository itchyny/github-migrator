package migrator

import (
	"fmt"
	"html"
	"math"
	"strings"
	"time"

	"github.com/itchyny/github-migrator/github"
)

func (b *builder) buildImportEventComments() []*github.ImportComment {
	xs := make([]*github.ImportComment, 0, len(b.events))
	egs := groupEventsByCreated(b.events)
	for _, eg := range egs {
		if body := b.buildImportEventGroupBody(eg); body != "" {
			xs = append(xs, &github.ImportComment{
				Body:      b.buildUserActionBody(getEventUser(eg[0]), body, ""),
				CreatedAt: eg[0].CreatedAt,
			})
		}
	}
	return xs
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
		"closed":                1,
		"merged":                1,
		"reopened":              1,
		"labeled":               2,
		"unlabeled":             2,
		"rename":                3,
		"head_ref_deleted":      4,
		"head_ref_restored":     4,
		"head_ref_force_pushed": 5,
		"locked":                6,
		"unlocked":              6,
		"assigned":              7,
		"unassigned":            7,
		"review_requested":      8,
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

func (b *builder) buildImportEventGroupBody(eg []*github.Event) string {
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
					b.buildCommitLinkTag(b.target, e.CommitID),
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
		case "head_ref_force_pushed":
			actions = append(actions,
				fmt.Sprintf(
					"force-pushed the <code>%s</code> branch",
					html.EscapeString(b.pullReq.Head.Ref),
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
		case "assigned", "unassigned":
			if len(eg) == 1 && len(e.Assignees) <= 1 && e.Assigner.Login == e.Assignee.Login {
				if e.Event == "assigned" {
					return "self-assigned this"
				}
				return "removed their assignment"
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
		case "review_requested":
			if len(eg) == 1 && len(e.Reviewers) <= 1 && e.Actor.Login == e.Reviewer.Login {
				return "self-requested a review"
			}
			var targets []*github.User
			if len(e.Reviewers) > 0 {
				for _, u := range e.Reviewers {
					targets = append(targets, u)
				}
			} else {
				targets = append(targets, e.Reviewer)
			}
			actions = append(actions, "requested a review from "+b.mentionAll(targets))
		case "referenced", "mentioned", "subscribed":
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
		return action
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
	return action
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
