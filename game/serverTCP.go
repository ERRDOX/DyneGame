//go:build ignore

package game

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

const (
	CONN_HOST = "0.0.0.0"
	CONN_PORT = "8089"
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

func (a *Action) Joiner() {
	listener, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
	if err != nil {
		log.Fatalf("Error listening: %v\n", err)
	}
	defer listener.Close()
	log.Printf("Listening on %s:%s\n", CONN_HOST, CONN_PORT)

	// Graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		log.Println("Shutting down server...")
		listener.Close()
		os.Exit(0)
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v\n", err)
			continue
		}
		go handleRequest(conn, a)
	}
}

func handleRequest(conn net.Conn, a *Action) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	for {
		// Read until a newline or EOF
		data, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("Error reading: %v\n", err)
			break
		}
		data = data[:len(data)-1] // Remove newline character

		fmt.Printf("Received: %s\n", data)
		a.SetAct(data) // Update action with new data
		message := fmt.Sprintf("Message received: %s\n", data)
		conn.Write([]byte(message))
	}
}
