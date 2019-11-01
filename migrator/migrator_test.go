package migrator

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/itchyny/github-migrator/github"
	"github.com/itchyny/github-migrator/repo"
)

func init() {
	beforeImportIssueDuration = 0
	waitImportIssueInitialDuration = 0
}

func TestMigratorMigrate(t *testing.T) {
	source := repo.New(github.NewMockClient(
		github.MockGetRepo(func(path string) (*github.Repo, error) {
			return &github.Repo{
				Name:        "source",
				FullName:    "example/source",
				Description: "Source repository.",
				HTMLURL:     "http://localhost/example/source",
				Homepage:    "http://localhost/",
				Private:     false,
			}, nil
		}),
		github.MockListLabels(func(path string) github.Labels {
			return github.LabelsFromSlice([]*github.Label{
				{
					ID:          1,
					Name:        "bug",
					Description: "This is a bug.",
					Color:       "fc2929",
					Default:     false,
				},
				{
					ID:          2,
					Name:        "design",
					Description: "This is a design issue.",
					Color:       "bfdadc",
					Default:     false,
				},
			})
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
						Login: "bob",
					},
					Assignee: &github.User{
						Login: "bob",
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
		github.MockListReviewComments(func(path string, pullNumber int) github.ReviewComments {
			assert.Equal(t, path, "/repos/example/source/pulls/3/comments")
			assert.Equal(t, pullNumber, 3)
			return github.ReviewCommentsFromSlice([]*github.ReviewComment{
				{
					ID:       100,
					Path:     "sample.txt",
					DiffHunk: "@@ -0,0 +1 @@\n+foo",
					Body:     "Nice catch.",
					User: &github.User{
						Login: "sample-user-2",
					},
				},
				{
					ID:          200,
					Path:        "sample.txt",
					DiffHunk:    "@@ -0,0 +1 @@\n+foo",
					Body:        "@bob Thanks. bobb",
					InReplyToID: 100,
					User: &github.User{
						Login: "alice",
					},
				},
			})
		}),
	), "example/source")

	var assertImport func(string, *github.Import)
	target := repo.New(github.NewMockClient(
		github.MockListMembers(func(path string) github.Members {
			assert.Equal(t, path, "/orgs/example/members")
			return github.MembersFromSlice([]*github.Member{
				{
					Login: "sample-user-2",
				},
				{
					Login: "alice-1",
				},
			})
		}),
		github.MockGetRepo(func(path string) (*github.Repo, error) {
			return &github.Repo{
				Name:        "target",
				Description: "Target repository.",
				HTMLURL:     "http://localhost/example/target",
				Private:     true,
			}, nil
		}),
		github.MockUpdateRepo(func(path string, params *github.UpdateRepoParams) (*github.Repo, error) {
			assert.Equal(t, path, "/repos/example/target")
			assert.Equal(t, params.Name, "target")
			assert.Equal(t, params.Description, "Target repository.")
			assert.Equal(t, params.Homepage, "http://localhost/")
			assert.Equal(t, params.Private, true)
			return &github.Repo{}, nil
		}),
		github.MockListLabels(func(path string) github.Labels {
			return github.LabelsFromSlice([]*github.Label{
				{
					ID:          1,
					Name:        "bug",
					Description: "This is a bug.",
					Color:       "292929",
					Default:     false,
				},
			})
		}),
		github.MockCreateLabel(func(path string, params *github.CreateLabelParams) (*github.Label, error) {
			assert.Equal(t, path, "/repos/example/target/labels")
			assert.Equal(t, params.Name, "design")
			return nil, nil
		}),
		github.MockUpdateLabel(func(path, name string, params *github.UpdateLabelParams) (*github.Label, error) {
			assert.Equal(t, path, "/repos/example/target/labels/"+name)
			assert.Equal(t, params.Name, name)
			assert.Equal(t, params.Name, "bug")
			assert.Equal(t, params.Color, "fc2929")
			return nil, nil
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
		github.MockImport(func(path string, x *github.Import) (*github.ImportResult, error) {
			assertImport(path, x)
			return &github.ImportResult{
				ID:     12345,
				Status: "pending",
				URL:    "http://localhost/repo/example/target/import/issues/12345",
			}, nil
		}),
		github.MockGetImport(func(path string, id int) (*github.ImportResult, error) {
			assert.Equal(t, id, 12345)
			return &github.ImportResult{
				ID:     12345,
				Status: "imported",
				URL:    "http://localhost/repo/example/target/import/issues/12345",
			}, nil
		}),
	), "example/target")

	var importCount int
	assertImport = func(path string, x *github.Import) {
		switch importCount {
		case 0:
			assert.Equal(t, path, "/repos/example/target/import/issues")
			assert.Equal(t, x.Issue.Title, "Example title 2")
			assert.Equal(t, x.Issue.Body, `<table>
  <tr>
    <td>
      <img src="https://github.com/sample-user-2.png" width="35">
    </td>
    <td>
      @sample-user-2 created the original issue<br>imported from <a href="http://localhost/example/source/issues/2">example/source#2</a>
    </td>
  </tr>
</table>


Example body 2
See http://localhost/example/target/issues/1.`)
			assert.Equal(t, x.Issue.Assignee, "sample-user-2")
			assert.Equal(t, x.Issue.Labels, []string{"label1", "label2"})

			assert.Len(t, x.Comments, 2)
			assert.Equal(t, x.Comments[0].Body, `<table>
  <tr>
    <td>
      <img src="https://github.com/github.png" width="35">
    </td>
    <td>
      @sample-user-1 commented
    </td>
  </tr>
</table>


Example comment body 1`)
			assert.Equal(t, x.Comments[1].Body, `<table>
  <tr>
    <td>
      <img src="https://github.com/sample-user-2.png" width="35">
    </td>
    <td>
      @sample-user-2 commented
    </td>
  </tr>
</table>


Example comment body 2
Ref: http://localhost/example/target/issues/1.`)
		case 1:
			assert.Equal(t, path, "/repos/example/target/import/issues")
			assert.Equal(t, x.Issue.Title, "Example title 3")
			assert.Equal(t, x.Issue.Body, `<table>
  <tr>
    <td>
      <img src="https://github.com/github.png" width="35">
    </td>
    <td>
      @charlie created the original pull request<br>imported from <a href="http://localhost/example/source/pull/3">example/source#3</a>
    </td>
  </tr>
</table>


Example body 3`)
			assert.Equal(t, x.Issue.Assignee, "")
			assert.Equal(t, x.Issue.Labels, []string{})
			assert.Len(t, x.Comments, 1)
			assert.Equal(t, x.Comments[0].Body, "```"+`diff
# sample.txt
@@ -0,0 +1 @@
+foo
`+"```"+`

<table>
  <tr>
    <td>
      <img src="https://github.com/sample-user-2.png" width="35">
    </td>
    <td>
      @sample-user-2 commented
    </td>
  </tr>
</table>


Nice catch.

<table>
  <tr>
    <td>
      <img src="https://github.com/alice-1.png" width="35">
    </td>
    <td>
      @alice-1 commented
    </td>
  </tr>
</table>


@charlie Thanks. bobb`)
		}
		importCount++
	}

	mig := New(source, target, map[string]string{"bob": "charlie", "alice": "alice-1"})
	assert.Nil(t, mig.Migrate())
}
