package github

import (
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
func (ls Labels) Next() (*Label, error) {
	for x := range ls {
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
	ls := make(chan interface{})
	go func() {
		defer close(ls)
		for _, l := range xs {
			ls <- l
		}
	}()
	return ls
}

// LabelsToSlice collects Labels.
func LabelsToSlice(ls Labels) ([]*Label, error) {
	xs := []*Label{}
	for {
		l, err := ls.Next()
		if err != nil {
			if err != io.EOF {
				return nil, err
			}
			return xs, nil
		}
		xs = append(xs, l)
	}
}

// ListLabels lists the labels of an issue.
func (c *client) ListLabels(repo string) Labels {
	ls := make(chan interface{})
	go func() {
		defer close(ls)
		path := c.url(fmt.Sprintf("/repos/%s/labels?per_page=100", repo))
		for {
			var xs []*Label
			next, err := c.getList(path, &xs)
			if err != nil {
				ls <- fmt.Errorf("ListLabels %s: %w", repo, err)
				break
			}
			for _, x := range xs {
				ls <- x
			}
			if next == "" {
				break
			}
			path = next
		}
	}()
	return Labels(ls)
}

// CreateLabelParams represents the paramter for CreateLabel API.
type CreateLabelParams struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Color       string `json:"color"`
}

func (c *client) CreateLabel(repo string, params *CreateLabelParams) (*Label, error) {
	var r Label
	if err := c.post(c.url(fmt.Sprintf("/repos/%s/labels", repo)), params, &r); err != nil {
		return nil, fmt.Errorf("CreateLabel %s: %w", fmt.Sprintf("%s/labels", repo), err)
	}
	return &r, nil
}

// UpdateLabelParams represents the paramter for UpdateLabel API.
type UpdateLabelParams struct {
	Name        string `json:"new_name"`
	Description string `json:"description"`
	Color       string `json:"color"`
}

func (c *client) UpdateLabel(repo, name string, params *UpdateLabelParams) (*Label, error) {
	var r Label
	if err := c.patch(c.url(fmt.Sprintf("/repos/%s/labels/%s", repo, name)), params, &r); err != nil {
		return nil, fmt.Errorf("UpdateLabel %s: %w", fmt.Sprintf("%s/labels/%s", repo, name), err)
	}
	return &r, nil
}
