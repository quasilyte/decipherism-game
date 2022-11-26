package main

import (
	"math"

	"github.com/quasilyte/ge"
	"github.com/quasilyte/ge/gesignal"
	"github.com/quasilyte/ge/ui"
	"github.com/quasilyte/gmath"
)

type dialButton struct {
	pos       gmath.Vec
	centerPos gmath.Vec

	body          *ge.Sprite
	arrow         *ge.Sprite
	arrowRotation gmath.Rad
	arrowStep     gmath.Rad
	numStates     int
	state         int

	root   *ui.Root
	button *ui.Button

	EventActivated gesignal.Event[int]
}

func newDialButton(root *ui.Root, pos gmath.Vec, numStates int) *dialButton {
	return &dialButton{
		pos:       pos,
		centerPos: pos.Add(gmath.Vec{X: 56, Y: 50}),
		root:      root,
		numStates: numStates,
		arrowStep: math.Pi / gmath.Rad(numStates-1),
	}
}

func (b *dialButton) Init(scene *ge.Scene) {
	b.button = b.root.NewButton(invisibleButtonStyle.Resized(108, 108))
	b.button.Pos.Base = &b.pos
	b.button.EventActivated.Connect(nil, func(_ *ui.Button) {
		b.state++
		if b.state >= b.numStates {
			b.state = 0
			b.arrowRotation = math.Pi
		} else {
			b.arrowRotation += b.arrowStep
		}
		b.EventActivated.Emit(b.state)
	})
	scene.AddObject(b.button)

	b.body = scene.NewSprite(ImageDialButton)
	b.body.Pos.Base = &b.centerPos
	scene.AddGraphics(b.body)

	b.arrowRotation = math.Pi + (b.arrowStep * gmath.Rad(b.state))

	b.arrow = scene.NewSprite(ImageDialButtonArrow)
	b.arrow.Pos.Base = &b.centerPos
	b.arrow.Rotation = &b.arrowRotation
	scene.AddGraphics(b.arrow)
}

func (b *dialButton) IsDisposed() bool {
	return false
}

func (b *dialButton) Update(delta float64) {
}
