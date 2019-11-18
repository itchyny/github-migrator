package github

// MockClient represents a mock for GitHub client.
type MockClient struct {
	getUserCallback            func() (*User, error)
	listMembersCallback        func(string) Members
	getRepoCallback            func(string) (*Repo, error)
	updateRepoCallback         func(string, *UpdateRepoParams) (*Repo, error)
	listLabelsCallback         func(string) Labels
	createLabelCallback        func(string, *CreateLabelParams) (*Label, error)
	updateLabelCallback        func(string, string, *UpdateLabelParams) (*Label, error)
	listIssuesCallback         func(string, *ListIssuesParams) Issues
	getIssueCallback           func(string, int) (*Issue, error)
	listCommentsCallback       func(string, int) Comments
	listEventsCallback         func(string, int) Events
	listPullReqsCallback       func(string, *ListPullReqsParams) PullReqs
	getPullReqCallback         func(string, int) (*PullReq, error)
	listPullReqCommitsCallback func(string, int) Commits
	getDiffCallback            func(string, string) (string, error)
	getCompareCallback         func(string, string, string) (string, error)
	listReviewsCallback        func(string, int) Reviews
	listReviewCommentsCallback func(string, int) ReviewComments
	listProjectsCallback       func(string, *ListProjectsParams) Projects
	getProjectCallback         func(int) (*Project, error)
	importCallback             func(string, *Import) (*ImportResult, error)
	getImportCallback          func(string, int) (*ImportResult, error)
}

// MockClientOption is an option of mock client.
type MockClientOption func(*MockClient)

// NewMockClient creates a new MockClient.
func NewMockClient(opts ...MockClientOption) *MockClient {
	cli := &MockClient{}
	for _, opt := range opts {
		opt(cli)
	}
	return cli
}

// GetUser ...
func (c *MockClient) GetUser() (*User, error) {
	if c.getUserCallback != nil {
		return c.getUserCallback()
	}
	panic("MockClient#GetUser")
}

// MockGetUser ...
func MockGetUser(callback func() (*User, error)) MockClientOption {
	return func(c *MockClient) {
		c.getUserCallback = callback
	}
}

// ListMembers ...
func (c *MockClient) ListMembers(org string) Members {
	if c.listMembersCallback != nil {
		return c.listMembersCallback(listMembersPath(org))
	}
	panic("MockClient#ListMembers")
}

// MockListMembers ...
func MockListMembers(callback func(string) Members) MockClientOption {
	return func(c *MockClient) {
		c.listMembersCallback = callback
	}
}

// GetRepo ...
func (c *MockClient) GetRepo(repo string) (*Repo, error) {
	if c.getRepoCallback != nil {
		return c.getRepoCallback(getRepoPath(repo))
	}
	panic("MockClient#GetRepo")
}

// MockGetRepo ...
func MockGetRepo(callback func(string) (*Repo, error)) MockClientOption {
	return func(c *MockClient) {
		c.getRepoCallback = callback
	}
}

// UpdateRepo ...
func (c *MockClient) UpdateRepo(repo string, params *UpdateRepoParams) (*Repo, error) {
	if c.updateRepoCallback != nil {
		return c.updateRepoCallback(updateRepoPath(repo), params)
	}
	panic("MockClient#UpdateRepo")
}

// MockUpdateRepo ...
func MockUpdateRepo(callback func(string, *UpdateRepoParams) (*Repo, error)) MockClientOption {
	return func(c *MockClient) {
		c.updateRepoCallback = callback
	}
}

// ListLabels ...
func (c *MockClient) ListLabels(repo string) Labels {
	if c.listLabelsCallback != nil {
		return c.listLabelsCallback(listLabelsPath(repo))
	}
	panic("MockClient#ListLabels")
}

// MockListLabels ...
func MockListLabels(callback func(string) Labels) MockClientOption {
	return func(c *MockClient) {
		c.listLabelsCallback = callback
	}
}

