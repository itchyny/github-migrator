package migrator

import "github.com/itchyny/github-migrator/github"

func (m *migrator) lookupUser(name string) (*github.User, error) {
	if u, ok := m.userByName[name]; ok {
		return u, nil
	}
	if err, ok := m.errorUserByName[name]; ok {
		return nil, err
	}
	members, err := m.listTargetMembers()
	if err != nil {
		return nil, err
	}
	for _, mem := range members {
		if mem.Login == name {
			if m.userByName == nil {
				m.userByName = make(map[string]*github.User)
			}
			m.userByName[name] = mem.ToUser()
			return mem.ToUser(), nil
		}
	}
	u, err := m.target.GetUser(name)
	if err != nil {
		if m.errorUserByName == nil {
			m.errorUserByName = make(map[string]error)
		}
		m.errorUserByName[name] = err
		return nil, err
	}
	if m.userByName == nil {
		m.userByName = make(map[string]*github.User)
	}
	m.userByName[name] = u
	return u, nil
}
