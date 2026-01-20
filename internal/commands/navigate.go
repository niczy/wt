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

// Navigate handles the default command to enter a worktree directory
func Navigate(pattern string) error {
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
		selected, err = promptSelection(matches)
		if err != nil {
			return err
		}
	}

	wtHome, err := git.GetWTHome()
	if err != nil {
		return err
	}

	targetPath := filepath.Join(wtHome, selected)

	// Print the path for shell integration to capture
	// The shell wrapper function will read this and cd to the path
	fmt.Printf("WT_CD_PATH=%s\n", targetPath)

	return nil
}

// promptSelection prompts the user to select from multiple matches
func promptSelection(matches []fuzzy.Match) (string, error) {
	fmt.Fprintf(os.Stderr, "Multiple matches found:\n")
	for i, match := range matches {
		fmt.Fprintf(os.Stderr, "  [%d] %s\n", i+1, match.Text)
	}
	fmt.Fprintf(os.Stderr, "Enter selection (1-%d): ", len(matches))

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
