package github

import "fmt"

// Import represents an importing object.
type Import struct {
	Issue    *ImportIssue     `json:"issue"`
	Comments []*ImportComment `json:"comments"`
}

// ImportIssue represents an importing issue.
type ImportIssue struct {
	Title     string   `json:"title"`
	Body      string   `json:"body"`
	CreatedAt string   `json:"created_at"`
	UpdatedAt string   `json:"updated_at"`
	Closed    bool     `json:"closed"`
	ClosedAt  string   `json:"closed_at,omitempty"`
	Labels    []string `json:"labels,omitempty"`
	Assignee  string   `json:"assignee,omitempty"`
}

// ImportComment represents an importing comment.
type ImportComment struct {
	Body      string `json:"body"`
	CreatedAt string `json:"created_at"`
}

// ImportResult represents the result of import.
type ImportResult struct {
	ID              int    `json:"id"`
	Status          string `json:"status"`
	URL             string `json:"url"`
	ImportIssuesURL string `json:"import_issues_url"`
	RepositoryURL   string `json:"repository_url"`
	CreatedAt       string `json:"created_at"`
	UpdatedAt       string `json:"updated_at"`
}

// Import imports an importing object.
func (c *client) Import(repo string, params *Import) (*ImportResult, error) {
	var r ImportResult
	if err := c.post(c.url(fmt.Sprintf("/repos/%s/import/issues", repo)), params, &r); err != nil {
		return nil, fmt.Errorf("Import %s: %w", fmt.Sprintf("%s/import/issues", repo), err)
	}
	return &r, nil
}

// GetImport gets the importing status.
func (c *client) GetImport(repo string, id int) (*ImportResult, error) {
	var r ImportResult
	if err := c.get(c.url(fmt.Sprintf("/repos/%s/import/issues/%d", repo, id)), &r); err != nil {
		return nil, fmt.Errorf("GetImport %s: %w", fmt.Sprintf("%s/import/issues/%d", repo, id), err)
	}
	return &r, nil
}
