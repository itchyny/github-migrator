package github

// MockClient represents a mock for GitHub client.
type MockClient struct {
	listIssuesCallback func(string, *ListIssuesParams) Issues
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

// Login ...
func (c *MockClient) Login() (string, error) {
	return "mock", nil
}

// Hostname ...
func (c *MockClient) Hostname() string {
	return "api.github.com"
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
