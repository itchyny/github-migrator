package migrator

import (
	"io"

	"github.com/itchyny/github-migrator/github"
)

func (m *migrator) migrateProjectColumns(sourceID, targetID int) error {
	sourceColumns := m.source.ListProjectColumns(sourceID)
	targetColumns, err := github.ProjectColumnsToSlice(
		m.target.ListProjectColumns(targetID),
	)
	if err != nil {
		return err
	}
	for {
		c, err := sourceColumns.Next()
		if err != nil {
			if err != io.EOF {
				return err
			}
			return nil
		}
		d := lookupProjectColumn(targetColumns, c)
		if d == nil {
			if d, err = m.target.CreateProjectColumn(targetID, c.Name); err != nil {
				return err
			}
		}
	}
}

func lookupProjectColumn(ps []*github.ProjectColumn, c *github.ProjectColumn) *github.ProjectColumn {
	for _, d := range ps {
		if c.Name == d.Name {
			return d
		}
	}
	return nil
}
