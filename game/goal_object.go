package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// represents a "goal" object that the player is trying to launch particles to
// essentially a charged object, but may not have any charge
type GoalObject struct {
	ChargedObject
	borderColor   color.RGBA
	borderWidthPx float64
}

func (o *GoalObject) Update() error {
	// nothing to update.  Goal bounds-checking will be done in level class
	return nil
}

func (o *GoalObject) Draw(screen *ebiten.Image) {
	// first FillCircle is likely redundant, because this object gets added to objects array
	vector.FillCircle(screen, float32(o.positionPx.X), float32(o.positionPx.Y), float32(o.radius), o.color, true)
	vector.StrokeCircle(screen, float32(o.positionPx.X), float32(o.positionPx.Y), float32(o.radius), float32(o.borderWidthPx), o.borderColor, true)
}
