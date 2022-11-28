package main

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/quasilyte/ge"
	"github.com/quasilyte/ge/gesignal"
	"github.com/quasilyte/ge/input"
	"github.com/quasilyte/gmath"
)

const maxInputLen = 10

type componentInput struct {
	pos              gmath.Vec
	input            *input.Handler
	pressedRunes     []rune
	text             []byte
	lcdLabel         *lcdLabel
	cursorPos        int
	cursor           *ge.Label
	cursorBlinkDelay float64
	advancedOps      bool

	EventOnTextChanged gesignal.Event[gesignal.Void]
}

func newComponentInput(h *input.Handler, pos gmath.Vec, advancedOps bool) *componentInput {
	return &componentInput{input: h, pos: pos, advancedOps: advancedOps}
}

func (i *componentInput) SetText(s string) {
	i.text = []byte(s)
	i.cursorPos = -1
	i.onTextChanged()
}

func (i *componentInput) Init(scene *ge.Scene) {
	i.text = []byte("abc")
	textColor := ge.RGB(0x2eb43c)
	i.lcdLabel = newLCDLabel(i.pos, textColor, string(i.text))
	scene.AddObject(i.lcdLabel)

	i.cursor = scene.NewLabel(FontLCDSmall)
	i.cursor.Visible = false
	i.cursor.Height = 64
	i.cursor.Pos.Base = &i.pos
	i.cursor.Pos.Offset.Y += 6
	i.cursor.Text = "_"
	i.cursor.AlignVertical = ge.AlignVerticalCenter
	i.cursor.ColorScale.SetColor(textColor)
	scene.AddGraphics(i.cursor)

	i.cursorPos = -1
	i.cursorBlinkDelay = 0.2

	i.onTextChanged()
}

func (i *componentInput) IsDisposed() bool { return false }

func (i *componentInput) Update(delta float64) {
	i.cursorBlinkDelay -= delta
	if i.cursorBlinkDelay <= 0 {
		i.cursorBlinkDelay = 1
		i.cursor.Visible = !i.cursor.Visible
	}

	if i.advancedOps && len(i.text) != 0 {
		if i.cursorPos != -1 && i.input.ActionIsJustPressed(ActionCharInc) {
			i.text[i.cursorPos] = incChar(i.text[i.cursorPos])
			i.onTextChanged()
			return
		}
		if i.cursorPos != -1 && i.input.ActionIsJustPressed(ActionCharDec) {
			i.text[i.cursorPos] = decChar(i.text[i.cursorPos])
			i.onTextChanged()
			return
		}
		if i.input.ActionIsJustPressed(ActionRotateLeft) {
			rotateCharsLeft(i.text)
			i.onTextChanged()
			return
		}
		if i.input.ActionIsJustPressed(ActionRotateRight) {
			rotateCharsRight(i.text)
			i.onTextChanged()
			return
		}
	}

	if len(i.text) != 0 && (i.cursorPos > 0 || i.cursorPos == -1) && i.input.ActionIsJustPressed(ActionCursorLeft) {
		if i.cursorPos == -1 {
			i.cursorPos = len(i.text) - 1
		} else {
			i.cursorPos--
		}
		i.cursorBlinkDelay = 0.5
		i.cursor.Visible = true
		i.onCursorChanged()
	}
	if i.cursorPos != -1 && i.input.ActionIsJustPressed(ActionCursorRight) {
		i.cursorPos++
		i.cursorBlinkDelay = 0.5
		i.cursor.Visible = true
		i.onCursorChanged()
		if i.cursorPos >= len(i.text) {
			i.cursorPos = -1
		}
	}
	if len(i.text) != 0 && i.input.ActionIsJustPressed(ActionRemoveCurrentChar) {
		if i.cursorPos != -1 {
			deletePos := i.cursorPos
			i.text = append(i.text[:deletePos], i.text[deletePos+1:]...)
			i.onTextChanged()
			if len(i.text) == 0 || i.cursorPos == len(i.text) {
				i.cursorPos = -1
			}
			return
		}
	}
	if len(i.text) != 0 && i.input.ActionIsJustPressed(ActionRemovePrevChar) {
		if i.cursorPos != 0 {
			if i.cursorPos == -1 {
				// Remove the last letter.
				i.text = i.text[:len(i.text)-1]
			} else {
				// Remove the letter behind the cursor.
				deletePos := i.cursorPos - 1
				i.text = append(i.text[:deletePos], i.text[deletePos+1:]...)
				i.cursorPos--
			}
			i.onTextChanged()
			return
		}
	}
	if len(i.text) < maxInputLen && !ebiten.IsKeyPressed(ebiten.KeyControl) {
		i.pressedRunes = ebiten.AppendInputChars(i.pressedRunes[:0])
		if len(i.pressedRunes) != 0 {
			changed := false
			for _, r := range i.pressedRunes {
				if r >= 'a' && r <= 'z' {
					if i.cursorPos == -1 {
						// Append the char to the end of the text.
						i.text = append(i.text, byte(r))
					} else {
						// Prepend the char.
						insertPos := i.cursorPos
						i.text = append(i.text[:insertPos], append([]byte{byte(r)}, i.text[insertPos:]...)...)
						i.cursorPos++
					}
					changed = true
				}
			}
			if changed {
				i.onTextChanged()
			}
		}
	}
}

func (i *componentInput) onTextChanged() {
	i.lcdLabel.text = string(i.text)
	i.onCursorChanged()
	i.EventOnTextChanged.Emit(gesignal.Void{})
}

func (i *componentInput) onCursorChanged() {
	cursorPos := i.cursorPos
	if cursorPos < 0 {
		cursorPos = len(i.text)
	}
	const letterWidth = 26
	textWidth := letterWidth * len(i.text)
	firstLetterOffset := 160 - textWidth/2
	i.cursor.Pos.Offset.X = float64(firstLetterOffset + cursorPos*letterWidth + 2)
}
