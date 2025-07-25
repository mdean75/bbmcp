# Bitbucket MCP Server

A Model Context Protocol (MCP) server that provides tools for interacting with Bitbucket Server pull requests. Built using the [mcp-go](https://github.com/mark3labs/mcp-go) library.

## Project Structure

The project is organized for simplicity and maintainability:

```
bbcli/
├── main.go                     # Application entry point and server setup
├── tools.go                    # Tool registrations and handlers  
├── pkg/
│   └── bitbucket/             # Bitbucket API client
│       ├── types.go           # Bitbucket API data structures
│       └── client.go          # HTTP client for Bitbucket Server API
├── go.mod
└── README.md
```

## Key Features

- **Framework-based**: Uses mcp-go library for robust MCP protocol handling
- **Automatic versioning**: Merge/decline operations automatically fetch current PR versions to prevent conflicts
- **Type-safe**: Leverages mcp-go's type-safe tool definitions and parameter validation
- **Clean separation**: Bitbucket API logic separated from MCP server concerns

## Tools Available

- List pull requests for a repository
- Get detailed information about specific pull requests
- View pull request activity (comments, approvals, etc.)
- Get raw diff for pull requests
- Add comments to pull requests (general and inline comments)
- Create new pull requests
- Approve/unapprove pull requests
- Merge pull requests (with automatic version handling)
- Decline pull requests (with automatic version handling)
- List repositories in a project
- Get pull request configuration settings

## Environment Variables

Set the following environment variables:

```bash
export BITBUCKET_BASE_URL="https://your-bitbucket-server.com"
export BITBUCKET_USERNAME="your-username"
export BITBUCKET_PASSWORD="your-app-password"
# OR use a token instead of username/password
export BITBUCKET_TOKEN="your-personal-access-token"

# Optional: Set a default project key to avoid specifying it in every tool call
export BITBUCKET_DEFAULT_PROJECT_KEY="MYPROJ"
```

### Authentication Options

You can authenticate using either:
1. **Username and Password/App Password**: Set both `BITBUCKET_USERNAME` and `BITBUCKET_PASSWORD`
2. **Personal Access Token**: Set `BITBUCKET_TOKEN` (preferred for security)

### Default Project Key

When `BITBUCKET_DEFAULT_PROJECT_KEY` is set, all tools that require a `project_key` parameter will use this default value when the parameter is not explicitly provided. This simplifies usage when working primarily with repositories in a single project.

**Benefits:**
- Reduces repetitive parameter specification
- Maintains backward compatibility - explicit `project_key` parameters still work
- Clear error messages when no project key is available

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
- `project_key` (optional): The project key (uses default if BITBUCKET_DEFAULT_PROJECT_KEY is set)
- `repo_slug` (required): The repository slug
- `state` (optional): Filter by state (OPEN, MERGED, DECLINED)
- `limit` (optional): Maximum number of results (1-100, default: 25)

### get_pull_request
Get details of a specific pull request.

**Parameters:**
- `project_key` (optional): The project key (uses default if BITBUCKET_DEFAULT_PROJECT_KEY is set)
- `repo_slug` (required): The repository slug
- `pull_request_id` (required): The pull request ID

### get_pull_request_activity
Get activity for a pull request (comments, approvals, etc.).

**Parameters:**
- `project_key` (optional): The project key (uses default if BITBUCKET_DEFAULT_PROJECT_KEY is set)
- `repo_slug` (required): The repository slug
- `pull_request_id` (required): The pull request ID

### get_pull_request_diff
Get the raw diff for a pull request.

**Parameters:**
- `project_key` (optional): The project key (uses default if BITBUCKET_DEFAULT_PROJECT_KEY is set)
- `repo_slug` (required): The repository slug
- `pull_request_id` (required): The pull request ID
- `context_lines` (optional): Number of context lines around changes
- `whitespace` (optional): Whitespace handling ('ignore-all', 'ignore-space-at-eol', 'ignore-space-change', 'ignore-trailing-space')
- `since` (optional): Base commit hash to diff from
- `until` (optional): End commit hash to diff to

### create_pull_request_comment
Add a comment to a pull request.

**Parameters:**
- `project_key` (optional): The project key (uses default if BITBUCKET_DEFAULT_PROJECT_KEY is set)
- `repo_slug` (required): The repository slug
- `pull_request_id` (required): The pull request ID
- `text` (required): The comment text
- `anchor_json` (optional): JSON-encoded anchor object for inline comments with properties:
  - `line`: Line number for inline comment
  - `line_type`: Line type (ADDED, REMOVED, CONTEXT)
  - `path`: File path for inline comment
  - `file_type`: File type (FROM, TO)
  - `from_hash`: Source commit hash
  - `to_hash`: Target commit hash
  - `src_path`: Source file path (for renames)
  - `dst_path`: Destination file path (for renames)
  - `diff_type`: Diff type (EFFECTIVE, RANGE, COMMIT) - auto-set to RANGE when commit hashes provided
  - `orphaned_type`: Orphaned comment type

### create_pull_request
Create a new pull request.

**Parameters:**
- `project_key` (optional): The project key (uses default if BITBUCKET_DEFAULT_PROJECT_KEY is set)
- `repo_slug` (required): The repository slug
- `title` (required): The pull request title
- `from_branch` (required): Source branch name
- `to_branch` (required): Target branch name
- `description` (optional): The pull request description

### approve_pull_request
Approve a pull request.

**Parameters:**
- `project_key` (optional): The project key (uses default if BITBUCKET_DEFAULT_PROJECT_KEY is set)
- `repo_slug` (required): The repository slug
- `pull_request_id` (required): The pull request ID

### unapprove_pull_request
Remove approval from a pull request.

**Parameters:**
- `project_key` (optional): The project key (uses default if BITBUCKET_DEFAULT_PROJECT_KEY is set)
- `repo_slug` (required): The repository slug
- `pull_request_id` (required): The pull request ID

### merge_pull_request
Merge a pull request (automatically fetches current version for optimistic locking).

**Parameters:**
- `project_key` (optional): The project key (uses default if BITBUCKET_DEFAULT_PROJECT_KEY is set)
- `repo_slug` (required): The repository slug
- `pull_request_id` (required): The pull request ID

### decline_pull_request
Decline a pull request (automatically fetches current version for optimistic locking).

**Parameters:**
- `project_key` (optional): The project key (uses default if BITBUCKET_DEFAULT_PROJECT_KEY is set)
- `repo_slug` (required): The repository slug
- `pull_request_id` (required): The pull request ID

### get_repos
Get a list of repositories in a project.

**Parameters:**
- `project_key` (optional): The project key (uses default if BITBUCKET_DEFAULT_PROJECT_KEY is set)
- `limit` (optional): Maximum number of results to return (default: 25)
- `start` (optional): Starting index for pagination (default: 0)

### get_pull_request_settings
Get pull request configuration settings for a repository.

**Parameters:**
- `project_key` (optional): The project key (uses default if BITBUCKET_DEFAULT_PROJECT_KEY is set)
- `repo_slug` (required): The repository slug

## Usage with MCP Clients

This server communicates via STDIO using the Model Context Protocol. It can be used with any MCP-compatible client such as Claude Desktop or VS Code with MCP support.

### Example Claude Desktop Configuration

```json
{
  "mcpServers": {
    "filesystem": {
      "command": "npx",
      "args": [
        "-y",
        "@modelcontextprotocol/server-filesystem",
        "/path/to/allow"
      ]
    },
    "bbcli": {
      "command": "/Users/mdeangelo/projects/bbcli/bbcli",
      "args": [],
      "env": {
        "BITBUCKET_BASE_URL": "http://localhost:7990",
        "BITBUCKET_USERNAME": "your-username",
        "BITBUCKET_PASSWORD": "your-app-password",
        "BITBUCKET_DEFAULT_PROJECT_KEY": "MYPROJ"
      }
    }
  }
}
```

## Implementation Notes

- **Built with mcp-go**: Uses the official mcp-go library for robust MCP protocol implementation
- **Automatic version management**: Merge and decline operations automatically fetch the current PR version to prevent optimistic locking conflicts
- **Simplified anchor handling**: Inline comment anchors are passed as JSON strings for easier client integration
- **Error handling**: Comprehensive error handling with descriptive messages

## Security

- Uses HTTP Basic Authentication with Bitbucket Server
- Requires valid Bitbucket Server credentials
- All API requests are made over HTTPS (when configured)
