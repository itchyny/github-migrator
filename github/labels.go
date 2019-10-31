package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

// Label represents a label.
type Label struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Color       string `json:"color"`
	Default     bool   `json:"default"`
}

// Labels represents a collection of labels.
type Labels <-chan interface{}

// Next emits the next Label.
func (cs Labels) Next() (*Label, error) {
	for x := range cs {
		switch x := x.(type) {
		case error:
			return nil, x
		case *Label:
			return x, nil
		}
		break
	}
	return nil, io.EOF
}

// LabelsFromSlice creates Labels from a slice.
func LabelsFromSlice(xs []*Label) Labels {
	cs := make(chan interface{})
	go func() {
		defer close(cs)
		for _, i := range xs {
			cs <- i
		}
	}()
	return cs
}

// LabelsToSlice collects Labels.
func LabelsToSlice(cs Labels) ([]*Label, error) {
	var xs []*Label
	for {
		i, err := cs.Next()
		if err != nil {
			if err != io.EOF {
				return nil, err
			}
			return xs, nil
		}
		xs = append(xs, i)
	}
}

func listLabelsPath(repo string) string {
	return newPath(fmt.Sprintf("/repos/%s/labels", repo)).
		String()
}

// ListLabels lists the labels of an issue.
func (c *client) ListLabels(repo string) Labels {
	cs := make(chan interface{})
	go func() {
		defer close(cs)
		path := c.url(listLabelsPath(repo))
		for {
			xs, next, err := c.listLabels(path)
			if err != nil {
				cs <- err
				break
			}
			for _, x := range xs {
				cs <- x
			}
			if next == "" {
				break
			}
			path = next
		}
	}()
	return Labels(cs)
}

func (c *client) listLabels(path string) ([]*Label, string, error) {
	res, err := c.get(path)
	if err != nil {
		return nil, "", err
	}
	defer res.Body.Close()

	var r []*Label
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, "", err
	}

	return r, getNext(res.Header), nil
}

// CreateLabelParams represents the paramter for CreateLabel API.
type CreateLabelParams struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Color       string `json:"color"`
}

func createLabelsPath(repo string) string {
	return newPath(fmt.Sprintf("/repos/%s/labels", repo)).
		String()
}

func (c *client) CreateLabel(repo string, params *CreateLabelParams) (*Label, error) {
	bs, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	body := bytes.NewReader(bs)
	res, err := c.post(c.url(createLabelsPath(repo)), body)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var r Label
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, err
	}

	return &r, nil
}
