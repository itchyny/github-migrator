package main

import (
	"fmt"
	"os"

	"github.com/itchyny/github-migrator/github"
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
	source, target := args[0], args[1]
	sourceCli, err := createGitHubClient(
		"GITHUB_MIGRATOR_SOURCE_TOKEN",
		"GITHUB_MIGRATOR_SOURCE_API_ENDPOINT",
	)
	if err != nil {
		return err
	}
	sourceName, err := sourceCli.Login()
	if err != nil {
		return err
	}
	fmt.Printf("%s\n", sourceName)
	targetCli, err := createGitHubClient(
		"GITHUB_MIGRATOR_TARGET_TOKEN",
		"GITHUB_MIGRATOR_TARGET_API_ENDPOINT",
	)
	if err != nil {
		return err
	}
	targetName, err := targetCli.Login()
	if err != nil {
		return err
	}
	fmt.Printf("%s\n", targetName)
	fmt.Printf("[%s] %s => [%s] %s\n", sourceCli.Hostname(), source, targetCli.Hostname(), target)
	return err
}

func createGitHubClient(tokenEnv, endpointEnv string) (github.Client, error) {
	token := os.Getenv(tokenEnv)
	if token == "" {
		return nil, fmt.Errorf("GitHub token not found (specify %s)", tokenEnv)
	}
	endpoint := os.Getenv(endpointEnv)
	if endpoint == "" {
		return nil, fmt.Errorf("GitHub endpoint not found (specify %s)", endpointEnv)
	}
	return github.New(token, endpoint), nil
}
