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

func (s *MCPServer) handleInitialize(params InitializeParams) *MCPResponse {
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
		Result:  result,
	}
}

func (s *MCPServer) handleCallTool(params CallToolParams) *MCPResponse {
	switch params.Name {
	case "list_pull_requests":
		return s.handleListPullRequests(params.Arguments)
	case "get_pull_request":
		return s.handleGetPullRequest(params.Arguments)
	case "get_pull_request_activity":
		return s.handleGetPullRequestActivity(params.Arguments)
	case "create_pull_request":
		return s.handleCreatePullRequest(params.Arguments)
	case "approve_pull_request":
		return s.handleApprovePullRequest(params.Arguments)
	case "unapprove_pull_request":
		return s.handleUnapprovePullRequest(params.Arguments)
	case "merge_pull_request":
		return s.handleMergePullRequest(params.Arguments)
	case "decline_pull_request":
		return s.handleDeclinePullRequest(params.Arguments)
	default:
		return &MCPResponse{
			JSONRPC: "2.0",
			Error: &MCPError{
				Code:    -32601,
				Message: "Unknown tool",
			},
		}
	}
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
		return s.handleInitialize(params)

	case "initialized":
		// No-op for initialized notification
		return &MCPResponse{
			JSONRPC: "2.0",
			ID:      request.ID,
			Result:  map[string]interface{}{},
		}

	case "tools/list":
		return s.handleToolsList()

	case "tools/call":
		var params CallToolParams
		if request.Params != nil {
			if paramsBytes, err := json.Marshal(request.Params); err == nil {
				if err := json.Unmarshal(paramsBytes, &params); err != nil {
					return s.errorResponse(-32602, "Invalid parameters")
				}
			}
		}
		return s.handleCallTool(params)

	default:
		return s.errorResponse(-32601, "Method not found")
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
		response.ID = request.ID

		if err := encoder.Encode(response); err != nil {
			log.Printf("Error encoding response: %v", err)
		}
	}
}

func (s *MCPServer) handleToolsList() *MCPResponse {
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
	}

	result := ToolListResult{Tools: tools}
	return &MCPResponse{
		JSONRPC: "2.0",
		Result:  result,
	}
}
