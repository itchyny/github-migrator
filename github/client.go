package github

type Client interface {
}

func New(token, root string) *client {
	return &client{token: token, root: root}
}

type client struct {
	token, root string
}
