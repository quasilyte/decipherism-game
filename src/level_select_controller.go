package main

import (
	"fmt"
	"strings"

	"github.com/quasilyte/ge"
	"github.com/quasilyte/ge/tiled"
	"github.com/quasilyte/ge/ui"
	"github.com/quasilyte/ge/xslices"
	"github.com/quasilyte/gmath"
)

type levelSelectController struct {
	gameState      *gameState
	scene          *ge.Scene
	secretKeywords []string
}

func newLevelSelectController(s *gameState) *levelSelectController {
	return &levelSelectController{gameState: s}
}

func (c *levelSelectController) Init(scene *ge.Scene) {
	c.scene = scene

	bg := scene.NewSprite(ImagePaperBg)
	bg.Centered = false
	scene.AddGraphics(bg)

	layer := ge.NewShaderLayer()
	layer.Shader = scene.NewShader(ShaderHandwriting)

	chapter := c.gameState.chapter
	levelStrings := make([]string, len(chapter.levels))
	for i := range c.gameState.chapter.levels {
		if chapter.IsBonus() {
			levelStrings[i] = fmt.Sprintf("[  ]       Component %d", i+1)
		} else {
			levelStrings[i] = fmt.Sprintf("[  ][  ]   Component %d", i+1)
		}
	}

	var encodedKeyword string
	c.secretKeywords = make([]string, len(chapter.levels))
	if !chapter.IsBonus() {
		tileset, err := tiled.UnmarshalTileset(scene.LoadRaw(RawComponentSchemaTilesetJSON))
		if err != nil {
			panic(err)
		}
		runner := newSchemaRunner()
		inputData := []byte(chapter.keyword)
		for i, levelName := range chapter.levels {
			level := theStoryModeMap.levels[levelName]
			levelData := scene.LoadRaw(level.id)
			schema := decodeSchema(gmath.Vec{}, tileset, levelData)
			completionData := c.gameState.GetLevelCompletionData(levelName)
			if completionData != nil && completionData.SecretKeyword {
				levelStrings[i] += "  (" + strings.ToUpper(string(inputData)) + ")"
			}
			inputData = []byte(runner.Exec(schema, string(inputData)))
			c.secretKeywords[i] = string(inputData)
		}
		encodedKeyword = strings.ToUpper(string(inputData))
	}
	labelText := "Block " + c.gameState.chapter.label + "\n\n" + strings.Join(levelStrings, "\n\n")
	if !chapter.IsBonus() {
		labelText += "\n\n" + encodedKeyword
	}
	offset := gmath.Vec{X: 322, Y: 76}
	l := scene.NewLabel(FontHandwritten)
	l.Text = labelText
	l.Pos.Offset = offset
	l.ColorScale.SetRGBA(30, 30, 60, 220)
	layer.AddGraphics(l)

	c.initUI(offset)

	scene.AddGraphics(layer)
}

func (c *levelSelectController) initDecipherConfig(content contentStatus, config *decipherConfig) {
	config.terminalUpgrades.valueInspector = xslices.Contains(content.techLevelFeatures, "value inspector")
	config.terminalUpgrades.textBuffer = xslices.Contains(content.techLevelFeatures, "text buffer")
	config.terminalUpgrades.branchingInfo = xslices.Contains(content.techLevelFeatures, "branching info")
	config.terminalUpgrades.ioLog = xslices.Contains(content.techLevelFeatures, "i/o logs")
	config.terminalUpgrades.outputPredictor = xslices.Contains(content.techLevelFeatures, "output predictor")
	config.advancedInput = xslices.Contains(content.techLevelFeatures, "advanced input")
}

func (c *levelSelectController) initUI(offset gmath.Vec) {
	uiRoot := ui.NewRoot(c.scene.Context(), c.gameState.input)
	uiRoot.ActivationAction = ActionMenuConfirm
	uiRoot.NextInputAction = ActionMenuNext
	uiRoot.PrevInputAction = ActionMenuPrev
	c.scene.AddObject(uiRoot)

	offset = offset.Add(gmath.Vec{X: 192, Y: 156})

	var bgroup buttonGroup
	chapter := c.gameState.chapter
	for i, levelName := range chapter.levels {
		level := theStoryModeMap.levels[levelName]
		secretKeyword := c.secretKeywords[i]
		b := uiRoot.NewButton(outlineButtonStyle.Resized(454, 80))
		bgroup.AddButton(b)
		b.EventActivated.Connect(nil, func(_ *ui.Button) {
			content := calculateContentStatus(c.gameState)
			c.gameState.level = level
			c.gameState.content = content
			config := decipherConfig{
				secretKeyword: secretKeyword,
				storyMode:     true,
			}
			c.initDecipherConfig(content, &config)
			tileset, err := tiled.UnmarshalTileset(c.scene.LoadRaw(RawComponentSchemaTilesetJSON))
			if err != nil {
				panic(err)
			}
			m, err := tiled.UnmarshalMap(c.scene.LoadRaw(c.gameState.level.id))
			if err != nil {
				panic(err)
			}
			config.levelTemplate = tilemapToTemplate(tileset, m)
			c.scene.Context().ChangeScene(newDecipherController(c.gameState, config))
		})
		completionData := c.gameState.GetLevelCompletionData(levelName)
		if completionData != nil {
			checkmark := c.scene.NewSprite(ImageCompleteMark)
			checkmark.Pos.Offset = offset.Add(gmath.Vec{X: -154, Y: 40})
			c.scene.AddGraphics(checkmark)
			if completionData.SecretKeyword {
				checkmark2 := c.scene.NewSprite(ImageCompleteMark)
				checkmark2.Pos.Offset = offset.Add(gmath.Vec{X: -70, Y: 40})
				c.scene.AddGraphics(checkmark2)
			}
		}
		b.Pos.Offset = offset
		c.scene.AddObject(b)
		if !chapter.IsBonus() {
			arrow := c.scene.NewSprite(ImagePipelineArrow)
			arrow.Pos.Offset = offset.Add(gmath.Vec{X: -236, Y: 112})
			c.scene.AddGraphics(arrow)
		}
		offset.Y += 166
	}
	bgroup.Connect(uiRoot)
	bgroup.FocusFirst()
}

func (c *levelSelectController) Update(delta float64) {
	if c.gameState.input.ActionIsJustPressed(ActionLeave) {
		c.leave()
		return
	}
}

func (c *levelSelectController) leave() {
	c.scene.Context().ChangeScene(newChapterSelectController(c.gameState))
}
