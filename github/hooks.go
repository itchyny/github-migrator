package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

// Hook represents a hook.
type Hook struct {
	Type      string      `json:"type"`
	ID        int         `json:"id"`
	Name      string      `json:"name"`
	Active    bool        `json:"active"`
	Events    []string    `json:"events"`
	Config    *HookConfig `json:"config"`
	CreatedAt string      `json:"created_at"`
	UpdatedAt string      `json:"updated_at"`
}

// HookConfig ...
type HookConfig struct {
	ContentType string `json:"content_type"`
	URL         string `json:"url"`
	InsecureSsl string `json:"insecure_ssl"`
	Secret      string `json:"secret,omitempty"`
}

// Hooks represents a collection of hooks.
type Hooks <-chan interface{}

// Next emits the next Hook.
func (hs Hooks) Next() (*Hook, error) {
	for x := range hs {
		switch x := x.(type) {
		case error:
			return nil, x
		case *Hook:
			return x, nil
		}
		break
	}
	return nil, io.EOF
}

// HooksFromSlice creates Hooks from a slice.
func HooksFromSlice(xs []*Hook) Hooks {
	hs := make(chan interface{})
	go func() {
		defer close(hs)
		for _, h := range xs {
			hs <- h
		}
	}()
	return hs
}

// HooksToSlice collects Hooks.
func HooksToSlice(hs Hooks) ([]*Hook, error) {
	xs := []*Hook{}
	for {
		h, err := hs.Next()
		if err != nil {
			if err != io.EOF {
				return nil, err
			}
			return xs, nil
		}
		xs = append(xs, h)
	}
}

func listHooksPath(repo string) string {
	return newPath("/repos/"+repo+"/hooks").
		query("per_page", "100").
		String()
}

// ListHooks lists the hooks.
func (c *client) ListHooks(repo string) Hooks {
	hs := make(chan interface{})
	go func() {
		defer close(hs)
		path := c.url(listHooksPath(repo))
		for {
			xs, next, err := c.listHooks(path)
			if err != nil {
				hs <- err
				break
			}
			for _, x := range xs {
				hs <- x
			}
			if next == "" {
				break
			}
			path = next
		}
	}()
	return Hooks(hs)
}

func (c *client) listHooks(path string) ([]*Hook, string, error) {
	res, err := c.get(path)
	if err != nil {
		return nil, "", err
	}
	defer res.Body.Close()

	var r []*Hook
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, "", err
	}

	return r, getNext(res.Header), nil
}

func getHookPath(repo string, hookID int) string {
	return newPath(fmt.Sprintf("/repos/%s/hooks/%d", repo, hookID)).
		String()
}

type hookOrError struct {
	Hook
	Message string `json:"message"`
}

// GetHook gets the hook.
func (c *client) GetHook(repo string, hookID int) (*Hook, error) {
	res, err := c.get(c.url(getHookPath(repo, hookID)))
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var r hookOrError
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, err
	}

	if r.Message != "" {
		return nil, fmt.Errorf("%s: %s", r.Message, "/hooks/"+fmt.Sprint(hookID))
	}

	return &r.Hook, nil
}

// CreateHookParams represents the paramter for CreateHook API.
type CreateHookParams struct {
	Name   string      `json:"name"`
	Active bool        `json:"active"`
	Events []string    `json:"events"`
	Config *HookConfig `json:"config"`
}

func createHookPath(repo string) string {
	return newPath(fmt.Sprintf("/repos/%s/hooks", repo)).
		String()
}

// CreateHook creates the hook.
func (c *client) CreateHook(repo string, params *CreateHookParams) (*Hook, error) {
	params.Name = "web"
	bs, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	body := bytes.NewReader(bs)
	res, err := c.post(c.url(createHookPath(repo)), body)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var r hookOrError
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, err
	}

	if r.Message != "" {
		return nil, fmt.Errorf("%s: %s", r.Message, "/hooks")
	}

	return &r.Hook, nil
}

// UpdateHookParams represents the paramter for UpdateHook API.
type UpdateHookParams struct {
	Active bool        `json:"active"`
	Events []string    `json:"events"`
	Config *HookConfig `json:"config"`
}

func updateHookPath(repo string, hookID int) string {
	return newPath(fmt.Sprintf("/repos/%s/hooks/%d", repo, hookID)).
		String()
}

// UpdateHook updates the hook.
func (c *client) UpdateHook(repo string, hookID int, params *UpdateHookParams) (*Hook, error) {
	bs, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	body := bytes.NewReader(bs)
	res, err := c.patch(c.url(updateHookPath(repo, hookID)), body)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var r hookOrError
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, err
	}

	if r.Message != "" {
		return nil, fmt.Errorf("%s: %s", r.Message, "/hooks/"+fmt.Sprint(hookID))
	}

	return &r.Hook, nil
}
