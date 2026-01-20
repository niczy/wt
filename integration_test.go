// +build integration

package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// Integration tests require a real git repository environment
// Run with: go test -tags=integration ./...

func TestIntegration_CreateListNavigateDelete(t *testing.T) {
	// Skip if not in a git repository
	if _, err := exec.Command("git", "rev-parse", "--git-dir").Output(); err != nil {
		t.Skip("skipping integration test: not in a git repository")
	}

	// Create a temporary WT_HOME
	tmpDir, err := os.MkdirTemp("", "wt-integration-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Set WT_HOME environment variable
	originalWTHome := os.Getenv("WT_HOME")
	os.Setenv("WT_HOME", tmpDir)
	defer func() {
		if originalWTHome != "" {
			os.Setenv("WT_HOME", originalWTHome)
		} else {
			os.Unsetenv("WT_HOME")
		}
	}()

	// Build the wt binary
	binaryPath := filepath.Join(tmpDir, "wt")
	buildCmd := exec.Command("go", "build", "-o", binaryPath, ".")
	if output, err := buildCmd.CombinedOutput(); err != nil {
		t.Fatalf("failed to build wt: %v\n%s", err, output)
	}

	worktreeName := "test-integration-wt"

	// Test create
	t.Run("Create", func(t *testing.T) {
		cmd := exec.Command(binaryPath, "-c", worktreeName)
		cmd.Env = append(os.Environ(), "WT_HOME="+tmpDir)
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("create failed: %v\n%s", err, output)
		}
		if !strings.Contains(string(output), "Created worktree") {
			t.Errorf("expected 'Created worktree' message, got: %s", output)
		}
	})

	// Test list
	t.Run("List", func(t *testing.T) {
		cmd := exec.Command(binaryPath, "-l")
		cmd.Env = append(os.Environ(), "WT_HOME="+tmpDir)
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("list failed: %v\n%s", err, output)
		}
		if !strings.Contains(string(output), worktreeName) {
			t.Errorf("expected worktree '%s' in list, got: %s", worktreeName, output)
		}
	})

	// Test navigate (fuzzy search)
	t.Run("Navigate", func(t *testing.T) {
		cmd := exec.Command(binaryPath, "integration")
		cmd.Env = append(os.Environ(), "WT_HOME="+tmpDir)
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("navigate failed: %v\n%s", err, output)
		}
		if !strings.Contains(string(output), "WT_CD_PATH=") {
			t.Errorf("expected 'WT_CD_PATH=' in output, got: %s", output)
		}
	})

	// Test delete
	t.Run("Delete", func(t *testing.T) {
		cmd := exec.Command(binaryPath, "-d", "integration")
		cmd.Env = append(os.Environ(), "WT_HOME="+tmpDir)
		cmd.Stdin = strings.NewReader("y\n")
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("delete failed: %v\n%s", err, output)
		}
		if !strings.Contains(string(output), "Deleted worktree") {
			t.Errorf("expected 'Deleted worktree' message, got: %s", output)
		}
	})

	// Verify worktree was deleted
	t.Run("VerifyDeleted", func(t *testing.T) {
		cmd := exec.Command(binaryPath, "-l")
		cmd.Env = append(os.Environ(), "WT_HOME="+tmpDir)
		output, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("list failed: %v\n%s", err, output)
		}
		if strings.Contains(string(output), worktreeName) {
			t.Errorf("worktree '%s' should have been deleted, but found in: %s", worktreeName, output)
		}
	})
}

func TestIntegration_Help(t *testing.T) {
	// Build the wt binary
	tmpDir, err := os.MkdirTemp("", "wt-integration-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	binaryPath := filepath.Join(tmpDir, "wt")
	buildCmd := exec.Command("go", "build", "-o", binaryPath, ".")
	if output, err := buildCmd.CombinedOutput(); err != nil {
		t.Fatalf("failed to build wt: %v\n%s", err, output)
	}

	cmd := exec.Command(binaryPath, "-h")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("help failed: %v\n%s", err, output)
	}

	expectedStrings := []string{
		"wt - Git Worktree Manager",
		"wt -c <name>",
		"wt -d <name>",
		"wt -l",
		"WT_HOME",
	}

	for _, expected := range expectedStrings {
		if !strings.Contains(string(output), expected) {
			t.Errorf("expected '%s' in help output, got: %s", expected, output)
		}
	}
}

func TestIntegration_ErrorCases(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "wt-integration-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	binaryPath := filepath.Join(tmpDir, "wt")
	buildCmd := exec.Command("go", "build", "-o", binaryPath, ".")
	if output, err := buildCmd.CombinedOutput(); err != nil {
		t.Fatalf("failed to build wt: %v\n%s", err, output)
	}

	t.Run("NavigateNoMatch", func(t *testing.T) {
		cmd := exec.Command(binaryPath, "nonexistent-pattern-xyz")
		cmd.Env = append(os.Environ(), "WT_HOME="+tmpDir)
		output, err := cmd.CombinedOutput()
		// Should exit with error
		if err == nil {
			t.Error("expected error for non-matching navigate")
		}
		if !strings.Contains(string(output), "no worktrees found") {
			t.Errorf("expected 'no worktrees found' error, got: %s", output)
		}
	})

	t.Run("DeleteNoMatch", func(t *testing.T) {
		cmd := exec.Command(binaryPath, "-d", "nonexistent-pattern-xyz")
		cmd.Env = append(os.Environ(), "WT_HOME="+tmpDir)
		output, err := cmd.CombinedOutput()
		// Should exit with error
		if err == nil {
			t.Error("expected error for non-matching delete")
		}
		if !strings.Contains(string(output), "no worktrees found") {
			t.Errorf("expected 'no worktrees found' error, got: %s", output)
		}
	})

	t.Run("CreateNotInGitRepo", func(t *testing.T) {
		// Create a non-git directory
		nonGitDir := filepath.Join(tmpDir, "not-a-repo")
		os.MkdirAll(nonGitDir, 0755)

		cmd := exec.Command(binaryPath, "-c", "test")
		cmd.Dir = nonGitDir
		cmd.Env = append(os.Environ(), "WT_HOME="+tmpDir)
		output, err := cmd.CombinedOutput()
		// Should exit with error
		if err == nil {
			t.Error("expected error when creating worktree outside git repo")
		}
		if !strings.Contains(string(output), "not in a git repository") {
			t.Errorf("expected 'not in a git repository' error, got: %s", output)
		}
	})
}
