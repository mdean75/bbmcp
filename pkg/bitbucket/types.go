package bitbucket

// Configuration for Bitbucket Server API
type Config struct {
	BaseURL           string
	Username          string
	Password          string // App password or personal access token
	Token             string
	DefaultProjectKey string
}

// Bitbucket API structures
type PullRequest struct {
	ID           int                    `json:"id"`
	Version      int                    `json:"version"`
	Title        string                 `json:"title"`
	Description  string                 `json:"description"`
	State        string                 `json:"state"`
	Open         bool                   `json:"open"`
	Closed       bool                   `json:"closed"`
	CreatedDate  int64                  `json:"createdDate"`
	UpdatedDate  int64                  `json:"updatedDate"`
	FromRef      PullRequestRef         `json:"fromRef"`
	ToRef        PullRequestRef         `json:"toRef"`
	Locked       bool                   `json:"locked"`
	Author       User                   `json:"author"`
	Reviewers    []Reviewer             `json:"reviewers"`
	Participants []Participant          `json:"participants"`
	Properties   map[string]interface{} `json:"properties"`
	Links        map[string]interface{} `json:"links"`
}

type PullRequestRef struct {
	ID           string     `json:"id"`
	DisplayID    string     `json:"displayId"`
	LatestCommit string     `json:"latestCommit"`
	Repository   Repository `json:"repository"`
}

type Repository struct {
	Slug          string                 `json:"slug"`
	ID            int                    `json:"id"`
	Name          string                 `json:"name"`
	ScmID         string                 `json:"scmId"`
	State         string                 `json:"state"`
	StatusMessage string                 `json:"statusMessage"`
	Forkable      bool                   `json:"forkable"`
	Project       Project                `json:"project"`
	Public        bool                   `json:"public"`
	Links         map[string]interface{} `json:"links"`
}

