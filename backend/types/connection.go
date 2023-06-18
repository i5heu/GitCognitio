package types

import (
	"sync"

	"github.com/gorilla/websocket"
)

type Connection struct {
	Id         string
	Authorized bool
	*websocket.Conn
	sync.Mutex
}

func (c *Connection) SafeWrite(mt int, payload []byte) error {
	c.Lock()
	defer c.Unlock()
	return c.WriteMessage(mt, payload)
}
