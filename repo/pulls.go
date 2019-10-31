package repo

import "github.com/itchyny/github-migrator/github"

// ListPullReqs lists the pull requests.
func (r *repo) ListPullReqs() github.PullReqs {
	return r.cli.ListPullReqs(r.path, &github.ListPullReqsParams{
		State:     github.ListPullReqsParamStateAll,
		Direction: github.ListPullReqsParamDirectionAsc,
	})
}
