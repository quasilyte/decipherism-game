package main

import (
	"fmt"
	"strings"

	"github.com/quasilyte/decipherism-game/leveldata"
	"github.com/quasilyte/ge/tiled"
	"github.com/quasilyte/gmath"
)

const (
	numSchemaCols int = 12
	numSchemaRows int = 8
)

type schemaBuilder struct {
	template    *leveldata.SchemaTemplate
	schema      *componentSchema
	offset      gmath.Vec
	elemByIndex [numSchemaCols * numSchemaRows]*schemaElem
}

type elemShape struct {
	inputs     [4]gmath.Vec
	outputs    [4]gmath.Vec
	numInputs  int
	numOutputs int
}

func (shape *elemShape) CanConnectTo(otherShape *elemShape) bool {
	for i := 0; i < shape.numOutputs; i++ {
		for j := 0; j < otherShape.numInputs; j++ {
			if shape.outputs[i].EqualApprox(otherShape.inputs[j]) {
				return true
			}
		}
	}
	return false
}

func decodeSchema(offset gmath.Vec, tileset *tiled.Tileset, data []byte) *componentSchema {
	m, err := tiled.UnmarshalMap(data)
	if err != nil {
		panic(err)
	}
	t, err := leveldata.TilemapToTemplate(tileset, m)
	if err != nil {
		panic(err) // Used only for builtin levels, should never panic
	}
	return newSchemaBuilder(offset, t).Build()
}

func newSchemaBuilder(offset gmath.Vec, t *leveldata.SchemaTemplate) *schemaBuilder {
	return &schemaBuilder{
		offset:   offset,
		template: t,
		schema:   &componentSchema{},
	}
}

func (b *schemaBuilder) Build() *componentSchema {
	b.build()
	return b.schema
}

func (b *schemaBuilder) elemShape(e *schemaElem) elemShape {
	return getElemShape(e.tileClass, e.pos, e.rotation, e.extraData)
}

func (b *schemaBuilder) rowcolByPos(pos gmath.Vec) (int, int) {
	col := int(pos.X-b.offset.X) / int(b.template.Tileset.TileWidth)
	row := int(pos.Y-b.offset.Y) / int(b.template.Tileset.TileHeight)
	return row, col
}

func (b *schemaBuilder) indexByPos(pos gmath.Vec) int {
	row, col := b.rowcolByPos(pos)
	return row*numSchemaCols + col
}

func (b *schemaBuilder) visitNeighbours(elem *schemaElem, f func(*schemaElem)) {
	elemRow, elemCol := b.rowcolByPos(elem.pos)
	toVisit := [...][2]int{
		{elemRow - 1, elemCol},
		{elemRow, elemCol - 1},
		{elemRow, elemCol + 1},
		{elemRow + 1, elemCol},
	}
	for _, rowcol := range toVisit {
		row := rowcol[0]
		col := rowcol[1]
		if (row < 0 || row >= numSchemaRows) || (col < 0 || col >= numSchemaCols) {
			continue
		}
		if e := b.elemByIndex[row*numSchemaCols+col]; e != nil {
			f(e)
		} else {
		}
	}
}

func (b *schemaBuilder) visitConnectedElems(elem *schemaElem, f func(e *schemaElem, outgoing bool)) {
	shape := b.elemShape(elem)
	b.visitNeighbours(elem, func(e *schemaElem) {
		otherShape := b.elemShape(e)
		if otherShape.CanConnectTo(&shape) {
			f(e, false)
			return
		}
		if shape.CanConnectTo(&otherShape) {
			f(e, true)
			return
		}
	})
}

func (b *schemaBuilder) connectElems(elem *schemaElem, minIncoming, wantOutgoing int) {
	numOutgoing := 0
	numIncoming := 0

	if wantOutgoing > 0 {
		elem.next = make([]*schemaElem, 0, wantOutgoing)
	}
	b.visitConnectedElems(elem, func(e *schemaElem, outgoing bool) {
		if !outgoing && e.next == nil {
			e.next = []*schemaElem{elem}
		}
		if outgoing {
			elem.next = append(elem.next, e)
			numOutgoing++
		} else {
			numIncoming++
		}
	})

	if wantOutgoing != -1 && numOutgoing != wantOutgoing {
		panic(fmt.Sprintf("%s: expected %d outgoing pipe(s), found %d", elem.tileClass, wantOutgoing, numOutgoing))
	}
	if minIncoming != -1 && numIncoming < minIncoming {
		panic(fmt.Sprintf("%s: expected at least %d incoming pipe(s), found %d", elem.tileClass, minIncoming, numIncoming))
	}

	if len(elem.next) > 1 {
		for i, e := range elem.next[1:] {
			if strings.HasPrefix(e.tileClass, "special_") {
				elem.next[0], elem.next[i+1] = elem.next[i+1], elem.next[0]
				break
			}
		}
	}
}

