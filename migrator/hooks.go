package migrator

import (
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
		var exists bool
		for _, targetHook := range targetHooks {
			if sourceHook.Name == targetHook.Name &&
				sourceHook.Config.URL == targetHook.Config.URL {
				if sourceHook.Active != targetHook.Active ||
					!reflect.DeepEqual(sourceHook.Events, targetHook.Events) ||
					!reflect.DeepEqual(sourceHook.Config, targetHook.Config) {
					if _, err := m.target.UpdateHook(targetHook.ID, &github.UpdateHookParams{
						Active: sourceHook.Active,
						Events: sourceHook.Events,
						Config: sourceHook.Config,
					}); err != nil {
						return err
					}
				}
				exists = true
				break
			}
		}
		if exists {
			continue
		}
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
