package game

import (
	"math"
	"time"

	"github.com/hajimehoshi/ebiten/v2"

	"dynegame/assets"
	"dynegame/utils"
)

const (
	shootCooldown     = time.Millisecond * 400
	rotationPerSecond = 1.1 * math.Pi

	bulletSpawnOffset = 1.0
	sprintSpeed       = 4
)

type Player struct {
	game     *Game
	score    int
	position utils.Vector
	rotation float64
	sprite   []*ebiten.Image
	bullet   []*utils.Bullet

	animationSpeed      float64
	animationTimer      float64
	playerFramePosition int
	shootCooldown       *utils.Timer
}

func NewPlayer(game *Game) *Player {
	sprite := assets.PlanePlayer

	bounds := sprite[1].Bounds()
	halfW := float64(bounds.Dx()) / 2
	halfH := float64(bounds.Dy()) / 2

	pos := utils.Vector{
		X: 4 * halfW,
		Y: screenHeight/2 - halfH,
	}

	return &Player{
		game:           game,
		position:       pos,
		rotation:       0,
		sprite:         sprite,
		animationSpeed: 0.1,
		animationTimer: 0,
		shootCooldown:  utils.NewTimer(shootCooldown),
	}
}

// Update updates the player's position, rotation
func (p *Player) Update(g *Game) {
	// p.animationTimer += 1.2 / float64(ebiten.TPS())
	// println(p.animationTimer, "animation speed : ", p.animationSpeed)
	// if p.animationTimer >= p.animationSpeed {
	// 	p.animationTimer = 0

	// 	if ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyD) {
	// 		p.playerFramePosition = (p.playerFramePosition + 1) % len(p.sprite)
	// 	}
	// }

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
		halfH := float64(bounds.Dy()) / 2

		spawnPos := utils.Vector{
			X: p.position.X + halfW + math.Sin(p.rotation)*bulletSpawnOffset,
			Y: p.position.Y + halfH + math.Cos(p.rotation)*-bulletSpawnOffset,
		}

		bullet := utils.NewBullet(spawnPos, p.rotation)
		// p.bullet = append(p.bullet, bullet)
		p.game.AddBulletPlayer(bullet)
	}

}

// Draw renders the player sprite onto the screen.
//
// It takes a screen image as input, and uses the player's position and sprite
// to calculate the appropriate drawing options. The function then draws the
// sprite onto the screen using these options.

func (p *Player) Draw(screen *ebiten.Image) {
	bounds := p.sprite[0].Bounds()
	halfW := float64(bounds.Dx()) / 2
	halfH := float64(bounds.Dy()) / 2

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-halfW, -halfH)
	op.GeoM.Rotate(p.rotation)
	op.GeoM.Translate(halfW, halfH)
	op.GeoM.Translate(p.position.X, p.position.Y)

	// op.GeoM.Translate(m.position.X, m.position.Y)
	// op.GeoM.Translate(p.position.X, p.position.Y)
	// op.GeoM.Rotate(p.rotation)
	screen.DrawImage(p.sprite[p.playerFramePosition], op)
}

// Collider returns a Rect representing the bounds of the player's sprite.
//
// It does not take any parameters.
// It returns a Rect.
func (p *Player) Collider(BoundsDecreaseRatio float64) utils.Rect {
	bounds := p.sprite[p.playerFramePosition].Bounds()

	return utils.NewRect(
		p.position.X,
		p.position.Y,
		float64(bounds.Dx())*BoundsDecreaseRatio,
		float64(bounds.Dy())*BoundsDecreaseRatio,
	)
}

// playerObstacleCollisions checks for collisions between the player and obstacles in the game.
//
// Parameters:
// - g: A pointer to the Game struct.
//
// Returns:
// - A boolean value indicating whether a collision occurred or not.
func (p *Player) playerObstacleCollisions(g *Game) bool {
	for _, m := range g.obstacle {
		if m.Collider().Intersects(p.Collider(humanBoundsDecreaseRatio)) {
			return true
		}
	}
	return false

}

// peripheralCollision checks if the player is colliding with the screen borders.
//
// It takes no parameters.
// Returns a boolean value indicating whether there is a collision or not.
func (p *Player) peripheralCollision() bool {
	bounds := p.sprite[0].Bounds()
	halfW := float64(bounds.Dx())
	halfH := float64(bounds.Dy())
	if p.position.Y < 0 || p.position.Y > screenHeight-halfH {
		return true
	}
	if p.position.X < 0 || p.position.X > screenWidth-halfW {
		return true
	}
	return false
}
