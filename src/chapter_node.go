package main

import (
	"math"

	"github.com/quasilyte/ge"
	"github.com/quasilyte/gmath"
)

type chapterNode struct {
	pos       gmath.Vec
	name      string
	rotation  gmath.Rad
	available bool
	completed bool
}

func newChapterNode(pos gmath.Vec, name string, available, completed bool) *chapterNode {
	return &chapterNode{
		pos:       pos,
		name:      name,
		available: available,
		completed: completed,
	}
}

func (n *chapterNode) Init(scene *ge.Scene) {
	if n.available {
		s := scene.NewSprite(ImageBlueMarker)
		s.Centered = false
		s.Pos.Base = &n.pos
		if scene.Rand().Bool() {
			n.rotation = math.Pi
			s.Pos.Offset = gmath.Vec{X: -42.5 * 4, Y: -42.5 * 4}
		}
		s.Rotation = &n.rotation
		scene.AddGraphics(s)
	}
	if n.completed {
		s := scene.NewSprite(ImageCompleteMark)
		s.Centered = false
		s.Pos.Base = &n.pos
		s.Pos.Offset = gmath.Vec{X: -42.5, Y: 42.5 * 3}
		scene.AddGraphics(s)
	}
}

func (n *chapterNode) IsDisposed() bool { return false }

func (n *chapterNode) Update(delta float64) {}
