//go:build ignore
// +build ignore

// go:build exclude
package game

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/gorilla/websocket"
)

const (
	CONN_HOST = "localhost"
	CONN_PORT = "8080"
	CONN_TYPE = "tcp"
)

type Action struct {
	mu  sync.Mutex
	Act string
}

func NewAction() *Action {
	return &Action{Act: ""}
}

func (a *Action) SetAct(act string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.Act = act
}

func (a *Action) GetAct() string {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.Act
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func (a *Action) Joiner() {

	http.HandleFunc("/ws", a.handleConnections)
	log.Println("Starting server on %s:%s\n", CONN_HOST, CONN_PORT)
	log.Fatal(http.ListenAndServe(":"+CONN_PORT, nil))
	log.Printf("Listening on %s:%s\n", CONN_HOST, CONN_PORT)

	// Graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		log.Println("Shutting down server...")
		os.Exit(0)
	}()

	// for {
	// 	conn, err := listener.Accept()
	// 	if err != nil {
	// 		log.Printf("Error accepting connection: %v\n", err)
	// 		continue
	// 	}
	// 	go handleRequest(conn, a)
	// }
}

func (a *Action) handleConnections(w http.ResponseWriter, r *http.Request) {
	// var upgrader = websocket.Upgrader{
	// 	ReadBufferSize:  32,
	// 	WriteBufferSize: 32,
	// 	CheckOrigin:     func(r *http.Request) bool { return true },
	// }
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()
	for {
		// Read until a newline or EOF
		_, data, err := ws.ReadMessage()
		if err != nil {
			log.Printf("Error reading: %v\n", err)
			break
		}
		// data = data[:len(data)-1] // Remove newline character

		fmt.Printf("Received: %s\n", string(data[:]))
		println(data[:])
		a.SetAct(string(data[:])) // Update action with new data
		if err := ws.WriteMessage(websocket.TextMessage, data); err != nil {
			log.Println("write:", err)
			break
		}
	}
}
