package main

import (
	"os"
	"runtime"

	"github.com/quasilyte/ge"
	"github.com/quasilyte/ge/ui"
	"github.com/quasilyte/gmath"
)

type mainMenuController struct {
	gameState *gameState
	scene     *ge.Scene
}

func newMainMenuController(s *gameState) *mainMenuController {
	return &mainMenuController{gameState: s}
}

func (c *mainMenuController) Init(scene *ge.Scene) {
	c.scene = scene

	scene.Audio().SetGroupVolume(SoundGroupMusic, volumeMultiplier(c.gameState.data.Options.MusicVolumeLevel))
	scene.Audio().SetGroupVolume(SoundGroupEffect, volumeMultiplier(c.gameState.data.Options.EffectsVolumeLevel))
	if c.gameState.data.Options.MusicVolumeLevel != 0 {
		scene.Audio().ContinueMusic(AudioMenuMusic)
	}

	bg := scene.NewSprite(ImagePaperBg)
	bg.Centered = false
	scene.AddGraphics(bg)

	layer := ge.NewShaderLayer()
	layer.Shader = scene.NewShader(ShaderHandwriting)

	offset := gmath.Vec{X: 512, Y: 76}
	l := scene.NewLabel(FontHandwritten)
	l.Text = `What should I do?

  -  Get  back  to  work

  -  Review  the  notes

  -  Adjust  the  options`
	if runtime.GOARCH != "wasm" {
		l.Text += "\n\n  -  Run  a  custom  simulation"
		l.Text += "\n\n  -  Call  it  a  day"
	}
	l.Pos.Offset = offset
	l.ColorScale.SetRGBA(30, 30, 60, 220)
	layer.AddGraphics(l)

	c.initUI(offset)

	scene.AddGraphics(layer)
}

func (c *mainMenuController) initUI(offset gmath.Vec) {
	uiRoot := ui.NewRoot(c.scene.Context(), c.gameState.input)
	uiRoot.ActivationAction = ActionMenuConfirm
	uiRoot.NextInputAction = ActionMenuNext
	uiRoot.PrevInputAction = ActionMenuPrev
	c.scene.AddObject(uiRoot)

	var bgroup buttonGroup

	offset = offset.Add(gmath.Vec{X: 96, Y: 156})

	storyModeButton := uiRoot.NewButton(outlineButtonStyle.Resized(648, 80))
	bgroup.AddButton(storyModeButton)
	storyModeButton.Pos.Offset = offset
	storyModeButton.EventActivated.Connect(nil, func(_ *ui.Button) {
		c.scene.Context().ChangeScene(newChapterSelectController(c.gameState))
	})
	c.scene.AddObject(storyModeButton)

	offset.Y += 166

	reviewNotesButton := uiRoot.NewButton(outlineButtonStyle.Resized(632, 80))
	bgroup.AddButton(reviewNotesButton)
	reviewNotesButton.Pos.Offset = offset
	reviewNotesButton.EventActivated.Connect(nil, func(_ *ui.Button) {
		c.scene.Context().ChangeScene(newManualControler(c.gameState, ""))
	})
	c.scene.AddObject(reviewNotesButton)

	offset.Y += 166

	optionsButton := uiRoot.NewButton(outlineButtonStyle.Resized(668, 80))
	bgroup.AddButton(optionsButton)
	optionsButton.Pos.Offset = offset
	optionsButton.EventActivated.Connect(nil, func(_ *ui.Button) {
		c.scene.Context().ChangeScene(newOptionsController(c.gameState))
	})
	c.scene.AddObject(optionsButton)

	offset.Y += 166

	if runtime.GOARCH != "wasm" {
		customModeButton := uiRoot.NewButton(outlineButtonStyle.Resized(848, 80))
		bgroup.AddButton(customModeButton)
		customModeButton.Pos.Offset = offset
		customModeButton.EventActivated.Connect(nil, func(_ *ui.Button) {
			c.scene.Context().ChangeScene(newCustomLevelSelectController(c.gameState))
		})
		c.scene.AddObject(customModeButton)
		offset.Y += 166

		exitButton := uiRoot.NewButton(outlineButtonStyle.Resized(460, 80))
		bgroup.AddButton(exitButton)
		exitButton.Pos.Offset = offset
		exitButton.EventActivated.Connect(nil, func(_ *ui.Button) {
			os.Exit(0)
		})
		c.scene.AddObject(exitButton)
	}

	bgroup.Connect(uiRoot)
	bgroup.FocusFirst()
}

func (c *mainMenuController) Update(delta float64) {

}
