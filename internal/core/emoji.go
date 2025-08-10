package core

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"os"
	"path/filepath"
	"strings"

	"github.com/benoitkugler/textlayout/fonts"
)

// Emoji rendering support with color fonts and fallback

// EmojiRenderer handles emoji rendering with proper color font support
type EmojiRenderer struct {
	ColorFonts   []string
	FallbackFont string
	EmojiSize    float64
	EnableSVG    bool
	EnableBitmap bool
	Cache        map[string]*image.RGBA

	// Color font support
	colorFontFaces map[string]fonts.Face
	loadedFonts    []string
}

// NewEmojiRenderer creates a new emoji renderer
func NewEmojiRenderer() *EmojiRenderer {
	er := &EmojiRenderer{
		ColorFonts: []string{
			"Apple Color Emoji",
			"Segoe UI Emoji",
			"Noto Color Emoji",
			"Android Emoji",
		},
		FallbackFont:   "Arial",
		EmojiSize:      16,
		EnableSVG:      true,
		EnableBitmap:   true,
		Cache:          make(map[string]*image.RGBA),
		colorFontFaces: make(map[string]fonts.Face),
		loadedFonts:    make([]string, 0),
	}

	// Try to load system emoji fonts
	er.loadSystemEmojiFonts()

	return er
}

// EmojiInfo represents information about an emoji
type EmojiInfo struct {
	Codepoint   string
	Name        string
	Category    string
	Subcategory string
	Keywords    []string
	SkinTones   []string
	ZWJSequence bool
	Modifiable  bool
}

// Common emoji categories
const (
	EmojiCategorySmileys    = "Smileys & Emotion"
	EmojiCategoryPeople     = "People & Body"
	EmojiCategoryAnimals    = "Animals & Nature"
	EmojiCategoryFood       = "Food & Drink"
	EmojiCategoryActivities = "Activities"
	EmojiCategoryTravel     = "Travel & Places"
	EmojiCategoryObjects    = "Objects"
	EmojiCategorySymbols    = "Symbols"
	EmojiCategoryFlags      = "Flags"
)

// IsEmoji checks if a rune is an emoji
func IsEmoji(r rune) bool {
	// Basic emoji ranges
	return (r >= 0x1F600 && r <= 0x1F64F) || // Emoticons
		(r >= 0x1F300 && r <= 0x1F5FF) || // Misc Symbols and Pictographs
		(r >= 0x1F680 && r <= 0x1F6FF) || // Transport and Map
		(r >= 0x1F1E6 && r <= 0x1F1FF) || // Regional indicators
		(r >= 0x2600 && r <= 0x26FF) || // Misc symbols
		(r >= 0x2700 && r <= 0x27BF) || // Dingbats
		(r >= 0xFE00 && r <= 0xFE0F) || // Variation selectors
		r == 0x200D // Zero width joiner
}

// ParseEmojiSequence parses an emoji sequence including ZWJ sequences
func (er *EmojiRenderer) ParseEmojiSequence(text string) []EmojiSequence {
	var sequences []EmojiSequence
	runes := []rune(text)

	i := 0
	for i < len(runes) {
		if IsEmoji(runes[i]) {
			sequence := er.extractEmojiSequence(runes, i)
			sequences = append(sequences, sequence)
			i += len(sequence.Runes)
		} else {
			i++
		}
	}

	return sequences
}

// EmojiSequence represents a complete emoji sequence
type EmojiSequence struct {
	Runes       []rune
	Text        string
	IsZWJ       bool
	HasModifier bool
	SkinTone    string
	Category    string
}

