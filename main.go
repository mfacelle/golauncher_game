package main

import (
	"log"

    "github.com/hajimehoshi/ebiten/v2"

	"golauncher_game/game"
)


func main() {
	log.Println("start")
	gameObj := game.NewGame()

	err := ebiten.RunGame(gameObj)
	if err != nil {
		panic(err)
	}
}
