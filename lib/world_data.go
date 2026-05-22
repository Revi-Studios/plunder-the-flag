package lib

import collider "github.com/melonfunction/ebiten-collider"

type WorldData struct {
	Gravity float64
	Hash    *collider.SpatialHash
}
