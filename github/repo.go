package github

import (
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

func (c *client) GetRepo(repo string) (*Repo, error) {
	var r Repo
	if err := c.get(c.url(fmt.Sprintf("/repos/%s", repo)), &r); err != nil {
		return nil, fmt.Errorf("GetRepo %s: %w", repo, err)
	}
	return &r, nil
}

// UpdateRepoParams represents a parameter on updating a repository.
type UpdateRepoParams struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Homepage    string `json:"homepage"`
	Private     bool   `json:"private"`
}

// UpdateRepo updates a repository.
func (c *client) UpdateRepo(repo string, params *UpdateRepoParams) (*Repo, error) {
	var r Repo
	if err := c.patch(c.url(fmt.Sprintf("/repos/%s", repo)), params, &r); err != nil {
		return nil, fmt.Errorf("UpdateRepo %s: %w", repo, err)
	}
	return &r, nil
}
