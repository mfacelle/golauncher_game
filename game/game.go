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

type Game struct {
	player    *Player
	objects   []ChargedObject
	particles []*BasicParticle
}

// need to eventually udpate this to either handle multiple types of objects,
// or split into multiple json types
type objectJSON struct {
	PositionPx Vector     `json:"positionPx"`
	Radius     float64    `json:"radius"`
	Color      color.RGBA `json:"color"`
	Charge     float64    `json:"charge"`
}

func NewGame() *Game {
	log.Println("new game")

	game := &Game{}

	// look into constructing this stuff better
	game.player = NewPlayer()
	game.objects = loadObjects("objects.json")

	return game
}

// want this to load all types of objects defining a level, not just ChargedObject.
func loadObjects(path string) []ChargedObject {
	file, err := os.ReadFile(path)
	if err != nil {
		log.Printf("failed to read objects file: %v", err)
		return nil
	}

	var objectData []objectJSON
	if err := json.Unmarshal(file, &objectData); err != nil {
		log.Printf("failed to parse objects file: %v", err)
		return nil
	}

	objects := make([]ChargedObject, 0, len(objectData))
	for _, data := range objectData {
		objects = append(objects, ChargedObject{
			positionPx: data.PositionPx,
			radius:     data.Radius,
			color:      data.Color,
			charge:     data.Charge,
		})
	}

	return objects
}

func (game *Game) Update() error {
	// should eventually make this some kind of state machine (main menu, load screen, actual game, etc)

	// fire projectile when the left mouse button is pressed
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		playerCenter := game.player.Center()
		projectile := NewBasicParticle(playerCenter, game.player.rotationRad)
		game.particles = append(game.particles, projectile)
	}

	// upate game objects

	// update objects
	game.player.Update()
	for _, object := range game.objects {
		object.Update()
	}

	// update projectiles
	remainingProjectiles := make([]*BasicParticle, 0, len(game.particles))
	for _, projectile := range game.particles {
		collision := projectile.Update(game.objects)
		// only keep projectiles that have not collided with objects and are within the screen bounds
		if !collision && isWithinScreenBounds(projectile.positionPx) {
			remainingProjectiles = append(remainingProjectiles, projectile)
		}
		// may want some kind of special "destroy" function for projectiles that encounter collision?
	}
	game.particles = remainingProjectiles

	return nil
}

func isWithinScreenBounds(position Vector) bool {
	return position.X >= -50 && position.X <= screenWidth+50 && position.Y >= -50 && position.Y <= screenHeight+50
}

func (game *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0x10, 0x10, 0x20, 0xff})
	ebitenutil.DebugPrint(screen, "WASD: move, mouse: aim, left click: fire")

	game.player.Draw(screen)
	for _, object := range game.objects {
		object.Draw(screen)
	}
	for _, projectile := range game.particles {
		projectile.Draw(screen)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return outsideWidth, outsideHeight
}
