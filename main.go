package main

import (
	"github.com/hajimehoshi/ebiten/v2"

	"dynegame/game"
)

// todo: create cobra cli to handle cli variable to run client and host
func main() {
	g := game.NewGame()
	// ebiten.SetFullscreen(true)

	err := ebiten.RunGame(g)
	if err != nil {
		panic(err)
	}
}
