package game

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
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

func (self *Flag) Pick() *Flag {
	log.Println("Flag Picked!")
	self.IsPicked = true
	self.worldData.Hash.Remove(self.collisionShape)
	return self
}

func (self *Flag) Drop(x, y float64) {
	log.Println("Flag Dropped!")
	self.IsPicked = false
	self.X = x
	self.Y = y

	self.collisionShape = self.worldData.Hash.NewRectangleShape(self.X, self.Y, 16, 32)
	self.collisionShape.SetParent(self)
}

func (self Flag) New(team int, world *lib.WorldData, x, y float64) *Flag {
	img, _, err := ebitenutil.NewImageFromFile("assets/images/red-flag.png")
	if err != nil {
		log.Fatal(err)
	}

	flag := &Flag{
		worldData: world,
		sprite:    ebiten.NewImage(16, 32),
		Team:      team,
		X:         x,
		Y:         y,
		IsPicked:  false,
	}
	flag.animator = lib.Animator{}.New(flag.sprite, map[string]lib.Animation{
		"default": {
			Speed:         30,
			Source_sprite: img,
			Length:        3,
			Width:         16,
			Height:        32,
		},
	})
	flag.sprite = flag.animator.Start("default")
	flag.collisionShape = flag.worldData.Hash.NewRectangleShape(flag.X, flag.Y, 16, 32)
	flag.collisionShape.SetParent(flag)
	return flag
}

func (self *Flag) Draw(screen *ebiten.Image) {
	self.animator.UpdateAndDraw()
	if !self.IsPicked {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(self.X-float64(self.sprite.Bounds().Dx())/2, self.Y-float64(self.sprite.Bounds().Dy()))
		screen.DrawImage(self.sprite, op)
	}
}
