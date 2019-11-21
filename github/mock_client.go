package github

// MockClient represents a mock for GitHub client.
type MockClient struct {
	getLoginCallback            func() (*User, error)
	listUsersCallback           func() Users
	getUserCallback             func(string) (*User, error)
	listMembersCallback         func(string) Members
	getRepoCallback             func(string) (*Repo, error)
	updateRepoCallback          func(string, *UpdateRepoParams) (*Repo, error)
	listLabelsCallback          func(string) Labels
	createLabelCallback         func(string, *CreateLabelParams) (*Label, error)
	updateLabelCallback         func(string, string, *UpdateLabelParams) (*Label, error)
	listIssuesCallback          func(string, *ListIssuesParams) Issues
	getIssueCallback            func(string, int) (*Issue, error)
	listCommentsCallback        func(string, int) Comments
	listEventsCallback          func(string, int) Events
	listPullReqsCallback        func(string, *ListPullReqsParams) PullReqs
	getPullReqCallback          func(string, int) (*PullReq, error)
	listPullReqCommitsCallback  func(string, int) Commits
	getDiffCallback             func(string, string) (string, error)
	getCompareCallback          func(string, string, string) (string, error)
	listReviewsCallback         func(string, int) Reviews
	getReviewCallback           func(string, int, int) (*Review, error)
	listReviewCommentsCallback  func(string, int) ReviewComments
	listProjectsCallback        func(string, *ListProjectsParams) Projects
	getProjectCallback          func(int) (*Project, error)
	createProjectCallback       func(string, *CreateProjectParams) (*Project, error)
	updateProjectCallback       func(int, *UpdateProjectParams) (*Project, error)
	listProjectColumnsCallback  func(int) ProjectColumns
	getProjectColumnCallback    func(int) (*ProjectColumn, error)
	createProjectColumnCallback func(int, string) (*ProjectColumn, error)
	updateProjectColumnCallback func(int, string) (*ProjectColumn, error)
	listProjectCardsCallback    func(int) ProjectCards
	getProjectCardCallback      func(int) (*ProjectCard, error)
	createProjectCardCallback   func(int, *CreateProjectCardParams) (*ProjectCard, error)
	updateProjectCardCallback   func(int, *UpdateProjectCardParams) (*ProjectCard, error)
	moveProjectCardCallback     func(int, *MoveProjectCardParams) (*ProjectCard, error)
	listHooksCallback           func(string) Hooks
	getHookCallback             func(string, int) (*Hook, error)
	createHookCallback          func(string, *CreateHookParams) (*Hook, error)
	updateHookCallback          func(string, int, *UpdateHookParams) (*Hook, error)
	importCallback              func(string, *Import) (*ImportResult, error)
	getImportCallback           func(string, int) (*ImportResult, error)
}

// MockClientOption is an option of mock client.
type MockClientOption func(*MockClient)

// NewMockClient creates a new MockClient.
func NewMockClient(opts ...MockClientOption) *MockClient {
	cli := &MockClient{}
	for _, opt := range opts {
		opt(cli)
	}
	return cli
}

// GetLogin ...
func (c *MockClient) GetLogin() (*User, error) {
	if c.getLoginCallback != nil {
		return c.getLoginCallback()
	}
	panic("MockClient#GetLogin")
}

// MockGetLogin ...
func MockGetLogin(callback func() (*User, error)) MockClientOption {
	return func(c *MockClient) {
		c.getLoginCallback = callback
	}
}

// ListUsers ...
func (c *MockClient) ListUsers() Users {
	if c.listUsersCallback != nil {
		return c.listUsersCallback()
	}
	panic("MockClient#ListUsers")
}

// MockListUsers ...
func MockListUsers(callback func() Users) MockClientOption {
	return func(c *MockClient) {
		c.listUsersCallback = callback
	}
}

// GetUser ...
func (c *MockClient) GetUser(name string) (*User, error) {
	if c.getUserCallback != nil {
		return c.getUserCallback(name)
	}
	panic("MockClient#GetUser")
}

// MockGetUser ...
func MockGetUser(callback func(string) (*User, error)) MockClientOption {
	return func(c *MockClient) {
		c.getUserCallback = callback
	}
}

// ListMembers ...
func (c *MockClient) ListMembers(org string) Members {
	if c.listMembersCallback != nil {
		return c.listMembersCallback(org)
	}
	panic("MockClient#ListMembers")
}

