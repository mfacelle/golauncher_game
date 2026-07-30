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
	stateLevelCleared
	stateGameOver
	stateVictory
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
	totalNumLevels    int
	durationTickCount int
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
		totalNumLevels:    CountTotalNumLevels("levels/level*.json"),
		durationTickCount: 0,
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

func (g *Game) nextLevel() {
	nextLevel := g.loadNextLevel()
	if nextLevel == nil {
		log.Println("Failed to load next level, returning to main menu")
		g.mainMenu()
	} else {
		g.currentLevel = nextLevel
		g.currentState = statePlaying
	}
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

func (g *Game) levelCleared() {
	// state update logic will handle moving to next level
	g.currentState = stateLevelCleared
	g.durationTickCount = 0
}

func (g *Game) gameOver() {
	// state update logic will handle returning to main menu
	g.currentState = stateGameOver
	g.durationTickCount = 0
}

func (g *Game) victory() {
	// state update logic will handle returning to main menu
	g.currentState = stateVictory
	g.durationTickCount = 0
}

// eventually break up main menu and pause menu into separate classes
func (g *Game) Update() error {

	// additional states to add:
	// level clear (show brief screen, load next level)
	// game over (show brief screen, return to main menu)
	// victory screen (all levels cleared, show brief screen, return to main menu)

	// this should probably be broken up into calling Update on each object, based on state.
	// for now, this works, though.  but maybe split them up into separate functions?
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
				// check if this was the last level, and load victory screen if it was.
				// otherwise, load next level
				g.info.LevelsCleared++
				if g.currentLevelValue >= g.totalNumLevels {
					g.victory()
				} else {
					g.levelCleared()
				}
			} else if g.currentLevel.GameOver {
				// show some kind of game over screen
				log.Println("Level failed!")
				g.gameOver()
			}
		}
	case statePaused:
		// unpause the level if player presses pause button again
		if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
			g.resume()
		}
	case stateLevelCleared:
		// consider also updating level, to keep any particles moving
		// update count of ticks and move to next level if tick count reached
		g.durationTickCount++
		if g.durationTickCount > LevelClearedDurationS*ebiten.TPS() {
			g.nextLevel()
		}
	case stateGameOver:
		// consider also updating level, to keep any particles moving
		// update count of ticks and go back to main menu if tick count reached
		g.durationTickCount++
		if g.durationTickCount > GameOverScreenDurationS*ebiten.TPS() {
			g.mainMenu()
		}
	case stateVictory:
		// consider also updating level, to keep any particles moving
		// update count of ticks and go back to main menu if tick count reached
		g.durationTickCount++
		if g.durationTickCount > VictoryScreenDurationS*ebiten.TPS() {
			g.mainMenu()
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
	screen.Fill(color.RGBA{53, 128, 52, 255})

	// if main menu, return to prevent drawing current level
	if g.currentState == stateMainMenu {
		g.drawMainMenu(screen)
		return
	}

	// draw current level
	if g.currentLevel != nil {
		g.currentLevel.Draw(screen)
	}

	// draw states that overlay on top of the current level
	if g.currentState == statePaused {
		g.drawPauseMenu(screen)
	}

	if g.currentState == stateLevelCleared {
		g.drawLevelCleared(screen)
	} else if g.currentState == stateGameOver {
		g.drawGameOver(screen)
	} else if g.currentState == stateVictory {
		g.drawVictory(screen)
	}

	// add some kind of brief "level clear" display state
}

// keep at top-level "game" class
func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return WindowWidth, WindowHeight
}

