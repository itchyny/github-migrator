package migrator

import (
	"fmt"

	"github.com/itchyny/github-migrator/repo"
)

// Migrator represents a GitHub migrator.
type Migrator interface {
	Migrate() error
}

// New creates a new Migrator.
func New(source, target repo.Repo) Migrator {
	fmt.Printf("%s => %s\n", source.Name(), target.Name())
	return &migrator{source: source, target: target}
}

type migrator struct {
	source, target repo.Repo
}

func (m *migrator) Migrate() error {
	is, err := m.source.ListIssues()
	if err != nil {
		return err
	}
	for _, i := range is {
		fmt.Printf("%#v\n", i)
	}
	return nil
}
