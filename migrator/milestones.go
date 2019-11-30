package migrator

import (
	"fmt"

	"github.com/itchyny/github-migrator/github"
)

func (m *migrator) migrateMilestones() error {
	sourceMilestones, err := github.MilestonesToSlice(
		m.source.ListMilestones(&github.ListMilestonesParams{
			State: github.ListMilestonesParamStateAll,
		}),
	)
	if err != nil {
		return err
	}
	targetMilestones, err := github.MilestonesToSlice(
		m.target.ListMilestones(&github.ListMilestonesParams{
			State: github.ListMilestonesParamStateAll,
		}),
	)
	if err != nil {
		return err
	}
	var largestMilestoneNumber int
	for _, l := range targetMilestones {
		if largestMilestoneNumber < l.Number {
			largestMilestoneNumber = l.Number
		}
	}
	for _, l := range sourceMilestones {
		fmt.Printf("[=>] migrating a milestone: %s\n", l.Title)
		for l.Number > largestMilestoneNumber+1 {
			n, err := m.target.CreateMilestone(&github.CreateMilestoneParams{
				Title: "[Deleted milestone]",
				State: github.MilestoneStateClosed,
			})
			if err != nil {
				return err
			}
			largestMilestoneNumber = n.Number
			if err := m.target.DeleteMilestone(n.Number); err != nil {
				return err
			}
		}
		n := lookupMilestone(targetMilestones, l)
		if n == nil {
			fmt.Printf("[>>] creating a new milestone: %s\n", l.Title)
			if n, err = m.target.CreateMilestone(&github.CreateMilestoneParams{
				Title: l.Title, Description: l.Description,
				State: l.State, DueOn: l.DueOn,
			}); err != nil {
				return err
			}
			largestMilestoneNumber = n.Number
		}
		if l.Description != n.Description || l.State != n.State || l.DueOn != n.DueOn {
			fmt.Printf("[|>] updating an existing milestone: %s\n", l.Title)
			if n, err = m.target.UpdateMilestone(n.Number, &github.UpdateMilestoneParams{
				Title:       l.Title,
				Description: l.Description,
				State:       l.State,
				DueOn:       l.DueOn,
			}); err != nil {
				return err
			}
		}
	}
	targetMilestones, err = github.MilestonesToSlice(
		m.target.ListMilestones(&github.ListMilestonesParams{
			State: github.ListMilestonesParamStateAll,
		}),
	)
	if err != nil {
		return err
	}
	m.milestoneByTitle = make(map[string]*github.Milestone)
	for _, l := range targetMilestones {
		m.milestoneByTitle[l.Title] = l
	}
	return nil
}

func lookupMilestone(ps []*github.Milestone, l *github.Milestone) *github.Milestone {
	for _, n := range ps {
		if l.Title == n.Title {
			return n
		}
	}
	return nil
}
