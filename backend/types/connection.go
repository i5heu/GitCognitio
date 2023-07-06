package types

import (
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Connection struct {
	*websocket.Conn
	Id                  uuid.UUID `json:"id"`
	sendQueue           chan Message
	authorized          bool
	authorizedMu        sync.Mutex
	sendWorkerRunning   bool
	sendWorkerRunningMu sync.Mutex
}

func NewConnection(conn *websocket.Conn) (*Connection, error) {
	thisCon := &Connection{
		Id:         uuid.New(),
		authorized: false,
		Conn:       conn,
		sendQueue:  make(chan Message, 256),
	}

	err := thisCon.createSendWorker()
	if err != nil {
		fmt.Println("error creating send worker", err)
	}

	return thisCon, err
}

func (c *Connection) Send(message Message) {
	if !message.isPublic && !c.authorized {
		return
	}

	c.sendQueue <- message
}

func (c *Connection) IsAuthorized() bool {
	c.authorizedMu.Lock()
	defer c.authorizedMu.Unlock()

	return c.authorized
}

func (c *Connection) Authorize(risk string, broadcastChannel *chan Message, conn *Connection) {
	if "this will authorize the connection for all data" == risk {
		c.authorizedMu.Lock()
		defer c.authorizedMu.Unlock()
		c.authorized = true

		*broadcastChannel <- Message{
			ID:   conn.GetId().String(),
			Type: "message",
			Data: conn.Conn.RemoteAddr().String() + " is now authorized",
		}

		conn.Send(Message{
			ID:   conn.GetId().String(),
			Type: "message",
			Data: "You are now authorized",
		})

	} else {
		panic("the developer was not aware of the risk of this function")
	}
}

func (c *Connection) GetId() uuid.UUID {
	return c.Id
}

func (c *Connection) createSendWorker() error {
	c.sendWorkerRunningMu.Lock()

	if c.sendWorkerRunning {
		c.sendWorkerRunningMu.Unlock()
		return errors.New("send worker is already running")
	}

	c.sendWorkerRunning = true
	c.sendWorkerRunningMu.Unlock()

	go func() {
		defer selfDestruct(c)

		for message := range c.sendQueue {
			b, err := json.Marshal(message)
			if err != nil {
				fmt.Println("error marshalling message", err)
				continue
			}

			err = c.WriteMessage(websocket.TextMessage, b)
			if err != nil {
				fmt.Println("error writing message", err)
				continue
			}
		}
	}()

	return nil
}

func selfDestruct(c *Connection) {
	c.sendWorkerRunningMu.Lock()
	defer c.sendWorkerRunningMu.Unlock()

	if c.sendWorkerRunning {
		close(c.sendQueue)
		c.sendWorkerRunning = false
	} else {
		fmt.Println("send worker was not set as running although it should be")
	}
}
