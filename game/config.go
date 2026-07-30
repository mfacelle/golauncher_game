package game

import "image/color"

// top-level definitions used across classes/packages.
// trying to keep them all stored here, for ease of finding them (for now)
const (
	WindowWidth             = 800
	WindowHeight            = 800
	LevelClearedDurationS   = 1
	GameOverScreenDurationS = 3
	VictoryScreenDurationS  = 3
	InitParticleCount       = 100
	InitScore               = 0
	UiBarHeightPx           = 50
)

// commonly used colors, in RGB, for use in various menus (also some unused, for reference)
// making these "const global" vars. Not really const or all that safe, but it'll work for now
var ColorOffWhite = color.RGBA{R: 255, G: 255, B: 229, A: 255}
var ColorBrown = color.RGBA{R: 90, G: 81, B: 63, A: 255}
var ColorLightBorder = color.RGBA{R: 77, G: 77, B: 69, A: 255}
var ColorBorder = color.RGBA{R: 50, G: 50, B: 45, A: 255}
var ColorDarkBorder = color.RGBA{R: 26, G: 26, B: 24, A: 255}
