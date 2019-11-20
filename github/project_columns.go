package github

import (
	"fmt"
	"io"
)

// ProjectColumn represents a project column.
type ProjectColumn struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	URL       string `json:"url"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// ProjectColumns represents a collection of project columns.
type ProjectColumns <-chan interface{}

// Next emits the next ProjectColumn.
func (ps ProjectColumns) Next() (*ProjectColumn, error) {
	for x := range ps {
		switch x := x.(type) {
		case error:
			return nil, x
		case *ProjectColumn:
			return x, nil
		}
		break
	}
	return nil, io.EOF
}

// ProjectColumnsFromSlice creates ProjectColumns from a slice.
func ProjectColumnsFromSlice(xs []*ProjectColumn) ProjectColumns {
	ps := make(chan interface{})
	go func() {
		defer close(ps)
		for _, p := range xs {
			ps <- p
		}
	}()
	return ps
}

// ProjectColumnsToSlice collects ProjectColumns.
func ProjectColumnsToSlice(ps ProjectColumns) ([]*ProjectColumn, error) {
	xs := []*ProjectColumn{}
	for {
		p, err := ps.Next()
		if err != nil {
			if err != io.EOF {
				return nil, err
			}
			return xs, nil
		}
		xs = append(xs, p)
	}
}

// ListProjectColumns lists the project columns.
func (c *client) ListProjectColumns(projectID int) ProjectColumns {
	ps := make(chan interface{})
	go func() {
		defer close(ps)
		path := c.url(fmt.Sprintf("/projects/%d/columns?per_page=100", projectID))
		for {
			var xs []*ProjectColumn
			next, err := c.getList(path, &xs)
			if err != nil {
				ps <- fmt.Errorf("ListProjectColumns %d: %w", projectID, err)
				break
			}
			for _, x := range xs {
				ps <- x
			}
			if next == "" {
				break
			}
			path = next
		}
	}()
	return ProjectColumns(ps)
}

func (c *client) GetProjectColumn(projectColumnID int) (*ProjectColumn, error) {
	var r ProjectColumn
	if err := c.get(c.url(fmt.Sprintf("/projects/columns/%d", projectColumnID)), &r); err != nil {
		return nil, fmt.Errorf("GetProjectColumn %s: %w", fmt.Sprintf("projects/columns/%d", projectColumnID), err)
	}
	return &r, nil
}

// CreateProjectColumn creates a project column.
func (c *client) CreateProjectColumn(projectID int, name string) (*ProjectColumn, error) {
	var r ProjectColumn
	if err := c.post(c.url(fmt.Sprintf("/projects/%d/columns", projectID)), map[string]string{"name": name}, &r); err != nil {
		return nil, err
	}
	return &r, nil
}

// UpdateProjectColumn updates the project column.
func (c *client) UpdateProjectColumn(projectColumnID int, name string) (*ProjectColumn, error) {
	var r ProjectColumn
	if err := c.patch(c.url(fmt.Sprintf("/projects/columns/%d", projectColumnID)), map[string]string{"name": name}, &r); err != nil {
		return nil, fmt.Errorf("UpdateProjectColumn %s: %w", fmt.Sprintf("projects/columns/%d", projectColumnID), err)
	}
	return &r, nil
}
