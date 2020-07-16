package repo

import "github.com/itchyny/github-migrator/github"

// ListHooks lists the hooks.
func (r *Repo) ListHooks() github.Hooks {
	return r.cli.ListHooks(r.path)
}

// GetHook gets the hook.
func (r *Repo) GetHook(hookID int) (*github.Hook, error) {
	return r.cli.GetHook(r.path, hookID)
}

// CreateHook creates a hook.
func (r *Repo) CreateHook(params *github.CreateHookParams) (*github.Hook, error) {
	return r.cli.CreateHook(r.path, params)
}

// UpdateHook updates the hook.
func (r *Repo) UpdateHook(hookID int, params *github.UpdateHookParams) (*github.Hook, error) {
	return r.cli.UpdateHook(r.path, hookID, params)
}
