package migrator

import "github.com/itchyny/github-migrator/github"

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
		var exists bool
		for _, targetLabel := range targetLabels {
			if sourceLabel.Name == targetLabel.Name {
				if sourceLabel.Description != targetLabel.Description ||
					sourceLabel.Color != targetLabel.Color {
					if _, err := m.target.UpdateLabel(sourceLabel.Name, &github.UpdateLabelParams{
						Name:        sourceLabel.Name,
						Description: sourceLabel.Description,
						Color:       sourceLabel.Color,
					}); err != nil {
						return err
					}
				}
				exists = true
				break
			}
		}
		if exists {
			continue
		}
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
