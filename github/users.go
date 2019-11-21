package github

import (
	"fmt"
	"io"
)

// User represents a user.
type User struct {
	Login   string `json:"login"`
	HTMLURL string `json:"html_url"`
}

// GetLogin ...
func (c *client) GetLogin() (*User, error) {
	var r User
	if err := c.get(c.url("/user"), &r); err != nil {
		return nil, fmt.Errorf("GetLogin %s: %w", "/user", err)
	}
	return &r, nil
}

// Users represents a collection of users.
type Users <-chan interface{}

// Next emits the next User.
func (cs Users) Next() (*User, error) {
	for x := range cs {
		switch x := x.(type) {
		case error:
			return nil, x
		case *User:
			return x, nil
		}
		break
	}
	return nil, io.EOF
}

// UsersFromSlice creates Users from a slice.
func UsersFromSlice(xs []*User) Users {
	cs := make(chan interface{})
	go func() {
		defer close(cs)
		for _, p := range xs {
			cs <- p
		}
	}()
	return cs
}

// UsersToSlice collects Users.
func UsersToSlice(cs Users) ([]*User, error) {
	xs := []*User{}
	for {
		p, err := cs.Next()
		if err != nil {
			if err != io.EOF {
				return nil, err
			}
			return xs, nil
		}
		xs = append(xs, p)
	}
}

// ListUsers lists all the users.
func (c *client) ListUsers() Users {
	cs := make(chan interface{})
	go func() {
		defer close(cs)
		path := c.url("/users?per_page=100")
		for {
			var xs []*User
			next, err := c.getList(path, &xs)
			if err != nil {
				cs <- fmt.Errorf("ListUsers /users: %w", err)
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
	return Users(cs)
}

// GetUser ...
func (c *client) GetUser(name string) (*User, error) {
	var r User
	if err := c.get(c.url(fmt.Sprintf("/users/%s", name)), &r); err != nil {
		return nil, fmt.Errorf("GetUser %s: %w", fmt.Sprintf("/user/%s", name), err)
	}
	return &r, nil
}
