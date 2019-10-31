package github

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
)

// Client represents a GitHub client.
type Client interface {
	Login() (string, error)
	GetRepo(string) (*Repo, error)
	ListIssues(string, *ListIssuesParams) Issues
}

// New creates a new GitHub client.
func New(token, endpoint string) Client {
	cli := &http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: endpoint != "https://api.github.com",
		},
	}}
	return &client{token: token, endpoint: endpoint, client: cli}
}

type client struct {
	token, endpoint string
	client          *http.Client
}

func (c *client) url(path string) string {
	return c.endpoint + path
}

func (c *client) get(path string) (*http.Response, error) {
	req, err := c.request("GET", path, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "token "+c.token)
	fmt.Printf("fetching: %s\n", req.URL)
	return c.client.Do(req)
}

func (c *client) request(method, path string, body io.Reader) (*http.Request, error) {
	return http.NewRequest(method, path, body)
}
