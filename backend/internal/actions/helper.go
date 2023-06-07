package actions

import "github.com/i5heu/GitCognitio/types"

func broadcastMessage(broadcastChannel *chan types.Message, message types.Message) {
	*broadcastChannel <- message
}
