package main

import (
	"fmt"
	"strings"

	"github.com/quasilyte/ge/tiled"
	"github.com/quasilyte/gmath"
)

type schemaTemplate struct {
	tileset     *tiled.Tileset
	elems       []schemaTemplateElem
	numKeywords int
	keywords    []string
}

type schemaTemplateElem struct {
	class     string
	classID   int
	rotation  gmath.Rad
	extraData any
	pos       gmath.Vec
}

func tilemapToTemplate(tileset *tiled.Tileset, m *tiled.Map) *schemaTemplate {
	calcObjectPos := func(o tiled.Object) gmath.Vec {
		pos := gmath.Vec{X: float64(o.X) + tileset.TileWidth/2, Y: float64(o.Y) - tileset.TileHeight/2}
		switch o.Rotation {
		case 90:
			pos.Y += tileset.TileHeight
		case 180:
			pos.X -= tileset.TileWidth
			pos.Y += tileset.TileHeight
		case 270:
			pos.X -= tileset.TileWidth
		}
		return pos
	}

	var result schemaTemplate

	result.tileset = tileset

	elemList := make([]schemaTemplateElem, 0, 24)

	ref := m.Tilesets[0]
	layer := m.Layers[0]

	foundSettings := false
	for _, o := range layer.Objects {
		id := o.GID - ref.FirstGID
		t := tileset.TileByID(id)
		if t.Class == "settings" {
			if foundSettings {
				panic("more than one settings object?")
			}
			foundSettings = true
			allKeywords := strings.TrimSpace(o.GetStringProp("keywords", ""))
			if allKeywords == "" {
				panic("settings.keywords is not set")
			}
			keywordList := strings.Split(allKeywords, "\n")
			result.numKeywords = o.GetIntProp("num_keywords", 0)
			result.keywords = keywordList
			continue
		}
		pos := calcObjectPos(o)
		elem := schemaTemplateElem{
			pos:      pos,
			classID:  t.Index,
			class:    t.Class,
			rotation: gmath.DegToRad(float64(o.Rotation)),
		}
		switch elem.class {
		case "angle_pipe", "special_angle_pipe":
			if o.FlippedVertically() {
				panic(fmt.Sprintf("%s: vertical flipping is obsolete", elem.pos))
			}
			elem.extraData = &angleElemExtra{
				flipHorizontally: o.FlippedHorizontally(),
			}
		case "elem_countdown0", "elem_countdown1", "elem_countdown2", "elem_countdown3":
			extra := &countdownElemExtra{
				initialValue: t.Index - int(ImageElemCountdown0-componentSchemaImageOffset) + 1,
			}
			elem.extraData = extra
		case "elem_if", "elem_ifnot":
			extra := &ifElemExtra{
				condKind:  o.GetStringProp("cond_kind", ""),
				stringArg: o.GetStringProp("string_arg", ""),
				intArg:    o.GetIntProp("int_arg", 0),
			}
			if extra.condKind == "" {
				panic(fmt.Sprintf("elem_if with empty cond_kind"))
			}
			elem.extraData = extra

		}
		elemList = append(elemList, elem)
	}

	result.elems = elemList

	return &result
}
