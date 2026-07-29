package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// represents an object that has a charge that can attract/repel particles launched by the player.
type ChargedObject struct {
	positionPx    Vector
	radius        float64
	color         color.RGBA
	borderColor   color.RGBA
	borderWidthPx float64
	charge        float64
	sprite        *ebiten.Image
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
	// vector.FillCircle(screen, float32(o.positionPx.X), float32(o.positionPx.Y), float32(o.radius), o.color, true)
	vector.StrokeCircle(screen, float32(o.positionPx.X), float32(o.positionPx.Y), float32(o.radius), float32(o.borderWidthPx), o.borderColor, true)

	bounds := o.sprite.Bounds()
	halfWidth := float64(bounds.Dx()) / 2
	halfHeight := float64(bounds.Dy()) / 2

	drawOpts := &ebiten.DrawImageOptions{}
	// scale sprite up to fit radius of object
	drawOpts.GeoM.Translate(-halfWidth, -halfHeight)
	drawOpts.GeoM.Scale((o.radius*2.0)/float64(bounds.Dx()), (o.radius*2.0)/float64(bounds.Dy()))
	drawOpts.GeoM.Translate(halfWidth, halfHeight)

	// log.Printf("dx=%d, dy=%d, radius*2=%.2f, scaledX=%.2f", bounds.Dx(), bounds.Dy(), o.radius*2.0, (o.radius*2.0)/float64(bounds.Dx()))

	drawOpts.GeoM.Translate(o.positionPx.X-float64(bounds.Dx()), o.positionPx.Y-float64(bounds.Dy()))
	screen.DrawImage(o.sprite, drawOpts)
}
