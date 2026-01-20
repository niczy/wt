package commands

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Helper function to temporarily set WT_HOME
func withWTHome(t *testing.T, path string, fn func()) {
	original := os.Getenv("WT_HOME")
	os.Setenv("WT_HOME", path)
	defer func() {
		if original != "" {
			os.Setenv("WT_HOME", original)
		} else {
			os.Unsetenv("WT_HOME")
		}
	}()
	fn()
}

// Helper function to create a mock worktree directory
func createMockWorktree(t *testing.T, wtHome, name string) string {
	wtPath := filepath.Join(wtHome, name)
	gitPath := filepath.Join(wtPath, ".git")
	if err := os.MkdirAll(gitPath, 0755); err != nil {
		t.Fatalf("failed to create mock worktree: %v", err)
	}
	return wtPath
}

// Helper to capture stdout
func captureOutput(fn func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	fn()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	return buf.String()
}

func TestList_EmptyDir(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "wt-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	withWTHome(t, tmpDir, func() {
		output := captureOutput(func() {
			err := List()
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
		if !strings.Contains(output, "No worktrees found") {
			t.Errorf("expected 'No worktrees found' message, got: %s", output)
		}
	})
}

func TestList_WithWorktrees(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "wt-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create mock worktrees
	createMockWorktree(t, tmpDir, "repo-feature")
	createMockWorktree(t, tmpDir, "repo-bugfix")

	withWTHome(t, tmpDir, func() {
		output := captureOutput(func() {
			err := List()
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
		if !strings.Contains(output, "repo-feature") {
			t.Errorf("expected 'repo-feature' in output, got: %s", output)
		}
		if !strings.Contains(output, "repo-bugfix") {
			t.Errorf("expected 'repo-bugfix' in output, got: %s", output)
		}
	})
}

func TestNavigate_SingleMatch(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "wt-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	createMockWorktree(t, tmpDir, "repo-feature")
	createMockWorktree(t, tmpDir, "repo-bugfix")

	withWTHome(t, tmpDir, func() {
		output := captureOutput(func() {
			err := Navigate("feature")
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
		expectedPath := filepath.Join(tmpDir, "repo-feature")
		if !strings.Contains(output, "WT_CD_PATH="+expectedPath) {
			t.Errorf("expected WT_CD_PATH=%s, got: %s", expectedPath, output)
		}
	})
}

func TestNavigate_NoMatch(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "wt-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	createMockWorktree(t, tmpDir, "repo-feature")

	withWTHome(t, tmpDir, func() {
		err := Navigate("nonexistent")
		if err == nil {
			t.Error("expected error for non-matching pattern")
		}
		if !strings.Contains(err.Error(), "no worktree matching") {
			t.Errorf("expected 'no worktree matching' error, got: %v", err)
		}
	})
}

func TestNavigate_EmptyWorktrees(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "wt-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	withWTHome(t, tmpDir, func() {
		err := Navigate("something")
		if err == nil {
			t.Error("expected error for empty worktrees")
		}
		if !strings.Contains(err.Error(), "no worktrees found") {
			t.Errorf("expected 'no worktrees found' error, got: %v", err)
		}
	})
}

func TestCreate_NotInGitRepo(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "wt-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Change to a non-git directory
	originalDir, _ := os.Getwd()
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to change to temp dir: %v", err)
	}
	defer func() { _ = os.Chdir(originalDir) }()

	withWTHome(t, tmpDir, func() {
		err := Create("test-worktree")
		if err == nil {
			t.Error("expected error when not in git repo")
		}
		if !strings.Contains(err.Error(), "not in a git repository") {
			t.Errorf("expected 'not in a git repository' error, got: %v", err)
		}
	})
}

func TestDelete_NoMatch(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "wt-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	createMockWorktree(t, tmpDir, "repo-feature")

	withWTHome(t, tmpDir, func() {
		err := Delete("nonexistent")
		if err == nil {
			t.Error("expected error for non-matching pattern")
		}
		if !strings.Contains(err.Error(), "no worktree matching") {
			t.Errorf("expected 'no worktree matching' error, got: %v", err)
		}
	})
}

func TestDelete_EmptyWorktrees(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "wt-test-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	withWTHome(t, tmpDir, func() {
		err := Delete("something")
		if err == nil {
			t.Error("expected error for empty worktrees")
		}
		if !strings.Contains(err.Error(), "no worktrees found") {
			t.Errorf("expected 'no worktrees found' error, got: %v", err)
		}
	})
}
