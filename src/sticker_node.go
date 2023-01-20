package main

import (
	"github.com/quasilyte/ge"
	"github.com/quasilyte/gmath"
)

type stickerNode struct {
	pos gmath.Vec

	sprite *ge.Sprite
	label  *ge.Label

	text string
}

func newStickerNode(pos gmath.Vec, text string) *stickerNode {
	return &stickerNode{
		pos:  pos,
		text: text,
	}
}

func (s *stickerNode) Init(scene *ge.Scene) {
	s.sprite = scene.NewSprite(ImageHintSticker)
	s.sprite.Centered = false
	s.sprite.Pos.Base = &s.pos
	scene.AddGraphics(s.sprite)

	s.label = scene.NewLabel(FontHandwrittenSmall)
	s.label.Pos.Base = &s.pos
	s.label.Pos.Offset = gmath.Vec{X: 34, Y: 72}
	s.label.Text = s.text
	s.label.ColorScale.SetRGBA(30, 30, 60, 220)
	scene.AddGraphics(s.label)
}

func (s *stickerNode) IsDisposed() bool { return false }

func (s *stickerNode) Update(delta float64) {}
