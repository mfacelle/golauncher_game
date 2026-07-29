package game

import (
	"fmt"
	"image/color"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"golauncher_game/assets"
)

type gameState int

const (
	stateMainMenu gameState = iota
	statePlaying
	statePaused
)

// get gameState string, for debugging
func (gs gameState) String() string {
	// Boundary check to prevent panic
	if gs < stateMainMenu || gs > statePaused {
		return fmt.Sprintf("Invalid(%d)", gs)
	}

	// Fast lookup using an array
	return [...]string{"MainMenu", "Playing", "Paused"}[gs]
}

// top-level game class, managing all states of the game (menus and levels)
type Game struct {
	currentState      gameState
	currentLevel      *Level
	currentLevelValue int
	highScoreText     string
	info              PlayerInfo
}

// create a new game object
func NewGame() *Game {
	log.Println("new game")

	// start at level 1.  eventually, add some kind of level select menu
	g := &Game{
		currentState:      stateMainMenu,
		currentLevel:      nil,
		currentLevelValue: 0,
		highScoreText:     "High Score: 0",
		info: PlayerInfo{
			NumParticles:  InitParticleCount,
			LevelsCleared: InitScore,
		},
	}

	return g
}

func (g *Game) loadNextLevel() *Level {
	// check for the "debug" level.  if found, always load that one
	_, err := os.Stat("levels/level0.json")
	if err == nil {
		g.currentLevelValue = 0
	} else {
		g.currentLevelValue++
	}

	levelFileName := fmt.Sprintf("levels/level%d.json", g.currentLevelValue)
	log.Printf("Loading %s\n", levelFileName)
	return loadLevelFromJSON(levelFileName, &g.info)
}

func (g *Game) start() {
	log.Println("starting game")
	g.currentState = statePlaying
	g.currentLevel = g.loadNextLevel()
}

func (g *Game) mainMenu() {
	log.Println("back to main menu")
	g.highScoreText = fmt.Sprintf("High Score: %d", g.info.CalculateScore())
	g.currentState = stateMainMenu
	g.currentLevelValue = 0
	// reset player info
	g.info = PlayerInfo{
		NumParticles:  InitParticleCount,
		LevelsCleared: InitScore,
	}
}

func (g *Game) pause() {
	if g.currentState == statePlaying {
		g.currentState = statePaused
	}
}

func (g *Game) resume() {
	if g.currentState == statePaused {
		g.currentState = statePlaying
	}
}

// eventually break up main menu and pause menu into separate classes
func (g *Game) Update() error {

	// additional states to add:
	// level clear (show brief screen, load next level)
	// game over (show brief screen, return to main menu)
	// victory screen (all levels cleared, show brief screen, return to main menu)

	// this should probably be broken up into calling Update on each object, based on state.
	// for now, this works, though
	switch g.currentState {
	case stateMainMenu:
		// start the level if player presses start button or clicks mouse
		if inpututil.IsKeyJustPressed(ebiten.KeyEnter) ||
			inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			g.start()
		}
	case statePlaying:
		// pause the level if player presses pause, otherwise, run the level
		if g.currentLevel == nil {
			log.Println("Current level is nil.  Returning to main menu")
			g.mainMenu()
		} else if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
			g.pause()
		} else {
			err := g.currentLevel.Update()
			if err != nil {
				return err
			}

			if g.currentLevel.IsCleared {
				// load next level - may want to make this a separate game state,
				// which will allow displaying some kind of "level cleared!" message
				// probably not great to set currentLevel right here
				g.info.LevelsCleared++
				nextLevel := g.loadNextLevel()
				if nextLevel == nil {
					log.Println("Failed to load next level, returning to main menu")
					g.mainMenu()
				} else {
					g.currentLevel = nextLevel
				}
			} else if g.currentLevel.GameOver {
				// show some kind of game over screen
				log.Println("Level failed!")
				g.mainMenu()
			}
		}
	case statePaused:
		// unpause the level if player presses pause button again
		if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
			g.resume()
		}
	}

	return nil
}

func isWithinScreenBounds(position Vector) bool {
	// buffer of 5 px around screen edges
	bufferPx := 5.0
	return position.X >= -bufferPx && position.X <= WindowWidth+bufferPx && position.Y >= -bufferPx && position.Y <= WindowHeight+bufferPx
}

// eventually break up main menu and pause menu into separate classes (or at least seperate functions)
func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0x20, 0x20, 0x20, 0xff})

	if g.currentState == stateMainMenu {
		g.drawMainMenu(screen)
		return
	}

	if g.currentLevel != nil {
		g.currentLevel.Draw(screen)
	}

	if g.currentState == statePaused {
		g.drawPauseMenu(screen)
	}

	// add some kind of brief "level clear" display state
}

