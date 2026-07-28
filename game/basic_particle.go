package game

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type BasicParticle struct {
	positionPx Vector
	velocityPx Vector
	charge     float64
	radiusPx   float64
	color      color.RGBA
}

const (
	initialVelocity       = 180.0
	electrostaticConstant = 100 // really 8.98e9
	softening             = 1.0
	maxDistance           = 800.0
)

func NewBasicParticle(posPx Vector, rotationRad float64) *BasicParticle {
	return &BasicParticle{
		positionPx: posPx,
		velocityPx: Vector{
			X: initialVelocity * math.Cos(rotationRad),
			Y: initialVelocity * math.Sin(rotationRad),
		},
		charge:   -10,
		radiusPx: 2,
		color:    color.RGBA{R: 100, G: 100, B: 255, A: 255},
	}
}

func (p *BasicParticle) Update(objects []ChargedObject) bool {
	// this could definitely be better
	collision := false

	accelerationX := 0.0
	accelerationY := 0.0

	for _, object := range objects {
		// check for collision and don't update if one occurs
		if object.IntersectsProjectile(p) {
			p.velocityPx = Vector{}
			collision = true
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

	// update velocity, then position
	p.velocityPx.X += accelerationX / float64(ebiten.TPS())
	p.velocityPx.Y += accelerationY / float64(ebiten.TPS())
	p.positionPx.X += p.velocityPx.X / float64(ebiten.TPS())
	p.positionPx.Y += p.velocityPx.Y / float64(ebiten.TPS())

	return collision
}

func (p *BasicParticle) Draw(screen *ebiten.Image) {
	vector.FillCircle(screen, float32(p.positionPx.X), float32(p.positionPx.Y), float32(p.radiusPx), p.color, false)
}
