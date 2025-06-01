package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

// MCP Protocol structures
type MCPRequest struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params,omitempty"`
}

type MCPResponse struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      interface{} `json:"id"`
	Result  interface{} `json:"result,omitempty"`
	Error   *MCPError   `json:"error,omitempty"`
}

type MCPError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type InitializeParams struct {
	ProtocolVersion string                 `json:"protocolVersion"`
	Capabilities    map[string]interface{} `json:"capabilities"`
	ClientInfo      map[string]interface{} `json:"clientInfo"`
}

type ServerCapabilities struct {
	Tools     map[string]interface{} `json:"tools,omitempty"`
	Resources map[string]interface{} `json:"resources,omitempty"`
}

type InitializeResult struct {
	ProtocolVersion string             `json:"protocolVersion"`
	Capabilities    ServerCapabilities `json:"capabilities"`
	ServerInfo      map[string]string  `json:"serverInfo"`
}

type Tool struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	InputSchema map[string]interface{} `json:"inputSchema"`
}

type ToolListResult struct {
	Tools []Tool `json:"tools"`
}

type CallToolParams struct {
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments,omitempty"`
}

type CallToolResult struct {
	Content []ToolContent `json:"content"`
}

type ToolContent struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// Bitbucket API structures
type BitbucketConfig struct {
	BaseURL  string
	Username string
	Password string // App password or personal access token
}

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
	PermittedOperations []string               `json:"permittedOperations"`
}

