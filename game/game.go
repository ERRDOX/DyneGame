package game

import (
	"fmt"
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"

	"dynegame/assets"
)

const (
	ACT_SERVER_CONN_HOST    = "localhost"
	ACT_SERVER_CONN_PORT    = "27199" //socket port
	ACT_SERVER_CONN_TYPE    = "tcp"
	STATUS_SERVER_CONN_HOST = "localhost"
	STATUS_SERVER_CONN_PORT = "27198" //socket port
	STATUS_SERVER_CONN_TYPE = "tcp"
	MAP                     = "DragonMap"

	screenWidth  = 1800
	screenHeight = 900
	// screenWidth     = 800
	// screenHeight    = 600
	meteorSpawnTime = 1 * time.Second

	baseMeteorVelocity  = 0.25
	meteorSpeedUpAmount = 0.01
	meteorSpeedUpTime   = 5 * time.Second
	//   1 >=BoundsDecreaseRatio> 0
	humanBoundsDecreaseRatio  = 1.0
	objectBoundsDecreaseRatio = 1.0
	bulletBoundsDecreaseRatio = 1.0 //nothing have changed if you set 1.0
	FIXBoundsDecreaseRatio    = 1.0
)

type Game struct {
	player           *Player
	SecondPlayer     *SecondPlayer
	Action           *Action
	meteorSpawnTimer *Timer
	Explosion        []*Explosion
	obstacle         []*Obstacle
	meteors          []*Meteor
	bullets          []*Bullet
	score            int
	baseVelocity     float64
	velocityTimer    *Timer
}

func NewGame() *Game {
	g := &Game{
		meteorSpawnTimer: NewTimer(meteorSpawnTime),
		baseVelocity:     baseMeteorVelocity,
		velocityTimer:    NewTimer(meteorSpeedUpTime),
	}

	g.obstacle = append(g.obstacle, NewMaptoObstacle(DragonMap)...)
	g.Action = NewAction()
	g.SecondPlayer = NewSecondPlayer(g)
	g.player = NewPlayer(g)
	// go g.Action.Server()
	go g.respBullet()

	return g
}

