package migrator

import (
	"fmt"

	"github.com/itchyny/github-migrator/github"
)

func (m *migrator) checkRepos() error {
	sourceRepo, err := m.getSourceRepo()
	if err != nil {
		return err
	}

	targetRepo, err := m.getTargetRepo()
	if err != nil {
		return err
	}

	fmt.Printf(
		"migrating: %s (%s) => %s (%s)\n",
		sourceRepo.Name, sourceRepo.HTMLURL,
		targetRepo.Name, targetRepo.HTMLURL,
	)

	return nil
}

func (m *migrator) getSourceRepo() (*github.Repo, error) {
	if m.sourceRepo != nil {
		return m.sourceRepo, nil
	}
	repo, err := m.source.Get()
	if err != nil {
		return nil, err
	}
	m.sourceRepo = repo
	return repo, nil
}

func (m *migrator) getTargetRepo() (*github.Repo, error) {
	if m.targetRepo != nil {
		return m.targetRepo, nil
	}
	repo, err := m.target.Get()
	if err != nil {
		return nil, err
	}
	m.targetRepo = repo
	return repo, nil
}
