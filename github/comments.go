package github

import (
	"encoding/json"
	"fmt"
	"io"
)

// Comment represents a comment.
type Comment struct {
	Body      string `json:"body"`
	HTMLURL   string `json:"html_url"`
	User      *User  `json:"user"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// Comments represents a collection of comments.
type Comments <-chan interface{}

// Next emits the next Comment.
func (cs Comments) Next() (*Comment, error) {
	for x := range cs {
		switch x := x.(type) {
		case error:
			return nil, x
		case *Comment:
			return x, nil
		}
		break
	}
	return nil, io.EOF
}

// CommentsFromSlice creates Comments from a slice.
func CommentsFromSlice(xs []*Comment) Comments {
	cs := make(chan interface{})
	go func() {
		defer close(cs)
		for _, i := range xs {
			cs <- i
		}
	}()
	return cs
}

// CommentsToSlice collects Comments.
func CommentsToSlice(cs Comments) ([]*Comment, error) {
	var xs []*Comment
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

func listCommentsPath(repo string, issueNumber int) string {
	return newPath(fmt.Sprintf("/repos/%s/issues/%d/comments", repo, issueNumber)).
		String()
}

// ListComments lists the comments of an issue.
func (c *client) ListComments(repo string, issueNumber int) Comments {
	cs := make(chan interface{})
	go func() {
		defer close(cs)
		path := c.url(listCommentsPath(repo, issueNumber))
		for {
			xs, next, err := c.listComments(path)
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
	return Comments(cs)
}

func (c *client) listComments(path string) ([]*Comment, string, error) {
	res, err := c.get(path)
	if err != nil {
		return nil, "", err
	}
	defer res.Body.Close()

	var r []*Comment
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, "", err
	}

	return r, getNext(res.Header), nil
}
