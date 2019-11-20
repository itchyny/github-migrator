package migrator

import (
	"strings"

	"github.com/itchyny/github-migrator/github"
	"github.com/itchyny/github-migrator/repo"
)

// Migrator represents a GitHub migrator.
type Migrator interface {
	Migrate() error
}

// New creates a new Migrator.
func New(source, target repo.Repo, userMapping map[string]string) Migrator {
	return &migrator{source: source, target: target, userMapping: userMapping}
}

type migrator struct {
	source, target         repo.Repo
	userMapping            map[string]string
	sourceRepo, targetRepo *github.Repo
	commentFilters         commentFilters
	targetMembers          []*github.Member
	targetProjects         []*github.Project
	projectByIDs           map[int]*github.Project
	userByNames            map[string]*github.User
	errorUserByNames       map[string]error
}

// Migrate the repository.
func (m *migrator) Migrate() (err error) {
	if m.sourceRepo, err = m.source.Get(); err != nil {
		return err
	}
	if m.targetRepo, err = m.target.Get(); err != nil {
		return err
	}
	m.commentFilters = newCommentFilters(
		newRepoURLFilter(m.sourceRepo, m.targetRepo),
		newUserMappingFilter(m.userMapping),
	)
	if m.targetMembers, err = github.MembersToSlice(m.target.ListMembers()); err != nil {
		return err
	}
	if err = m.migrateRepo(); err != nil {
		return err
	}
	if err = m.migrateLabels(); err != nil {
		return err
	}
	// projects should be imported before issues
	if err = m.migrateProjects(); err != nil {
		return err
	}
	if projects, err := github.ProjectsToSlice(m.target.ListProjects()); err != nil {
		if !strings.Contains(err.Error(), "Projects are disabled for this repository") {
			return err
		}
		m.targetProjects = []*github.Project{}
	} else {
		m.targetProjects = projects
	}
	if err = m.migrateIssues(); err != nil {
		return err
	}
	if err = m.migrateHooks(); err != nil {
		return err
	}
	return nil
}
