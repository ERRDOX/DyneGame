package utils

import (
	"math"

	"github.com/hajimehoshi/ebiten/v2"

	"dynegame/assets"
)

const (
	bulletSpeedPerSecond = 900.0
)

type Bullet struct {
	Position Vector
	Rotation float64
	Sprite   *ebiten.Image
}

func NewBullet(pos Vector, rotation float64) *Bullet {
	sprite := assets.PoleSprite

	bounds := sprite.Bounds()
	halfW := float64(bounds.Dx()) / 2
	halfH := float64(bounds.Dy()) / 2

	pos.X -= halfW
	pos.Y -= halfH

	b := &Bullet{
		Position: pos,
		Rotation: rotation,
		Sprite:   sprite,
	}

	return b
}

func (b *Bullet) Update() {
	speed := bulletSpeedPerSecond / float64(ebiten.TPS())

	b.Position.X += math.Sin(b.Rotation) * speed
	b.Position.Y += math.Cos(b.Rotation) * -speed
}

func (b *Bullet) Draw(screen *ebiten.Image) {
	bounds := b.Sprite.Bounds()
	halfW := float64(bounds.Dx()) / 2
	halfH := float64(bounds.Dy()) / 2

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-halfW, -halfH)
	op.GeoM.Rotate(b.Rotation)
	op.GeoM.Translate(halfW, halfH)

	op.GeoM.Translate(b.Position.X, b.Position.Y)

	screen.DrawImage(b.Sprite, op)
}

func (b *Bullet) Collider(BoundsDecreaseRatio float64) Rect {
	bounds := b.Sprite.Bounds()

	return NewRect(
		b.Position.X,
		b.Position.Y,
		float64(bounds.Dx())*BoundsDecreaseRatio,
		float64(bounds.Dy())*BoundsDecreaseRatio,
	)
}