// extractEmojiSequence extracts a complete emoji sequence
func (er *EmojiRenderer) extractEmojiSequence(runes []rune, start int) EmojiSequence {
	sequence := EmojiSequence{
		Runes: make([]rune, 0),
	}

	i := start
	for i < len(runes) && (IsEmoji(runes[i]) || er.isModifier(runes[i])) {
		sequence.Runes = append(sequence.Runes, runes[i])

		// Check for ZWJ sequence
		if runes[i] == 0x200D {
			sequence.IsZWJ = true
		}

		// Check for skin tone modifiers
		if er.isSkinToneModifier(runes[i]) {
			sequence.HasModifier = true
			sequence.SkinTone = er.getSkinToneName(runes[i])
		}

		i++

		// Stop if we hit a non-emoji, non-modifier character
		if i < len(runes) && !IsEmoji(runes[i]) && !er.isModifier(runes[i]) {
			break
		}
	}

	sequence.Text = string(sequence.Runes)
	sequence.Category = er.GetEmojiCategory(sequence.Runes[0])

	return sequence
}

// isModifier checks if a rune is an emoji modifier
func (er *EmojiRenderer) isModifier(r rune) bool {
	return (r >= 0x1F3FB && r <= 0x1F3FF) || // Skin tone modifiers
		(r >= 0xFE00 && r <= 0xFE0F) || // Variation selectors
		r == 0x200D // Zero width joiner
}

// isSkinToneModifier checks if a rune is a skin tone modifier
func (er *EmojiRenderer) isSkinToneModifier(r rune) bool {
	return r >= 0x1F3FB && r <= 0x1F3FF
}

// getSkinToneName gets the name of a skin tone modifier
func (er *EmojiRenderer) getSkinToneName(r rune) string {
	switch r {
	case 0x1F3FB:
		return "light"
	case 0x1F3FC:
		return "medium-light"
	case 0x1F3FD:
		return "medium"
	case 0x1F3FE:
		return "medium-dark"
	case 0x1F3FF:
		return "dark"
	default:
		return ""
	}
}

// GetEmojiCategory gets the category of an emoji
func (er *EmojiRenderer) GetEmojiCategory(r rune) string {
	switch {
	case r >= 0x1F600 && r <= 0x1F64F:
		return EmojiCategorySmileys
	case r >= 0x1F466 && r <= 0x1F487:
		return EmojiCategoryPeople
	case r >= 0x1F400 && r <= 0x1F43F:
		return EmojiCategoryAnimals
	case r >= 0x1F32D && r <= 0x1F37F:
		return EmojiCategoryFood
	case r >= 0x1F3A0 && r <= 0x1F3FF:
		return EmojiCategoryActivities
	case r >= 0x1F680 && r <= 0x1F6FF:
		return EmojiCategoryTravel
	case r >= 0x1F4A0 && r <= 0x1F4FF:
		return EmojiCategoryObjects
	case r >= 0x1F500 && r <= 0x1F5FF:
		return EmojiCategorySymbols
	case r >= 0x1F1E6 && r <= 0x1F1FF:
		return EmojiCategoryFlags
	default:
		return EmojiCategorySymbols
	}
}

// RenderEmoji renders an emoji sequence
func (er *EmojiRenderer) RenderEmoji(sequence EmojiSequence, size float64) *image.RGBA {
	// Check cache first
	cacheKey := sequence.Text + "_" + string(rune(int(size)))
	if cached, exists := er.Cache[cacheKey]; exists {
		return cached
	}

	var result *image.RGBA

	// Try color font rendering first
	if er.EnableSVG || er.EnableBitmap {
		result = er.renderColorEmoji(sequence, size)
	}

	// Fallback to text rendering
	if result == nil {
		result = er.renderFallbackEmoji(sequence, size)
	}

	// Cache the result
	if result != nil {
		er.Cache[cacheKey] = result
	}

	return result
}

