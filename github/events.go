package github

import (
	"encoding/json"
	"fmt"
	"io"
)

// Event represents an event.
type Event struct {
	ID        int         `json:"id"`
	Actor     *User       `json:"actor"`
	Event     string      `json:"event"`
	Label     *EventLabel `json:"label"`
	CreatedAt string      `json:"created_at"`
}

// EventLabel ...
type EventLabel struct {
	Name  string `json:"name"`
	Color string `json:"color"`
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

func listEventsPath(repo string, issueNumber int) string {
	return newPath(fmt.Sprintf("/repos/%s/issues/%d/events", repo, issueNumber)).
		String()
}

// ListEvents lists the events of an issue.
func (c *client) ListEvents(repo string, issueNumber int) Events {
	es := make(chan interface{})
	go func() {
		defer close(es)
		path := c.url(listEventsPath(repo, issueNumber))
		for {
			xs, next, err := c.listEvents(path)
			if err != nil {
				es <- err
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

func (c *client) listEvents(path string) ([]*Event, string, error) {
	res, err := c.get(path)
	if err != nil {
		return nil, "", err
	}
	defer res.Body.Close()

	var r []*Event
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, "", err
	}

	return r, getNext(res.Header), nil
}
