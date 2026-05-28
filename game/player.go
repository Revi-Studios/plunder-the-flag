package game

import (
	_ "embed"
	"fmt"
	"image/color"
	_ "image/png"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
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

	last_jumped        int
	last_flag_picked   int
	last_debug_toggled int
	facing_left        bool
	onGround           bool
}

func (Player) New(worldData *lib.WorldData, x, y float64) *Player {
	idle, _, err := ebitenutil.NewImageFromFile("assets/images/pirate-pink/pirate-pink-idle.png")
	if err != nil {
		log.Fatal(err)
	}
	run, _, err := ebitenutil.NewImageFromFile("assets/images/pirate-pink/pirate-pink-run.png")
	if err != nil {
		log.Fatal(err)
	}

	p := &Player{
		Name:       "Pirate in Pink",
		worldData:  worldData,
		Jump_force: 800,
		speed:      170,
		X:          x,
		Y:          y,
	}

	p.animator = lib.Animator{}.New(map[string]lib.Animation{
		"idle": {
			Speed:         20,
			Source_sprite: idle,
			Length:        5,
			Width:         50,
			Height:        67,
		},
		"run": {
			Speed:         5,
			Source_sprite: run,
			Length:        5,
			Width:         50,
			Height:        67,
		},
	})

	p.Sprite = p.animator.InitImage()
	p.animator.Play("idle")
	p.collisionShape = p.worldData.Hash.NewRectangleShape(p.X, p.Y, 40, 50)
	p.collisionShape.SetParent(p)

	return p
}

func (self *Player) Update(delta float64) {
	self.onGround = false
	self.last_flag_picked++
	self.last_debug_toggled++

	self.collisionShape.Move(0, 1)
	for _, collision := range self.worldData.Hash.CheckCollisions(self.collisionShape) {
		if collision.Other.GetParent() == "ground" {
			pushY := collision.SeparatingVector.Y
			if pushY == 0 {
				pushY = -1
			}
			self.collisionShape.Move(0, pushY)

			self.X = self.collisionShape.Pos.X
			self.Y = self.collisionShape.Pos.Y + self.collisionShape.Height/2

			if pushY < 0 {
				self.yv = 0
				self.onGround = true
			}
		}
	}

	if !self.onGround {
		self.yv += self.worldData.Gravity * delta
	}

	switch true {
	case ebiten.IsKeyPressed(ebiten.KeyE) && self.last_flag_picked > 30:
		if self.flag != nil {
			self.flag.Drop(self.X, self.Y)
			self.flag = nil

			self.last_flag_picked = 0
		} else if collisions := self.worldData.Hash.CheckCollisions(self.collisionShape); len(collisions) > 0 {
			for _, collision := range collisions {
				if flag, ok := collision.Other.GetParent().(*Flag); ok && self.flag == nil {
					self.flag = flag.Pick()

					self.last_flag_picked = 0
				}
			}
		}

	case (ebiten.IsKeyPressed(ebiten.KeySpace) || ebiten.IsKeyPressed(ebiten.KeyArrowUp)) && self.onGround:
		self.yv -= self.Jump_force * delta
		self.Y--

	case ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyArrowLeft):
		self.facing_left = true
		self.xv = -self.speed * delta

		self.animator.Play("run")

	case ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyArrowRight):
		self.facing_left = false
		self.xv = self.speed * delta

		self.animator.Play("run")

	case ebiten.IsKeyPressed(ebiten.KeyR):
		self.X = 0
		self.Y = 0
		self.xv = 0
		self.yv = 0

	case ebiten.IsKeyPressed(ebiten.Key1) && self.last_debug_toggled > 30:
		self.last_debug_toggled = 0
		self.worldData.Debug = !self.worldData.Debug

	default:
		self.xv = 0

		self.animator.Play("idle")

	}

	self.X += self.xv
	self.Y += self.yv
	self.collisionShape.MoveTo(self.X, self.Y-self.collisionShape.Height/2)
}

func (self *Player) Draw(screen *ebiten.Image) {
	self.animator.UpdateAndDraw()

	// Draw player
	op := &ebiten.DrawImageOptions{}
	if self.facing_left {
		op.GeoM.Scale(-1, 1)
		op.GeoM.Translate(float64(self.Sprite.Bounds().Dx()), 0)
	}
	op.GeoM.Translate(self.X-float64(self.Sprite.Bounds().Dx())/2, self.Y-float64(self.Sprite.Bounds().Dy()))

	screen.DrawImage(self.Sprite, op)

	if self.worldData.Debug {
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("[X: %.0f, Y: %.0f]", self.X, self.Y), int(self.X), int(self.Y))
		vector.StrokeRect(screen, float32(self.collisionShape.Pos.X-self.collisionShape.Width/2), float32(self.collisionShape.Pos.Y-self.collisionShape.Height/2), float32(self.collisionShape.Width), float32(self.collisionShape.Height), 2.0, color.RGBA{R: 200, G: 10, B: 10, A: 255}, true)
	}

	// Draw player name
	face := &text.GoTextFace{Source: self.worldData.Font, Size: 15}

	tOp := &text.DrawOptions{}
	tOp.GeoM.Translate(self.X-(text.Advance(self.Name, face))/2, self.Y-float64(self.Sprite.Bounds().Dy())-10)
	tOp.Filter = ebiten.FilterLinear

	text.Draw(screen, self.Name, face, tOp)
}

func (self *Player) Pos() (x, y float64) {
	return self.X, self.Y
}
