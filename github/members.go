package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
)

// Member represents a member.
type Member User

// Members represents a collection of comments.
type Members <-chan interface{}

// Next emits the next Member.
func (cs Members) Next() (*Member, error) {
	for x := range cs {
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
	cs := make(chan interface{})
	go func() {
		defer close(cs)
		for _, i := range xs {
			cs <- i
		}
	}()
	return cs
}

// MembersToSlice collects Members.
func MembersToSlice(cs Members) ([]*Member, error) {
	xs := []*Member{}
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

func listMembersPath(org string) string {
	return newPath(fmt.Sprintf("/orgs/%s/members", org)).
		String()
}

// ListMembers lists the members of the organization.
func (c *client) ListMembers(org string) Members {
	cs := make(chan interface{})
	go func() {
		defer close(cs)
		path := c.url(listMembersPath(org))
		for {
			xs, next, err := c.listMembers(path)
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
	return Members(cs)
}

func (c *client) listMembers(path string) ([]*Member, string, error) {
	res, err := c.get(path)
	if err != nil {
		return nil, "", err
	}
	defer res.Body.Close()

	bs, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, "", err
	}

	var r []*Member
	if err := json.NewDecoder(bytes.NewReader(bs)).Decode(&r); err != nil {
		var em map[string]string
		if err := json.NewDecoder(bytes.NewReader(bs)).Decode(&em); err == nil {
			if mess, ok := em["message"]; ok && mess == "Not Found" {
				return nil, "", nil
			}
		}
		return nil, "", err
	}

	return r, getNext(res.Header), nil
}
