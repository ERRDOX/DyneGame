package utils

import (
	"time"

	"github.com/hajimehoshi/ebiten/v2"

	"dynegame/assets"
)

type Explosion struct {
	position Vector
	sprite   []*ebiten.Image
}

func NewExplosion(pos Vector) *Explosion {
	sprite := assets.Explosion

	e := &Explosion{
		position: pos,
		sprite:   sprite,
	}
	return e
}

func (e *Explosion) Update() {
	// do nothing
}

func (e *Explosion) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}

	op.GeoM.Translate(e.position.X, e.position.Y)
	// for i := range e.sprite {
	screen.DrawImage(e.sprite[3], op)
	time.Sleep(70 * time.Millisecond)
	// println(i)
	// }
}
