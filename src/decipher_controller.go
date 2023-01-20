package main

import (
	"time"

	"github.com/quasilyte/ge"
	"github.com/quasilyte/ge/gesignal"
	"github.com/quasilyte/ge/input"
	"github.com/quasilyte/ge/ui"
	"github.com/quasilyte/ge/xslices"
	"github.com/quasilyte/gmath"
)

type decipherController struct {
	gameState *gameState
	scene     *ge.Scene

	config decipherConfig

	startTime     time.Time
	secretDecoded bool

	keywords       []string
	keywordToggles []*ge.Sprite
	keywordState   []bool
	numDecoded     int

	schema       *componentSchema
	schemaNodes  []*schemaElemNode
	stickerNodes []*stickerNode
	ioLogs       []string

	schemaBg     *ge.Sprite
	terminalBg   *ge.Sprite
	terminalNode *terminalNode

	signalNode      *signalNode
	signalNodeSpeed float64
	simulationInput string
	runner          *schemaRunner
	termRunner      *schemaRunner

	paused bool

	componentInput *componentInput
	outputLabel    *lcdLabel
	statusLabel    *lcdLabel
}

type decipherConfig struct {
	secretKeyword    string
	terminalUpgrades terminalUpgrades
	advancedInput    bool
	storyMode        bool
	levelTemplate    *schemaTemplate
}

func newDecipherController(s *gameState, config decipherConfig) *decipherController {
	return &decipherController{
		gameState:  s,
		runner:     newSchemaRunner(),
		termRunner: newSchemaRunner(),
		config:     config,
	}
}

func (c *decipherController) newLabel(pos gmath.Vec, text string) *ge.Label {
	l := c.scene.NewLabel(FontSmall)
	l.Text = text
	l.Pos.Offset = pos
	l.Height = 64
	l.ColorScale.SetRGBA(0xe7, 0xe7, 0xe6, 220)
	l.AlignVertical = ge.AlignVerticalCenter
	return l
}

