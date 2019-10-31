package github

// MockClient represents a mock for GitHub client.
type MockClient struct {
	getUserCallback      func() (*User, error)
	getRepoCallback      func(string) (*Repo, error)
	updateRepoCallback   func(string, *UpdateRepoParams) (*Repo, error)
	listLabelsCallback   func(string) Labels
	listIssuesCallback   func(string, *ListIssuesParams) Issues
	listCommentsCallback func(string, int) Comments
	listPullReqsCallback func(string, *ListPullReqsParams) PullReqs
	importCallback       func(string, *Import) error
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

// Import ...
func (c *MockClient) Import(repo string, issue *Import) error {
	if c.importCallback != nil {
		return c.importCallback(issueImportPath(repo), issue)
	}
	panic("MockClient#Import")
}

// MockImport ...
func MockImport(callback func(string, *Import) error) MockClientOption {
	return func(c *MockClient) {
		c.importCallback = callback
	}
}
