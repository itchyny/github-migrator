package repo

import (
	"testing"

	"github.com/itchyny/github-migrator/github"
)

func TestRepo(t *testing.T) {
	cli := github.New("token", "https://github.com")
	var _ Repo = New(cli, "example/test")
}
