package main

type contentStatus struct {
	techLevelFeatures    []string
	chapters             []string
	manualPages          []string
	levelsCompleted      int
	bonusLevelsCompleted int
	solvedShift          bool
	solvedIncDec         bool
	solvedAtbash         bool
	solvedRot13          bool
	solvedNegation       bool
	hackedEverything     bool
	usedCheats           bool
	usedHiddenKeybinds   bool
	sawCollision         bool
}

func calculateContentStatus(state *gameState) contentStatus {
	result := contentStatus{
		solvedShift:        state.data.SolvedShift,
		solvedIncDec:       state.data.SolvedIncDec,
		solvedAtbash:       state.data.SolvedAtbash,
		solvedRot13:        state.data.SolvedRot13,
		solvedNegation:     state.data.SolvedNegation,
		usedCheats:         state.data.UsedCheats,
		usedHiddenKeybinds: state.data.UsedHiddenKeybinds,
		sawCollision:       state.data.SawCollision,
	}

	chaptersCleared := 0
	for i := range theStoryModeMap.chapters {
		chapter := &theStoryModeMap.chapters[i]
		completionData := state.GetChapterCompletionData(chapter)
		if completionData.allLevelsCompleted && !chapter.IsBonus() {
			chaptersCleared++
		}
		for _, levelName := range chapter.levels {
			levelCompletionData := state.GetLevelCompletionData(levelName)
			if levelCompletionData != nil {
				if chapter.IsBonus() {
					result.bonusLevelsCompleted++
				} else {
					result.levelsCompleted++
				}
			}
		}
		available := true
		if chapter.requires != "" {
			otherChapter := theStoryModeMap.getChapter(chapter.requires)
			otherChapterStatus := state.GetChapterCompletionData(otherChapter)
			if chapter.IsBonus() {
				available = otherChapterStatus.secretDecoded
			} else {
				available = otherChapterStatus.partiallyCompleted
			}
		}
		if available {
			result.chapters = append(result.chapters, chapter.name)
		}
	}
	result.hackedEverything = result.levelsCompleted+result.bonusLevelsCompleted == len(theStoryModeMap.levels)

	techLevel := chaptersCleared
	if techLevel >= 1 {
		result.techLevelFeatures = append(result.techLevelFeatures, "value inspector")
	}
	if techLevel >= 2 {
		result.techLevelFeatures = append(result.techLevelFeatures, "branching info")
	}
	if techLevel >= 3 {
		result.techLevelFeatures = append(result.techLevelFeatures, "text buffer")
	}
	if techLevel >= 4 {
		result.techLevelFeatures = append(result.techLevelFeatures, "i/o logs")
	}
	if techLevel >= 5 {
		result.techLevelFeatures = append(result.techLevelFeatures, "output predictor")
	}
	if techLevel >= 6 {
		result.techLevelFeatures = append(result.techLevelFeatures, "advanced input")
	}

	for _, p := range theGameManual.pages {
		if p.cond(&result) {
			result.manualPages = append(result.manualPages, p.title)
		}
	}

	return result
}
