package migrator

import (
	"strings"

	"github.com/itchyny/github-migrator/github"
)

func (m *migrator) isTargetMember(name string) (bool, error) {
	if strings.HasPrefix(m.targetRepo.FullName, name+"/") {
		return true, nil
	}
	for _, member := range m.targetMembers {
		if member.Login == name {
			return true, nil
		}
	}
	return false, nil
}

func (m *migrator) lookupUser(name string) (*github.User, error) {
	if u, ok := m.userByNames[name]; ok {
		return u, nil
	}
	if err, ok := m.errorUserByNames[name]; ok {
		return nil, err
	}
	for _, member := range m.targetMembers {
		if member.Login == name {
			return member.ToUser(), nil
		}
	}
	u, err := m.target.GetUser(name)
	if err != nil {
		if m.errorUserByNames == nil {
			m.errorUserByNames = make(map[string]error)
		}
		m.errorUserByNames[name] = err
		return nil, err
	}
	if m.userByNames == nil {
		m.userByNames = make(map[string]*github.User)
	}
	m.userByNames[name] = u
	return u, nil
}
