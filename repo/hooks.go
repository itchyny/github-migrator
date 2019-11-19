package repo

import "github.com/itchyny/github-migrator/github"

// ListHooks lists the hooks.
func (r *repo) ListHooks() github.Hooks {
	return r.cli.ListHooks(r.path)
}

// GetHook gets the hook.
func (r *repo) GetHook(hookID int) (*github.Hook, error) {
	return r.cli.GetHook(r.path, hookID)
}

// CreateHook creates a hook.
func (r *repo) CreateHook(params *github.CreateHookParams) (*github.Hook, error) {
	return r.cli.CreateHook(r.path, params)
}

// UpdateHook updates the hook.
func (r *repo) UpdateHook(hookID int, params *github.UpdateHookParams) (*github.Hook, error) {
	return r.cli.UpdateHook(r.path, hookID, params)
}
