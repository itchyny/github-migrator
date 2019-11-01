package github

import (
	"encoding/json"
	"fmt"
	"io"
)

// ReviewComment represents a review comment.
type ReviewComment struct {
	ID          int    `json:"id"`
	Path        string `json:"path"`
	Line        int    `json:"line"`
	Body        string `json:"body"`
	DiffHunk    string `json:"diff_hunk"`
	HTMLURL     string `json:"html_url"`
	User        *User  `json:"user"`
	InReplyToID int    `json:"in_reply_to_id"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// ReviewComments represents a collection of review comments.
type ReviewComments <-chan interface{}

// Next emits the next ReviewComment.
func (cs ReviewComments) Next() (*ReviewComment, error) {
	for x := range cs {
		switch x := x.(type) {
		case error:
			return nil, x
		case *ReviewComment:
			return x, nil
		}
		break
	}
	return nil, io.EOF
}

// ReviewCommentsFromSlice creates ReviewComments from a slice.
func ReviewCommentsFromSlice(xs []*ReviewComment) ReviewComments {
	cs := make(chan interface{})
	go func() {
		defer close(cs)
		for _, i := range xs {
			cs <- i
		}
	}()
	return cs
}

// ReviewCommentsToSlice collects ReviewComments.
func ReviewCommentsToSlice(cs ReviewComments) ([]*ReviewComment, error) {
	xs := []*ReviewComment{}
	for {
		i, err := cs.Next()
		if err != nil {
			if err != io.EOF {
				return nil, err
			}
			return xs, nil
		}
		xs = append(xs, i)
	}
}

func listReviewCommentsPath(repo string, pullNumber int) string {
	return newPath(fmt.Sprintf("/repos/%s/pulls/%d/comments", repo, pullNumber)).
		String()
}

// ListReviewComments lists the review comments of a pull request.
func (c *client) ListReviewComments(repo string, pullNumber int) ReviewComments {
	cs := make(chan interface{})
	go func() {
		defer close(cs)
		path := c.url(listReviewCommentsPath(repo, pullNumber))
		for {
			xs, next, err := c.listReviewComments(path)
			if err != nil {
				cs <- err
				break
			}
			for _, x := range xs {
				cs <- x
			}
			if next == "" {
				break
			}
			path = next
		}
	}()
	return ReviewComments(cs)
}

func (c *client) listReviewComments(path string) ([]*ReviewComment, string, error) {
	res, err := c.get(path)
	if err != nil {
		return nil, "", err
	}
	defer res.Body.Close()

	var r []*ReviewComment
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, "", err
	}

	return r, getNext(res.Header), nil
}
