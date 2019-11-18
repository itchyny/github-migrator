package github

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"
)

// Project represents a project.
type Project struct {
	ID         int          `json:"id"`
	Name       string       `json:"name"`
	Body       string       `json:"body"`
	Number     int          `json:"number"`
	State      ProjectState `json:"state"`
	OwnerURL   string       `json:"owner_url"`
	URL        string       `json:"url"`
	HTMLURL    string       `json:"html_url"`
	ColumnsURL string       `json:"columns_url"`
	Creator    *User        `json:"creator"`
	CreatedAt  string       `json:"created_at"`
	UpdatedAt  string       `json:"updated_at"`
}

// ProjectState ...
type ProjectState int

// ProjectState ...
const (
	ProjectStateOpen ProjectState = iota + 1
	ProjectStateClosed
)

var stringToProjectState = map[string]ProjectState{
	"open":   ProjectStateOpen,
	"closed": ProjectStateClosed,
}

var projectStateToString = map[ProjectState]string{
	ProjectStateOpen:   "open",
	ProjectStateClosed: "closed",
}

// UnmarshalJSON implements json.Unmarshaler
func (t *ProjectState) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	if x, ok := stringToProjectState[s]; ok {
		*t = x
		return nil
	}
	return fmt.Errorf("unknown project state: %q", s)
}

// MarshalJSON implements json.Marshaler
func (t ProjectState) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

// String implements Stringer
func (t ProjectState) String() string {
	return projectStateToString[t]
}

// GoString implements GoString
func (t ProjectState) GoString() string {
	return strconv.Quote(t.String())
}

// Projects represents a collection of projects.
type Projects <-chan interface{}

// Next emits the next Project.
func (ps Projects) Next() (*Project, error) {
	for x := range ps {
		switch x := x.(type) {
		case error:
			return nil, x
		case *Project:
			return x, nil
		}
		break
	}
	return nil, io.EOF
}

// ProjectsFromSlice creates Projects from a slice.
func ProjectsFromSlice(xs []*Project) Projects {
	ps := make(chan interface{})
	go func() {
		defer close(ps)
		for _, p := range xs {
			ps <- p
		}
	}()
	return ps
}

// ProjectsToSlice collects Projects.
func ProjectsToSlice(ps Projects) ([]*Project, error) {
	xs := []*Project{}
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

// ListProjectsParams represents the paramter for ListProjects API.
type ListProjectsParams struct {
	State ListProjectsParamState
}

// ListProjectsParamState ...
type ListProjectsParamState int

// ListProjectsParamState ...
const (
	ListProjectsParamStateDefault ListProjectsParamState = iota + 1
	ListProjectsParamStateOpen
	ListProjectsParamStateClosed
	ListProjectsParamStateAll
)

func (f ListProjectsParamState) String() string {
	switch f {
	case ListProjectsParamStateOpen:
		return "open"
	case ListProjectsParamStateClosed:
		return "closed"
	case ListProjectsParamStateAll:
		return "all"
	default:
		return ""
	}
}

func listProjectsPath(repo string, params *ListProjectsParams) string {
	return newPath("/repos/"+repo+"/projects").
		query("state", params.State.String()).
		query("per_page", "100").
		String()
}

// ListProjects lists the projects.
func (c *client) ListProjects(repo string, params *ListProjectsParams) Projects {
	ps := make(chan interface{})
	go func() {
		defer close(ps)
		path := c.url(listProjectsPath(repo, params))
		for {
			xs, next, err := c.listProjects(path)
			if err != nil {
				ps <- err
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
	return Projects(ps)
}

func (c *client) listProjects(path string) ([]*Project, string, error) {
	res, err := c.get(path)
	if err != nil {
		return nil, "", err
	}
	defer res.Body.Close()

	var r []*Project
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, "", err
	}

	return r, getNext(res.Header), nil
}

func getProjectPath(repo string, projectID int) string {
	return newPath(fmt.Sprintf("/repos/%s/projects/%d", repo, projectID)).
		String()
}

type projectOrError struct {
	Project
	Message string `json:"message"`
}

func (c *client) GetProject(repo string, projectID int) (*Project, error) {
	res, err := c.get(c.url(getProjectPath(repo, projectID)))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var r projectOrError
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, err
	}

	if r.Message != "" {
		return nil, fmt.Errorf("%s: %s", r.Message, "/projects/"+fmt.Sprint(projectID))
	}

	return &r.Project, nil
}
