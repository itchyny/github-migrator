package github

import "io"

// Commit represents a commit.
type Commit struct {
	URL     string `json:"url"`
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
