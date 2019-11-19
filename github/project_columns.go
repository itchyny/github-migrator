package github

import (
	"bytes"
	"encoding/json"
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

func listProjectColumnsPath(projectID int) string {
	return newPath(fmt.Sprintf("/projects/%d/columns", projectID)).
		query("per_page", "100").
		String()
}

// ListProjectColumns lists the project columns.
func (c *client) ListProjectColumns(projectID int) ProjectColumns {
	ps := make(chan interface{})
	go func() {
		defer close(ps)
		path := c.url(listProjectColumnsPath(projectID))
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

func getProjectColumnPath(projectColumnID int) string {
	return newPath(fmt.Sprintf("/projects/columns/%d", projectColumnID)).
		String()
}

type projectColumnOrError struct {
	ProjectColumn
	Message string `json:"message"`
}

func (c *client) GetProjectColumn(projectColumnID int) (*ProjectColumn, error) {
	res, err := c.get(c.url(getProjectColumnPath(projectColumnID)))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var r projectColumnOrError
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, err
	}

	if r.Message != "" {
		return nil, fmt.Errorf("GetProjectColumn %s: %s", fmt.Sprintf("projects/columns/%d", projectColumnID), r.Message)
	}

	return &r.ProjectColumn, nil
}

// CreateProjectColumnParams represents the paramter for CreateProjectColumn API.
type CreateProjectColumnParams struct {
	Name string `json:"name"`
}

func createProjectColumnPath(projectID int) string {
	return newPath(fmt.Sprintf("/projects/%d/columns", projectID)).
		String()
}

// CreateProjectColumn creates a project column.
func (c *client) CreateProjectColumn(projectID int, name string) (*ProjectColumn, error) {
	bs, err := json.Marshal(map[string]string{"name": name})
	if err != nil {
		return nil, err
	}
	body := bytes.NewReader(bs)
	res, err := c.post(c.url(createProjectColumnPath(projectID)), body)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var r projectColumnOrError
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, err
	}

	if r.Message != "" {
		return nil, fmt.Errorf("CreateProjectColumn %s: %s", fmt.Sprintf("projects/%d/columns", projectID), r.Message)
	}

	return &r.ProjectColumn, nil
}

func updateProjectColumnPath(projectColumnID int) string {
	return newPath(fmt.Sprintf("/projects/columns/%d", projectColumnID)).
		String()
}

// UpdateProjectColumn updates the project column.
func (c *client) UpdateProjectColumn(projectColumnID int, name string) (*ProjectColumn, error) {
	bs, err := json.Marshal(map[string]string{"name": name})
	if err != nil {
		return nil, err
	}
	body := bytes.NewReader(bs)
	res, err := c.patch(c.url(updateProjectColumnPath(projectColumnID)), body)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var r projectColumnOrError
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, err
	}

	if r.Message != "" {
		return nil, fmt.Errorf("UpdateProjectColumn %s: %s", fmt.Sprintf("projects/columns/%d", projectColumnID), r.Message)
	}

	return &r.ProjectColumn, nil
}
