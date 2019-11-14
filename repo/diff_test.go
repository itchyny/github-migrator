package repo

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/itchyny/github-migrator/github"
)

func TestRepoGetDiff(t *testing.T) {
	expected := `diff --git a/README.md b/README.md
index 1234567..89abcde 100644
--- a/README.md
+++ b/README.md
@@ -1,6 +1,16 @@
 # README
-deleted
+added
`
	repo := New(github.NewMockClient(
		github.MockGetDiff(func(path string, sha string) (string, error) {
			assert.Contains(t, path, "/repos/example/test/commits/"+sha)
			return expected, nil
		}),
	), "example/test")
	got, err := repo.GetDiff("xxxyyy")
	assert.Nil(t, err)
	assert.Equal(t, got, expected)
}

func TestRepoGetCompare(t *testing.T) {
	expected := `diff --git a/README.md b/README.md
index 1234567..89abcde 100644
--- a/README.md
+++ b/README.md
@@ -1,6 +1,16 @@
 # README
-deleted
+added
`
	repo := New(github.NewMockClient(
		github.MockGetCompare(func(path string, base, head string) (string, error) {
			assert.Contains(t, path, "/repos/example/test/compare/"+base+"..."+head)
			return expected, nil
		}),
	), "example/test")
	got, err := repo.GetCompare("xxxyyy", "zzzwww")
	assert.Nil(t, err)
	assert.Equal(t, got, expected)
}
