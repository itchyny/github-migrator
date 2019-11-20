package migrator

import (
	"fmt"
	"io"

	"github.com/itchyny/github-migrator/github"
)

func (m *migrator) migrateProjects() error {
	sourceProjects := m.source.ListProjects()
	targetProjects, err := github.ProjectsToSlice(m.target.ListProjects())
	if err != nil {
		return err
	}
	for {
		p, err := sourceProjects.Next()
		if err != nil {
			if err != io.EOF {
				return err
			}
			return nil
		}
		fmt.Printf("[=>] migrating a project: %s\n", p.Name)
		q := lookupProject(targetProjects, p)
		if q == nil {
			fmt.Printf("[>>] creating a new project: %s\n", p.Name)
			if q, err = m.target.CreateProject(&github.CreateProjectParams{
				Name: p.Name, Body: p.Body,
			}); err != nil {
				return err
			}
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
}

func lookupProject(ps []*github.Project, p *github.Project) *github.Project {
	for _, q := range ps {
		if p.Name == q.Name {
			return q
		}
	}
	return nil
}

func (m *migrator) listTargetProjects() ([]*github.Project, error) {
	if m.projects != nil {
		return m.projects, nil
	}
	projects, err := github.ProjectsToSlice(m.target.ListProjects())
	if err != nil {
		return nil, err
	}
	m.projects = projects
	return projects, nil
}
