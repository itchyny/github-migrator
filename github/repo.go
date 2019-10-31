package github

import (
	"encoding/json"
	"errors"
)

// Repo represents a repository.
type Repo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	HTMLURL     string `json:"html_url"`
}

type repoOrError struct {
	Repo
	Message string `json:"message"`
}

func getRepoPath(repo string) string {
	return newPath("/repos/" + repo).
		String()
}

func (c *client) GetRepo(path string) (*Repo, error) {
	res, err := c.get(c.url(getRepoPath(path)))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var r repoOrError
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, err
	}

	if r.Message != "" {
		return nil, errors.New(r.Message)
	}

	return &r.Repo, nil
}
