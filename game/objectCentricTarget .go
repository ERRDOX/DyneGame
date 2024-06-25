//go:build ignore
// +build ignore

/*
Target Position:
The target is set to the center of the screen. This is calculated by dividing the screen width and height by 2. The meteor will eventually move towards this point.

Random Angle Calculation:
An angle is randomly generated, which is used to place the meteor at a random position around the screen. The angle is in radians (as rand.Float64() * 2 * math.Pi gives a value between 0 and 2π).

Initial Position of Meteor:
The meteor's initial position (pos) is calculated using the random angle and a radius (r) set to half of the screen width. This places the meteor on an imaginary circle around the center of the screen.

Velocity Calculation:
The meteor's velocity is calculated by adding a base velocity (baseVelocity) to a random value multiplied by 1.5. This adds variability to how fast different meteors move.

Direction Towards Target:
The direction vector is calculated by subtracting the meteor's position from the target position (the center of the screen). This vector is then normalized (converted to a unit vector) to get the normalizedDirection.

Movement Vector:
The movement vector is calculated by multiplying the normalizedDirection by the velocity. This vector determines how the meteor will move in each frame.

Random Rotation Speed:
A random rotation speed is calculated for the meteor. It lies between rotationSpeedMin and rotationSpeedMax, providing a small rotational movement to the meteor.

Sprite Selection:
A sprite is randomly selected from the assets.MeteorSprites array to visually represent the meteor.

Meteor Object Creation:
A new Meteor struct is instantiated with the calculated position, movement, rotation speed, and selected sprite. This struct represents a single meteor in the game.

Return Meteor Object:
Finally, the newly created Meteor object is returned.
*/
package game

import (
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"

	"game/assets"
)

const (
	rotationSpeedMin = -0.01
	rotationSpeedMax = 0.04
)

type Meteor struct {
	position      Vector
	rotation      float64
	movement      Vector
	rotationSpeed float64
	sprite        *ebiten.Image
}

func NewMeteor(baseVelocity float64) *Meteor {
	target := Vector{
		X: screenWidth / 2,
		Y: screenHeight / 2,
	}
	// Position the meteor at a random Y position on the right edge of the screen
	//    <--- o	|
	//    <--- o	|
	//    <--- o	|
	// posY := rand.Float64() * float64(screenHeight)
	// pos := Vector{
	//     X: float64(screenWidth), // Right edge of the screen
	//     Y: posY,
	// }
	angle := rand.Float64() * 2 * math.Pi
	r := screenWidth / 2.0

	pos := Vector{
		X: target.X + math.Cos(angle)*r,
		Y: target.Y + math.Sin(angle)*r,
	}

	velocity := baseVelocity + rand.Float64()*1.5

	direction := Vector{
		X: target.X - pos.X,
		Y: target.Y - pos.Y,
	}
	normalizedDirection := direction.Normalize()

	movement := Vector{
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

func (m *Meteor) Collider() Rect {
	bounds := m.sprite.Bounds()

	return NewRect(
		m.position.X,
		m.position.Y,
		float64(bounds.Dx()),
		float64(bounds.Dy()),
	)
}
