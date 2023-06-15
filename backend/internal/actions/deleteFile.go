package actions

import (
	"fmt"

	"github.com/i5heu/GitCognitio/internal/gitio"
	"github.com/i5heu/GitCognitio/types"
)

func DeleteFile(message types.Message, broadcastChannel *chan types.Message, rm *gitio.RepoManager) {

	err := gitio.DeleteFile(message.Path)
	if err != nil {
		fmt.Println("error deleting file", err)
		broadcastMessage(broadcastChannel, types.Message{
			ID:   message.ID,
			Type: "error",
			Data: "error deleting file",
		})

		return
	}

	err = rm.Commit("Delete file: " + message.Path)
	if err != nil {
		fmt.Println("error commit file delete", err)
		broadcastMessage(broadcastChannel, types.Message{
			ID:   message.ID,
			Type: "error",
			Data: "error commit file delete",
		})

		return
	}

	rm.PushNonBlock()

	broadcastMessage(broadcastChannel, types.Message{
		ID:   message.ID,
		Type: "thread-delete",
		Path: message.Path,
	})
}
