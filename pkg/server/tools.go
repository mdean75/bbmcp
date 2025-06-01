package server

import (
	"encoding/json"

	"bbcli/pkg/types"
)

func (s *MCPServer) HandleRequest(request *types.MCPRequest) *types.MCPResponse {
	switch request.Method {
	case "initialize":
		var params types.InitializeParams
		if request.Params != nil {
			if paramsBytes, err := json.Marshal(request.Params); err == nil {
				json.Unmarshal(paramsBytes, &params)
			}
		}
		return s.HandleInitialize(params, request.ID)

	case "notifications/initialized":
		// No response for notifications
		return nil

	case "tools/list":
		return s.handleToolsList(request.ID)

	case "tools/call":
		var params types.CallToolParams
		if request.Params != nil {
			if paramsBytes, err := json.Marshal(request.Params); err == nil {
				if err := json.Unmarshal(paramsBytes, &params); err != nil {
					response := s.errorResponse(-32602, "Invalid parameters")
					response.ID = request.ID
					return response
				}
			}
		}
		return s.HandleCallTool(params, request.ID)

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

func (s *MCPServer) handleToolsList(id interface{}) *types.MCPResponse {
	tools := []types.Tool{
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
			Description: "Merge a pull request (automatically fetches current version for optimistic locking)",
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
			Name:        "decline_pull_request",
			Description: "Decline a pull request (automatically fetches current version for optimistic locking)",
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
							"diff_type": map[string]interface{}{
								"type":        "string",
								"description": "Diff type (EFFECTIVE, RANGE, COMMIT) - defaults to RANGE when commit hashes are provided",
								"enum":        []string{"EFFECTIVE", "RANGE", "COMMIT"},
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

	result := types.ToolListResult{Tools: tools}
	return &types.MCPResponse{
		JSONRPC: "2.0",
		ID:      id,
		Result:  result,
	}
}