// MockListMembers ...
func MockListMembers(callback func(string) Members) MockClientOption {
	return func(c *MockClient) {
		c.listMembersCallback = callback
	}
}

// GetRepo ...
func (c *MockClient) GetRepo(repo string) (*Repo, error) {
	if c.getRepoCallback != nil {
		return c.getRepoCallback(repo)
	}
	panic("MockClient#GetRepo")
}

// MockGetRepo ...
func MockGetRepo(callback func(string) (*Repo, error)) MockClientOption {
	return func(c *MockClient) {
		c.getRepoCallback = callback
	}
}

// UpdateRepo ...
func (c *MockClient) UpdateRepo(repo string, params *UpdateRepoParams) (*Repo, error) {
	if c.updateRepoCallback != nil {
		return c.updateRepoCallback(repo, params)
	}
	panic("MockClient#UpdateRepo")
}

// MockUpdateRepo ...
func MockUpdateRepo(callback func(string, *UpdateRepoParams) (*Repo, error)) MockClientOption {
	return func(c *MockClient) {
		c.updateRepoCallback = callback
	}
}

// ListLabels ...
func (c *MockClient) ListLabels(repo string) Labels {
	if c.listLabelsCallback != nil {
		return c.listLabelsCallback(repo)
	}
	panic("MockClient#ListLabels")
}

// MockListLabels ...
func MockListLabels(callback func(string) Labels) MockClientOption {
	return func(c *MockClient) {
		c.listLabelsCallback = callback
	}
}

// CreateLabel ...
func (c *MockClient) CreateLabel(repo string, params *CreateLabelParams) (*Label, error) {
	if c.createLabelCallback != nil {
		return c.createLabelCallback(repo, params)
	}
	panic("MockClient#CreateLabel")
}

// MockCreateLabel ...
func MockCreateLabel(callback func(string, *CreateLabelParams) (*Label, error)) MockClientOption {
	return func(c *MockClient) {
		c.createLabelCallback = callback
	}
}

// UpdateLabel ...
func (c *MockClient) UpdateLabel(repo, name string, params *UpdateLabelParams) (*Label, error) {
	if c.updateLabelCallback != nil {
		return c.updateLabelCallback(repo, name, params)
	}
	panic("MockClient#UpdateLabel")
}

// MockUpdateLabel ...
func MockUpdateLabel(callback func(string, string, *UpdateLabelParams) (*Label, error)) MockClientOption {
	return func(c *MockClient) {
		c.updateLabelCallback = callback
	}
}

// ListIssues ...
func (c *MockClient) ListIssues(repo string, params *ListIssuesParams) Issues {
	if c.listIssuesCallback != nil {
		return c.listIssuesCallback(repo, params)
	}
	panic("MockClient#ListIssues")
}

// MockListIssues ...
func MockListIssues(callback func(string, *ListIssuesParams) Issues) MockClientOption {
	return func(c *MockClient) {
		c.listIssuesCallback = callback
	}
}

// GetIssue ...
func (c *MockClient) GetIssue(repo string, issueNumber int) (*Issue, error) {
	if c.getIssueCallback != nil {
		return c.getIssueCallback(repo, issueNumber)
	}
	panic("MockClient#GetIssue")
}

// MockGetIssue ...
func MockGetIssue(callback func(string, int) (*Issue, error)) MockClientOption {
	return func(c *MockClient) {
		c.getIssueCallback = callback
	}
}

// ListComments ...
func (c *MockClient) ListComments(repo string, issueNumber int) Comments {
	if c.listCommentsCallback != nil {
		return c.listCommentsCallback(repo, issueNumber)
	}
	panic("MockClient#ListComments")
}

// MockListComments ...
func MockListComments(callback func(string, int) Comments) MockClientOption {
	return func(c *MockClient) {
		c.listCommentsCallback = callback
	}
}

// ListEvents ...
func (c *MockClient) ListEvents(repo string, issueNumber int) Events {
	if c.listEventsCallback != nil {
		return c.listEventsCallback(repo, issueNumber)
	}
	panic("MockClient#ListEvents")
}

// MockListEvents ...
func MockListEvents(callback func(string, int) Events) MockClientOption {
	return func(c *MockClient) {
		c.listEventsCallback = callback
	}
}

