package main

import "github.com/quasilyte/ge/ui"

type buttonGroup struct {
	buttons []*ui.Button
}

func (g *buttonGroup) AddButton(b *ui.Button) {
	g.buttons = append(g.buttons, b)
}

func (g *buttonGroup) FocusFirst() {
	if len(g.buttons) == 0 {
		return
	}
	g.buttons[0].SetFocus(true)
}

func (g *buttonGroup) Connect(root *ui.Root) {
	if len(g.buttons) < 2 {
		return
	}
	for i := 0; i < len(g.buttons)-1; i++ {
		root.ConnectInputs(g.buttons[i+0], g.buttons[i+1])
	}
	root.ConnectInputs(g.buttons[len(g.buttons)-1], g.buttons[0])
}
