package github

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"
)

// Milestone represents a milestone.
type Milestone struct {
	ID               int            `json:"id"`
	HTMLURL          string         `json:"html_url"`
	Number           int            `json:"number"`
	Title            string         `json:"title"`
	Description      string         `json:"description"`
	State            MilestoneState `json:"state"`
	OpenMilestones   int            `json:"open_milestones"`
	ClosedMilestones int            `json:"closed_milestones"`
	Creator          *User          `json:"creator"`
	CreatedAt        string         `json:"created_at"`
	UpdatedAt        string         `json:"updated_at"`
	ClosedAt         string         `json:"closed_at"`
	DueOn            string         `json:"due_on"`
}

// MilestoneState ...
type MilestoneState int

// MilestoneState ...
const (
	MilestoneStateOpen MilestoneState = iota + 1
	MilestoneStateClosed
)

var stringToMilestoneState = map[string]MilestoneState{
	"open":   MilestoneStateOpen,
	"closed": MilestoneStateClosed,
}

var milestoneStateToString = map[MilestoneState]string{
	MilestoneStateOpen:   "open",
	MilestoneStateClosed: "closed",
}

// UnmarshalJSON implements json.Unmarshaler
func (t *MilestoneState) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	if x, ok := stringToMilestoneState[s]; ok {
		*t = x
		return nil
	}
	return fmt.Errorf("unknown milestone state: %q", s)
}

// MarshalJSON implements json.Marshaler
func (t MilestoneState) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

// String implements Stringer
func (t MilestoneState) String() string {
	return milestoneStateToString[t]
}

// GoString implements GoString
func (t MilestoneState) GoString() string {
	return strconv.Quote(t.String())
}

// Milestones represents a collection of milestones.
type Milestones <-chan interface{}

// Next emits the next Milestone.
func (ms Milestones) Next() (*Milestone, error) {
	for x := range ms {
		switch x := x.(type) {
		case error:
			return nil, x
		case *Milestone:
			return x, nil
		}
		break
	}
	return nil, io.EOF
}

// MilestonesFromSlice creates Milestones from a slice.
func MilestonesFromSlice(xs []*Milestone) Milestones {
	ms := make(chan interface{})
	go func() {
		defer close(ms)
		for _, p := range xs {
			ms <- p
		}
	}()
	return ms
}

// MilestonesToSlice collects Milestones.
func MilestonesToSlice(ms Milestones) ([]*Milestone, error) {
	xs := []*Milestone{}
	for {
		p, err := ms.Next()
		if err != nil {
			if err != io.EOF {
				return nil, err
			}
			return xs, nil
		}
		xs = append(xs, p)
	}
}

// ListMilestonesParams represents the paramter for ListMilestones API.
type ListMilestonesParams struct {
	State     ListMilestonesParamState
	Direction ListMilestonesParamDirection
	Sort      ListMilestonesParamSort
}

// ListMilestonesParamState ...
type ListMilestonesParamState int

// ListMilestonesParamState ...
const (
	ListMilestonesParamStateDefault ListMilestonesParamState = iota + 1
	ListMilestonesParamStateOpen
	ListMilestonesParamStateClosed
	ListMilestonesParamStateAll
)

func (f ListMilestonesParamState) String() string {
	switch f {
	case ListMilestonesParamStateOpen:
		return "open"
	case ListMilestonesParamStateClosed:
		return "closed"
	case ListMilestonesParamStateAll:
		return "all"
	default:
		return ""
	}
}

// ListMilestonesParamSort ...
type ListMilestonesParamSort int

// ListMilestonesParamSort ...
const (
	ListMilestonesParamSortDefault ListMilestonesParamSort = iota + 1
	ListMilestonesParamSortDueOn
	ListMilestonesParamSortCompleteness
)

func (f ListMilestonesParamSort) String() string {
	switch f {
	case ListMilestonesParamSortDueOn:
		return "due_on"
	case ListMilestonesParamSortCompleteness:
		return "completeness"
	default:
		return ""
	}
}

// ListMilestonesParamDirection ...
type ListMilestonesParamDirection int

// ListMilestonesParamDirection ...
const (
	ListMilestonesParamDirectionDefault ListMilestonesParamDirection = iota + 1
	ListMilestonesParamDirectionAsc
	ListMilestonesParamDirectionDesc
)

func (f ListMilestonesParamDirection) String() string {
	switch f {
	case ListMilestonesParamDirectionAsc:
		return "asc"
	case ListMilestonesParamDirectionDesc:
		return "desc"
	default:
		return ""
	}
}

func listMilestonesPath(repo string, params *ListMilestonesParams) string {
	return newPath(fmt.Sprintf("/repos/%s/milestones", repo)).
		query("state", params.State.String()).
		query("sort", params.Sort.String()).
		query("direction", params.Direction.String()).
		query("per_page", "100").
		String()
}

// ListMilestones lists the milestones.
func (c *client) ListMilestones(repo string, params *ListMilestonesParams) Milestones {
	ms := make(chan interface{})
	go func() {
		defer close(ms)
		path := c.url(listMilestonesPath(repo, params))
		for {
			var xs []*Milestone
			next, err := c.getList(path, &xs)
			if err != nil {
				ms <- fmt.Errorf("ListMilestones %s: %w", repo, err)
				break
			}
			for _, x := range xs {
				ms <- x
			}
			if next == "" {
				break
			}
			path = next
		}
	}()
	return Milestones(ms)
}

func (c *client) GetMilestone(repo string, milestoneNumber int) (*Milestone, error) {
	var r Milestone
	if err := c.get(c.url(fmt.Sprintf("/repos/%s/milestones/%d", repo, milestoneNumber)), &r); err != nil {
		return nil, fmt.Errorf("GetMilestone %s: %w", fmt.Sprintf("%s/milestones/%d", repo, milestoneNumber), err)
	}
	return &r, nil
}

// CreateMilestoneParams represents the paramter for CreateMilestone API.
type CreateMilestoneParams struct {
	Title       string         `json:"title"`
	Description string         `json:"description"`
	State       MilestoneState `json:"state"`
	DueOn       string         `json:"due_on"`
}

// CreateMilestone creates a milestone.
func (c *client) CreateMilestone(repo string, params *CreateMilestoneParams) (*Milestone, error) {
	var r Milestone
	if err := c.post(c.url(fmt.Sprintf("/repos/%s/milestones", repo)), params, &r); err != nil {
		return nil, fmt.Errorf("CreateMilestone %s: %w", fmt.Sprintf("%s/milestones", repo), err)
	}
	return &r, nil
}

// UpdateMilestoneParams represents the paramter for UpdateMilestone API.
type UpdateMilestoneParams CreateMilestoneParams

// UpdateMilestone updates the milestone.
func (c *client) UpdateMilestone(repo string, milestoneNumber int, params *UpdateMilestoneParams) (*Milestone, error) {
	var r Milestone
	if err := c.patch(c.url(fmt.Sprintf("/repos/%s/milestones/%d", repo, milestoneNumber)), params, &r); err != nil {
		return nil, fmt.Errorf("UpdateMilestone %s: %w", fmt.Sprintf("%s/milestones/%d", repo, milestoneNumber), err)
	}
	return &r, nil
}

// DeleteMilestone deletes the milestone.
func (c *client) DeleteMilestone(repo string, milestoneNumber int) error {
	if err := c.delete(c.url(fmt.Sprintf("/repos/%s/milestones/%d", repo, milestoneNumber))); err != nil {
		return fmt.Errorf("DeleteMilestone %s: %w", fmt.Sprintf("%s/milestones/%d", repo, milestoneNumber), err)
	}
	return nil
}
