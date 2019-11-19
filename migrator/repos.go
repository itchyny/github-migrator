package migrator

import (
	"fmt"

	"github.com/itchyny/github-migrator/github"
)

func (m *migrator) migrateRepo() error {
	sourceRepo, err := m.getSourceRepo()
	if err != nil {
		return err
	}

	targetRepo, err := m.getTargetRepo()
	if err != nil {
		return err
	}

	fmt.Printf(
		"[=>] migrating: %s (%s) => %s (%s)\n",
		sourceRepo.Name, sourceRepo.HTMLURL,
		targetRepo.Name, targetRepo.HTMLURL,
	)

	if params, ok := buildUpdateRepoParams(sourceRepo, targetRepo); ok {
		fmt.Printf("[|>] updating the repository: %s\n", targetRepo.HTMLURL)
		_, err = m.target.Update(params)
		if err != nil {
			return err
		}
	}
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

func buildUpdateRepoParams(sourceRepo, targetRepo *github.Repo) (*github.UpdateRepoParams, bool) {
	var update bool
	params := &github.UpdateRepoParams{
		Name:        targetRepo.Name,
		Description: targetRepo.Description,
		Homepage:    targetRepo.Homepage,
		Private:     targetRepo.Private,
	}
	if params.Description == "" && sourceRepo.Description != "" {
		params.Description = sourceRepo.Description
		update = true
	}
	if params.Homepage == "" && sourceRepo.Homepage != "" {
		params.Homepage = sourceRepo.Homepage
		update = true
	}
	// other fields should not be overwritten unless examined carefully
	return params, update
}