func (c *decipherController) Init(scene *ge.Scene) {
	c.scene = scene

	c.startTime = time.Now()

	if c.gameState.data.Options.MusicVolumeLevel != 0 {
		scene.Audio().PauseCurrentMusic()
		scene.Audio().PlayMusic(AudioDecipherMusic)
	}

	bg := scene.NewSprite(ImageDecipherBg)
	bg.Centered = false
	scene.AddGraphics(bg)

	c.schemaBg = scene.NewSprite(ImageSchemaBg)
	c.schemaBg.Centered = false
	c.schemaBg.Pos.Offset.X = 96
	c.schemaBg.Pos.Offset.Y = 96
	scene.AddGraphics(c.schemaBg)

	c.terminalBg = scene.NewSprite(ImageTerminalBg)
	c.terminalBg.Centered = false
	c.terminalBg.Pos.Offset.X = 96
	c.terminalBg.Pos.Offset.Y = 96
	c.terminalBg.Visible = false
	scene.AddGraphics(c.terminalBg)

	outputTitle := c.newLabel(gmath.Vec{X: 1568 + 64 + 16, Y: 96 * 2}, "OUTPUT")
	scene.AddGraphics(outputTitle)

	inputTitle := c.newLabel(gmath.Vec{X: 1568 + 64 + 16, Y: 96}, "INPUT")
	scene.AddGraphics(inputTitle)

	c.outputLabel = newLCDLabel(gmath.Vec{X: 1568 - 256 - 16, Y: 96 * 2}, defaultLCDColor, "?")
	scene.AddObject(c.outputLabel)

	statusTitle := c.newLabel(gmath.Vec{X: 1568 + 64 + 16, Y: 96 * 3}, "STATUS")
	scene.AddGraphics(statusTitle)

	c.statusLabel = newLCDLabel(gmath.Vec{X: 1568 - 256 - 16, Y: 96 * 3}, defaultLCDColor, "ready")
	scene.AddObject(c.statusLabel)

	uiRoot := ui.NewRoot(scene.Context(), c.gameState.input)
	uiRoot.ActivationAction = ActionButton

	slowTitleLabel := c.newLabel(gmath.Vec{X: 96, Y: 916 + 32}, "SLOW")
	slowTitleLabel.Width = 160
	scene.AddGraphics(slowTitleLabel)
	fastTitleLabel := c.newLabel(gmath.Vec{X: 96 + 160 + 138, Y: 916 + 32}, "FAST")
	fastTitleLabel.Width = 160
	scene.AddGraphics(fastTitleLabel)
	speedDial := newDialButton(uiRoot, gmath.Vec{X: 96 + 160, Y: 896 + 32}, 5)
	speedDial.state = 2
	scene.AddObject(speedDial)
	c.signalNodeSpeed = 120 * float64(speedDial.state+1)
	speedDial.EventActivated.Connect(nil, func(speedLevel int) {
		c.signalNodeSpeed = 120 * float64(speedLevel+1)
		if c.signalNode != nil {
			c.signalNode.speed = c.signalNodeSpeed
		}
	})

	c.initComponentSchema(c.schemaBg.Pos.Offset)

	for _, e := range c.schema.elems {
		node := newSchemaElemNode(e, c.gameState.data.Options.CrtShader)
		scene.AddObject(node)
		c.schemaNodes = append(c.schemaNodes, node)
	}

	for _, h := range c.config.levelTemplate.hints {
		hintNode := newStickerNode(h.pos, h.text)
		scene.AddObject(hintNode)
		c.stickerNodes = append(c.stickerNodes, hintNode)
	}

	c.componentInput = newComponentInput(c.gameState.input, gmath.Vec{X: 1568 - 256 - 16, Y: 96}, c.config.advancedInput)
	c.componentInput.EventOnTextChanged.Connect(nil, c.onInputTextChanged)
	scene.AddObject(c.componentInput)

	keywordsTitle := c.newLabel(gmath.Vec{X: 1568 - 256 - 16, Y: 96*4 + 32}, "ENCODED KEYWORDS")
	scene.AddGraphics(keywordsTitle)

	c.schema.encodedKeywords = make([]string, len(c.keywords))
	for i, keyword := range c.keywords {
		c.schema.encodedKeywords[i] = c.encodeKeyword(keyword)
	}

	offset := gmath.Vec{X: 1568 - 256 - 16, Y: 96*5 + 32}
	for _, keyword := range c.schema.encodedKeywords {
		l := newLCDLabel(offset, defaultLCDColor, keyword)
		scene.AddObject(l)
		offset.Y += 96

		toggle := scene.NewSprite(ImageOnOffButton)
		toggle.Centered = false
		toggle.Pos.Offset = offset.Add(gmath.Vec{X: 256 + 80, Y: -88})
		scene.AddGraphics(toggle)
		c.keywordToggles = append(c.keywordToggles, toggle)
		c.keywordState = append(c.keywordState, false)
	}

	var branches []string
	for _, e := range c.schema.elems {
		info, ok := e.extraData.(*ifElemExtra)
		if !ok {
			continue
		}
		branches = append(branches, info.condKind)
	}
	gmath.Shuffle(scene.Rand(), branches)
	if len(branches) > 3 {
		branches = branches[:3]
	}

	c.terminalNode = newTerminalNode(terminalConfig{
		username:    "quasilyte",
		branchHints: branches,
		upgrades:    c.config.terminalUpgrades,
	})
	scene.AddObject(c.terminalNode)
	c.terminalNode.UpdateInfo(statusInfo{})
}

func (c *decipherController) encodeKeyword(k string) string {
	return c.termRunner.Exec(c.schema, k)
}

func (c *decipherController) clearLevel() {
	if !c.config.storyMode {
		c.scene.Context().ChangeScene(newMainMenuController(c.gameState))
		return
	}

	completionData := xslices.Find(c.gameState.data.CompletedLevels, func(d *completedLevelData) bool {
		return d.Name == c.gameState.level.name
	})
	if completionData == nil {
		c.gameState.data.CompletedLevels = append(c.gameState.data.CompletedLevels, completedLevelData{
			Name:          c.gameState.level.name,
			SecretKeyword: c.secretDecoded,
		})
		c.gameState.data.CompletionTime += time.Since(c.startTime)

	} else if c.secretDecoded {
		completionData.SecretKeyword = true
	}

	c.gameState.data.SolvedCondTransform = c.gameState.data.SolvedCondTransform || c.schema.hasCondTransform
	c.gameState.data.SolvedPolygraphic = c.gameState.data.SolvedPolygraphic || c.schema.hasPolygraphic
	c.gameState.data.SolvedAtbash = c.gameState.data.SolvedAtbash || c.schema.hasAtbash
	c.gameState.data.SolvedRot13 = c.gameState.data.SolvedRot13 || c.schema.hasRot13
	c.gameState.data.SolvedIncDec = c.gameState.data.SolvedIncDec || c.schema.hasIncDec
	c.gameState.data.SolvedShift = c.gameState.data.SolvedShift || c.schema.hasShift
	c.gameState.data.SolvedNegation = c.gameState.data.SolvedNegation || c.schema.hasNegation

	c.scene.Context().SaveGameData("save", *c.gameState.data)
	c.scene.Context().ChangeScene(newResultsController(c.gameState))
}

