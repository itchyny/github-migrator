package github

import "net/url"

func (c *client) Hostname() string {
	u, _ := url.Parse(c.root)
	return u.Hostname()
}
