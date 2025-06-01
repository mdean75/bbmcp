package server

import (
	"encoding/json"
	"fmt"

	"bbcli/pkg/bitbucket"
	"bbcli/pkg/types"
)

func (s *MCPServer) handleListPullRequests(args map[string]interface{}) *types.MCPResponse {
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

	result := types.CallToolResult{
		Content: []types.ToolContent{
			{
				Type: "text",
				Text: string(content),
			},
		},
	}

	return &types.MCPResponse{
		JSONRPC: "2.0",
		Result:  result,
	}
}

func (s *MCPServer) handleGetPullRequest(args map[string]interface{}) *types.MCPResponse {
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

	result := types.CallToolResult{
		Content: []types.ToolContent{
			{
				Type: "text",
				Text: string(content),
			},
		},
	}

	return &types.MCPResponse{
		JSONRPC: "2.0",
		Result:  result,
	}
}

func (s *MCPServer) handleGetPullRequestActivity(args map[string]interface{}) *types.MCPResponse {
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

	result := types.CallToolResult{
		Content: []types.ToolContent{
			{
				Type: "text",
				Text: string(content),
			},
		},
	}

	return &types.MCPResponse{
		JSONRPC: "2.0",
		Result:  result,
	}
}

func (s *MCPServer) handleCreatePullRequest(args map[string]interface{}) *types.MCPResponse {
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
	pr := &bitbucket.PullRequest{
		Title:       title,
		Description: description,
		FromRef: bitbucket.PullRequestRef{
			ID: fromBranch,
			Repository: bitbucket.Repository{
				Slug: repoSlug,
				Project: bitbucket.Project{
					Key: projectKey,
				},
			},
		},
		ToRef: bitbucket.PullRequestRef{
			ID: toBranch,
			Repository: bitbucket.Repository{
				Slug: repoSlug,
				Project: bitbucket.Project{
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
					pr.Reviewers = append(pr.Reviewers, bitbucket.Reviewer{
						User: bitbucket.User{
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

	result := types.CallToolResult{
		Content: []types.ToolContent{
			{
				Type: "text",
				Text: string(content),
			},
		},
	}

	return &types.MCPResponse{
		JSONRPC: "2.0",
		Result:  result,
	}
}

func (s *MCPServer) handleApprovePullRequest(args map[string]interface{}) *types.MCPResponse {
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

	result := types.CallToolResult{
		Content: []types.ToolContent{
			{
				Type: "text",
				Text: "Pull request approved successfully",
			},
		},
	}

	return &types.MCPResponse{
		JSONRPC: "2.0",
		Result:  result,
	}
}

func (s *MCPServer) handleUnapprovePullRequest(args map[string]interface{}) *types.MCPResponse {
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

	result := types.CallToolResult{
		Content: []types.ToolContent{
			{
				Type: "text",
				Text: "Pull request approval removed successfully",
			},
		},
	}

	return &types.MCPResponse{
		JSONRPC: "2.0",
		Result:  result,
	}
}

func (s *MCPServer) handleMergePullRequest(args map[string]interface{}) *types.MCPResponse {
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

	// Get current PR to obtain the latest version for optimistic locking
	currentPR, err := s.bitbucket.GetPullRequest(projectKey, repoSlug, int(pullRequestID))
	if err != nil {
		return s.errorResponse(-32000, fmt.Sprintf("Failed to get current pull request version: %v", err))
	}

	mergedPR, err := s.bitbucket.MergePullRequest(projectKey, repoSlug, int(pullRequestID), currentPR.Version)
	if err != nil {
		return s.errorResponse(-32000, fmt.Sprintf("Failed to merge pull request: %v", err))
	}

	content, err := json.MarshalIndent(mergedPR, "", "  ")
	if err != nil {
		return s.errorResponse(-32000, fmt.Sprintf("Failed to marshal response: %v", err))
	}

	result := types.CallToolResult{
		Content: []types.ToolContent{
			{
				Type: "text",
				Text: string(content),
			},
		},
	}

	return &types.MCPResponse{
		JSONRPC: "2.0",
		Result:  result,
	}
}

func (s *MCPServer) handleDeclinePullRequest(args map[string]interface{}) *types.MCPResponse {
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

	// Get current PR to obtain the latest version for optimistic locking
	currentPR, err := s.bitbucket.GetPullRequest(projectKey, repoSlug, int(pullRequestID))
	if err != nil {
		return s.errorResponse(-32000, fmt.Sprintf("Failed to get current pull request version: %v", err))
	}

	declinedPR, err := s.bitbucket.DeclinePullRequest(projectKey, repoSlug, int(pullRequestID), currentPR.Version)
	if err != nil {
		return s.errorResponse(-32000, fmt.Sprintf("Failed to decline pull request: %v", err))
	}

	content, err := json.MarshalIndent(declinedPR, "", "  ")
	if err != nil {
		return s.errorResponse(-32000, fmt.Sprintf("Failed to marshal response: %v", err))
	}

	result := types.CallToolResult{
		Content: []types.ToolContent{
			{
				Type: "text",
				Text: string(content),
			},
		},
	}

	return &types.MCPResponse{
		JSONRPC: "2.0",
		Result:  result,
	}
}

func (s *MCPServer) handleCreatePullRequestComment(args map[string]interface{}) *types.MCPResponse {
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
	var anchor *bitbucket.CommentAnchor
	if anchorData, ok := args["anchor"]; ok {
		if anchorMap, ok := anchorData.(map[string]interface{}); ok {
			anchor = &bitbucket.CommentAnchor{}

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
			if diffType, ok := anchorMap["diff_type"].(string); ok {
				anchor.DiffType = diffType
			}
			if orphanedType, ok := anchorMap["orphaned_type"].(string); ok {
				anchor.OrphanedType = orphanedType
			}

			// Auto-set diffType when commit hashes are provided but diffType is not specified
			if (anchor.FromHash != "" || anchor.ToHash != "") && anchor.DiffType == "" {
				anchor.DiffType = "RANGE" // Default to RANGE when commit hashes are provided
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

	result := types.CallToolResult{
		Content: []types.ToolContent{
			{
				Type: "text",
				Text: string(content),
			},
		},
	}

	return &types.MCPResponse{
		JSONRPC: "2.0",
		Result:  result,
	}
}

func (s *MCPServer) handleGetPullRequestDiff(args map[string]interface{}) *types.MCPResponse {
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

	result := types.CallToolResult{
		Content: []types.ToolContent{
			{
				Type: "text",
				Text: diff,
			},
		},
	}

	return &types.MCPResponse{
		JSONRPC: "2.0",
		Result:  result,
	}
}

func (s *MCPServer) handleHelloWorld(args map[string]interface{}) *types.MCPResponse {
	name := "World"
	if nameVal, ok := args["name"].(string); ok && nameVal != "" {
		name = nameVal
	}

	result := types.CallToolResult{
		Content: []types.ToolContent{
			{
				Type: "text",
				Text: fmt.Sprintf("Hello, %s!", name),
			},
		},
	}

	return &types.MCPResponse{
		JSONRPC: "2.0",
		Result:  result,
	}
}

func (s *MCPServer) errorResponse(code int, message string) *types.MCPResponse {
	return &types.MCPResponse{
		JSONRPC: "2.0",
		Error: &types.MCPError{
			Code:    code,
			Message: message,
		},
	}
}
