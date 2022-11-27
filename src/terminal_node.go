package main

import (
	"fmt"
	"strings"

	"github.com/quasilyte/ge"
	"github.com/quasilyte/gmath"
)

type terminalNode struct {
	config terminalConfig

	text *ge.Label

	offset gmath.Vec

	textBuffer string
}

type statusInfo struct {
	value           string
	ioLogs          []string
	predictedOutput string
}

type terminalConfig struct {
	username string

	branchHints []string

	upgrades terminalUpgrades
}

type terminalUpgrades struct {
	ioLog           bool
	branchingInfo   bool
	textBuffer      bool
	valueInspector  bool
	outputPredictor bool
}

func newTerminalNode(config terminalConfig) *terminalNode {
	return &terminalNode{
		config: config,
		offset: gmath.Vec{X: 96, Y: 96},
	}
}

func (n *terminalNode) Init(scene *ge.Scene) {
	n.text = scene.NewLabel(FontLCDSmall)
	n.text.Pos.Offset = n.offset.Add(gmath.Vec{X: 16, Y: 16})
	n.text.ColorScale.SetColor(defaultLCDColor)
	n.text.Visible = false
	scene.AddGraphics(n.text)
}

func (n *terminalNode) SetVisible(visible bool) {
	n.text.Visible = visible
}

func (n *terminalNode) IsDisposed() bool { return false }

func (n *terminalNode) Update(delta float64) {}

func (n *terminalNode) GetBufferText() (string, bool) {
	if n.config.upgrades.textBuffer {
		return n.textBuffer, true
	}
	return "", false
}

func (n *terminalNode) UpdateInfo(info statusInfo) {
	greeting := n.config.username + "@decodeos $ statusdump"

	textlines := []string{
		greeting,
	}

	if n.config.upgrades.valueInspector {
		textlines = append(textlines, "", fmt.Sprintf("current value: %s", info.value))
	} else {
		textlines = append(textlines, "", "current value: unavailable")
	}

	if n.config.upgrades.branchingInfo {
		if len(n.config.branchHints) == 0 {
			textlines = append(textlines, "", "branching info: no branches")
		} else {
			textlines = append(textlines, "")
			for i, b := range n.config.branchHints {
				if i == 0 {
					textlines = append(textlines, "branching info: * "+b)
				} else {
					textlines = append(textlines, "                * "+b)
				}
			}
		}
	} else {
		textlines = append(textlines, "", "branching info: unavailable")
	}

	if n.config.upgrades.ioLog {
		if len(info.ioLogs) == 0 {
			textlines = append(textlines, "", "i/o logs: no data")
		} else {
			textlines = append(textlines, "")
			for i, l := range info.ioLogs {
				if i == 0 {
					textlines = append(textlines, "i/o logs: * "+l)
				} else {
					textlines = append(textlines, "          * "+l)
				}
			}
		}
	} else {
		textlines = append(textlines, "", "i/o logs: unavailable")
	}

	if n.config.upgrades.textBuffer {
		textBuffer := n.textBuffer
		if textBuffer == "" {
			textBuffer = "<empty>"
		}
		textlines = append(textlines, "", "text buffer: "+textBuffer)
	} else {
		textlines = append(textlines, "", "text buffer: unavailable")
	}

	if n.config.upgrades.outputPredictor {
		textlines = append(textlines, "", "output prediction: "+info.predictedOutput)
	} else {
		textlines = append(textlines, "", "output prediction: unavailable")
	}

	n.text.Text = strings.Join(textlines, "\n")
}
