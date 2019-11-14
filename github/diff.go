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
	req, err := c.request("GET", c.url(getDiffPath(repo, sha)), nil)
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

	if !bytes.HasPrefix(bs, []byte("diff --git")) {
		var r struct {
			Message string `json:"message"`
		}
		if err := json.NewDecoder(bytes.NewReader(bs)).Decode(&r); err != nil {
			return "", err
		}
		if r.Message != "" {
			return "", errors.New(r.Message)
		}
		return "", fmt.Errorf("failed to get diff of: %s, %s", repo, sha)
	}

	return string(bs), nil
}
