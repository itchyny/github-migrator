package migrator

import (
	"fmt"
	"io"
	"time"

	"github.com/itchyny/github-migrator/github"
)

var waitProjectColumnDuration = 100 * time.Millisecond

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
		fmt.Printf("[=>] migrating a project column: %s\n", c.Name)
		d := lookupProjectColumn(targetColumns, c)
		if d == nil {
			fmt.Printf("[>>] creating a new project column: %s\n", c.Name)
			if d, err = m.target.CreateProjectColumn(targetID, c.Name); err != nil {
				return err
			}
		}
		time.Sleep(waitProjectColumnDuration)
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
