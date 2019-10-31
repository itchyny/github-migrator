package main

import (
	"fmt"
	"os"

	"github.com/itchyny/github-migrator/github"
	"github.com/itchyny/github-migrator/migrator"
	"github.com/itchyny/github-migrator/repo"
)

const name = "github-migrator"

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", name, err)
		os.Exit(1)
	}
}

func run(args []string) error {
	if len(args) != 2 {
		return fmt.Errorf("usage: %s <source> <target>", name)
	}
	mig, err := createMigrator(args[0], args[1])
	if err != nil {
		return err
	}
	return mig.Migrate()
}

func createGitHubClient(tokenEnv, endpointEnv string) (github.Client, error) {
	token := os.Getenv(tokenEnv)
	if token == "" {
		return nil, fmt.Errorf("GitHub token not found (specify %s)", tokenEnv)
	}
	endpoint := os.Getenv(endpointEnv)
	if endpoint == "" {
		endpoint = "https://api.github.com"
	}
	cli := github.New(token, endpoint)
	name, err := cli.Login()
	if err != nil {
		return nil, fmt.Errorf("%s (or you may want to set %s)", err, endpointEnv)
	}
	fmt.Printf("login succeeded: %s\n", name)
	return cli, nil
}

func createMigrator(sourcePath, targetPath string) (migrator.Migrator, error) {
	sourceCli, err := createGitHubClient(
		"GITHUB_MIGRATOR_SOURCE_TOKEN",
		"GITHUB_MIGRATOR_SOURCE_API_ENDPOINT",
	)
	if err != nil {
		return nil, err
	}
	targetCli, err := createGitHubClient(
		"GITHUB_MIGRATOR_TARGET_TOKEN",
		"GITHUB_MIGRATOR_TARGET_API_ENDPOINT",
	)
	if err != nil {
		return nil, err
	}
	source := repo.New(sourceCli, sourcePath)
	target := repo.New(targetCli, targetPath)
	return migrator.New(source, target), nil
}
