package main

import (
	"bytes"
	"fmt"

	"github.com/quasilyte/gmath"
)

type schemaRunner struct {
	schema   *componentSchema
	current  *schemaElem
	input    []byte
	data     []byte
	counters [numSchemaCols * numSchemaRows]uint8
	lastCond bool
}

func newSchemaRunner() *schemaRunner {
	return &schemaRunner{
		data: make([]byte, 0, 16),
	}
}

func (r *schemaRunner) Exec(s *componentSchema, k string) string {
	r.Reset(s, []byte(k))
	for {
		_, hasMore := r.RunStep()
		if !hasMore {
			break
		}
	}
	return string(r.data)
}

func (r *schemaRunner) Reset(s *componentSchema, input []byte) {
	r.schema = s
	r.current = s.entry
	r.lastCond = false

	r.counters = [numSchemaCols * numSchemaRows]uint8{}
	for _, e := range r.schema.elems {
		countdownData, ok := e.extraData.(*countdownElemExtra)
		if ok {
			r.counters[e.elemID] = uint8(countdownData.initialValue)
		}
	}

	r.input = input
	r.data = r.data[:0]
	r.data = append(r.data, input...)
}

func (r *schemaRunner) RunStep() (gmath.Vec, bool) {
	if r.current.tileClass == "elem_output" {
		return gmath.Vec{}, false
	}

	var dst gmath.Vec
	switch r.current.kind {
	case transformElem:
		r.runTransformElem()
		r.current = r.current.next[0]
		dst = r.current.pos
	case muxElem:
		r.current = r.current.next[0]
		dst = r.current.pos
	case simplePipeElem, pipeConnect2Elem, inputElem:
		r.current = r.current.next[0]
		dst = r.current.pos
	case ifElem:
		switch r.current.tileClass {
		case "elem_ifnot":
			r.current = r.runIfNot()
			dst = r.current.pos
		case "elem_if":
			r.current = r.runIf()
			dst = r.current.pos
		case "elem_repeater":
			r.current = r.runRepeater()
			dst = r.current.pos
		case "elem_inv_repeater":
			r.current = r.runInvRepeater()
			dst = r.current.pos
		case "elem_countdown0", "elem_countdown1", "elem_countdown2", "elem_countdown3":
			value := r.counters[r.current.elemID]
			if value != 0 {
				r.counters[r.current.elemID]--
				r.current = r.current.next[0]
			} else {
				r.current = r.current.next[1]
			}
			dst = r.current.pos
		default:
			panic(fmt.Sprintf("unexpected if: %q", r.current.tileClass))
		}
	default:
		panic(fmt.Sprintf("unhandled %q", r.current.tileClass))
	}

	return dst, true
}

func (r *schemaRunner) runTransformElem() {
	switch r.current.tileClass {
	case "apply_reverse":
		r.runElemReverse()
	case "apply_swap_halves":
		r.runSwapHalves()
	case "apply_rotate_right":
		rotateCharsRight(r.data)
	case "apply_rotate_right_butfirst":
		rotateCharsRight(r.data[1:])
	case "apply_rotate_left":
		rotateCharsLeft(r.data)
	case "apply_rotate_left_butfirst":
		rotateCharsLeft(r.data[1:])
	case "apply_rot13":
		mapChars(r.data, r.rot13Char)
	case "apply_rot13_butfirst":
		mapCharsButfirst(r.data, r.rot13Char)
	case "apply_rot13_butlast":
		mapCharsButlast(r.data, r.rot13Char)
	case "apply_rot13_first":
		r.data[0] = r.rot13Char(r.data[0])
	case "apply_polygraphic_atbash":
		polygraphicAtbash(r.data)
	case "apply_atbash":
		mapChars(r.data, r.atbashChar)
	case "apply_atbash_butlast":
		mapCharsButlast(r.data, r.atbashChar)
	case "apply_atbash_first":
		r.data[0] = r.atbashChar(r.data[0])
	case "apply_add":
		r.runAdd()
	case "apply_add_butfirst":
		mapCharsButfirst(r.data, incChar)
	case "apply_add_last":
		r.data[len(r.data)-1] = incChar(r.data[len(r.data)-1])
	case "apply_add_first":
		r.data[0] = incChar(r.data[0])
	case "apply_add_nowrap":
		mapChars(r.data, r.incCharNowrap)
	case "apply_add_butfirst_nowrap":
		mapCharsButfirst(r.data, r.incCharNowrap)
	case "apply_add_dotted":
		mapChars(r.data, incCharDotted)
	case "apply_add_odd":
		mapOddChars(r.data, incChar)
	case "apply_add_even":
		mapEvenChars(r.data, incChar)

	case "apply_sub_first":
		r.data[0] = decChar(r.data[0])
	case "apply_sub_last":
		r.data[len(r.data)-1] = decChar(r.data[len(r.data)-1])
	case "apply_sub_undotted":
		mapChars(r.data, decCharUndotted)
	case "apply_sub_odd":
		mapOddChars(r.data, decChar)
	case "apply_sub_even":
		mapEvenChars(r.data, decChar)

	case "apply_sub":
		r.runSub()
	case "apply_sub_butlast":
		mapCharsButlast(r.data, decChar)
	case "apply_sub_nowrap":
		r.runSubNowrap()
	case "apply_hardshift_left":
		mapChars(r.data, r.hardshiftLeftChar)
	case "apply_hardshift_right":
		mapChars(r.data, r.hardshiftRightChar)
	case "apply_zigzag":
		r.runZigzag(r.data)

	default:
		panic(fmt.Sprintf("unexpected transform: %q", r.current.tileClass))
	}
}

