package github

type Client interface {
}

func New(token string) *client {
	return &client{token: token}
}

type client struct {
	token string
}
