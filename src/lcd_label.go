package main

import (
	"image/color"

	"github.com/quasilyte/ge"
	"github.com/quasilyte/gmath"
)

type lcdLabel struct {
	pos     gmath.Vec
	text    string
	clr     color.RGBA
	label   *ge.Label
	labelBg *ge.Rect
}

var (
	defaultLCDColor   = ge.RGB(0x2a9535)
	collisionLCDColor = ge.RGB(0xa32828)
	successLCDColor   = ge.RGB(0xcec844)
)

func newLCDLabel(pos gmath.Vec, clr color.RGBA, text string) *lcdLabel {
	return &lcdLabel{pos: pos, text: text, clr: clr}
}

func (l *lcdLabel) Init(scene *ge.Scene) {
	l.labelBg = ge.NewRect(scene.Context(), 328, 64)
	l.labelBg.Centered = false
	l.labelBg.Pos.Base = &l.pos
	l.labelBg.OutlineWidth = 4
	l.labelBg.FillColorScale.SetColor(ge.RGB(0x151917))
	l.labelBg.OutlineColorScale.SetRGBA(0x14, 0x12, 0x1b, 80)
	scene.AddGraphics(l.labelBg)

	l.label = scene.NewLabel(FontLCDSmall)
	l.label.Text = l.text
	l.label.Pos = l.labelBg.Pos
	l.label.ColorScale.SetColor(l.clr)
	l.label.Width = 320
	l.label.Height = 64
	l.label.AlignHorizontal = ge.AlignHorizontalCenter
	l.label.AlignVertical = ge.AlignVerticalCenter
	scene.AddGraphics(l.label)
}

func (l *lcdLabel) IsDisposed() bool { return false }

func (l *lcdLabel) SetColor(clr color.RGBA) {
	l.label.ColorScale.SetColor(clr)
}

func (l *lcdLabel) Update(delta float64) {
	l.label.Text = l.text
}
