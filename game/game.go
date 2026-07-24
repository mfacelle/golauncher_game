package game

import (
	"encoding/json"
	"image/color"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	screenWidth  = 600
	screenHeight = 480
)

// need to include state machine stuff and different game states
type Game struct {
	player    *Player
	objects   []ChargedObject
	obstacles []Obstacle
	particles []*BasicParticle
}

// move to a sepatate "level loader" class
// need to eventually udpate this to either handle multiple types of objects,
// or split into multiple json types
type objectJSON struct {
	Type       string     `json:"type"`
	PositionPx Vector     `json:"positionPx"`
	Radius     float64    `json:"radius"`
	WidthPx    float64    `json:"widthPx"`
	HeightPx   float64    `json:"heightPx"`
	Color      color.RGBA `json:"color"`
	Charge     float64    `json:"charge"`
}

// keep at top-level "game" class
// but update to initialize state machine and different objects
func NewGame() *Game {
	log.Println("new game")

	game := &Game{}

	// look into constructing this stuff better
	game.player = NewPlayer()
	game.objects, game.obstacles = loadLevelObjects("objects.json")

	return game
}

// move to a sepatate "level loader" class
func loadLevelObjects(path string) ([]ChargedObject, []Obstacle) {
	file, err := os.ReadFile(path)
	if err != nil {
		log.Printf("failed to read objects file: %v", err)
		return nil, nil
	}

	var objectData []objectJSON
	if err := json.Unmarshal(file, &objectData); err != nil {
		log.Printf("failed to parse objects file: %v", err)
		return nil, nil
	}

	objects := make([]ChargedObject, 0, len(objectData))
	obstacles := make([]Obstacle, 0, len(objectData))
	for _, data := range objectData {
		switch data.Type {
		case "obstacle":
			obstacles = append(obstacles, Obstacle{
				positionPx: data.PositionPx,
				widthPx:    data.WidthPx,
				heightPx:   data.HeightPx,
				color:      data.Color,
			})
		default:
			objects = append(objects, ChargedObject{
				positionPx: data.PositionPx,
				radius:     data.Radius,
				color:      data.Color,
				charge:     data.Charge,
			})
		}
	}

	return objects, obstacles
}

// move to a specific "level" class.  top-level "game" update should be a state machine to display
// different screens (i.e. pause, main menu, level, etc)
func (game *Game) Update() error {

	// fire projectile when the left mouse button is pressed
	// eventually look into how this is allocated.  pointer vs array (heap/stack?)
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		playerCenter := game.player.Center()
		projectile := NewBasicParticle(playerCenter, game.player.rotationRad)
		game.particles = append(game.particles, projectile)
	}

	// update objects
	game.player.Update()
	for _, object := range game.objects {
		object.Update()
	}
	for _, obstacle := range game.obstacles {
		obstacle.Update()
	}

	// update projectiles, and remove any that have a collision or leave the screen
	remainingProjectiles := make([]*BasicParticle, 0, len(game.particles))
	for _, projectile := range game.particles {
		// update projectile against all charged objects
		collision := projectile.Update(game.objects)

		// check projectile for collision with any obstacles
		for _, obstacle := range game.obstacles {
			if obstacle.IntersectsProjectile(projectile) {
				collision = true
				break
			}
		}

		// only keep projectiles that have not collided with objects and are within the screen bounds
		if !collision && isWithinScreenBounds(projectile.positionPx) {
			remainingProjectiles = append(remainingProjectiles, projectile)
		}

		// may want some kind of special "destroy" function for projectiles that encounter collision?
	}
	game.particles = remainingProjectiles

	return nil
}

// keep at top-level "game" class
func isWithinScreenBounds(position Vector) bool {
	return position.X >= -50 && position.X <= screenWidth+50 && position.Y >= -50 && position.Y <= screenHeight+50
}

// move to a specific "level" class.  top-level "game" update should be a state machine to display
// different screens (i.e. pause, main menu, level, etc)func (game *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0x10, 0x10, 0x20, 0xff})
	ebitenutil.DebugPrint(screen, "WASD: move, mouse: aim, left click: fire")

	game.player.Draw(screen)
	for _, object := range game.objects {
		object.Draw(screen)
	}
	for _, obstacle := range game.obstacles {
		obstacle.Draw(screen)
	}
	for _, projectile := range game.particles {
		projectile.Draw(screen)
	}
}

// keep at top-level "game" class
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}
