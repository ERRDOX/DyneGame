package game

import "github.com/hajimehoshi/ebiten/v2"

func (p *Player) DrawShadow(screen *ebiten.Image) {
	bounds := p.sprite[1].Bounds()
	halfW := float64(bounds.Dx()) / 2
	halfH := float64(bounds.Dy()) / 2

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-halfW, -halfH)
	op.GeoM.Rotate(p.rotation)
	op.GeoM.Translate(p.position.X-20, p.position.Y+10)
	op.ColorScale.Scale(0, 0, 0, 0.5)
	screen.DrawImage(p.sprite[0], op)
}
