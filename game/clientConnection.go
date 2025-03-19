package game

import (
	"dynegame/utils"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/gorilla/websocket"
)

// ClientConnection handles all WebSocket connections for the client
type ClientConnection struct {
	game *Game
	// Action connection for sending player actions
	actionConn *websocket.Conn
	// Position connection for receiving host position
	hostPosConn *websocket.Conn
	// Bullet connection for receiving host bullets
	hostBulletConn *websocket.Conn
	// Position connection for receiving client position
	clientPosConn *websocket.Conn
	// Bullet connection for receiving client bullets
	clientBulletConn *websocket.Conn
}

func NewClientConnection(game *Game) *ClientConnection {
	return &ClientConnection{
		game: game,
	}
}

// Connect establishes all necessary WebSocket connections
func (c *ClientConnection) Connect() error {
	// Connect to action server
	actionURL := fmt.Sprintf("ws://%s:%s/ws/action", ACT_SERVER_CONN_HOST, ACT_SERVER_CONN_PORT)
	actionConn, _, err := websocket.DefaultDialer.Dial(actionURL, nil)
	if err != nil {
		return fmt.Errorf("failed to connect to action server: %v", err)
	}
	c.actionConn = actionConn

	// Connect to host position server
	hostPosURL := fmt.Sprintf("ws://%s:%s/ws/position/host", STATUS_SERVER_CONN_HOST, STATUS_SERVER_CONN_PORT)
	hostPosConn, _, err := websocket.DefaultDialer.Dial(hostPosURL, nil)
	if err != nil {
		return fmt.Errorf("failed to connect to host position server: %v", err)
	}
	c.hostPosConn = hostPosConn

	// Connect to host bullet server
	hostBulletURL := fmt.Sprintf("ws://%s:%s/ws/position/bullet/host", STATUS_SERVER_CONN_HOST, STATUS_SERVER_CONN_PORT)
	hostBulletConn, _, err := websocket.DefaultDialer.Dial(hostBulletURL, nil)
	if err != nil {
		return fmt.Errorf("failed to connect to host bullet server: %v", err)
	}
	c.hostBulletConn = hostBulletConn

	// Connect to client position server
	clientPosURL := fmt.Sprintf("ws://%s:%s/ws/position/client", STATUS_SERVER_CONN_HOST, STATUS_SERVER_CONN_PORT)
	clientPosConn, _, err := websocket.DefaultDialer.Dial(clientPosURL, nil)
	if err != nil {
		return fmt.Errorf("failed to connect to client position server: %v", err)
	}
	c.clientPosConn = clientPosConn

	// Connect to client bullet server
	clientBulletURL := fmt.Sprintf("ws://%s:%s/ws/position/bullet/client", STATUS_SERVER_CONN_HOST, STATUS_SERVER_CONN_PORT)
	clientBulletConn, _, err := websocket.DefaultDialer.Dial(clientBulletURL, nil)
	if err != nil {
		return fmt.Errorf("failed to connect to client bullet server: %v", err)
	}
	c.clientBulletConn = clientBulletConn

	// Start receiving data in goroutines
	go c.receiveHostPosition()
	go c.receiveHostBullets()
	go c.receiveClientPosition()
	go c.receiveClientBullets()

	return nil
}

// receiveHostPosition continuously receives host player position updates
func (c *ClientConnection) receiveHostPosition() {
	for {
		_, message, err := c.hostPosConn.ReadMessage()
		if err != nil {
			log.Printf("Error reading host position: %v", err)
			return
		}

		// Parse position data (format: "x,y,rotation")
		parts := strings.Split(string(message), ",")
		if len(parts) != 3 {
			log.Printf("Invalid position data format: %s", message)
			continue
		}

		x, err := strconv.ParseFloat(parts[0], 64)
		if err != nil {
			log.Printf("Error parsing X position: %v", err)
			continue
		}

		y, err := strconv.ParseFloat(parts[1], 64)
		if err != nil {
			log.Printf("Error parsing Y position: %v", err)
			continue
		}

		rotation, err := strconv.ParseFloat(parts[2], 64)
		if err != nil {
			log.Printf("Error parsing rotation: %v", err)
			continue
		}

		// Update host player position
		c.game.player.position.X = x
		c.game.player.position.Y = y
		c.game.player.rotation = rotation
	}
}

