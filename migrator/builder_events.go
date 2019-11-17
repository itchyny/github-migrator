package migrator

import (
	"html"
	"strings"

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
		"labeled":   1,
		"unlabeled": 1,
	}
	for _, x := range xs {
		var appended bool
		for i, es := range ess {
			if es[0].Actor.Login == x.Actor.Login &&
				es[0].CreatedAt == x.CreatedAt &&
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

func (b *builder) buildImportEventGroupBody(eg []*github.Event) string {
	var addedLabels []string
	var removedLabels []string
	for _, e := range eg {
		switch e.Event {
		case "labeled":
			addedLabels = append(addedLabels, e.Label.Name)
		case "unlabeled":
			removedLabels = append(removedLabels, e.Label.Name)
		}
	}

	var action string
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
