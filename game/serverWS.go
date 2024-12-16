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

type Action struct {
	mu  sync.Mutex
	Act map[string]bool
}

func NewAction() *Action {
	return &Action{Act: make(map[string]bool)}
}

func (a *Action) SetAct(act string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.Act[act] = true
}
func (a *Action) RemoveAct(act string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.Act[act] = false
}

func (a *Action) GetAct() map[string]bool {
	a.mu.Lock()
	defer a.mu.Unlock()
	return a.Act
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func (g *Game) resposeServer() {
	http.HandleFunc("/ws/loc", g.respHandleConnections)
	log.Fatalln(http.ListenAndServe(":"+STATUS_SERVER_CONN_PORT, nil))
	fmt.Printf("Listening on %s:%s\n", STATUS_SERVER_CONN_HOST, STATUS_SERVER_CONN_PORT)

}
func (g *Game) respHandleConnections(w http.ResponseWriter, r *http.Request) {
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
		if string(data[:]) == "pos" {
			fmt.Printf("%f,%f,%f", g.player.position.X, g.player.position.Y, g.player.rotation)
			ws.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("%f,%f,%f", g.player.position.X, g.player.position.Y, g.player.rotation)))
			ws.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("%f,%f,%f", g.SecondPlayer.position.X, g.SecondPlayer.position.Y, g.SecondPlayer.rotation)))
			fmt.Println("Received: " + string(data[:]))
			continue
		}
	}
}

// Server is a function that starts the server intializing the websocket for the start of the game
func (a *Action) Server() {
	fmt.Printf("Listening on %s:%s\n", ACT_SERVER_CONN_HOST, ACT_SERVER_CONN_PORT)
	http.HandleFunc("/ws", a.serverHandleConnections)
	log.Fatalln(http.ListenAndServe(":"+ACT_SERVER_CONN_PORT, nil))
	fmt.Printf("Listening on %s:%s\n", ACT_SERVER_CONN_HOST, ACT_SERVER_CONN_PORT)

	// Graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		log.Println("Shutting down server...")
		os.Exit(0)
	}()

}

// the rival function is the same as the server function
// but it is used to connect to the server for the second player
func (a *Action) Rival() {

	http.HandleFunc("/ws/rival", a.rivalHandleConnections)
	log.Printf("Starting rival on %s:%s\n", ACT_SERVER_CONN_HOST, ACT_SERVER_CONN_PORT)
	log.Fatalln(http.ListenAndServe(":"+ACT_SERVER_CONN_PORT, nil))

	// Graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		log.Println("Shutting down server...")
		os.Exit(0)
	}()

}
func (a *Action) rivalHandleConnections(w http.ResponseWriter, r *http.Request) {
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
		// debug the map
		if string(data[:]) == "map" {
			ws.WriteMessage(websocket.TextMessage, []byte("DragonMap"))
			break
		}
		// data = data[:len(data)-1] // Remove newline character

		fmt.Println("Received: " + string(data[:]))
		if data[0] == 112 {
			a.SetAct(string(data[1:]))
		} else {
			a.RemoveAct(string(data[1:]))
		}
		if err := ws.WriteMessage(websocket.TextMessage, data); err != nil {
			log.Println("write:", err)
			break
		}
	}
}

func (a *Action) serverHandleConnections(w http.ResponseWriter, r *http.Request) {
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

		if string(data[:]) == "map" {
			ws.WriteMessage(websocket.TextMessage, []byte(MAP))
			fmt.Println("Received: " + string(data[:]))
			continue
		}
		// data = data[:len(data)-1] // Remove newline character

		fmt.Println("Received: " + string(data[:]))
		if data[0] == 112 {
			a.SetAct(string(data[1:]))
		} else {
			a.RemoveAct(string(data[1:]))
		}
		if err := ws.WriteMessage(websocket.TextMessage, data); err != nil {
			log.Println("write:", err)
			break
		}
	}
}

// call server to get the map
func (a *Action) GetMap() string {
	conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:8080/ws", nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer conn.Close()
	conn.WriteMessage(websocket.TextMessage, []byte("map"))
	_, message, err := conn.ReadMessage()
	if err != nil {
		log.Println("read:", err)
		return ""
	}
	return string(message)
}

//call server to get the server player act and location of both players

// func (a *Action) GetActLoc() map[string]bool {
// 	conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:8080/ws", nil)
// 	if err != nil {
// 		log.Fatal("dial:", err)
// 	}
// 	defer conn.Close()
// 	conn.WriteMessage(websocket.TextMessage, []byte("pos"))
// 	_, message, err := conn.ReadMessage()
// 	if err != nil {
// 		log.Println("read:", err)
// 		return nil
// 	}
// 	return a.GetAct()
// }
