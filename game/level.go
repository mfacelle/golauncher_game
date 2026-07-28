package game

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
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
	isCleared       bool
}

func (level Level) IsCleared() bool {
	return level.isCleared
}

func (level *Level) Update() error {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		playerCenter := level.player.Center()
		projectile := NewBasicParticle(playerCenter, level.player.rotationRad)
		level.particles = append(level.particles, projectile)
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
			level.isCleared = true
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

	ebitenutil.DebugPrint(screen, "WASD: move, mouse: aim, left click: fire")
	ebitenutil.DebugPrint(screen, "\n"+level.name)

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
}