func (c *decipherController) onProgramCompleted(output string) {
	if c.signalNode != nil {
		c.signalNode.Dispose()
		c.signalNode = nil
	}
	c.outputLabel.text = output
	c.statusLabel.text = "READY"
	c.outputLabel.SetColor(defaultLCDColor)

	if len(c.ioLogs) < 3 {
		c.ioLogs = append(c.ioLogs, c.simulationInput+" -> "+output)
	} else {
		c.ioLogs[0] = c.ioLogs[1]
		c.ioLogs[1] = c.ioLogs[2]
		c.ioLogs[2] = c.simulationInput + " -> " + output
	}

	if c.simulationInput == c.config.secretKeyword {
		c.scene.Context().Audio.PlaySound(AudioSecretUnlocked)
		c.secretDecoded = true
		completionData := xslices.Find(c.gameState.data.CompletedLevels, func(d *completedLevelData) bool {
			return d.Name == c.gameState.level.name
		})
		if completionData != nil {
			completionData.SecretKeyword = true
			c.outputLabel.SetColor(successLCDColor)
			c.scene.Context().SaveGameData("save", *c.gameState.data)
		}
		return
	}

	i := xslices.Index(c.keywords, c.simulationInput)
	if i != -1 && !c.keywordState[i] {
		c.keywordState[i] = true
		c.keywordToggles[i].FrameOffset.X = 96
		c.outputLabel.SetColor(successLCDColor)
		c.scene.Context().Audio.PlaySound(AudioDecodingSuccess)
		c.numDecoded++
		if c.numDecoded == len(c.keywordToggles) {
			c.scene.DelayedCall(1.0, c.clearLevel)
		}
		return
	}

	if i := xslices.Index(c.schema.encodedKeywords, output); i != -1 {
		c.gameState.data.SawCollision = true
		c.outputLabel.SetColor(collisionLCDColor)
		c.scene.Context().Audio.PlaySound(AudioCollision)
		return
	}
}

func (c *decipherController) prepareNextStep(sig *signalNode) {
	dst, hasMore := c.runner.RunStep()
	if !hasMore {
		c.onProgramCompleted(string(c.runner.data))
		return
	}
	sig.dst = dst
}

func (c *decipherController) nextStep(sig *signalNode) {
	var clr ge.ColorScale
	clr.SetRGBA(0xd1, 0xc2, 0x73, 255)
	c.scene.AddObject(newPingEffectNode(sig.pos, clr))
	if c.paused {
		return
	}
	c.prepareNextStep(sig)
}

func (c *decipherController) pauseProgram(paused bool) {
	c.paused = paused
	if !c.paused {
		c.statusLabel.text = "RUNNING"
		c.prepareNextStep(c.signalNode)
	} else {
		c.statusLabel.text = "PAUSED"
	}
}

func (c *decipherController) isInTerminalMode() bool {
	return c.terminalBg.Visible
}

func (c *decipherController) makeStatusInfo() statusInfo {
	value := "?"
	if c.signalNode != nil {
		value = string(c.runner.data)
	}

	var predictedOutput string
	if len(c.componentInput.text) != 0 {
		b := []byte(c.encodeKeyword(string(c.componentInput.text)))
		maskedLetter := fnvhash(b) % uint64(len(b))
		b[maskedLetter] = '?'
		predictedOutput = string(b)
	} else {
		predictedOutput = "?"
	}

	return statusInfo{
		value:           value,
		ioLogs:          c.ioLogs,
		predictedOutput: string(predictedOutput),
	}
}

func (c *decipherController) onInputTextChanged(gesignal.Void) {
	if !c.isInTerminalMode() {
		return
	}
	c.terminalNode.UpdateInfo(c.makeStatusInfo())
}

