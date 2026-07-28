package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type PlayerZone struct {
	positionPx    Vector
	widthPx       float64
	heightPx      float64
	color         color.RGBA
	borderColor   color.RGBA
	borderWidthPx float32
}

func NewPlayerZone(x, y, width, height float64) *PlayerZone {
	return &PlayerZone{
		positionPx:    Vector{X: x, Y: y},
		widthPx:       width,
		heightPx:      height,
		color:         color.RGBA{R: 20, G: 20, B: 40, A: 255},
		borderColor:   color.RGBA{R: 120, G: 120, B: 120, A: 255},
		borderWidthPx: 2,
	}
}

func (z *PlayerZone) Clamp(player *Player) {
	if player == nil {
		return
	}

	minX := z.positionPx.X
	maxX := z.positionPx.X + z.widthPx
	minY := z.positionPx.Y
	maxY := z.positionPx.Y + z.heightPx

	playerCenter := player.Center()

	// bound x, based on center of player
	if playerCenter.X < minX {
		player.positionPx.X = minX - float64(player.sprite.Bounds().Dx()/2)
		player.velocityPx.X = 0
	}
	if playerCenter.X > maxX {
		player.positionPx.X = maxX - float64(player.sprite.Bounds().Dx()/2)
		player.velocityPx.X = 0
	}

	// bound y, based on center of player
	if playerCenter.Y < minY {
		player.positionPx.Y = minY - float64(player.sprite.Bounds().Dy()/2)
		player.velocityPx.Y = 0
	}
	if playerCenter.Y > maxY {
		player.positionPx.Y = maxY - float64(player.sprite.Bounds().Dy()/2)
		player.velocityPx.Y = 0
	}
}

func (z *PlayerZone) Draw(screen *ebiten.Image) {
	vector.FillRect(screen, float32(z.positionPx.X), float32(z.positionPx.Y), float32(z.widthPx), float32(z.heightPx), z.color, true)
	vector.StrokeRect(screen, float32(z.positionPx.X), float32(z.positionPx.Y), float32(z.widthPx), float32(z.heightPx), z.borderWidthPx, z.borderColor, true)
}
