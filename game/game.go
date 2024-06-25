//go:build ignore
// +build ignore

package game

import (
	"fmt"
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"

	"game/assets"
)

const (
	screenWidth  = 1920
	screenHeight = 1080
	// screenWidth     = 800
	// screenHeight    = 600
	meteorSpawnTime = 1 * time.Second

	baseMeteorVelocity  = 0.25
	meteorSpeedUpAmount = 0.1
	meteorSpeedUpTime   = 5 * time.Second
	//   1 >=BoundsDecreaseRatio> 0
	humanBoundsDecreaseRatio  = 0.8
	objectBoundsDecreaseRatio = 0.7
	bulletBoundsDecreaseRatio = 1.0 //nothing have changed if you set 1.0
)

type Game struct {
	player           *Player
	meteorSpawnTimer *Timer
	meteors          []*Meteor
	bullets          []*Bullet
	score            int

	baseVelocity  float64
	velocityTimer *Timer
}

func NewGame() *Game {
	g := &Game{
		meteorSpawnTimer: NewTimer(meteorSpawnTime),
		baseVelocity:     baseMeteorVelocity,
		velocityTimer:    NewTimer(meteorSpeedUpTime),
	}

	g.player = NewPlayer()

	return g
}

func (g *Game) Update() error {
	g.velocityTimer.Update()
	if g.velocityTimer.IsReady() {
		g.velocityTimer.Reset()
		g.baseVelocity += meteorSpeedUpAmount
	}

	g.player.Update()

	g.meteorSpawnTimer.Update()
	if g.meteorSpawnTimer.IsReady() {
		g.meteorSpawnTimer.Reset()

		m := NewMeteor(g.baseVelocity, g.player)
		g.meteors = append(g.meteors, m)
	}

	for _, m := range g.meteors {
		m.Update()
	}

	for _, b := range g.bullets {
		b.Update()
	}
	bulletOutofScreen := func() {
		for i, b := range g.bullets {
			if 0 > b.position.X || b.position.X > screenWidth || 0 > b.position.Y || b.position.X > screenHeight {
				g.bullets = append(g.bullets[:i], g.bullets[i+1:]...)
			}
		}
	}
	go bulleutofScreen()
	meteorOutofScreen := func() {
		for i, b := range g.meteors {
			if 0 > b.position.X || b.position.X > screenWidth || 0 > b.position.Y || b.position.X > screenHeight {
				g.meteors = append(g.meteors[:i], g.meteors[i+1:]...)
			}
		}
	}
	go meteorOutofScreen()
	// Check for bullet collisions
	bulletObjectCollisions := func() {
		for i, m := range g.meteors {
			for j, b := range g.bullets {
				if m.Collider(objectBoundsDecreaseRatio).Intersects(b.Collider(bulletBoundsDecreaseRatio)) {
					g.meteors = append(g.meteors[:i], g.meteors[i+1:]...)
					g.bullets = append(g.bullets[:j], g.bullets[j+1:]...)
					g.score++
				}
			}
		}
	}
	go bulletObjectCollisions()
	// Check for meteor/player collisions
	objectHumanCollisions := func() {
		for _, m := range g.meteors {
			if m.Collider(objectBoundsDecreaseRatio).Intersects(g.player.Collider(humanBoundsDecreaseRatio)) {
				g.Reset()
				break
			}
		}
	}
	go objectHumanCollisions()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.DrawImage(assets.BackGround, nil)
	g.player.Draw(screen)

	for _, m := range g.meteors {
		m.Draw(screen)
	}

	for _, b := range g.bullets {
		b.Draw(screen)
	}

	text.Draw(screen, fmt.Sprintf("%06d", g.score), assets.ScoreFont, screenWidth/2-100, 50, color.White)
}

func (g *Game) AddBullet(b *Bullet) {
	g.bullets = append(g.bullets, b)
}

func (g *Game) Reset() {
	g.obstacle = NewObstacle()
	g.player = NewPlayer(g)
	g.meteors = nil
	g.bullets = nil
	g.score = 0
	g.meteorSpawnTimer.Reset()
	g.baseVelocity = baseMeteorVelocity
	g.velocityTimer.Reset()
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}
