package game

import (
	"golauncher_game/assets"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	rotationRadsPerSec = math.Pi
	movementPxPerSec   = 75
	movementAccelPx    = 2
)

type Player struct {
	positionPx  Vector
	velocityPx  Vector
	rotationRad float64
	sprite      *ebiten.Image
}

func NewPlayer(positionPx Vector) *Player {
	return &Player{
		positionPx:  positionPx,
		rotationRad: 0,
		sprite:      assets.PlayerSprite,
	}
}

func (p *Player) Update() error {
	// apply acceleration in position and rotation when moving
	dt := 1.0 / float64(ebiten.TPS())
	maxMoveSpeedPx := movementPxPerSec * dt
	moveAccelPx := movementAccelPx * dt

	// apply position acceleration
	moveX := 0.0
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		moveX -= 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		moveX += 1
	}
	if moveX != 0 {
		p.velocityPx.X += moveX * moveAccelPx
		if p.velocityPx.X > maxMoveSpeedPx {
			p.velocityPx.X = maxMoveSpeedPx
		} else if p.velocityPx.X < -maxMoveSpeedPx {
			p.velocityPx.X = -maxMoveSpeedPx
		}
	} else {
		p.velocityPx.X = approachZero(p.velocityPx.X, moveAccelPx)
	}

	moveY := 0.0
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		moveY -= 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		moveY += 1
	}
	if moveY != 0 {
		p.velocityPx.Y += moveY * moveAccelPx
		if p.velocityPx.Y > maxMoveSpeedPx {
			p.velocityPx.Y = maxMoveSpeedPx
		} else if p.velocityPx.Y < -maxMoveSpeedPx {
			p.velocityPx.Y = -maxMoveSpeedPx
		}
	} else {
		p.velocityPx.Y = approachZero(p.velocityPx.Y, moveAccelPx)
	}

	// update position
	p.positionPx.X += p.velocityPx.X
	p.positionPx.Y += p.velocityPx.Y

	// point toward the mouse position
	// eventually, would be nice to accelerate towards this position,
	// but need a good way to prevent wobble when overshooting angle
	mouseX, mouseY := ebiten.CursorPosition()
	playerCenter := p.Center()
	dx := float64(mouseX) - playerCenter.X
	dy := float64(mouseY) - playerCenter.Y
	if dx != 0 || dy != 0 {
		p.rotationRad = math.Atan2(dy, dx)
	}

	return nil
}

// let value approach 0 from negative or positive direction, using provided step
func approachZero(value, step float64) float64 {
	switch {
	case value > 0:
		if value > step {
			return value - step
		}
		return 0
	case value < 0:
		if value < -step {
			return value + step
		}
		return 0
	default:
		return 0
	}
}

func (p *Player) Center() Vector {
	bounds := p.sprite.Bounds()
	halfWidth := float64(bounds.Dx()) / 2
	halfHeight := float64(bounds.Dy()) / 2

	return Vector{
		X: p.positionPx.X + halfWidth,
		Y: p.positionPx.Y + halfHeight,
	}
}

func (p *Player) Draw(screen *ebiten.Image) {
	bounds := p.sprite.Bounds()
	halfWidth := float64(bounds.Dx()) / 2
	halfHeight := float64(bounds.Dy()) / 2

	// translate to center of object to apply rotation, then translate back
	drawOpts := &ebiten.DrawImageOptions{}
	drawOpts.GeoM.Translate(-halfWidth, -halfHeight)
	drawOpts.GeoM.Rotate(p.rotationRad)
	drawOpts.GeoM.Translate(halfWidth, halfHeight)

	drawOpts.GeoM.Translate(p.positionPx.X, p.positionPx.Y)

	screen.DrawImage(p.sprite, drawOpts)
}
