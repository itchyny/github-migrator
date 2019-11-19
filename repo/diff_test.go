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
		github.MockGetDiff(func(string, string) (string, error) {
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
		github.MockGetCompare(func(string, string, string) (string, error) {
			return expected, nil
		}),
	), "example/test")
	got, err := repo.GetCompare("xxxyyy", "zzzwww")
	assert.Nil(t, err)
	assert.Equal(t, got, expected)
}
