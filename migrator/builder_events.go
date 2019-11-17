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
				Body:      b.buildUserActionBody(eg[0].Actor, body, ""),
				CreatedAt: eg[0].CreatedAt,
			})
		}
	}
	return xs
}

func groupEventsByCreated(xs []*github.Event) [][]*github.Event {
	ess := make([][]*github.Event, 0, len(xs))
	mergeTypes := map[string]int{
		"closed":    1,
		"merged":    1,
		"reopened":  1,
		"labeled":   2,
		"unlabeled": 2,
	}
	for _, x := range xs {
		var appended bool
		for i, es := range ess {
			if es[0].Actor.Login == x.Actor.Login &&
				nearTime(es[0].CreatedAt, x.CreatedAt) &&
				mergeTypes[es[0].Event] == mergeTypes[x.Event] {
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
		}
	}

	var action string
	if len(actions) > 0 {
		for i, a := range actions {
			if i > 0 {
				action += ", "
				if i == len(actions)-1 {
					action += "and "
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

func quoteLabels(xs []string) string {
	ys := make([]string, len(xs))
	for i, x := range xs {
		ys[i] = "<code>" + html.EscapeString(x) + "</code>"
	}
	return strings.Join(ys, " ")
}
