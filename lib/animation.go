package lib

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type Animation struct {
	Source_sprite *ebiten.Image
	Length        int
	Speed         int

	Width  int
	Height int

	Start int

	frames []*ebiten.Image
	cached bool
}

func (animation *Animation) CacheFrames() {
	animation.frames = make([]*ebiten.Image, 0, animation.Length)
	for i := range animation.Length {
		animation.frames = append(animation.frames, animation.Source_sprite.SubImage(image.Rect(
			animation.Width*i+animation.Start,
			0,
			animation.Width*i+animation.Width+animation.Start,
			animation.Height,
		)).(*ebiten.Image))
	}
	animation.cached = true
}

func (animation *Animation) GetCachedFrame(index int) (*ebiten.Image, bool) {
	if index < animation.Length {
		return animation.frames[index], true
	}
	return nil, false
}
