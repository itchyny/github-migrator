package migrator

import (
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
		q := lookupProject(targetProjects, p)
		if q == nil {
			if q, err = m.target.CreateProject(&github.CreateProjectParams{
				Name: p.Name, Body: p.Body,
			}); err != nil {
				return err
			}
		}
		if p.Body != q.Body || p.State != q.State {
			if q, err = m.target.UpdateProject(q.ID, &github.UpdateProjectParams{
				// Do not update name.
				Body:  p.Body,
				State: p.State,
			}); err != nil {
				return err
			}
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
