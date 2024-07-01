package game

import (
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"

	"dynegame/assets"
)

type Obstacle struct {
	position Vector
	sprite   *ebiten.Image
}

// has base
//     _         _     _
//____|_|_______|_|___|_|___
// func NewObstacle(n int) []*Obstacle {
// 	obstacles := make([]*Obstacle, n)
// 	pos := Vector{
// 		X: screenWidth / 2,
// 		Y: screenHeight,
// 	}
// 	for i := 0; i < n; i++ {
// 		sprite := assets.Obstacle[rand.Intn(len(assets.Obstacle))]
// 		o := &Obstacle{
// 			position: pos,
// 			sprite:   sprite,
// 		}
// 		chunk := screenWidth / n
// 		o.position.Y = o.position.Y - float64(sprite.Bounds().Dy())
// 		nextNum := rand.Intn(chunk)
// 		o.position.X = float64(i*chunk + nextNum)
// 		obstacles[i] = o
// 	}

//		return obstacles
//	}
func NewObstacle(n int) []*Obstacle {
	obstacles := make([]*Obstacle, n)
	pos := Vector{
		X: screenWidth / 2,
		Y: screenHeight,
	}
	for i := 0; i < n; i++ {
		sprite := assets.Obstacle[rand.Intn(len(assets.Obstacle))]
		o := &Obstacle{
			position: pos,
			sprite:   sprite,
		}
		chunk := screenWidth / n
		o.position.Y = o.position.Y - float64(sprite.Bounds().Dy())
		nextNum := rand.Intn(chunk)
		o.position.X = float64(i*chunk + nextNum)
		o.position.Y = float64(i*chunk+nextNum) * 0.9
		obstacles[i] = o
	}

	return obstacles
}
func NewMaptoObstacle(m [][]uint8) []*Obstacle {
	sprite := assets.Obstacle[rand.Intn(len(assets.Obstacle))]
	// height := screenHeight / 100
	// width := screenWidth / 100

	// screenWidth, screenHeight := sprite.Bounds().Dx(), sprite.Bounds().Dy()
	var obstacles []*Obstacle
	// pos := Vector{
	// 	X: screenWidth,
	// 	Y: screenHeight,
	// }
	// o := &Obstacle{
	// 	position: pos,
	// 	sprite:   sprite,
	// }
	for i, row := range m {
		// o.position.Y = o.position.Y - float64(sprite.Bounds().Dy())
		for j, value := range row {
			if value == 1 {
				obstacles = append(obstacles, &Obstacle{
					position: Vector{
						X: float64(j*100 - sprite.Bounds().Dx()),
						Y: float64(i*100 - sprite.Bounds().Dy()),
					},
					sprite: sprite,
				})
			}
		}
	}
	println(&obstacles)

	return obstacles

}
func (o *Obstacle) Update() {

}

func (o *Obstacle) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}

	op.GeoM.Translate(o.position.X, o.position.Y)

	screen.DrawImage(o.sprite, op)
}

func (o *Obstacle) Collider() Rect {
	bounds := o.sprite.Bounds()
	return NewRect(
		o.position.X,
		o.position.Y,
		float64(bounds.Dx()),
		float64(bounds.Dy()),
	)
}