// receiveHostBullets continuously receives host player bullet updates
func (c *ClientConnection) receiveHostBullets() {
	for {
		_, message, err := c.hostBulletConn.ReadMessage()
		if err != nil {
			log.Printf("Error reading host bullets: %v", err)
			return
		}

		// Parse bullet data (format: "x,y,rotation")
		parts := strings.Split(string(message), ",")
		if len(parts) != 3 {
			log.Printf("Invalid bullet data format: %s", message)
			continue
		}

		x, err := strconv.ParseFloat(parts[0], 64)
		if err != nil {
			log.Printf("Error parsing bullet X position: %v", err)
			continue
		}

		y, err := strconv.ParseFloat(parts[1], 64)
		if err != nil {
			log.Printf("Error parsing bullet Y position: %v", err)
			continue
		}

		rotation, err := strconv.ParseFloat(parts[2], 64)
		if err != nil {
			log.Printf("Error parsing bullet rotation: %v", err)
			continue
		}

		// Create new bullet at received position
		bullet := utils.NewBullet(utils.Vector{X: x, Y: y}, rotation)
		c.game.AddBulletPlayer(bullet)
	}
}

// receiveClientPosition continuously receives client player position updates
func (c *ClientConnection) receiveClientPosition() {
	for {
		_, message, err := c.clientPosConn.ReadMessage()
		if err != nil {
			log.Printf("Error reading client position: %v", err)
			return
		}

		// Parse position data (format: "x,y,rotation")
		parts := strings.Split(string(message), ",")
		if len(parts) != 3 {
			log.Printf("Invalid position data format: %s", message)
			continue
		}

		x, err := strconv.ParseFloat(parts[0], 64)
		if err != nil {
			log.Printf("Error parsing X position: %v", err)
			continue
		}

		y, err := strconv.ParseFloat(parts[1], 64)
		if err != nil {
			log.Printf("Error parsing Y position: %v", err)
			continue
		}

		rotation, err := strconv.ParseFloat(parts[2], 64)
		if err != nil {
			log.Printf("Error parsing rotation: %v", err)
			continue
		}

		// Update client player position
		c.game.SecondPlayer.position.X = x
		c.game.SecondPlayer.position.Y = y
		c.game.SecondPlayer.rotation = rotation
	}
}

// receiveClientBullets continuously receives client player bullet updates
func (c *ClientConnection) receiveClientBullets() {
	for {
		_, message, err := c.clientBulletConn.ReadMessage()
		if err != nil {
			log.Printf("Error reading client bullets: %v", err)
			return
		}

		// Parse bullet data (format: "x,y,rotation")
		parts := strings.Split(string(message), ",")
		if len(parts) != 3 {
			log.Printf("Invalid bullet data format: %s", message)
			continue
		}

		x, err := strconv.ParseFloat(parts[0], 64)
		if err != nil {
			log.Printf("Error parsing bullet X position: %v", err)
			continue
		}

		y, err := strconv.ParseFloat(parts[1], 64)
		if err != nil {
			log.Printf("Error parsing bullet Y position: %v", err)
			continue
		}

		rotation, err := strconv.ParseFloat(parts[2], 64)
		if err != nil {
			log.Printf("Error parsing bullet rotation: %v", err)
			continue
		}

		// Create new bullet at received position
		bullet := utils.NewBullet(utils.Vector{X: x, Y: y}, rotation)
		c.game.AddBulletSecondPlayer(bullet)
	}
}

// Close closes all WebSocket connections
func (c *ClientConnection) Close() {
	if c.actionConn != nil {
		c.actionConn.Close()
	}
	if c.hostPosConn != nil {
		c.hostPosConn.Close()
	}
	if c.hostBulletConn != nil {
		c.hostBulletConn.Close()
	}
	if c.clientPosConn != nil {
		c.clientPosConn.Close()
	}
	if c.clientBulletConn != nil {
		c.clientBulletConn.Close()
	}
}
