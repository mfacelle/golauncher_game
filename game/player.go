package game

import (
	"golauncher_game/assets"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	defaultX           = 100
	defaultY           = 200
	rotationRadsPerSec = math.Pi
	movementPxPerSec   = 50
)

type Player struct {
	positionPx  Vector
	rotationRad float64
	sprite      *ebiten.Image
}

func NewPlayer() *Player {
	return &Player{
		positionPx:  Vector{X: defaultX, Y: defaultY},
		rotationRad: 0,
		sprite:      assets.PlayerSprite,
	}
}

func (p *Player) Update() error {
	// temp value for now (eventually make pixels per second and split rotation)
	moveSpeedPx := movementPxPerSec / float64(ebiten.TPS())
	rotationRads := rotationRadsPerSec / float64(ebiten.TPS())

	// handle movement with WASD
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		p.positionPx.X -= moveSpeedPx
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		p.positionPx.X += moveSpeedPx
	}
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		p.positionPx.Y -= moveSpeedPx
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		p.positionPx.Y += moveSpeedPx
	}

	// handle rotation with arrows keys (left/right only)
	if ebiten.IsKeyPressed(ebiten.KeyLeft) {
		p.rotationRad -= rotationRads
	}
	if ebiten.IsKeyPressed(ebiten.KeyRight) {
		p.rotationRad += rotationRads
	}

	return nil
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
