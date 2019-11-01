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

func listLabelsPath(repo string) string {
	return newPath(fmt.Sprintf("/repos/%s/labels", repo)).
		String()
}

// ListLabels lists the labels of an issue.
func (c *client) ListLabels(repo string) Labels {
	ls := make(chan interface{})
	go func() {
		defer close(ls)
		path := c.url(listLabelsPath(repo))
		for {
			xs, next, err := c.listLabels(path)
			if err != nil {
				ls <- err
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

// UpdateLabelParams represents the paramter for UpdateLabel API.
type UpdateLabelParams struct {
	Name        string `json:"new_name"`
	Description string `json:"description"`
	Color       string `json:"color"`
}

func updateLabelsPath(repo, name string) string {
	return newPath(fmt.Sprintf("/repos/%s/labels/%s", repo, name)).
		String()
}

func (c *client) UpdateLabel(repo, name string, params *UpdateLabelParams) (*Label, error) {
	bs, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	body := bytes.NewReader(bs)
	res, err := c.patch(c.url(updateLabelsPath(repo, name)), body)
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
