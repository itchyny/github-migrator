package github

import "net/url"

type path struct {
	path   string
	params url.Values
}

func newPath(x string) path {
	return path{path: x, params: url.Values{}}
}

func (p path) query(key, value string) path {
	if value != "" {
		p.params.Add(key, value)
	}
	return p
}

func (p path) String() string {
	if len(p.params) == 0 {
		return p.path
	}
	return p.path + "?" + p.params.Encode()
}
