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
