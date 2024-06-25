//go:build ignore
// +build ignore

package game

import (
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"

	"game/assets"
)

const (
	shootCooldown     = time.Millisecond * 400
	rotationPerSecond = 1.2 * math.Pi

	bulletSpawnOffset = 1.0
	sprintSpeed       = 2
)

type Player struct {
	game *Game

	position Vector
	rotation float64
	sprite   []*ebiten.Image

	animationSpeed      float64
	animationTimer      float64
	playerFramePosition int

	shootCooldown *Timer
}

func NewPlayer(game *Game) *Player {
	sprite := assets.Humanplayer

	bounds := sprite[0].Bounds()
	halfW := float64(bounds.Dx()) / 2
	halfH := float64(bounds.Dy()) / 2

	pos := Vector{
		X: screenWidth/2 - halfW,
		Y: screenHeight/2 - halfH,
	}

	return &Player{
		game:           game,
		position:       pos,
		rotation:       0,
		sprite:         sprite,
		animationSpeed: 0.1,
		animationTimer: 0,
		shootCooldown:  NewTimer(shootCooldown),
	}
}

// Update updates the player's position, rotation
func (p *Player) Update(g *Game) {
	p.animationTimer += 1.2 / float64(ebiten.TPS())
	if p.animationTimer >= p.animationSpeed {
		p.animationTimer = 0

		if ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyD) {
			p.playerFramePosition = (p.playerFramePosition + 1) % len(p.sprite)
		}
	}

	rotateSpeed := rotationPerSecond / float64(ebiten.TPS())

	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		p.rotation -= rotateSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		p.rotation += rotateSpeed
	}

	if ebiten.IsKeyPressed(ebiten.KeyA) {
		p.position.X -= sprintSpeed
		if p.playerObstacleCollisions(g) || p.peripheralCollision() {
			p.position.X += sprintSpeed
		}

	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		p.position.X += sprintSpeed
		if p.playerObstacleCollisions(g) || p.peripheralCollision() {
			p.position.X -= sprintSpeed
		}
	}

	if ebiten.IsKeyPressed(ebiten.KeyW) {
		p.position.Y -= sprintSpeed
		if p.playerObstacleCollisions(g) || p.peripheralCollision() {
			p.position.Y += sprintSpeed
		}
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		p.position.Y += sprintSpeed
		if p.playerObstacleCollisions(g) || p.peripheralCollision() {
			p.position.Y -= sprintSpeed
		}
	}
	p.shootCooldown.Update()
	if p.shootCooldown.IsReady() && ebiten.IsKeyPressed(ebiten.KeySpace) {
		p.shootCooldown.Reset()

		bounds := p.sprite[1].Bounds()
		halfW := float64(bounds.Dx()) / 2
		halfH := float64(bounds.Dy()) / 4

		spawnPos := Vector{
			X: p.position.X + halfW + math.Sin(p.rotation)*bulletSpawnOffset,
			Y: p.position.Y + halfH + math.Cos(p.rotation)*-bulletSpawnOffset,
		}

		bullet := NewBullet(spawnPos, p.rotation)
		p.game.AddBullet(bullet)
	}
}

// Draw renders the player sprite onto the screen.
//
// It takes a screen image as input, and uses the player's position and sprite
// to calculate the appropriate drawing options. The function then draws the
// sprite onto the screen using these options.
func (p *Player) Draw(screen *ebiten.Image) {

	op := &ebiten.DrawImageOptions{}

	op.GeoM.Translate(p.position.X, p.position.Y)

	screen.DrawImage(p.sprite[p.playerFramePosition], op)
}

// Collider returns a Rect representing the bounds of the player's sprite.
//
// It does not take any parameters.
// It returns a Rect.
func (p *Player) Collider(BoundsDecreaseRatio float64) Rect {
	bounds := p.sprite[p.playerFramePosition].Bounds()

	return NewRect(
		p.position.X,
		p.position.Y,
		float64(bounds.Dx())*BoundsDecreaseRatio,
		float64(bounds.Dy())*BoundsDecreaseRatio,
	)
}
func (p *Player) playerObstacleCollisions(g *Game) bool {
	for _, m := range g.obstacle {
		if m.Collider().Intersects(p.Collider(humanBoundsDecreaseRatio)) {
			return true
		}
	}
	return false
}

func (p *Player) peripheralCollision() bool {
	if p.position.Y < 0 || p.position.Y > screenHeight {
		return true
	}
	if p.position.X < 0 || p.position.X > screenHeight {
		return true
	}
	return false
}
