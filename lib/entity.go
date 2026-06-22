package lib

import "github.com/hajimehoshi/ebiten/v2"

type Entity interface {
	Update(delta float64)
	Draw(screen *ebiten.Image)
	Pos() (x, y float64)
}
