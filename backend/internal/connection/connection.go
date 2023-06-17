package connection

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/i5heu/GitCognitio/internal/actions"
	"github.com/i5heu/GitCognitio/internal/config"
	"github.com/i5heu/GitCognitio/internal/gitio"
	"github.com/i5heu/GitCognitio/types"
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

var Upgrader = websocket.Upgrader{
	ReadBufferSize:  config.ReadBufferSize,
	WriteBufferSize: config.WriteBufferSize,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// HandleConnection handles the connection received from the HTTP server.
func HandleConnection(w http.ResponseWriter, r *http.Request, rm *gitio.RepoManager, connections *[]*Connection, connectionsMutex *sync.Mutex, broadcastChannel chan types.Message) {
	conn, err := UpgradeConnection(w, r)
	if err != nil {
		log.Printf("error upgrading connection: %v\n", err)
		return
	}
	defer func() {
		if err := conn.Close(); err != nil {
			log.Printf("error closing connection: %v\n", err)
		}
	}()

	// AddConnection safely adds a new connection to the existing connection pool.
	AddConnection(connections, connectionsMutex, conn)

	for {
		_, byteMessage, err := conn.ReadMessage()
		if err != nil {
			log.Printf("error reading message: %v\n", err)
			break
		}

		var message types.Message
		err = json.Unmarshal(byteMessage, &message)
		if err != nil {
			log.Printf("error unmarshalling message: %v\n", err)
			conn.WriteMessage(websocket.TextMessage, []byte("error unmarshalling message"))
			continue
		}

		if conn.Authorized {
			HandleMessage(message, broadcastChannel, rm)
		} else {
			err = AuthenticateMessage(message, &conn.Authorized)
			if err != nil {
				conn.WriteMessage(websocket.TextMessage, []byte(err.Error()))
				continue
			}
		}
	}
}

// UpgradeConnection upgrades the HTTP server connection to the WebSocket protocol.
func UpgradeConnection(w http.ResponseWriter, r *http.Request) (*Connection, error) {
	conn, err := Upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}
	return &Connection{Conn: conn}, nil
}

// AddConnection adds a new connection to the connection pool.
func AddConnection(connections *[]*Connection, connectionsMutex *sync.Mutex, conn *Connection) {
	uuid := uuid.New()

	// prepare connection
	uuidString := uuid.String()
	conn.Id = uuidString
	conn.Authorized = false

	connectionsMutex.Lock()
	*connections = append(*connections, conn)
	connectionsMutex.Unlock()
}

// AuthenticateMessage checks the provided message for correct authentication data.
func AuthenticateMessage(message types.Message, authorized *bool) error {
	password := config.PassWord

	if message.Type == "message" && message.Data == "!pwd "+password {
		*authorized = true
		return nil
	}
	return errors.New("auth error")
}

// HandleMessage performs action based on the message type.
func HandleMessage(message types.Message, broadcastChannel chan types.Message, rm *gitio.RepoManager) {
	switch message.Type {
	case "message":
		actions.NewMdFile(message, &broadcastChannel, rm)
	case "delete":
		actions.DeleteFile(message, &broadcastChannel, rm)
	case "typing":
		BroadcastMessage(broadcastChannel, types.Message{
			ID:   message.ID,
			Type: message.Type,
			Data: message.Data,
		})
	default:
		BroadcastMessage(broadcastChannel, types.Message{
			ID:   message.ID,
			Type: "error",
			Data: "unknown message type",
		})
	}
}

func BroadcastMessage(broadcastChannel chan types.Message, message types.Message) {
	broadcastChannel <- message
}

func BroadcastMessageWorker(broadcastChannel <-chan types.Message, connections *[]*Connection, connectionsMutex *sync.Mutex) {
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
				err := (*connections)[i].SafeWrite(websocket.TextMessage, b)
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
