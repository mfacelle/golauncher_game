package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
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

func (o *ChargedObject) IntersectsParticle(projectile *BasicParticle) bool {
	if projectile == nil {
		return false
	}

	dx := projectile.positionPx.X - o.positionPx.X
	dy := projectile.positionPx.Y - o.positionPx.Y
	distanceSq := dx*dx + dy*dy
	// ignoring projectile size for now. looks better with sprites this way
	collisionRadius := o.radius //+ projectile.radiusPx

	return distanceSq <= collisionRadius*collisionRadius
}

func (o *ChargedObject) Draw(screen *ebiten.Image) {
	// vector.FillCircle(screen, float32(o.positionPx.X), float32(o.positionPx.Y), float32(o.radius), o.color, true)
	// vector.StrokeCircle(screen, float32(o.positionPx.X), float32(o.positionPx.Y), float32(o.radius), float32(o.borderWidthPx), o.borderColor, true)

	// scale sprite up to fit radius of object
	drawOpts := &ebiten.DrawImageOptions{}
	bounds := o.sprite.Bounds()
	drawOpts.GeoM.Scale((o.radius*2.0)/float64(bounds.Dx()), (o.radius*2.0)/float64(bounds.Dy()))

	// need to offset by radius. circle is draw with position at center,
	// sprite drawn with position at top-left
	drawOpts.GeoM.Translate(o.positionPx.X-o.radius, o.positionPx.Y-o.radius)
	screen.DrawImage(o.sprite, drawOpts)
}
