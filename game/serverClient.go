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

func SendClientActionToServer() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		log.Println("Shutting down server...")
		os.Exit(0)
	}()

	url := fmt.Sprintf("ws://%s:%s/ws", ACT_SERVER_CONN_HOST, ACT_SERVER_CONN_PORT)
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		log.Fatal("Error connecting to WebSocket server:", err)
	}
	defer conn.Close()
	log.Printf("Client connected to server %s\n", url)

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
			wasPressed := keyStates[k.key]
			if isPressed && !wasPressed {
				msg := fmt.Sprintf("p:%s", k.message)
				if err := conn.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
					log.Println("Error sending key press action:", err)
					return
				}
				keyStates[k.key] = true
			} else if !isPressed && wasPressed {
				// If the key was pressed and is now released, send the release event.
				msg := fmt.Sprintf("r:%s", k.message)
				if err := conn.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
					log.Println("Error sending key release action:", err)
					return
				}
				keyStates[k.key] = false
			}
		}

		// Small sleep to prevent a busy loop.
		time.Sleep(30 * time.Millisecond)
	}
}
