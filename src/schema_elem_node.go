package main

import (
	"github.com/quasilyte/decipherism-game/leveldata"
	resource "github.com/quasilyte/ebitengine-resource"
	"github.com/quasilyte/ge"
	"github.com/quasilyte/gmath"
)

type schemaElemNode struct {
	data *leveldata.SchemaElem

	rotation gmath.Rad

	shaderEnabled   bool
	shaderTick      float64
	shaderTickDelay float64
	shaderSeed      int
	shaderStep      int

	sprite *ge.Sprite
}

func newSchemaElemNode(data *leveldata.SchemaElem, shaderEnabled bool) *schemaElemNode {
	return &schemaElemNode{
		data:          data,
		rotation:      data.Rotation,
		shaderEnabled: shaderEnabled,
	}
}

func (n *schemaElemNode) Init(scene *ge.Scene) {
	imageID := resource.ImageID(n.data.TileClassID) + componentSchemaImageOffset + 1
	n.sprite = scene.NewSprite(imageID)
	n.sprite.Pos.Base = &n.data.Pos
	n.sprite.Rotation = &n.rotation
	if extra, ok := n.data.ExtraData.(*leveldata.AngleElemExtra); ok {
		n.sprite.FlipHorizontal = extra.FlipHorizontally
	}
	scene.AddGraphics(n.sprite)
	if n.shaderEnabled {
		n.sprite.Shader = scene.NewShader(ShaderVideoDistortion)
		n.shaderStep = scene.Rand().IntRange(1, 4)
		n.sprite.Shader.SetIntValue("Seed", scene.Rand().IntRange(10, 500))
		n.shaderTickDelay = 0.2
		n.shaderTick = float64(scene.Rand().IntRange(0, 2))
	}
}

func (n *schemaElemNode) IsDisposed() bool {
	return n.sprite.IsDisposed()
}

func (n *schemaElemNode) Update(delta float64) {
	if !n.shaderEnabled {
		return
	}
	n.shaderTickDelay -= delta
	if n.shaderTickDelay <= 0 {
		n.shaderTickDelay = 0.3
		n.shaderTick += 1.0
		if n.shaderTick > 3.0 {
			n.shaderTick = 0
		}
		n.shaderSeed += n.shaderStep
		if n.shaderSeed > 1000 {
			n.shaderSeed = 10
		}
		n.sprite.Shader.SetIntValue("Seed", n.shaderSeed+1)
		n.sprite.Shader.SetFloatValue("Tick", n.shaderTick)
	}
}
