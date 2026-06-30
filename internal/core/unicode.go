package core

import (
	"fmt"
	"strings"
	"unicode"

	"bytes"

	"github.com/go-text/typesetting/di"
	gt_font "github.com/go-text/typesetting/font"
	"github.com/go-text/typesetting/language"
	"github.com/go-text/typesetting/shaping"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

// Unicode shaping and complex script support with proper HarfBuzz integration

// TextDirection represents text direction
type TextDirection int

const (
	TextDirectionLTR TextDirection = iota // Left-to-Right
	TextDirectionRTL                      // Right-to-Left
	TextDirectionTTB                      // Top-to-Bottom
	TextDirectionBTT                      // Bottom-to-Top
)

// ScriptType represents different script types
type ScriptType int

const (
	ScriptLatin ScriptType = iota
	ScriptArabic
	ScriptHebrew
	ScriptDevanagari
	ScriptThai
	ScriptChinese
	ScriptJapanese
	ScriptKorean
	ScriptCyrillic
	ScriptGreek
	ScriptBengali
	ScriptTamil
	ScriptTelugu
	ScriptKhmer
	ScriptMyanmar
)

// TextShaper handles complex text shaping using HarfBuzz
type TextShaper struct {
	Direction       TextDirection
	Script          ScriptType
	Language        string
	EnableLigatures bool
	EnableKerning   bool

	// go-text/typesetting integration
	shaper   *shaping.HarfbuzzShaper
	fontFace *gt_font.Face

	// Standard Go font face for fallback
	GoFontFace font.Face

	// Per-script font overrides for multilingual support
	scriptFonts map[ScriptType]*scriptFont

	// Font size for shaping (matches rendering scale)
	fontSize float64
}

// scriptFont holds font data for a specific script
type scriptFont struct {
	face     *gt_font.Face
	fontData []byte
	ttfFont  *truetype.Font // freetype font for glyph rendering
}

// NewTextShaper creates a new text shaper
func NewTextShaper() *TextShaper {
	return &TextShaper{
		Direction:       TextDirectionLTR,
		Script:          ScriptLatin,
		Language:        "en",
		EnableLigatures: true,
		EnableKerning:   true,
	}
}

// ShapedText represents shaped text with positioning
type ShapedText struct {
	Glyphs    []ShapedGlyph
	Direction TextDirection
	Width     float64
	Height    float64
}

// ShapedGlyph represents a positioned glyph
type ShapedGlyph struct {
	GlyphID   uint32
	X, Y      float64
	AdvanceX  float64
	AdvanceY  float64
	Cluster   int
	Character rune

	// Font data for this glyph (nil means use context's default font)
	font *truetype.Font
}

// SetFont sets the font for shaping
func (ts *TextShaper) SetFont(fontFace *gt_font.Face) error {
	if fontFace == nil {
		return fmt.Errorf("font face cannot be nil")
	}

	ts.fontFace = fontFace
	ts.shaper = &shaping.HarfbuzzShaper{}
	return nil
}

// SetFontSize sets the font size for shaping (should match rendering scale)
func (ts *TextShaper) SetFontSize(size float64) {
	ts.fontSize = size
}

// SetGoFontFace sets the standard Go font face for fallback shaping
func (ts *TextShaper) SetGoFontFace(face font.Face) {
	ts.GoFontFace = face
}

// HasFont returns true if a font is loaded for shaping
func (ts *TextShaper) HasFont() bool {
	return ts.fontFace != nil
}

// SetFontBytes sets the font for shaping from raw SFNT bytes
func (ts *TextShaper) SetFontBytes(fontData []byte, points float64) error {
	f, err := gt_font.ParseTTF(bytes.NewReader(fontData))
	if err != nil {
		return err
	}

	ts.fontFace = f
	ts.shaper = &shaping.HarfbuzzShaper{}
	ts.fontSize = points * 72 / 96 // Match rendering scale (fontHeight)

	return nil
}

// SetScriptFontBytes sets a font for a specific script (e.g., Bengali, Arabic)
func (ts *TextShaper) SetScriptFontBytes(script ScriptType, fontData []byte, points float64) error {
	f, err := gt_font.ParseTTF(bytes.NewReader(fontData))
	if err != nil {
		return err
	}

	ttfFont, err := truetype.Parse(fontData)
	if err != nil {
		return err
	}

	if ts.scriptFonts == nil {
		ts.scriptFonts = make(map[ScriptType]*scriptFont)
	}
	ts.scriptFonts[script] = &scriptFont{
		face:     f,
		fontData: fontData,
		ttfFont:  ttfFont,
	}
	ts.fontSize = points * 72 / 96 // Match rendering scale (fontHeight)

	return nil
}

// GetScriptFont returns the font for a specific script, or nil if not set
func (ts *TextShaper) GetScriptFont(script ScriptType) *scriptFont {
	if ts.scriptFonts == nil {
		return nil
	}
	return ts.scriptFonts[script]
}

// ShapeText shapes text according to Unicode rules using proper BiDi and HarfBuzz
func (ts *TextShaper) ShapeText(text string) *ShapedText {
	if ts.shaper == nil || ts.fontFace == nil {
		// Fallback to simple shaping if no font is set
		return ts.fallbackShapeText(text)
	}

	// Apply bidirectional algorithm for proper text ordering
	runs := ts.segmentBidiRuns(text)

	shaped := &ShapedText{
		Glyphs: make([]ShapedGlyph, 0),
	}

	currentX := 0.0

	// Shape each run separately
	for _, run := range runs {
		runShaped := ts.shapeRun(run)

		// Adjust positions
		for i := range runShaped.Glyphs {
			runShaped.Glyphs[i].X += currentX
		}

		shaped.Glyphs = append(shaped.Glyphs, runShaped.Glyphs...)
		currentX += runShaped.Width
	}

	shaped.Width = currentX
	shaped.Height = 16.0 // Default height, will be updated when we have proper font metrics

	return shaped
}

// BidiRun represents a run of text with consistent direction and script
type BidiRun struct {
	Text      string
	Direction TextDirection
	Script    ScriptType
	Language  string
	Level     int
}

// segmentBidiRuns segments text into bidirectional runs
func (ts *TextShaper) segmentBidiRuns(text string) []BidiRun {
	// For now, use simple segmentation until we properly integrate bidi
	// The golang.org/x/text/unicode/bidi API is complex and needs proper setup

	// Simple fallback: treat entire text as single run
	direction := TextDirectionLTR
	script := ts.detectScript(text)

	// Basic RTL detection
	for _, r := range text {
		if unicode.Is(unicode.Arabic, r) || unicode.Is(unicode.Hebrew, r) {
			direction = TextDirectionRTL
			break
		}
	}

	return []BidiRun{
		{
			Text:      text,
			Direction: direction,
			Script:    script,
			Language:  ts.Language,
			Level:     0,
		},
	}
}

// detectScript detects the primary script in a text run
func (ts *TextShaper) detectScript(text string) ScriptType {
	arabicCount := 0
	hebrewCount := 0
	devanagariCount := 0
	thaiCount := 0
	cjkCount := 0
	cyrillicCount := 0
	greekCount := 0
	bengaliCount := 0
	tamilCount := 0
	teluguCount := 0
	khmerCount := 0
	myanmarCount := 0
	totalCount := 0

	for _, r := range text {
		if unicode.Is(unicode.Arabic, r) {
			arabicCount++
		} else if unicode.Is(unicode.Hebrew, r) {
			hebrewCount++
		} else if unicode.In(r, unicode.Devanagari) {
			devanagariCount++
		} else if unicode.In(r, unicode.Thai) {
			thaiCount++
		} else if unicode.In(r, unicode.Bengali) {
			bengaliCount++
		} else if unicode.In(r, unicode.Tamil) {
			tamilCount++
		} else if unicode.In(r, unicode.Telugu) {
			teluguCount++
		} else if unicode.In(r, unicode.Khmer) {
			khmerCount++
		} else if unicode.In(r, unicode.Myanmar) {
			myanmarCount++
		} else if unicode.In(r, unicode.Han) || unicode.In(r, unicode.Hiragana) || unicode.In(r, unicode.Katakana) || unicode.In(r, unicode.Hangul) {
			cjkCount++
		} else if unicode.In(r, unicode.Cyrillic) {
			cyrillicCount++
		} else if unicode.In(r, unicode.Greek) {
			greekCount++
		}
		totalCount++
	}

	// Determine primary script
	if arabicCount > 0 {
		return ScriptArabic
	} else if hebrewCount > 0 {
		return ScriptHebrew
	} else if devanagariCount > 0 {
		return ScriptDevanagari
	} else if thaiCount > 0 {
		return ScriptThai
	} else if bengaliCount > 0 {
		return ScriptBengali
	} else if tamilCount > 0 {
		return ScriptTamil
	} else if teluguCount > 0 {
		return ScriptTelugu
	} else if khmerCount > 0 {
		return ScriptKhmer
	} else if myanmarCount > 0 {
		return ScriptMyanmar
	} else if cjkCount > 0 {
		// Simplified CJK detection
		return ScriptChinese
	} else if cyrillicCount > 0 {
		return ScriptCyrillic
	} else if greekCount > 0 {
		return ScriptGreek
	}

	return ScriptLatin
}

// shapeRun shapes a single bidirectional run using HarfBuzz
func (ts *TextShaper) shapeRun(run BidiRun) *ShapedText {
	// Use script-specific font if available, otherwise fall back to main font
	face := ts.fontFace
	scriptFont := ts.GetScriptFont(run.Script)
	if scriptFont != nil {
		face = scriptFont.face
	}

	if face == nil || ts.shaper == nil {
		return ts.fallbackShapeRun(run)
	}

	// Configure direction
	direction := di.DirectionLTR
	if run.Direction == TextDirectionRTL {
		direction = di.DirectionRTL
	}

	// Configure script
	var script language.Script
	switch run.Script {
	case ScriptArabic:
		script = language.Arabic
	case ScriptHebrew:
		script = language.Hebrew
	case ScriptDevanagari:
		script = language.Devanagari
	case ScriptThai:
		script = language.Thai
	case ScriptChinese:
		script = language.Han
	case ScriptCyrillic:
		script = language.Cyrillic
	case ScriptGreek:
		script = language.Greek
	case ScriptBengali:
		script = language.Bengali
	case ScriptTamil:
		script = language.Tamil
	case ScriptTelugu:
		script = language.Telugu
	case ScriptKhmer:
		script = language.Khmer
	case ScriptMyanmar:
		script = language.Myanmar
	default:
		script = language.Latin
	}

	// Configure language
	var lang language.Language
	switch run.Script {
	case ScriptBengali:
		lang = language.NewLanguage("BEN")
	case ScriptArabic:
		lang = language.NewLanguage("ARA")
	case ScriptDevanagari:
		lang = language.NewLanguage("HIN")
	case ScriptTamil:
		lang = language.NewLanguage("TAM")
	case ScriptTelugu:
		lang = language.NewLanguage("TEL")
	case ScriptThai:
		lang = language.NewLanguage("THA")
	}

	// Create input for shaping
	runes := []rune(run.Text)
	shapingSize := fixed.Int26_6(ts.fontSize * 64)
	if shapingSize == 0 {
		shapingSize = fixed.I(32) // Default fallback
	}
	input := shaping.Input{
		Text:      runes,
		RunStart:  0,
		RunEnd:    len(runes),
		Direction: direction,
		Face:      face,
		Size:      shapingSize,
		Script:    script,
		Language:  lang,
	}

	// Shape the text
	output := ts.shaper.Shape(input)

	shaped := &ShapedText{
		Direction: run.Direction,
		Glyphs:    make([]ShapedGlyph, len(output.Glyphs)),
	}

	currentX := 0.0
	for i, glyph := range output.Glyphs {
		glyphFont := (*truetype.Font)(nil)
		if scriptFont != nil {
			glyphFont = scriptFont.ttfFont
		}

		// Use HarfBuzz GPOS offsets (already scaled to output units)
		xOffset := float64(glyph.XOffset) / 64.0
		yOffset := float64(glyph.YOffset) / 64.0
		advance := float64(glyph.Advance) / 64.0

		shaped.Glyphs[i] = ShapedGlyph{
			GlyphID:   uint32(glyph.GlyphID),
			X:         currentX + xOffset,
			Y:         yOffset,
			AdvanceX:  advance,
			AdvanceY:  0,
			Cluster:   int(glyph.ClusterIndex),
			Character: runes[glyph.ClusterIndex],
			font:      glyphFont,
		}

		currentX += advance
	}

	shaped.Width = currentX
	shaped.Height = 16.0 // Will be updated with proper font metrics

	return shaped
}

// fallbackShapeText provides simple shaping when HarfBuzz is not available
func (ts *TextShaper) fallbackShapeText(text string) *ShapedText {
	shaped := &ShapedText{
		Direction: ts.Direction,
		Glyphs:    make([]ShapedGlyph, 0),
	}

	x := 0.0
	cluster := 0

	// Use standard font face metrics if available
	hasFont := ts.GoFontFace != nil

	for _, r := range text {
		advance := 12.0 // Default fallback

		if hasFont {
			if adv, ok := ts.GoFontFace.GlyphAdvance(r); ok {
				advance = float64(adv) / 64.0
			}
		}

		glyph := ShapedGlyph{
			GlyphID:   uint32(r),
			X:         x,
			Y:         0,
			AdvanceX:  advance,
			AdvanceY:  0,
			Cluster:   cluster,
			Character: r,
		}

		shaped.Glyphs = append(shaped.Glyphs, glyph)
		x += glyph.AdvanceX
		cluster++
	}

	shaped.Width = x
	shaped.Height = 16

	return shaped
}

// fallbackShapeRun provides simple shaping for a run when HarfBuzz is not available
func (ts *TextShaper) fallbackShapeRun(run BidiRun) *ShapedText {
	shaped := &ShapedText{
		Direction: run.Direction,
		Glyphs:    make([]ShapedGlyph, 0),
	}

	x := 0.0
	cluster := 0

	// Use standard font face metrics if available
	hasFont := ts.GoFontFace != nil

	for _, r := range run.Text {
		advance := 12.0 // Default fallback

		if hasFont {
			if adv, ok := ts.GoFontFace.GlyphAdvance(r); ok {
				advance = float64(adv) / 64.0
			}
		}

		glyph := ShapedGlyph{
			GlyphID:   uint32(r),
			X:         x,
			Y:         0,
			AdvanceX:  advance,
			AdvanceY:  0,
			Cluster:   cluster,
			Character: r,
		}

		shaped.Glyphs = append(shaped.Glyphs, glyph)
		x += glyph.AdvanceX
		cluster++
	}

	shaped.Width = x
	shaped.Height = 16

	return shaped
}

// isRTLChar checks if a character is RTL
func (ts *TextShaper) isRTLChar(r rune) bool {
	return unicode.Is(unicode.Arabic, r) || unicode.Is(unicode.Hebrew, r)
}

// shapeArabicText shapes Arabic text with contextual forms
func (ts *TextShaper) shapeArabicText(text string) *ShapedText {
	shaped := &ShapedText{
		Direction: TextDirectionRTL,
		Glyphs:    make([]ShapedGlyph, 0),
	}

	runes := []rune(text)
	x := 0.0

	for i, r := range runes {
		// Determine contextual form
		form := ts.getArabicForm(runes, i)

		glyph := ShapedGlyph{
			GlyphID:   uint32(r), // Simplified - should map to actual glyph
			X:         x,
			Y:         0,
			AdvanceX:  12, // Simplified advance
			AdvanceY:  0,
			Cluster:   i,
			Character: r,
		}

		// Apply contextual shaping
		glyph.GlyphID = ts.applyArabicShaping(glyph.GlyphID, form)

		shaped.Glyphs = append(shaped.Glyphs, glyph)
		x += glyph.AdvanceX
	}

	shaped.Width = x
	shaped.Height = 16 // Simplified

	return shaped
}

// ArabicForm represents Arabic contextual forms
type ArabicForm int

const (
	ArabicIsolated ArabicForm = iota
	ArabicInitial
	ArabicMedial
	ArabicFinal
)

// getArabicForm determines the contextual form of an Arabic character
func (ts *TextShaper) getArabicForm(runes []rune, index int) ArabicForm {
	if index < 0 || index >= len(runes) {
		return ArabicIsolated
	}

	current := runes[index]
	if !unicode.Is(unicode.Arabic, current) {
		return ArabicIsolated
	}

	canConnectBefore := index > 0 && ts.canConnect(runes[index-1], current)
	canConnectAfter := index < len(runes)-1 && ts.canConnect(current, runes[index+1])

	if canConnectBefore && canConnectAfter {
		return ArabicMedial
	} else if canConnectBefore {
		return ArabicFinal
	} else if canConnectAfter {
		return ArabicInitial
	}

	return ArabicIsolated
}

// canConnect checks if two Arabic characters can connect
func (ts *TextShaper) canConnect(left, right rune) bool {
	// Simplified connection rules
	return unicode.Is(unicode.Arabic, left) && unicode.Is(unicode.Arabic, right)
}

// applyArabicShaping applies contextual shaping to Arabic characters
func (ts *TextShaper) applyArabicShaping(glyphID uint32, form ArabicForm) uint32 {
	// This would map to actual shaped glyphs in a real implementation
	// For now, return the same glyph ID
	return glyphID
}

// shapeHebrewText shapes Hebrew text
func (ts *TextShaper) shapeHebrewText(text string) *ShapedText {
	shaped := &ShapedText{
		Direction: TextDirectionRTL,
		Glyphs:    make([]ShapedGlyph, 0),
	}

	x := 0.0
	cluster := 0

	for _, r := range text {
		glyph := ShapedGlyph{
			GlyphID:   uint32(r),
			X:         x,
			Y:         0,
			AdvanceX:  12,
			AdvanceY:  0,
			Cluster:   cluster,
			Character: r,
		}

		shaped.Glyphs = append(shaped.Glyphs, glyph)
		x += glyph.AdvanceX
		cluster++
	}

	shaped.Width = x
	shaped.Height = 16

	return shaped
}

// shapeDevanagariText shapes Devanagari text with conjuncts
func (ts *TextShaper) shapeDevanagariText(text string) *ShapedText {
	shaped := &ShapedText{
		Direction: TextDirectionLTR,
		Glyphs:    make([]ShapedGlyph, 0),
	}

	// Process conjuncts and reordering
	processed := ts.processDevanagariConjuncts(text)

	x := 0.0
	cluster := 0

	for _, r := range processed {
		glyph := ShapedGlyph{
			GlyphID:   uint32(r),
			X:         x,
			Y:         0,
			AdvanceX:  12,
			AdvanceY:  0,
			Cluster:   cluster,
			Character: r,
		}

		shaped.Glyphs = append(shaped.Glyphs, glyph)
		x += glyph.AdvanceX
		cluster++
	}

	shaped.Width = x
	shaped.Height = 16

	return shaped
}

// processDevanagariConjuncts processes Devanagari conjuncts and reordering
func (ts *TextShaper) processDevanagariConjuncts(text string) string {
	// Simplified Devanagari processing
	// In a real implementation, this would handle:
	// - Conjunct formation
	// - Vowel reordering
	// - Matra positioning
	return text
}

// shapeThaiText shapes Thai text with word breaking
func (ts *TextShaper) shapeThaiText(text string) *ShapedText {
	shaped := &ShapedText{
		Direction: TextDirectionLTR,
		Glyphs:    make([]ShapedGlyph, 0),
	}

	x := 0.0
	cluster := 0

	for _, r := range text {
		glyph := ShapedGlyph{
			GlyphID:   uint32(r),
			X:         x,
			Y:         0,
			AdvanceX:  12,
			AdvanceY:  0,
			Cluster:   cluster,
			Character: r,
		}

		shaped.Glyphs = append(shaped.Glyphs, glyph)
		x += glyph.AdvanceX
		cluster++
	}

	shaped.Width = x
	shaped.Height = 16

	return shaped
}

// shapeLatinText shapes Latin text with ligatures and kerning
func (ts *TextShaper) shapeLatinText(text string) *ShapedText {
	shaped := &ShapedText{
		Direction: TextDirectionLTR,
		Glyphs:    make([]ShapedGlyph, 0),
	}

	// Apply ligatures if enabled
	if ts.EnableLigatures {
		text = ts.applyLatinLigatures(text)
	}

	x := 0.0
	cluster := 0

	runes := []rune(text)
	for i, r := range runes {
		glyph := ShapedGlyph{
			GlyphID:   uint32(r),
			X:         x,
			Y:         0,
			AdvanceX:  12,
			AdvanceY:  0,
			Cluster:   cluster,
			Character: r,
		}

		// Apply kerning if enabled
		if ts.EnableKerning && i > 0 {
			kerning := ts.getKerning(runes[i-1], r)
			glyph.X += kerning
			x += kerning
		}

		shaped.Glyphs = append(shaped.Glyphs, glyph)
		x += glyph.AdvanceX
		cluster++
	}

	shaped.Width = x
	shaped.Height = 16

	return shaped
}

// applyLatinLigatures applies common Latin ligatures
func (ts *TextShaper) applyLatinLigatures(text string) string {
	// Common ligatures
	ligatures := map[string]string{
		"fi":  "ﬁ",
		"fl":  "ﬂ",
		"ff":  "ﬀ",
		"ffi": "ﬃ",
		"ffl": "ﬄ",
	}

	result := text
	for from, to := range ligatures {
		result = strings.ReplaceAll(result, from, to)
	}

	return result
}

// getKerning gets kerning between two characters
func (ts *TextShaper) getKerning(left, right rune) float64 {
	// Simplified kerning table
	kerningPairs := map[string]float64{
		"AV": -1.5,
		"AW": -1.0,
		"AY": -1.5,
		"AT": -1.0,
		"VA": -1.5,
		"WA": -1.0,
		"YA": -1.5,
		"TA": -1.0,
	}

	pair := string([]rune{left, right})
	if kerning, exists := kerningPairs[pair]; exists {
		return kerning
	}

	return 0
}

// Context integration

// SetTextShaper sets the text shaper for the context
func (dc *Context) SetTextShaper(shaper *TextShaper) {
	dc.textShaper = shaper
}

// GetTextShaper returns the current text shaper
func (dc *Context) GetTextShaper() *TextShaper {
	if dc.textShaper == nil {
		dc.textShaper = NewTextShaper()
	}
	return dc.textShaper
}

// DrawShapedString draws shaped text
func (dc *Context) DrawShapedString(text string, x, y float64) {
	shaper := dc.GetTextShaper()
	shaped := shaper.ShapeText(text)

	// Draw each glyph
	for _, glyph := range shaped.Glyphs {
		glyphX := x + glyph.X
		glyphY := y + glyph.Y

		// In a real implementation, this would render the actual glyph
		dc.DrawString(string(glyph.Character), glyphX, glyphY)
	}
}
