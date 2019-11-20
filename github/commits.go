package github

import (
	"fmt"
	"io"
)

// Commit represents a commit.
type Commit struct {
	SHA     string `json:"sha"`
	HTMLURL string `json:"html_url"`
	Commit  struct {
		Author    *CommitUser `json:"author"`
		Committer *CommitUser `json:"committer"`
		Message   string      `json:"message"`
	} `json:"commit"`
	Author    *User `json:"author"`
	Committer *User `json:"committer"`
	Parents   []struct {
		URL string `json:"url"`
		SHA string `json:"sha"`
	} `json:"parents"`
}

// CommitUser represents a commit user.
type CommitUser struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Date  string `json:"date"`
}

// Commits represents a collection of commits.
type Commits <-chan interface{}

// Next emits the next Commit.
func (cs Commits) Next() (*Commit, error) {
	for x := range cs {
		switch x := x.(type) {
		case error:
			return nil, x
		case *Commit:
			return x, nil
		}
		break
	}
	return nil, io.EOF
}

// CommitsFromSlice creates Commits from a slice.
func CommitsFromSlice(xs []*Commit) Commits {
	cs := make(chan interface{})
	go func() {
		defer close(cs)
		for _, p := range xs {
			cs <- p
		}
	}()
	return cs
}

// CommitsToSlice collects Commits.
func CommitsToSlice(cs Commits) ([]*Commit, error) {
	xs := []*Commit{}
	for {
		p, err := cs.Next()
		if err != nil {
			if err != io.EOF {
				return nil, err
			}
			return xs, nil
		}
		xs = append(xs, p)
	}
}

// ListPullReqCommits lists the commits of a pull request.
func (c *client) ListPullReqCommits(repo string, pullNumber int) Commits {
	cs := make(chan interface{})
	go func() {
		defer close(cs)
		path := c.url(fmt.Sprintf("/repos/%s/pulls/%d/commits?per_page=100", repo, pullNumber))
		for {
			var xs []*Commit
			next, err := c.getList(path, &xs)
			if err != nil {
				cs <- fmt.Errorf("ListPullReqCommits %s/pull/%d: %w", repo, pullNumber, err)
				break
			}
			for _, x := range xs {
				cs <- x
			}
			if next == "" {
				break
			}
			path = next
		}
	}()
	return Commits(cs)
}
