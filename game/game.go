package game

import (
	"bytes"
	"fmt"
	"image/color"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	collider "github.com/melonfunction/ebiten-collider"
	"github.com/revi-studios/plunder-the-flag/lib"
)

type Game struct {
	Player *Player

	flag  *Flag
	flag2 *Flag
	title *ebiten.Image

	WorldData *lib.WorldData

	world      *ebiten.Image
	worldScale float64
	fontScale  float64
	ground     *collider.RectangleShape
}

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
			Font:    fontSource,
			Debug:   false,
		},
		title:     title,
		world:     ebiten.NewImage(100, 100),
		worldZoom: 1,
	}
	game.flag = Flag{}.New(0, game.WorldData, 20, 100)
	game.flag2 = Flag{}.New(1, game.WorldData, 80, 100)
	game.Player = Player{}.New(game.WorldData, 20, 0)
	game.ground = game.WorldData.Hash.NewRectangleShape(0, 200, 800, 200)
	game.ground.SetParent("ground")

	return &game
}

type Game struct {
	Player *Player

	flag  *Flag
	title *ebiten.Image

	WorldData *lib.WorldData

	world     *ebiten.Image
	worldZoom int
	ground    *collider.RectangleShape
}

func (g *Game) Update() error {

	g.Player.Update(float64(1) / 60)

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.world.Clear()
	// g.world.Fill(color.RGBA{R: 88, G: 127, B: 232})

	op := &ebiten.DrawImageOptions{}

	// Title
	op.GeoM.Reset()
	op.GeoM.Translate(100, 20)
	op.GeoM.Scale(2, 2)
	g.world.DrawImage(g.title, op)

	// Flags
	g.flag.Draw(g.world)
	g.flag2.Draw(g.world)

	// Player
	g.Player.Draw(g.world)

	if g.WorldData.Debug {
		g.WorldData.Hash.Draw(g.world)
		vector.StrokeRect(g.world, float32(g.ground.Pos.X-g.ground.Width/2), float32(g.ground.Pos.Y-g.ground.Height/2), float32(g.ground.Width), float32(g.ground.Height), 2.0, color.RGBA{R: 200, G: 10, B: 10, A: 255}, true)
	}

	op.GeoM.Reset()
	op.GeoM.Scale(float64(g.worldZoom), float64(g.worldZoom))
	screen.DrawImage(g.world, op)

	// Cords Text
	tOp := &text.DrawOptions{}
	tOp.GeoM.Translate(20, 30)
	tOp.Filter = ebiten.FilterLinear
	text.Draw(screen, fmt.Sprintf("x: %.1f, y: %.1f", g.Player.X, g.Player.Y), &text.GoTextFace{Source: g.WorldData.Font, Size: 25}, tOp)

	tOp.GeoM.Reset()
	tOp.GeoM.Translate(20, 0)
	tOp.Filter = ebiten.FilterLinear
	text.Draw(screen, fmt.Sprintf("FPS: %.2f", ebiten.ActualFPS()), &text.GoTextFace{Source: g.flag.worldData.Font, Size: 25}, tOp)

	tOp.GeoM.Reset()
	tOp.GeoM.Translate(20, 60)
	tOp.Filter = ebiten.FilterLinear
	hasFlag := g.Player.flag != nil
	text.Draw(screen, fmt.Sprintf("Has a flag: %v", hasFlag), &text.GoTextFace{Source: g.flag.worldData.Font, Size: 25}, tOp)

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	if g.world.Bounds().Dx() != outsideWidth || g.world.Bounds().Dy() != outsideHeight {
		g.world = ebiten.NewImage(outsideWidth/g.worldZoom, outsideHeight/g.worldZoom)
	}

	return outsideWidth, outsideHeight
}