// renderColorEmoji renders emoji using color fonts
func (er *EmojiRenderer) renderColorEmoji(sequence EmojiSequence, size float64) *image.RGBA {
	// This would integrate with actual color font rendering
	// For now, create a placeholder colored emoji

	img := image.NewRGBA(image.Rect(0, 0, int(size), int(size)))

	// Create a simple colored representation based on category
	var baseColor color.RGBA
	switch sequence.Category {
	case EmojiCategorySmileys:
		baseColor = color.RGBA{255, 220, 100, 255} // Yellow
	case EmojiCategoryPeople:
		baseColor = color.RGBA{255, 200, 150, 255} // Skin tone
	case EmojiCategoryAnimals:
		baseColor = color.RGBA{150, 100, 50, 255} // Brown
	case EmojiCategoryFood:
		baseColor = color.RGBA{255, 150, 100, 255} // Orange
	case EmojiCategoryActivities:
		baseColor = color.RGBA{100, 150, 255, 255} // Blue
	case EmojiCategoryTravel:
		baseColor = color.RGBA{100, 255, 150, 255} // Green
	case EmojiCategoryObjects:
		baseColor = color.RGBA{200, 200, 200, 255} // Gray
	case EmojiCategorySymbols:
		baseColor = color.RGBA{255, 100, 255, 255} // Magenta
	case EmojiCategoryFlags:
		baseColor = color.RGBA{255, 100, 100, 255} // Red
	default:
		baseColor = color.RGBA{128, 128, 128, 255} // Default gray
	}

	// Apply skin tone modification if present
	if sequence.HasModifier {
		baseColor = er.applySkinToneModification(baseColor, sequence.SkinTone)
	}

	// Draw a simple emoji representation
	er.drawSimpleEmoji(img, baseColor, sequence, size)

	return img
}

// applySkinToneModification applies skin tone modification to color
func (er *EmojiRenderer) applySkinToneModification(base color.RGBA, skinTone string) color.RGBA {
	switch skinTone {
	case "light":
		return color.RGBA{255, 220, 177, base.A}
	case "medium-light":
		return color.RGBA{240, 195, 140, base.A}
	case "medium":
		return color.RGBA{198, 134, 66, base.A}
	case "medium-dark":
		return color.RGBA{140, 85, 33, base.A}
	case "dark":
		return color.RGBA{92, 51, 23, base.A}
	default:
		return base
	}
}

// drawSimpleEmoji draws a simple emoji representation
func (er *EmojiRenderer) drawSimpleEmoji(img *image.RGBA, baseColor color.RGBA, sequence EmojiSequence, size float64) {
	bounds := img.Bounds()
	center := image.Point{bounds.Dx() / 2, bounds.Dy() / 2}
	radius := int(size / 3)

	// Draw based on emoji type
	if sequence.Category == EmojiCategorySmileys {
		// Draw a simple smiley face
		er.drawCircle(img, center, radius, baseColor)

		// Eyes
		eyeColor := color.RGBA{0, 0, 0, 255}
		leftEye := image.Point{center.X - radius/3, center.Y - radius/3}
		rightEye := image.Point{center.X + radius/3, center.Y - radius/3}
		er.drawCircle(img, leftEye, radius/6, eyeColor)
		er.drawCircle(img, rightEye, radius/6, eyeColor)

		// Smile
		er.drawArc(img, center, radius/2, 0, 180, eyeColor)
	} else {
		// Draw a simple colored circle for other categories
		er.drawCircle(img, center, radius, baseColor)
	}
}

// drawCircle draws a filled circle
func (er *EmojiRenderer) drawCircle(img *image.RGBA, center image.Point, radius int, c color.RGBA) {
	bounds := img.Bounds()

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			dx := x - center.X
			dy := y - center.Y
			if dx*dx+dy*dy <= radius*radius {
				img.SetRGBA(x, y, c)
			}
		}
	}
}

// drawArc draws an arc (simplified)
func (er *EmojiRenderer) drawArc(img *image.RGBA, center image.Point, radius int, startAngle, endAngle float64, c color.RGBA) {
	// Simplified arc drawing for smile
	for y := center.Y; y < center.Y+radius/2; y++ {
		for x := center.X - radius/2; x < center.X+radius/2; x++ {
			dx := x - center.X
			dy := y - center.Y
			if dx*dx+dy*dy >= (radius/3)*(radius/3) && dx*dx+dy*dy <= (radius/2)*(radius/2) {
				img.SetRGBA(x, y, c)
			}
		}
	}
}

