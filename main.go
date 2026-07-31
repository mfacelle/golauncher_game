package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"

	"staticsheep/game"
)

func main() {
	log.Println("start")
	gameObj := game.NewGame()

	ebiten.SetWindowSize(game.WindowWidth, game.WindowHeight)

	err := ebiten.RunGame(gameObj)

	if err != nil {
		panic(err)
	}
}
