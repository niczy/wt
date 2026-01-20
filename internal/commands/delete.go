package commands

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/niczy/wt/internal/fuzzy"
	"github.com/niczy/wt/internal/git"
)

// Delete handles the -d flag to delete a worktree
func Delete(pattern string) error {
	worktrees, err := git.ListWorktrees()
	if err != nil {
		return err
	}

	if len(worktrees) == 0 {
		return fmt.Errorf("no worktrees found in WT_HOME")
	}

	matches := fuzzy.FuzzyMatch(pattern, worktrees)
	if len(matches) == 0 {
		return fmt.Errorf("no worktree matching '%s' found", pattern)
	}

	var selected string
	if len(matches) == 1 {
		selected = matches[0].Text
	} else {
		// Multiple matches, ask user to choose
		selected, err = promptDeleteSelection(matches)
		if err != nil {
			return err
		}
	}

	wtHome, err := git.GetWTHome()
	if err != nil {
		return err
	}

	targetPath := filepath.Join(wtHome, selected)

	// Confirm deletion
	fmt.Fprintf(os.Stderr, "Delete worktree '%s'? [y/N]: ", selected)
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return fmt.Errorf("failed to read input: %w", err)
	}

	input = strings.TrimSpace(strings.ToLower(input))
	if input != "y" && input != "yes" {
		fmt.Println("Deletion cancelled")
		return nil
	}

	// Remove the worktree
	if err := git.RemoveWorktree(targetPath); err != nil {
		// If git worktree remove fails, try to remove the directory manually
		fmt.Fprintf(os.Stderr, "Warning: git worktree remove failed, attempting manual removal: %v\n", err)
		if err := os.RemoveAll(targetPath); err != nil {
			return fmt.Errorf("failed to remove worktree directory: %w", err)
		}
	}

	fmt.Printf("Deleted worktree: %s\n", selected)
	return nil
}

// promptDeleteSelection prompts the user to select from multiple matches for deletion
func promptDeleteSelection(matches []fuzzy.Match) (string, error) {
	fmt.Fprintf(os.Stderr, "Multiple matches found:\n")
	for i, match := range matches {
		fmt.Fprintf(os.Stderr, "  [%d] %s\n", i+1, match.Text)
	}
	fmt.Fprintf(os.Stderr, "Enter selection to delete (1-%d): ", len(matches))

	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("failed to read input: %w", err)
	}

	input = strings.TrimSpace(input)
	selection, err := strconv.Atoi(input)
	if err != nil || selection < 1 || selection > len(matches) {
		return "", fmt.Errorf("invalid selection: %s", input)
	}

	return matches[selection-1].Text, nil
}
