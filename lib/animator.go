package lib

import (
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

type Animator struct {
	tick               int
	ticks_before_reset int

	Animations map[string]Animation

	animation Animation

	sprite         *ebiten.Image
	sprite_current int
}

func (Animator) New(sprite *ebiten.Image, animations map[string]Animation) *Animator {

	a := &Animator{
		Animations: animations,
		sprite:     sprite,
	}

	return a
}

func (a *Animator) Start(play string) *ebiten.Image {
	if _, ok := a.Animations[play]; !ok {
		log.Println("Animation", play, "wasn't found in animations")
	}

	a.animation = a.Animations[play]

	a.ticks_before_reset = a.animation.Speed * a.animation.Length

	a.sprite.Deallocate()

	a.sprite = ebiten.NewImage(a.animation.Width, a.animation.Height)

	if !a.animation.cached {
		a.animation.CacheFrames()
	}

	a.ForceDraw()

	return a.sprite
}

// Updates and draws a new frame if needed
func (a *Animator) UpdateAndDraw() {
	a.tick++
	if a.tick >= a.ticks_before_reset {
		a.tick = 0
	}

	if new := a.tick / a.animation.Speed; new != a.sprite_current {

		a.sprite.Clear()
		op := &ebiten.DrawImageOptions{}

		if cached, ok := a.animation.GetCachedFrame(new); ok {
			a.sprite.DrawImage(cached, op)
		} else {
			a.sprite.DrawImage(a.animation.Source_sprite.SubImage(image.Rect(
				a.animation.Width*new,
				0,
				a.animation.Width*new+a.animation.Width,
				a.animation.Height,
			)).(*ebiten.Image), op)
		}

		a.sprite_current = new
	}
}

// Forces a render without updating the internal animation tick
func (a *Animator) ForceDraw() {
	a.sprite.Clear()
	op := &ebiten.DrawImageOptions{}

	a.sprite.DrawImage(a.animation.Source_sprite.SubImage(image.Rect(
		a.animation.Width*a.sprite_current,
		0,
		a.animation.Width*a.sprite_current+a.animation.Width,
		a.animation.Height,
	)).(*ebiten.Image), op)
}
