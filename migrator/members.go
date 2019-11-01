package migrator

import "github.com/itchyny/github-migrator/github"

func (m *migrator) listTargetMembers() ([]*github.Member, error) {
	if m.members != nil {
		return m.members, nil
	}
	members, err := github.MembersToSlice(m.target.ListMembers())
	if err != nil {
		return nil, err
	}
	m.members = members
	return members, nil
}

func (m *migrator) isTargetMember(name string) (bool, error) {
	members, err := m.listTargetMembers()
	if err != nil {
		return false, err
	}
	m.members = members
	for _, m := range members {
		if m.Login == name {
			return true, nil
		}
	}
	return false, nil
}
