package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/nicholas/wt/internal/git"
)

// Create handles the -c flag to create a new worktree
func Create(worktreeName string) error {
	// Get WT_HOME
	wtHome, err := git.GetWTHome()
	if err != nil {
		return err
	}

	// Ensure WT_HOME exists
	if err := os.MkdirAll(wtHome, 0755); err != nil {
		return fmt.Errorf("failed to create WT_HOME directory: %w", err)
	}

	// Get current repo name
	repoName, err := git.GetRepoName()
	if err != nil {
		return err
	}

	// Construct target path: {WT_HOME}/{REPO_NAME}-{WORKTREE_NAME}
	targetDirName := fmt.Sprintf("%s-%s", repoName, worktreeName)
	targetPath := filepath.Join(wtHome, targetDirName)

	// Check if worktree already exists
	if _, err := os.Stat(targetPath); err == nil {
		return fmt.Errorf("worktree already exists at: %s", targetPath)
	}

	// Create the worktree with the worktree name as branch name
	if err := git.CreateWorktree(targetPath, worktreeName); err != nil {
		return err
	}

	fmt.Printf("Created worktree at: %s\n", targetPath)
	fmt.Printf("To enter the worktree, run: cd %s\n", targetPath)
	
	// Print the path for shell integration
	fmt.Printf("WT_CD_PATH=%s\n", targetPath)

	return nil
}
