package game

import (
	"fmt"
	"golauncher_game/assets"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// represents a "level" which has a goal, obstacles, objects, the player zone, and particles that
// the player will launch.
// look into reorganizing objects.  Probably want arrays of objects, not pointers
type Level struct {
	player          *Player
	name            string
	goal            *GoalObject
	objects         []ChargedObject
	obstacles       []Obstacle
	playZone        *PlayerZone
	particles       []*BasicParticle
	backgroundColor color.RGBA
	IsCleared       bool
	GameOver        bool
	info            *PlayerInfo
}

func (level *Level) Update() error {

	// if no particles remaining, fail the level
	if level.info.NumParticles == 0 && len(level.particles) == 0 {
		level.GameOver = true
	}

	// eventually:
	// - include options to switch to different types of particles (and add text for it)
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) &&
		level.info.NumParticles > 0 {
		playerCenter := level.player.Center()
		particle := NewBasicParticle(playerCenter, level.player.rotationRad)
		level.particles = append(level.particles, particle)

		level.info.NumParticles--
	}

	level.player.Update()
	if level.playZone != nil {
		level.playZone.Clamp(level.player)
	}
	for _, object := range level.objects {
		object.Update()
	}
	for _, obstacle := range level.obstacles {
		obstacle.Update()
	}

	// consider parallelizing this with goroutines?
	remainingParticles := make([]*BasicParticle, 0, len(level.particles))
	for _, particle := range level.particles {
		// update particle position and detect collisions with charged objects or obstacles
		collision := particle.Update(level.objects)
		for _, obstacle := range level.obstacles {
			if obstacle.IntersectsParticle(particle) {
				collision = true
				break
			}
		}

		// check for goal collision, indicating level is cleared
		if level.goal.IntersectsParticle(particle) {
			collision = true
			level.IsCleared = true
			log.Println("Level cleared!")
		}

		if !collision && isWithinScreenBounds(particle.positionPx) {
			remainingParticles = append(remainingParticles, particle)
		}

	}
	level.particles = remainingParticles

	return nil
}

// NOTE: need to create some kind of toolbar for UI things like fluff count
func (level *Level) Draw(screen *ebiten.Image) {
	screen.Fill(level.backgroundColor)

	if level.playZone != nil {
		level.playZone.Draw(screen)
	}

	for _, object := range level.objects {
		object.Draw(screen)
	}

	for _, obstacle := range level.obstacles {
		obstacle.Draw(screen)
	}

	for _, particle := range level.particles {
		particle.Draw(screen)
	}

	level.player.Draw(screen)

	level.goal.Draw(screen)

	// draw toolbar at top for UI
	// code for this can be made better later
	// should also split this into a separate function
	vector.FillRect(screen,
		float32(0.0),
		float32(0.0),
		float32(WindowWidth),
		float32(UiBarHeightPx),
		ColorBrown,
		true)
	vector.StrokeRect(screen,
		float32(0.0),
		float32(0.0),
		float32(WindowWidth),
		float32(UiBarHeightPx),
		float32(2),
		ColorBorder,
		true)

	ebitenutil.DebugPrint(screen, "WASD: move | mouse: aim | left click: fire | "+level.name)

	// draw remaining particle count and label
	// note: some of this stuff is static and could be set up in constructor
	particleCountText := fmt.Sprintf("%d", level.info.NumParticles)
	textFont := &text.GoTextFace{Source: assets.MainFont, Size: 48}
	textWidth, _ := text.Measure(particleCountText, textFont, 0)
	textOp := &text.DrawOptions{}
	// align to bottom of UI bar
	textOp.GeoM.Translate(0, UiBarHeightPx)
	textOp.PrimaryAlign = text.AlignStart
	textOp.SecondaryAlign = text.AlignEnd
	textOp.ColorScale.ScaleWithColor(ColorOffWhite)
	text.Draw(screen, particleCountText, textFont, textOp)

	particleCountLabel := "Fluff Remaining"
	labelFont := &text.GoTextFace{Source: assets.MainFont, Size: 18}
	labelOp := &text.DrawOptions{}
	// align to bottom of UI bar
	labelOp.GeoM.Translate(textWidth+10, UiBarHeightPx)
	labelOp.PrimaryAlign = text.AlignStart
	labelOp.SecondaryAlign = text.AlignEnd
	labelOp.ColorScale.ScaleWithColor(ColorOffWhite)
	text.Draw(screen, particleCountLabel, labelFont, labelOp)
}
