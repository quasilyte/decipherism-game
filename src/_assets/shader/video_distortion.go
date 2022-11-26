//go:build ignore
// +build ignore

package main

var Tick float
var Seed float

func Fragment(pos vec4, texCoord vec2, _ vec4) vec4 {
	c := imageSrc0At(texCoord)

	colorMultiplier := vec4(1, 1, 1, 1)
	if c[3] != 0.0 {
		if int(pos.y+Tick)%4 != int(0) {
			colorMultiplier = vec4(0.6, 0.6, 0.6, 0.6)
		}
	}

	pixSize := imageSrcTextureSize()
	originTexPos, _ := imageSrcRegionOnTexture()
	actualTexPos := vec2(texCoord.x-originTexPos.x, texCoord.y-originTexPos.y)
	actualPixPos := actualTexPos * pixSize
	intSeed := int(Seed)
	if c[3] != 0.0 {
		p := actualPixPos
		pixelOffset := int(actualPixPos.x) + int(actualPixPos.y*pixSize.x)
		seedMod := pixelOffset % intSeed
		pixelOffset += seedMod

		posHash := pixelOffset + intSeed
		dir := posHash % 5
		dist := 1.0
		if seedMod == int(0) {
			dist = 2.0
		}
		if dir == int(1) {
			p.x += dist
		} else if dir == int(2) {
			p.x -= dist
		} else if dir == int(3) {
			p.y += dist
		} else if dir == int(4) {
			p.y -= dist
		}
		return imageSrc0At(p/pixSize+originTexPos) * colorMultiplier
	}

	return c
}
