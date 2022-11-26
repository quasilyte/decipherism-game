package main

import (
	"github.com/quasilyte/ge"
	"github.com/quasilyte/ge/gesignal"
	"github.com/quasilyte/gmath"
)

type signalNode struct {
	pos    gmath.Vec
	sprite *ge.Sprite
	dst    gmath.Vec
	speed  float64

	EventDestinationReached gesignal.Event[*signalNode]
}

func newSignalNode(pos gmath.Vec) *signalNode {
	return &signalNode{
		pos:   pos,
		speed: 160,
	}
}

func (s *signalNode) Init(scene *ge.Scene) {
	s.sprite = scene.NewSprite(ImageSignal)
	s.sprite.Pos.Base = &s.pos
	scene.AddGraphics(s.sprite)
}

func (s *signalNode) IsDisposed() bool {
	return s.sprite.IsDisposed()
}

func (s *signalNode) Dispose() {
	s.sprite.Dispose()
}

func (s *signalNode) Update(delta float64) {
	if s.dst.IsZero() {
		return
	}

	step := s.speed * delta
	if s.pos.DistanceTo(s.dst) < step {
		s.pos = s.dst
		s.dst = gmath.Vec{}
		s.EventDestinationReached.Emit(s)
	} else {
		s.pos = s.pos.MoveTowards(s.dst, step)
	}
}
