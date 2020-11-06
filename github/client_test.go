package github

import "testing"

func TestNew(t *testing.T) {
	var _ Client = New("token", "http://localhost", "")
}
