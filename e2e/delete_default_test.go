// delete_default_test.go contains tests for default branch deletion protection.
package e2e

import (
	"strings"
	"testing"

	"github.com/k1LoW/exec"
	"github.com/k1LoW/git-wt/testutil"
)

func TestE2E_DeleteDefaultBranch(t *testing.T) {
	t.Parallel()
	binPath := buildBinary(t)

	t.Run("blocks_safe_delete_of_default_branch", func(t *testing.T) {
		t.Parallel()
		repo := testutil.NewTestRepo(t)
		repo.CreateFile("README.md", "# Test")
		repo.Commit("initial commit")

		// Create another branch to be on when trying to delete main
		cmd := exec.Command("git", "checkout", "-b", "other")
		cmd.Dir = repo.Root
		if err := cmd.Run(); err != nil {
			t.Fatalf("failed to create other branch: %v", err)
		}

		// Try to delete main branch with -d
		out, err := runGitWt(t, binPath, repo.Root, "-d", "main")
		if err == nil {
			t.Fatal("should fail when deleting default branch")
		}
		if !strings.Contains(out, "cannot delete default branch") {
			t.Errorf("error should mention default branch protection, got: %s", out)
		}
		if !strings.Contains(out, "--allow-delete-default") {
			t.Errorf("error should suggest --allow-delete-default, got: %s", out)
		}
	})

	t.Run("blocks_force_delete_of_default_branch", func(t *testing.T) {
		t.Parallel()
		repo := testutil.NewTestRepo(t)
		repo.CreateFile("README.md", "# Test")
		repo.Commit("initial commit")

		// Create another branch to be on when trying to delete main
		cmd := exec.Command("git", "checkout", "-b", "other")
		cmd.Dir = repo.Root
		if err := cmd.Run(); err != nil {
			t.Fatalf("failed to create other branch: %v", err)
		}

		// Try to delete main branch with -D
		out, err := runGitWt(t, binPath, repo.Root, "-D", "main")
		if err == nil {
			t.Fatal("should fail when force deleting default branch")
		}
		if !strings.Contains(out, "cannot delete default branch") {
			t.Errorf("error should mention default branch protection, got: %s", out)
		}
	})

	t.Run("allows_delete_with_override_flag", func(t *testing.T) {
		t.Parallel()
		repo := testutil.NewTestRepo(t)
		repo.CreateFile("README.md", "# Test")
		repo.Commit("initial commit")

		// Create another branch to be on when deleting main
		cmd := exec.Command("git", "checkout", "-b", "other")
		cmd.Dir = repo.Root
		if err := cmd.Run(); err != nil {
			t.Fatalf("failed to create other branch: %v", err)
		}

		// Delete main branch with override flag
		out, err := runGitWt(t, binPath, repo.Root, "-D", "--allow-delete-default", "main")
		if err != nil {
			t.Fatalf("should allow deleting default branch with override: %v\noutput: %s", err, out)
		}

		// Verify branch was deleted
		cmd = exec.Command("git", "branch", "--list", "main")
		cmd.Dir = repo.Root
		branchOut, err := cmd.Output()
		if err != nil {
			t.Fatalf("git branch --list failed: %v", err)
		}
		if strings.Contains(string(branchOut), "main") {
			t.Error("main branch should have been deleted")
		}
	})

	t.Run("blocks_default_branch_in_multiple_args", func(t *testing.T) {
		t.Parallel()
		repo := testutil.NewTestRepo(t)
		repo.CreateFile("README.md", "# Test")
		repo.Commit("initial commit")

		// Create other branches
		cmd := exec.Command("git", "branch", "feature-a")
		cmd.Dir = repo.Root
		if err := cmd.Run(); err != nil {
			t.Fatalf("failed to create feature-a branch: %v", err)
		}
		cmd = exec.Command("git", "branch", "feature-b")
		cmd.Dir = repo.Root
		if err := cmd.Run(); err != nil {
			t.Fatalf("failed to create feature-b branch: %v", err)
		}

		// Create another branch to be on when trying to delete
		cmd = exec.Command("git", "checkout", "-b", "other")
		cmd.Dir = repo.Root
		if err := cmd.Run(); err != nil {
			t.Fatalf("failed to create other branch: %v", err)
		}

		// Try to delete multiple branches including main
		out, err := runGitWt(t, binPath, repo.Root, "-D", "feature-a", "main", "feature-b")
		if err == nil {
			t.Fatal("should fail when deleting multiple branches including default")
		}
		if !strings.Contains(out, "cannot delete default branch") {
			t.Errorf("error should mention default branch protection, got: %s", out)
		}

		// feature-a should NOT be deleted (check happens before any deletion)
		cmd = exec.Command("git", "branch", "--list", "feature-a")
		cmd.Dir = repo.Root
		branchOut, err := cmd.Output()
		if err != nil {
			t.Fatalf("git branch --list failed: %v", err)
		}
		if !strings.Contains(string(branchOut), "feature-a") {
			t.Error("feature-a should still exist (deletion should not happen if default branch is in args)")
		}
	})

	t.Run("blocks_worktree_with_default_branch", func(t *testing.T) {
		t.Parallel()
		repo := testutil.NewTestRepo(t)
		repo.CreateFile("README.md", "# Test")
		repo.Commit("initial commit")

		// Create a worktree for main branch
		// First, checkout to another branch
		cmd := exec.Command("git", "checkout", "-b", "other")
		cmd.Dir = repo.Root
		if err := cmd.Run(); err != nil {
			t.Fatalf("failed to create other branch: %v", err)
		}

		// Create worktree for main
		out, err := runGitWt(t, binPath, repo.Root, "main")
		if err != nil {
			t.Fatalf("failed to create worktree for main: %v\noutput: %s", err, out)
		}

		// Try to delete the worktree with -d
		out, err = runGitWt(t, binPath, repo.Root, "-d", "main")
		if err == nil {
			t.Fatal("should fail when deleting worktree of default branch")
		}
		if !strings.Contains(out, "cannot delete default branch") {
			t.Errorf("error should mention default branch protection, got: %s", out)
		}
	})
}
