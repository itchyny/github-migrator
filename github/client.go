package github

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/tomnomnom/linkheader"
)

// Client represents a GitHub client.
type Client interface {
	GetLogin() (*User, error)
	ListUsers() Users
	GetUser(string) (*User, error)
	ListMembers(string) Members
	GetRepo(string) (*Repo, error)
	UpdateRepo(string, *UpdateRepoParams) (*Repo, error)
	ListLabels(string) Labels
	CreateLabel(string, *CreateLabelParams) (*Label, error)
	UpdateLabel(string, string, *UpdateLabelParams) (*Label, error)
	ListIssues(string, *ListIssuesParams) Issues
	GetIssue(string, int) (*Issue, error)
	AddAssignees(string, int, []string) error
	ListComments(string, int) Comments
	ListEvents(string, int) Events
	ListPullReqs(string, *ListPullReqsParams) PullReqs
	GetPullReq(string, int) (*PullReq, error)
	ListPullReqCommits(string, int) Commits
	GetDiff(string, string) (string, error)
	GetCompare(string, string, string) (string, error)
	ListReviews(string, int) Reviews
	GetReview(string, int, int) (*Review, error)
	ListReviewComments(string, int) ReviewComments
	ListProjects(string, *ListProjectsParams) Projects
	GetProject(int) (*Project, error)
	CreateProject(string, *CreateProjectParams) (*Project, error)
	UpdateProject(int, *UpdateProjectParams) (*Project, error)
	DeleteProject(int) error
	ListProjectColumns(int) ProjectColumns
	GetProjectColumn(int) (*ProjectColumn, error)
	CreateProjectColumn(int, string) (*ProjectColumn, error)
	UpdateProjectColumn(int, string) (*ProjectColumn, error)
	ListProjectCards(int) ProjectCards
	GetProjectCard(int) (*ProjectCard, error)
	CreateProjectCard(int, *CreateProjectCardParams) (*ProjectCard, error)
	UpdateProjectCard(int, *UpdateProjectCardParams) (*ProjectCard, error)
	MoveProjectCard(int, *MoveProjectCardParams) (*ProjectCard, error)
	ListMilestones(string, *ListMilestonesParams) Milestones
	GetMilestone(string, int) (*Milestone, error)
	CreateMilestone(string, *CreateMilestoneParams) (*Milestone, error)
	UpdateMilestone(string, int, *UpdateMilestoneParams) (*Milestone, error)
	DeleteMilestone(string, int) error
	ListHooks(string) Hooks
	GetHook(string, int) (*Hook, error)
	CreateHook(string, *CreateHookParams) (*Hook, error)
	UpdateHook(string, int, *UpdateHookParams) (*Hook, error)
	Import(string, *Import) (*ImportResult, error)
	GetImport(string, int) (*ImportResult, error)
}

// New creates a new GitHub client.
func New(token, endpoint string, proxy string, opts ...ClientOption) Client {
	cli := &http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: endpoint != "https://api.github.com",
		},
	}}
	if proxy != "" {
		proxyURL, _ := url.Parse(proxy)
		cli = &http.Client{Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		}}
	}
	c := &client{token, endpoint, cli, &Logger{}}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

// ClientOption is an option of  client.
type ClientOption func(*client)

// ClientLogger returns a client option to set the logger.
func ClientLogger(l *Logger) ClientOption {
	return func(c *client) {
		c.logger = l
	}
}

type client struct {
	token, endpoint string
	client          *http.Client
	logger          *Logger
}

func (c *client) url(path string) string {
	return c.endpoint + path
}

func (c *client) do(method, path string, body interface{}) (*http.Response, error) {
	var retryCnt int
	duration := time.Minute
	for {
		res, retry, err := c.doOnce(method, path, body)
		if err == nil || !retry || retryCnt >= 7 {
			return res, err
		}
		retryCnt++
		if retryCnt > 2 {
			duration *= 2
			if duration > 10*time.Minute {
				duration = 10 * time.Minute
			}
		}
		time.Sleep(duration)
	}
}

func (c *client) doOnce(method, path string, body interface{}) (*http.Response, bool, error) {
	var b io.Reader
	if body != nil {
		bs, err := json.Marshal(body)
		if err != nil {
			return nil, false, err
		}
		b = bytes.NewReader(bs)
	}
	req, err := c.request(method, path, b)
	if err != nil {
		return nil, false, err
	}
	return c.doReq(req)
}

func (c *client) request(method, path string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, path, body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "token "+c.token)
	req.Header.Add("Accept", "application/vnd.github.golden-comet-preview+json")
	req.Header.Add("Accept", "application/vnd.github.symmetra-preview+json")
	req.Header.Add("Accept", "application/vnd.github.comfort-fade-preview+json")
	req.Header.Add("Accept", "application/vnd.github.sailor-v-preview+json")
	req.Header.Add("Accept", "application/vnd.github.starfox-preview+json")
	req.Header.Add("Accept", "application/vnd.github.inertia-preview+json")
	req.Header.Add("User-Agent", "github-migrator")
	return req, nil
}

func (c *client) doReq(req *http.Request) (*http.Response, bool, error) {
	c.logger.preRequest(req)
	res, err := c.client.Do(req)
	c.logger.postRequest(res, err)
	if err != nil {
		return nil, true, err
	}
	if res.StatusCode < 200 || 400 <= res.StatusCode {
		return nil, 500 <= res.StatusCode, getError(res)
	}
	return res, false, nil
}

func getError(res *http.Response) error {
	defer res.Body.Close()
	var r struct {
		Message string    `json:"message"`
		Errors  apiErrors `json:"errors"`
	}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return err
	}
	if len(r.Errors) == 0 {
		return errors.New(r.Message)
	}
	return fmt.Errorf("%s: %w", r.Message, r.Errors)
}

func (c *client) get(path string, v interface{}) error {
	res, err := c.do("GET", path, nil)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if err := json.NewDecoder(res.Body).Decode(&v); err != nil {
		return err
	}
	return nil
}

func (c *client) post(path string, body, v interface{}) error {
	res, err := c.do("POST", path, body)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if err := json.NewDecoder(res.Body).Decode(&v); err != nil {
		return err
	}
	return nil
}

func (c *client) patch(path string, body, v interface{}) error {
	res, err := c.do("PATCH", path, body)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if err := json.NewDecoder(res.Body).Decode(&v); err != nil {
		return err
	}
	return nil
}

func (c *client) delete(path string) error {
	res, err := c.do("DELETE", path, nil)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	return nil
}

func (c *client) getList(path string, v interface{}) (string, error) {
	res, err := c.do("GET", path, nil)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	if err := json.NewDecoder(res.Body).Decode(&v); err != nil {
		return "", err
	}
	return getNext(res.Header), nil
}

func getNext(header http.Header) string {
	xs := header["Link"]
	if len(xs) == 0 {
		return ""
	}
	links := linkheader.Parse(xs[0])
	for _, link := range links {
		if link.Rel == "next" {
			return link.URL
		}
	}
	return ""
}
