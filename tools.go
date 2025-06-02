package main

import (
	"context"
	"encoding/json"
	"fmt"

	"bbcli/pkg/bitbucket"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func registerListPullRequestsTool(s *server.MCPServer, bb *bitbucket.Server) {
	listPRTool := mcp.NewTool("list_pull_requests",
		mcp.WithDescription("List pull requests for a repository"),
		mcp.WithString("project_key",
			mcp.Required(),
			mcp.Description("The project key"),
		),
		mcp.WithString("repo_slug",
			mcp.Required(),
			mcp.Description("The repository slug"),
		),
		mcp.WithString("state",
			mcp.Description("Filter by state (OPEN, MERGED, DECLINED)"),
			mcp.Enum("OPEN", "MERGED", "DECLINED"),
		),
		mcp.WithNumber("limit",
			mcp.Description("Maximum number of results to return (1-100)"),
		),
	)

	s.AddTool(listPRTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()

		projectKey, _ := args["project_key"].(string)
		repoSlug, _ := args["repo_slug"].(string)
		state, _ := args["state"].(string)

		limit := 25
		if limitVal, ok := args["limit"].(float64); ok {
			limit = int(limitVal)
		}

		prs, err := bb.GetPullRequests(projectKey, repoSlug, state, limit)
		if err != nil {
			return nil, fmt.Errorf("failed to get pull requests: %v", err)
		}

		content, err := json.MarshalIndent(prs, "", "  ")
		if err != nil {
			return nil, fmt.Errorf("failed to marshal response: %v", err)
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: string(content),
				},
			},
		}, nil
	})
}

func registerGetPullRequestTool(s *server.MCPServer, bb *bitbucket.Server) {
	getPRTool := mcp.NewTool("get_pull_request",
		mcp.WithDescription("Get details of a specific pull request"),
		mcp.WithString("project_key",
			mcp.Required(),
			mcp.Description("The project key"),
		),
		mcp.WithString("repo_slug",
			mcp.Required(),
			mcp.Description("The repository slug"),
		),
		mcp.WithNumber("pull_request_id",
			mcp.Required(),
			mcp.Description("The pull request ID"),
		),
	)

	s.AddTool(getPRTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()

		projectKey, _ := args["project_key"].(string)
		repoSlug, _ := args["repo_slug"].(string)
		pullRequestID, _ := args["pull_request_id"].(float64)

		pr, err := bb.GetPullRequest(projectKey, repoSlug, int(pullRequestID))
		if err != nil {
			return nil, fmt.Errorf("failed to get pull request: %v", err)
		}

		content, err := json.MarshalIndent(pr, "", "  ")
		if err != nil {
			return nil, fmt.Errorf("failed to marshal response: %v", err)
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: string(content),
				},
			},
		}, nil
	})
}

func registerGetPullRequestActivityTool(s *server.MCPServer, bb *bitbucket.Server) {
	getActivityTool := mcp.NewTool("get_pull_request_activity",
		mcp.WithDescription("Get activity (comments, approvals, etc.) for a pull request"),
		mcp.WithString("project_key",
			mcp.Required(),
			mcp.Description("The project key"),
		),
		mcp.WithString("repo_slug",
			mcp.Required(),
			mcp.Description("The repository slug"),
		),
		mcp.WithNumber("pull_request_id",
			mcp.Required(),
			mcp.Description("The pull request ID"),
		),
	)

	s.AddTool(getActivityTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()

		projectKey, _ := args["project_key"].(string)
		repoSlug, _ := args["repo_slug"].(string)
		pullRequestID, _ := args["pull_request_id"].(float64)

		activity, err := bb.GetPullRequestActivity(projectKey, repoSlug, int(pullRequestID))
		if err != nil {
			return nil, fmt.Errorf("failed to get pull request activity: %v", err)
		}

		content, err := json.MarshalIndent(activity, "", "  ")
		if err != nil {
			return nil, fmt.Errorf("failed to marshal response: %v", err)
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: string(content),
				},
			},
		}, nil
	})
}

