package github

import (
	"fmt"
	"io"
)

// Member represents a member.
type Member User

// Members represents a collection of comments.
type Members <-chan interface{}

// Next emits the next Member.
func (ms Members) Next() (*Member, error) {
	for x := range ms {
		switch x := x.(type) {
		case error:
			return nil, x
		case *Member:
			return x, nil
		}
		break
	}
	return nil, io.EOF
}

// MembersFromSlice creates Members from a slice.
func MembersFromSlice(xs []*Member) Members {
	ms := make(chan interface{})
	go func() {
		defer close(ms)
		for _, m := range xs {
			ms <- m
		}
	}()
	return ms
}

// MembersToSlice collects Members.
func MembersToSlice(ms Members) ([]*Member, error) {
	xs := []*Member{}
	for {
		m, err := ms.Next()
		if err != nil {
			if err != io.EOF {
				return nil, err
			}
			return xs, nil
		}
		xs = append(xs, m)
	}
}

func listMembersPath(org string) string {
	return newPath(fmt.Sprintf("/orgs/%s/members", org)).
		String()
}

// ListMembers lists the members of the organization.
func (c *client) ListMembers(org string) Members {
	ms := make(chan interface{})
	go func() {
		defer close(ms)
		path := c.url(listMembersPath(org))
		for {
			var xs []*Member
			next, err := c.getList(path, &xs)
			if err != nil {
				if err.Error() != "Not Found" {
					ms <- fmt.Errorf("ListMembers %s: %w", org, err)
				}
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
	return Members(ms)
}
