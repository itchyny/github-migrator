package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
)

// ProjectCard represents a project card.
type ProjectCard struct {
	ID         int    `json:"id"`
	Note       string `json:"note"`
	Archived   bool   `json:"archived"`
	Creator    *User  `json:"creator"`
	ContentURL string `json:"content_url"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

// ProjectCards represents a collection of project cards.
type ProjectCards <-chan interface{}

// Next emits the next ProjectCard.
func (ps ProjectCards) Next() (*ProjectCard, error) {
	for x := range ps {
		switch x := x.(type) {
		case error:
			return nil, x
		case *ProjectCard:
			return x, nil
		}
		break
	}
	return nil, io.EOF
}

// ProjectCardsFromSlice creates ProjectCards from a slice.
func ProjectCardsFromSlice(xs []*ProjectCard) ProjectCards {
	ps := make(chan interface{})
	go func() {
		defer close(ps)
		for _, p := range xs {
			ps <- p
		}
	}()
	return ps
}

// ProjectCardsToSlice collects ProjectCards.
func ProjectCardsToSlice(ps ProjectCards) ([]*ProjectCard, error) {
	xs := []*ProjectCard{}
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

func listProjectCardsPath(columnID int) string {
	return newPath(fmt.Sprintf("/projects/columns/%d/cards", columnID)).
		query("per_page", "100").
		String()
}

// ListProjectCards lists the project cards.
func (c *client) ListProjectCards(columnID int) ProjectCards {
	ps := make(chan interface{})
	go func() {
		defer close(ps)
		path := c.url(listProjectCardsPath(columnID))
		for {
			var xs []*ProjectCard
			next, err := c.getList(path, &xs)
			if err != nil {
				ps <- fmt.Errorf("ListProjectCards %d: %w", columnID, err)
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
	return ProjectCards(ps)
}

func getProjectCardPath(projectCardID int) string {
	return newPath(fmt.Sprintf("/projects/columns/cards/%d", projectCardID)).
		String()
}

type projectCardOrError struct {
	ProjectCard
	Message string `json:"message"`
}

func (c *client) GetProjectCard(projectCardID int) (*ProjectCard, error) {
	res, err := c.get(c.url(getProjectCardPath(projectCardID)))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var r projectCardOrError
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, err
	}

	if r.Message != "" {
		return nil, fmt.Errorf("GetProjectCard %s: %s", fmt.Sprintf("projects/columns/cards/%d", projectCardID), r.Message)
	}

	return &r.ProjectCard, nil
}

// CreateProjectCardParams represents the paramter for CreateProjectCard API.
type CreateProjectCardParams struct {
	Note        string                 `json:"note"`
	ContentID   int                    `json:"content_id"`
	ContentType ProjectCardContentType `json:"content_type"`
}

// ProjectCardContentType ...
type ProjectCardContentType int

// ProjectCardContentType ...
const (
	ProjectCardContentTypeIssue ProjectCardContentType = iota + 1
	ProjectCardContentTypePullRequest
)

var stringToProjectCardContentType = map[string]ProjectCardContentType{
	"Issue":       ProjectCardContentTypeIssue,
	"PullRequest": ProjectCardContentTypePullRequest,
}

var projectCardContentTypeToString = map[ProjectCardContentType]string{
	ProjectCardContentTypeIssue:       "Issue",
	ProjectCardContentTypePullRequest: "PullRequest",
}

// UnmarshalJSON implements json.Unmarshaler
func (t *ProjectCardContentType) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	if x, ok := stringToProjectCardContentType[s]; ok {
		*t = x
		return nil
	}
	return fmt.Errorf("unknown project card content type: %q", s)
}

// MarshalJSON implements json.Marshaler
func (t ProjectCardContentType) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

// String implements Stringer
func (t ProjectCardContentType) String() string {
	return projectCardContentTypeToString[t]
}

// GoString implements GoString
func (t ProjectCardContentType) GoString() string {
	return strconv.Quote(t.String())
}

func createProjectCardPath(columnID int) string {
	return newPath(fmt.Sprintf("/projects/columns/%d/cards", columnID)).
		String()
}

// CreateProjectCard creates a project card.
func (c *client) CreateProjectCard(columnID int, params *CreateProjectCardParams) (*ProjectCard, error) {
	bs, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	body := bytes.NewReader(bs)
	res, err := c.post(c.url(createProjectCardPath(columnID)), body)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var r projectCardOrError
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, err
	}

	if r.Message != "" {
		return nil, fmt.Errorf("CreateProjectCard %s: %s", fmt.Sprintf("projects/columns/%d/cards", columnID), r.Message)
	}

	return &r.ProjectCard, nil
}

// UpdateProjectCardParams represents the paramter for UpdateProjectCard API.
type UpdateProjectCardParams struct {
	Note     string `json:"note,omitempty"`
	Archived bool   `json:"archived,omitempty"`
}

func updateProjectCardPath(projectCardID int) string {
	return newPath(fmt.Sprintf("/projects/columns/cards/%d", projectCardID)).
		String()
}

// UpdateProjectCard updates the project card.
func (c *client) UpdateProjectCard(projectCardID int, params *UpdateProjectCardParams) (*ProjectCard, error) {
	bs, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	body := bytes.NewReader(bs)
	res, err := c.patch(c.url(updateProjectCardPath(projectCardID)), body)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var r projectCardOrError
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, err
	}

	if r.Message != "" {
		return nil, fmt.Errorf("UpdateProjectCard %s: %s", fmt.Sprintf("projects/columns/cards/%d", projectCardID), r.Message)
	}

	return &r.ProjectCard, nil
}

// MoveProjectCardParams represents the paramter for MoveProjectCard API.
type MoveProjectCardParams struct {
	Position string `json:"position"`
	ColumnID bool   `json:"column_id,omitempty"`
}

func moveProjectCardPath(projectCardID int) string {
	return newPath(fmt.Sprintf("/projects/columns/cards/%d/moves", projectCardID)).
		String()
}

// MoveProjectCard moves the project card.
func (c *client) MoveProjectCard(projectCardID int, params *MoveProjectCardParams) (*ProjectCard, error) {
	bs, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	body := bytes.NewReader(bs)
	res, err := c.post(c.url(moveProjectCardPath(projectCardID)), body)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var r projectCardOrError
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, err
	}

	if r.Message != "" {
		return nil, fmt.Errorf("MoveProjectCard %s: %s", fmt.Sprintf("projects/columns/cards/%d/moves", projectCardID), r.Message)
	}

	return &r.ProjectCard, nil
}
