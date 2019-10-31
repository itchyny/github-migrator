package migrator

import (
	"testing"

	"github.com/itchyny/github-migrator/github"
	"github.com/itchyny/github-migrator/repo"
)

func TestMigrator(t *testing.T) {
	sourceCli := github.NewMockClient()
	targetCli := github.NewMockClient()
	source := repo.New(sourceCli, "example/source")
	target := repo.New(targetCli, "example/target")
	var _ Migrator = New(source, target)
}
