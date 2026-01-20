package git

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetWTHome_Default(t *testing.T) {
	// Unset WT_HOME to test default
	original := os.Getenv("WT_HOME")
	os.Unsetenv("WT_HOME")
	defer func() {
		if original != "" {
			os.Setenv("WT_HOME", original)
		}
	}()

	wtHome, err := GetWTHome()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	home, _ := os.UserHomeDir()
	expected := filepath.Join(home, "worktrees")
	if wtHome != expected {
		t.Errorf("expected %s, got %s", expected, wtHome)
	}
}

func TestGetWTHome_EnvVar(t *testing.T) {
	original := os.Getenv("WT_HOME")
	defer func() {
		if original != "" {
			os.Setenv("WT_HOME", original)
		} else {
			os.Unsetenv("WT_HOME")
		}
	}()

	expected := "/custom/worktrees/path"
	os.Setenv("WT_HOME", expected)

	wtHome, err := GetWTHome()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if wtHome != expected {
		t.Errorf("expected %s, got %s", expected, wtHome)
	}
}

func TestGetRepoName_InGitRepo(t *testing.T) {
	// This test assumes we're running in a git repository
	name, err := GetRepoName()
	if err != nil {
		t.Skipf("skipping test, not in a git repository: %v", err)
	}

	if name == "" {
		t.Error("expected non-empty repo name")
	}
}

func TestGetRepoRoot_InGitRepo(t *testing.T) {
	// This test assumes we're running in a git repository
	root, err := GetRepoRoot()
	if err != nil {
		t.Skipf("skipping test, not in a git repository: %v", err)
	}

	if root == "" {
		t.Error("expected non-empty repo root")
	}

	// Root should be a directory
	info, err := os.Stat(root)
	if err != nil {
		t.Fatalf("repo root does not exist: %v", err)
	}
	if !info.IsDir() {
		t.Error("repo root should be a directory")
	}
}

func TestGetCurrentBranch_InGitRepo(t *testing.T) {
	// This test assumes we're running in a git repository
	branch, err := GetCurrentBranch()
	if err != nil {
		t.Skipf("skipping test, not in a git repository: %v", err)
	}

	if branch == "" {
		t.Error("expected non-empty branch name")
	}
}

func TestListWorktrees_EmptyDir(t *testing.T) {
	// Create a temporary directory for WT_HOME
	tmpDir, err := os.MkdirTemp("", "wt-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Set WT_HOME to temp directory
	original := os.Getenv("WT_HOME")
	os.Setenv("WT_HOME", tmpDir)
	defer func() {
		if original != "" {
			os.Setenv("WT_HOME", original)
		} else {
			os.Unsetenv("WT_HOME")
		}
	}()

	worktrees, err := ListWorktrees()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(worktrees) != 0 {
		t.Errorf("expected empty worktree list, got %d", len(worktrees))
	}
}

func TestListWorktrees_NonExistentDir(t *testing.T) {
	// Set WT_HOME to non-existent directory
	original := os.Getenv("WT_HOME")
	os.Setenv("WT_HOME", "/nonexistent/path/that/does/not/exist")
	defer func() {
		if original != "" {
			os.Setenv("WT_HOME", original)
		} else {
			os.Unsetenv("WT_HOME")
		}
	}()

	worktrees, err := ListWorktrees()
	if err != nil {
		t.Fatalf("unexpected error for non-existent dir: %v", err)
	}

	if len(worktrees) != 0 {
		t.Errorf("expected empty worktree list for non-existent dir, got %d", len(worktrees))
	}
}

func TestListWorktrees_WithMockWorktrees(t *testing.T) {
	// Create a temporary directory for WT_HOME
	tmpDir, err := os.MkdirTemp("", "wt-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create mock worktree directories with .git
	mockWorktrees := []string{"repo-feature", "repo-bugfix"}
	for _, wt := range mockWorktrees {
		wtPath := filepath.Join(tmpDir, wt)
		gitPath := filepath.Join(wtPath, ".git")
		if err := os.MkdirAll(gitPath, 0755); err != nil {
			t.Fatalf("failed to create mock worktree: %v", err)
		}
	}

	// Create a non-worktree directory (no .git)
	nonWtPath := filepath.Join(tmpDir, "not-a-worktree")
	if err := os.MkdirAll(nonWtPath, 0755); err != nil {
		t.Fatalf("failed to create non-worktree dir: %v", err)
	}

	// Set WT_HOME to temp directory
	original := os.Getenv("WT_HOME")
	os.Setenv("WT_HOME", tmpDir)
	defer func() {
		if original != "" {
			os.Setenv("WT_HOME", original)
		} else {
			os.Unsetenv("WT_HOME")
		}
	}()

	worktrees, err := ListWorktrees()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(worktrees) != 2 {
		t.Errorf("expected 2 worktrees, got %d", len(worktrees))
	}

	// Verify we got the right worktrees
	found := make(map[string]bool)
	for _, wt := range worktrees {
		found[wt] = true
	}
	for _, expected := range mockWorktrees {
		if !found[expected] {
			t.Errorf("expected worktree '%s' not found", expected)
		}
	}
	if found["not-a-worktree"] {
		t.Error("non-worktree directory should not be listed")
	}
}
