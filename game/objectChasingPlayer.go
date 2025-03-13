package game

import (
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"

	"dynegame/assets"
	"dynegame/utils"
)

const (
	rotationSpeedMin = -0.01
	rotationSpeedMax = 0.04
)

type Meteor struct {
	position      utils.Vector
	rotation      float64
	movement      utils.Vector
	rotationSpeed float64
	sprite        *ebiten.Image
}

func ChasingMovement(m *Meteor, baseVelocity float64, p *Player) *Meteor {
	pos := m.position
	target := p.position

	velocity := baseVelocity + rand.Float64()*1.5

	direction := utils.Vector{
		X: target.X - pos.X,
		Y: target.Y - pos.Y,
	}
	normalizedDirection := direction.Normalize()

	m.movement = utils.Vector{
		X: normalizedDirection.X * velocity,
		Y: normalizedDirection.Y * velocity,
	}
	return m
}
func NewMeteor(baseVelocity float64, p *Player) *Meteor {
	target := p.position

	// the meteor Start from the random Y position on the right edge of the screen
	//    <--- o	|
	//    <--- o	|
	//    <--- o	|
	// posY := rand.Float64() * float64(screenHeight)
	// pos := Vector{
	// 	X: float64(screenWidth), // Right edge of the screen
	// 	Y: posY,
	// }
	posX := rand.Float64() * float64(screenHeight)
	pos := utils.Vector{
		X: posX,
		Y: 0, // upper edge of the screen
	}
	velocity := baseVelocity + rand.Float64()*1.5

	direction := utils.Vector{
		X: target.X - pos.X,
		Y: target.Y - pos.Y,
	}
	normalizedDirection := direction.Normalize()

	movement := utils.Vector{
		X: normalizedDirection.X * velocity,
		Y: normalizedDirection.Y * velocity,
	}

	sprite := assets.MeteorSprites[rand.Intn(len(assets.MeteorSprites))]

	m := &Meteor{
		position:      pos,
		movement:      movement,
		rotationSpeed: rotationSpeedMin + rand.Float64()*(rotationSpeedMax-rotationSpeedMin),
		sprite:        sprite,
	}
	return m
}

func (m *Meteor) Update() {
	m.position.X += m.movement.X
	m.position.Y += m.movement.Y
	m.rotation += m.rotationSpeed
}

func (m *Meteor) Draw(screen *ebiten.Image) {
	bounds := m.sprite.Bounds()
	halfW := float64(bounds.Dx()) / 2
	halfH := float64(bounds.Dy()) / 2

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-halfW, -halfH)
	op.GeoM.Rotate(m.rotation)
	op.GeoM.Translate(halfW, halfH)

	op.GeoM.Translate(m.position.X, m.position.Y)

	screen.DrawImage(m.sprite, op)
}

func (m *Meteor) Collider(BoundsDecreaseRatio float64) utils.Rect {
	bounds := m.sprite.Bounds()

	return utils.NewRect(
		m.position.X,
		m.position.Y,
		float64(bounds.Dx())*BoundsDecreaseRatio,
		float64(bounds.Dy())*BoundsDecreaseRatio,
	)
}
