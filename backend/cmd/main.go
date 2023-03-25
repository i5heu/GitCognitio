package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow any origin for this example
	},
}

type Connection struct {
	*websocket.Conn
	sync.Mutex
}

func (c *Connection) safeWrite(mt int, payload []byte) error {
	c.Lock()
	defer c.Unlock()
	return c.WriteMessage(mt, payload)
}

func main() {
	http.HandleFunc("/", handleConnection)

	fmt.Println("Server started on :8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}

func handleConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("error upgrading connection: %v\n", err)
		return
	}

	c := &Connection{Conn: conn}
	defer c.Close()

	connectionsMutex.Lock()
	connections = append(connections, c)
	connectionsMutex.Unlock()

	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Printf("error reading message: %v\n", err)
			break
		}

		broadcast(mt, message)
	}
}

var connections = make([]*Connection, 0)
var connectionsMutex = &sync.Mutex{}

func broadcast(mt int, message []byte) {
	connectionsMutex.Lock()
	defer connectionsMutex.Unlock()

	for i := 0; i < len(connections); i++ {
		err := connections[i].safeWrite(mt, message)
		if err != nil {
			connections[i].Close()
			connections = append(connections[:i], connections[i+1:]...)
			i--
		}
	}
}
