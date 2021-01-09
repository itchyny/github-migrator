package github

import (
	"fmt"
	"io"
	"io/ioutil"
	"strings"
)

const maxDiffSize = 1 * 1024 * 1024

func (c *client) GetDiff(repo string, sha string) (string, error) {
	return c.getDiff("GetDiff", fmt.Sprintf("/repos/%s/commits/%s", repo, sha))
}

func (c *client) GetCompare(repo string, base, head string) (string, error) {
	return c.getDiff("GetCompare", fmt.Sprintf("/repos/%s/compare/%s...%s", repo, base, head))
}

func (c *client) getDiff(name, path string) (string, error) {
	req, err := c.request("GET", c.url(path), nil)
	if err != nil {
		return "", err
	}
	req.Header.Add("Accept", "application/vnd.github.v3.diff")
	res, _, err := c.doReq(req)
	if err != nil {
		return "", fmt.Errorf("%s %s: %w", name, strings.TrimPrefix(path, "/repos/"), err)
	}
	defer res.Body.Close()

	bs, err := ioutil.ReadAll(&io.LimitedReader{R: res.Body, N: maxDiffSize})
	if err != nil {
		return "", err
	}
	return string(bs), nil
}