func (r *schemaRunner) hardshiftLeftChar(b byte) byte {
	if b < 'n' {
		return b
	}
	return r.atbashChar(b)
}

func (r *schemaRunner) hardshiftRightChar(b byte) byte {
	if b < 'n' {
		return r.atbashChar(b)
	}
	return b
}

func (r *schemaRunner) rot13Char(b byte) byte {
	if b < 'n' {
		return 'n' + (b - 'a')
	}
	return 'a' + (b - 'n')
}

func (r *schemaRunner) atbashChar(b byte) byte {
	return 'a' + (25 - (b - 'a'))
}

func (r *schemaRunner) runAdd() {
	for i, b := range r.data {
		r.data[i] = incChar(b)
	}
}

func (r *schemaRunner) runSub() {
	for i, b := range r.data {
		r.data[i] = decChar(b)
	}
}

func (r *schemaRunner) runSubNowrap() {
	for i, b := range r.data {
		r.data[i] = r.decCharNowrap(b)
	}
}

func (r *schemaRunner) runZigzag(data []byte) {
	for i := 0; i < len(data)-1; i += 2 {
		data[i], data[i+1] = data[i+1], data[i]
	}
}

func (r *schemaRunner) runSwapHalves() {
	if len(r.data) < 2 {
		return
	}
	mid := len(r.data) / 2
	end := mid
	offset := 0
	if len(r.data)%2 != 0 {
		offset = 1
	}
	for i := 0; i < end; i++ {
		j := mid + i + offset
		r.data[i], r.data[j] = r.data[j], r.data[i]
	}
}

func (r *schemaRunner) runElemReverse() {
	b := r.data
	for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}
}

func (r *schemaRunner) evalIfCond() bool {
	extra := r.current.extraData.(*ifElemExtra)
	result := false
	switch extra.condKind {
	case "anagram":
		result = checkAnagram(r.data, []byte(extra.stringArg))
	case "eq":
		result = bytes.Equal(r.data, []byte(extra.stringArg))
	case "substr_count":
		result = bytes.Count(r.data, []byte(extra.stringArg)) == extra.intArg
	case "contains_letter":
		result = bytes.ContainsAny(r.data, extra.stringArg)
	case "contains":
		result = bytes.Contains(r.data, []byte(extra.stringArg))
	case "has_prefix":
		result = bytes.HasPrefix(r.data, []byte(extra.stringArg))
	case "has_suffix":
		result = bytes.HasSuffix(r.data, []byte(extra.stringArg))
	case "last_gt":
		result = r.data[len(r.data)-1] > extra.stringArg[0]
	case "len_even":
		result = len(r.data)%2 == 0
	case "fnv_even":
		result = fnvhash(r.data)%2 == 0
	case "len_eq":
		result = len(r.data) == extra.intArg
	case "len_lt":
		result = len(r.data) < extra.intArg
	case "len_gt":
		result = len(r.data) > extra.intArg
	case "unchanged":
		result = bytes.Equal(r.data, r.input)
	case "fixed_cond":
		result = extra.intArg == 1
	default:
		panic(fmt.Sprintf("unexpected %q elem_if cond kind", extra.condKind))
	}
	return result
}

func (r *schemaRunner) runRepeater() *schemaElem {
	if r.lastCond {
		return r.current.next[0]
	}
	return r.current.next[1]
}

func (r *schemaRunner) runInvRepeater() *schemaElem {
	if !r.lastCond {
		return r.current.next[0]
	}
	return r.current.next[1]
}

func (r *schemaRunner) runIfNot() *schemaElem {
	r.lastCond = !r.evalIfCond()
	if r.lastCond {
		return r.current.next[0]
	}
	return r.current.next[1]
}

func (r *schemaRunner) runIf() *schemaElem {
	r.lastCond = r.evalIfCond()
	if r.lastCond {
		return r.current.next[0]
	}
	return r.current.next[1]
}
