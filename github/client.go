package github

import (
	"crypto/tls"
	"io"
	"net/http"
)

// Client represents a GitHub client.
type Client interface {
	Login() (string, error)
	Hostname() string
	ListIssues(string, *ListIssuesParams) ([]*Issue, error)
}

// New creates a new GitHub client.
func New(token, root string) Client {
	cli := &http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}}
	return &client{token: token, root: root, client: cli}
}

type client struct {
	token, root string
	client      *http.Client
}

func (c *client) url(path string) string {
	return c.root + path
}

func (c *client) get(path string) (*http.Response, error) {
	req, err := c.request("GET", path, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "token "+c.token)
	return c.client.Do(req)
}

func (c *client) request(method, path string, body io.Reader) (*http.Request, error) {
	return http.NewRequest(method, path, body)
}
