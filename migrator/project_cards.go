package migrator

import (
	"fmt"
	"io"
	"strings"

	"github.com/itchyny/github-migrator/github"
)

func (m *migrator) migrateProjectCards() error {
	sourceProjects, err := github.ProjectsToSlice(m.source.ListProjects())
	if err != nil {
		if strings.Contains(err.Error(), "Projects are disabled for this repository") {
			return nil // do nothing
		}
		return err
	}
	if len(sourceProjects) == 0 {
		return nil
	}
	targetProjects, err := github.ProjectsToSlice(m.target.ListProjects())
	if err != nil {
		return err
	}
	for _, p := range sourceProjects {
		fmt.Printf("[=>] migrating cards in a project: %s\n", p.Name)
		q := lookupProject(targetProjects, p)
		if q == nil {
			return fmt.Errorf("project not found: %s", p.Name)
		}
		if err := m.migrateProjectCardsInProject(p.ID, q.ID); err != nil {
			return err
		}
	}
	return nil
}

func (m *migrator) migrateProjectCardsInProject(sourceID, targetID int) error {
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
		fmt.Printf("[=>] migrating cards in a project column: %s\n", c.Name)
		d := lookupProjectColumn(targetColumns, c)
		if d == nil {
			return fmt.Errorf("project card not found: %s", c.Name)
		}
		if err := m.migrateProjectCardsInColumn(c.ID, d.ID); err != nil {
			return err
		}
	}
}

func (m *migrator) migrateProjectCardsInColumn(sourceID, targetID int) error {
	sourceCards, err := github.ProjectCardsToSlice(
		m.source.ListProjectCards(sourceID),
	)
	if err != nil {
		return err
	}
	targetCards, err := github.ProjectCardsToSlice(
		m.target.ListProjectCards(targetID),
	)
	if err != nil {
		return err
	}
	reverseProjectCards(sourceCards)
	for _, c := range sourceCards {
		fmt.Printf("[=>] migrating a card: %s\n", m.getCardInfo(c))
		if lookupProjectCard(targetCards, c) != nil {
			fmt.Printf("[--] skipping: %s (already exists)\n", m.getCardInfo(c))
			continue
		}
		fmt.Printf("[>>] creating a new card: %s\n", m.getCardInfo(c))
		var params *github.CreateProjectCardParams
		if issueNumber := c.GetIssueNumber(); issueNumber > 0 {
			id, err := m.getTargetIssueID(issueNumber)
			if err != nil {
				return err
			}
			params = &github.CreateProjectCardParams{
				ContentID:   id,
				ContentType: github.ProjectCardContentTypeIssue,
			}
		} else {
			params = &github.CreateProjectCardParams{
				Note: c.Note,
			}
		}
		if _, err := m.target.CreateProjectCard(targetID, params); err != nil {
			return err
		}
	}
	return nil
}

func lookupProjectCard(cs []*github.ProjectCard, c *github.ProjectCard) *github.ProjectCard {
	for _, d := range cs {
		if c.Note != "" && c.Note == d.Note || c.GetIssueNumber() == d.GetIssueNumber() {
			return d
		}
	}
	return nil
}

func reverseProjectCards(cs []*github.ProjectCard) {
	for left, right := 0, len(cs)-1; left < right; left, right = left+1, right-1 {
		cs[left], cs[right] = cs[right], cs[left]
	}
}

func (m *migrator) getCardInfo(c *github.ProjectCard) string {
	if issueNumber := c.GetIssueNumber(); issueNumber > 0 {
		return fmt.Sprintf("%s/issues/%d", m.targetRepo.FullName, issueNumber)
	}
	xs := strings.Split(c.Note, "\n")
	if len(xs) > 0 {
		return xs[0]
	}
	return c.Note
}
