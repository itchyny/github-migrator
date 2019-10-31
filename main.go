package main

import (
	"fmt"
	"os"

	"github.com/itchyny/github-migrator/github"
)

const name = "github-migrator"

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s: %s\n", name, err)
		os.Exit(1)
	}
}

func run() error {
	cli, err := createGitHubClient("GITHUB_MIGRATOR_SOURCE_TOKEN")
	if err != nil {
		return err
	}
	fmt.Printf("%#v\n", cli)
	return err
}

func createGitHubClient(envName string) (github.Client, error) {
	token := os.Getenv(envName)
	if token == "" {
		return nil, fmt.Errorf("GitHub token not found (specify %s)", envName)
	}
	client := github.New(token)
	return client, nil
}
