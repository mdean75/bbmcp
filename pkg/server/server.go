package server

import (
	"log"
	"os"

	"bbcli/pkg/bitbucket"
	"bbcli/pkg/types"
)

// MCPServer implements the Model Context Protocol
type MCPServer struct {
	bitbucket *bitbucket.Server
}

// NewMCPServer creates a new MCP server instance
func NewMCPServer() *MCPServer {
	config := &bitbucket.Config{
		BaseURL:  os.Getenv("BITBUCKET_BASE_URL"),
		Username: os.Getenv("BITBUCKET_USERNAME"),
		Password: os.Getenv("BITBUCKET_PASSWORD"),
	}

	if config.BaseURL == "" || config.Username == "" || config.Password == "" {
		log.Fatal("Missing required environment variables: BITBUCKET_BASE_URL, BITBUCKET_USERNAME, BITBUCKET_PASSWORD")
	}

	return &MCPServer{
		bitbucket: bitbucket.NewServer(config),
	}
}

func (s *MCPServer) HandleInitialize(params types.InitializeParams, id interface{}) *types.MCPResponse {
	result := types.InitializeResult{
		ProtocolVersion: "2024-11-05",
		Capabilities: types.ServerCapabilities{
			Tools: map[string]interface{}{
				"listChanged": true,
			},
		},
		ServerInfo: map[string]string{
			"name":    "bbcli",
			"version": "1.0.0",
		},
	}

	return &types.MCPResponse{
		JSONRPC: "2.0",
		ID:      id,
		Result:  result,
	}
}

func (s *MCPServer) HandleCallTool(params types.CallToolParams, id interface{}) *types.MCPResponse {
	var response *types.MCPResponse
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
		response = &types.MCPResponse{
			JSONRPC: "2.0",
			Error: &types.MCPError{
				Code:    -32601,
				Message: "Unknown tool",
			},
		}
	}
	response.ID = id
	return response
}
