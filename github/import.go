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
	return newPath("/repos/" + repo + "/import/issues").
		String()
}

type importRes struct {
	Status string `json:"status"`
}

type importOrError struct {
	importRes
	Message string `json:"message"`
}

// Import imports an importing object.
func (c *client) Import(path string, x *Import) error {
	bs, err := json.Marshal(x)
	if err != nil {
		return err
	}
	body := bytes.NewReader(bs)
	res, err := c.post(c.url(issueImportPath(path)), body)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	var r importOrError
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return err
	}
	if r.Message != "" {
		return fmt.Errorf("%s: %s", r.Message, issueImportPath(path))
	}
	if r.importRes.Status != "imported" && r.importRes.Status != "pending" {
		return fmt.Errorf("%s: %s (%#v)", r.importRes.Status, issueImportPath(path), r)
	}
	return nil
}
