package leveldata

import (
	"fmt"
	"strings"

	"github.com/quasilyte/ge/tiled"
	"github.com/quasilyte/gmath"
)

type SchemaTemplate struct {
	Tileset     *tiled.Tileset
	Elems       []SchemaTemplateElem
	NumKeywords int
	Keywords    []string
	Hints       []SchemaHintTemplate
}

type SchemaHintTemplate struct {
	Text string
	Pos  gmath.Vec
}

type SchemaTemplateElem struct {
	Class     string
	ClassID   int
	Rotation  gmath.Rad
	ExtraData any
	Pos       gmath.Vec
}

type AngleElemExtra struct {
	FlipHorizontally bool
}

type CountdownElemExtra struct {
	InitialValue int
}

type IfElemExtra struct {
	CondKind  string
	StringArg string
	IntArg    int
}

func ValidateLevelData(tileset *tiled.Tileset, levelData []byte) error {
	tmpl, err := LoadLevelTemplate(tileset, levelData)
	if err != nil {
		return err
	}
	_, err = NewSchemaBuilder(gmath.Vec{}, tmpl).Build()
	return err
}

func LoadLevelTemplate(tileset *tiled.Tileset, levelData []byte) (*SchemaTemplate, error) {
	m, err := tiled.UnmarshalMap(levelData)
	if err != nil {
		return nil, err
	}
	return TilemapToTemplate(tileset, m)
}

func TilemapToTemplate(tileset *tiled.Tileset, m *tiled.Map) (*SchemaTemplate, error) {
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

	var result SchemaTemplate

	result.Tileset = tileset

	elemList := make([]SchemaTemplateElem, 0, 24)

	ref := m.Tilesets[0]
	layer := m.Layers[0]

	foundSettings := false
	for _, o := range layer.Objects {
		id := o.GID - ref.FirstGID
		t := tileset.TileByID(id)
		pos := calcObjectPos(o)
		if t.Class == "settings" {
			if foundSettings {
				return nil, fmt.Errorf("%v: found more than one settings element", pos)
			}
			foundSettings = true
			allKeywords := strings.TrimSpace(o.GetStringProp("keywords", ""))
			if allKeywords == "" {
				return nil, fmt.Errorf("%v: settings.keywords property is empty", pos)
			}
			keywordList := strings.Split(allKeywords, "\n")
			result.NumKeywords = o.GetIntProp("num_keywords", 0)
			result.Keywords = keywordList
			continue
		}
		if t.Class == "hint" {
			pos := gmath.Vec{X: float64(o.X), Y: float64(o.Y)}
			result.Hints = append(result.Hints, SchemaHintTemplate{
				Text: o.GetStringProp("text", ""),
				Pos:  pos,
			})
			continue
		}
		elem := SchemaTemplateElem{
			Pos:      pos,
			ClassID:  t.Index,
			Class:    t.Class,
			Rotation: gmath.DegToRad(float64(o.Rotation)),
		}
		switch elem.Class {
		case "angle_pipe", "special_angle_pipe":
			if o.FlippedVertically() {
				return nil, fmt.Errorf("%v: vertical flipping is obsolete", pos)
			}
			elem.ExtraData = &AngleElemExtra{
				FlipHorizontally: o.FlippedHorizontally(),
			}
		case "elem_countdown0", "elem_countdown1", "elem_countdown2", "elem_countdown3":
			extra := &CountdownElemExtra{
				InitialValue: 3,
			}
			switch elem.Class {
			case "elem_countdown0":
				extra.InitialValue = 0
			case "elem_countdown1":
				extra.InitialValue = 1
			case "elem_countdown2":
				extra.InitialValue = 2
			}
			elem.ExtraData = extra
		case "elem_if", "elem_ifnot":
			extra := &IfElemExtra{
				CondKind:  o.GetStringProp("cond_kind", ""),
				StringArg: o.GetStringProp("string_arg", ""),
				IntArg:    o.GetIntProp("int_arg", 0),
			}
			if extra.CondKind == "" {
				return nil, fmt.Errorf("%v: elem_if cond property is empty", pos)
			}
			elem.ExtraData = extra
		}
		elemList = append(elemList, elem)
	}

	result.Elems = elemList

	return &result, nil
}
