package main

import (
	"embed"
	"io"
	"time"

	"github.com/quasilyte/ge"
	"github.com/quasilyte/ge/input"
	"github.com/quasilyte/ge/resource"

	_ "image/png"
)

//go:embed all:_assets
var gameAssets embed.FS

const (
	AudioNone resource.AudioID = iota
	AudioDecodingSuccess
	AudioSecretUnlocked
	AudioCollision
	AudioMenuMusic
	AudioDecipherMusic
)

const (
	FontLCDSmall resource.FontID = iota
	FontLCDNormal
	FontSmall
	FontHandwritten
)

const (
	ImageNone resource.ImageID = iota
	ImagePaperBg
	ImageDecipherBg
	ImageSchemaBg
	ImageTerminalBg
	ImageSignal
	ImagePingParticle
	ImageDialButton
	ImageDialButtonArrow
	ImageOnOffButton
	ImageManualCheetah
	ImageManualRTFM
	ImageManualWat
	ImageManualUILegend
	ImageManualHiddenKeybinds
	ImageManualMajorBlock
	ImageManualFullCompletionBonus
	ImageManualShapes
	ImageManualScopes
	ImageConditionalTransformations
	ImageManualOutputCollision
	ImageManualOptimizedDecoding
	ImageManualNegation
	ImageManualTranspositionCiphers
	ImageManualSubstitutionCiphers
	ImageManualPolygraphicCiphers
	ImageManualHintAtbash
	ImageManualHintRot13
	ImageManualValueInspector
	ImageManualTextBuffer
	ImageManualBranchingInfo
	ImageManualIOLogs
	ImageManualOutputPredictor
	ImageManualAdvancedInput
	ImageManualCipherPipeline
	ImageManualHackingTooMuchTime
	ImageManualFinalTreasure
	ImageBlueMarker
	ImageChapterSelectOutline
	ImageCompleteMark
	ImagePipelineArrow
	componentSchemaImageOffset
	ImagePipeConnect2
	ImageSpecialAnglePipe
	ImageAnglePipe
	ImageSpecialPipe
	ImagePipe
	ImageElemInput
	ImageElemOutput
	ImageElemMux
	ImageElemIf
	ImageElemIfNot
	ImageElemRepeater
	ImageElemInvRepeater
	ImageElemCountdown0
	ImageElemCountdown1
	ImageElemCountdown2
	ImageElemCountdown3
	ImageElemReverse
	ImageElemSwapHalves
	ImageElemRotateLeft
	ImageElemRotateLeftButfirst
	ImageElemRotateRight
	ImageElemRotateRightButfirst
	ImageElemAdd
	ImageElemAddButfirst
	ImageElemAddFirst
	ImageElemAddLast
	ImageElemAddNowrap
	ImageElemAddButfirstNowrap
	ImageElemAddDotted
	ImageElemAddButfirstDotted
	ImageElemAddEven
	ImageElemAddOdd
	ImageElemSub
	ImageElemSubButlast
	ImageElemSubFirst
	ImageElemSubLast
	ImageElemSubNowrap
	ImageElemSubUndotted
	ImageElemSubEven
	ImageElemSubOdd
	ImageElemAtbash
	ImageElemAtbashButlast
	ImageElemAtbashFirst
	ImageElemPolygraphicAtbash
	ImageElemRot13
	ImageElemRot13Butfirst
	ImageElemRot13Butlast
	ImageElemRot13First
	ImageElemHardshiftLeft
	ImageElemHardshiftRight
	ImageElemZigzag
)

const (
	RawNone resource.RawID = iota
	RawComponentSchemaTilesetJSON
	rawLastID
	RawLevel1JSON
	RawLevel2JSON
	RawLevel3JSON
	RawLevel4JSON
	RawLevel5JSON
	RawLevel6JSON
	RawLevel7JSON
	RawLevel8JSON
	RawLevel9JSON
	RawLevel10JSON
	RawLevel11JSON
	RawLevel12JSON
	RawLevel13JSON
	RawLevel14JSON
	RawLevel15JSON
	RawLevel16JSON
)

const (
	ActionUnknown input.Action = iota
	ActionClearStage
	ActionMenuConfirm
	ActionMenuNext
	ActionMenuPrev
	ActionMenuNextPage
	ActionMenuPrevPage
	ActionLeave
	ActionModeSwap
	ActionPasteKeyword1
	ActionPasteKeyword2
	ActionPasteKeyword3
	ActionPasteKeyword4
	ActionInstantRunProgram
	ActionStartProgram
	ActionPauseProgram
	ActionButton
	ActionCursorLeft
	ActionCursorRight
	ActionRemovePrevChar
	ActionRemoveCurrentChar
	ActionBufferCut
	ActionBufferCopy
	ActionBufferPaste
	ActionCharInc
	ActionCharDec
	ActionRotateLeft
	ActionRotateRight
)

