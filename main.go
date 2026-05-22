package main

import (
	"bytes"
	_ "embed"
	"image"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/revi-studios/plunder-the-flag/game"
)

//go:embed assets/images/red-flag.png
var iconBytes []byte

func main() {
	img, _, err := image.Decode(bytes.NewReader(iconBytes))
	if err != nil {
		log.Fatal(err)
	}

	ebiten.SetWindowIcon([]image.Image{img})
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Plunder the Flag")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)

	if err := ebiten.RunGame(game.NewGame()); err != nil {
		log.Fatal(err)
	}
}
