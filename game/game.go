package game

import (
	"bytes"
	"fmt"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	collider "github.com/melonfunction/ebiten-collider"
	"github.com/revi-studios/plunder-the-flag/lib"
)

func NewGame() *Game {
	title, _, err := ebitenutil.NewImageFromFile("assets/images/plunder-the-flag-title.png")
	if err != nil {
		log.Fatal(err)
	}
	fontBytes, err := os.ReadFile("assets/fonts/pirataone-regular.ttf")
	if err != nil {
		log.Fatalf("failed to read font file: %v", err)
	}
	fontSource, err := text.NewGoTextFaceSource(bytes.NewReader(fontBytes))
	if err != nil {
		log.Fatalf("failed to parse font: %v", err)
	}

	game := Game{
		WorldData: &lib.WorldData{
			Gravity: 70,
			Hash:    collider.NewSpatialHash(180),
		},
		font:      fontSource,
		title:     title,
		world:     ebiten.NewImage(100, 100),
		worldZoom: 2,
	}
	game.flag = Flag{}.New(0, game.WorldData, 20, 100)
	game.Player = Player{}.New(game.WorldData, 20, 0)

	return &game
}

type Game struct {
	Player *Player

	flag  *Flag
	title *ebiten.Image

	WorldData *lib.WorldData
	font      *text.GoTextFaceSource

	world     *ebiten.Image
	worldZoom int
}

func (g *Game) Update() error {

	g.Player.Update(float64(1) / 60)

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.world.Clear()
	// g.world.Fill(color.RGBA{R: 88, G: 127, B: 232})

	op := &ebiten.DrawImageOptions{}

	// Flag
	g.flag.Draw(screen)

	// Title
	op.GeoM.Reset()
	op.GeoM.Translate(100, 20)
	g.world.DrawImage(g.title, op)

	// Player
	g.Player.Draw(screen)

	op.GeoM.Reset()
	op.GeoM.Scale(float64(g.worldZoom), float64(g.worldZoom))
	screen.DrawImage(g.world, op)

	// Cords Text
	tOp := &text.DrawOptions{}
	tOp.GeoM.Translate(20, 30)
	tOp.Filter = ebiten.FilterLinear
	text.Draw(screen, fmt.Sprintf("x: %.1f, y: %.1f", g.Player.X, g.Player.Y), &text.GoTextFace{Source: g.font, Size: 25}, tOp)

	tOp.GeoM.Reset()
	tOp.GeoM.Translate(20, 0)
	tOp.Filter = ebiten.FilterLinear
	text.Draw(screen, fmt.Sprintf("FPS: %.2f", ebiten.ActualFPS()), &text.GoTextFace{Source: g.font, Size: 25}, tOp)

	tOp.GeoM.Reset()
	tOp.GeoM.Translate(20, 60)
	tOp.Filter = ebiten.FilterLinear
	hasFlag := g.Player.flag != nil
	text.Draw(screen, fmt.Sprintf("Has a flag: %v", hasFlag), &text.GoTextFace{Source: g.font, Size: 25}, tOp)

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	if g.world.Bounds().Dx() != outsideWidth || g.world.Bounds().Dy() != outsideHeight {
		g.world = ebiten.NewImage(outsideWidth/g.worldZoom, outsideHeight/g.worldZoom)
	}

	return outsideWidth, outsideHeight
}