// renderFallbackEmoji renders emoji as text fallback
func (er *EmojiRenderer) renderFallbackEmoji(sequence EmojiSequence, size float64) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, int(size), int(size)))

	// Fill with white background
	draw.Draw(img, img.Bounds(), &image.Uniform{color.White}, image.Point{}, draw.Src)

	// Draw the emoji text (simplified)
	// In a real implementation, this would use proper text rendering
	textColor := color.RGBA{0, 0, 0, 255}
	er.drawText(img, sequence.Text, textColor, size)

	return img
}

// drawText draws text on image (simplified)
func (er *EmojiRenderer) drawText(img *image.RGBA, text string, c color.RGBA, size float64) {
	// This is a very simplified text drawing
	// In a real implementation, use proper font rendering
	bounds := img.Bounds()

	// Just draw a simple representation
	centerX := bounds.Dx() / 2
	centerY := bounds.Dy() / 2

	// Draw a simple square to represent text
	for y := centerY - int(size/4); y < centerY+int(size/4); y++ {
		for x := centerX - int(size/4); x < centerX+int(size/4); x++ {
			if x >= 0 && x < bounds.Dx() && y >= 0 && y < bounds.Dy() {
				img.SetRGBA(x, y, c)
			}
		}
	}
}

// GetEmojiInfo gets information about an emoji
func (er *EmojiRenderer) GetEmojiInfo(emojiText string) *EmojiInfo {
	sequences := er.ParseEmojiSequence(emojiText)
	if len(sequences) == 0 {
		return nil
	}

	sequence := sequences[0]

	return &EmojiInfo{
		Codepoint:   er.getCodepoint(sequence.Runes),
		Name:        er.getEmojiName(sequence.Runes[0]),
		Category:    sequence.Category,
		Keywords:    er.getEmojiKeywords(sequence.Runes[0]),
		SkinTones:   er.getSupportedSkinTones(sequence.Runes[0]),
		ZWJSequence: sequence.IsZWJ,
		Modifiable:  er.isModifiable(sequence.Runes[0]),
	}
}

// getCodepoint gets the Unicode codepoint string
func (er *EmojiRenderer) getCodepoint(runes []rune) string {
	var parts []string
	for _, r := range runes {
		parts = append(parts, string(rune(r)))
	}
	return strings.Join(parts, " ")
}

// getEmojiName gets the name of an emoji
func (er *EmojiRenderer) getEmojiName(r rune) string {
	// Simplified emoji names
	names := map[rune]string{
		0x1F600: "grinning face",
		0x1F601: "beaming face with smiling eyes",
		0x1F602: "face with tears of joy",
		0x1F603: "grinning face with big eyes",
		0x1F604: "grinning face with smiling eyes",
		0x1F605: "grinning face with sweat",
		0x1F606: "grinning squinting face",
		0x1F607: "smiling face with halo",
		// Add more as needed
	}

	if name, exists := names[r]; exists {
		return name
	}

	return "unknown emoji"
}

// getEmojiKeywords gets keywords for an emoji
func (er *EmojiRenderer) getEmojiKeywords(r rune) []string {
	// Simplified keyword mapping
	switch {
	case r >= 0x1F600 && r <= 0x1F64F:
		return []string{"face", "emotion", "smile"}
	case r >= 0x1F400 && r <= 0x1F43F:
		return []string{"animal", "nature"}
	case r >= 0x1F32D && r <= 0x1F37F:
		return []string{"food", "drink"}
	default:
		return []string{"emoji"}
	}
}

// getSupportedSkinTones gets supported skin tones for an emoji
func (er *EmojiRenderer) getSupportedSkinTones(r rune) []string {
	// Check if emoji supports skin tone modifiers
	if er.isModifiable(r) {
		return []string{"light", "medium-light", "medium", "medium-dark", "dark"}
	}
	return []string{}
}