// ListPullReqs ...
func (c *MockClient) ListPullReqs(repo string, params *ListPullReqsParams) PullReqs {
	if c.listPullReqsCallback != nil {
		return c.listPullReqsCallback(repo, params)
	}
	panic("MockClient#ListPullReqs")
}

// MockListPullReqs ...
func MockListPullReqs(callback func(string, *ListPullReqsParams) PullReqs) MockClientOption {
	return func(c *MockClient) {
		c.listPullReqsCallback = callback
	}
}

// GetPullReq ...
func (c *MockClient) GetPullReq(repo string, pullNumber int) (*PullReq, error) {
	if c.getPullReqCallback != nil {
		return c.getPullReqCallback(repo, pullNumber)
	}
	panic("MockClient#GetPullReq")
}

// MockGetPullReq ...
func MockGetPullReq(callback func(string, int) (*PullReq, error)) MockClientOption {
	return func(c *MockClient) {
		c.getPullReqCallback = callback
	}
}

// ListPullReqCommits ...
func (c *MockClient) ListPullReqCommits(repo string, pullNumber int) Commits {
	if c.listPullReqCommitsCallback != nil {
		return c.listPullReqCommitsCallback(repo, pullNumber)
	}
	panic("MockClient#ListPullReqCommits")
}

// MockListPullReqCommits ...
func MockListPullReqCommits(callback func(string, int) Commits) MockClientOption {
	return func(c *MockClient) {
		c.listPullReqCommitsCallback = callback
	}
}

// GetDiff ...
func (c *MockClient) GetDiff(repo string, sha string) (string, error) {
	if c.getDiffCallback != nil {
		return c.getDiffCallback(repo, sha)
	}
	panic("MockClient#GetDiff")
}

// MockGetDiff ...
func MockGetDiff(callback func(string, string) (string, error)) MockClientOption {
	return func(c *MockClient) {
		c.getDiffCallback = callback
	}
}

// GetCompare ...
func (c *MockClient) GetCompare(repo string, base, head string) (string, error) {
	if c.getCompareCallback != nil {
		return c.getCompareCallback(repo, base, head)
	}
	panic("MockClient#GetCompare")
}

// MockGetCompare ...
func MockGetCompare(callback func(string, string, string) (string, error)) MockClientOption {
	return func(c *MockClient) {
		c.getCompareCallback = callback
	}
}

// ListReviews ...
func (c *MockClient) ListReviews(repo string, pullNumber int) Reviews {
	if c.listReviewsCallback != nil {
		return c.listReviewsCallback(repo, pullNumber)
	}
	panic("MockClient#ListReviews")
}

// MockListReviews ...
func MockListReviews(callback func(string, int) Reviews) MockClientOption {
	return func(c *MockClient) {
		c.listReviewsCallback = callback
	}
}

// GetReview ...
func (c *MockClient) GetReview(repo string, pullNumber, reviewID int) (*Review, error) {
	if c.getReviewCallback != nil {
		return c.getReviewCallback(repo, pullNumber, reviewID)
	}
	panic("MockClient#GetReview")
}

// MockGetReview ...
func MockGetReview(callback func(string, int, int) (*Review, error)) MockClientOption {
	return func(c *MockClient) {
		c.getReviewCallback = callback
	}
}

// ListReviewComments ...
func (c *MockClient) ListReviewComments(repo string, pullNumber int) ReviewComments {
	if c.listReviewCommentsCallback != nil {
		return c.listReviewCommentsCallback(repo, pullNumber)
	}
	panic("MockClient#ListReviewComments")
}

// MockListReviewComments ...
func MockListReviewComments(callback func(string, int) ReviewComments) MockClientOption {
	return func(c *MockClient) {
		c.listReviewCommentsCallback = callback
	}
}

// ListProjects ...
func (c *MockClient) ListProjects(repo string, params *ListProjectsParams) Projects {
	if c.listProjectsCallback != nil {
		return c.listProjectsCallback(repo, params)
	}
	panic("MockClient#ListProjects")
}

// MockListProjects ...
func MockListProjects(callback func(string, *ListProjectsParams) Projects) MockClientOption {
	return func(c *MockClient) {
		c.listProjectsCallback = callback
	}
}

// GetProject ...
func (c *MockClient) GetProject(projectID int) (*Project, error) {
	if c.getProjectCallback != nil {
		return c.getProjectCallback(projectID)
	}
	panic("MockClient#GetProject")
}

