package github

import (
	"encoding/json"
	"fmt"
	"io"
)

// PullReq represents a pull request.
type PullReq struct {
	Issue
	Merged   bool  `json:"merged"`
	MergedBy *User `json:"merged_by"`
	Draft    bool  `json:"draft"`
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
	ListPullReqsParamStateDefault ListPullReqsParamState = iota
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
	ListPullReqsParamSortDefault ListPullReqsParamSort = iota
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
	ListPullReqsParamDirectionDefault ListPullReqsParamDirection = iota
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
	return newPath("/repos/"+repo+"/pulls").
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
			xs, next, err := c.listPullReqs(path)
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
	return PullReqs(ps)
}

func (c *client) listPullReqs(path string) ([]*PullReq, string, error) {
	res, err := c.get(path)
	if err != nil {
		return nil, "", err
	}
	defer res.Body.Close()

	var r []*PullReq
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, "", err
	}

	return r, getNext(res.Header), nil
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
		return nil, fmt.Errorf("%s: %s", r.Message, "/pulls/"+fmt.Sprint(pullNumber))
	}

	return &r.PullReq, nil
}
