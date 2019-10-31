package main

import (
	"fmt"
	"net/url"
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
	source, err := url.Parse(args[0])
	if err != nil {
		return err
	}
	target, err := url.Parse(args[1])
	if err != nil {
		return err
	}
	sourceCli, err := createGitHubClient("GITHUB_MIGRATOR_SOURCE_TOKEN", urlToRoot(source))
	if err != nil {
		return err
	}
	sourceName, err := sourceCli.Login()
	if err != nil {
		return err
	}
	fmt.Printf("%#v\n", sourceName)
	targetCli, err := createGitHubClient("GITHUB_MIGRATOR_TARGET_TOKEN", urlToRoot(target))
	if err != nil {
		return err
	}
	targetName, err := targetCli.Login()
	if err != nil {
		return err
	}
	fmt.Printf("%#v\n", targetName)
	return err
}

func urlToRoot(u *url.URL) string {
	var userinfo string
	if u.User != nil {
		userinfo = u.User.String() + "@"
	}
	return fmt.Sprintf("%s://%s%s", u.Scheme, userinfo, u.Host)
}

func createGitHubClient(envName, root string) (github.Client, error) {
	token := os.Getenv(envName)
	if token == "" {
		return nil, fmt.Errorf("GitHub token not found (specify %s)", envName)
	}
	client := github.New(token, root)
	return client, nil
}
