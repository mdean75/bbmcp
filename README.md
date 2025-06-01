# Bitbucket MCP Server

A Model Context Protocol (MCP) server that provides tools for interacting with Bitbucket Server pull requests.

## Features

- List pull requests for a repository
- Get detailed information about specific pull requests
- View pull request activity (comments, approvals, etc.)
- Create new pull requests
- Approve/unapprove pull requests
- Merge pull requests
- Decline pull requests

## Environment Variables

Set the following environment variables:

```bash
export BITBUCKET_BASE_URL="https://your-bitbucket-server.com"
export BITBUCKET_USERNAME="your-username"
export BITBUCKET_PASSWORD="your-app-password"
```

## Build and Run

```bash
# Build the project
go build -o bbcli

# Run the MCP server
./bbcli
```

## MCP Tools

The server provides the following MCP tools:

### list_pull_requests
List pull requests for a repository.

**Parameters:**
- `project_key` (required): The project key
- `repo_slug` (required): The repository slug
- `state` (optional): Filter by state (OPEN, MERGED, DECLINED)
- `limit` (optional): Maximum number of results (1-100, default: 25)

### get_pull_request
Get details of a specific pull request.

**Parameters:**
- `project_key` (required): The project key
- `repo_slug` (required): The repository slug
- `pull_request_id` (required): The pull request ID

### get_pull_request_activity
Get activity for a pull request (comments, approvals, etc.).

**Parameters:**
- `project_key` (required): The project key
- `repo_slug` (required): The repository slug
- `pull_request_id` (required): The pull request ID

### create_pull_request
Create a new pull request.

**Parameters:**
- `project_key` (required): The project key
- `repo_slug` (required): The repository slug
- `title` (required): The pull request title
- `from_branch` (required): Source branch name
- `to_branch` (required): Target branch name
- `description` (optional): The pull request description
- `reviewers` (optional): Array of reviewer usernames

### approve_pull_request
Approve a pull request.

**Parameters:**
- `project_key` (required): The project key
- `repo_slug` (required): The repository slug
- `pull_request_id` (required): The pull request ID

### unapprove_pull_request
Remove approval from a pull request.

**Parameters:**
- `project_key` (required): The project key
- `repo_slug` (required): The repository slug
- `pull_request_id` (required): The pull request ID

### merge_pull_request
Merge a pull request.

**Parameters:**
- `project_key` (required): The project key
- `repo_slug` (required): The repository slug
- `pull_request_id` (required): The pull request ID
- `version` (required): The pull request version for optimistic locking

### decline_pull_request
Decline a pull request.

**Parameters:**
- `project_key` (required): The project key
- `repo_slug` (required): The repository slug
- `pull_request_id` (required): The pull request ID
- `version` (required): The pull request version for optimistic locking

## Usage with MCP Clients

This server communicates via STDIO using the Model Context Protocol. It can be used with any MCP-compatible client such as Claude Desktop or VS Code with MCP support.

## Security

- Uses HTTP Basic Authentication with Bitbucket Server
- Requires valid Bitbucket Server credentials
- All API requests are made over HTTPS (when configured)
