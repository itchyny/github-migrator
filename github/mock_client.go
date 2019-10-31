package github

// MockClient represents a mock for GitHub client.
type MockClient struct {
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