// isModifiable checks if an emoji supports skin tone modification
func (er *EmojiRenderer) isModifiable(r rune) bool {
	// People emojis typically support skin tone modification
	return (r >= 0x1F466 && r <= 0x1F487) || // People
		(r >= 0x1F574 && r <= 0x1F575) || // Detective, etc.
		(r >= 0x1F57A && r <= 0x1F57A) || // Man dancing
		(r >= 0x1F590 && r <= 0x1F590) // Hand with fingers splayed
}

// Context integration

// SetEmojiRenderer sets the emoji renderer for the context
func (dc *Context) SetEmojiRenderer(renderer *EmojiRenderer) {
	dc.emojiRenderer = renderer
}

// GetEmojiRenderer returns the current emoji renderer
func (dc *Context) GetEmojiRenderer() *EmojiRenderer {
	if dc.emojiRenderer == nil {
		dc.emojiRenderer = NewEmojiRenderer()
	}
	return dc.emojiRenderer
}

// DrawStringWithEmoji draws text with emoji support
func (dc *Context) DrawStringWithEmoji(text string, x, y float64) {
	renderer := dc.GetEmojiRenderer()

	// Parse text for emoji sequences
	runes := []rune(text)
	currentX := x

	i := 0
	for i < len(runes) {
		if IsEmoji(runes[i]) {
			// Extract and render emoji
			sequence := renderer.extractEmojiSequence(runes, i)
			emojiImg := renderer.RenderEmoji(sequence, renderer.EmojiSize)

			if emojiImg != nil {
				dc.DrawImage(emojiImg, int(currentX), int(y-renderer.EmojiSize))
			}

			currentX += renderer.EmojiSize
			i += len(sequence.Runes)
		} else {
			// Regular text character
			dc.DrawString(string(runes[i]), currentX, y)
			currentX += 8 // Simplified character width
			i++
		}
	}
}

// loadSystemEmojiFonts attempts to load system emoji fonts
func (er *EmojiRenderer) loadSystemEmojiFonts() {
	// Common system font paths for emoji fonts
	fontPaths := []string{
		// macOS
		"/System/Library/Fonts/Apple Color Emoji.ttc",
		"/Library/Fonts/Apple Color Emoji.ttc",
		// Windows
		"C:/Windows/Fonts/seguiemj.ttf", // Segoe UI Emoji
		// Linux
		"/usr/share/fonts/truetype/noto-color-emoji/NotoColorEmoji.ttf",
		"/usr/share/fonts/noto-color-emoji/NotoColorEmoji.ttf",
		"/usr/share/fonts/TTF/NotoColorEmoji.ttf",
		// Android
		"/system/fonts/NotoColorEmoji.ttf",
	}

	for _, fontPath := range fontPaths {
		if er.loadEmojiFont(fontPath) {
			fmt.Printf("Loaded emoji font: %s\n", fontPath)
		}
	}
}

// loadEmojiFont loads a single emoji font file
func (er *EmojiRenderer) loadEmojiFont(fontPath string) bool {
	// Check if file exists
	if _, err := os.Stat(fontPath); os.IsNotExist(err) {
		return false
	}

	// Read font file (for future COLR/CPAL parsing)
	_, err := os.Stat(fontPath) // Just verify it's readable for now
	if err != nil {
		return false
	}

	// Parse font - for now, just mark as loaded without parsing
	// TODO: Implement proper COLR/CPAL font parsing for color emoji
	// The textlayout library has different font interfaces than expected

	// Store the font path for now
	fontName := filepath.Base(fontPath)
	er.loadedFonts = append(er.loadedFonts, fontName)

	return true
}

// hasColorFont checks if any color fonts are loaded
func (er *EmojiRenderer) hasColorFont() bool {
	return len(er.loadedFonts) > 0
}

// getColorFont gets the first available color font
func (er *EmojiRenderer) getColorFont() fonts.Face {
	if len(er.loadedFonts) > 0 {
		return er.colorFontFaces[er.loadedFonts[0]]
	}
	return nil
}