func registerCreatePullRequestTool(s *server.MCPServer, bb *bitbucket.Server) {
	createPRTool := mcp.NewTool("create_pull_request",
		mcp.WithDescription("Create a new pull request"),
		mcp.WithString("project_key",
			mcp.Required(),
			mcp.Description("The project key"),
		),
		mcp.WithString("repo_slug",
			mcp.Required(),
			mcp.Description("The repository slug"),
		),
		mcp.WithString("title",
			mcp.Required(),
			mcp.Description("The pull request title"),
		),
		mcp.WithString("from_branch",
			mcp.Required(),
			mcp.Description("Source branch name"),
		),
		mcp.WithString("to_branch",
			mcp.Required(),
			mcp.Description("Target branch name"),
		),
		mcp.WithString("description",
			mcp.Description("The pull request description"),
		),
	)

	s.AddTool(createPRTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()

		projectKey, _ := args["project_key"].(string)
		repoSlug, _ := args["repo_slug"].(string)
		title, _ := args["title"].(string)
		fromBranch, _ := args["from_branch"].(string)
		toBranch, _ := args["to_branch"].(string)
		description, _ := args["description"].(string)

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

		createdPR, err := bb.CreatePullRequest(projectKey, repoSlug, pr)
		if err != nil {
			return nil, fmt.Errorf("failed to create pull request: %v", err)
		}

		content, err := json.MarshalIndent(createdPR, "", "  ")
		if err != nil {
			return nil, fmt.Errorf("failed to marshal response: %v", err)
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: string(content),
				},
			},
		}, nil
	})
}

func registerApprovePullRequestTool(s *server.MCPServer, bb *bitbucket.Server) {
	approveTool := mcp.NewTool("approve_pull_request",
		mcp.WithDescription("Approve a pull request"),
		mcp.WithString("project_key",
			mcp.Required(),
			mcp.Description("The project key"),
		),
		mcp.WithString("repo_slug",
			mcp.Required(),
			mcp.Description("The repository slug"),
		),
		mcp.WithNumber("pull_request_id",
			mcp.Required(),
			mcp.Description("The pull request ID"),
		),
	)

	s.AddTool(approveTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()

		projectKey, _ := args["project_key"].(string)
		repoSlug, _ := args["repo_slug"].(string)
		pullRequestID, _ := args["pull_request_id"].(float64)

		err := bb.ApprovePullRequest(projectKey, repoSlug, int(pullRequestID))
		if err != nil {
			return nil, fmt.Errorf("failed to approve pull request: %v", err)
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: "Pull request approved successfully",
				},
			},
		}, nil
	})
}

func registerUnapprovePullRequestTool(s *server.MCPServer, bb *bitbucket.Server) {
	unapproveTool := mcp.NewTool("unapprove_pull_request",
		mcp.WithDescription("Remove approval from a pull request"),
		mcp.WithString("project_key",
			mcp.Required(),
			mcp.Description("The project key"),
		),
		mcp.WithString("repo_slug",
			mcp.Required(),
			mcp.Description("The repository slug"),
		),
		mcp.WithNumber("pull_request_id",
			mcp.Required(),
			mcp.Description("The pull request ID"),
		),
	)

	s.AddTool(unapproveTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()

		projectKey, _ := args["project_key"].(string)
		repoSlug, _ := args["repo_slug"].(string)
		pullRequestID, _ := args["pull_request_id"].(float64)

		err := bb.UnapprovalPullRequest(projectKey, repoSlug, int(pullRequestID))
		if err != nil {
			return nil, fmt.Errorf("failed to unapprove pull request: %v", err)
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: "Pull request approval removed successfully",
				},
			},
		}, nil
	})
}

