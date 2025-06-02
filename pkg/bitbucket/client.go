package bitbucket

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// Server handles Bitbucket Server API operations
type Server struct {
	config *Config
	client *http.Client
}

// NewServer creates a new Bitbucket Server API client
func NewServer(config *Config) *Server {
	return &Server{
		config: config,
		client: &http.Client{},
	}
}

func (bs *Server) makeRequest(method, endpoint string, body io.Reader) (*http.Response, error) {
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

func (bs *Server) GetPullRequests(projectKey, repoSlug string, state string, limit int) ([]PullRequest, error) {
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

func (bs *Server) GetPullRequest(projectKey, repoSlug string, pullRequestID int) (*PullRequest, error) {
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

func (bs *Server) GetPullRequestActivity(projectKey, repoSlug string, pullRequestID int) (*PullRequestActivity, error) {
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

func (bs *Server) CreatePullRequest(projectKey, repoSlug string, pr *PullRequest) (*PullRequest, error) {
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

func (bs *Server) ApprovePullRequest(projectKey, repoSlug string, pullRequestID int) error {
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

func (bs *Server) UnapprovalPullRequest(projectKey, repoSlug string, pullRequestID int) error {
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

func (bs *Server) MergePullRequest(projectKey, repoSlug string, pullRequestID int, version int) (*PullRequest, error) {
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

func (bs *Server) DeclinePullRequest(projectKey, repoSlug string, pullRequestID int, version int) (*PullRequest, error) {
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

func (bs *Server) GetPullRequestDiff(projectKey, repoSlug string, pullRequestID int, contextLines int, whitespace string, since string, until string) (string, error) {
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

func (bs *Server) CreatePullRequestComment(projectKey, repoSlug string, pullRequestID int, text string, anchor *CommentAnchor) (*Comment, error) {
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

type RepositoryResponse struct {
	Size       int          `json:"size"`
	Limit      int          `json:"limit"`
	IsLastPage bool         `json:"isLastPage"`
	Values     []Repository `json:"values"`
	Start      int          `json:"start"`
}

func (bs *Server) GetRepos(projectKey string, limit, start int) ([]Repository, error) {
	var allRepos []Repository

	for {
		endpoint := fmt.Sprintf("/projects/%s/repos?start=%d&limit=%d", projectKey, start, limit)

		resp, err := bs.makeRequest("GET", endpoint, nil)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("API request failed with status %d", resp.StatusCode)
		}

		var pageResponse RepositoryResponse
		if err := json.NewDecoder(resp.Body).Decode(&pageResponse); err != nil {
			return nil, err
		}

		allRepos = append(allRepos, pageResponse.Values...)

		if pageResponse.IsLastPage {
			break
		}
		start = start + pageResponse.Limit
	}

	return allRepos, nil
}

func (bs *Server) GetPullRequestSettings(projectKey, repoSlug string) (*PullRequestSettings, error) {
	endpoint := fmt.Sprintf("/projects/%s/repos/%s/settings/pull-requests", projectKey, repoSlug)

	resp, err := bs.makeRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d", resp.StatusCode)
	}

	var settings PullRequestSettings
	if err := json.NewDecoder(resp.Body).Decode(&settings); err != nil {
		return nil, err
	}

	return &settings, nil
}
