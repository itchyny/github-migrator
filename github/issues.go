package github

import "encoding/json"

// Issue represents an issue.
type Issue struct {
	Number  int    `json:"number"`
	Title   string `json:"title"`
	State   string `json:"state"`
	Body    string `json:"body"`
	HTMLURL string `json:"html_url"`
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
		String()
}

// ListIssues lists the issues.
func (c *client) ListIssues(repo string, params *ListIssuesParams) ([]*Issue, error) {
	res, err := c.get(listIssuesPath(repo, params))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var r []*Issue
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, err
	}

	return r, nil
}
