package actions

import (
	"fmt"

	"github.com/i5heu/GitCognitio/types"
)

func QrLoginApprove(message types.Message, broadcastChannel *chan types.Message, connections *[]*types.Connection) {
	fmt.Println("QrLoginApprove")
	for _, conn := range *connections {
		if conn.GetId().String() == message.Data {
			fmt.Println("Authorized", conn.Id)

			conn.Authorize("this will authorize the connection for all data")

			*broadcastChannel <- types.Message{
				ID:   conn.GetId().String(),
				Type: "message",
				Data: conn.Conn.RemoteAddr().String() + " is now authorized",
			}

			conn.Send(types.Message{
				ID:   conn.GetId().String(),
				Type: "message",
				Data: "You are now authorized",
			})
		}
	}
}
