package main

import (
	"math"

	"github.com/quasilyte/ge"
	"github.com/quasilyte/gmath"
)

type pingEffectParticle struct {
	velocity gmath.Vec
	pos      gmath.Vec
	rotation gmath.Rad
	sprite   *ge.Sprite
}

type pingEffectNode struct {
	pos       gmath.Vec
	c         ge.ColorScale
	particles [4]pingEffectParticle
}

func newPingEffectNode(pos gmath.Vec, c ge.ColorScale) *pingEffectNode {
	return &pingEffectNode{
		pos: pos,
		c:   c,
	}
}

func (e *pingEffectNode) Init(scene *ge.Scene) {
	rotation := gmath.Rad(0)
	for i := range e.particles {
		p := &e.particles[i]
		p.rotation = rotation
		p.pos = e.pos.MoveInDirection(32, rotation-(math.Pi/4))
		p.velocity = gmath.RadToVec(rotation - (math.Pi / 4)).Mulf(20)
		p.sprite = scene.NewSprite(ImagePingParticle)
		p.sprite.Rotation = &p.rotation
		p.sprite.Pos.Base = &p.pos
		p.sprite.SetColorScale(e.c)
		scene.AddGraphics(p.sprite)
		rotation += math.Pi / 2
	}
}

func (e *pingEffectNode) IsDisposed() bool {
	return e.particles[0].sprite.IsDisposed()
}

func (e *pingEffectNode) Update(delta float64) {
	alpha := e.particles[0].sprite.GetAlpha()
	if alpha < 0.1 {
		for _, p := range e.particles {
			p.sprite.Dispose()
		}
		return
	}
	for i := range e.particles {
		p := &e.particles[i]
		p.sprite.SetAlpha(alpha - float32(delta*2))
		p.pos = p.pos.Add(p.velocity.Mulf(delta))
	}
}