type Project struct {
	Key         string                 `json:"key"`
	ID          int                    `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Public      bool                   `json:"public"`
	Type        string                 `json:"type"`
	Links       map[string]interface{} `json:"links"`
}

type User struct {
	Name         string `json:"name"`
	EmailAddress string `json:"emailAddress"`
	ID           int    `json:"id"`
	DisplayName  string `json:"displayName"`
	Active       bool   `json:"active"`
	Slug         string `json:"slug"`
	Type         string `json:"type"`
}

type Reviewer struct {
	User               User   `json:"user"`
	Role               string `json:"role"`
	Approved           bool   `json:"approved"`
	Status             string `json:"status"`
	LastReviewedCommit string `json:"lastReviewedCommit"`
}

type Participant struct {
	User               User   `json:"user"`
	Role               string `json:"role"`
	Approved           bool   `json:"approved"`
	Status             string `json:"status"`
	LastReviewedCommit string `json:"lastReviewedCommit"`
}

type PullRequestActivity struct {
	Values     []Activity `json:"values"`
	Size       int        `json:"size"`
	Limit      int        `json:"limit"`
	IsLastPage bool       `json:"isLastPage"`
	Start      int        `json:"start"`
}

type Activity struct {
	ID               int         `json:"id"`
	CreatedDate      int64       `json:"createdDate"`
	User             User        `json:"user"`
	Action           string      `json:"action"`
	CommentAction    string      `json:"commentAction,omitempty"`
	Comment          *Comment    `json:"comment,omitempty"`
	FromHash         string      `json:"fromHash,omitempty"`
	PreviousFromHash string      `json:"previousFromHash,omitempty"`
	PreviousToHash   string      `json:"previousToHash,omitempty"`
	ToHash           string      `json:"toHash,omitempty"`
	Added            *CommitList `json:"added,omitempty"`
	Removed          *CommitList `json:"removed,omitempty"`
}

type Comment struct {
	Properties          map[string]interface{} `json:"properties"`
	ID                  int                    `json:"id"`
	Version             int                    `json:"version"`
	Text                string                 `json:"text"`
	Author              User                   `json:"author"`
	CreatedDate         int64                  `json:"createdDate"`
	UpdatedDate         int64                  `json:"updatedDate"`
	Comments            []Comment              `json:"comments"`
	Tasks               []Task                 `json:"tasks"`
	PermittedOperations interface{}            `json:"permittedOperations"`
}

type Task struct {
	Anchor              TaskAnchor  `json:"anchor"`
	Author              User        `json:"author"`
	CreatedDate         int64       `json:"createdDate"`
	ID                  int         `json:"id"`
	PermittedOperations interface{} `json:"permittedOperations"`
	State               string      `json:"state"`
	Text                string      `json:"text"`
}

type TaskAnchor struct {
	ID         int                    `json:"id"`
	Version    int                    `json:"version"`
	FileType   string                 `json:"fileType"`
	FromHash   string                 `json:"fromHash"`
	ToHash     string                 `json:"toHash"`
	Line       int                    `json:"line"`
	LineType   string                 `json:"lineType"`
	Path       string                 `json:"path"`
	Properties map[string]interface{} `json:"properties"`
}

type CommentAnchor struct {
	Line         int    `json:"line,omitempty"`
	LineType     string `json:"lineType,omitempty"`
	Path         string `json:"path,omitempty"`
	FileType     string `json:"fileType,omitempty"`
	FromHash     string `json:"fromHash,omitempty"`
	ToHash       string `json:"toHash,omitempty"`
	SrcPath      string `json:"srcPath,omitempty"`
	DstPath      string `json:"dstPath,omitempty"`
	DiffType     string `json:"diffType,omitempty"`
	OrphanedType string `json:"orphanedType,omitempty"`
}

type CommitList struct {
	Values     []Commit `json:"values"`
	Size       int      `json:"size"`
	Limit      int      `json:"limit"`
	IsLastPage bool     `json:"isLastPage"`
	Start      int      `json:"start"`
}

type Commit struct {
	ID                 string         `json:"id"`
	DisplayID          string         `json:"displayId"`
	Author             Person         `json:"author"`
	AuthorTimestamp    int64          `json:"authorTimestamp"`
	Committer          Person         `json:"committer"`
	CommitterTimestamp int64          `json:"committerTimestamp"`
	Message            string         `json:"message"`
	Parents            []CommitParent `json:"parents"`
}

type Person struct {
	Name         string `json:"name"`
	EmailAddress string `json:"emailAddress"`
}

type CommitParent struct {
	ID        string `json:"id"`
	DisplayID string `json:"displayId"`
}

type PullRequestSettings struct {
	MergeConfig              *MergeConfig `json:"mergeConfig,omitempty"`
	RequiredApprovers        int          `json:"requiredApprovers,omitempty"`
	RequiredAllApprovers     bool         `json:"requiredAllApprovers,omitempty"`
	RequiredAllTasksComplete bool         `json:"requiredAllTasksComplete,omitempty"`
	RequiredSuccessfulBuilds int          `json:"requiredSuccessfulBuilds,omitempty"`
	NeedsWork                bool         `json:"needsWork,omitempty"`
	// Plugin-specific settings with namespaced keys
	RequiredApproversHook *HookConfig `json:"com.atlassian.bitbucket.server.bitbucket-bundled-hooks:requiredApprovers,omitempty"`
	RequiredBuildsHook    *HookConfig `json:"com.atlassian.bitbucket.server.bitbucket-build:requiredBuilds,omitempty"`
}

type RequiredApprovalPath struct {
	ID                int     `json:"id"`
	Path              string  `json:"path"`
	RequiredApprovals int     `json:"requiredApprovals"`
	Users             []User  `json:"users"`
	Groups            []Group `json:"groups"`
}

type DefaultReviewer struct {
	ID                int        `json:"id"`
	SourceRefMatcher  RefMatcher `json:"sourceRefMatcher"`
	TargetRefMatcher  RefMatcher `json:"targetRefMatcher"`
	Reviewers         []User     `json:"reviewers"`
	RequiredApprovals int        `json:"requiredApprovals"`
}

type RefMatcher struct {
	ID        string `json:"id"`
	Type      string `json:"type"`
	Active    bool   `json:"active"`
	DisplayID string `json:"displayId"`
}

type MergeConfig struct {
	CommitSummaries int             `json:"commitSummaries"`
	DefaultStrategy MergeStrategy   `json:"defaultStrategy"`
	Strategies      []MergeStrategy `json:"strategies"`
	Type            string          `json:"type"`
}

type MergeStrategy struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Enabled     bool   `json:"enabled"`
	Flag        string `json:"flag"`
}

type HookConfig struct {
	Enable bool `json:"enable"`
	Count  int  `json:"count"`
}

type RequiredReviewerGroup struct {
	ID                int        `json:"id"`
	Name              string     `json:"name"`
	RequiredApprovals int        `json:"requiredApprovals"`
	SourceRefMatcher  RefMatcher `json:"sourceRefMatcher"`
	TargetRefMatcher  RefMatcher `json:"targetRefMatcher"`
}

type Group struct {
	Name string `json:"name"`
}

type BranchPermission struct {
	ID         int        `json:"id"`
	Type       string     `json:"type"`
	Matcher    RefMatcher `json:"matcher"`
	Users      []User     `json:"users"`
	Groups     []Group    `json:"groups"`
	AccessKeys []string   `json:"accessKeys"`
}
