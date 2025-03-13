package game

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

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

// ServerInit is a function that starts the server intializing the websocket
func (a *Action) ServerGetClientAction() {
	fmt.Printf("Listening on %s:%s\n", ACT_SERVER_CONN_HOST, ACT_SERVER_CONN_PORT)
	http.HandleFunc("/ws", a.serverHandlClientAction)
	log.Fatalln(http.ListenAndServe(":"+ACT_SERVER_CONN_PORT, nil))
	// fmt.Printf("Listening on %s:%s\n", ACT_SERVER_CONN_HOST, ACT_SERVER_CONN_PORT)
	// Graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		log.Println("Shutting down server...")
		os.Exit(0)
	}()
}

// this handler is used to send the player action to the client
func (g *Game) ServeBulletPosAndPlayerPos() {
	fmt.Printf("Serving bullets Listening on %s:%s\n", STATUS_SERVER_CONN_HOST, STATUS_SERVER_CONN_PORT)
	go func() {
		http.HandleFunc("/ws/position/bullet/client", g.respClientBulletHandler)
	}()
	go func() {
		http.HandleFunc("/ws/position/bullet/host", g.respHostBulletHandler)
	}()
	go func() {
		http.HandleFunc("/ws/position/host", g.respHostPosHandler)
	}()
	go func() {
		http.HandleFunc("/ws/position/client", g.respClientPosHandler)
	}()
	log.Fatalln(http.ListenAndServe(":"+STATUS_SERVER_CONN_PORT, nil))
}

func (g *Game) respHostPosHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()
	for {
		time.Sleep(120 * time.Millisecond)
		fmt.Printf("%f,%f,%f", g.player.position.X, g.player.position.Y, g.player.rotation)
		ws.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("%f,%f,%f", g.player.position.X, g.player.position.Y, g.player.rotation)))
		fmt.Println("Received: ")

	}
}
func (g *Game) respClientPosHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()
	for {
		time.Sleep(120 * time.Millisecond)
		fmt.Printf("%f,%f,%f", g.SecondPlayer.position.X, g.SecondPlayer.position.Y, g.SecondPlayer.rotation)
		ws.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("%f,%f,%f", g.SecondPlayer.position.X, g.SecondPlayer.position.Y, g.SecondPlayer.rotation)))
		fmt.Println("Received: ")

	}
}

func (g *Game) respClientBulletHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()
	for {
		time.Sleep(120 * time.Millisecond)
		if len(g.SecondPlayer.bullet) == 0 {
			continue
		}
		for _, b := range g.player.bullet {
			ws.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("%f,%f,%f", b.Position.X, b.Position.Y, b.Rotation)))
			fmt.Printf("Sent: %f, %f, %f\n", b.Position.X, b.Position.Y, b.Rotation)
		}
		if err != nil {
			log.Printf("Error marshalling bullet data: %v\n", err)
			break
		}
	}
}

func (g *Game) respHostBulletHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()
	for {
		time.Sleep(120 * time.Millisecond)
		if len(g.player.bullet) == 0 {
			continue
		}
		for _, b := range g.player.bullet {
			ws.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("%f,%f,%f", b.Position.X, b.Position.Y, b.Rotation)))
			fmt.Printf("Sent: %f, %f, %f\n", b.Position.X, b.Position.Y, b.Rotation)
		}
		if err != nil {
			log.Printf("Error marshalling bullet data: %v\n", err)
			break
		}
	}
}

// the ServerHostAct function is the same as the server function
// but it is used to connect to the server for the client player
func (a *Action) ServerHostAct() {
	http.HandleFunc("/ws/ServerHostAct", a.rivalHandleConnections)
	log.Printf("Starting Host on %s:%s\n", ACT_SERVER_CONN_HOST, ACT_SERVER_CONN_PORT)
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

func (a *Action) serverHandlClientAction(w http.ResponseWriter, r *http.Request) {
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
	conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:8080/ws/map", nil)
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
