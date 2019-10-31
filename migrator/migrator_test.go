package migrator

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/itchyny/github-migrator/github"
	"github.com/itchyny/github-migrator/repo"
)

func TestMigratorMigrate(t *testing.T) {
	source := repo.New(github.NewMockClient(
		github.MockGetRepo(func(path string) (*github.Repo, error) {
			return &github.Repo{
				Name:        "source",
				Description: "Source repository.",
				HTMLURL:     "http://localhost/example/source",
			}, nil
		}),
		github.MockListIssues(func(path string, _ *github.ListIssuesParams) github.Issues {
			return github.IssuesFromSlice([]*github.Issue{
				{
					Number:  1,
					Title:   "Example title 1",
					State:   "closed",
					Body:    "Example body 1",
					HTMLURL: "http://localhost/example/source/issues/1",
				},
				{
					Number:  2,
					Title:   "Example title 2",
					State:   "open",
					Body:    "Example body 2",
					HTMLURL: "http://localhost/example/source/issues/2",
				},
			})
		}),
		github.MockListComments(func(path string, issueNumber int) github.Comments {
			assert.Equal(t, issueNumber, 2)
			return github.CommentsFromSlice([]*github.Comment{
				{
					Body:    "Example body 1",
					HTMLURL: "http://localhost/example/source/issues/1#issuecomment-1",
				},
				{
					Body:    "Example body 2",
					HTMLURL: "http://localhost/example/source/issues/1#issuecomment-2",
				},
			})
		}),
	), "example/source")

	target := repo.New(github.NewMockClient(
		github.MockGetRepo(func(path string) (*github.Repo, error) {
			return &github.Repo{
				Name:        "target",
				Description: "target repository.",
				HTMLURL:     "http://localhost/example/target",
			}, nil
		}),
		github.MockListIssues(func(path string, _ *github.ListIssuesParams) github.Issues {
			return github.IssuesFromSlice([]*github.Issue{
				{
					Number:  1,
					Title:   "Example title 1",
					State:   "closed",
					Body:    "Example body 1",
					HTMLURL: "http://localhost/example/target/issues/1",
				},
			})
		}),
	), "example/target")

	mig := New(source, target)
	assert.Nil(t, mig.Migrate())
}
