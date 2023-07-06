package types

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

func TestNewConnection(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		upgrader := websocket.Upgrader{}
		conn, err := upgrader.Upgrade(writer, request, nil)
		if err != nil {
			t.Fatal(err)
		}
		// The connection is closed when the handler returns,
		// so we need to read from it in a separate goroutine
		go func() {
			for {
				if _, _, err := conn.ReadMessage(); err != nil {
					return
				}
			}
		}()
	}))
	defer s.Close()

	wsURL := "ws" + s.URL[4:]
	d := websocket.Dialer{}
	conn, _, err := d.Dial(wsURL, nil)
	if err != nil {
		t.Fatal(err)
	}

	c, err := NewConnection(conn)
	assert.NoError(t, err)
	assert.NotNil(t, c)
	assert.False(t, c.IsAuthorized())
}

func TestConnectionAuthorize(t *testing.T) {
	s := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		upgrader := websocket.Upgrader{}
		conn, err := upgrader.Upgrade(writer, request, nil)
		if err != nil {
			t.Fatal(err)
		}

		// Read messages in a loop until connection is closed.
		// This represents your WebSocket server's behavior.
		go func() {
			for {
				msgType, msgData, err := conn.ReadMessage()
				if err != nil {
					return
				}

				// Reply with a message, indicating authorization status.
				// This should mirror your actual server's behavior.
				if msgType == websocket.TextMessage {
					message := string(msgData)
					if message == "this will authorize the connection for all data" {
						conn.WriteMessage(websocket.TextMessage, []byte("authorized"))
					} else {
						conn.WriteMessage(websocket.TextMessage, []byte("unauthorized"))
					}
				}
			}
		}()
	}))
	defer s.Close()

	wsURL := "ws" + s.URL[4:]
	d := websocket.Dialer{}
	conn, _, err := d.Dial(wsURL, nil)
	if err != nil {
		t.Fatal(err)
	}

	c, err := NewConnection(conn)
	assert.NoError(t, err)

	// make bordcast channel
	broadcastChannel := make(chan Message, 256)

	c.Authorize("this will authorize the connection for all data", &broadcastChannel, c)

	assert.True(t, c.IsAuthorized())
}

func TestConnectionSend(t *testing.T) {
	c := Connection{
		Id:         uuid.New(),
		authorized: true,
		Conn:       nil, // this is just a dummy test, so it's fine to be nil
		sendQueue:  make(chan Message, 256),
	}

	msg := Message{
		isPublic: true,
		// fill in the other fields here...
	}

	c.Send(msg)

	sentMsg := <-c.sendQueue

	assert.Equal(t, msg, sentMsg)
}
