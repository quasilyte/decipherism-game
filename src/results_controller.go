package main

import (
	"fmt"
	"strings"

	"github.com/quasilyte/ge"
	"github.com/quasilyte/ge/xslices"
	"github.com/quasilyte/gmath"
)

type resultsController struct {
	gameState *gameState
	scene     *ge.Scene
}

func newResultsController(s *gameState) *resultsController {
	return &resultsController{gameState: s}
}

func (c *resultsController) Init(scene *ge.Scene) {
	c.scene = scene

	ctx := scene.Context()
	rect := ge.NewRect(ctx, ctx.WindowWidth, ctx.WindowWidth)
	rect.Centered = false
	rect.FillColorScale.SetRGBA(0x14, 0x18, 0x13, 0xff)
	scene.AddGraphics(rect)

	var textLines []string
	addList := func(title string, lines []string) {
		if len(lines) == 0 {
			return
		}
		textLines = append(textLines, "\n"+title+":\n")
		for _, l := range lines {
			textLines = append(textLines, "      * "+l)
		}
	}

	percent := gmath.Percentage(len(c.gameState.data.CompletedLevels), len(theStoryModeMap.levels))
	textLines = append(textLines, fmt.Sprintf("success! hexagon is now %d%% hacked", percent))

	newContent := calculateContentStatus(c.gameState)

	newFeatures := xslices.Diff(c.gameState.content.techLevelFeatures, newContent.techLevelFeatures)
	if len(newFeatures) > 1 {
		panic("more than one feature is unlocked in on level?")
	}
	if len(newFeatures) != 0 {
		textLines = append(textLines, "\n> unlocked the "+newFeatures[0]+" feature")
	}

	newManualPages := xslices.Diff(c.gameState.content.manualPages, newContent.manualPages)
	addList("> new notes", newManualPages)

	if !xslices.Equal(c.gameState.content.chapters, newContent.chapters) {
		textLines = append(textLines, "\n> new blocks are accessible")
	}
	textLines = append(textLines, "\npress [enter] to continue")

	l := scene.NewLabel(FontLCDNormal)
	l.ColorScale.SetColor(defaultLCDColor)
	l.Pos.Offset = gmath.Vec{X: 64, Y: 64}
	l.Text = strings.Join(textLines, "\n")
	scene.AddGraphics(l)
}

func (c *resultsController) Update(delta float64) {
	if c.gameState.input.ActionIsJustPressed(ActionMenuConfirm) {
		c.scene.Audio().PauseCurrentMusic()
		c.scene.Context().ChangeScene(newLevelSelectController(c.gameState))
	}
}
