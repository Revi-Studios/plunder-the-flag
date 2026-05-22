package game

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	collider "github.com/melonfunction/ebiten-collider"
	"github.com/revi-studios/plunder-the-flag/lib"
)

type Player struct {
	Name string

	X float64
	Y float64

	worldData      *lib.WorldData
	collisionShape *collider.RectangleShape

	Jump_force float64
	speed      float64

	flag *Flag

	Sprite   *ebiten.Image
	animator *lib.Animator

	xv float64
	yv float64

	tick               int
	ticks_before_reset int

	last_jumped      int
	last_flag_picked int
	facing_left      bool
}

func (Player) New(worldData *lib.WorldData, x, y float64) *Player {
	img, _, err := ebitenutil.NewImageFromFile("assets/images/pirate-pink.png")
	if err != nil {
		log.Fatal(err)
	}

	log.Print("Player sprite size:", img.Bounds().Dx(), img.Bounds().Dy())

	p := &Player{
		Name:       "Pirate in Pink",
		Sprite:     ebiten.NewImage(50, 67),
		worldData:  worldData,
		Jump_force: 800,
		speed:      170,
		X:          x,
		Y:          y,
	}

	p.animator = lib.Animator{}.New(p.Sprite, map[string]lib.Animation{
		"idle": {
			Speed:         30,
			Source_sprite: img,
			Length:        5,
			Width:         50,
			Height:        67,
		},
	})

	p.Sprite = p.animator.Start("idle")
	p.collisionShape = p.worldData.Hash.NewRectangleShape(p.X, p.Y, 16, 32)
	p.collisionShape.SetParent(p)

	return p
}

func (self *Player) Update(delta float64) error {
	self.last_flag_picked++

	self.collisionShape.MoveTo(self.X, self.Y)

	switch true {
	case ebiten.IsKeyPressed(ebiten.KeyE) && self.last_flag_picked > 30:
		if self.flag != nil {
			self.flag.Drop(self.X, self.Y)
			log.Println("Dropped flag:", self.flag, "on team", self.flag.Team)
			self.flag = nil

			self.last_flag_picked = 0
		} else if collisions := self.worldData.Hash.CheckCollisions(self.collisionShape); len(collisions) > 0 {
			for _, collision := range collisions {
				if flag, ok := collision.Other.GetParent().(*Flag); ok && self.flag == nil {
					self.flag = flag.Pick()
					log.Println("Picked flag:", flag, "on team", flag.Team)

					self.last_flag_picked = 0
				}
			}
		}

	case ebiten.IsKeyPressed(ebiten.KeySpace) && self.Y >= 100:
		self.yv -= self.Jump_force * delta
		self.Y--

	case ebiten.IsKeyPressed(ebiten.KeyA):
		self.facing_left = true
		self.xv = -self.speed * delta

	case ebiten.IsKeyPressed(ebiten.KeyD):
		self.facing_left = false
		self.xv = self.speed * delta

		self.animator.ForceDraw()
	default:
		self.xv = 0

	}

	if ebiten.IsKeyPressed(ebiten.KeyR) {
		self.X = 0
		self.Y = 0
	}

	if !(self.Y >= 100) && !(self.Y+self.yv >= 100) {
		self.yv += self.worldData.Gravity * delta
	} else {
		self.yv = 0
		self.Y = 100
	}

	self.X += self.xv
	self.Y += self.yv

	return nil
}

func (self *Player) Draw(screen *ebiten.Image) {
	self.animator.UpdateAndDraw()

	op := &ebiten.DrawImageOptions{}

	if self.facing_left {
		op.GeoM.Scale(-1, 1)
		op.GeoM.Translate(float64(self.Sprite.Bounds().Dx()), 0)
	}

	op.GeoM.Translate(self.X-float64(self.Sprite.Bounds().Dx())/2, self.Y-float64(self.Sprite.Bounds().Dy()))

	screen.DrawImage(self.Sprite, op)
}

func (self *Player) Pos() (x, y float64) {
	return self.X, self.Y
}
