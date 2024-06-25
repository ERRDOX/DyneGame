package game

import (
	"game/assets"

	"github.com/hajimehoshi/ebiten/v2"
)

func Background(screen *ebiten.Image) {
	screen.DrawImage(assets.BackGround, nil)

}
