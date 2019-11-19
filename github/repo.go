package github

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// Repo represents a repository.
type Repo struct {
	Name        string `json:"name"`
	FullName    string `json:"full_name"`
	Description string `json:"description"`
	Homepage    string `json:"homepage"`
	HTMLURL     string `json:"html_url"`
	Private     bool   `json:"private"`
}

type repoOrError struct {
	Repo
	Message string `json:"message"`
}

func getRepoPath(repo string) string {
	return newPath(fmt.Sprintf("/repos/%s", repo)).
		String()
}

func (c *client) GetRepo(repo string) (*Repo, error) {
	res, err := c.get(c.url(getRepoPath(repo)))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var r repoOrError
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, err
	}

	if r.Message != "" {
		return nil, fmt.Errorf("GetRepo %s: %s", repo, r.Message)
	}

	return &r.Repo, nil
}

// UpdateRepoParams represents a parameter on updating a repository.
type UpdateRepoParams struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Homepage    string `json:"homepage"`
	Private     bool   `json:"private"`
}

func updateRepoPath(repo string) string {
	return newPath(fmt.Sprintf("/repos/%s", repo)).
		String()
}

// UpdateRepo updates a repository.
func (c *client) UpdateRepo(repo string, params *UpdateRepoParams) (*Repo, error) {
	bs, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	body := bytes.NewReader(bs)
	res, err := c.patch(c.url(updateRepoPath(repo)), body)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var r repoOrError
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, err
	}

	if r.Message != "" {
		return nil, fmt.Errorf("UpdateRepo %s: %s", repo, r.Message)
	}

	return &r.Repo, nil
}