const (
	ShaderNone resource.ShaderID = iota
	ShaderHandwriting
	ShaderVideoDistortion
)

func main() {
	ctx := ge.NewContext()
	ctx.Rand.SetSeed(time.Now().Unix())
	ctx.GameName = "cipher_cracker"
	ctx.WindowTitle = "Cipher Cracker"
	ctx.WindowWidth = 1920
	ctx.WindowHeight = 1080
	ctx.FullScreen = true

	ctx.Loader.OpenAssetFunc = func(path string) io.ReadCloser {
		f, err := gameAssets.Open("_assets/" + path)
		if err != nil {
			ctx.OnCriticalError(err)
		}
		return f
	}

	// Associate audio resources.
	audioResources := map[resource.AudioID]resource.Audio{
		AudioDecodingSuccess: {Path: "audio/decoding_success.wav", Volume: -0.4},
		AudioSecretUnlocked:  {Path: "audio/secret_unlocked.wav", Volume: -0.35},
		AudioCollision:       {Path: "audio/collision.wav", Volume: -0.2},
		AudioMenuMusic:       {Path: "audio/menu.ogg", Volume: -0.5},
		AudioDecipherMusic:   {Path: "audio/hack.ogg", Volume: -0.55},
	}
	for id, res := range audioResources {
		ctx.Loader.AudioRegistry.Set(id, res)
		ctx.Loader.PreloadAudio(id)
	}

	// Associate font resources.
	fontResources := map[resource.FontID]resource.Font{
		FontLCDSmall:    {Path: "font.ttf", Size: 32},
		FontLCDNormal:   {Path: "font.ttf", Size: 40},
		FontSmall:       {Path: "sector_017.otf", Size: 36},
		FontHandwritten: {Path: "TidyHand.ttf", Size: 50},
	}
	for id, res := range fontResources {
		ctx.Loader.FontRegistry.Set(id, res)
		ctx.Loader.PreloadFont(id)
	}

	// Associate image resources.
	imageResources := map[resource.ImageID]resource.ImageInfo{
		ImagePaperBg:         {Path: "paper_bg.png"},
		ImageDecipherBg:      {Path: "decipher_bg.png"},
		ImageSchemaBg:        {Path: "schema_bg.png"},
		ImageTerminalBg:      {Path: "terminal_bg.png"},
		ImageSignal:          {Path: "signal.png"},
		ImagePingParticle:    {Path: "ping_particle.png"},
		ImageDialButton:      {Path: "dial_button.png"},
		ImageDialButtonArrow: {Path: "dial_button_arrow.png"},
		ImageOnOffButton:     {Path: "onoff_toggle.png", FrameWidth: 96},

		ImageManualCheetah:              {Path: "manual/cheetah.png"},
		ImageManualRTFM:                 {Path: "manual/rtfm.png"},
		ImageManualWat:                  {Path: "manual/wat.png"},
		ImageManualUILegend:             {Path: "manual/ui_legend.png"},
		ImageManualHiddenKeybinds:       {Path: "manual/hidden_keybinds.png"},
		ImageManualMajorBlock:           {Path: "manual/major_block.png"},
		ImageManualFullCompletionBonus:  {Path: "manual/full_completion_bonus.png"},
		ImageManualOptimizedDecoding:    {Path: "manual/optimized_decoding.png"},
		ImageManualShapes:               {Path: "manual/shapes.png"},
		ImageManualScopes:               {Path: "manual/scopes.png"},
		ImageConditionalTransformations: {Path: "manual/conditional_transformations.png"},
		ImageManualNegation:             {Path: "manual/negation.png"},
		ImageManualOutputCollision:      {Path: "manual/collisions.png"},
		ImageManualTranspositionCiphers: {Path: "manual/hint_transposition_ciphers.png"},
		ImageManualSubstitutionCiphers:  {Path: "manual/hint_substitution_ciphers.png"},
		ImageManualPolygraphicCiphers:   {Path: "manual/hint_polygraphic_ciphers.png"},
		ImageManualHintAtbash:           {Path: "manual/hint_atbash.png"},
		ImageManualHintRot13:            {Path: "manual/hint_rot13.png"},
		ImageManualValueInspector:       {Path: "manual/value_inspector.png"},
		ImageManualTextBuffer:           {Path: "manual/text_buffer.png"},
		ImageManualBranchingInfo:        {Path: "manual/branching_info.png"},
		ImageManualIOLogs:               {Path: "manual/io_logs.png"},
		ImageManualOutputPredictor:      {Path: "manual/output_predictor.png"},
		ImageManualAdvancedInput:        {Path: "manual/advanced_input.png"},
		ImageManualCipherPipeline:       {Path: "manual/cipher_pipeline.png"},
		ImageManualHackingTooMuchTime:   {Path: "manual/hacking_too_much_time.png"},
		ImageManualFinalTreasure:        {Path: "manual/final_treasure.png"},

		ImageBlueMarker:           {Path: "blue_marker.png"},
		ImageChapterSelectOutline: {Path: "chapter_select_outline.png"},
		ImageCompleteMark:         {Path: "complete_mark.png"},
		ImagePipelineArrow:        {Path: "pipeline_arrow.png"},

		ImagePipeConnect2:            {Path: "elements/pipe_connect2.png"},
		ImageSpecialAnglePipe:        {Path: "elements/special_angle_pipe.png"},
		ImageAnglePipe:               {Path: "elements/angle_pipe.png"},
		ImagePipe:                    {Path: "elements/pipe.png"},
		ImageSpecialPipe:             {Path: "elements/special_pipe.png"},
		ImageElemInput:               {Path: "elements/elem_input.png"},
		ImageElemOutput:              {Path: "elements/elem_output.png"},
		ImageElemMux:                 {Path: "elements/elem_mux.png"},
		ImageElemIf:                  {Path: "elements/elem_if.png"},
		ImageElemIfNot:               {Path: "elements/elem_ifnot.png"},
		ImageElemRepeater:            {Path: "elements/elem_repeater.png"},
		ImageElemInvRepeater:         {Path: "elements/elem_inv_repeater.png"},
		ImageElemCountdown0:          {Path: "elements/elem_countdown0.png"},
		ImageElemCountdown1:          {Path: "elements/elem_countdown1.png"},
		ImageElemCountdown2:          {Path: "elements/elem_countdown2.png"},
		ImageElemCountdown3:          {Path: "elements/elem_countdown3.png"},
		ImageElemReverse:             {Path: "elements/elem_reverse.png"},
		ImageElemSwapHalves:          {Path: "elements/elem_swap_halves.png"},
		ImageElemRotateLeft:          {Path: "elements/elem_rotate_left.png"},
		ImageElemRotateLeftButfirst:  {Path: "elements/elem_rotate_left_butfirst.png"},
		ImageElemRotateRight:         {Path: "elements/elem_rotate_right.png"},
		ImageElemRotateRightButfirst: {Path: "elements/elem_rotate_right_butfirst.png"},
		ImageElemAdd:                 {Path: "elements/elem_add.png"},
		ImageElemAddButfirst:         {Path: "elements/elem_add_butfirst.png"},
		ImageElemAddFirst:            {Path: "elements/elem_add_first.png"},
		ImageElemAddLast:             {Path: "elements/elem_add_last.png"},
		ImageElemAddNowrap:           {Path: "elements/elem_add_nowrap.png"},
		ImageElemAddButfirstNowrap:   {Path: "elements/elem_add_butfirst_nowrap.png"},
		ImageElemAddButfirstDotted:   {Path: "elements/elem_add_butfirst_dotted.png"},
		ImageElemAddDotted:           {Path: "elements/elem_add_dotted.png"},
		ImageElemAddEven:             {Path: "elements/elem_add_even.png"},
		ImageElemAddOdd:              {Path: "elements/elem_add_odd.png"},
		ImageElemSub:                 {Path: "elements/elem_sub.png"},
		ImageElemSubButlast:          {Path: "elements/elem_sub_butlast.png"},
		ImageElemSubFirst:            {Path: "elements/elem_sub_first.png"},
		ImageElemSubLast:             {Path: "elements/elem_sub_last.png"},
		ImageElemSubNowrap:           {Path: "elements/elem_sub_nowrap.png"},
		ImageElemSubUndotted:         {Path: "elements/elem_sub_undotted.png"},
		ImageElemSubEven:             {Path: "elements/elem_sub_even.png"},
		ImageElemSubOdd:              {Path: "elements/elem_sub_odd.png"},
		ImageElemAtbash:              {Path: "elements/elem_atbash.png"},
		ImageElemAtbashButlast:       {Path: "elements/elem_atbash_butlast.png"},
		ImageElemAtbashFirst:         {Path: "elements/elem_atbash_first.png"},
		ImageElemPolygraphicAtbash:   {Path: "elements/elem_polygraphic_atbash.png"},
		ImageElemRot13:               {Path: "elements/elem_rot13.png"},
		ImageElemRot13Butfirst:       {Path: "elements/elem_rot13_butfirst.png"},
		ImageElemRot13Butlast:        {Path: "elements/elem_rot13_butlast.png"},
		ImageElemRot13First:          {Path: "elements/elem_rot13_first.png"},
		ImageElemHardshiftLeft:       {Path: "elements/elem_hardshift_left.png"},
		ImageElemHardshiftRight:      {Path: "elements/elem_hardshift_right.png"},
		ImageElemZigzag:              {Path: "elements/elem_zigzag.png"},
	}
	for id, res := range imageResources {
		ctx.Loader.ImageRegistry.Set(id, res)
		ctx.Loader.PreloadImage(id)
	}

	prepareAssets(ctx)

	// Associate other resources.
	rawResources := map[resource.RawID]resource.Raw{
		RawComponentSchemaTilesetJSON: {Path: "schemas.tsj"},
	}
	for id, res := range rawResources {
		ctx.Loader.RawRegistry.Set(id, res)
		ctx.Loader.PreloadRaw(id)
	}

	// Associate shader resources.
	shaderResources := map[resource.ShaderID]resource.ShaderInfo{
		ShaderVideoDistortion: {Path: "shader/video_distortion.go"},
		ShaderHandwriting:     {Path: "shader/handwriting.go"},
	}
	for id, res := range shaderResources {
		ctx.Loader.ShaderRegistry.Set(id, res)
		ctx.Loader.PreloadShader(id)
	}

	state := &gameState{
		data: &persistentGameData{
			Options: gameOptions{
				Music:     true,
				CrtShader: true,
			},
		},
	}
	ctx.LoadGameData("save", &state.data)

	// Bind controls.
	keymap := input.Keymap{
		ActionMenuConfirm:       {input.KeyMouseLeft, input.KeyEnter},
		ActionMenuNext:          {input.KeyDown},
		ActionMenuPrev:          {input.KeyUp},
		ActionMenuNextPage:      {input.KeyDown, input.KeyRight},
		ActionMenuPrevPage:      {input.KeyUp, input.KeyLeft},
		ActionLeave:             {input.KeyEscape},
		ActionModeSwap:          {input.KeyTab},
		ActionPasteKeyword1:     {input.KeyWithModifier(input.Key1, input.ModControl)},
		ActionPasteKeyword2:     {input.KeyWithModifier(input.Key2, input.ModControl)},
		ActionPasteKeyword3:     {input.KeyWithModifier(input.Key3, input.ModControl)},
		ActionPasteKeyword4:     {input.KeyWithModifier(input.Key4, input.ModControl)},
		ActionInstantRunProgram: {input.KeyWithModifier(input.KeyEnter, input.ModControl)},
		ActionStartProgram:      {input.KeyEnter},
		ActionPauseProgram:      {input.KeySpace},
		ActionButton:            {input.KeyMouseLeft},
		ActionRemovePrevChar:    {input.KeyBackspace},
		ActionRemoveCurrentChar: {input.KeyDelete},
		ActionCursorLeft:        {input.KeyLeft},
		ActionCursorRight:       {input.KeyRight},
		ActionBufferCut:         {input.KeyWithModifier(input.KeyX, input.ModControl)},
		ActionBufferCopy:        {input.KeyWithModifier(input.KeyC, input.ModControl)},
		ActionBufferPaste:       {input.KeyWithModifier(input.KeyV, input.ModControl)},
		ActionCharInc:           {input.KeyWithModifier(input.KeyEqual, input.ModControl)},
		ActionCharDec:           {input.KeyWithModifier(input.KeyMinus, input.ModControl)},
		ActionRotateLeft:        {input.KeyWithModifier(input.KeyLeft, input.ModControl)},
		ActionRotateRight:       {input.KeyWithModifier(input.KeyRight, input.ModControl)},
		ActionClearStage:        {input.KeyWithModifier(input.KeyBackquote, input.ModShift)},
	}
	state.input = ctx.Input.NewHandler(0, keymap)

	if err := ge.RunGame(ctx, newMainMenuController(state)); err != nil {
		panic(err)
	}
}
