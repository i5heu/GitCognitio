package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/i5heu/GitCognitio/internal/actions"
	"github.com/i5heu/GitCognitio/internal/gitio"
	"github.com/i5heu/GitCognitio/types"
)

const (
	repoURL          = "git@github.com:i5heu/Tyche-Test.git"
	ReadBufferSize   = 1024
	WriteBufferSize  = 1024
	WebSocketAddress = ":8081"
)

type Connection struct {
	*websocket.Conn
	sync.Mutex
}

func (c *Connection) safeWrite(mt int, payload []byte) error {
	c.Lock()
	defer c.Unlock()
	return c.WriteMessage(mt, payload)
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  ReadBufferSize,
	WriteBufferSize: WriteBufferSize,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow any origin for this example
	},
}

func handleConnection(w http.ResponseWriter, r *http.Request, rm *gitio.RepoManager, connections *[]*Connection, connectionsMutex *sync.Mutex, broadcastChannel chan types.Message) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("error upgrading connection: %v\n", err)
		return
	}

	c := &Connection{Conn: conn}
	defer func() {
		if err := c.Close(); err != nil {
			log.Printf("error closing connection: %v\n", err)
		}
	}()

	connectionsMutex.Lock()
	//needs authentication
	*connections = append(*connections, c)
	connectionsMutex.Unlock()

	authorized := false
	password := "pass123"

	for {
		_, byteMessage, err := c.ReadMessage()
		if err != nil {
			log.Printf("error reading message: %v\n", err)
			break
		}

		var message types.Message
		err = json.Unmarshal(byteMessage, &message)
		if err != nil {
			log.Printf("error unmarshalling message: %v\n", err)
			c.WriteMessage(websocket.TextMessage, []byte("error unmarshalling message"))
			continue
		}

		if !authorized {
			if message.Type == "message" && message.Data == "!pwd "+password {
				authorized = true
				c.WriteMessage(websocket.TextMessage, []byte("auth ok"))
			} else {
				c.WriteMessage(websocket.TextMessage, []byte("auth error"))
			}
			continue
		}

		switch message.Type {
		case "message":
			actions.NewMdFile(message, &broadcastChannel, rm)
		case "delete":
			actions.DeleteFile(message, &broadcastChannel, rm)
		case "typing":
			broadcastMessage(broadcastChannel, types.Message{
				ID:   message.ID,
				Type: message.Type,
				Data: message.Data,
			})
		default:
			broadcastMessage(broadcastChannel, types.Message{
				ID:   message.ID,
				Type: "error",
				Data: "unknown message type",
			})
			continue
		}

	}
}

func broadcastMessage(broadcastChannel chan types.Message, message types.Message) {
	broadcastChannel <- message
}

func broadcastMessageWorker(broadcastChannel <-chan types.Message, connections *[]*Connection, connectionsMutex *sync.Mutex) {
	go func() {
		for message := range broadcastChannel {
			fmt.Println("broadcasting message", message)

			b, err := json.Marshal(message)
			if err != nil {
				log.Printf("error marshalling message: %v\n", err)
				return
			}
			connectionsMutex.Lock()

			for i := 0; i < len(*connections); i++ {
				err := (*connections)[i].safeWrite(websocket.TextMessage, b)
				if err != nil {
					if err := (*connections)[i].Close(); err != nil {
						log.Printf("error closing connection: %v\n", err)
					}
					*connections = append((*connections)[:i], (*connections)[i+1:]...)
					i--
				}
			}

			connectionsMutex.Unlock()
		}
	}()
}

func main() {
	// Get the home directory
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Failed to get home directory: %s", err)
	}

	// Set the path to the .GitCognito directory in the home directory
	path := filepath.Join(home, ".GitCognito")

	// Initialize RepoManager
	rm, err := gitio.NewRepoManager(repoURL, path, filepath.Join(home, ".ssh", "id_rsa"))
	if err != nil {
		log.Fatalf("Failed to initialize RepoManager: %s", err)
	}
	rm.StartPushListener()

	connections := make([]*Connection, 0)
	connectionsMutex := &sync.Mutex{}
	broadcastChannel := make(chan types.Message, 100)
	broadcastMessageWorker(broadcastChannel, &connections, connectionsMutex)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handleConnection(w, r, rm, &connections, connectionsMutex, broadcastChannel)
	})

	fmt.Println("Server started on", WebSocketAddress)
	log.Fatal(http.ListenAndServe(WebSocketAddress, nil))
}