func registerMergePullRequestTool(s *server.MCPServer, bb *bitbucket.Server) {
	mergeTool := mcp.NewTool("merge_pull_request",
		mcp.WithDescription("Merge a pull request (automatically fetches current version for optimistic locking)"),
		mcp.WithString("project_key",
			mcp.Required(),
			mcp.Description("The project key"),
		),
		mcp.WithString("repo_slug",
			mcp.Required(),
			mcp.Description("The repository slug"),
		),
		mcp.WithNumber("pull_request_id",
			mcp.Required(),
			mcp.Description("The pull request ID"),
		),
	)

	s.AddTool(mergeTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()

		projectKey, _ := args["project_key"].(string)
		repoSlug, _ := args["repo_slug"].(string)
		pullRequestID, _ := args["pull_request_id"].(float64)

		// Get current PR to obtain the latest version for optimistic locking
		currentPR, err := bb.GetPullRequest(projectKey, repoSlug, int(pullRequestID))
		if err != nil {
			return nil, fmt.Errorf("failed to get current pull request version: %v", err)
		}

		mergedPR, err := bb.MergePullRequest(projectKey, repoSlug, int(pullRequestID), currentPR.Version)
		if err != nil {
			return nil, fmt.Errorf("failed to merge pull request: %v", err)
		}

		content, err := json.MarshalIndent(mergedPR, "", "  ")
		if err != nil {
			return nil, fmt.Errorf("failed to marshal response: %v", err)
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: string(content),
				},
			},
		}, nil
	})
}

func registerDeclinePullRequestTool(s *server.MCPServer, bb *bitbucket.Server) {
	declineTool := mcp.NewTool("decline_pull_request",
		mcp.WithDescription("Decline a pull request (automatically fetches current version for optimistic locking)"),
		mcp.WithString("project_key",
			mcp.Required(),
			mcp.Description("The project key"),
		),
		mcp.WithString("repo_slug",
			mcp.Required(),
			mcp.Description("The repository slug"),
		),
		mcp.WithNumber("pull_request_id",
			mcp.Required(),
			mcp.Description("The pull request ID"),
		),
	)

	s.AddTool(declineTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()

		projectKey, _ := args["project_key"].(string)
		repoSlug, _ := args["repo_slug"].(string)
		pullRequestID, _ := args["pull_request_id"].(float64)

		// Get current PR to obtain the latest version for optimistic locking
		currentPR, err := bb.GetPullRequest(projectKey, repoSlug, int(pullRequestID))
		if err != nil {
			return nil, fmt.Errorf("failed to get current pull request version: %v", err)
		}

		declinedPR, err := bb.DeclinePullRequest(projectKey, repoSlug, int(pullRequestID), currentPR.Version)
		if err != nil {
			return nil, fmt.Errorf("failed to decline pull request: %v", err)
		}

		content, err := json.MarshalIndent(declinedPR, "", "  ")
		if err != nil {
			return nil, fmt.Errorf("failed to marshal response: %v", err)
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: string(content),
				},
			},
		}, nil
	})
}

func registerGetPullRequestDiffTool(s *server.MCPServer, bb *bitbucket.Server) {
	getDiffTool := mcp.NewTool("get_pull_request_diff",
		mcp.WithDescription("Get the raw diff for a pull request"),
		mcp.WithString("project_key",
			mcp.Required(),
			mcp.Description("The project key"),
		),
		mcp.WithString("repo_slug",
			mcp.Required(),
			mcp.Description("The repository slug"),
		),
		mcp.WithNumber("pull_request_id",
			mcp.Required(),
			mcp.Description("The pull request ID"),
		),
		mcp.WithNumber("context_lines",
			mcp.Description("Number of context lines around changes (optional)"),
		),
		mcp.WithString("whitespace",
			mcp.Description("Whitespace handling"),
			mcp.Enum("ignore-all", "ignore-space-at-eol", "ignore-space-change", "ignore-trailing-space"),
		),
		mcp.WithString("since",
			mcp.Description("Base commit hash to diff from (optional)"),
		),
		mcp.WithString("until",
			mcp.Description("End commit hash to diff to (optional)"),
		),
	)

	s.AddTool(getDiffTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()

		projectKey, _ := args["project_key"].(string)
		repoSlug, _ := args["repo_slug"].(string)
		pullRequestID, _ := args["pull_request_id"].(float64)

		// Optional parameters
		contextLines := 0
		if contextVal, ok := args["context_lines"].(float64); ok {
			contextLines = int(contextVal)
		}

		whitespace, _ := args["whitespace"].(string)
		since, _ := args["since"].(string)
		until, _ := args["until"].(string)

		diff, err := bb.GetPullRequestDiff(projectKey, repoSlug, int(pullRequestID), contextLines, whitespace, since, until)
		if err != nil {
			return nil, fmt.Errorf("failed to get pull request diff: %v", err)
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: diff,
				},
			},
		}, nil
	})
}

