package migrator

import "fmt"

func (m *migrator) checkRepos() error {
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

	return nil
}
