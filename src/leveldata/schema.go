package leveldata

import (
	"math"
	"strings"

	"github.com/quasilyte/gmath"
)

type ComponentSchema struct {
	Entry *SchemaElem

	Elems []*SchemaElem

	NumKeywords     int
	Keywords        []string
	EncodedKeywords []string

	HasCondTransform bool
	HasPolygraphic   bool
	HasAtbash        bool
	HasRot13         bool
	HasIncDec        bool
	HasShift         bool
	HasNegation      bool
}

type SchemaElemKind int

const (
	UnknownElem SchemaElemKind = iota
	InputElem
	OutputElem
	MuxElem
	SimplePipeElem
	PipeConnect2Elem
	IfElem
	TransformElem
)

type SchemaElem struct {
	ElemID      int
	TileClassID int
	TileClass   string
	Kind        SchemaElemKind

	Next []*SchemaElem

	Pos gmath.Vec

	Rotation gmath.Rad

	ExtraData any
}

func getElemShape(class string, pos gmath.Vec, rotation gmath.Rad, extraData any) elemShape {
	var shape elemShape
	switch class {
	case "pipe_connect2":
		// Two inputs, one output.
		// The default orientation is:
		//
		//   v
		//   |
		//   .-->
		//   |
		//   ^
		//
		startRotation1 := rotation + (math.Pi / 2)
		startRotation2 := rotation - (math.Pi / 2)
		endRotation := rotation
		shape.inputs[0] = pos.MoveInDirection(48, startRotation1-math.Pi)
		shape.inputs[1] = pos.MoveInDirection(48, startRotation2-math.Pi)
		shape.numInputs = 2
		shape.outputs[0] = pos.MoveInDirection(48, endRotation)
		shape.numOutputs = 1

	case "angle_pipe", "special_angle_pipe":
		// One input, out output. Can be flipped.
		// The default orientation is:
		//
		//   >--.
		//      |
		//      v
		//
		startRotation := rotation
		endRotation := rotation + (math.Pi / 2)
		extra := extraData.(*AngleElemExtra)
		if extra.FlipHorizontally {
			startRotation += math.Pi
		}
		shape.inputs[0] = pos.MoveInDirection(48, startRotation-math.Pi)
		shape.numInputs = 1
		shape.outputs[0] = pos.MoveInDirection(48, endRotation)
		shape.numOutputs = 1

	case "pipe", "special_pipe":
		// One input, out output.
		shape.inputs[0] = pos.MoveInDirection(48, rotation-math.Pi)
		shape.numInputs = 1
		shape.outputs[0] = pos.MoveInDirection(48, rotation)
		shape.numOutputs = 1

	default:
		// The default shape allows any kind of connection from any direction.
		left := gmath.Vec{X: -48}
		right := gmath.Vec{X: 48}
		up := gmath.Vec{Y: -48}
		down := gmath.Vec{Y: 48}
		shape.inputs = [4]gmath.Vec{left, right, up, down}
		shape.numInputs = 4
		shape.outputs = [4]gmath.Vec{left, right, up, down}
		shape.numOutputs = 4
		for i := 0; i < shape.numInputs; i++ {
			shape.inputs[i] = shape.inputs[i].Add(pos)
		}
		for i := 0; i < shape.numOutputs; i++ {
			shape.outputs[i] = shape.outputs[i].Add(pos)
		}
	}

	return shape
}

func getSchemaElemKind(tileClass string) SchemaElemKind {
	switch tileClass {
	case "elem_input":
		return InputElem
	case "elem_output":
		return OutputElem
	case "elem_mux":
		return MuxElem
	case "elem_if", "elem_ifnot", "elem_repeater", "elem_inv_repeater":
		return IfElem
	case "pipe_connect2":
		return PipeConnect2Elem
	}

	if strings.HasPrefix(tileClass, "apply_") {
		return TransformElem
	}
	if strings.Contains(tileClass, "_countdown") {
		return IfElem
	}
	if strings.Contains(tileClass, "pipe") {
		return SimplePipeElem
	}

	return UnknownElem
}
