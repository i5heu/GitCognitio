package gitio

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
)

// RepoManager is a structure that encapsulates the Git repository operations.
type RepoManager struct {
	Repo        *git.Repository
	auth        ssh.AuthMethod
	PushChannel chan bool
}

// RepoStats holds statistics for a Git repository.
type RepoStats struct {
	CommitCount   int
	BranchCount   int
	TagCount      int
	FileCount     int
	RepoSize      int64
	TextFilesSize int64
	GitFolderSize int64
}

// NewRepoManager initializes a new RepoManager instance.

func NewRepoManager(repoURL, path, sshPath string) (*RepoManager, error) {
	var repo *git.Repository
	var err error
	var sshAuth ssh.AuthMethod
	pushChannel := make(chan bool, 10)

	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		// Clone the repository

		// Create ssh auth
		sshAuth, err = ssh.NewPublicKeysFromFile("git", sshPath, "")
		if err != nil {
			return nil, err
		}

		repo, err = git.PlainClone(path, false, &git.CloneOptions{
			URL:  repoURL,
			Auth: sshAuth,
		})
	} else {
		// Open an existing repository
		repo, err = git.PlainOpen(path)
		if err == nil {
			// Pull changes from the remote
			w, err := repo.Worktree()
			if err != nil {
				return nil, err
			}

			// Create ssh auth
			sshAuth, err = ssh.NewPublicKeysFromFile("git", sshPath, "")
			if err != nil {
				return nil, err
			}

			err = w.Pull(&git.PullOptions{
				RemoteName: "origin",
				Auth:       sshAuth,
			})
			if err != nil && err != git.NoErrAlreadyUpToDate {
				return nil, err
			}
		}
	}

	if err != nil {
		return nil, err
	}

	return &RepoManager{Repo: repo, auth: sshAuth, PushChannel: pushChannel}, nil
}

// Commit commits changes to the repository.
func (rm *RepoManager) Commit(message string) error {
	w, err := rm.Repo.Worktree()
	if err != nil {
		return err
	}

	_, err = w.Add(".")
	if err != nil {
		return err
	}

	commit, err := w.Commit(message, &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Your Name",
			Email: "your.email@example.com",
			When:  time.Now(),
		},
	})
	if err != nil {
		return err
	}

	_, err = rm.Repo.CommitObject(commit)
	return err
}

// Pull pulls changes from the remote repository.
func (rm *RepoManager) Pull() error {
	w, err := rm.Repo.Worktree()
	if err != nil {
		return err
	}

	err = w.Pull(&git.PullOptions{
		RemoteName: "origin",
		Auth:       rm.auth,
	})
	if err != nil && err != git.NoErrAlreadyUpToDate {
		return err
	}

	return nil
}

// Push pushes changes to the remote repository.
func (rm *RepoManager) Push() error {
	err := rm.Repo.Push(&git.PushOptions{
		RemoteName: "origin",
		Auth:       rm.auth,
	})
	if err != nil && err != git.NoErrAlreadyUpToDate {
		return err
	}

	return nil
}

// GetRepoStats returns the statistics for the Git repository.
func (rm *RepoManager) GetRepoStats() (*RepoStats, error) {
	stats := &RepoStats{}

	// Count commits
	commitIter, err := rm.Repo.Log(&git.LogOptions{All: true})
	if err != nil {
		return nil, err
	}
	err = commitIter.ForEach(func(c *object.Commit) error {
		stats.CommitCount++
		return nil
	})
	if err != nil {
		return nil, err
	}

	// Count branches
	branchIter, err := rm.Repo.Branches()
	if err != nil {
		return nil, err
	}
	err = branchIter.ForEach(func(ref *plumbing.Reference) error {
		stats.BranchCount++
		return nil
	})
	if err != nil {
		return nil, err
	}

	// Count tags
	tagIter, err := rm.Repo.Tags()
	if err != nil {
		return nil, err
	}
	err = tagIter.ForEach(func(ref *plumbing.Reference) error {
		stats.TagCount++
		return nil
	})
	if err != nil {
		return nil, err
	}

	w, err := rm.Repo.Worktree()
	if err != nil {
		return nil, err
	}

	// Count and size of files, and size of text files
	err = filepath.Walk(w.Filesystem.Root(), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// Ignore .git directory
		if !info.IsDir() && !strings.HasPrefix(path, filepath.Join(w.Filesystem.Root(), ".git")) {
			stats.FileCount++
			stats.RepoSize += info.Size()

			// Check if the file is a text file and markdown files
			if strings.HasSuffix(info.Name(), ".md") ||
				strings.HasSuffix(info.Name(), ".markdown") ||
				strings.HasSuffix(info.Name(), ".txt") {

				stats.TextFilesSize += info.Size()
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	// Size of .git directory
	err = filepath.Walk(filepath.Join(w.Filesystem.Root(), ".git"), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			stats.GitFolderSize += info.Size()
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return stats, nil
}

func (rm *RepoManager) StartPushListener() {
	go func() {
		for range rm.PushChannel {
			fmt.Println("Pushing to remote")
			// Perform the push operation
			err := rm.Push()
			if err != nil {
				// Handle the error appropriately
				fmt.Println("Error pushing:", err)
			} else {
				fmt.Println("Pushed")
			}
		}
	}()
}

func (rm *RepoManager) PushNonBlock() {
	go func() {
		select {
		case rm.PushChannel <- true:
			// Value pushed to the channel
		default:
			// Channel is not ready to receive the value
			fmt.Println("Channel is not ready to receive the value")
		}
	}()
}
