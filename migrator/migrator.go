package migrator

import (
	"github.com/itchyny/github-migrator/github"
	"github.com/itchyny/github-migrator/repo"
)

// Migrator represents a GitHub migrator.
type Migrator interface {
	Migrate() error
}

// New creates a new Migrator.
func New(source, target repo.Repo) Migrator {
	return &migrator{source: source, target: target}
}

type migrator struct {
	source, target         repo.Repo
	sourceRepo, targetRepo *github.Repo
}

// Migrate the repository.
func (m *migrator) Migrate() error {
	if err := m.checkRepos(); err != nil {
		return err
	}
	if err := m.migrateIssues(); err != nil {
		return err
	}
	return nil
}
