package main

import (
	"log"
	"os"

	"bbcli/pkg/bitbucket"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	mcpServer := NewMCPServer()

	if err := server.ServeStdio(mcpServer); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

func NewMCPServer() *server.MCPServer {
	config := getBitbucketConfig()
	bbClient := bitbucket.NewServer(config)

	s := server.NewMCPServer(
		"bbcli",
		"1.0.0",
		server.WithToolCapabilities(true),
	)

	registerBitbucketTools(s, bbClient)
	return s
}

func getBitbucketConfig() *bitbucket.Config {
	config := &bitbucket.Config{
		BaseURL:  os.Getenv("BITBUCKET_BASE_URL"),
		Username: os.Getenv("BITBUCKET_USERNAME"),
		Password: os.Getenv("BITBUCKET_PASSWORD"),
	}

	if config.BaseURL == "" || config.Username == "" || config.Password == "" {
		log.Fatal("Missing required environment variables: BITBUCKET_BASE_URL, BITBUCKET_USERNAME, BITBUCKET_PASSWORD")
	}

	return config
}

func registerBitbucketTools(s *server.MCPServer, bb *bitbucket.Server) {
	registerListPullRequestsTool(s, bb)
	registerGetPullRequestTool(s, bb)
	registerGetPullRequestActivityTool(s, bb)
	registerCreatePullRequestTool(s, bb)
	registerApprovePullRequestTool(s, bb)
	registerUnapprovePullRequestTool(s, bb)
	registerMergePullRequestTool(s, bb)
	registerDeclinePullRequestTool(s, bb)
	registerGetPullRequestDiffTool(s, bb)
	registerCreatePullRequestCommentTool(s, bb)
	registerHelloWorldTool(s)
}
