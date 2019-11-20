package migrator

import (
	"strings"

	"github.com/itchyny/github-migrator/github"
)

// returns the user, isMember (or owner) and error
func (m *migrator) lookupUser(name string) (*github.User, bool, error) {
	isOwner := strings.HasPrefix(m.targetRepo.FullName, name+"/")
	if u, ok := m.userByName[name]; ok {
		return u, isOwner, nil
	}
	if err, ok := m.errorUserByName[name]; ok {
		return nil, false, err
	}
	members, err := m.listTargetMembers()
	if err != nil {
		return nil, false, err
	}
	for _, mem := range members {
		if mem.Login == name {
			return mem.ToUser(), true, nil
		}
	}
	u, err := m.target.GetUser(name)
	if err != nil {
		if m.errorUserByName == nil {
			m.errorUserByName = make(map[string]error)
		}
		m.errorUserByName[name] = err
		return nil, false, err
	}
	if m.userByName == nil {
		m.userByName = make(map[string]*github.User)
	}
	m.userByName[name] = u
	return u, isOwner, nil
}
