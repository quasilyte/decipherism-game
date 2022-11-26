package main

import (
	"fmt"
	"math"
	"strings"

	"github.com/quasilyte/gmath"
)

type schemaElemKind int

const (
	unknownElem schemaElemKind = iota
	inputElem
	outputElem
	muxElem
	simplePipeElem
	pipeConnect2Elem
	ifElem
	transformElem
)

type schemaElem struct {
	elemID      int
	tileClassID int
	tileClass   string
	kind        schemaElemKind

	next []*schemaElem

	pos gmath.Vec

	rotation gmath.Rad

	extraData any
}

type angleElemExtra struct {
	flipHorizontally bool
}

type countdownElemExtra struct {
	initialValue int
}

type ifElemExtra struct {
	condKind  string
	stringArg string
	intArg    int
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
		extra := extraData.(*angleElemExtra)
		if extra.flipHorizontally {
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

func getSchemaElemKind(tileClass string) schemaElemKind {
	switch tileClass {
	case "elem_input":
		return inputElem
	case "elem_output":
		return outputElem
	case "elem_mux":
		return muxElem
	case "elem_if", "elem_ifnot", "elem_repeater", "elem_inv_repeater":
		return ifElem
	case "pipe_connect2":
		return pipeConnect2Elem
	}

	if strings.HasPrefix(tileClass, "apply_") {
		return transformElem
	}
	if strings.Contains(tileClass, "_countdown") {
		return ifElem
	}
	if strings.Contains(tileClass, "pipe") {
		return simplePipeElem
	}

	panic(fmt.Sprintf("unexpected %q class", tileClass))
}
