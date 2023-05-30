package main

import (
	"encoding/json"
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

type Message struct {
	ID   string `json:"id"`
	Type string `json:"type"`
	Data string `json:"data"`
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
		_, byteMessage, err := c.ReadMessage()
		if err != nil {
			log.Printf("error reading message: %v\n", err)
			break
		}

		var message Message
		err = json.Unmarshal(byteMessage, &message)
		if err != nil {
			log.Printf("error unmarshalling message: %v\n", err)
			broadcastMessage(Message{
				ID:   "1",
				Type: "error",
				Data: "error unmarshalling message",
			})
			continue
		}

		broadcastMessage(Message{
			ID:   message.ID,
			Type: message.Type,
			Data: message.Data,
		})
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

func broadcastMessage(message Message) {

	b, err := json.Marshal(message)
	if err != nil {
		log.Printf("error marshalling message: %v\n", err)
		return
	}

	broadcast(websocket.TextMessage, b)
}