type Task struct {
	Anchor              TaskAnchor `json:"anchor"`
	Author              User       `json:"author"`
	CreatedDate         int64      `json:"createdDate"`
	ID                  int        `json:"id"`
	PermittedOperations []string   `json:"permittedOperations"`
	State               string     `json:"state"`
	Text                string     `json:"text"`
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

// BitbucketServer handles Bitbucket Server API operations
type BitbucketServer struct {
	config *BitbucketConfig
	client *http.Client
}

func NewBitbucketServer(config *BitbucketConfig) *BitbucketServer {
	return &BitbucketServer{
		config: config,
		client: &http.Client{},
	}
}

func (bs *BitbucketServer) makeRequest(method, endpoint string, body io.Reader) (*http.Response, error) {
	url := fmt.Sprintf("%s/rest/api/1.0%s", bs.config.BaseURL, endpoint)
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(bs.config.Username, bs.config.Password)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	return bs.client.Do(req)
}

func (bs *BitbucketServer) GetPullRequests(projectKey, repoSlug string, state string, limit int) ([]PullRequest, error) {
	endpoint := fmt.Sprintf("/projects/%s/repos/%s/pull-requests", projectKey, repoSlug)

	params := []string{}
	if state != "" {
		params = append(params, fmt.Sprintf("state=%s", state))
	}
	if limit > 0 {
		params = append(params, fmt.Sprintf("limit=%d", limit))
	}

	if len(params) > 0 {
		endpoint += "?" + strings.Join(params, "&")
	}

	resp, err := bs.makeRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d", resp.StatusCode)
	}

	var result struct {
		Values []PullRequest `json:"values"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Values, nil
}

func (bs *BitbucketServer) GetPullRequest(projectKey, repoSlug string, pullRequestID int) (*PullRequest, error) {
	endpoint := fmt.Sprintf("/projects/%s/repos/%s/pull-requests/%d", projectKey, repoSlug, pullRequestID)

	resp, err := bs.makeRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d", resp.StatusCode)
	}

	var pr PullRequest
	if err := json.NewDecoder(resp.Body).Decode(&pr); err != nil {
		return nil, err
	}

	return &pr, nil
}

func (bs *BitbucketServer) GetPullRequestActivity(projectKey, repoSlug string, pullRequestID int) (*PullRequestActivity, error) {
	endpoint := fmt.Sprintf("/projects/%s/repos/%s/pull-requests/%d/activities", projectKey, repoSlug, pullRequestID)

	resp, err := bs.makeRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d", resp.StatusCode)
	}

	var activity PullRequestActivity
	if err := json.NewDecoder(resp.Body).Decode(&activity); err != nil {
		return nil, err
	}

	return &activity, nil
}

func (bs *BitbucketServer) CreatePullRequest(projectKey, repoSlug string, pr *PullRequest) (*PullRequest, error) {
	endpoint := fmt.Sprintf("/projects/%s/repos/%s/pull-requests", projectKey, repoSlug)

	jsonData, err := json.Marshal(pr)
	if err != nil {
		return nil, err
	}

	resp, err := bs.makeRequest("POST", endpoint, strings.NewReader(string(jsonData)))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var createdPR PullRequest
	if err := json.NewDecoder(resp.Body).Decode(&createdPR); err != nil {
		return nil, err
	}

	return &createdPR, nil
}

func (bs *BitbucketServer) ApprovePullRequest(projectKey, repoSlug string, pullRequestID int) error {
	endpoint := fmt.Sprintf("/projects/%s/repos/%s/pull-requests/%d/approve", projectKey, repoSlug, pullRequestID)

	resp, err := bs.makeRequest("POST", endpoint, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	return nil
}

func (bs *BitbucketServer) UnapprovalPullRequest(projectKey, repoSlug string, pullRequestID int) error {
	endpoint := fmt.Sprintf("/projects/%s/repos/%s/pull-requests/%d/approve", projectKey, repoSlug, pullRequestID)

	resp, err := bs.makeRequest("DELETE", endpoint, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	return nil
}

func (bs *BitbucketServer) MergePullRequest(projectKey, repoSlug string, pullRequestID int, version int) (*PullRequest, error) {
	endpoint := fmt.Sprintf("/projects/%s/repos/%s/pull-requests/%d/merge?version=%d", projectKey, repoSlug, pullRequestID, version)

	resp, err := bs.makeRequest("POST", endpoint, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var mergedPR PullRequest
	if err := json.NewDecoder(resp.Body).Decode(&mergedPR); err != nil {
		return nil, err
	}

	return &mergedPR, nil
}

func (bs *BitbucketServer) DeclinePullRequest(projectKey, repoSlug string, pullRequestID int, version int) (*PullRequest, error) {
	endpoint := fmt.Sprintf("/projects/%s/repos/%s/pull-requests/%d/decline?version=%d", projectKey, repoSlug, pullRequestID, version)

	resp, err := bs.makeRequest("POST", endpoint, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var declinedPR PullRequest
	if err := json.NewDecoder(resp.Body).Decode(&declinedPR); err != nil {
		return nil, err
	}

	return &declinedPR, nil
}

func (bs *BitbucketServer) GetPullRequestDiff(projectKey, repoSlug string, pullRequestID int, contextLines int, whitespace string, since string, until string) (string, error) {
	endpoint := fmt.Sprintf("/projects/%s/repos/%s/pull-requests/%d/diff", projectKey, repoSlug, pullRequestID)

	// Add query parameters if provided
	params := []string{}
	if contextLines > 0 {
		params = append(params, fmt.Sprintf("contextLines=%d", contextLines))
	}
	if whitespace != "" {
		params = append(params, fmt.Sprintf("whitespace=%s", whitespace))
	}
	if since != "" {
		params = append(params, fmt.Sprintf("since=%s", since))
	}
	if until != "" {
		params = append(params, fmt.Sprintf("until=%s", until))
	}

	if len(params) > 0 {
		endpoint += "?" + strings.Join(params, "&")
	}

	resp, err := bs.makeRequest("GET", endpoint, nil)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	// Read the raw diff content as text
	diffBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(diffBytes), nil
}

func (bs *BitbucketServer) CreatePullRequestComment(projectKey, repoSlug string, pullRequestID int, text string, anchor *CommentAnchor) (*Comment, error) {
	endpoint := fmt.Sprintf("/projects/%s/repos/%s/pull-requests/%d/comments", projectKey, repoSlug, pullRequestID)

	// Create the comment request body
	commentRequest := map[string]interface{}{
		"text": text,
	}

	// Add anchor if provided (for inline comments)
	if anchor != nil {
		commentRequest["anchor"] = anchor
	}

	jsonData, err := json.Marshal(commentRequest)
	if err != nil {
		return nil, err
	}

	resp, err := bs.makeRequest("POST", endpoint, strings.NewReader(string(jsonData)))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var comment Comment
	if err := json.NewDecoder(resp.Body).Decode(&comment); err != nil {
		return nil, err
	}

	return &comment, nil
}

// MCPServer implements the Model Context Protocol
type MCPServer struct {
	bitbucket *BitbucketServer
}

func NewMCPServer() *MCPServer {
	config := &BitbucketConfig{
		BaseURL:  os.Getenv("BITBUCKET_BASE_URL"),
		Username: os.Getenv("BITBUCKET_USERNAME"),
		Password: os.Getenv("BITBUCKET_PASSWORD"),
	}

	if config.BaseURL == "" || config.Username == "" || config.Password == "" {
		log.Fatal("Missing required environment variables: BITBUCKET_BASE_URL, BITBUCKET_USERNAME, BITBUCKET_PASSWORD")
	}

	return &MCPServer{
		bitbucket: NewBitbucketServer(config),
	}
}

func (s *MCPServer) handleInitialize(params InitializeParams, id interface{}) *MCPResponse {
	result := InitializeResult{
		ProtocolVersion: "2024-11-05",
		Capabilities: ServerCapabilities{
			Tools: map[string]interface{}{
				"listChanged": true,
			},
		},
		ServerInfo: map[string]string{
			"name":    "bbcli",
			"version": "1.0.0",
		},
	}

	return &MCPResponse{
		JSONRPC: "2.0",
		ID:      id,
		Result:  result,
	}
}

func (s *MCPServer) handleCallTool(params CallToolParams, id interface{}) *MCPResponse {
	var response *MCPResponse
	switch params.Name {
	case "list_pull_requests":
		response = s.handleListPullRequests(params.Arguments)
	case "get_pull_request":
		response = s.handleGetPullRequest(params.Arguments)
	case "get_pull_request_activity":
		response = s.handleGetPullRequestActivity(params.Arguments)
	case "create_pull_request":
		response = s.handleCreatePullRequest(params.Arguments)
	case "approve_pull_request":
		response = s.handleApprovePullRequest(params.Arguments)
	case "unapprove_pull_request":
		response = s.handleUnapprovePullRequest(params.Arguments)
	case "merge_pull_request":
		response = s.handleMergePullRequest(params.Arguments)
	case "decline_pull_request":
		response = s.handleDeclinePullRequest(params.Arguments)
	case "hello_world":
		response = s.handleHelloWorld(params.Arguments)
	case "get_pull_request_diff":
		response = s.handleGetPullRequestDiff(params.Arguments)
	case "create_pull_request_comment":
		response = s.handleCreatePullRequestComment(params.Arguments)
	default:
		response = &MCPResponse{
			JSONRPC: "2.0",
			Error: &MCPError{
				Code:    -32601,
				Message: "Unknown tool",
			},
		}
	}
	response.ID = id
	return response
}

func (s *MCPServer) handleListPullRequests(args map[string]interface{}) *MCPResponse {
	projectKey, ok := args["project_key"].(string)
	if !ok {
		return s.errorResponse(-32602, "project_key is required and must be a string")
	}

	repoSlug, ok := args["repo_slug"].(string)
	if !ok {
		return s.errorResponse(-32602, "repo_slug is required and must be a string")
	}

	state := ""
	if stateVal, ok := args["state"].(string); ok {
		state = stateVal
	}

	limit := 25
	if limitVal, ok := args["limit"].(float64); ok {
		limit = int(limitVal)
	}

	prs, err := s.bitbucket.GetPullRequests(projectKey, repoSlug, state, limit)
	if err != nil {
		return s.errorResponse(-32000, fmt.Sprintf("Failed to get pull requests: %v", err))
	}

	content, err := json.MarshalIndent(prs, "", "  ")
	if err != nil {
		return s.errorResponse(-32000, fmt.Sprintf("Failed to marshal response: %v", err))
	}

	result := CallToolResult{
		Content: []ToolContent{
			{
				Type: "text",
				Text: string(content),
			},
		},
	}

	return &MCPResponse{
		JSONRPC: "2.0",
		Result:  result,
	}
}

func (s *MCPServer) handleGetPullRequest(args map[string]interface{}) *MCPResponse {
	projectKey, ok := args["project_key"].(string)
	if !ok {
		return s.errorResponse(-32602, "project_key is required and must be a string")
	}

	repoSlug, ok := args["repo_slug"].(string)
	if !ok {
		return s.errorResponse(-32602, "repo_slug is required and must be a string")
	}

	pullRequestID, ok := args["pull_request_id"].(float64)
	if !ok {
		return s.errorResponse(-32602, "pull_request_id is required and must be an integer")
	}

	pr, err := s.bitbucket.GetPullRequest(projectKey, repoSlug, int(pullRequestID))
	if err != nil {
		return s.errorResponse(-32000, fmt.Sprintf("Failed to get pull request: %v", err))
	}

	content, err := json.MarshalIndent(pr, "", "  ")
	if err != nil {
		return s.errorResponse(-32000, fmt.Sprintf("Failed to marshal response: %v", err))
	}

	result := CallToolResult{
		Content: []ToolContent{
			{
				Type: "text",
				Text: string(content),
			},
		},
	}

	return &MCPResponse{
		JSONRPC: "2.0",
		Result:  result,
	}
}

func (s *MCPServer) handleGetPullRequestActivity(args map[string]interface{}) *MCPResponse {
	projectKey, ok := args["project_key"].(string)
	if !ok {
		return s.errorResponse(-32602, "project_key is required and must be a string")
	}

	repoSlug, ok := args["repo_slug"].(string)
	if !ok {
		return s.errorResponse(-32602, "repo_slug is required and must be a string")
	}

	pullRequestID, ok := args["pull_request_id"].(float64)
	if !ok {
		return s.errorResponse(-32602, "pull_request_id is required and must be an integer")
	}

	activity, err := s.bitbucket.GetPullRequestActivity(projectKey, repoSlug, int(pullRequestID))
	if err != nil {
		return s.errorResponse(-32000, fmt.Sprintf("Failed to get pull request activity: %v", err))
	}

	content, err := json.MarshalIndent(activity, "", "  ")
	if err != nil {
		return s.errorResponse(-32000, fmt.Sprintf("Failed to marshal response: %v", err))
	}

	result := CallToolResult{
		Content: []ToolContent{
			{
				Type: "text",
				Text: string(content),
			},
		},
	}

	return &MCPResponse{
		JSONRPC: "2.0",
		Result:  result,
	}
}

func (s *MCPServer) handleCreatePullRequest(args map[string]interface{}) *MCPResponse {
	projectKey, ok := args["project_key"].(string)
	if !ok {
		return s.errorResponse(-32602, "project_key is required and must be a string")
	}

	repoSlug, ok := args["repo_slug"].(string)
	if !ok {
		return s.errorResponse(-32602, "repo_slug is required and must be a string")
	}

	title, ok := args["title"].(string)
	if !ok {
		return s.errorResponse(-32602, "title is required and must be a string")
	}

	fromBranch, ok := args["from_branch"].(string)
	if !ok {
		return s.errorResponse(-32602, "from_branch is required and must be a string")
	}

	toBranch, ok := args["to_branch"].(string)
	if !ok {
		return s.errorResponse(-32602, "to_branch is required and must be a string")
	}

	description := ""
	if desc, ok := args["description"].(string); ok {
		description = desc
	}

	// Create the pull request structure
	pr := &PullRequest{
		Title:       title,
		Description: description,
		FromRef: PullRequestRef{
			ID: fromBranch,
			Repository: Repository{
				Slug: repoSlug,
				Project: Project{
					Key: projectKey,
				},
			},
		},
		ToRef: PullRequestRef{
			ID: toBranch,
			Repository: Repository{
				Slug: repoSlug,
				Project: Project{
					Key: projectKey,
				},
			},
		},
	}

	// Add reviewers if provided
	if reviewersInterface, ok := args["reviewers"]; ok {
		if reviewersList, ok := reviewersInterface.([]interface{}); ok {
			for _, reviewerInterface := range reviewersList {
				if reviewer, ok := reviewerInterface.(string); ok {
					pr.Reviewers = append(pr.Reviewers, Reviewer{
						User: User{
							Name: reviewer,
						},
					})
				}
			}
		}
	}

	createdPR, err := s.bitbucket.CreatePullRequest(projectKey, repoSlug, pr)
	if err != nil {
		return s.errorResponse(-32000, fmt.Sprintf("Failed to create pull request: %v", err))
	}

	content, err := json.MarshalIndent(createdPR, "", "  ")
	if err != nil {
		return s.errorResponse(-32000, fmt.Sprintf("Failed to marshal response: %v", err))
	}

	result := CallToolResult{
		Content: []ToolContent{
			{
				Type: "text",
				Text: string(content),
			},
		},
	}

	return &MCPResponse{
		JSONRPC: "2.0",
		Result:  result,
	}
}

func (s *MCPServer) handleApprovePullRequest(args map[string]interface{}) *MCPResponse {
	projectKey, ok := args["project_key"].(string)
	if !ok {
		return s.errorResponse(-32602, "project_key is required and must be a string")
	}

	repoSlug, ok := args["repo_slug"].(string)
	if !ok {
		return s.errorResponse(-32602, "repo_slug is required and must be a string")
	}

	pullRequestID, ok := args["pull_request_id"].(float64)
	if !ok {
		return s.errorResponse(-32602, "pull_request_id is required and must be an integer")
	}

	err := s.bitbucket.ApprovePullRequest(projectKey, repoSlug, int(pullRequestID))
	if err != nil {
		return s.errorResponse(-32000, fmt.Sprintf("Failed to approve pull request: %v", err))
	}

	result := CallToolResult{
		Content: []ToolContent{
			{
				Type: "text",
				Text: "Pull request approved successfully",
			},
		},
	}

	return &MCPResponse{
		JSONRPC: "2.0",
		Result:  result,
	}
}

func (s *MCPServer) handleUnapprovePullRequest(args map[string]interface{}) *MCPResponse {
	projectKey, ok := args["project_key"].(string)
	if !ok {
		return s.errorResponse(-32602, "project_key is required and must be a string")
	}

	repoSlug, ok := args["repo_slug"].(string)
	if !ok {
		return s.errorResponse(-32602, "repo_slug is required and must be a string")
	}

	pullRequestID, ok := args["pull_request_id"].(float64)
	if !ok {
		return s.errorResponse(-32602, "pull_request_id is required and must be an integer")
	}

	err := s.bitbucket.UnapprovalPullRequest(projectKey, repoSlug, int(pullRequestID))
	if err != nil {
		return s.errorResponse(-32000, fmt.Sprintf("Failed to unapprove pull request: %v", err))
	}

	result := CallToolResult{
		Content: []ToolContent{
			{
				Type: "text",
				Text: "Pull request approval removed successfully",
			},
		},
	}

	return &MCPResponse{
		JSONRPC: "2.0",
		Result:  result,
	}
}

func (s *MCPServer) handleMergePullRequest(args map[string]interface{}) *MCPResponse {
	projectKey, ok := args["project_key"].(string)
	if !ok {
		return s.errorResponse(-32602, "project_key is required and must be a string")
	}

	repoSlug, ok := args["repo_slug"].(string)
	if !ok {
		return s.errorResponse(-32602, "repo_slug is required and must be a string")
	}

	pullRequestID, ok := args["pull_request_id"].(float64)
	if !ok {
		return s.errorResponse(-32602, "pull_request_id is required and must be an integer")
	}

	version, ok := args["version"].(float64)
	if !ok {
		return s.errorResponse(-32602, "version is required and must be an integer")
	}

	mergedPR, err := s.bitbucket.MergePullRequest(projectKey, repoSlug, int(pullRequestID), int(version))
	if err != nil {
		return s.errorResponse(-32000, fmt.Sprintf("Failed to merge pull request: %v", err))
	}

	content, err := json.MarshalIndent(mergedPR, "", "  ")
	if err != nil {
		return s.errorResponse(-32000, fmt.Sprintf("Failed to marshal response: %v", err))
	}

	result := CallToolResult{
		Content: []ToolContent{
			{
				Type: "text",
				Text: string(content),
			},
		},
	}

	return &MCPResponse{
		JSONRPC: "2.0",
		Result:  result,
	}
}

func (s *MCPServer) handleDeclinePullRequest(args map[string]interface{}) *MCPResponse {
	projectKey, ok := args["project_key"].(string)
	if !ok {
		return s.errorResponse(-32602, "project_key is required and must be a string")
	}

	repoSlug, ok := args["repo_slug"].(string)
	if !ok {
		return s.errorResponse(-32602, "repo_slug is required and must be a string")
	}

	pullRequestID, ok := args["pull_request_id"].(float64)
	if !ok {
		return s.errorResponse(-32602, "pull_request_id is required and must be an integer")
	}

	version, ok := args["version"].(float64)
	if !ok {
		return s.errorResponse(-32602, "version is required and must be an integer")
	}

	declinedPR, err := s.bitbucket.DeclinePullRequest(projectKey, repoSlug, int(pullRequestID), int(version))
	if err != nil {
		return s.errorResponse(-32000, fmt.Sprintf("Failed to decline pull request: %v", err))
	}

	content, err := json.MarshalIndent(declinedPR, "", "  ")
	if err != nil {
		return s.errorResponse(-32000, fmt.Sprintf("Failed to marshal response: %v", err))
	}

	result := CallToolResult{
		Content: []ToolContent{
			{
				Type: "text",
				Text: string(content),
			},
		},
	}

	return &MCPResponse{
		JSONRPC: "2.0",
		Result:  result,
	}
}

func (s *MCPServer) handleCreatePullRequestComment(args map[string]interface{}) *MCPResponse {
	projectKey, ok := args["project_key"].(string)
	if !ok {
		return s.errorResponse(-32602, "project_key is required and must be a string")
	}

	repoSlug, ok := args["repo_slug"].(string)
	if !ok {
		return s.errorResponse(-32602, "repo_slug is required and must be a string")
	}

	pullRequestID, ok := args["pull_request_id"].(float64)
	if !ok {
		return s.errorResponse(-32602, "pull_request_id is required and must be an integer")
	}

	text, ok := args["text"].(string)
	if !ok {
		return s.errorResponse(-32602, "text is required and must be a string")
	}

	// Optional anchor for inline comments
	var anchor *CommentAnchor
	if anchorData, ok := args["anchor"]; ok {
		if anchorMap, ok := anchorData.(map[string]interface{}); ok {
			anchor = &CommentAnchor{}

			if line, ok := anchorMap["line"].(float64); ok {
				anchor.Line = int(line)
			}
			if lineType, ok := anchorMap["line_type"].(string); ok {
				anchor.LineType = lineType
			}
			if path, ok := anchorMap["path"].(string); ok {
				anchor.Path = path
			}
			if fileType, ok := anchorMap["file_type"].(string); ok {
				anchor.FileType = fileType
			}
			if fromHash, ok := anchorMap["from_hash"].(string); ok {
				anchor.FromHash = fromHash
			}
			if toHash, ok := anchorMap["to_hash"].(string); ok {
				anchor.ToHash = toHash
			}
			if srcPath, ok := anchorMap["src_path"].(string); ok {
				anchor.SrcPath = srcPath
			}
			if dstPath, ok := anchorMap["dst_path"].(string); ok {
				anchor.DstPath = dstPath
			}
			if orphanedType, ok := anchorMap["orphaned_type"].(string); ok {
				anchor.OrphanedType = orphanedType
			}
		}
	}

	comment, err := s.bitbucket.CreatePullRequestComment(projectKey, repoSlug, int(pullRequestID), text, anchor)
	if err != nil {
		return s.errorResponse(-32000, fmt.Sprintf("Failed to create pull request comment: %v", err))
	}

	content, err := json.MarshalIndent(comment, "", "  ")
	if err != nil {
		return s.errorResponse(-32000, fmt.Sprintf("Failed to marshal response: %v", err))
	}

	result := CallToolResult{
		Content: []ToolContent{
			{
				Type: "text",
				Text: string(content),
			},
		},
	}

	return &MCPResponse{
		JSONRPC: "2.0",
		Result:  result,
	}
}

func (s *MCPServer) handleGetPullRequestDiff(args map[string]interface{}) *MCPResponse {
	projectKey, ok := args["project_key"].(string)
	if !ok {
		return s.errorResponse(-32602, "project_key is required and must be a string")
	}

	repoSlug, ok := args["repo_slug"].(string)
	if !ok {
		return s.errorResponse(-32602, "repo_slug is required and must be a string")
	}

	pullRequestID, ok := args["pull_request_id"].(float64)
	if !ok {
		return s.errorResponse(-32602, "pull_request_id is required and must be an integer")
	}

	// Optional parameters
	contextLines := 0
	if contextVal, ok := args["context_lines"].(float64); ok {
		contextLines = int(contextVal)
	}

	whitespace := ""
	if whitespaceVal, ok := args["whitespace"].(string); ok {
		whitespace = whitespaceVal
	}

	since := ""
	if sinceVal, ok := args["since"].(string); ok {
		since = sinceVal
	}

	until := ""
	if untilVal, ok := args["until"].(string); ok {
		until = untilVal
	}

	diff, err := s.bitbucket.GetPullRequestDiff(projectKey, repoSlug, int(pullRequestID), contextLines, whitespace, since, until)
	if err != nil {
		return s.errorResponse(-32000, fmt.Sprintf("Failed to get pull request diff: %v", err))
	}

	result := CallToolResult{
		Content: []ToolContent{
			{
				Type: "text",
				Text: diff,
			},
		},
	}

	return &MCPResponse{
		JSONRPC: "2.0",
		Result:  result,
	}
}

func (s *MCPServer) handleHelloWorld(args map[string]interface{}) *MCPResponse {
	name := "World"
	if nameVal, ok := args["name"].(string); ok && nameVal != "" {
		name = nameVal
	}

	result := CallToolResult{
		Content: []ToolContent{
			{
				Type: "text",
				Text: fmt.Sprintf("Hello, %s!", name),
			},
		},
	}

	return &MCPResponse{
		JSONRPC: "2.0",
		Result:  result,
	}
}

func (s *MCPServer) errorResponse(code int, message string) *MCPResponse {
	return &MCPResponse{
		JSONRPC: "2.0",
		Error: &MCPError{
			Code:    code,
			Message: message,
		},
	}
}

func (s *MCPServer) handleRequest(request *MCPRequest) *MCPResponse {
	switch request.Method {
	case "initialize":
		var params InitializeParams
		if request.Params != nil {
			if paramsBytes, err := json.Marshal(request.Params); err == nil {
				json.Unmarshal(paramsBytes, &params)
			}
		}
		return s.handleInitialize(params, request.ID)

	case "notifications/initialized":
		// No response for notifications
		return nil

	case "tools/list":
		return s.handleToolsList(request.ID)

	case "tools/call":
		var params CallToolParams
		if request.Params != nil {
			if paramsBytes, err := json.Marshal(request.Params); err == nil {
				if err := json.Unmarshal(paramsBytes, &params); err != nil {
					response := s.errorResponse(-32602, "Invalid parameters")
					response.ID = request.ID
					return response
				}
			}
		}
		return s.handleCallTool(params, request.ID)

	case "resources/list":
		// Not implemented but return proper error
		response := s.errorResponse(-32601, "Method not found")
		response.ID = request.ID
		return response

	case "prompts/list":
		// Not implemented but return proper error
		response := s.errorResponse(-32601, "Method not found")
		response.ID = request.ID
		return response

	default:
		response := s.errorResponse(-32601, "Method not found")
		response.ID = request.ID
		return response
	}
}

func (s *MCPServer) handleToolsList(id interface{}) *MCPResponse {
	tools := []Tool{
		{
			Name:        "list_pull_requests",
			Description: "List pull requests for a repository",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"project_key": map[string]interface{}{
						"type":        "string",
						"description": "The project key",
					},
					"repo_slug": map[string]interface{}{
						"type":        "string",
						"description": "The repository slug",
					},
					"state": map[string]interface{}{
						"type":        "string",
						"description": "Filter by state (OPEN, MERGED, DECLINED)",
						"enum":        []string{"OPEN", "MERGED", "DECLINED"},
					},
					"limit": map[string]interface{}{
						"type":        "integer",
						"description": "Maximum number of results to return",
						"minimum":     1,
						"maximum":     100,
					},
				},
				"required": []string{"project_key", "repo_slug"},
			},
		},
		{
			Name:        "get_pull_request",
			Description: "Get details of a specific pull request",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"project_key": map[string]interface{}{
						"type":        "string",
						"description": "The project key",
					},
					"repo_slug": map[string]interface{}{
						"type":        "string",
						"description": "The repository slug",
					},
					"pull_request_id": map[string]interface{}{
						"type":        "integer",
						"description": "The pull request ID",
					},
				},
				"required": []string{"project_key", "repo_slug", "pull_request_id"},
			},
		},
		{
			Name:        "get_pull_request_activity",
			Description: "Get activity (comments, approvals, etc.) for a pull request",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"project_key": map[string]interface{}{
						"type":        "string",
						"description": "The project key",
					},
					"repo_slug": map[string]interface{}{
						"type":        "string",
						"description": "The repository slug",
					},
					"pull_request_id": map[string]interface{}{
						"type":        "integer",
						"description": "The pull request ID",
					},
				},
				"required": []string{"project_key", "repo_slug", "pull_request_id"},
			},
		},
		{
			Name:        "create_pull_request",
			Description: "Create a new pull request",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"project_key": map[string]interface{}{
						"type":        "string",
						"description": "The project key",
					},
					"repo_slug": map[string]interface{}{
						"type":        "string",
						"description": "The repository slug",
					},
					"title": map[string]interface{}{
						"type":        "string",
						"description": "The pull request title",
					},
					"description": map[string]interface{}{
						"type":        "string",
						"description": "The pull request description",
					},
					"from_branch": map[string]interface{}{
						"type":        "string",
						"description": "Source branch name",
					},
					"to_branch": map[string]interface{}{
						"type":        "string",
						"description": "Target branch name",
					},
					"reviewers": map[string]interface{}{
						"type": "array",
						"items": map[string]interface{}{
							"type": "string",
						},
						"description": "List of reviewer usernames",
					},
				},
				"required": []string{"project_key", "repo_slug", "title", "from_branch", "to_branch"},
			},
		},
		{
			Name:        "approve_pull_request",
			Description: "Approve a pull request",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"project_key": map[string]interface{}{
						"type":        "string",
						"description": "The project key",
					},
					"repo_slug": map[string]interface{}{
						"type":        "string",
						"description": "The repository slug",
					},
					"pull_request_id": map[string]interface{}{
						"type":        "integer",
						"description": "The pull request ID",
					},
				},
				"required": []string{"project_key", "repo_slug", "pull_request_id"},
			},
		},
		{
			Name:        "unapprove_pull_request",
			Description: "Remove approval from a pull request",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"project_key": map[string]interface{}{
						"type":        "string",
						"description": "The project key",
					},
					"repo_slug": map[string]interface{}{
						"type":        "string",
						"description": "The repository slug",
					},
					"pull_request_id": map[string]interface{}{
						"type":        "integer",
						"description": "The pull request ID",
					},
				},
				"required": []string{"project_key", "repo_slug", "pull_request_id"},
			},
		},
		{
			Name:        "merge_pull_request",
			Description: "Merge a pull request",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"project_key": map[string]interface{}{
						"type":        "string",
						"description": "The project key",
					},
					"repo_slug": map[string]interface{}{
						"type":        "string",
						"description": "The repository slug",
					},
					"pull_request_id": map[string]interface{}{
						"type":        "integer",
						"description": "The pull request ID",
					},
					"version": map[string]interface{}{
						"type":        "integer",
						"description": "The pull request version for optimistic locking",
					},
				},
				"required": []string{"project_key", "repo_slug", "pull_request_id", "version"},
			},
		},
		{
			Name:        "decline_pull_request",
			Description: "Decline a pull request",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"project_key": map[string]interface{}{
						"type":        "string",
						"description": "The project key",
					},
					"repo_slug": map[string]interface{}{
						"type":        "string",
						"description": "The repository slug",
					},
					"pull_request_id": map[string]interface{}{
						"type":        "integer",
						"description": "The pull request ID",
					},
					"version": map[string]interface{}{
						"type":        "integer",
						"description": "The pull request version for optimistic locking",
					},
				},
				"required": []string{"project_key", "repo_slug", "pull_request_id", "version"},
			},
		},
		{
			Name:        "get_pull_request_diff",
			Description: "Get the raw diff for a pull request",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"project_key": map[string]interface{}{
						"type":        "string",
						"description": "The project key",
					},
					"repo_slug": map[string]interface{}{
						"type":        "string",
						"description": "The repository slug",
					},
					"pull_request_id": map[string]interface{}{
						"type":        "integer",
						"description": "The pull request ID",
					},
					"context_lines": map[string]interface{}{
						"type":        "integer",
						"description": "Number of context lines around changes (optional)",
						"minimum":     0,
					},
					"whitespace": map[string]interface{}{
						"type":        "string",
						"description": "Whitespace handling: 'ignore-all', 'ignore-space-at-eol', 'ignore-space-change', 'ignore-trailing-space' (optional)",
						"enum":        []string{"ignore-all", "ignore-space-at-eol", "ignore-space-change", "ignore-trailing-space"},
					},
					"since": map[string]interface{}{
						"type":        "string",
						"description": "Base commit hash to diff from (optional)",
					},
					"until": map[string]interface{}{
						"type":        "string",
						"description": "End commit hash to diff to (optional)",
					},
				},
				"required": []string{"project_key", "repo_slug", "pull_request_id"},
			},
		},
		{
			Name:        "create_pull_request_comment",
			Description: "Add a comment to a pull request",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"project_key": map[string]interface{}{
						"type":        "string",
						"description": "The project key",
					},
					"repo_slug": map[string]interface{}{
						"type":        "string",
						"description": "The repository slug",
					},
					"pull_request_id": map[string]interface{}{
						"type":        "integer",
						"description": "The pull request ID",
					},
					"text": map[string]interface{}{
						"type":        "string",
						"description": "The comment text",
					},
					"anchor": map[string]interface{}{
						"type":        "object",
						"description": "Optional anchor for inline comments",
						"properties": map[string]interface{}{
							"line": map[string]interface{}{
								"type":        "integer",
								"description": "Line number for inline comment",
							},
							"line_type": map[string]interface{}{
								"type":        "string",
								"description": "Line type (ADDED, REMOVED, CONTEXT)",
								"enum":        []string{"ADDED", "REMOVED", "CONTEXT"},
							},
							"path": map[string]interface{}{
								"type":        "string",
								"description": "File path for inline comment",
							},
							"file_type": map[string]interface{}{
								"type":        "string",
								"description": "File type (FROM, TO)",
								"enum":        []string{"FROM", "TO"},
							},
							"from_hash": map[string]interface{}{
								"type":        "string",
								"description": "Source commit hash",
							},
							"to_hash": map[string]interface{}{
								"type":        "string",
								"description": "Target commit hash",
							},
							"src_path": map[string]interface{}{
								"type":        "string",
								"description": "Source file path (for renames)",
							},
							"dst_path": map[string]interface{}{
								"type":        "string",
								"description": "Destination file path (for renames)",
							},
							"orphaned_type": map[string]interface{}{
								"type":        "string",
								"description": "Orphaned comment type",
							},
						},
					},
				},
				"required": []string{"project_key", "repo_slug", "pull_request_id", "text"},
			},
		},
		{
			Name:        "hello_world",
			Description: "Say hello to someone",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"name": map[string]interface{}{
						"type":        "string",
						"description": "Name of the person to greet",
					},
				},
				"required": []string{"name"},
			},
		},
	}

	result := ToolListResult{Tools: tools}
	return &MCPResponse{
		JSONRPC: "2.0",
		ID:      id,
		Result:  result,
	}
}

func main() {
	server := NewMCPServer()

	// Read from stdin and write to stdout (STDIO transport)
	decoder := json.NewDecoder(os.Stdin)
	encoder := json.NewEncoder(os.Stdout)

	for {
		var request MCPRequest
		if err := decoder.Decode(&request); err != nil {
			if err == io.EOF {
				break
			}
			log.Printf("Error decoding request: %v", err)
			continue
		}

		response := server.handleRequest(&request)
		// Only send response if it's not nil (notifications don't get responses)
		if response != nil {
			if err := encoder.Encode(response); err != nil {
				log.Printf("Error encoding response: %v", err)
			}
		}
	}
}