func (g *Game) Update() error {
	g.velocityTimer.Update()
	if g.velocityTimer.IsReady() {
		g.velocityTimer.Reset()
		g.baseVelocity += meteorSpeedUpAmount
	}

	// g.Action.Joiner()
	g.Explosion = nil
	g.SecondPlayer.Update(g)
	g.player.Update(g)

	// g.meteorSpawnTimer.Update()
	// if g.meteorSpawnTimer.IsReady() {
	// 	g.meteorSpawnTimer.Reset()

	// 	// m := NewMeteor(g.baseVelocity, g.player)
	// 	g.meteors = append(g.meteors, m)
	// }

	// for _, m := range g.meteors {
	// 	m.Update()
	// }
	for _, b := range g.bullets {
		b.Update()
	}
	for _, b := range g.player.bullet {
		b.Update()
	}
	for _, b := range g.SecondPlayer.bullet {
		b.Update()
	}
	// println("BULLET: ", len(g.bullets))
	// println("BULLETP: ", len(g.player.bullet))
	// println("BULLETSP: ", len(g.SecondPlayer.bullet))

	// bulletOutofScreen := func() {
	// 	for i, b := range g.bullets {
	// 		if 0 > b.position.X || b.position.X > screenHeight || 0 > b.position.Y || b.position.Y > screenWidth {
	// 			// println("BULLET X: ", g.meteors[i].position.X)
	// 			// println("BULLET Y: ", g.meteors[i].position.Y)
	// 			g.bullets = append(g.bullets[:i], g.bullets[i+1:]...)
	// 		}
	// 	}
	// }
	// go bulletOutofScreen()
	bulletOutofScreenPlayer := func() {
		time.Sleep(3 * time.Second)
		for i := len(g.player.bullet) - 1; i >= 0; i-- {
			b := g.player.bullet[i]
			if 0 > b.position.X || b.position.Y > screenHeight || 0 > b.position.Y || b.position.Y > screenWidth {
				println("BULLET X: %f ", g.player.bullet[i].position.X)
				println("BULLET Y: %f ", g.player.bullet[i].position.Y)
				g.player.bullet = append(g.player.bullet[:i], g.player.bullet[i+1:]...)
			}
		}
		for i := len(g.SecondPlayer.bullet) - 1; i >= 0; i-- {
			b := g.SecondPlayer.bullet[i]
			if 0 > b.position.X || b.position.Y > screenHeight || 0 > b.position.Y || b.position.Y > screenWidth {
				println("BULLET X: %f ", g.SecondPlayer.bullet[i].position.X)
				println("BULLET Y: %f ", g.SecondPlayer.bullet[i].position.Y)
				g.SecondPlayer.bullet = append(g.SecondPlayer.bullet[:i], g.SecondPlayer.bullet[i+1:]...)
			}
		}
	}
	go bulletOutofScreenPlayer()

	// bulletOutofScreenSecondPlayer := func() {
	// 	for i := len(g.SecondPlayer.bullet) - 1; i >= 0; i-- {
	// 		b := g.SecondPlayer.bullet[i]
	// 		if 0 > b.position.X || b.position.X > screenHeight || 0 > b.position.Y || b.position.Y > screenWidth {
	// 			println("BULLET X: ", g.SecondPlayer.bullet[i].position.X)
	// 			println("BULLET Y: ", g.SecondPlayer.bullet[i].position.Y)
	// 			g.SecondPlayer.bullet = append(g.SecondPlayer.bullet[:i], g.SecondPlayer.bullet[i+1:]...)
	// 		}
	// 	}
	// }
	// go bulletOutofScreenSecondPlayer()

	// meteorOutofScreen := func() {
	// 	for i, m := range g.meteors {
	// 		if 0 > m.position.X || m.position.X > screenHeight || 0 > m.position.Y || m.position.Y > screenWidth {
	// 			// println("METEOR X: ", g.meteors[i].position.X)
	// 			// println("METEOR Y: ", g.meteors[i].position.Y)
	// 			g.meteors = append(g.meteors[:i], g.meteors[i+1:]...)
	// 		}
	// 	}
	// }
	// go meteorOutofScreen()

	// Check for bullet collisions
	// bulletObjectCollisions := func() {
	// 	for i, m := range g.meteors {
	// 		for j, b := range g.bullets {
	// 			if m.Collider(objectBoundsDecreaseRatio).Intersects(b.Collider(bulletBoundsDecreaseRatio)) {
	// 				g.meteors = append(g.meteors[:i], g.meteors[i+1:]...)
	// 				g.bullets = append(g.bullets[:j], g.bullets[j+1:]...)
	// 				g.score++
	// 			}
	// 		}
	// 	}
	// }
	// go bulletObjectCollisions()
	// player hit secondplayer
	bulletPlayerCollisions := func() {

		for i, b := range g.player.bullet {
			if len(g.player.bullet) == 0 {
				continue
			}
			if b.Collider(bulletBoundsDecreaseRatio).Intersects(g.SecondPlayer.Collider(humanBoundsDecreaseRatio)) {
				g.player.bullet = append(g.player.bullet[:i], g.player.bullet[i+1:]...)
				g.Explosion = append(g.Explosion, NewExplosion(b.position))
				g.player.score++
			}
			for _, o := range g.obstacle {
				if b.Collider(bulletBoundsDecreaseRatio).Intersects(o.Collider()) {
					g.player.bullet = append(g.player.bullet[:i], g.player.bullet[i+1:]...)
					g.Explosion = append(g.Explosion, NewExplosion(b.position))
				}
			}

		}

	}
	go bulletPlayerCollisions()
	// SecondPalyer hit player
	bulletSecondPlayerCollisions := func() {
		for i, b := range g.SecondPlayer.bullet {
			if b.Collider(bulletBoundsDecreaseRatio).Intersects(g.player.Collider(humanBoundsDecreaseRatio)) {
				g.SecondPlayer.bullet = append(g.SecondPlayer.bullet[:i], g.SecondPlayer.bullet[i+1:]...)
				g.Explosion = append(g.Explosion, NewExplosion(b.position))
				g.SecondPlayer.score++
			}
			for _, o := range g.obstacle {
				if b.Collider(bulletBoundsDecreaseRatio).Intersects(o.Collider()) {
					g.SecondPlayer.bullet = append(g.SecondPlayer.bullet[:i], g.SecondPlayer.bullet[i+1:]...)
					g.Explosion = append(g.Explosion, NewExplosion(b.position))
				}
			}
		}
	}
	// if bullet contanct the obstackls
	// for i, b := range g.SecondPlayer.bullet {
	// 	for _, o := range g.obstacle {
	// 		if b.Collider(bulletBoundsDecreaseRatio).Intersects(o.Collider()) {
	// 			g.SecondPlayer.bullet = append(g.SecondPlayer.bullet[:i], g.SecondPlayer.bullet[i+1:]...)
	// 			// g.obstacle = append(g.obstacle[:j], g.obstacle[j+1:]...)
	// 		}
	// 	}
	// }

	go bulletSecondPlayerCollisions()

	// Check for meteor/player collisions
	// objectHumanCollisions := func() {
	// 	for _, m := range g.meteors {
	// 		if m.Collider(objectBoundsDecreaseRatio).Intersects(g.player.Collider(humanBoundsDecreaseRatio)) {
	// 			g.Reset()
	// 			break
	// 		}
	// 	}
	// }
	// go objectHumanCollisions()
	// TODO: return Termination
	return nil
}
func (g *Game) Draw(screen *ebiten.Image) {
	Background(screen)
	g.player.DrawShadow(screen)
	g.player.Draw(screen)

	g.SecondPlayer.DrawShadow(screen)
	g.SecondPlayer.Draw(screen)

	// for _, m := range g.meteors {
	// 	m.Draw(screen)
	// }

	for _, b := range g.player.bullet {
		b.Draw(screen)
	}
	for _, b := range g.SecondPlayer.bullet {
		b.Draw(screen)
	}
	for _, o := range g.obstacle {
		o.Draw(screen)
	}
	//TODO: Below for loop makes lag on the rendering send it to background
	// for _, e := range g.Explosion {
	// 	go e.Draw(screen)
	// }
	text.Draw(screen, fmt.Sprintf("%06d", g.SecondPlayer.score), assets.ScoreFont, screenWidth/4, 50, color.RGBA{128, 128, 128, 255})
	text.Draw(screen, fmt.Sprintf("%06d", g.player.score), assets.ScoreFont, 3*screenWidth/4, 50, color.RGBA{128, 128, 128, 255})
}

//	func (g *Game) AddBullet(b *Bullet) {
//		g.bullets = append(g.bullets, b)
//	}
func (g *Game) AddBulletPlayer(b *Bullet) {
	g.player.bullet = append(g.player.bullet, b)
}
func (g *Game) AddBulletSecondPlayer(b *Bullet) {
	g.SecondPlayer.bullet = append(g.SecondPlayer.bullet, b)
}

//	func (g *Game) AddBulletSP(b *Bullet) {
//		g.SecondPlayer.bullet = append(g.SecondPlayer.bullet, b)
//	}
func (g *Game) Reset() {
	g.player = NewPlayer(g)
	g.SecondPlayer = NewSecondPlayer(g)
	g.Action.Act = make(map[string]bool)
	g.meteors = nil
	g.player.bullet = nil
	g.SecondPlayer.bullet = nil
	g.Explosion = nil
	g.bullets = nil
	g.score = 0
	// g.meteorSpawnTimer.Reset()
	g.baseVelocity = baseMeteorVelocity
	g.velocityTimer.Reset()
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}
