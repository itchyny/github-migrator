package migrator

import (
	"fmt"

	"github.com/itchyny/github-migrator/repo"
)

type Migrator interface {
	Migrate() error
}

func New(source, target repo.Repo) *migrator {
	fmt.Printf("%s => %s\n", source.Name(), target.Name())
	return &migrator{source: source, target: target}
}

type migrator struct {
	source, target repo.Repo
}

func (m *migrator) Migrate() error {
	return nil
}
