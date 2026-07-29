package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// obstacle is similar to ChargedObject, but is rectangular and has no charge.
// considering adding a charge to some eventually, but would need to figure out how to best calculate
// attraction/repulsion forces
type Obstacle struct {
	positionPx    Vector
	widthPx       float64
	heightPx      float64
	color         color.RGBA
	borderColor   color.RGBA
	borderWidthPx float64
}

func (o *Obstacle) Update() error {
	return nil
}

func (o *Obstacle) IntersectsProjectile(projectile *BasicParticle) bool {
	if projectile == nil {
		return false
	}

	return projectile.positionPx.X >= o.positionPx.X &&
		projectile.positionPx.X <= o.positionPx.X+o.widthPx &&
		projectile.positionPx.Y >= o.positionPx.Y &&
		projectile.positionPx.Y <= o.positionPx.Y+o.heightPx
}

func (o *Obstacle) Draw(screen *ebiten.Image) {
	vector.FillRect(screen, float32(o.positionPx.X), float32(o.positionPx.Y), float32(o.widthPx), float32(o.heightPx), o.color, true)
	vector.StrokeRect(screen, float32(o.positionPx.X), float32(o.positionPx.Y), float32(o.widthPx), float32(o.heightPx), float32(o.borderWidthPx), o.borderColor, true)

}
