package game

import (
	"bytes"
	"fmt"
	"image/color"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	collider "github.com/melonfunction/ebiten-collider"
	"github.com/revi-studios/plunder-the-flag/lib"
)

type Game struct {
	Player   *Player
	entities *[]lib.Entity

	title *ebiten.Image

	WorldData *lib.WorldData
	Font      *text.GoTextFaceSource

	world  *ebiten.Image
	ground *collider.RectangleShape

	worldScale float64
	textScale  float64

	drawDebugMenu bool
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
		Font: fontSource,
		WorldData: &lib.WorldData{
			Gravity: 70,
			Hash:    collider.NewSpatialHash(180),
			Debug:   false,
		},
		title:    title,
		entities: &[]lib.Entity{},

		worldScale: 1.5,
		textScale:  1,
		world:      ebiten.NewImage(100, 100),
	}

	game.ground = game.WorldData.Hash.NewRectangleShape(0, 300, 712, 200)
	game.ground.SetParent("ground")
	*game.entities = append(*game.entities, Flag{}.New(0, game.WorldData, 20, 100))
	*game.entities = append(*game.entities, Flag{}.New(1, game.WorldData, 80, 100))
	*game.entities = append(*game.entities, Player{}.New(game.WorldData, game.Font, "Barron", 1, 20, 0))
	game.Player = Player{}.New(game.WorldData, game.Font, "Pirate in Pink", 0, 20, 0)

	return &game
}

func (g *Game) Update() error {
	delta := min(1/ebiten.ActualTPS(), 1.0/60)

	if inpututil.IsKeyJustPressed(ebiten.Key9) {
		g.worldScale -= 0.1
		log.Println("World scale changed to", g.worldScale)
	}
	if inpututil.IsKeyJustPressed(ebiten.Key0) {
		g.worldScale += 0.1
		log.Println("World scale changed to", g.worldScale)

	}
	if inpututil.IsKeyJustPressed(ebiten.Key3) {
		g.drawDebugMenu = !g.drawDebugMenu
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyR) {
		g.Player.X = 0
		g.Player.Y = 0
		g.Player.yv = 0
		for _, entity := range *g.entities {
			if e, ok := entity.(*Player); ok {
				e.X = 0
				e.Y = 0
				e.yv = 0
			}
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.Key1) {
		g.WorldData.Debug = !g.WorldData.Debug
	}
	g.Player.Update(delta)

	for _, entity := range *g.entities {
		entity.Update(delta)
	}

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

	// Entities
	for _, entity := range *g.entities {
		entity.Draw(g.world)
	}

	// Player
	g.Player.Draw(g.world)

	if g.WorldData.Debug {
		g.WorldData.Hash.Draw(g.world)
		vector.StrokeRect(g.world, float32(g.ground.Pos.X-g.ground.Width/2), float32(g.ground.Pos.Y-g.ground.Height/2), float32(g.ground.Width), float32(g.ground.Height), 2.0, color.RGBA{R: 200, G: 10, B: 10, A: 255}, true)
	}

	deviceScale := ebiten.Monitor().DeviceScaleFactor()
	worldScale := g.worldScale * deviceScale
	textScale := g.textScale * deviceScale

	op.GeoM.Reset()
	op.GeoM.Scale(worldScale, worldScale)
	op.Filter = ebiten.FilterNearest
	screen.DrawImage(g.world, op)

	if g.drawDebugMenu {
		DrawDebugMenu(screen, &text.GoTextFace{Source: g.Font, Size: 25}, textScale, 15, []struct {
			Name  string
			Value any
		}{
			{"Fps", ebiten.ActualFPS()},
			{"Tps", ebiten.ActualTPS()},
			{"Entities", len(*g.entities)},
			{},
			{"World Scale", g.worldScale},
			{"Text Scale", g.textScale},
			{},
			{"Font Family", g.Font.Metadata().Family},
			{"WorldImage", g.world.Bounds().Size()},
			{"Screen", screen.Bounds().Size()},
		}, []struct {
			Name  string
			Value any
		}{
			{"Name", g.Player.Name},
			{"X", g.Player.X},
			{"Y", g.Player.Y},
			{},
			{"On Ground", g.Player.onGround},
			{"Has Flag", g.Player.flag != nil},
			{"Flag", g.Player.flag},
		})
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	deviceScale := ebiten.Monitor().DeviceScaleFactor()

	highDPIWidth := int(float64(outsideWidth) * deviceScale)
	highDPIHeight := int(float64(outsideHeight) * deviceScale)

	w := int(float64(highDPIWidth) / g.worldScale)
	h := int(float64(highDPIHeight) / g.worldScale)

	if g.world.Bounds().Dx() != w || g.world.Bounds().Dy() != h {
		g.world = ebiten.NewImage(w, h)
	}

	return highDPIWidth, highDPIHeight
}

func DrawDebugMenu(screen *ebiten.Image, face text.Face, textScale float64, margin float64, left, right []struct {
	Name  string
	Value any
}) {
	var row float64
	for _, item := range left {

		var str string
		switch item.Value.(type) {
		case float64, float32:
			str = fmt.Sprintf("%s: %.1f", item.Name, item.Value)
		default:
			str = fmt.Sprintf("%s: %v", item.Name, item.Value)
		}
		if item.Name == "" {
			str = ""
		}

		op := &text.DrawOptions{}
		op.GeoM.Translate(margin, face.Metrics().HAscent*row)
		op.GeoM.Scale(textScale, textScale)
		op.Filter = ebiten.FilterLinear
		text.Draw(screen, str, face, op)

		row++
	}

	row = 0
	for _, item := range right {

		var str string
		switch item.Value.(type) {
		case float64, float32:
			str = fmt.Sprintf("%s: %.1f", item.Name, item.Value)
		default:
			str = fmt.Sprintf("%s: %v", item.Name, item.Value)
		}
		if item.Name == "" {
			str = ""
		}

		op := &text.DrawOptions{}
		op.GeoM.Translate((float64(screen.Bounds().Dx())/2)-text.Advance(str, face)-margin, face.Metrics().HAscent*row)
		op.GeoM.Scale(textScale, textScale)
		op.Filter = ebiten.FilterLinear
		text.Draw(screen, str, face, op)

		row++
	}
}
