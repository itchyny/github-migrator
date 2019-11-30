package migrator

import (
	"fmt"
	"strings"

	"github.com/itchyny/github-migrator/github"
)

func (m *migrator) migrateProjects() error {
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
	var largestProjectNumber int
	for _, l := range targetProjects {
		if largestProjectNumber < l.Number {
			largestProjectNumber = l.Number
		}
	}
	for _, p := range sourceProjects {
		fmt.Printf("[=>] migrating a project: %s\n", p.Name)
		for p.Number > largestProjectNumber+1 {
			q, err := m.target.CreateProject(&github.CreateProjectParams{
				Name: "[Deleted project]",
			})
			if err != nil {
				return err
			}
			largestProjectNumber = q.Number
			if err := m.target.DeleteProject(q.ID); err != nil {
				return err
			}
		}
		q := lookupProject(targetProjects, p)
		if q == nil {
			fmt.Printf("[>>] creating a new project: %s\n", p.Name)
			if q, err = m.target.CreateProject(&github.CreateProjectParams{
				Name: p.Name, Body: p.Body,
			}); err != nil {
				return err
			}
			largestProjectNumber = q.Number
		}
		if p.Body != q.Body || p.State != q.State {
			fmt.Printf("[|>] updating an existing project: %s\n", p.Name)
			if q, err = m.target.UpdateProject(q.ID, &github.UpdateProjectParams{
				// Do not update name.
				Body:  p.Body,
				State: p.State,
			}); err != nil {
				return err
			}
		}
		if err := m.migrateProjectColumns(p.ID, q.ID); err != nil {
			return err
		}
	}
	return nil
}

func (m *migrator) getProject(id int) (*github.Project, error) {
	if p, ok := m.projectByIDs[id]; ok {
		return p, nil
	}
	p, err := m.source.GetProject(id)
	if err != nil {
		return nil, err
	}
	if m.projectByIDs == nil {
		m.projectByIDs = make(map[int]*github.Project)
	}
	m.projectByIDs[id] = p
	return p, nil
}

func lookupProject(ps []*github.Project, p *github.Project) *github.Project {
	for _, q := range ps {
		if p.Name == q.Name {
			return q
		}
	}
	return nil
}