// CreateLabel ...
func (c *MockClient) CreateLabel(repo string, params *CreateLabelParams) (*Label, error) {
	if c.createLabelCallback != nil {
		return c.createLabelCallback(createLabelsPath(repo), params)
	}
	panic("MockClient#CreateLabel")
}

// MockCreateLabel ...
func MockCreateLabel(callback func(string, *CreateLabelParams) (*Label, error)) MockClientOption {
	return func(c *MockClient) {
		c.createLabelCallback = callback
	}
}

// UpdateLabel ...
func (c *MockClient) UpdateLabel(repo, name string, params *UpdateLabelParams) (*Label, error) {
	if c.updateLabelCallback != nil {
		return c.updateLabelCallback(updateLabelsPath(repo, name), name, params)
	}
	panic("MockClient#UpdateLabel")
}

// MockUpdateLabel ...
func MockUpdateLabel(callback func(string, string, *UpdateLabelParams) (*Label, error)) MockClientOption {
	return func(c *MockClient) {
		c.updateLabelCallback = callback
	}
}

// ListIssues ...
func (c *MockClient) ListIssues(repo string, params *ListIssuesParams) Issues {
	if c.listIssuesCallback != nil {
		return c.listIssuesCallback(listIssuesPath(repo, params), params)
	}
	panic("MockClient#ListIssues")
}

// MockListIssues ...
func MockListIssues(callback func(string, *ListIssuesParams) Issues) MockClientOption {
	return func(c *MockClient) {
		c.listIssuesCallback = callback
	}
}

// GetIssue ...
func (c *MockClient) GetIssue(repo string, issueNumber int) (*Issue, error) {
	if c.getIssueCallback != nil {
		return c.getIssueCallback(getIssuePath(repo, issueNumber), issueNumber)
	}
	panic("MockClient#GetIssue")
}

// MockGetIssue ...
func MockGetIssue(callback func(string, int) (*Issue, error)) MockClientOption {
	return func(c *MockClient) {
		c.getIssueCallback = callback
	}
}

// ListComments ...
func (c *MockClient) ListComments(repo string, issueNumber int) Comments {
	if c.listCommentsCallback != nil {
		return c.listCommentsCallback(listCommentsPath(repo, issueNumber), issueNumber)
	}
	panic("MockClient#ListComments")
}

// MockListComments ...
func MockListComments(callback func(string, int) Comments) MockClientOption {
	return func(c *MockClient) {
		c.listCommentsCallback = callback
	}
}

// ListEvents ...
func (c *MockClient) ListEvents(repo string, issueNumber int) Events {
	if c.listEventsCallback != nil {
		return c.listEventsCallback(listEventsPath(repo, issueNumber), issueNumber)
	}
	panic("MockClient#ListEvents")
}

// MockListEvents ...
func MockListEvents(callback func(string, int) Events) MockClientOption {
	return func(c *MockClient) {
		c.listEventsCallback = callback
	}
}

// ListPullReqs ...
func (c *MockClient) ListPullReqs(repo string, params *ListPullReqsParams) PullReqs {
	if c.listPullReqsCallback != nil {
		return c.listPullReqsCallback(listPullReqsPath(repo, params), params)
	}
	panic("MockClient#ListPullReqs")
}

// MockListPullReqs ...
func MockListPullReqs(callback func(string, *ListPullReqsParams) PullReqs) MockClientOption {
	return func(c *MockClient) {
		c.listPullReqsCallback = callback
	}
}

// GetPullReq ...
func (c *MockClient) GetPullReq(repo string, pullNumber int) (*PullReq, error) {
	if c.getPullReqCallback != nil {
		return c.getPullReqCallback(getPullReqPath(repo, pullNumber), pullNumber)
	}
	panic("MockClient#GetPullReq")
}

// MockGetPullReq ...
func MockGetPullReq(callback func(string, int) (*PullReq, error)) MockClientOption {
	return func(c *MockClient) {
		c.getPullReqCallback = callback
	}
}

// ListPullReqCommits ...
func (c *MockClient) ListPullReqCommits(repo string, pullNumber int) Commits {
	if c.listPullReqCommitsCallback != nil {
		return c.listPullReqCommitsCallback(listPullReqCommitsPath(repo, pullNumber), pullNumber)
	}
	panic("MockClient#ListPullReqCommits")
}

