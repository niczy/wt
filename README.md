# wt - Git Worktree Manager

A CLI tool to create and manage git worktrees with fuzzy search.

## Installation

```bash
go install github.com/niczy/wt@latest
```

Or build from source:

```bash
git clone https://github.com/niczy/wt.git
cd wt
go build -o wt .
```

## Usage

```
wt -c <name>      Create a new worktree with the given name
wt -d <name>      Delete a worktree (fuzzy search)
wt -l             List all worktrees
wt <name>         Navigate to a worktree (fuzzy search)
```

### Examples

```bash
# Create a worktree for a feature branch
wt -c feature-x
# Creates: $WT_HOME/{repo-name}-feature-x

# Navigate to a worktree using fuzzy search
wt feat
# Will cd to the matching worktree

# Delete a worktree
wt -d feature
# Will fuzzy search and delete the matching worktree

# List all worktrees
wt -l
```

## Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `WT_HOME` | Directory where worktrees are stored | `~/worktrees` |

## Shell Integration

To enable the `cd` functionality when navigating to worktrees, add this function to your shell config (`~/.bashrc` or `~/.zshrc`):

```bash
wt() {
  local output
  output=$(command wt "$@")
  local exit_code=$?
  
  # Check if output contains WT_CD_PATH
  if echo "$output" | grep -q "^WT_CD_PATH="; then
    local target_path
    target_path=$(echo "$output" | grep "^WT_CD_PATH=" | cut -d= -f2)
    echo "$output" | grep -v "^WT_CD_PATH="
    cd "$target_path"
  else
    echo "$output"
  fi
  
  return $exit_code
}
```

## How It Works

- **Create (`-c`)**: Creates a git worktree at `$WT_HOME/{repo-name}-{worktree-name}`. It first tries to checkout an existing branch with the name, or creates a new branch if it doesn't exist.

- **Navigate**: Uses fuzzy search to find matching worktrees. If multiple matches are found, prompts the user to select one.

- **Delete (`-d`)**: Uses fuzzy search to find the worktree, confirms with the user, then removes both the worktree and its directory.

## Features

- **Fuzzy Search**: Quickly find worktrees by partial name match
- **Multi-match Selection**: When multiple worktrees match, choose interactively
- **Shell Integration**: Seamlessly `cd` into worktree directories
- **Confirmation on Delete**: Prevents accidental deletion of worktrees

## Development

### Running Tests

```bash
# Run unit tests
go test -v ./...

# Run unit tests with coverage
go test -v -race -coverprofile=coverage.out ./...

# Run integration tests (requires git repository)
go test -v -tags=integration ./...

# Run all tests
go test -v -tags=integration ./...
```

### Linting

```bash
# Run go vet
go vet ./...

# Run golangci-lint (if installed)
golangci-lint run
```

### Building

```bash
# Build for current platform
go build -o wt .

# Cross-compile for different platforms
GOOS=linux GOARCH=amd64 go build -o wt-linux-amd64 .
GOOS=darwin GOARCH=amd64 go build -o wt-darwin-amd64 .
GOOS=windows GOARCH=amd64 go build -o wt-windows-amd64.exe .
```

## CI/CD

This project uses GitHub Actions for continuous integration. The workflow runs:

- Unit tests on Ubuntu and macOS with Go 1.21 and 1.22
- Integration tests on Ubuntu and macOS
- Build verification on Ubuntu, macOS, and Windows
- Linting with golangci-lint

## License

MIT License
