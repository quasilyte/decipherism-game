package main

import (
	"github.com/quasilyte/ge/resource"
	"github.com/quasilyte/ge/xslices"
)

type gameManual struct {
	pages []gameManualPage
}

type gameManualPage struct {
	title  string
	text   string
	image  resource.ImageID
	cond   func(*contentStatus) bool
	params func(*persistentGameData) []any
}

var theGameManual *gameManual

func init() {
	theGameManual = &gameManual{}

	theGameManual.pages = append(theGameManual.pages, gameManualPage{
		title: "Cheetah",
		text: `
			*cheat code sounds*
			\n
			Detected. (eye)
		`,
		image: ImageManualCheetah,
		cond:  func(content *contentStatus) bool { return content.usedCheats },
	})

	theGameManual.pages = append(theGameManual.pages, gameManualPage{
		title: "Let The Hackery Begin",
		text: `
			I am going to use the
			Decoder terminal to bypass
			the Hexagon security systems.
			\n
			It has no manual, so I'm going
			to write one myself.
			\n
			(There are more pages, read on.)
		`,
		image: ImageManualRTFM,
		cond:  func(content *contentStatus) bool { return true },
	})

	theGameManual.pages = append(theGameManual.pages, gameManualPage{
		title: "The Hexagon Security System",
		text: `
			There are six major blocks that form
			the defense perimeter of the Hexagon.
			\n
			Every block has at least 3 components.
			In order to hack the component,
			I need to decipher all of its keywords.
			\n
			The encoded keywords are known thanks
			to the Decoder machine. I need to input
			the decoded keywords to win.
		`,
		image: ImageManualWat,
		cond:  func(content *contentStatus) bool { return true },
	})

	theGameManual.pages = append(theGameManual.pages, gameManualPage{
		title: "Completionist Bonus",
		text: `
			Technically, it's not necessary for me to
			hack every component inside
			a major block to finish my job.
			\n
			Nonetheless, there is a reason to do it.
			Every block that is fully completed
			increases the power of the terminal.
			More features can become available.
		`,
		image: ImageManualFullCompletionBonus,
		cond:  func(content *contentStatus) bool { return content.levelsCompleted >= 6 },
	})

	theGameManual.pages = append(theGameManual.pages, gameManualPage{
		title: "Major Blocks For Dummies",
		text: `
			OK, it looks like I need to remind myself.
			\n
			Remember the previous page?
			Only the major blocks are required.
			There is no need for me to go
			through any other block!
			\n
			I will draw how they look like,
			so I will not get confused again.
		`,
		image: ImageManualMajorBlock,
		cond:  func(content *contentStatus) bool { return content.bonusLevelsCompleted >= 1 },
	})

	theGameManual.pages = append(theGameManual.pages, gameManualPage{
		title: "Decoder Controls",
		text: `
			[enter] runs the program
			[space] toggles the pause
			[tab] toggles the terminal view
			\n
			The input controls are similar to
			what I expected: arrows, backspace, etc.
			\n
			[ctrl]+[?] = ??? hmmm...
		`,
		image: ImageManualUILegend,
		cond:  func(content *contentStatus) bool { return true },
	})

	theGameManual.pages = append(theGameManual.pages, gameManualPage{
		title: "Hidden Controls",
		text: `
			I discovered the secret key combinations!
			\n
			[ctrl]+[enter] runs the instant simulation.
			\n
			[ctrl]+[num] replaces the input contents
			with the specified encoded keyword.
			[ctrl]+[1] inserts the first keyword, etc.
		`,
		image: ImageManualHiddenKeybinds,
		cond:  func(content *contentStatus) bool { return content.usedHiddenKeybinds },
	})

	theGameManual.pages = append(theGameManual.pages, gameManualPage{
		title: "Simualtion Speed",
		text: `
			There is a dial thing below
			the screen that can be used to
			change the evaluation speed.
			\n
			If I want to have a good score
			on a leaderboard, I should
			rotate the hell out of that knob.
		`,
		image: ImageManualHackingTooMuchTime,
		cond:  func(content *contentStatus) bool { return content.bonusLevelsCompleted >= 3 },
	})

	theGameManual.pages = append(theGameManual.pages, gameManualPage{
		title: "Element Shapes",
		text: `
			Every element has two main aspects:
			its inscription and its outline shape.
			\n
			- Circle: a simple enter/leave element
			- Square: a data transformer
			- Diamond: a conditional path branching
			\n
			The "+" sign shows the "true" path.
		`,
		image: ImageManualShapes,
		cond:  func(content *contentStatus) bool { return content.levelsCompleted >= 1 },
	})

	theGameManual.pages = append(theGameManual.pages, gameManualPage{
		title: "Transformation Scopes",
		text: `
			Most tranfrormation operations have
			scopes. They may operate on all
			letters or on some portion of them.
			\n
			The three rectangles represent letters.
			The first and last rectangles stand for
			the first and last letters respectively.
			The middle one is everything in between.
			\n
			Outline-only sections are not transformed.
		`,
		image: ImageManualScopes,
		cond:  func(content *contentStatus) bool { return content.levelsCompleted >= 4 },
	})

	theGameManual.pages = append(theGameManual.pages, gameManualPage{
		title: "Conditional Transformations",
		text: `
			A full transformation element format
			consists of these three parts:
			\n
			1. Condition
			2. Operation
			3. Scope
			\n
			It's possible to find a "+" element
			that changes only dot-marked letters.
		`,
		image: ImageConditionalTransformations,
		cond:  func(content *contentStatus) bool { return content.solvedCondTransform },
	})

	theGameManual.pages = append(theGameManual.pages, gameManualPage{
		title: "Output Collision",
		text: `
			Sometimes different inputs can
			produce the same output.
			This output may look exactly
			like the encoded keyword even if
			the input is not a keyword
			I'm looking for.
			\n
			The keywords are usually valid words.
			If input looks like some gibberish,
			more often than not it's close but no cigar.
		`,
		image: ImageManualOutputCollision,
		cond:  func(content *contentStatus) bool { return content.sawCollision },
	})

	theGameManual.pages = append(theGameManual.pages, gameManualPage{
		title: "Compound Pipeline Ciphers",
		text: `
			Each major block has its a
			cross-component keyword.
			\n
			Some secret word is given as an input
			to the 1st component. The output of
			the 1st component is then used as
			an input of the 2nd component, etc.
			\n
			In order to decode it, I need to
			backtrack the steps one by one.
		`,
		image: ImageManualCipherPipeline,
		cond:  func(content *contentStatus) bool { return content.levelsCompleted >= 9 },
	})

	theGameManual.pages = append(theGameManual.pages, gameManualPage{
		title: "Transposition Ciphers",
		text: `
			Some elements move the letters
			around without changing them.
			These operations are similar to
			transposition ciphers.
		`,
		image: ImageManualTranspositionCiphers,
		cond:  func(content *contentStatus) bool { return content.solvedShift },
	})

	theGameManual.pages = append(theGameManual.pages, gameManualPage{
		title: "Substitution Ciphers",
		text: `
			The elements like [+] and [-]
			apply a substitution cipher.
			\n
			They replace the characters from
			the text by some other characters.
			The decoding process involves the
			reversed substitution.
			\n
			The Caesar cipher is one of the
			examples of the substitution ciphers.
		`,
		image: ImageManualSubstitutionCiphers,
		cond:  func(content *contentStatus) bool { return content.solvedIncDec },
	})

	theGameManual.pages = append(theGameManual.pages, gameManualPage{
		title: "Polygraphic Ciphers",
		text: `
			Some ciphers may operate on
			groups of letters.
			\n
			A polygraphic atbash-like encoding
			illustrated here replaces the pairs
			of dot-marked letters.
		`,
		image: ImageManualPolygraphicCiphers,
		cond:  func(content *contentStatus) bool { return content.solvedPolygraphic },
	})

	theGameManual.pages = append(theGameManual.pages, gameManualPage{
		title: "Negation",
		text: `
			This peculiar Anglular thing is negation.
			\n
			The logical negation indicates that
			the truth value of the statement
			is being reversed.
			\n
			Practically speaking, if there is a
			branching with "cool" condition,
			a negated branching would react on
			"lame" condition.
		`,
		image: ImageManualNegation,
		cond:  func(content *contentStatus) bool { return content.solvedNegation },
	})

	theGameManual.pages = append(theGameManual.pages, gameManualPage{
		title: "Atbash Cipher",
		text: `
			A component labeled ZYX implements
			a cipher named Atbash.
			\n
			Atbash uses a reversed alphabet for
			encoding. The dotted patterns on
			the paper show the direct letters
			relation, so the decoding should be
			quite straightforward for me.
		`,
		image: ImageManualHintAtbash,
		cond:  func(content *contentStatus) bool { return content.solvedAtbash },
	})

	theGameManual.pages = append(theGameManual.pages, gameManualPage{
		title: "Rot-13 Cipher",
		text: `
			A ROT-13 uses an alphabet shifted
			by 13 positions.
			\n
			So, A becomes N, B becomes O.
			\n
			The dotted patterns can
			be used as guidelines. The sum of
			the circles is 8 for the matching pairs.
		`,
		image: ImageManualHintRot13,
		cond:  func(content *contentStatus) bool { return content.solvedRot13 },
	})

	theGameManual.pages = append(theGameManual.pages, gameManualPage{
		title: "Smart Decoding",
		text: `
			The encoded keyword is NQJGN.
			I start with reversing the two + shifts
			at once: LOHEL, it's a single step.
			\n
			At that point, I can guess the word
			already, but just to be sure, a double
			rotation would reveal HELLO.
			\n
			So it's possible to decode this
			schema in 2 steps instead of 4.
		`,
		image: ImageManualOptimizedDecoding,
		cond:  func(content *contentStatus) bool { return content.levelsCompleted >= 14 },
	})

	theGameManual.pages = append(theGameManual.pages, gameManualPage{
		title: "Value Inspector",
		text: `
			If I pause the simulation with [space]
			and open the terminal with [tab],
			I will see the current input state.
			\n
			I can use that to learn the effects
			of any individual transformation element.
		`,
		image: ImageManualValueInspector,
		cond: func(content *contentStatus) bool {
			return xslices.Contains(content.techLevelFeatures, "value inspector")
		},
	})

	theGameManual.pages = append(theGameManual.pages, gameManualPage{
		title: "Text Buffer",
		text: `
			When text buffer is available,
			[ctrl]+[c] stores the current
			input value in it. Pressing [ctrl]+[v] replaces
			the input contents with the previously saved text.
			\n
			[ctrl]+[x] works as expected too.
			\n
			I can also use it as a scratch pad
			when my terminal is open.
		`,
		image: ImageManualTextBuffer,
		cond: func(content *contentStatus) bool {
			return xslices.Contains(content.techLevelFeatures, "text buffer")
		},
	})

	theGameManual.pages = append(theGameManual.pages, gameManualPage{
		title: "Branching Info",
		text: `
			This routine dumps the branching condition
			of up to three conditional elements.
			\n
			It doesn't know the exact condition,
			only its basic operation. For instance, it may tell
			that some branch checks for the string length,
			but I need to find out that magic value on my own.
			\n
			The order of the entries seem to be random too.
		`,
		image: ImageManualBranchingInfo,
		cond: func(content *contentStatus) bool {
			return xslices.Contains(content.techLevelFeatures, "branching info")
		},
	})

	theGameManual.pages = append(theGameManual.pages, gameManualPage{
		title: "Conditions Vocab",
		text: `
			- len: length (number of letters)
			- substr: substring, a part of the word
			- eq: equals (=)
			- lt: less than (<)
			- gt: greater than (>)
			- even: numbers like 2, 4, 6, ...
			- odd: numbers like 1, 3, 5, ...
			\n
			There is also an FNV hash condition that
			can't realistically be predicted.
		`,
		image: ImageManualConditionsVocab,
		cond: func(content *contentStatus) bool {
			return xslices.Contains(content.techLevelFeatures, "branching info")
		},
	})

	theGameManual.pages = append(theGameManual.pages, gameManualPage{
		title: "I/O Logs",
		text: `
			Records the last three runs.
			It saves the input->out pairs.
			\n
			Could be useful when trying to analyze
			the several inputs and outputs together.
			\n
			As always, it's only accessible from the terminal.
		`,
		image: ImageManualIOLogs,
		cond: func(content *contentStatus) bool {
			return xslices.Contains(content.techLevelFeatures, "i/o logs")
		},
	})

	theGameManual.pages = append(theGameManual.pages, gameManualPage{
		title: "Output Predictor",
		text: `
			This thing estimates the encoding results
			in real time!
			\n
			I open the terminal, type something into the
			input and it prints out the answer for me.
			\n
			Some of the characters could still
			remain masked, but a keen mind
			would guess it anyway.
		`,
		image: ImageManualOutputPredictor,
		cond: func(content *contentStatus) bool {
			return xslices.Contains(content.techLevelFeatures, "output predictor")
		},
	})

	theGameManual.pages = append(theGameManual.pages, gameManualPage{
		title: "Advanced Input Commands",
		text: `
			The improved input system includes new moves:
			\n
			[ctrl][=] increment current letter
			[ctrl][-] decrement current letter
			[ctrl][left] shift letters left
			[ctrl][right] shift letters right
			\n
			This should make things a lot easier.
		`,
		image: ImageManualAdvancedInput,
		cond: func(content *contentStatus) bool {
			return xslices.Contains(content.techLevelFeatures, "advanced input")
		},
	})

	theGameManual.pages = append(theGameManual.pages, gameManualPage{
		title: "WOW I'm Such a Nerd",
		text: `
			I hacked throught... everything?
			\n
			It took me %d seconds (%d minutes)
			to clear all %d levels.
			\n
			And yes.
			I am a dog.
			I've always been a stooped dog.
		`,
		params: func(data *persistentGameData) []any {
			return []any{
				int(data.CompletionTime.Seconds()),
				int(data.CompletionTime.Minutes()),
				len(theStoryModeMap.levels),
			}
		},
		image: ImageManualFinalTreasure,
		cond:  func(content *contentStatus) bool { return content.hackedEverything },
	})
}
