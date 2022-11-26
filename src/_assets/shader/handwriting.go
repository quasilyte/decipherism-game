//go:build ignore
// +build ignore

package main

func Fragment(pos vec4, texCoord vec2, _ vec4) vec4 {
	c := imageSrc0At(texCoord)

	pixSize := imageSrcTextureSize()
	originTexPos, _ := imageSrcRegionOnTexture()
	actualTexPos := vec2(texCoord.x-originTexPos.x, texCoord.y-originTexPos.y)
	actualPixPos := actualTexPos * pixSize
	if c[3] != 0.0 {
		p := actualPixPos
		posHash := int(actualPixPos.x+actualPixPos.y) * int(actualPixPos.y*5)
		state := posHash % 15
		if state == int(1) {
			p.x += 1.0
		} else if state == int(2) {
			p.x -= 1.0
		} else if state == int(3) {
			p.y += 1.0
		} else if state == int(4) {
			p.y -= 1.0
		} else {
			return c
		}
		colorMultiplier := vec4(0.95, 0.95, 0.95, 0.95)
		return imageSrc0At(p/pixSize+originTexPos) * colorMultiplier
	}

	return c
}
