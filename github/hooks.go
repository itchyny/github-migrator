package github

import (
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

// ListHooks lists the hooks.
func (c *client) ListHooks(repo string) Hooks {
	hs := make(chan interface{})
	go func() {
		defer close(hs)
		path := c.url(fmt.Sprintf("/repos/%s/hooks?per_page=100", repo))
		for {
			var xs []*Hook
			next, err := c.getList(path, &xs)
			if err != nil {
				if err.Error() != "Not Found" {
					hs <- fmt.Errorf("ListHooks %s: %w", repo, err)
				}
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

// GetHook gets the hook.
func (c *client) GetHook(repo string, hookID int) (*Hook, error) {
	var r Hook
	if err := c.get(c.url(fmt.Sprintf("/repos/%s/hooks/%d", repo, hookID)), &r); err != nil {
		return nil, fmt.Errorf("GetHook %s: %w", fmt.Sprintf("%s/hooks/%d", repo, hookID), err)
	}
	return &r, nil
}

// CreateHookParams represents the paramter for CreateHook API.
type CreateHookParams struct {
	Name   string      `json:"name"`
	Active bool        `json:"active"`
	Events []string    `json:"events"`
	Config *HookConfig `json:"config"`
}

// CreateHook creates a hook.
func (c *client) CreateHook(repo string, params *CreateHookParams) (*Hook, error) {
	params.Name = "web"
	var r Hook
	if err := c.post(c.url(fmt.Sprintf("/repos/%s/hooks", repo)), params, &r); err != nil {
		return nil, fmt.Errorf("CreateHook %s: %w", fmt.Sprintf("%s/hooks", repo), err)
	}
	return &r, nil
}

// UpdateHookParams represents the paramter for UpdateHook API.
type UpdateHookParams struct {
	Active bool        `json:"active"`
	Events []string    `json:"events"`
	Config *HookConfig `json:"config"`
}

// UpdateHook updates the hook.
func (c *client) UpdateHook(repo string, hookID int, params *UpdateHookParams) (*Hook, error) {
	var r Hook
	if err := c.patch(c.url(fmt.Sprintf("/repos/%s/hooks/%d", repo, hookID)), params, &r); err != nil {
		return nil, fmt.Errorf("UpdateHook %s: %w", fmt.Sprintf("%s/hooks/%d", repo, hookID), err)
	}
	return &r, nil
}
