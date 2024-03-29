package main

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/quasilyte/decipherism-game/leveldata"
	"github.com/quasilyte/ge"
	"github.com/quasilyte/ge/tiled"
	"github.com/quasilyte/ge/ui"
	"github.com/quasilyte/gmath"
)

type customLevelSelectController struct {
	gameState *gameState

	scene *ge.Scene

	levelSlider  gmath.Slider
	allFilenames []string
	levelButtons []*levelButton

	totalCounter *ge.Label
}

type levelButton struct {
	node      *ui.Button
	fileIndex int
}

func newCustomLevelSelectController(gameState *gameState) *customLevelSelectController {
	return &customLevelSelectController{gameState: gameState}
}

func (c *customLevelSelectController) Init(scene *ge.Scene) {
	c.scene = scene
	ctx := scene.Context()

	buttonWidth := 1024.0
	offset := gmath.Vec{X: ctx.WindowWidth/2 - buttonWidth/2, Y: 164}
	var bgroup buttonGroup

	l := scene.NewLabel(FontLCDTiny)
	l.ColorScale.SetColor(defaultLCDColor)
	l.Pos.Offset = gmath.Vec{Y: 100}
	if c.gameState.userFolder != "" {
		l.Text = "scanning '" + c.gameState.userFolder + "' for levels"
	} else {
		l.Text = "$DECIPHERISM_DATA is unset"
	}
	l.Width = ctx.WindowWidth
	l.AlignHorizontal = ge.AlignHorizontalCenter
	scene.AddGraphics(l)

	uiRoot := ui.NewRoot(ctx, c.gameState.input)
	uiRoot.ActivationAction = ActionMenuConfirm
	uiRoot.NextInputAction = ActionMenuNext
	uiRoot.PrevInputAction = ActionMenuPrev
	scene.AddObject(uiRoot)

	allFilenames, err := c.scanCustomLevels()
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			l.Text = fmt.Sprintf("scan '%s/levels' for levels: non-existing path", c.gameState.userFolder)
		} else {
			l.Text = fmt.Sprintf("scan $DECIPHERISM_DATA: %v", err)
		}
	}
	c.allFilenames = allFilenames
	c.levelSlider.SetBounds(0, len(c.allFilenames)-1)

	for i := 0; i < 5; i++ {
		b := uiRoot.NewButton(optionsButtonStyle.Resized(buttonWidth, 80))
		bgroup.AddButton(b)
		buttonIndex := i
		b.EventActivated.Connect(nil, func(b *ui.Button) {
			fileIndex := c.levelButtons[buttonIndex].fileIndex
			selectedFilename := c.allFilenames[fileIndex]
			levelData, err := os.ReadFile(selectedFilename)
			if err != nil {
				panic(err) // TODO: better error handling
			}
			levelTemplate, err := loadLevelTemplate(c.scene, levelData)
			if err != nil {
				panic(err) // Should be already verified by this moment
			}
			config := decipherConfig{
				levelTemplate: levelTemplate,
			}
			c.scene.Context().ChangeScene(newDecipherController(c.gameState, config))
		})
		c.levelButtons = append(c.levelButtons, &levelButton{
			node: b,
		})
		b.Pos.Offset = offset
		scene.AddObject(b)
		offset.Y += 128
	}

	scrollButtonWidth := 320.0

	scrollBackButton := uiRoot.NewButton(optionsButtonStyle.Resized(scrollButtonWidth, 80))
	bgroup.AddButton(scrollBackButton)
	scrollBackButton.Text = "<"
	scrollBackButton.Pos.Offset = offset
	scrollBackButton.EventActivated.Connect(nil, func(_ *ui.Button) {
		c.updateSelectionPage()
	})
	scene.AddObject(scrollBackButton)

	scrollNextButton := uiRoot.NewButton(optionsButtonStyle.Resized(scrollButtonWidth, 80))
	bgroup.AddButton(scrollNextButton)
	scrollNextButton.Text = ">"
	scrollNextButton.Pos.Offset = offset.Add(gmath.Vec{X: +(buttonWidth - scrollButtonWidth)})
	scrollNextButton.EventActivated.Connect(nil, func(_ *ui.Button) {
		c.updateSelectionPage()
	})
	scene.AddObject(scrollNextButton)

	c.totalCounter = scene.NewLabel(FontLCDSmall)
	c.totalCounter.ColorScale.SetColor(defaultLCDColor)
	c.totalCounter.Width = buttonWidth
	c.totalCounter.Height = 80
	c.totalCounter.Pos.Offset = offset
	c.totalCounter.AlignHorizontal = ge.AlignHorizontalCenter
	c.totalCounter.AlignVertical = ge.AlignVerticalCenter
	c.totalCounter.Text = fmt.Sprintf("%d levels", len(allFilenames))
	scene.AddGraphics(c.totalCounter)

	offset.Y += 128

	backButtonWidth := 480.0
	backButton := uiRoot.NewButton(optionsButtonStyle.Resized(backButtonWidth, 80))
	bgroup.AddButton(backButton)
	backButton.Text = "back"
	backButton.Pos.Offset = offset.Add(gmath.Vec{X: (buttonWidth - backButtonWidth) / 2})
	backButton.EventActivated.Connect(nil, func(_ *ui.Button) {
		c.leave()
	})
	scene.AddObject(backButton)
	offset.Y += 128

	bgroup.Connect(uiRoot)
	bgroup.FocusFirst()

	c.updateSelectionPage()
}

func (c *customLevelSelectController) Update(delta float64) {
	if c.gameState.input.ActionIsJustPressed(ActionLeave) {
		c.leave()
		return
	}
}

func (c *customLevelSelectController) leave() {
	c.scene.Context().ChangeScene(newMainMenuController(c.gameState))
}

func (c *customLevelSelectController) updateSelectionPage() {
	for i, b := range c.levelButtons {
		b.fileIndex = c.levelSlider.Value()
		c.levelSlider.Inc()
		if i >= len(c.allFilenames) {
			b.node.Text = "empty"
			b.node.SetDisabled(true)
			continue
		}
		b.node.SetDisabled(false)
		filename := c.allFilenames[b.fileIndex]
		name := strings.TrimSuffix(filepath.Base(filename), ".json")
		name = strings.ReplaceAll(name, "_", " ")
		labelText := strconv.Itoa(b.fileIndex+1) + ". " + name
		if len(labelText) > 26 {
			labelText = labelText[:26] + "..."
		}
		b.node.Text = labelText
	}
}

func (c *customLevelSelectController) scanCustomLevels() ([]string, error) {
	levelsPath := filepath.Join(c.gameState.userFolder, "levels")

	var result []string
	files, err := os.ReadDir(levelsPath)
	if err != nil {
		return nil, err
	}
	tileset, err := tiled.UnmarshalTileset(c.scene.LoadRaw(RawComponentSchemaTilesetJSON).Data)
	if err != nil {
		panic(err)
	}
	for _, f := range files {
		if !strings.HasSuffix(f.Name(), ".json") {
			continue
		}
		fullName := filepath.Join(levelsPath, f.Name())
		data, err := os.ReadFile(fullName)
		if err != nil {
			fmt.Printf("[ERROR] load %q: %v\n", f.Name(), err)
			continue
		}
		if err := leveldata.ValidateLevelData(tileset, data); err != nil {
			fmt.Printf("[ERROR] load %q: %v\n", f.Name(), err)
			continue
		}
		result = append(result, fullName)
	}

	sort.Strings(result)

	return result, nil
}
