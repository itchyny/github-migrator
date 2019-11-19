package github

import (
	"encoding/json"
	"fmt"
	"io"
)

// PullReq represents a pull request.
type PullReq struct {
	Issue
	Merged         bool        `json:"merged"`
	MergedAt       string      `json:"merged_at"`
	MergedBy       *User       `json:"merged_by"`
	MergeCommitSHA string      `json:"merge_commit_sha"`
	Draft          bool        `json:"draft"`
	Head           *PullReqRef `json:"head"`
	Base           *PullReqRef `json:"base"`
	Commits        int         `json:"commits"`
	Additions      int         `json:"additions"`
	Deletions      int         `json:"deletions"`
	ChangedFiles   int         `json:"changed_files"`
}

// PullReqRef ...
type PullReqRef struct {
	SHA  string `json:"sha"`
	Ref  string `json:"ref"`
	User *User  `json:"user"`
	Repo *Repo  `json:"repo"`
}

// PullReqs represents a collection of pull requests.
type PullReqs <-chan interface{}

// Next emits the next PullReq.
func (ps PullReqs) Next() (*PullReq, error) {
	for x := range ps {
		switch x := x.(type) {
		case error:
			return nil, x
		case *PullReq:
			return x, nil
		}
		break
	}
	return nil, io.EOF
}

// PullReqsFromSlice creates PullReqs from a slice.
func PullReqsFromSlice(xs []*PullReq) PullReqs {
	ps := make(chan interface{})
	go func() {
		defer close(ps)
		for _, p := range xs {
			ps <- p
		}
	}()
	return ps
}

// PullReqsToSlice collects PullReqs.
func PullReqsToSlice(ps PullReqs) ([]*PullReq, error) {
	xs := []*PullReq{}
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

// ListPullReqsParams represents the paramter for ListPullReqs API.
type ListPullReqsParams struct {
	State     ListPullReqsParamState
	Head      string
	Base      string
	Sort      ListPullReqsParamSort
	Direction ListPullReqsParamDirection
}

// ListPullReqsParamState ...
type ListPullReqsParamState int

// ListPullReqsParamState ...
const (
	ListPullReqsParamStateDefault ListPullReqsParamState = iota + 1
	ListPullReqsParamStateOpen
	ListPullReqsParamStateClosed
	ListPullReqsParamStateAll
)

func (f ListPullReqsParamState) String() string {
	switch f {
	case ListPullReqsParamStateOpen:
		return "open"
	case ListPullReqsParamStateClosed:
		return "closed"
	case ListPullReqsParamStateAll:
		return "all"
	default:
		return ""
	}
}

// ListPullReqsParamSort ...
type ListPullReqsParamSort int

// ListPullReqsParamSort ...
const (
	ListPullReqsParamSortDefault ListPullReqsParamSort = iota + 1
	ListPullReqsParamSortCreated
	ListPullReqsParamSortUpdated
	ListPullReqsParamSortPopularity
	ListPullReqsParamSortLongRunning
)

func (f ListPullReqsParamSort) String() string {
	switch f {
	case ListPullReqsParamSortCreated:
		return "created"
	case ListPullReqsParamSortUpdated:
		return "updated"
	case ListPullReqsParamSortPopularity:
		return "popularity"
	case ListPullReqsParamSortLongRunning:
		return "long-running"
	default:
		return ""
	}
}

// ListPullReqsParamDirection ...
type ListPullReqsParamDirection int

// ListPullReqsParamDirection ...
const (
	ListPullReqsParamDirectionDefault ListPullReqsParamDirection = iota + 1
	ListPullReqsParamDirectionAsc
	ListPullReqsParamDirectionDesc
)

func (f ListPullReqsParamDirection) String() string {
	switch f {
	case ListPullReqsParamDirectionAsc:
		return "asc"
	case ListPullReqsParamDirectionDesc:
		return "desc"
	default:
		return ""
	}
}

func listPullReqsPath(repo string, params *ListPullReqsParams) string {
	return newPath(fmt.Sprintf("/repos/%s/pulls", repo)).
		query("state", params.State.String()).
		query("head", params.Head).
		query("base", params.Base).
		query("sort", params.Sort.String()).
		query("direction", params.Direction.String()).
		query("per_page", "100").
		String()
}

// ListPullReqs lists the pull requests.
func (c *client) ListPullReqs(repo string, params *ListPullReqsParams) PullReqs {
	ps := make(chan interface{})
	go func() {
		defer close(ps)
		path := c.url(listPullReqsPath(repo, params))
		for {
			var xs []*PullReq
			next, err := c.getList(path, &xs)
			if err != nil {
				ps <- fmt.Errorf("ListPullReqs %s: %w", repo, err)
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
	return PullReqs(ps)
}

func getPullReqPath(repo string, pullNumber int) string {
	return newPath(fmt.Sprintf("/repos/%s/pulls/%d", repo, pullNumber)).
		String()
}

type pullReqOrError struct {
	PullReq
	Message string `json:"message"`
}

func (c *client) GetPullReq(repo string, pullNumber int) (*PullReq, error) {
	res, err := c.get(c.url(getPullReqPath(repo, pullNumber)))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var r pullReqOrError
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, err
	}

	if r.Message != "" {
		return nil, fmt.Errorf("GetPullReq %s: %s", fmt.Sprintf("%s/pulls/%d", repo, pullNumber), r.Message)
	}

	return &r.PullReq, nil
}
