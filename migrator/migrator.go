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
func New(source, target repo.Repo, userMapping map[string]string) Migrator {
	return &migrator{source: source, target: target, userMapping: userMapping}
}

type migrator struct {
	source, target         repo.Repo
	sourceRepo, targetRepo *github.Repo
	userMapping            map[string]string
	members                []*github.Member
	userByName             map[string]*github.User
	errorUserByName        map[string]error
	projects               []*github.Project
}

// Migrate the repository.
func (m *migrator) Migrate() error {
	if err := m.migrateRepo(); err != nil {
		return err
	}
	if err := m.migrateLabels(); err != nil {
		return err
	}
	// projects should be imported before issues
	if err := m.migrateProjects(); err != nil {
		return err
	}
	if err := m.migrateIssues(); err != nil {
		return err
	}
	if err := m.migrateHooks(); err != nil {
		return err
	}
	return nil
}
