package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type ChargedObject struct {
	positionPx Vector
	radius     float64
	color      color.RGBA
	charge     float64
}

func (o *ChargedObject) Update() error {
	return nil
}

func (o *ChargedObject) IntersectsProjectile(projectile *BasicParticle) bool {
	if projectile == nil {
		return false
	}

	dx := projectile.positionPx.X - o.positionPx.X
	dy := projectile.positionPx.Y - o.positionPx.Y
	distanceSq := dx*dx + dy*dy
	collisionRadius := o.radius + projectile.radiusPx

	return distanceSq <= collisionRadius*collisionRadius
}

func (o *ChargedObject) Draw(screen *ebiten.Image) {
	vector.FillCircle(screen, float32(o.positionPx.X), float32(o.positionPx.Y), float32(o.radius), o.color, true)
}
