package leveldata

import (
	"fmt"
	"strings"

	"github.com/quasilyte/ge/tiled"
	"github.com/quasilyte/gmath"
)

const (
	NumSchemaCols int = 12
	NumSchemaRows int = 8
)

type SchemaBuilder struct {
	template    *SchemaTemplate
	schema      *ComponentSchema
	offset      gmath.Vec
	elemByIndex [NumSchemaCols * NumSchemaRows]*SchemaElem
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

func DecodeSchema(offset gmath.Vec, tileset *tiled.Tileset, data []byte) *ComponentSchema {
	m, err := tiled.UnmarshalMap(data)
	if err != nil {
		panic(err)
	}
	t, err := TilemapToTemplate(tileset, m)
	if err != nil {
		panic(err) // Used only for builtin levels, should never panic
	}
	result, err := NewSchemaBuilder(offset, t).Build()
	if err != nil {
		panic(err)
	}
	return result
}

func NewSchemaBuilder(offset gmath.Vec, t *SchemaTemplate) *SchemaBuilder {
	return &SchemaBuilder{
		offset:   offset,
		template: t,
		schema:   &ComponentSchema{},
	}
}

func (b *SchemaBuilder) Build() (result *ComponentSchema, err error) {
	defer func() {
		rv := recover()
		if rv == nil {
			return
		}
		e, ok := rv.(error)
		if ok {
			err = e
			return
		}
		panic(rv)
	}()

	b.build()
	return b.schema, err
}

func (b *SchemaBuilder) elemShape(e *SchemaElem) elemShape {
	return getElemShape(e.TileClass, e.Pos, e.Rotation, e.ExtraData)
}

func (b *SchemaBuilder) rowcolByPos(pos gmath.Vec) (int, int) {
	col := int(pos.X-b.offset.X) / int(b.template.Tileset.TileWidth)
	row := int(pos.Y-b.offset.Y) / int(b.template.Tileset.TileHeight)
	return row, col
}

func (b *SchemaBuilder) indexByPos(pos gmath.Vec) int {
	row, col := b.rowcolByPos(pos)
	return row*NumSchemaCols + col
}

func (b *SchemaBuilder) visitNeighbours(elem *SchemaElem, f func(*SchemaElem)) {
	elemRow, elemCol := b.rowcolByPos(elem.Pos)
	toVisit := [...][2]int{
		{elemRow - 1, elemCol},
		{elemRow, elemCol - 1},
		{elemRow, elemCol + 1},
		{elemRow + 1, elemCol},
	}
	for _, rowcol := range toVisit {
		row := rowcol[0]
		col := rowcol[1]
		if (row < 0 || row >= NumSchemaRows) || (col < 0 || col >= NumSchemaCols) {
			continue
		}
		if e := b.elemByIndex[row*NumSchemaCols+col]; e != nil {
			f(e)
		} else {
		}
	}
}

func (b *SchemaBuilder) visitConnectedElems(elem *SchemaElem, f func(e *SchemaElem, outgoing bool)) {
	shape := b.elemShape(elem)
	b.visitNeighbours(elem, func(e *SchemaElem) {
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

func (b *SchemaBuilder) connectElems(elem *SchemaElem, minIncoming, wantOutgoing int) {
	numOutgoing := 0
	numIncoming := 0

	if wantOutgoing > 0 {
		elem.Next = make([]*SchemaElem, 0, wantOutgoing)
	}
	b.visitConnectedElems(elem, func(e *SchemaElem, outgoing bool) {
		if !outgoing && e.Next == nil {
			e.Next = []*SchemaElem{elem}
		}
		if outgoing {
			elem.Next = append(elem.Next, e)
			numOutgoing++
		} else {
			numIncoming++
		}
	})

	if wantOutgoing != -1 && numOutgoing != wantOutgoing {
		b.errorf(elem, "expected %d outgoing pipe(s), found %d", wantOutgoing, numOutgoing)
	}
	if minIncoming != -1 && numIncoming < minIncoming {
		b.errorf(elem, "expected at least %d incoming pipe(s), found %d", minIncoming, numIncoming)
	}

	if len(elem.Next) > 1 {
		for i, e := range elem.Next[1:] {
			if strings.HasPrefix(e.TileClass, "special_") {
				elem.Next[0], elem.Next[i+1] = elem.Next[i+1], elem.Next[0]
				break
			}
		}
	}
}

func (b *SchemaBuilder) build() {
	s := b.schema
	s.NumKeywords = b.template.NumKeywords
	s.Keywords = b.template.Keywords

	numInputs := 0
	foundOutput := false
	elemList := make([]*SchemaElem, 0, 24)

	for _, t := range b.template.Elems {
		tileClassID := t.ClassID
		if tileClassID == -1 {
			tileClassID = b.template.Tileset.TileByClass(t.Class).Index
		}
		elemKind := getSchemaElemKind(t.Class)
		elem := &SchemaElem{
			Pos:         t.Pos.Add(b.offset),
			TileClassID: tileClassID,
			TileClass:   t.Class,
			Kind:        elemKind,
			Rotation:    t.Rotation,
			ExtraData:   t.ExtraData,
		}
		if elemKind == UnknownElem {
			b.errorf(elem, "unexpected elem class: %s", t.Class)
		}
		// TODO: use a tileset metadata for that.
		switch {
		case strings.HasSuffix(elem.TileClass, "_dotted") || strings.HasSuffix(elem.TileClass, "_undotted") || strings.HasSuffix(elem.TileClass, "_even") || strings.HasSuffix(elem.TileClass, "_odd"):
			s.HasCondTransform = true
		case strings.Contains(elem.TileClass, "polygraphic"):
			s.HasPolygraphic = true
		case strings.Contains(elem.TileClass, "atbash"):
			s.HasAtbash = true
		case strings.Contains(elem.TileClass, "rot13"):
			s.HasRot13 = true
		case strings.Contains(elem.TileClass, "add") || strings.Contains(elem.TileClass, "sub"):
			s.HasIncDec = true
		case strings.Contains(elem.TileClass, "rotate") || strings.Contains(elem.TileClass, "reverse"):
			s.HasShift = true
		case strings.Contains(elem.TileClass, "ifnot") || strings.Contains(elem.TileClass, "inv_repeater"):
			s.HasNegation = true
		}
		elemList = append(elemList, elem)
		b.elemByIndex[b.indexByPos(elem.Pos)] = elem
		if elem.TileClass == "elem_input" {
			s.Entry = elem
			numInputs++
		}
		if elem.TileClass == "elem_output" {
			foundOutput = true
		}
	}

	if !foundOutput {
		panic(fmt.Errorf("expected at least 1 OUT (output) element, found 0"))
	}
	if numInputs != 1 {
		panic(fmt.Errorf("expected exactly 1 IN (input) element, found %d", numInputs))
	}

	id := 0
	for _, elem := range elemList {
		elem.ElemID = id
		id++

		switch elem.Kind {
		case TransformElem, MuxElem:
			b.connectElems(elem, 1, 1)
		case OutputElem:
			// One or more inputs.
			numIncoming := 0
			b.visitConnectedElems(elem, func(e *SchemaElem, outgoing bool) {
				if outgoing {
					b.errorf(elem, "unexpected outgoing pipe")
				}
				numIncoming++
				if e.Next == nil {
					e.Next = []*SchemaElem{elem}
				}
			})
			if numIncoming == 0 {
				b.errorf(elem, "expected at least 1 incoming pipe")
			}
		case InputElem:
			// One simple outgoing pipe.
			numOutgoing := 0
			elem.Next = make([]*SchemaElem, 1)
			b.visitConnectedElems(elem, func(e *SchemaElem, outgoing bool) {
				if !outgoing {
					b.errorf(elem, "unexpected incoming pipe")
				}
				numOutgoing++
				elem.Next[0] = e
			})
			if numOutgoing != 1 {
				b.errorf(elem, "expected 1 outgoing pipe, found %d", numOutgoing)
			}
		case IfElem:
			// One incoming and two outgoing pipes.
			numOutgoing := 0
			numIncoming := 0
			elem.Next = make([]*SchemaElem, 2)
			b.visitConnectedElems(elem, func(e *SchemaElem, outgoing bool) {
				if !outgoing {
					e.Next = []*SchemaElem{elem}
				}
				if outgoing {
					numOutgoing++
					if strings.HasPrefix(e.TileClass, "special_") {
						elem.Next[0] = e
					} else {
						elem.Next[1] = e
					}
				} else {
					numIncoming++
				}
			})
			if numOutgoing != 2 {
				b.errorf(elem, "expected 2 outgoing pipes, found %d", numOutgoing)
			}
			if elem.Next[0] == nil || elem.Next[1] == nil {
				b.errorf(elem, "invalid combination of outgoing pipes")
			}
			if numIncoming < 1 {
				b.errorf(elem, "expected at least 1 incoming pipe, found %d", numIncoming)
			}
		}
	}

	for _, elem := range elemList {
		switch elem.Kind {
		case SimplePipeElem:
			if elem.Next == nil {
				b.connectElems(elem, 1, 1)
			}
		case PipeConnect2Elem:
			if elem.Next == nil {
				b.connectElems(elem, 2, 1)
			}
		}
	}

	for _, elem := range elemList {
		if elem.Next == nil && elem.TileClass != "elem_output" {
			b.errorf(elem, "elem is not properly connected")
		}
	}

	s.Elems = elemList
}

func (b *SchemaBuilder) errorf(elem *SchemaElem, format string, args ...any) {
	panic(fmt.Errorf("%v: %s: %s", elem.Pos, elem.TileClass, fmt.Sprintf(format, args...)))
}
