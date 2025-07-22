package main

import (
	"log"
	"os"
	"strings"

	"bbcli/pkg/bitbucket"
	"bbcli/pkg/tools"
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
		BaseURL:           os.Getenv("BITBUCKET_BASE_URL"),
		Username:          os.Getenv("BITBUCKET_USERNAME"),
		Password:          os.Getenv("BITBUCKET_PASSWORD"),
		Token:             os.Getenv("BITBUCKET_TOKEN"),
		DefaultProjectKey: os.Getenv("BITBUCKET_DEFAULT_PROJECT_KEY"),
	}

	var missing []string
	if config.BaseURL == "" {
		missing = append(missing, "BITBUCKET_BASE_URL")
	}

	if config.Token == "" && (config.Username == "" || config.Password == "") {
		missing = append(missing, "BITBUCKET_TOKEN or both BITBUCKET_USERNAME and BITBUCKET_PASSWORD")
	}

	if len(missing) > 0 {
		log.Fatalf("Missing required environment variables: %s", strings.Join(missing, ", "))
	}

	return config
}

func registerBitbucketTools(s *server.MCPServer, bb *bitbucket.Server) {
	tools.RegisterListPullRequests(s, bb)
	tools.RegisterGetPullRequest(s, bb)
	tools.RegisterGetPullRequestActivity(s, bb)
	tools.RegisterCreatePullRequest(s, bb)
	tools.RegisterApprovePullRequest(s, bb)
	tools.RegisterUnapprovePullRequest(s, bb)
	tools.RegisterMergePullRequest(s, bb)
	tools.RegisterDeclinePullRequest(s, bb)
	tools.RegisterGetPullRequestDiff(s, bb)
	tools.RegisterCreatePullRequestComment(s, bb)

	tools.RegisterGetRepos(s, bb)
	tools.RegisterGetPullRequestSettings(s, bb)

	tools.RegisterHelloWorld(s)
}
