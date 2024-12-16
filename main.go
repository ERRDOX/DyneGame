package main

import (
	"github.com/hajimehoshi/ebiten/v2"

	"dynegame/game"
)

func main() {
	g := game.NewGame()
	// ebiten.SetFullscreen(true)

	err := ebiten.RunGame(g)
	if err != nil {
		panic(err)
	}
}
