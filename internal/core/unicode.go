package core

import (
	"strings"
	"unicode"
)

// Unicode shaping and complex script support

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
)

// TextShaper handles complex text shaping
type TextShaper struct {
	Direction       TextDirection
	Script          ScriptType
	Language        string
	EnableLigatures bool
	EnableKerning   bool
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
}

// ShapeText shapes text according to Unicode rules
func (ts *TextShaper) ShapeText(text string) *ShapedText {
	// Detect script and direction if not set
	if ts.Script == ScriptLatin && ts.Direction == TextDirectionLTR {
		ts.detectScriptAndDirection(text)
	}

	// Apply bidirectional algorithm for mixed scripts
	if ts.containsMixedDirections(text) {
		text = ts.applyBidiAlgorithm(text)
	}

	// Shape the text
	switch ts.Script {
	case ScriptArabic:
		return ts.shapeArabicText(text)
	case ScriptHebrew:
		return ts.shapeHebrewText(text)
	case ScriptDevanagari:
		return ts.shapeDevanagariText(text)
	case ScriptThai:
		return ts.shapeThaiText(text)
	default:
		return ts.shapeLatinText(text)
	}
}

// detectScriptAndDirection detects the primary script and direction
func (ts *TextShaper) detectScriptAndDirection(text string) {
	arabicCount := 0
	hebrewCount := 0
	devanagariCount := 0
	thaiCount := 0
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
		}
		totalCount++
	}

	// Determine primary script
	if arabicCount > totalCount/2 {
		ts.Script = ScriptArabic
		ts.Direction = TextDirectionRTL
	} else if hebrewCount > totalCount/2 {
		ts.Script = ScriptHebrew
		ts.Direction = TextDirectionRTL
	} else if devanagariCount > totalCount/2 {
		ts.Script = ScriptDevanagari
		ts.Direction = TextDirectionLTR
	} else if thaiCount > totalCount/2 {
		ts.Script = ScriptThai
		ts.Direction = TextDirectionLTR
	}
}

// containsMixedDirections checks if text contains mixed directions
func (ts *TextShaper) containsMixedDirections(text string) bool {
	hasLTR := false
	hasRTL := false

	for _, r := range text {
		if unicode.Is(unicode.Arabic, r) || unicode.Is(unicode.Hebrew, r) {
			hasRTL = true
		} else if unicode.IsLetter(r) {
			hasLTR = true
		}

		if hasLTR && hasRTL {
			return true
		}
	}

	return false
}

// applyBidiAlgorithm applies the Unicode Bidirectional Algorithm (simplified)
func (ts *TextShaper) applyBidiAlgorithm(text string) string {
	// This is a simplified implementation
	// In production, use a proper Unicode BiDi implementation

	runes := []rune(text)
	result := make([]rune, len(runes))

	// Find RTL segments and reverse them
	i := 0
	for i < len(runes) {
		if ts.isRTLChar(runes[i]) {
			// Find end of RTL segment
			start := i
			for i < len(runes) && ts.isRTLChar(runes[i]) {
				i++
			}

			// Reverse the RTL segment
			for j := start; j < i; j++ {
				result[j] = runes[i-1-(j-start)]
			}
		} else {
			result[i] = runes[i]
			i++
		}
	}

	return string(result)
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
