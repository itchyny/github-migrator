package github

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
)

func getDiffPath(repo string, sha string) string {
	return newPath(fmt.Sprintf("/repos/%s/commits/%s", repo, sha)).
		String()
}

func (c *client) GetDiff(repo string, sha string) (string, error) {
	return c.getDiff(getDiffPath(repo, sha))
}

func getComparePath(repo string, base, head string) string {
	return newPath(fmt.Sprintf("/repos/%s/compare/%s...%s", repo, base, head)).
		String()
}

func (c *client) GetCompare(repo string, base, head string) (string, error) {
	return c.getDiff(getComparePath(repo, base, head))
}

func (c *client) getDiff(path string) (string, error) {
	req, err := c.request("GET", c.url(path), nil)
	if err != nil {
		return "", err
	}
	fmt.Printf("fetching: %s\n", req.URL)
	req.Header.Add("Accept", "application/vnd.github.v3.diff")
	res, err := c.client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	bs, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	if len(bs) > 0 && !bytes.HasPrefix(bs, []byte("diff --git")) {
		var r struct {
			Message string `json:"message"`
		}
		if err := json.NewDecoder(bytes.NewReader(bs)).Decode(&r); err != nil {
			return "", err
		}
		if r.Message != "" {
			return "", errors.New(r.Message)
		}
		return "", fmt.Errorf("failed to get %s", path)
	}

	return string(bs), nil
}
