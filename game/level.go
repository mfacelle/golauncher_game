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
		projectile := NewBasicParticle(playerCenter, level.player.rotationRad)
		level.particles = append(level.particles, projectile)

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

	remainingProjectiles := make([]*BasicParticle, 0, len(level.particles))
	for _, projectile := range level.particles {
		// update projectile position and detect collisions with charged objects or obstacles
		collision := projectile.Update(level.objects)
		for _, obstacle := range level.obstacles {
			if obstacle.IntersectsProjectile(projectile) {
				collision = true
				break
			}
		}

		// check for goal collision, indicating level is cleared
		if level.goal.IntersectsProjectile(projectile) {
			collision = true
			level.IsCleared = true
			log.Println("Level cleared!")
		}

		if !collision && isWithinScreenBounds(projectile.positionPx) {
			remainingProjectiles = append(remainingProjectiles, projectile)
		}

	}
	level.particles = remainingProjectiles

	return nil
}

func (level *Level) Draw(screen *ebiten.Image) {
	screen.Fill(level.backgroundColor)

	if level.playZone != nil {
		level.playZone.Draw(screen)
	}

	ebitenutil.DebugPrint(screen, "WASD: move | mouse: aim | left click: fire | "+level.name)

	level.player.Draw(screen)

	for _, object := range level.objects {
		object.Draw(screen)
	}

	for _, obstacle := range level.obstacles {
		obstacle.Draw(screen)
	}

	for _, projectile := range level.particles {
		projectile.Draw(screen)
	}

	level.goal.Draw(screen)

	// draw remaining particle count and label
	// note: some of this stuff is static and could be set up in constructor
	particleCountText := fmt.Sprintf("%d", level.info.NumParticles)
	textFont := &text.GoTextFace{Source: assets.MainFont, Size: 48}
	_, textHeight := text.Measure(particleCountText, textFont, 0)
	textOp := &text.DrawOptions{}
	// padding with 20 pixels
	textOp.GeoM.Translate(0, 20)
	textOp.PrimaryAlign = text.AlignStart
	textOp.SecondaryAlign = text.AlignStart
	textOp.ColorScale.ScaleWithColor(color.White)
	text.Draw(screen, particleCountText, textFont, textOp)

	particleCountLabel := "Particles\nRemaining"
	labelFont := &text.GoTextFace{Source: assets.MainFont, Size: 18}
	labelOp := &text.DrawOptions{}
	labelOp.GeoM.Translate(0, textHeight+10)
	labelOp.PrimaryAlign = text.AlignStart
	labelOp.SecondaryAlign = text.AlignStart
	labelOp.LineSpacing = 12
	labelOp.ColorScale.ScaleWithColor(color.White)
	text.Draw(screen, particleCountLabel, labelFont, labelOp)
}
