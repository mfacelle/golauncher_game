package assets

import (
	"bytes"
	"embed"
	"image"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

//go:embed *.png
var assets embed.FS

//go:embed font.ttf
var fontData []byte

var PlayerSprite = mustLoadImage("player.png")
var SheepGoalSprite = mustLoadImage("sheep_goal.png")
var SheepRedSprite = mustLoadImage("sheep_red.png")
var SheepBlueSprite = mustLoadImage("sheep_blue.png")
var FluffWhiteSprite = mustLoadImage("fluff_white.png")
var FluffRedSprite = mustLoadImage("fluff_red.png")
var FluffBlueSprite = mustLoadImage("fluff_blue.png")
var MainFont = mustLoadFont("font.ttf")

// load an image, or error if unavailable
func mustLoadImage(name string) *ebiten.Image {
	file, err := assets.Open(name)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		panic(err)
	}

	return ebiten.NewImageFromImage(img)
}

func mustLoadFont(name string) *text.GoTextFaceSource {
	fontSource, err := text.NewGoTextFaceSource(bytes.NewReader(fontData))
	if err != nil {
		log.Fatal(err)
	}

	// Create a usable text.Face with a specific size
	return fontSource
}
