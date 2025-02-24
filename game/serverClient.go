package game

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/websocket"
	"github.com/hajimehoshi/ebiten/v2"
)

func sendClientActionToServer() {
	// Graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		log.Println("Shutting down server...")
		os.Exit(0)
	}()

	// Establish WebSocket connection
	url := fmt.Sprintf("ws://%s:%s/ws", ACT_SERVER_CONN_HOST, ACT_SERVER_CONN_PORT)
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal("Error connecting to WebSocket server:", err)
	}
	defer conn.Close()
	log.Printf("Client connected to server %s\n", url)

	// Send initial "map" request and log the response
	if err := conn.WriteMessage(websocket.TextMessage, []byte("map")); err != nil {
		log.Println("Error sending map request:", err)
		return
	}
	_, message, err := conn.ReadMessage()
	if err != nil {
		log.Println("Error reading map response:", err)
		return
	}
	log.Printf("Map is %s\n", string(message))

	// Define the keys to track and their corresponding messages.
	keysToSend := []struct {
		key     ebiten.Key
		message string
	}{
		{ebiten.KeyW, "w"},
		{ebiten.KeyA, "a"},
		{ebiten.KeyS, "s"},
		{ebiten.KeyD, "d"},
		{ebiten.KeyUp, "up"},
		{ebiten.KeyDown, "down"},
		{ebiten.KeyLeft, "left"},
		{ebiten.KeyRight, "right"},
		{ebiten.KeySpace, "space"},
	}

	// Map to track the current state of each key (true if pressed)
	keyStates := make(map[ebiten.Key]bool)

	// Main loop: only send events when a key is pressed or released.
	for {
		for _, k := range keysToSend {
			isPressed := ebiten.IsKeyPressed(k.key)
			// wasPressed := keyStates[k.key] // defaults to false if not present

			// If the key has just been pressed, send the press event.
			if isPressed {
				msg := fmt.Sprintf("p: ", k.message)
				if err := conn.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
					log.Println("Error sending key press action:", err)
					return
				}
				keyStates[k.key] = true
			} else if !isPressed {
				// If the key was pressed and is now released, send the release event.
				msg := fmt.Sprintf("r:", k.message)
				if err := conn.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
					log.Println("Error sending key release action:", err)
					return
				}
				keyStates[k.key] = false
			}
		}

		// Small sleep to prevent a busy loop.
		time.Sleep(10 * time.Millisecond)
	}
}
