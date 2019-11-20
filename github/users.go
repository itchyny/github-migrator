package github

import "fmt"

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

// GetUser ...
func (c *client) GetUser(name string) (*User, error) {
	var r User
	if err := c.get(c.url(fmt.Sprintf("/users/%s", name)), &r); err != nil {
		return nil, fmt.Errorf("GetUser %s: %w", fmt.Sprintf("/user/%s", name), err)
	}
	return &r, nil
}
