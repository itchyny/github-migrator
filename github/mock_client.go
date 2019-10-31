package github

type MockClient struct {
}

type MockClientOption func(*MockClient)

func NewMockClient(opts ...MockClientOption) *MockClient {
	cli := &MockClient{}
	for _, opt := range opts {
		opt(cli)
	}
	return cli
}

func (c *MockClient) Login() (string, error) {
	return "mock", nil
}

func (c *MockClient) Hostname() string {
	return "api.github.com"
}
