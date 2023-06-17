package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/i5heu/GitCognitio/internal/config"
	"github.com/i5heu/GitCognitio/internal/connection"
	"github.com/i5heu/GitCognitio/internal/gitio"
	"github.com/i5heu/GitCognitio/types"
)

func main() {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Failed to get home directory: %s", err)
	}

	path := filepath.Join(home, ".GitCognito")

	rm, err := gitio.NewRepoManager(config.RepoURL, path, filepath.Join(home, ".ssh", "id_rsa"))
	if err != nil {
		log.Fatalf("Failed to initialize RepoManager: %s", err)
	}
	rm.StartPushListener()

	connections := make([]*connection.Connection, 0)
	connectionsMutex := &sync.Mutex{}
	broadcastChannel := make(chan types.Message, 100)
	connection.BroadcastMessageWorker(broadcastChannel, &connections, connectionsMutex)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		connection.HandleConnection(w, r, rm, &connections, connectionsMutex, broadcastChannel)
	})

	fmt.Println("Server started on", config.WebSocketAddress)
	log.Fatal(http.ListenAndServe(config.WebSocketAddress, nil))
}
