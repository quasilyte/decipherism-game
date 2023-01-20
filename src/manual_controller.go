package main

import (
	"fmt"
	"strings"

	"github.com/quasilyte/ge"
	"github.com/quasilyte/ge/xslices"
	"github.com/quasilyte/gmath"
)

type manualController struct {
	scene     *ge.Scene
	gameState *gameState

	pageSlider     gmath.Slider
	pagesAvailable []gameManualPage

	bg           *ge.Sprite
	illustration *ge.Sprite
	text         *ge.Label

	open string
}

func newManualControler(s *gameState, open string) *manualController {
	return &manualController{
		gameState: s,
		open:      open,
	}
}

func (c *manualController) Init(scene *ge.Scene) {
	c.scene = scene

	content := calculateContentStatus(c.gameState)
	for _, p := range theGameManual.pages {
		if xslices.Contains(content.manualPages, p.title) {
			c.pagesAvailable = append(c.pagesAvailable, p)
		}
	}
	c.pageSlider.SetBounds(0, len(c.pagesAvailable)-1)

	c.bg = scene.NewSprite(ImagePaperBg)
	c.bg.Centered = false
	scene.AddGraphics(c.bg)

	layer := ge.NewShaderLayer()
	layer.Shader = scene.NewShader(ShaderHandwriting)

	c.text = scene.NewLabel(FontHandwritten)
	c.text.Pos.Offset = gmath.Vec{X: 176, Y: 72}
	c.text.ColorScale.SetRGBA(30, 30, 60, 220)
	layer.AddGraphics(c.text)

	scene.AddGraphics(layer)

	c.illustration = ge.NewSprite(scene.Context())
	c.illustration.Centered = false
	c.illustration.Visible = false
	c.illustration.Pos.Offset.X = 1370
	scene.AddGraphics(c.illustration)

	pageIndex := 0
	if c.open != "" {
		for i, p := range c.pagesAvailable {
			if p.title == c.open {
				pageIndex = i
				break
			}
		}
	}
	c.pageSlider.TrySetValue(pageIndex)
	c.flipPage(c.pagesAvailable[pageIndex])
}

func (c *manualController) flipPage(p gameManualPage) {
	flipHorizontal := c.bg.FlipHorizontal
	flipVertical := c.bg.FlipVertical
	for {
		c.bg.FlipHorizontal = c.scene.Rand().Bool()
		c.bg.FlipVertical = c.scene.Rand().Bool()
		if c.bg.FlipHorizontal != flipHorizontal || c.bg.FlipVertical != flipVertical {
			break
		}
	}

	if p.image != ImageNone {
		c.illustration.SetImage(c.scene.LoadImage(p.image))
		c.illustration.Visible = true
	} else {
		c.illustration.Visible = false
	}

	var buf strings.Builder
	if p.title != "" {
		buf.WriteString(p.title)
		buf.WriteRune('\n')
		buf.WriteRune('\n')
	}
	text := p.text
	if p.params != nil {
		text = fmt.Sprintf(text, p.params(c.gameState.data)...)
	}
	lines := strings.Split(text, "\n")
	for _, l := range lines {
		l = strings.TrimSpace(l)
		if l == "" {
			continue
		}
		if strings.HasPrefix(l, "> ") {
			l = strings.TrimPrefix(l, "> ")
			buf.WriteString("    ")
		}
		if l == "\\n" {
			buf.WriteByte('\n')
			continue
		}
		buf.WriteString(l)
		buf.WriteByte('\n')
	}
	c.text.Text = buf.String()
}

func (c *manualController) Update(delta float64) {
	if c.gameState.input.ActionIsJustPressed(ActionLeave) {
		c.leave()
		return
	}
	if c.gameState.input.ActionIsJustPressed(ActionMenuNextPage) || c.gameState.input.ActionIsJustPressed(ActionMenuConfirm) {
		c.nextPage()
		return
	}
	if c.gameState.input.ActionIsJustPressed(ActionMenuPrevPage) {
		c.prevPage()
		return
	}
}

func (c *manualController) nextPage() {
	c.pageSlider.Inc()
	c.flipPage(c.pagesAvailable[c.pageSlider.Value()])
}

func (c *manualController) prevPage() {
	c.pageSlider.Dec()
	c.flipPage(c.pagesAvailable[c.pageSlider.Value()])
}

func (c *manualController) leave() {
	c.scene.Context().ChangeScene(newMainMenuController(c.gameState))
}
