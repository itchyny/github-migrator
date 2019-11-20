package github

import (
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
		for _, c := range xs {
			cs <- c
		}
	}()
	return cs
}

// CommentsToSlice collects Comments.
func CommentsToSlice(cs Comments) ([]*Comment, error) {
	xs := []*Comment{}
	for {
		c, err := cs.Next()
		if err != nil {
			if err != io.EOF {
				return nil, err
			}
			return xs, nil
		}
		xs = append(xs, c)
	}
}

// ListComments lists the comments of an issue.
func (c *client) ListComments(repo string, issueNumber int) Comments {
	cs := make(chan interface{})
	go func() {
		defer close(cs)
		path := c.url(fmt.Sprintf("/repos/%s/issues/%d/comments?per_page=100", repo, issueNumber))
		for {
			var xs []*Comment
			next, err := c.getList(path, &xs)
			if err != nil {
				cs <- fmt.Errorf("ListComments %s/issues/%d: %w", repo, issueNumber, err)
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
