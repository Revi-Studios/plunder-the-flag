package lib

import (
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	collider "github.com/melonfunction/ebiten-collider"
)

type WorldData struct {
	Gravity float64
	Hash    *collider.SpatialHash
	Font    *text.GoTextFaceSource
	Debug   bool
}
