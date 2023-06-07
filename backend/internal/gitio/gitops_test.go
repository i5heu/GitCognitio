package gitio

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func TestGetRepoStats(t *testing.T) {
	// Create a temporary directory for the test repository
	dir, err := ioutil.TempDir("", "gitops_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %s", err)
	}
	defer os.RemoveAll(dir) // clean up

	// Initialize a new repository
	repo, err := git.PlainInit(dir, false)
	if err != nil {
		t.Fatalf("Failed to initialize repository: %s", err)
	}

	// Create a new file in the repository
	filePath := filepath.Join(dir, "test.txt")
	err = ioutil.WriteFile(filePath, []byte("Test content"), 0644)
	if err != nil {
		t.Fatalf("Failed to write file: %s", err)
	}

	// Create a second file in the repository
	filePath2 := filepath.Join(dir, "test2.md")
	err = ioutil.WriteFile(filePath2, []byte("Test content"), 0644)
	if err != nil {
		t.Fatalf("Failed to write file: %s", err)
	}

	// Add the file to the repository
	w, err := repo.Worktree()
	if err != nil {
		t.Fatalf("Failed to get worktree: %s", err)
	}
	_, err = w.Add("test.txt")
	if err != nil {
		t.Fatalf("Failed to add file: %s", err)
	}

	// Commit the change
	_, err = w.Commit("Initial commit", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Test User",
			Email: "test@example.com",
			When:  time.Now(),
		},
	})
	if err != nil {
		t.Fatalf("Failed to commit: %s", err)
	}

	// Create the RepoManager
	rm := &RepoManager{Repo: repo}

	// Call the function we're testing
	stats, err := rm.GetRepoStats()
	if err != nil {
		t.Fatalf("Failed to get repository stats: %s", err)
	}

	// Check that the stats match what we expect
	if stats.CommitCount != 1 {
		t.Errorf("Expected 1 commit, got %d", stats.CommitCount)
	}
	if stats.FileCount != 2 {
		t.Errorf("Expected 1 file, got %d", stats.FileCount)
	}
	if stats.RepoSize <= 0 {
		t.Errorf("Expected repo size to be greater than 0, got %d", stats.RepoSize)
	}
	if stats.TextFilesSize <= 0 {
		t.Errorf("Expected text file size to be greater than 0, got %d", stats.TextFilesSize)
	}
	if stats.GitFolderSize <= 0 {
		t.Errorf("Expected .git directory size to be greater than 0, got %d", stats.GitFolderSize)
	}
}
