package game

import (
	"fmt"
	"image/color"
	"log"
	"strconv"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
	collider "github.com/melonfunction/ebiten-collider"
	"github.com/revi-studios/plunder-the-flag/lib"
)

type Flag struct {
	Team int
	X    float64
	Y    float64

	IsPicked bool

	worldData      *lib.WorldData
	collisionShape *collider.RectangleShape

	sprite   *ebiten.Image
	animator *lib.Animator
}

func (Flag) New(team int, world *lib.WorldData, x, y float64) *Flag {
	img, _, err := ebitenutil.NewImageFromFile("assets/images/red-flag.png")
	if err != nil {
		log.Fatal(err)
	}

	var start int
	switch team {
	case 1:
		start = 48
	default:
		start = 0
	}

	flag := &Flag{
		worldData: world,
		Team:      team,
		IsPicked:  false,
	}
	flag.animator = lib.Animator{}.New(map[string]lib.Animation{
		"default": {
			Start:         start,
			Speed:         30,
			Source_sprite: img,
			Length:        3,
			Width:         16,
			Height:        22,
		},
	})
	flag.sprite = flag.animator.InitImage()

	flag.animator.Play("default")

	flag.Drop(x, y)

	return flag
}

func (self *Flag) Draw(screen *ebiten.Image) {
	self.animator.UpdateAndDraw()
	if !self.IsPicked {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(self.X-float64(self.sprite.Bounds().Dx())/2, self.Y-float64(self.sprite.Bounds().Dy()))
		screen.DrawImage(self.sprite, op)
	}
	if self.worldData.Debug {
		if !self.IsPicked {
			ebitenutil.DebugPrintAt(screen, fmt.Sprintf("[X: %.0f, Y: %.0f]", self.X, self.Y), int(self.X), int(self.Y))
			vector.StrokeRect(screen, float32(self.collisionShape.Pos.X-self.collisionShape.Width/2), float32(self.collisionShape.Pos.Y-self.collisionShape.Height/2), float32(self.collisionShape.Width), float32(self.collisionShape.Height), 2.0, color.RGBA{R: 200, G: 10, B: 10, A: 255}, true)
		}
	}

}

func (self *Flag) Pick() *Flag {
	log.Println("Flag Picked!" + " (Team: " + strconv.Itoa(self.Team) + ")")
	self.IsPicked = true
	self.worldData.Hash.Remove(self.collisionShape)
	return self
}

func (self *Flag) Drop(x, y float64) {
	self.X = x
	self.Y = y

	self.collisionShape = self.worldData.Hash.NewRectangleShape(self.X, 8, 16, 32)
	self.collisionShape.SetParent(self)

	var found bool
FindGround:
	for range 400 {
		self.collisionShape.Move(0, 1)

		for _, collision := range self.worldData.Hash.CheckCollisions(self.collisionShape) {
			if collision.Other.GetParent() == "ground" {
				found = true
				self.collisionShape.Move(collision.SeparatingVector.X, collision.SeparatingVector.Y)

				self.Y = self.collisionShape.Pos.Y + self.collisionShape.Height/2
				break FindGround
			}
		}
	}

	if !found {
		self.collisionShape.MoveTo(self.X, self.Y-self.collisionShape.Height/2)
	}

	log.Println("Flag Dropped!" + " (Team: " + strconv.Itoa(self.Team) + ")")
	self.IsPicked = false
}
