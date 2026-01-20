package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/nicholas/wt/internal/commands"
)

const usage = `wt - Git Worktree Manager

Usage:
  wt -c <name>      Create a new worktree with the given name
  wt -d <name>      Delete a worktree (fuzzy search)
  wt -l             List all worktrees
  wt <name>         Navigate to a worktree (fuzzy search)

Environment Variables:
  WT_HOME           Directory where worktrees are stored (default: ~/worktrees)

Examples:
  wt -c feature-x   Create worktree at $WT_HOME/{repo}-feature-x
  wt feat           Navigate to worktree matching "feat"
  wt -d feature     Delete worktree matching "feature"

Shell Integration:
  To enable 'cd' functionality, add this to your shell config:

  # For bash/zsh:
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
`

func main() {
	createFlag := flag.String("c", "", "Create a new worktree with the given name")
	deleteFlag := flag.String("d", "", "Delete a worktree (fuzzy search)")
	listFlag := flag.Bool("l", false, "List all worktrees")
	helpFlag := flag.Bool("h", false, "Show help")

	flag.Usage = func() {
		fmt.Print(usage)
	}

	flag.Parse()

	if *helpFlag {
		flag.Usage()
		os.Exit(0)
	}

	var err error

	switch {
	case *createFlag != "":
		err = commands.Create(*createFlag)
	case *deleteFlag != "":
		err = commands.Delete(*deleteFlag)
	case *listFlag:
		err = commands.List()
	case flag.NArg() == 1:
		err = commands.Navigate(flag.Arg(0))
	case flag.NArg() == 0:
		err = commands.List()
	default:
		flag.Usage()
		os.Exit(1)
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