func registerCreatePullRequestCommentTool(s *server.MCPServer, bb *bitbucket.Server) {
	commentTool := mcp.NewTool("create_pull_request_comment",
		mcp.WithDescription("Add a comment to a pull request"),
		mcp.WithString("project_key",
			mcp.Required(),
			mcp.Description("The project key"),
		),
		mcp.WithString("repo_slug",
			mcp.Required(),
			mcp.Description("The repository slug"),
		),
		mcp.WithNumber("pull_request_id",
			mcp.Required(),
			mcp.Description("The pull request ID"),
		),
		mcp.WithString("text",
			mcp.Required(),
			mcp.Description("The comment text"),
		),
		mcp.WithString("anchor_json",
			mcp.Description("Optional JSON-encoded anchor for inline comments (contains: line, line_type, path, file_type, from_hash, to_hash, src_path, dst_path, diff_type, orphaned_type)"),
		),
	)

	s.AddTool(commentTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()

		projectKey, _ := args["project_key"].(string)
		repoSlug, _ := args["repo_slug"].(string)
		pullRequestID, _ := args["pull_request_id"].(float64)
		text, _ := args["text"].(string)

		// Optional anchor for inline comments
		var anchor *bitbucket.CommentAnchor
		if anchorJSON, ok := args["anchor_json"].(string); ok && anchorJSON != "" {
			anchor = &bitbucket.CommentAnchor{}
			var anchorData map[string]interface{}
			if err := json.Unmarshal([]byte(anchorJSON), &anchorData); err == nil {
				if line, ok := anchorData["line"].(float64); ok {
					anchor.Line = int(line)
				}
				if lineType, ok := anchorData["line_type"].(string); ok {
					anchor.LineType = lineType
				}
				if path, ok := anchorData["path"].(string); ok {
					anchor.Path = path
				}
				if fileType, ok := anchorData["file_type"].(string); ok {
					anchor.FileType = fileType
				}
				if fromHash, ok := anchorData["from_hash"].(string); ok {
					anchor.FromHash = fromHash
				}
				if toHash, ok := anchorData["to_hash"].(string); ok {
					anchor.ToHash = toHash
				}
				if srcPath, ok := anchorData["src_path"].(string); ok {
					anchor.SrcPath = srcPath
				}
				if dstPath, ok := anchorData["dst_path"].(string); ok {
					anchor.DstPath = dstPath
				}
				if diffType, ok := anchorData["diff_type"].(string); ok {
					anchor.DiffType = diffType
				}
				if orphanedType, ok := anchorData["orphaned_type"].(string); ok {
					anchor.OrphanedType = orphanedType
				}

				// Auto-set diffType when commit hashes are provided but diffType is not specified
				if (anchor.FromHash != "" || anchor.ToHash != "") && anchor.DiffType == "" {
					anchor.DiffType = "RANGE"
				}
			}
		}

		comment, err := bb.CreatePullRequestComment(projectKey, repoSlug, int(pullRequestID), text, anchor)
		if err != nil {
			return nil, fmt.Errorf("failed to create pull request comment: %v", err)
		}

		content, err := json.MarshalIndent(comment, "", "  ")
		if err != nil {
			return nil, fmt.Errorf("failed to marshal response: %v", err)
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: string(content),
				},
			},
		}, nil
	})
}

func registerHelloWorldTool(s *server.MCPServer) {
	helloTool := mcp.NewTool("hello_world",
		mcp.WithDescription("Say hello to someone"),
		mcp.WithString("name",
			mcp.Required(),
			mcp.Description("Name of the person to greet"),
		),
	)

	s.AddTool(helloTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		args := request.GetArguments()

		name := "World"
		if nameVal, ok := args["name"].(string); ok && nameVal != "" {
			name = nameVal
		}

		return &mcp.CallToolResult{
			Content: []mcp.Content{
				mcp.TextContent{
					Type: "text",
					Text: fmt.Sprintf("Hello, %s!", name),
				},
			},
		}, nil
	})
}