// MockGetProject ...
func MockGetProject(callback func(int) (*Project, error)) MockClientOption {
	return func(c *MockClient) {
		c.getProjectCallback = callback
	}
}

// CreateProject ...
func (c *MockClient) CreateProject(repo string, params *CreateProjectParams) (*Project, error) {
	if c.createProjectCallback != nil {
		return c.createProjectCallback(repo, params)
	}
	panic("MockClient#CreateProject")
}

// MockCreateProject ...
func MockCreateProject(callback func(string, *CreateProjectParams) (*Project, error)) MockClientOption {
	return func(c *MockClient) {
		c.createProjectCallback = callback
	}
}

// UpdateProject ...
func (c *MockClient) UpdateProject(projectID int, params *UpdateProjectParams) (*Project, error) {
	if c.updateProjectCallback != nil {
		return c.updateProjectCallback(projectID, params)
	}
	panic("MockClient#UpdateProject")
}

// MockUpdateProject ...
func MockUpdateProject(callback func(int, *UpdateProjectParams) (*Project, error)) MockClientOption {
	return func(c *MockClient) {
		c.updateProjectCallback = callback
	}
}

// ListProjectColumns ...
func (c *MockClient) ListProjectColumns(projectID int) ProjectColumns {
	if c.listProjectColumnsCallback != nil {
		return c.listProjectColumnsCallback(projectID)
	}
	panic("MockClient#ListProjectColumns")
}

// MockListProjectColumns ...
func MockListProjectColumns(callback func(int) ProjectColumns) MockClientOption {
	return func(c *MockClient) {
		c.listProjectColumnsCallback = callback
	}
}

// GetProjectColumn ...
func (c *MockClient) GetProjectColumn(projectColumnID int) (*ProjectColumn, error) {
	if c.getProjectColumnCallback != nil {
		return c.getProjectColumnCallback(projectColumnID)
	}
	panic("MockClient#GetProjectColumn")
}

// MockGetProjectColumn ...
func MockGetProjectColumn(callback func(int) (*ProjectColumn, error)) MockClientOption {
	return func(c *MockClient) {
		c.getProjectColumnCallback = callback
	}
}

// CreateProjectColumn ...
func (c *MockClient) CreateProjectColumn(projectID int, name string) (*ProjectColumn, error) {
	if c.createProjectColumnCallback != nil {
		return c.createProjectColumnCallback(projectID, name)
	}
	panic("MockClient#CreateProjectColumn")
}

// MockCreateProjectColumn ...
func MockCreateProjectColumn(callback func(int, string) (*ProjectColumn, error)) MockClientOption {
	return func(c *MockClient) {
		c.createProjectColumnCallback = callback
	}
}

// UpdateProjectColumn ...
func (c *MockClient) UpdateProjectColumn(projectColumnID int, name string) (*ProjectColumn, error) {
	if c.updateProjectColumnCallback != nil {
		return c.updateProjectColumnCallback(projectColumnID, name)
	}
	panic("MockClient#UpdateProjectColumn")
}

// MockUpdateProjectColumn ...
func MockUpdateProjectColumn(callback func(int, string) (*ProjectColumn, error)) MockClientOption {
	return func(c *MockClient) {
		c.updateProjectColumnCallback = callback
	}
}

// ListProjectCards ...
func (c *MockClient) ListProjectCards(columnID int) ProjectCards {
	if c.listProjectCardsCallback != nil {
		return c.listProjectCardsCallback(columnID)
	}
	panic("MockClient#ListProjectCards")
}

// MockListProjectCards ...
func MockListProjectCards(callback func(int) ProjectCards) MockClientOption {
	return func(c *MockClient) {
		c.listProjectCardsCallback = callback
	}
}

// GetProjectCard ...
func (c *MockClient) GetProjectCard(projectCardID int) (*ProjectCard, error) {
	if c.getProjectCardCallback != nil {
		return c.getProjectCardCallback(projectCardID)
	}
	panic("MockClient#GetProjectCard")
}

// MockGetProjectCard ...
func MockGetProjectCard(callback func(int) (*ProjectCard, error)) MockClientOption {
	return func(c *MockClient) {
		c.getProjectCardCallback = callback
	}
}

