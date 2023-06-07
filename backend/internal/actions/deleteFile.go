package actions

import (
	"fmt"

	"github.com/i5heu/GitCognitio/internal/gitio"
	"github.com/i5heu/GitCognitio/types"
)

func DeleteFile(message types.Message, broadcastChannel *chan types.Message, rm *gitio.RepoManager) {

	err := gitio.DeleteFile(message.Path)
	if err != nil {
		fmt.Println("error creating file", err)
		broadcastMessage(broadcastChannel, types.Message{
			ID:   message.ID,
			Type: "error",
			Data: "error creating file",
		})
	}

	err = rm.Commit("Delete file: " + message.Path)
	if err != nil {
		fmt.Println("error committing file", err)
		broadcastMessage(broadcastChannel, types.Message{
			ID:   message.ID,
			Type: "error",
			Data: "error creating file",
		})
	}

	rm.PushNonBlock()

	broadcastMessage(broadcastChannel, types.Message{
		ID:   message.ID,
		Type: "success",
		Path: message.Path,
	})
}