func (b *schemaBuilder) build() {
	s := b.schema
	s.numKeywords = b.template.NumKeywords
	s.keywords = b.template.Keywords

	elemList := make([]*schemaElem, 0, 24)

	for _, t := range b.template.Elems {
		tileClassID := t.ClassID
		if tileClassID == -1 {
			tileClassID = b.template.Tileset.TileByClass(t.Class).Index
		}
		elem := &schemaElem{
			pos:         t.Pos.Add(b.offset),
			tileClassID: tileClassID,
			tileClass:   t.Class,
			kind:        getSchemaElemKind(t.Class),
			rotation:    t.Rotation,
			extraData:   t.ExtraData,
		}
		// TODO: use a tileset metadata for that.
		switch {
		case strings.HasSuffix(elem.tileClass, "_dotted") || strings.HasSuffix(elem.tileClass, "_undotted") || strings.HasSuffix(elem.tileClass, "_even") || strings.HasSuffix(elem.tileClass, "_odd"):
			s.hasCondTransform = true
		case strings.Contains(elem.tileClass, "polygraphic"):
			s.hasPolygraphic = true
		case strings.Contains(elem.tileClass, "atbash"):
			s.hasAtbash = true
		case strings.Contains(elem.tileClass, "rot13"):
			s.hasRot13 = true
		case strings.Contains(elem.tileClass, "add") || strings.Contains(elem.tileClass, "sub"):
			s.hasIncDec = true
		case strings.Contains(elem.tileClass, "rotate") || strings.Contains(elem.tileClass, "reverse"):
			s.hasShift = true
		case strings.Contains(elem.tileClass, "ifnot") || strings.Contains(elem.tileClass, "inv_repeater"):
			s.hasNegation = true
		}
		elemList = append(elemList, elem)
		b.elemByIndex[b.indexByPos(elem.pos)] = elem
		if elem.tileClass == "elem_input" {
			s.entry = elem
		}
	}

	if s.entry == nil {
		panic("a component schema without input?")
	}

	id := 0
	for _, elem := range elemList {
		elem.elemID = id
		id++

		switch elem.kind {
		case transformElem, muxElem:
			b.connectElems(elem, 1, 1)
		case outputElem:
			// One or more inputs.
			numIncoming := 0
			b.visitConnectedElems(elem, func(e *schemaElem, outgoing bool) {
				if outgoing {
					panic(fmt.Sprintf("%s: unexpected outgoing pipe", elem.tileClass))
				}
				numIncoming++
				if e.next == nil {
					e.next = []*schemaElem{elem}
				}
			})
			if numIncoming == 0 {
				panic(fmt.Sprintf("%s: expected at least 1 incoming pipe", elem.tileClass))
			}
		case inputElem:
			// One simple outgoing pipe.
			numOutgoing := 0
			elem.next = make([]*schemaElem, 1)
			b.visitConnectedElems(elem, func(e *schemaElem, outgoing bool) {
				if !outgoing {
					panic(fmt.Sprintf("%s: unexpected incoming pipe", elem.tileClass))
				}
				numOutgoing++
				elem.next[0] = e
			})
			if numOutgoing != 1 {
				panic(fmt.Sprintf("%s: expected 1 outgoing pipe, found %d", elem.tileClass, numOutgoing))
			}
		case ifElem:
			// One incoming and two outgoing pipes.
			numOutgoing := 0
			numIncoming := 0
			elem.next = make([]*schemaElem, 2)
			b.visitConnectedElems(elem, func(e *schemaElem, outgoing bool) {
				if !outgoing {
					e.next = []*schemaElem{elem}
				}
				if outgoing {
					numOutgoing++
					if strings.HasPrefix(e.tileClass, "special_") {
						elem.next[0] = e
					} else {
						elem.next[1] = e
					}
				} else {
					numIncoming++
				}
			})
			if numOutgoing != 2 {
				panic(fmt.Sprintf("%s: expected 2 outgoing pipes, found %d", elem.tileClass, numOutgoing))
			}
			if elem.next[0] == nil || elem.next[1] == nil {
				panic(fmt.Sprintf("%s: invalid combination of outgoing pipes", elem.tileClass))
			}
			if numIncoming < 1 {
				panic(fmt.Sprintf("%s: expected at least 1 incoming pipe, found %d", elem.tileClass, numIncoming))
			}
		}
	}

	for _, elem := range elemList {
		switch elem.kind {
		case simplePipeElem:
			if elem.next == nil {
				b.connectElems(elem, 1, 1)
			}
		case pipeConnect2Elem:
			if elem.next == nil {
				b.connectElems(elem, 2, 1)
			}
		}
	}

	for _, elem := range elemList {
		if elem.next == nil {
			fmt.Println(elem.tileClass)
		}
	}

	s.elems = elemList
}
