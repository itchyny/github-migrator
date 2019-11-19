package github

import (
	"bytes"
	"encoding/json"
	"fmt"
)

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

func issueImportPath(repo string) string {
	return newPath(fmt.Sprintf("/repos/%s/import/issues", repo)).
		String()
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

type importResultOrError struct {
	ImportResult
	Message string    `json:"message"`
	Errors  apiErrors `json:"errors"`
}

// Import imports an importing object.
func (c *client) Import(repo string, x *Import) (*ImportResult, error) {
	bs, err := json.Marshal(x)
	if err != nil {
		return nil, err
	}
	body := bytes.NewReader(bs)
	res, err := c.post(c.url(issueImportPath(repo)), body)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var r importResultOrError
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, err
	}
	if r.Message != "" {
		return nil, fmt.Errorf("Import %s: %s: %s", fmt.Sprintf("%s/import/issues", repo), r.Message, r.Errors)
	}
	if r.Errors != nil {
		return nil, r.Errors
	}
	return &r.ImportResult, nil
}

func getImportPath(repo string, id int) string {
	return newPath(fmt.Sprintf("/repos/%s/import/issues/%d", repo, id)).
		String()
}

// GetImport gets the importing status.
func (c *client) GetImport(repo string, id int) (*ImportResult, error) {
	res, err := c.get(c.url(getImportPath(repo, id)))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var r importResultOrError
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, err
	}
	if r.Message != "" || r.Errors != nil {
		errorPrefix := "GetImport " + fmt.Sprintf("%s/import/issues/%d", repo, id)
		if r.Message != "" {
			if r.Errors == nil {
				return nil, fmt.Errorf("%s: %s", errorPrefix, r.Message)
			}
			return nil, fmt.Errorf("%s: %s: %s", errorPrefix, r.Message, r.Errors)
		}
		return nil, fmt.Errorf("%s: %s", errorPrefix, r.Errors.Error())
	}
	return &r.ImportResult, nil
}
