package game

import (
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Projectile struct {
	positionPx Vector
	velocityPx Vector
	charge     float64
	radiusPx   float64
	color      color.RGBA
}

const (
	initialVelocity       = 180.0
	electrostaticConstant = 100 // really 8.98e9
	minDistance           = 1.0
	maxDistance           = 500.0
	maxVelocity           = 250.0
)

func NewProjectile(posPx Vector, rotationRad float64) *Projectile {
	return &Projectile{
		positionPx: posPx,
		velocityPx: Vector{
			X: initialVelocity * math.Sin(rotationRad),
			Y: initialVelocity * -math.Cos(rotationRad),
		},
		charge:   -10,
		radiusPx: 2,
		color:    color.RGBA{R: 255, G: 255, B: 100, A: 255},
	}
}

func (p *Projectile) Update(objects []ChargedObject) bool {
	// this could definitely be better
	collision := false

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
		dx := object.positionPx.X - p.positionPx.X
		dy := object.positionPx.Y - p.positionPx.Y
		distanceSq := dx*dx + dy*dy
		if minDistance*minDistance < distanceSq && distanceSq < maxDistance*maxDistance {
			distance := math.Sqrt(distanceSq)
			force := math.Abs(electrostaticConstant * p.charge * object.charge / distanceSq)
			if !object.isAttract {
				force *= -1
			}
			p.velocityPx.X += force * dx / distance / float64(ebiten.TPS())
			p.velocityPx.Y += force * dy / distance / float64(ebiten.TPS())
			//log.Print("distanceSq: ", distanceSq, "force: ", force, " distance: ", distance, " dx: ", dx, " dy: ", dy)
		}
	}

	p.positionPx.X += p.velocityPx.X / float64(ebiten.TPS())
	p.positionPx.Y += p.velocityPx.Y / float64(ebiten.TPS())

	return collision
}

func (p *Projectile) Draw(screen *ebiten.Image) {
	vector.FillCircle(screen, float32(p.positionPx.X), float32(p.positionPx.Y), float32(p.radiusPx), p.color, false)
}
