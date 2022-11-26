package main

import (
	"github.com/quasilyte/ge"
	"github.com/quasilyte/ge/ui"
)

var invisibleButtonStyle ui.ButtonStyle
var outlineButtonStyle ui.ButtonStyle
var labelButtonStyle ui.ButtonStyle
var optionsButtonStyle ui.ButtonStyle

func init() {
	optionsButtonStyle = ui.DefaultButtonStyle()
	optionsButtonStyle.Font = FontLCDNormal
	optionsButtonStyle.TextColor.SetColor(defaultLCDColor)
	optionsButtonStyle.BorderWidth = 4
	optionsButtonStyle.BackgroundColor.A = 0
	optionsButtonStyle.BorderColor.SetColor(defaultLCDColor)
	optionsButtonStyle.FocusedBackgroundColor = optionsButtonStyle.BackgroundColor
	optionsButtonStyle.FocusedBorderColor = optionsButtonStyle.BorderColor
	optionsButtonStyle.FocusedTextColor.SetColor(ge.RGB(0x2eb43c))

	invisibleButtonStyle = ui.DefaultButtonStyle()
	invisibleButtonStyle.BorderWidth = 0
	invisibleButtonStyle.BackgroundColor.A = 0
	invisibleButtonStyle.FocusedBackgroundColor.A = 0
	invisibleButtonStyle.DisabledBackgroundColor.A = 0
	invisibleButtonStyle.BorderColor.A = 0
	invisibleButtonStyle.FocusedBorderColor.A = 0
	invisibleButtonStyle.DisabledBorderColor.A = 0
	invisibleButtonStyle.DisabledTextColor.A = 0

	outlineButtonStyle = ui.DefaultButtonStyle()
	outlineButtonStyle.BorderWidth = 4
	outlineButtonStyle.BackgroundColor.A = 0
	outlineButtonStyle.FocusedBackgroundColor.A = 0
	outlineButtonStyle.DisabledBackgroundColor.A = 0
	outlineButtonStyle.BorderColor.A = 0
	outlineButtonStyle.FocusedBorderColor.SetRGBA(0x25, 0x25, 0x40, 200)
	outlineButtonStyle.DisabledBorderColor.A = 0
	outlineButtonStyle.DisabledTextColor.A = 0

	labelButtonStyle = ui.DefaultButtonStyle()
	labelButtonStyle.Font = FontHandwritten
	labelButtonStyle.TextColor.SetRGBA(30, 30, 60, 200)
	labelButtonStyle.FocusedTextColor.SetRGBA(60, 60, 120, 255)
	labelButtonStyle.BorderWidth = 0
	labelButtonStyle.BackgroundColor.A = 0
	labelButtonStyle.FocusedBackgroundColor.A = 0
	labelButtonStyle.DisabledBackgroundColor.A = 0
	labelButtonStyle.BorderColor.A = 0
	labelButtonStyle.FocusedBorderColor.A = 0
	labelButtonStyle.DisabledBorderColor.A = 0
	labelButtonStyle.DisabledTextColor.A = 0
}
