package actions

import (
	"fmt"

	"github.com/i5heu/GitCognitio/types"
)

func QrLoginApprove(message types.Message, broadcastChannel *chan types.Message, connections *[]*types.Connection) {
	fmt.Println("QrLoginApprove")
	for _, conn := range *connections {
		if conn.Id == message.Data {
			fmt.Println("Authorized", conn.Id)
			conn.Lock()
			conn.Authorized = true
			conn.Unlock()
			*broadcastChannel <- types.Message{
				ID:   conn.Id,
				Type: "message",
				Data: conn.Conn.RemoteAddr().String() + " is now authorized",
			}

			conn.Conn.WriteJSON(types.Message{
				ID:   conn.Id,
				Type: "message",
				Data: "You are now authorized",
			})
		}
	}
}
