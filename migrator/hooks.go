package migrator

import (
	"fmt"
	"reflect"

	"github.com/itchyny/github-migrator/github"
)

func (m *migrator) migrateHooks() error {
	sourceHooks, err := github.HooksToSlice(m.source.ListHooks())
	if err != nil {
		return err
	}
	targetHooks, err := github.HooksToSlice(m.target.ListHooks())
	if err != nil {
		return err
	}
	for _, sourceHook := range sourceHooks {
		fmt.Printf("[=>] migrating a hook: %s\n", sourceHook.Config.URL)
		var exists bool
		for _, targetHook := range targetHooks {
			if sourceHook.Name == targetHook.Name &&
				sourceHook.Config.URL == targetHook.Config.URL {
				if sourceHook.Active != targetHook.Active ||
					!reflect.DeepEqual(sourceHook.Events, targetHook.Events) ||
					!reflect.DeepEqual(sourceHook.Config, targetHook.Config) {
					fmt.Printf("[|>] updating an existing hook: %s\n", targetHook.Config.URL)
					if _, err := m.target.UpdateHook(targetHook.ID, &github.UpdateHookParams{
						Active: sourceHook.Active,
						Events: sourceHook.Events,
						Config: sourceHook.Config,
					}); err != nil {
						return err
					}
				} else {
					fmt.Printf("[--] skipping: %s (already exists)\n", sourceHook.Config.URL)
				}
				exists = true
				break
			}
		}
		if exists {
			continue
		}
		fmt.Printf("[>>] creating a new hook: %s\n", sourceHook.Config.URL)
		if _, err := m.target.CreateHook(&github.CreateHookParams{
			Active: sourceHook.Active,
			Events: sourceHook.Events,
			Config: sourceHook.Config,
		}); err != nil {
			return err
		}
	}
	return nil
}
