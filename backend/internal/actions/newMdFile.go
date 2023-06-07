package actions

import (
	"fmt"

	"github.com/i5heu/GitCognitio/internal/gitio"
	"github.com/i5heu/GitCognitio/types"
)

func NewMdFile(message types.Message, broadcastChannel *chan types.Message, rm *gitio.RepoManager) {

	path := "messages/" + message.ID + ".md"

	err := gitio.CreateFile(path, message.Data)
	if err != nil {
		fmt.Println("error creating file", err)
		broadcastMessage(broadcastChannel, types.Message{
			ID:   message.ID,
			Type: "error",
			Data: "error creating file",
		})

		return
	}

	err = rm.Commit("New file: " + path)
	if err != nil {
		fmt.Println("error committing file", err)
		broadcastMessage(broadcastChannel, types.Message{
			ID:   message.ID,
			Type: "error",
			Data: "error creating file",
		})

		return
	}

	rm.PushNonBlock()

	broadcastMessage(broadcastChannel, types.Message{
		ID:   message.ID,
		Type: "message",
		Path: path,
		Data: message.Data,
	})
}
