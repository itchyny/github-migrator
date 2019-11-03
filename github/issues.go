package github

import (
	"encoding/json"
	"io"
)

// Issue represents an issue.
type Issue struct {
	Number      int               `json:"number"`
	Title       string            `json:"title"`
	State       string            `json:"state"`
	Body        string            `json:"body"`
	HTMLURL     string            `json:"html_url"`
	User        *User             `json:"user"`
	Assignee    *User             `json:"assignee"`
	CreatedAt   string            `json:"created_at"`
	UpdatedAt   string            `json:"updated_at"`
	ClosedAt    string            `json:"closed_at,omitempty"`
	Labels      []*Label          `json:"labels"`
	PullRequest *IssuePullRequest `json:"pull_request"`
}

// IssueType ...
type IssueType int

// IssueType ...
const (
	IssueTypeIssue IssueType = iota
	IssueTypePullReq
)

func (t IssueType) String() string {
	switch t {
	case IssueTypeIssue:
		return "issue"
	case IssueTypePullReq:
		return "pull request"
	default:
		return ""
	}
}

// Type returns IssueTypePullReq ("pull request") or IssueTypeIssue ("issue").
func (i *Issue) Type() IssueType {
	if i.PullRequest != nil {
		return IssueTypePullReq
	}
	return IssueTypeIssue
}

// IssuePullRequest represents the pull request information of an issue.
type IssuePullRequest struct {
	URL      string `json:"url"`
	HTMLURL  string `json:"html_url"`
	DiffURL  string `json:"diff_url"`
	PatchURL string `json:"patch_url"`
}

// Issues represents a collection of issues.
type Issues <-chan interface{}

// Next emits the next Issue.
func (is Issues) Next() (*Issue, error) {
	for x := range is {
		switch x := x.(type) {
		case error:
			return nil, x
		case *Issue:
			return x, nil
		}
		break
	}
	return nil, io.EOF
}

// IssuesFromSlice creates Issues from a slice.
func IssuesFromSlice(xs []*Issue) Issues {
	is := make(chan interface{})
	go func() {
		defer close(is)
		for _, i := range xs {
			is <- i
		}
	}()
	return is
}

// IssuesToSlice collects Issues.
func IssuesToSlice(is Issues) ([]*Issue, error) {
	xs := []*Issue{}
	for {
		i, err := is.Next()
		if err != nil {
			if err != io.EOF {
				return nil, err
			}
			return xs, nil
		}
		xs = append(xs, i)
	}
}

// ListIssuesParams represents the paramter for ListIssues API.
type ListIssuesParams struct {
	Filter    ListIssuesParamFilter
	State     ListIssuesParamState
	Sort      ListIssuesParamSort
	Direction ListIssuesParamDirection
}

// ListIssuesParamFilter ...
type ListIssuesParamFilter int

// ListIssuesParamFilter ...
const (
	ListIssuesParamFilterDefault ListIssuesParamFilter = iota
	ListIssuesParamFilterAssigned
	ListIssuesParamFilterCreated
	ListIssuesParamFilterMentioned
	ListIssuesParamFilterSubscribed
	ListIssuesParamFilterAll
)

func (f ListIssuesParamFilter) String() string {
	switch f {
	case ListIssuesParamFilterAssigned:
		return "assigned"
	case ListIssuesParamFilterCreated:
		return "created"
	case ListIssuesParamFilterMentioned:
		return "mentioned"
	case ListIssuesParamFilterSubscribed:
		return "subscribed"
	case ListIssuesParamFilterAll:
		return "all"
	default:
		return ""
	}
}

// ListIssuesParamState ...
type ListIssuesParamState int

// ListIssuesParamState ...
const (
	ListIssuesParamStateDefault ListIssuesParamState = iota
	ListIssuesParamStateOpen
	ListIssuesParamStateClosed
	ListIssuesParamStateAll
)

func (f ListIssuesParamState) String() string {
	switch f {
	case ListIssuesParamStateOpen:
		return "open"
	case ListIssuesParamStateClosed:
		return "closed"
	case ListIssuesParamStateAll:
		return "all"
	default:
		return ""
	}
}

// ListIssuesParamSort ...
type ListIssuesParamSort int

// ListIssuesParamSort ...
const (
	ListIssuesParamSortDefault ListIssuesParamSort = iota
	ListIssuesParamSortCreated
	ListIssuesParamSortUpdated
	ListIssuesParamSortComments
)

func (f ListIssuesParamSort) String() string {
	switch f {
	case ListIssuesParamSortCreated:
		return "created"
	case ListIssuesParamSortUpdated:
		return "updated"
	case ListIssuesParamSortComments:
		return "comments"
	default:
		return ""
	}
}

// ListIssuesParamDirection ...
type ListIssuesParamDirection int

// ListIssuesParamDirection ...
const (
	ListIssuesParamDirectionDefault ListIssuesParamDirection = iota
	ListIssuesParamDirectionAsc
	ListIssuesParamDirectionDesc
)

func (f ListIssuesParamDirection) String() string {
	switch f {
	case ListIssuesParamDirectionAsc:
		return "asc"
	case ListIssuesParamDirectionDesc:
		return "desc"
	default:
		return ""
	}
}

func listIssuesPath(repo string, params *ListIssuesParams) string {
	return newPath("/repos/"+repo+"/issues").
		query("filter", params.Filter.String()).
		query("state", params.State.String()).
		query("sort", params.Sort.String()).
		query("direction", params.Direction.String()).
		query("per_page", "100").
		String()
}

// ListIssues lists the issues.
func (c *client) ListIssues(repo string, params *ListIssuesParams) Issues {
	is := make(chan interface{})
	go func() {
		defer close(is)
		path := c.url(listIssuesPath(repo, params))
		for {
			xs, next, err := c.listIssues(path)
			if err != nil {
				is <- err
				break
			}
			for _, x := range xs {
				is <- x
			}
			if next == "" {
				break
			}
			path = next
		}
	}()
	return Issues(is)
}

func (c *client) listIssues(path string) ([]*Issue, string, error) {
	res, err := c.get(path)
	if err != nil {
		return nil, "", err
	}
	defer res.Body.Close()

	var r []*Issue
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, "", err
	}

	return r, getNext(res.Header), nil
}