// -----
// the following functions need to be revisited... lots of static stuff that can be done on initialization.
// should also probably move these to their own classes.  keeping here as "good enough" for now

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
	titleTextOp.ColorScale.ScaleWithColor(ColorOffWhite)
	text.Draw(screen, titleText, &text.GoTextFace{Source: assets.MainFont, Size: 48}, titleTextOp)

	subtitleText := "(work in progress)"
	subtitleTextOp := &text.DrawOptions{}
	subtitleTextOp.GeoM.Translate(screenCenterX, screenCenterY+50)
	subtitleTextOp.PrimaryAlign = text.AlignCenter
	subtitleTextOp.SecondaryAlign = text.AlignCenter
	subtitleTextOp.ColorScale.ScaleWithColor(ColorOffWhite)
	text.Draw(screen, subtitleText, &text.GoTextFace{Source: assets.MainFont, Size: 18}, subtitleTextOp)

	instructionText := "Press ENTER or click mouse to start"
	instrTextOp := &text.DrawOptions{}
	instrTextOp.GeoM.Translate(screenCenterX, screenCenterY+100)
	instrTextOp.PrimaryAlign = text.AlignCenter
	instrTextOp.SecondaryAlign = text.AlignCenter
	instrTextOp.ColorScale.ScaleWithColor(ColorOffWhite)
	text.Draw(screen, instructionText, &text.GoTextFace{Source: assets.MainFont, Size: 24}, instrTextOp)

	highScoreTextOp := &text.DrawOptions{}
	highScoreTextOp.GeoM.Translate(screenCenterX, screenCenterY+200)
	highScoreTextOp.PrimaryAlign = text.AlignCenter
	highScoreTextOp.SecondaryAlign = text.AlignCenter
	highScoreTextOp.ColorScale.ScaleWithColor(ColorOffWhite)
	text.Draw(screen, g.highScoreText, &text.GoTextFace{Source: assets.MainFont, Size: 24}, highScoreTextOp)
}

// draws the pause menu (eventually make separate class if this gets any more invovled)
// note: this code is essentially reused for level cleared, game over, and victory.  consider
// making this into a generic "draw menu" type of function or class
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
		ColorLightBorder,
		true)
	vector.StrokeRect(screen,
		float32(screenCenterX-pauseBoxWidthPx/2),
		float32(screenCenterY-pauseBoxHeightPx/2),
		float32(pauseBoxWidthPx),
		float32(pauseBoxHeightPx),
		float32(4.0),
		ColorBorder,
		true)

	pauseText := "PAUSED"
	pauseTextOp := &text.DrawOptions{}
	pauseTextOp.GeoM.Translate(screenCenterX, screenCenterY)
	pauseTextOp.PrimaryAlign = text.AlignCenter
	pauseTextOp.SecondaryAlign = text.AlignCenter
	pauseTextOp.ColorScale.ScaleWithColor(ColorOffWhite)
	text.Draw(screen, pauseText, &text.GoTextFace{Source: assets.MainFont, Size: 48}, pauseTextOp)
}

// draws the level cleared message
func (g *Game) drawLevelCleared(screen *ebiten.Image) {
	// probably would be better to set most of this up in constructor, since it's mostly static.
	// good enough here for now

	screenCenterX := float64(screen.Bounds().Dx() / 2)
	screenCenterY := float64(screen.Bounds().Dy() / 2)
	textBoxWidthPx := 400.0
	textBoxHeightPx := 150.0

	vector.FillRect(screen,
		float32(screenCenterX-textBoxWidthPx/2),
		float32(screenCenterY-textBoxHeightPx/2),
		float32(textBoxWidthPx),
		float32(textBoxHeightPx),
		ColorLightBorder,
		true)
	vector.StrokeRect(screen,
		float32(screenCenterX-textBoxWidthPx/2),
		float32(screenCenterY-textBoxHeightPx/2),
		float32(textBoxWidthPx),
		float32(textBoxHeightPx),
		float32(4.0),
		ColorBorder,
		true)

	menuText := "LEVEL CLEARED!"
	textOp := &text.DrawOptions{}
	textOp.GeoM.Translate(screenCenterX, screenCenterY)
	textOp.PrimaryAlign = text.AlignCenter
	textOp.SecondaryAlign = text.AlignCenter
	textOp.ColorScale.ScaleWithColor(ColorOffWhite)
	text.Draw(screen, menuText, &text.GoTextFace{Source: assets.MainFont, Size: 48}, textOp)
}

