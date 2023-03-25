package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Message struct {
	ID   int         `json:"id"`
	Data interface{} `json:"data"`
}

type Response struct {
	ID   int         `json:"id"`
	Data interface{} `json:"data"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}
		defer conn.Close()

		for {
			messageType, message, err := conn.ReadMessage()
			if err != nil {
				log.Println(err)
				return
			}
			log.Printf("Received message: %s\n", message)

			//unmarshal message
			var msg Message
			err = json.Unmarshal(message, &msg)
			if err != nil {
				log.Println(err)
				return
			}

			//create response
			var resp Response
			resp.ID = msg.ID
			resp.Data = msg.Data.(string) + " from Server"

			// Echo message back to client
			message, err = json.Marshal(resp)
			if err != nil {
				log.Println(err)
				return
			}
			err = conn.WriteMessage(messageType, message)
			if err != nil {
				log.Println(err)
				return
			}
		}
	})

	fmt.Println("Starting server on :8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}
