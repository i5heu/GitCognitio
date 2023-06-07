package actions

import (
	"github.com/i5heu/GitCognitio/internal/gitio"
)

func NewMdFile(content string, path string, rm *gitio.RepoManager) error {
	err := gitio.CreateFile(path, content)
	if err != nil {
		return err
	}

	err = rm.Commit("New file: " + path)
	if err != nil {
		return err
	}

	rm.PushNonBlock()
	return nil
}
