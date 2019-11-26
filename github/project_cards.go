package github

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"
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

// GetIssueNumber ...
func (c *ProjectCard) GetIssueNumber() int {
	if i := strings.LastIndexByte(c.ContentURL, '/'); i >= 0 {
		j, err := strconv.Atoi(c.ContentURL[i+1:])
		if err != nil {
			return -1
		}
		return j
	}
	return -1
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

// ListProjectCards lists the project cards.
func (c *client) ListProjectCards(columnID int) ProjectCards {
	ps := make(chan interface{})
	go func() {
		defer close(ps)
		path := c.url(fmt.Sprintf("/projects/columns/%d/cards?per_page=100", columnID))
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

func (c *client) GetProjectCard(projectCardID int) (*ProjectCard, error) {
	var r ProjectCard
	if err := c.get(c.url(fmt.Sprintf("/projects/columns/cards/%d", projectCardID)), &r); err != nil {
		return nil, fmt.Errorf("GetProjectCard %s: %w", fmt.Sprintf("projects/columns/cards/%d", projectCardID), err)
	}
	return &r, nil
}

// CreateProjectCardParams represents the paramter for CreateProjectCard API.
type CreateProjectCardParams struct {
	Note        string                 `json:"note,omitempty"`
	ContentID   int                    `json:"content_id,omitempty"`
	ContentType ProjectCardContentType `json:"content_type,omitempty"`
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

// CreateProjectCard creates a project card.
func (c *client) CreateProjectCard(columnID int, params *CreateProjectCardParams) (*ProjectCard, error) {
	var r ProjectCard
	if err := c.post(c.url(fmt.Sprintf("/projects/columns/%d/cards", columnID)), params, &r); err != nil {
		return nil, fmt.Errorf("CreateProjectCard %s: %w", fmt.Sprintf("projects/columns/%d/cards", columnID), err)
	}
	return &r, nil
}

// UpdateProjectCardParams represents the paramter for UpdateProjectCard API.
type UpdateProjectCardParams struct {
	Note     string `json:"note,omitempty"`
	Archived bool   `json:"archived,omitempty"`
}

// UpdateProjectCard updates the project card.
func (c *client) UpdateProjectCard(projectCardID int, params *UpdateProjectCardParams) (*ProjectCard, error) {
	var r ProjectCard
	if err := c.patch(c.url(fmt.Sprintf("/projects/columns/cards/%d", projectCardID)), params, &r); err != nil {
		return nil, fmt.Errorf("UpdateProjectCard %s: %w", fmt.Sprintf("projects/columns/cards/%d", projectCardID), err)
	}
	return &r, nil
}

// MoveProjectCardParams represents the paramter for MoveProjectCard API.
type MoveProjectCardParams struct {
	Position string `json:"position"`
	ColumnID bool   `json:"column_id,omitempty"`
}

// MoveProjectCard moves the project card.
func (c *client) MoveProjectCard(projectCardID int, params *MoveProjectCardParams) (*ProjectCard, error) {
	var r ProjectCard
	if err := c.post(c.url(fmt.Sprintf("/projects/columns/cards/%d/moves", projectCardID)), params, &r); err != nil {
		return nil, fmt.Errorf("MoveProjectCard %s: %w", fmt.Sprintf("projects/columns/cards/%d/moves", projectCardID), err)
	}
	return &r, nil
}