// draws the game over message
func (g *Game) drawGameOver(screen *ebiten.Image) {
	// probably would be better to set most of this up in constructor, since it's mostly static.
	// good enough here for now

	screenCenterX := float64(screen.Bounds().Dx() / 2)
	screenCenterY := float64(screen.Bounds().Dy() / 2)
	textBoxWidthPx := 400.0
	textBoxHeightPx := 400.0

	vector.FillRect(screen,
		float32(screenCenterX-textBoxWidthPx/2),
		float32(screenCenterY-textBoxHeightPx/2),
		float32(textBoxWidthPx),
		float32(textBoxHeightPx),
		ColorLightBorder,
		true)
	vector.StrokeRect(screen,
		float32(screenCenterX-textBoxWidthPx/2),
		float32(screenCenterY-textBoxHeightPx/2),
		float32(textBoxWidthPx),
		float32(textBoxHeightPx),
		float32(4.0),
		ColorBorder,
		true)

	menuText := "GAME OVER"
	textOp := &text.DrawOptions{}
	textOp.GeoM.Translate(screenCenterX, screenCenterY)
	textOp.PrimaryAlign = text.AlignCenter
	textOp.SecondaryAlign = text.AlignCenter
	textOp.ColorScale.ScaleWithColor(ColorOffWhite)
	text.Draw(screen, menuText, &text.GoTextFace{Source: assets.MainFont, Size: 48}, textOp)

	currScoreText := fmt.Sprintf("Score: %d", g.info.CalculateScore())
	currScoreTextOp := &text.DrawOptions{}
	currScoreTextOp.GeoM.Translate(screenCenterX, screenCenterY+50)
	currScoreTextOp.PrimaryAlign = text.AlignCenter
	currScoreTextOp.SecondaryAlign = text.AlignCenter
	currScoreTextOp.ColorScale.ScaleWithColor(ColorOffWhite)
	text.Draw(screen, currScoreText, &text.GoTextFace{Source: assets.MainFont, Size: 24}, currScoreTextOp)

	highScoreTextOp := &text.DrawOptions{}
	highScoreTextOp.GeoM.Translate(screenCenterX, screenCenterY+100)
	highScoreTextOp.PrimaryAlign = text.AlignCenter
	highScoreTextOp.SecondaryAlign = text.AlignCenter
	highScoreTextOp.ColorScale.ScaleWithColor(ColorOffWhite)
	text.Draw(screen, g.highScoreText, &text.GoTextFace{Source: assets.MainFont, Size: 24}, highScoreTextOp)
}

// draws the victory message
func (g *Game) drawVictory(screen *ebiten.Image) {
	// probably would be better to set most of this up in constructor, since it's mostly static.
	// good enough here for now

	screenCenterX := float64(screen.Bounds().Dx() / 2)
	screenCenterY := float64(screen.Bounds().Dy() / 2)
	textBoxWidthPx := 400.0
	textBoxHeightPx := 400.0

	vector.FillRect(screen,
		float32(screenCenterX-textBoxWidthPx/2),
		float32(screenCenterY-textBoxHeightPx/2),
		float32(textBoxWidthPx),
		float32(textBoxHeightPx),
		ColorLightBorder,
		true)
	vector.StrokeRect(screen,
		float32(screenCenterX-textBoxWidthPx/2),
		float32(screenCenterY-textBoxHeightPx/2),
		float32(textBoxWidthPx),
		float32(textBoxHeightPx),
		float32(4.0),
		ColorBorder,
		true)

	menuText := "VICTORY!"
	textOp := &text.DrawOptions{}
	textOp.GeoM.Translate(screenCenterX, screenCenterY)
	textOp.PrimaryAlign = text.AlignCenter
	textOp.SecondaryAlign = text.AlignCenter
	textOp.ColorScale.ScaleWithColor(ColorOffWhite)
	text.Draw(screen, menuText, &text.GoTextFace{Source: assets.MainFont, Size: 48}, textOp)

	clearText := "All levels cleared"
	clearTextOp := &text.DrawOptions{}
	clearTextOp.GeoM.Translate(screenCenterX, screenCenterY+50)
	clearTextOp.PrimaryAlign = text.AlignCenter
	clearTextOp.SecondaryAlign = text.AlignCenter
	clearTextOp.ColorScale.ScaleWithColor(ColorOffWhite)
	text.Draw(screen, clearText, &text.GoTextFace{Source: assets.MainFont, Size: 24}, clearTextOp)

	currScoreText := fmt.Sprintf("Score: %d", g.info.CalculateScore())
	currScoreTextOp := &text.DrawOptions{}
	currScoreTextOp.GeoM.Translate(screenCenterX, screenCenterY+100)
	currScoreTextOp.PrimaryAlign = text.AlignCenter
	currScoreTextOp.SecondaryAlign = text.AlignCenter
	currScoreTextOp.ColorScale.ScaleWithColor(ColorOffWhite)
	text.Draw(screen, currScoreText, &text.GoTextFace{Source: assets.MainFont, Size: 24}, currScoreTextOp)
}
