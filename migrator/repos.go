package migrator

import (
	"fmt"

	"github.com/itchyny/github-migrator/github"
)

func (m *migrator) migrateRepo() error {
	fmt.Printf(
		"[=>] migrating: %s (%s) => %s (%s)\n",
		m.sourceRepo.Name, m.sourceRepo.HTMLURL,
		m.targetRepo.Name, m.targetRepo.HTMLURL,
	)

	if params, ok := buildUpdateRepoParams(m.sourceRepo, m.targetRepo); ok {
		fmt.Printf("[|>] updating the repository: %s\n", m.targetRepo.HTMLURL)
		if _, err := m.target.Update(params); err != nil {
			return err
		}
	}
	return nil
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
