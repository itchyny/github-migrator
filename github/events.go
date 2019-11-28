package github

import (
	"fmt"
	"io"
)

// Event represents an event.
type Event struct {
	ID              int                   `json:"id"`
	Actor           *User                 `json:"actor"`
	Event           string                `json:"event"`
	Label           *EventLabel           `json:"label"`
	CommitID        string                `json:"commit_id"`
	Rename          *EventRename          `json:"rename"`
	LockReason      string                `json:"lock_reason"`
	Assignee        *User                 `json:"assignee"`
	Assignees       []*User               `json:"assignees"`
	Assigner        *User                 `json:"assigner"`
	Reviewer        *User                 `json:"requested_reviewer"`
	Reviewers       []*User               `json:"requested_reviewers"`
	RequestedTeam   *EventTeam            `json:"requested_team"`
	DismissedReview *EventDismissedReview `json:"dismissed_review"`
	ProjectCard     *EventProjectCard     `json:"project_card"`
	Milestone       *EventMilestone       `json:"milestone"`
	CreatedAt       string                `json:"created_at"`
}

// EventLabel ...
type EventLabel struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}

// EventRename ...
type EventRename struct {
	From string `json:"from"`
	To   string `json:"to"`
}

// EventTeam ...
type EventTeam struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// EventDismissedReview ...
type EventDismissedReview struct {
	State            string `json:"state"`
	ReviewID         int    `json:"review_id"`
	DismissalMessage string `json:"dismissal_message"`
}

// EventProjectCard ...
type EventProjectCard struct {
	ID                 int    `json:"id"`
	ProjectID          int    `json:"project_id"`
	ColumnName         string `json:"column_name"`
	PreviousColumnName string `json:"previous_column_name"`
}

// EventMilestone ...
type EventMilestone struct {
	Title string `json:"title"`
}

// Events represents a collection of events.
type Events <-chan interface{}

// Next emits the next Event.
func (es Events) Next() (*Event, error) {
	for x := range es {
		switch x := x.(type) {
		case error:
			return nil, x
		case *Event:
			return x, nil
		}
		break
	}
	return nil, io.EOF
}

// EventsFromSlice creates Events from a slice.
func EventsFromSlice(xs []*Event) Events {
	es := make(chan interface{})
	go func() {
		defer close(es)
		for _, e := range xs {
			es <- e
		}
	}()
	return es
}

// EventsToSlice collects Events.
func EventsToSlice(es Events) ([]*Event, error) {
	xs := []*Event{}
	for {
		e, err := es.Next()
		if err != nil {
			if err != io.EOF {
				return nil, err
			}
			return xs, nil
		}
		xs = append(xs, e)
	}
}

// ListEvents lists the events of an issue.
func (c *client) ListEvents(repo string, issueNumber int) Events {
	es := make(chan interface{})
	go func() {
		defer close(es)
		path := c.url(fmt.Sprintf("/repos/%s/issues/%d/events?per_page=100", repo, issueNumber))
		for {
			var xs []*Event
			next, err := c.getList(path, &xs)
			if err != nil {
				es <- fmt.Errorf("ListEvents %s/issues/%d: %w", repo, issueNumber, err)
				break
			}
			for _, x := range xs {
				es <- x
			}
			if next == "" {
				break
			}
			path = next
		}
	}()
	return Events(es)
}
