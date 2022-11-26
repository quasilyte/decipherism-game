package main

import (
	"github.com/quasilyte/ge"
	"github.com/quasilyte/ge/ui"
	"github.com/quasilyte/ge/xslices"
	"github.com/quasilyte/gmath"
)

type chapterSelectController struct {
	scene     *ge.Scene
	gameState *gameState
}

func newChapterSelectController(s *gameState) *chapterSelectController {
	return &chapterSelectController{gameState: s}
}

func (c *chapterSelectController) Init(scene *ge.Scene) {
	c.scene = scene

	bg := scene.NewSprite(ImagePaperBg)
	bg.Centered = false
	bg.FlipHorizontal = true
	scene.AddGraphics(bg)

	uiRoot := ui.NewRoot(c.scene.Context(), c.gameState.input)
	uiRoot.ActivationAction = ActionMenuConfirm
	c.scene.AddObject(uiRoot)

	nodeOffset := gmath.Vec{X: 112 + (42 * 4), Y: 36 + (42.5 * 3)}
	content := calculateContentStatus(c.gameState)
	for i := range theStoryModeMap.chapters {
		chapter := &theStoryModeMap.chapters[i]
		chapterStatus := c.gameState.GetChapterCompletionData(chapter)
		available := xslices.Contains(content.chapters, chapter.name)
		pos := nodeOffset.Add(gmath.Vec{X: 42 * (chapter.gridPos.X * 7), Y: 42 * (chapter.gridPos.Y * 7)})
		if !available {
			continue
		}
		n := newChapterNode(pos, chapter.name, available, chapterStatus.fullyCompleted)
		scene.AddObject(n)

		button := uiRoot.NewButton(labelButtonStyle.Resized(42*4, 42.5*4))
		button.Pos.Offset = pos
		button.Text = chapter.label
		button.EventActivated.Connect(nil, func(_ *ui.Button) {
			c.gameState.chapter = chapter
			c.scene.Context().ChangeScene(newLevelSelectController(c.gameState))
		})
		c.scene.AddObject(button)
	}

	// {
	// 	offset := gmath.Vec{X: 100, Y: 30}
	// 	pos := offset.Add(gmath.Vec{X: 43 * 4, Y: 43 * 3})
	// 	n := newChapterNode(pos, "x", true, false)
	// 	scene.AddObject(n)
	// }
	// {
	// 	offset := gmath.Vec{X: 100, Y: 30}
	// 	pos := offset.Add(gmath.Vec{X: 42.5 * ((7 * 3) + 4), Y: 42.5 * 3})
	// 	n := newChapterNode(pos, "x", true, true)
	// 	scene.AddObject(n)
	// }

	outline := scene.NewSprite(ImageChapterSelectOutline)
	outline.Centered = false
	scene.AddGraphics(outline)
}

func (c *chapterSelectController) Update(delta float64) {
	if c.gameState.input.ActionIsJustPressed(ActionLeave) {
		c.leave()
		return
	}
}

func (c *chapterSelectController) leave() {
	c.scene.Context().ChangeScene(newMainMenuController(c.gameState))
}
