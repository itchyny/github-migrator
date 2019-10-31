package migrator

import (
	"fmt"
	"io"

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
	source, target repo.Repo
}

func (m *migrator) Migrate() error {
	sourceRepo, err := m.source.Get()
	if err != nil {
		return err
	}
	targetRepo, err := m.target.Get()
	if err != nil {
		return err
	}
	fmt.Printf(
		"migrating: %s (%s) => %s (%s)\n",
		sourceRepo.Name, sourceRepo.HTMLURL,
		targetRepo.Name, targetRepo.HTMLURL,
	)

	is := m.source.ListIssues()
	for {
		i, err := is.Next()
		if err != nil {
			if err != io.EOF {
				return err
			}
			break
		}
		fmt.Printf("%#v\n", i)
	}
	return nil
}