func (c *decipherController) swapMode() {
	if c.isInTerminalMode() {
		c.terminalBg.Visible = false
		c.terminalNode.SetVisible(false)
		c.schemaBg.Visible = true
		for _, e := range c.schemaNodes {
			e.sprite.SetAlpha(1)
		}
		for _, hint := range c.stickerNodes {
			hint.sprite.SetAlpha(1)
		}
		if c.signalNode != nil {
			c.signalNode.sprite.SetAlpha(1)
		}
	} else {
		c.terminalBg.Visible = true
		c.terminalNode.UpdateInfo(c.makeStatusInfo())
		c.terminalNode.SetVisible(true)
		c.schemaBg.Visible = false
		for _, e := range c.schemaNodes {
			e.sprite.SetAlpha(0.2)
		}
		for _, hint := range c.stickerNodes {
			hint.sprite.SetAlpha(0.3)
		}
		if c.signalNode != nil {
			c.signalNode.sprite.SetAlpha(0.5)
		}
	}
}

func (c *decipherController) leave() {
	c.scene.Audio().PauseCurrentMusic()
	if c.config.storyMode {
		c.scene.Context().ChangeScene(newLevelSelectController(c.gameState))
	} else {
		c.scene.Context().ChangeScene(newMainMenuController(c.gameState))
	}
}

func (c *decipherController) Update(delta float64) {
	if c.gameState.input.ActionIsJustPressed(ActionLeave) {
		c.leave()
		return
	}

	if c.gameState.input.ActionIsJustPressed(ActionClearStage) {
		c.gameState.data.UsedCheats = true
		c.clearLevel()
		return
	}

	if !c.isInTerminalMode() {
		for i := 0; i < 4; i++ {
			action := ActionPasteKeyword1 + input.Action(i)
			if c.gameState.input.ActionIsJustPressed(action) {
				c.gameState.data.UsedHiddenKeybinds = true
				if len(c.schema.encodedKeywords) > i {
					c.componentInput.SetText(c.schema.encodedKeywords[i])
				}
				return
			}
		}
	}

	if !c.isInTerminalMode() && c.gameState.input.ActionIsJustPressed(ActionPauseProgram) && c.signalNode != nil {
		c.pauseProgram(!c.paused)
		return
	}

	if (c.paused || c.signalNode == nil) && c.gameState.input.ActionIsJustPressed(ActionModeSwap) {
		c.swapMode()
		return
	}

	if len(c.componentInput.text) != 0 && c.gameState.input.ActionIsJustPressed(ActionBufferCut) {
		c.terminalNode.textBuffer = string(c.componentInput.text)
		if c.isInTerminalMode() {
			c.terminalNode.UpdateInfo(c.makeStatusInfo())
		}
		c.componentInput.SetText("")
		return
	}
	if len(c.componentInput.text) != 0 && c.gameState.input.ActionIsJustPressed(ActionBufferCopy) {
		c.terminalNode.textBuffer = string(c.componentInput.text)
		if c.isInTerminalMode() {
			c.terminalNode.UpdateInfo(c.makeStatusInfo())
		}
		return
	}
	if c.gameState.input.ActionIsJustPressed(ActionBufferPaste) {
		v, ok := c.terminalNode.GetBufferText()
		if ok {
			c.componentInput.SetText(v)
		}
		return
	}

	if !c.isInTerminalMode() && len(c.componentInput.text) != 0 {
		if c.gameState.input.ActionIsJustPressed(ActionInstantRunProgram) {
			c.gameState.data.UsedHiddenKeybinds = true
			c.simulationInput = string(c.componentInput.text)
			c.onProgramCompleted(c.encodeKeyword(c.simulationInput))
			return
		}
		if c.gameState.input.ActionIsJustPressed(ActionStartProgram) {
			if c.signalNode != nil {
				c.signalNode.Dispose()
			}
			c.paused = false
			c.statusLabel.text = "RUNNING"
			c.outputLabel.text = "?"
			c.outputLabel.SetColor(defaultLCDColor)
			c.simulationInput = string(c.componentInput.text)
			c.runner.Reset(c.schema, c.componentInput.text)
			c.signalNode = newSignalNode(c.schema.entry.pos)
			c.signalNode.speed = c.signalNodeSpeed
			c.prepareNextStep(c.signalNode)
			c.signalNode.EventDestinationReached.Connect(nil, c.nextStep)
			c.scene.AddObject(c.signalNode)
			return
		}
	}
}

func (c *decipherController) initComponentSchema(offset gmath.Vec) {
	c.schema = newSchemaBuilder(offset, c.config.levelTemplate).Build()
	c.keywords = append([]string{}, c.schema.keywords...)
	gmath.Shuffle(c.scene.Rand(), c.keywords)
	c.keywords = c.keywords[:c.schema.numKeywords]
}
