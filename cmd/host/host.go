package main

import (
	"github.com/hajimehoshi/ebiten/v2"

	"dynegame/game"
)

func main() {
	g := game.NewGameHost()

	err := ebiten.RunGame(g)
	if err != nil {
		panic(err)
	}
}
