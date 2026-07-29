package game

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// represents a "goal" object that the player is trying to launch particles to
// essentially a charged object, but may not have any charge
type GoalObject struct {
	ChargedObject
}

func (o *GoalObject) Update() error {
	// nothing to update.  Goal bounds-checking will be done in level class
	return nil
}

func (o *GoalObject) Draw(screen *ebiten.Image) {
	// redundant becuase this object should have been added to list of ChargedObjects in the level.
	// ...but this class doesn't know that, so this is kind of weird.
	// vector.FillCircle(screen, float32(o.positionPx.X), float32(o.positionPx.Y), float32(o.radius), o.color, true)
	// vector.StrokeCircle(screen, float32(o.positionPx.X), float32(o.positionPx.Y), float32(o.radius), float32(o.borderWidthPx), o.borderColor, true)
}
