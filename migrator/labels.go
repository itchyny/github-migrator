package migrator

import (
	"fmt"
	"strings"

	"github.com/itchyny/github-migrator/github"
)

func (m *migrator) migrateLabels() error {
	sourceLabels, err := github.LabelsToSlice(m.source.ListLabels())
	if err != nil {
		return err
	}
	targetLabels, err := github.LabelsToSlice(m.target.ListLabels())
	if err != nil {
		return err
	}
	for _, sourceLabel := range sourceLabels {
		fmt.Printf("[=>] migrating a label: %s\n", sourceLabel.Name)
		var exists bool
		for _, targetLabel := range targetLabels {
			if strings.EqualFold(sourceLabel.Name, targetLabel.Name) {
				if sourceLabel.Description != targetLabel.Description ||
					sourceLabel.Color != targetLabel.Color {
					fmt.Printf("[|>] updating an existing label: %s\n", targetLabel.Name)
					if _, err := m.target.UpdateLabel(targetLabel.Name, &github.UpdateLabelParams{
						Name:        sourceLabel.Name,
						Description: sourceLabel.Description,
						Color:       sourceLabel.Color,
					}); err != nil {
						return err
					}
				} else {
					fmt.Printf("[--] skipping: %s (already exists)\n", sourceLabel.Name)
				}
				exists = true
				break
			}
		}
		if exists {
			continue
		}
		fmt.Printf("[>>] creating a new label: %s\n", sourceLabel.Name)
		if _, err := m.target.CreateLabel(&github.CreateLabelParams{
			Name:        sourceLabel.Name,
			Description: sourceLabel.Description,
			Color:       sourceLabel.Color,
		}); err != nil {
			return err
		}
	}
	return nil
}