// keep at top-level "game" class
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return WindowWidth, WindowHeight
}

// draws the main menu (eventually make separate class if this gets any more invovled)
func (g *Game) drawMainMenu(screen *ebiten.Image) {
	// good enough for now
	ebitenutil.DebugPrint(screen, "Main Menu\nPress Enter or click mouse to start")

	// probably would be better to set most of this up in constructor, since it's mostly static.
	// good enough here for now

	screenCenterX := float64(screen.Bounds().Dx() / 2)
	screenCenterY := float64(screen.Bounds().Dy() / 2)
	titleText := "Static Sheep"
	titleTextOp := &text.DrawOptions{}
	titleTextOp.GeoM.Translate(screenCenterX, screenCenterY)
	titleTextOp.PrimaryAlign = text.AlignCenter
	titleTextOp.SecondaryAlign = text.AlignCenter
	titleTextOp.ColorScale.ScaleWithColor(color.White)
	text.Draw(screen, titleText, &text.GoTextFace{Source: assets.MainFont, Size: 48}, titleTextOp)

	subtitleText := "(work in progress, no sheep yet)"
	subtitleTextOp := &text.DrawOptions{}
	subtitleTextOp.GeoM.Translate(screenCenterX, screenCenterY+50)
	subtitleTextOp.PrimaryAlign = text.AlignCenter
	subtitleTextOp.SecondaryAlign = text.AlignCenter
	subtitleTextOp.ColorScale.ScaleWithColor(color.White)
	text.Draw(screen, subtitleText, &text.GoTextFace{Source: assets.MainFont, Size: 18}, subtitleTextOp)

	instructionText := "Press ENTER to start"
	instrTextOp := &text.DrawOptions{}
	instrTextOp.GeoM.Translate(screenCenterX, screenCenterY+100)
	instrTextOp.PrimaryAlign = text.AlignCenter
	instrTextOp.SecondaryAlign = text.AlignCenter
	instrTextOp.ColorScale.ScaleWithColor(color.White)
	text.Draw(screen, instructionText, &text.GoTextFace{Source: assets.MainFont, Size: 24}, instrTextOp)

	highScoreTextOp := &text.DrawOptions{}
	highScoreTextOp.GeoM.Translate(screenCenterX, screenCenterY+200)
	highScoreTextOp.PrimaryAlign = text.AlignCenter
	highScoreTextOp.SecondaryAlign = text.AlignCenter
	highScoreTextOp.ColorScale.ScaleWithColor(color.White)
	text.Draw(screen, g.highScoreText, &text.GoTextFace{Source: assets.MainFont, Size: 24}, highScoreTextOp)
}

// draws the pause menu (eventually make separate class if this gets any more invovled)
func (g *Game) drawPauseMenu(screen *ebiten.Image) {
	// move to be actual text in the middle of the screen (with a background box probably)
	ebitenutil.DebugPrint(screen, "PAUSED\nPress Escape to resume")

	// probably would be better to set most of this up in constructor, since it's mostly static.
	// good enough here for now

	screenCenterX := float64(screen.Bounds().Dx() / 2)
	screenCenterY := float64(screen.Bounds().Dy() / 2)
	pauseBoxWidthPx := 300.0
	pauseBoxHeightPx := 200.0

	vector.FillRect(screen,
		float32(screenCenterX-pauseBoxWidthPx/2),
		float32(screenCenterY-pauseBoxHeightPx/2),
		float32(pauseBoxWidthPx),
		float32(pauseBoxHeightPx),
		color.RGBA{R: 40, G: 40, B: 80, A: 255},
		true)
	vector.StrokeRect(screen,
		float32(screenCenterX-pauseBoxWidthPx/2),
		float32(screenCenterY-pauseBoxHeightPx/2),
		float32(pauseBoxWidthPx),
		float32(pauseBoxHeightPx),
		float32(4.0),
		color.RGBA{R: 200, G: 200, B: 200, A: 255},
		true)

	pauseText := "PAUSED"
	pauseTextOp := &text.DrawOptions{}
	pauseTextOp.GeoM.Translate(screenCenterX, screenCenterY)
	pauseTextOp.PrimaryAlign = text.AlignCenter
	pauseTextOp.SecondaryAlign = text.AlignCenter
	pauseTextOp.ColorScale.ScaleWithColor(color.White)
	text.Draw(screen, pauseText, &text.GoTextFace{Source: assets.MainFont, Size: 48}, pauseTextOp)
}
