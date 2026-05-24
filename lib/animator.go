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

	animation         Animation
	playing_animation string

	sprite         *ebiten.Image
	sprite_current int
}

func (Animator) New(animations map[string]Animation) *Animator {

	a := &Animator{
		Animations: animations,
	}

	return a
}

func (a *Animator) InitImage() *ebiten.Image {

	var width, height int

	for _, animation := range a.Animations {
		if animation.Width > width {
			width = animation.Width
		}
		if animation.Height > height {
			height = animation.Height
		}
	}

	a.sprite = ebiten.NewImage(width, height)

	return a.sprite
}

func (a *Animator) Play(play string) {
	if a.playing_animation == play {
		return
	}

	if _, ok := a.Animations[play]; !ok {
		log.Println("Animation", play, "wasn't found in animations")
	}

	a.playing_animation = play

	a.animation = a.Animations[play]

	a.ticks_before_reset = a.animation.Speed * a.animation.Length

	if !a.animation.cached {
		a.animation.CacheFrames()
	}

	a.ForceDraw()

}

func (a *Animator) GetPlaying() string {
	return a.playing_animation
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