// MockListPullReqCommits ...
func MockListPullReqCommits(callback func(string, int) Commits) MockClientOption {
	return func(c *MockClient) {
		c.listPullReqCommitsCallback = callback
	}
}

// GetDiff ...
func (c *MockClient) GetDiff(repo string, sha string) (string, error) {
	if c.getDiffCallback != nil {
		return c.getDiffCallback(getDiffPath(repo, sha), sha)
	}
	panic("MockClient#GetDiff")
}

// MockGetDiff ...
func MockGetDiff(callback func(string, string) (string, error)) MockClientOption {
	return func(c *MockClient) {
		c.getDiffCallback = callback
	}
}

// GetCompare ...
func (c *MockClient) GetCompare(repo string, base, head string) (string, error) {
	if c.getCompareCallback != nil {
		return c.getCompareCallback(getComparePath(repo, base, head), base, head)
	}
	panic("MockClient#GetCompare")
}

// MockGetCompare ...
func MockGetCompare(callback func(string, string, string) (string, error)) MockClientOption {
	return func(c *MockClient) {
		c.getCompareCallback = callback
	}
}

// ListReviews ...
func (c *MockClient) ListReviews(repo string, pullNumber int) Reviews {
	if c.listReviewsCallback != nil {
		return c.listReviewsCallback(listReviewsPath(repo, pullNumber), pullNumber)
	}
	panic("MockClient#ListReviews")
}

// MockListReviews ...
func MockListReviews(callback func(string, int) Reviews) MockClientOption {
	return func(c *MockClient) {
		c.listReviewsCallback = callback
	}
}

// ListReviewComments ...
func (c *MockClient) ListReviewComments(repo string, pullNumber int) ReviewComments {
	if c.listReviewCommentsCallback != nil {
		return c.listReviewCommentsCallback(listReviewCommentsPath(repo, pullNumber), pullNumber)
	}
	panic("MockClient#ListReviewComments")
}

// MockListReviewComments ...
func MockListReviewComments(callback func(string, int) ReviewComments) MockClientOption {
	return func(c *MockClient) {
		c.listReviewCommentsCallback = callback
	}
}

// ListProjects ...
func (c *MockClient) ListProjects(repo string, params *ListProjectsParams) Projects {
	if c.listProjectsCallback != nil {
		return c.listProjectsCallback(listProjectsPath(repo, params), params)
	}
	panic("MockClient#ListProjects")
}

// MockListProjects ...
func MockListProjects(callback func(string, *ListProjectsParams) Projects) MockClientOption {
	return func(c *MockClient) {
		c.listProjectsCallback = callback
	}
}

// GetProject ...
func (c *MockClient) GetProject(projectID int) (*Project, error) {
	if c.getProjectCallback != nil {
		return c.getProjectCallback(projectID)
	}
	panic("MockClient#GetProject")
}

// MockGetProject ...
func MockGetProject(callback func(int) (*Project, error)) MockClientOption {
	return func(c *MockClient) {
		c.getProjectCallback = callback
	}
}

// Import ...
func (c *MockClient) Import(repo string, issue *Import) (*ImportResult, error) {
	if c.importCallback != nil {
		return c.importCallback(issueImportPath(repo), issue)
	}
	panic("MockClient#Import")
}

// MockImport ...
func MockImport(callback func(string, *Import) (*ImportResult, error)) MockClientOption {
	return func(c *MockClient) {
		c.importCallback = callback
	}
}

// GetImport ...
func (c *MockClient) GetImport(repo string, id int) (*ImportResult, error) {
	if c.getImportCallback != nil {
		return c.getImportCallback(getImportPath(repo, id), id)
	}
	panic("MockClient#MockGetImport")
}

// MockGetImport ...
func MockGetImport(callback func(string, int) (*ImportResult, error)) MockClientOption {
	return func(c *MockClient) {
		c.getImportCallback = callback
	}
}
