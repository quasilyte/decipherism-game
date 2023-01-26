package main

import (
	"strconv"

	"github.com/quasilyte/ge"
	"github.com/quasilyte/ge/ui"
	"github.com/quasilyte/gmath"
)

type optionsController struct {
	gameState *gameState
	scene     *ge.Scene
}

func newOptionsController(s *gameState) *optionsController {
	return &optionsController{gameState: s}
}

func (c *optionsController) Init(scene *ge.Scene) {
	c.scene = scene

	ctx := scene.Context()
	rect := ge.NewRect(ctx, ctx.WindowWidth, ctx.WindowWidth)
	rect.Centered = false
	rect.FillColorScale.SetRGBA(0x14, 0x18, 0x13, 0xff)
	scene.AddGraphics(rect)

	buttonWidth := 640.0
	offset := gmath.Vec{X: ctx.WindowWidth/2 - buttonWidth/2, Y: 256 - 64}
	uiRoot := ui.NewRoot(ctx, c.gameState.input)
	uiRoot.ActivationAction = ActionMenuConfirm
	uiRoot.NextInputAction = ActionMenuNext
	uiRoot.PrevInputAction = ActionMenuPrev
	scene.AddObject(uiRoot)

	onoffText := func(v bool) string {
		if v {
			return "on"
		}
		return "off"
	}

	options := &c.gameState.data.Options

	var bgroup buttonGroup

	{
		musicToggle := uiRoot.NewButton(optionsButtonStyle.Resized(buttonWidth, 80))
		bgroup.AddButton(musicToggle)
		var musicToggleValue gmath.Slider
		musicToggleValue.SetBounds(0, 5)
		musicToggleValue.TrySetValue(options.MusicVolumeLevel)
		musicToggle.Text = "music volume: " + strconv.Itoa(options.MusicVolumeLevel)
		musicToggle.Pos.Offset = offset
		musicToggle.EventActivated.Connect(nil, func(_ *ui.Button) {
			musicToggleValue.Inc()
			options.MusicVolumeLevel = musicToggleValue.Value()
			musicToggle.Text = "music volume: " + strconv.Itoa(options.MusicVolumeLevel)
			if options.MusicVolumeLevel != 0 {
				scene.Audio().SetGroupVolume(SoundGroupMusic, volumeMultiplier(options.MusicVolumeLevel))
				scene.Audio().PauseCurrentMusic()
				scene.Audio().PlayMusic(AudioMenuMusic)
			} else {
				scene.Audio().PauseCurrentMusic()
			}
		})
		scene.AddObject(musicToggle)
	}
	offset.Y += 128

	{
		effectToggle := uiRoot.NewButton(optionsButtonStyle.Resized(buttonWidth, 80))
		bgroup.AddButton(effectToggle)
		var effectToggleValue gmath.Slider
		effectToggleValue.SetBounds(0, 5)
		effectToggleValue.TrySetValue(options.EffectsVolumeLevel)
		effectToggle.Text = "effects volume: " + strconv.Itoa(options.EffectsVolumeLevel)
		effectToggle.Pos.Offset = offset
		effectToggle.EventActivated.Connect(nil, func(_ *ui.Button) {
			effectToggleValue.Inc()
			options.EffectsVolumeLevel = effectToggleValue.Value()
			effectToggle.Text = "effects volume: " + strconv.Itoa(options.EffectsVolumeLevel)
			if options.EffectsVolumeLevel != 0 {
				scene.Audio().SetGroupVolume(SoundGroupEffect, volumeMultiplier(options.EffectsVolumeLevel))
				scene.Audio().PlaySound(AudioSecretUnlocked)
			}
		})
		scene.AddObject(effectToggle)
	}
	offset.Y += 128

	shaderToggle := uiRoot.NewButton(optionsButtonStyle.Resized(buttonWidth, 80))
	bgroup.AddButton(shaderToggle)
	shaderToggle.Text = "crt shaders: " + onoffText(options.CrtShader)
	shaderToggle.Pos.Offset = offset
	shaderToggle.EventActivated.Connect(nil, func(_ *ui.Button) {
		options.CrtShader = !options.CrtShader
		shaderToggle.Text = "crt shaders: " + onoffText(options.CrtShader)
	})
	scene.AddObject(shaderToggle)
	offset.Y += 128

	backButton := uiRoot.NewButton(optionsButtonStyle.Resized(buttonWidth, 80))
	bgroup.AddButton(backButton)
	backButton.Text = "back"
	backButton.Pos.Offset = offset
	backButton.EventActivated.Connect(nil, func(_ *ui.Button) {
		c.leave()
	})
	scene.AddObject(backButton)
	offset.Y += 128

	bgroup.Connect(uiRoot)
	bgroup.FocusFirst()

	version := scene.NewLabel(FontLCDSmall)
	version.ColorScale.SetColor(defaultLCDColor)
	version.Text = "build " + buildVersion
	version.Pos.Offset = offset
	version.Width = buttonWidth
	version.Height = 80
	version.AlignHorizontal = ge.AlignHorizontalCenter
	version.AlignVertical = ge.AlignVerticalCenter
	scene.AddGraphics(version)
}

func (c *optionsController) leave() {
	c.scene.Context().SaveGameData("save", *c.gameState.data)
	c.scene.Context().ChangeScene(newMainMenuController(c.gameState))
}

func (c *optionsController) Update(delta float64) {
	if c.gameState.input.ActionIsJustPressed(ActionLeave) {
		c.leave()
	}
}
