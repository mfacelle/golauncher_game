package game

import (
	"image/color"
	"math"
	"staticsheep/assets"

	"github.com/hajimehoshi/ebiten/v2"
)

type BasicParticle struct {
	positionPx Vector
	velocityPx Vector
	charge     float64
	radiusPx   float64
	color      color.RGBA
	sprite     *ebiten.Image
}

const (
	initialVelocity       = 180.0
	electrostaticConstant = 100 // 8.98e9 in the real world
	softening             = 1.0
	maxDistance           = 600.0
)

// probably want to return a regular object, not a pointer. need to look into this more
func NewBasicParticle(posPx Vector, rotationRad float64) *BasicParticle {
	return &BasicParticle{
		positionPx: posPx,
		velocityPx: Vector{
			X: initialVelocity * math.Cos(rotationRad),
			Y: initialVelocity * math.Sin(rotationRad),
		},
		charge:   -10,
		radiusPx: 4,
		color:    color.RGBA{R: 100, G: 100, B: 255, A: 255},
		sprite:   assets.FluffBigSprite,
	}
}

func (p *BasicParticle) Update(objects []ChargedObject) bool {
	// this could definitely be better
	collision := false

	accelerationX := 0.0
	accelerationY := 0.0

	for _, object := range objects {
		// check for collision and stop updating if one occurs
		if object.IntersectsParticle(p) {
			p.velocityPx = Vector{}
			collision = true
			accelerationX = 0.0
			accelerationY = 0.0
			break
		}

		// skip over any obstacles with no charge
		if object.charge == 0 {
			continue
		}

		// get diff in distance, and applying softening to avoid out of control forces
		dx := object.positionPx.X - p.positionPx.X
		dy := object.positionPx.Y - p.positionPx.Y
		distanceSq := dx*dx + dy*dy + softening*softening
		if distanceSq < maxDistance*maxDistance {
			distance := math.Sqrt(distanceSq)

			force := electrostaticConstant * math.Abs(p.charge*object.charge) / distanceSq

			// if charges are the same, force should repel
			if p.charge*object.charge > 0 {
				force *= -1
			}

			// update acceleration, to be applied after all objects evaluated
			accelerationX += force * dx / distance
			accelerationY += force * dy / distance
			//log.Print("distanceSq: ", distanceSq, "force: ", force, " distance: ", distance, " dx: ", dx, " dy: ", dy)
		}
	}

	// update velocity, then position (convert from per second to per tick)
	p.velocityPx.X += accelerationX / float64(ebiten.TPS())
	p.velocityPx.Y += accelerationY / float64(ebiten.TPS())
	p.positionPx.X += p.velocityPx.X / float64(ebiten.TPS())
	p.positionPx.Y += p.velocityPx.Y / float64(ebiten.TPS())

	return collision
}

func (p *BasicParticle) Draw(screen *ebiten.Image) {
	// vector.FillCircle(screen, float32(p.positionPx.X), float32(p.positionPx.Y), float32(p.radiusPx), p.color, false)

	drawOpts := &ebiten.DrawImageOptions{}
	// shift by radius to center sprite on particle
	drawOpts.GeoM.Translate(p.positionPx.X-p.radiusPx, p.positionPx.Y-p.radiusPx)
	screen.DrawImage(p.sprite, drawOpts)
}
