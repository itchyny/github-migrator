package github

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"
)

// Issue represents an issue.
type Issue struct {
	Number      int               `json:"number"`
	Title       string            `json:"title"`
	State       IssueState        `json:"state"`
	Body        string            `json:"body"`
	HTMLURL     string            `json:"html_url"`
	User        *User             `json:"user"`
	Assignee    *User             `json:"assignee"`
	CreatedAt   string            `json:"created_at"`
	UpdatedAt   string            `json:"updated_at"`
	ClosedAt    string            `json:"closed_at,omitempty"`
	ClosedBy    *User             `json:"closed_by,omitempty"`
	Labels      []*Label          `json:"labels"`
	PullRequest *IssuePullRequest `json:"pull_request"`
}

// IssueState ...
type IssueState int

// IssueState ...
const (
	IssueStateOpen IssueState = iota + 1
	IssueStateClosed
)

var stringToIssueState = map[string]IssueState{
	"open":   IssueStateOpen,
	"closed": IssueStateClosed,
}

var issueStateToString = map[IssueState]string{
	IssueStateOpen:   "open",
	IssueStateClosed: "closed",
}

// UnmarshalJSON implements json.Unmarshaler
func (t *IssueState) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	if x, ok := stringToIssueState[s]; ok {
		*t = x
		return nil
	}
	return fmt.Errorf("unknown issue state: %q", s)
}

// MarshalJSON implements json.Marshaler
func (t IssueState) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

// String implements Stringer
func (t IssueState) String() string {
	return issueStateToString[t]
}

// GoString implements GoString
func (t IssueState) GoString() string {
	return strconv.Quote(t.String())
}

// IssueType ...
type IssueType int

// IssueType ...
const (
	IssueTypeIssue IssueType = iota + 1
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
	ListIssuesParamFilterDefault ListIssuesParamFilter = iota + 1
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
	ListIssuesParamStateDefault ListIssuesParamState = iota + 1
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
	ListIssuesParamSortDefault ListIssuesParamSort = iota + 1
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
	ListIssuesParamDirectionDefault ListIssuesParamDirection = iota + 1
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
	return newPath(fmt.Sprintf("/repos/%s/issues", repo)).
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
			var xs []*Issue
			next, err := c.getList(path, &xs)
			if err != nil {
				is <- fmt.Errorf("ListIssues %s: %w", repo, err)
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

func (c *client) GetIssue(repo string, issueNumber int) (*Issue, error) {
	var r Issue
	if err := c.get(c.url(fmt.Sprintf("/repos/%s/issues/%d", repo, issueNumber)), &r); err != nil {
		return nil, fmt.Errorf("GetIssue %s: %w", fmt.Sprintf("%s/issues/%d", repo, issueNumber), err)
	}
	return &r, nil
}
