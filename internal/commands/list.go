package commands

import (
	"fmt"
	"path/filepath"

	"github.com/nicholas/wt/internal/git"
)

// List shows all worktrees in WT_HOME
func List() error {
	wtHome, err := git.GetWTHome()
	if err != nil {
		return err
	}

	worktrees, err := git.ListWorktrees()
	if err != nil {
		return err
	}

	if len(worktrees) == 0 {
		fmt.Printf("No worktrees found in %s\n", wtHome)
		return nil
	}

	fmt.Printf("Worktrees in %s:\n", wtHome)
	for _, wt := range worktrees {
		fmt.Printf("  %s\n", filepath.Join(wtHome, wt))
	}

	return nil
}
