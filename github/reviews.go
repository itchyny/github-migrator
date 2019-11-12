package github

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"
)

// Review represents a review.
type Review struct {
	ID          int         `json:"id"`
	State       ReviewState `json:"state"`
	Body        string      `json:"body"`
	HTMLURL     string      `json:"html_url"`
	User        *User       `json:"user"`
	CommitID    string      `json:"commit_id"`
	SubmittedAt string      `json:"submitted_at"`
}

// ReviewState ...
type ReviewState int

// ReviewState ...
const (
	ReviewStateApproved ReviewState = iota + 1
	ReviewStateChangesRequested
	ReviewStateCommented
	ReviewStatePending
)

var stringToReviewState = map[string]ReviewState{
	"APPROVED":          ReviewStateApproved,
	"CHANGES_REQUESTED": ReviewStateChangesRequested,
	"COMMENTED":         ReviewStateCommented,
	"PENDING":           ReviewStatePending,
}

var reviewStateToString = map[ReviewState]string{
	ReviewStateApproved:         "APPROVED",
	ReviewStateChangesRequested: "CHANGES_REQUESTED",
	ReviewStateCommented:        "COMMENTED",
	ReviewStatePending:          "PENDING",
}

// UnmarshalJSON implements json.Unmarshaler
func (t *ReviewState) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	if x, ok := stringToReviewState[s]; ok {
		*t = x
		return nil
	}
	return fmt.Errorf("unknown review state: %q", s)
}

// MarshalJSON implements json.Marshaler
func (t ReviewState) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.String())
}

// String implements Stringer
func (t ReviewState) String() string {
	return reviewStateToString[t]
}

// GoString implements GoString
func (t ReviewState) GoString() string {
	return strconv.Quote(t.String())
}

// Reviews represents a collection of reviews.
type Reviews <-chan interface{}

// Next emits the next Review.
func (rs Reviews) Next() (*Review, error) {
	for x := range rs {
		switch x := x.(type) {
		case error:
			return nil, x
		case *Review:
			return x, nil
		}
		break
	}
	return nil, io.EOF
}

// ReviewsFromSlice creates Reviews from a slice.
func ReviewsFromSlice(xs []*Review) Reviews {
	rs := make(chan interface{})
	go func() {
		defer close(rs)
		for _, p := range xs {
			rs <- p
		}
	}()
	return rs
}

// ReviewsToSlice collects Reviews.
func ReviewsToSlice(rs Reviews) ([]*Review, error) {
	xs := []*Review{}
	for {
		p, err := rs.Next()
		if err != nil {
			if err != io.EOF {
				return nil, err
			}
			return xs, nil
		}
		xs = append(xs, p)
	}
}

func listReviewsPath(repo string, pullNumber int) string {
	return newPath(fmt.Sprintf("/repos/%s/pulls/%d/reviews", repo, pullNumber)).
		String()
}

// ListReviews lists the reviews.
func (c *client) ListReviews(repo string, pullNumber int) Reviews {
	rs := make(chan interface{})
	go func() {
		defer close(rs)
		path := c.url(listReviewsPath(repo, pullNumber))
		for {
			xs, next, err := c.listReviews(path)
			if err != nil {
				rs <- err
				break
			}
			for _, x := range xs {
				rs <- x
			}
			if next == "" {
				break
			}
			path = next
		}
	}()
	return Reviews(rs)
}

func (c *client) listReviews(path string) ([]*Review, string, error) {
	res, err := c.get(path)
	if err != nil {
		return nil, "", err
	}
	defer res.Body.Close()

	var r []*Review
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, "", err
	}

	return r, getNext(res.Header), nil
}