// CreateProjectCard ...
func (c *MockClient) CreateProjectCard(columnID int, params *CreateProjectCardParams) (*ProjectCard, error) {
	if c.createProjectCardCallback != nil {
		return c.createProjectCardCallback(columnID, params)
	}
	panic("MockClient#CreateProjectCard")
}

// MockCreateProjectCard ...
func MockCreateProjectCard(callback func(int, *CreateProjectCardParams) (*ProjectCard, error)) MockClientOption {
	return func(c *MockClient) {
		c.createProjectCardCallback = callback
	}
}

// UpdateProjectCard ...
func (c *MockClient) UpdateProjectCard(projectCardID int, params *UpdateProjectCardParams) (*ProjectCard, error) {
	if c.updateProjectCardCallback != nil {
		return c.updateProjectCardCallback(projectCardID, params)
	}
	panic("MockClient#UpdateProjectCard")
}

// MockUpdateProjectCard ...
func MockUpdateProjectCard(callback func(int, *UpdateProjectCardParams) (*ProjectCard, error)) MockClientOption {
	return func(c *MockClient) {
		c.updateProjectCardCallback = callback
	}
}

// MoveProjectCard ...
func (c *MockClient) MoveProjectCard(projectCardID int, params *MoveProjectCardParams) (*ProjectCard, error) {
	if c.moveProjectCardCallback != nil {
		return c.moveProjectCardCallback(projectCardID, params)
	}
	panic("MockClient#MoveProjectCard")
}

// MockMoveProjectCard ...
func MockMoveProjectCard(callback func(int, *MoveProjectCardParams) (*ProjectCard, error)) MockClientOption {
	return func(c *MockClient) {
		c.moveProjectCardCallback = callback
	}
}

// ListHooks ...
func (c *MockClient) ListHooks(repo string) Hooks {
	if c.listHooksCallback != nil {
		return c.listHooksCallback(repo)
	}
	panic("MockClient#ListHooks")
}

// MockListHooks ...
func MockListHooks(callback func(string) Hooks) MockClientOption {
	return func(c *MockClient) {
		c.listHooksCallback = callback
	}
}

// GetHook ...
func (c *MockClient) GetHook(repo string, hookID int) (*Hook, error) {
	if c.getHookCallback != nil {
		return c.getHookCallback(repo, hookID)
	}
	panic("MockClient#GetHook")
}

// MockGetHook ...
func MockGetHook(callback func(string, int) (*Hook, error)) MockClientOption {
	return func(c *MockClient) {
		c.getHookCallback = callback
	}
}

// CreateHook ...
func (c *MockClient) CreateHook(repo string, params *CreateHookParams) (*Hook, error) {
	if c.createHookCallback != nil {
		return c.createHookCallback(repo, params)
	}
	panic("MockClient#CreateHook")
}

// MockCreateHook ...
func MockCreateHook(callback func(string, *CreateHookParams) (*Hook, error)) MockClientOption {
	return func(c *MockClient) {
		c.createHookCallback = callback
	}
}

// UpdateHook ...
func (c *MockClient) UpdateHook(repo string, hookID int, params *UpdateHookParams) (*Hook, error) {
	if c.updateHookCallback != nil {
		return c.updateHookCallback(repo, hookID, params)
	}
	panic("MockClient#UpdateHook")
}

// MockUpdateHook ...
func MockUpdateHook(callback func(string, int, *UpdateHookParams) (*Hook, error)) MockClientOption {
	return func(c *MockClient) {
		c.updateHookCallback = callback
	}
}

// Import ...
func (c *MockClient) Import(repo string, issue *Import) (*ImportResult, error) {
	if c.importCallback != nil {
		return c.importCallback(repo, issue)
	}
	panic("MockClient#Import")
}

// MockImport ...
func MockImport(callback func(string, *Import) (*ImportResult, error)) MockClientOption {
	return func(c *MockClient) {
		c.importCallback = callback
	}
}

// GetImport ...
func (c *MockClient) GetImport(repo string, id int) (*ImportResult, error) {
	if c.getImportCallback != nil {
		return c.getImportCallback(repo, id)
	}
	panic("MockClient#MockGetImport")
}

// MockGetImport ...
func MockGetImport(callback func(string, int) (*ImportResult, error)) MockClientOption {
	return func(c *MockClient) {
		c.getImportCallback = callback
	}
}
