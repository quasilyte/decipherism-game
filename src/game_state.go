package main

import (
	"time"

	"github.com/quasilyte/ge/input"
	"github.com/quasilyte/ge/resource"
	"github.com/quasilyte/ge/xslices"
	"github.com/quasilyte/gmath"
)

type gameState struct {
	input   *input.Handler
	chapter *storyModeChapter
	level   storyModeLevel
	data    *persistentGameData
	content contentStatus
}

type chapterCompletionData struct {
	partiallyCompleted bool
	allLevelsCompleted bool
	fullyCompleted     bool
	secretDecoded      bool
}

func (state *gameState) GetLevelCompletionData(name string) *completedLevelData {
	return xslices.Find(state.data.CompletedLevels, func(l *completedLevelData) bool {
		return l.Name == name
	})
}

func (state *gameState) GetChapterCompletionData(c *storyModeChapter) chapterCompletionData {
	var result chapterCompletionData
	levelsCompleted := 0
	keywordsSolved := 0
	for i, levelName := range c.levels {
		levelData := xslices.Find(state.data.CompletedLevels, func(l *completedLevelData) bool {
			return l.Name == levelName
		})
		if levelData != nil {
			levelsCompleted++
			if levelData.SecretKeyword {
				keywordsSolved++
				if i == 0 {
					result.secretDecoded = true
				}
			}
		}
	}
	result.fullyCompleted = levelsCompleted == len(c.levels)
	if !c.IsBonus() {
		result.fullyCompleted = result.fullyCompleted && keywordsSolved == len(c.levels)
	}
	result.allLevelsCompleted = levelsCompleted == len(c.levels)
	result.partiallyCompleted = levelsCompleted != 0 && levelsCompleted >= (len(c.levels)-1)
	return result
}

type persistentGameData struct {
	CompletedLevels    []completedLevelData
	SolvedAtbash       bool
	SolvedRot13        bool
	SolvedIncDec       bool
	SolvedShift        bool
	SolvedNegation     bool
	UsedCheats         bool
	UsedHiddenKeybinds bool
	SawCollision       bool
	CompletionTime     time.Duration
	Options            gameOptions
}

type gameOptions struct {
	Music     bool
	CrtShader bool
}

type completedLevelData struct {
	Name          string
	SecretKeyword bool
}

type storyModeMap struct {
	chapters []storyModeChapter
	levels   map[string]storyModeLevel
}

func (m *storyModeMap) getChapter(name string) *storyModeChapter {
	return xslices.Find(m.chapters, func(c *storyModeChapter) bool {
		return c.name == name
	})
}

type storyModeChapter struct {
	name     string
	label    string
	keyword  string
	levels   []string
	requires string
	gridPos  gmath.Vec
}

func (c *storyModeChapter) IsBonus() bool { return c.keyword == "" }

type storyModeLevel struct {
	name string
	id   resource.RawID
}

var theStoryModeMap = &storyModeMap{
	chapters: []storyModeChapter{
		{
			label:    "1+",
			name:     "bonus1",
			requires: "story1",
			levels: []string{
				"double_negation",
				"spellbook",
				"lossy_conversion",
			},
			gridPos: gmath.Vec{X: 0, Y: 1},
		},
		{
			label:    "2+",
			name:     "bonus2",
			requires: "story2",
			levels: []string{
				"swap_shifter",
				"sub_loop",
				"branchless_encoder",
			},
			gridPos: gmath.Vec{X: 1, Y: 1},
		},
		{
			label:    "3+",
			name:     "bonus3",
			requires: "story3",
			levels: []string{
				"claws",
				"clear_head",
			},
			gridPos: gmath.Vec{X: 2, Y: 1},
		},
		{
			label:    "4+",
			name:     "bonus4",
			requires: "story4",
			levels: []string{
				"stuttering",
				"even_odd_add",
			},
			gridPos: gmath.Vec{X: 3, Y: 1},
		},
		{
			label:    "5+",
			name:     "bonus5",
			requires: "story5",
			levels: []string{
				"double_zigzag",
				"the_best_number",
			},
			gridPos: gmath.Vec{X: 4, Y: 1},
		},
		{
			label:    "6+",
			name:     "bonus6",
			requires: "story6",
			levels: []string{
				"pyramid",
				"mission_impossible",
			},
			gridPos: gmath.Vec{X: 3, Y: 2},
		},

		{
			label:   "1",
			name:    "story1",
			keyword: "rain",
			levels: []string{
				"hello_world",
				"rinse_repeat",
				"add_or_sub",
			},
			gridPos: gmath.Vec{X: 0, Y: 0},
		},
		{
			label:    "2",
			name:     "story2",
			requires: "story1",
			keyword:  "storm",
			levels: []string{
				"vowel_shifter",
				"efforts_negated",
				"loop",
			},
			gridPos: gmath.Vec{X: 1, Y: 0},
		},
		{
			label:    "3",
			name:     "story3",
			requires: "story2",
			keyword:  "thunder",
			levels: []string{
				"atbash",
				"polygraphic_atbash",
				"determination",
			},
			gridPos: gmath.Vec{X: 2, Y: 0},
		},
		{
			label:    "4",
			name:     "story4",
			requires: "story3",
			keyword:  "tsunami",
			levels: []string{
				"ladder",
				"red_herring",
				"switch",
			},
			gridPos: gmath.Vec{X: 3, Y: 0},
		},
		{
			label:    "5",
			name:     "story5",
			requires: "story4",
			keyword:  "whirlwind",
			levels: []string{
				"dotmask",
				"symmetry",
				"deduction",
				"odd_evening",
			},
			gridPos: gmath.Vec{X: 4, Y: 0},
		},
		{
			label:    "6",
			name:     "story6",
			requires: "story5",
			keyword:  "earthquake",
			levels: []string{
				"single_key",
				"rot13",
				"fixed_cond",
				"spiral",
			},
			gridPos: gmath.Vec{X: 4, Y: 2},
		},
	},
}