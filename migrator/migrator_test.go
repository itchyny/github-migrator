package migrator

import (
	"fmt"
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
				FullName:    "example/source",
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
					User: &github.User{
						Login: "sample-user-1",
					},
				},
				{
					Number:  2,
					Title:   "Example title 2",
					State:   "open",
					Body:    "Example body 2\nSee http://localhost/example/source/issues/1.",
					HTMLURL: "http://localhost/example/source/issues/2",
					User: &github.User{
						Login: "sample-user-2",
					},
					Assignee: &github.User{
						Login: "sample-user-2",
					},
					Labels: []*github.Label{
						{Name: "label1"},
						{Name: "label2"},
					},
				},
				{
					Number:  3,
					Title:   "Example title 3",
					State:   "open",
					Body:    "Example body 3",
					HTMLURL: "http://localhost/example/source/pull/3",
					User: &github.User{
						Login: "sample-user-3",
					},
					PullRequest: &github.IssuePullRequest{
						URL:     "http://localhost/example/source/pulls/3",
						HTMLURL: "http://localhost/example/source/pull/3",
					},
				},
			})
		}),
		github.MockListComments(func(path string, issueNumber int) github.Comments {
			switch issueNumber {
			case 2:
				return github.CommentsFromSlice([]*github.Comment{
					{
						Body:    "Example comment body 1",
						HTMLURL: "http://localhost/example/source/issues/1#issuecomment-1",
						User: &github.User{
							Login: "sample-user-1",
						},
					},
					{
						Body:    "Example comment body 2\nRef: http://localhost/example/source/issues/1.",
						HTMLURL: "http://localhost/example/source/issues/1#issuecomment-2",
						User: &github.User{
							Login: "sample-user-2",
						},
					},
				})
			case 3:
				return github.CommentsFromSlice([]*github.Comment{})
			default:
				assert.Nil(t, fmt.Errorf("unexpected issue number: %d", issueNumber))
				return nil
			}
		}),
	), "example/source")

	var assertImport func(string, *github.Import)
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
					User: &github.User{
						Login: "sample-user-1",
					},
				},
			})
		}),
		github.MockImport(func(path string, x *github.Import) error {
			assertImport(path, x)
			return nil
		}),
	), "example/target")

	var importCount int
	assertImport = func(path string, x *github.Import) {
		switch importCount {
		case 0:
			assert.Equal(t, path, "/repos/example/target/import/issues")
			assert.Equal(t, x.Issue.Title, "Example title 2")
			assert.Contains(t, x.Issue.Body, `<img src="https://github.com/sample-user-2.png" width="35">`)
			assert.Contains(t, x.Issue.Body, `Original issue by @sample-user-2 - imported from <a href="http://localhost/example/source/issues/2">example/source#2</a>`)
			assert.Contains(t, x.Issue.Body, `Example body 2`)
			assert.Contains(t, x.Issue.Body, `See http://localhost/example/target/issues/1.`)
			assert.Equal(t, x.Issue.Assignee, "sample-user-2")
			assert.Equal(t, x.Issue.Labels, []string{"label1", "label2"})

			assert.Len(t, x.Comments, 2)
			assert.Contains(t, x.Comments[0].Body, `<img src="https://github.com/sample-user-1.png" width="35">`)
			assert.Contains(t, x.Comments[0].Body, `@sample-user-1 commented`)
			assert.Contains(t, x.Comments[0].Body, `Example comment body 1`)
			assert.Contains(t, x.Comments[1].Body, `<img src="https://github.com/sample-user-2.png" width="35">`)
			assert.Contains(t, x.Comments[1].Body, `@sample-user-2 commented`)
			assert.Contains(t, x.Comments[1].Body, `Example comment body 2`)
			assert.Contains(t, x.Comments[1].Body, `Ref: http://localhost/example/target/issues/1.`)
		case 1:
			assert.Equal(t, path, "/repos/example/target/import/issues")
			assert.Equal(t, x.Issue.Title, "Example title 3")
			assert.Contains(t, x.Issue.Body, `<img src="https://github.com/sample-user-3.png" width="35">`)
			assert.Contains(t, x.Issue.Body, `Original pull request by @sample-user-3 - imported from <a href="http://localhost/example/source/pull/3">example/source#3</a>`)
			assert.Contains(t, x.Issue.Body, `Example body 3`)
			assert.Equal(t, x.Issue.Assignee, "")
			assert.Equal(t, x.Issue.Labels, []string{})
			assert.Len(t, x.Comments, 0)
		}
		importCount++
	}

	mig := New(source, target)
	assert.Nil(t, mig.Migrate())
}
