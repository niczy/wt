package git

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// GetWTHome returns the WT_HOME directory, defaulting to ~/worktrees
func GetWTHome() (string, error) {
	wtHome := os.Getenv("WT_HOME")
	if wtHome == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to get home directory: %w", err)
		}
		wtHome = filepath.Join(home, "worktrees")
	}
	return wtHome, nil
}

// GetRepoName returns the name of the current git repository
func GetRepoName() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("not in a git repository: %w", err)
	}
	repoPath := strings.TrimSpace(string(output))
	return filepath.Base(repoPath), nil
}

// GetRepoRoot returns the root path of the current git repository
func GetRepoRoot() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("not in a git repository: %w", err)
	}
	return strings.TrimSpace(string(output)), nil
}

// GetCurrentBranch returns the current branch name
func GetCurrentBranch() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get current branch: %w", err)
	}
	return strings.TrimSpace(string(output)), nil
}

// CreateWorktree creates a new git worktree at the specified path
func CreateWorktree(targetPath, branchName string) error {
	// First, try to create the worktree with an existing branch
	cmd := exec.Command("git", "worktree", "add", targetPath, branchName)
	if err := cmd.Run(); err != nil {
		// If branch doesn't exist, create a new branch
		cmd = exec.Command("git", "worktree", "add", "-b", branchName, targetPath)
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to create worktree: %w", err)
		}
	}
	return nil
}

// RemoveWorktree removes a git worktree
func RemoveWorktree(worktreePath string) error {
	cmd := exec.Command("git", "worktree", "remove", worktreePath, "--force")
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to remove worktree: %s", string(output))
	}
	return nil
}

// ListWorktrees returns a list of all worktree paths under WT_HOME
func ListWorktrees() ([]string, error) {
	wtHome, err := GetWTHome()
	if err != nil {
		return nil, err
	}

	var worktrees []string
	entries, err := os.ReadDir(wtHome)
	if err != nil {
		if os.IsNotExist(err) {
			return worktrees, nil
		}
		return nil, fmt.Errorf("failed to read WT_HOME directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			worktreePath := filepath.Join(wtHome, entry.Name())
			// Check if it's a valid git worktree
			gitDir := filepath.Join(worktreePath, ".git")
			if _, err := os.Stat(gitDir); err == nil {
				worktrees = append(worktrees, entry.Name())
			}
		}
	}

	return worktrees, nil
}
