package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/i5heu/GitCognitio/internal/actions"
	"github.com/i5heu/GitCognitio/internal/gitio"
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

var upgrader = websocket.Upgrader{
	ReadBufferSize:  ReadBufferSize,
	WriteBufferSize: WriteBufferSize,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow any origin for this example
	},
}

func handleConnection(w http.ResponseWriter, r *http.Request, rm *gitio.RepoManager, connections *[]*Connection, connectionsMutex *sync.Mutex) {
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
	*connections = append(*connections, c)
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
			}, connections, connectionsMutex)
			continue
		}

		stat, err := rm.GetRepoStats()
		if err != nil {
			log.Printf("error getting repo stats: %v\n", err)
			continue
		}

		if message.Type == "message" {
			err = actions.NewMdFile(message.Data, "test.md", rm)
			if err != nil {
				log.Printf("error creating file: %v\n", err)
				broadcastMessage(Message{
					ID:   "1",
					Type: "error",
					Data: "error creating file",
				}, connections, connectionsMutex)
				continue
			}
		}

		broadcastMessage(Message{
			ID:   message.ID,
			Type: message.Type,
			Data: strconv.FormatInt(stat.RepoSize, 10),
		}, connections, connectionsMutex)
	}
}

func broadcastMessage(message Message, connections *[]*Connection, connectionsMutex *sync.Mutex) {

	b, err := json.Marshal(message)
	if err != nil {
		log.Printf("error marshalling message: %v\n", err)
		return
	}

	connectionsMutex.Lock()
	defer connectionsMutex.Unlock()

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

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		handleConnection(w, r, rm, &connections, connectionsMutex)
	})

	fmt.Println("Server started on", WebSocketAddress)
	log.Fatal(http.ListenAndServe(WebSocketAddress, nil))
}
