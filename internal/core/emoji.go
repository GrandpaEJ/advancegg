package core

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"os"
	"strings"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
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

	// Font support for actual emoji rendering
	emojiFont     *truetype.Font
	emojiFontFace font.Face
	fontLoaded    bool
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
		FallbackFont: "Arial",
		EmojiSize:    16,
		EnableSVG:    true,
		EnableBitmap: true,
		Cache:        make(map[string]*image.RGBA),
		fontLoaded:   false,
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

// renderColorEmoji renders emoji using actual font glyphs
func (er *EmojiRenderer) renderColorEmoji(sequence EmojiSequence, size float64) *image.RGBA {
	// Try to render using actual emoji font first
	if er.fontLoaded && er.emojiFontFace != nil {
		return er.renderWithFont(sequence, size)
	}

	// Fallback to simple colored representation
	return er.renderSimpleEmoji(sequence, size)
}

// renderWithFont renders emoji using the loaded emoji font
func (er *EmojiRenderer) renderWithFont(sequence EmojiSequence, size float64) *image.RGBA {
	// For now, always fall back to simple rendering since color emoji fonts
	// require special handling that freetype doesn't support well
	// TODO: Implement proper COLR/CPAL color emoji font rendering
	return er.renderSimpleEmoji(sequence, size)

	/*
		// This code would work for regular fonts but not color emoji fonts
		img := image.NewRGBA(image.Rect(0, 0, int(size), int(size)))

		// Fill with transparent background
		draw.Draw(img, img.Bounds(), &image.Uniform{color.RGBA{0, 0, 0, 0}}, image.ZP, draw.Src)

		// Create freetype context
		c := freetype.NewContext()
		c.SetDPI(72)
		c.SetFont(er.emojiFont)
		c.SetFontSize(size * 0.8) // Slightly smaller than the canvas
		c.SetClip(img.Bounds())
		c.SetDst(img)
		c.SetSrc(image.NewUniform(color.RGBA{0, 0, 0, 255})) // Black color for now

		// Calculate position to center the emoji
		pt := freetype.Pt(int(size*0.1), int(size*0.8))

		// Draw the emoji text
		_, err := c.DrawString(sequence.Text, pt)
		if err != nil {
			// If font rendering fails, fall back to simple rendering
			return er.renderSimpleEmoji(sequence, size)
		}

		return img
	*/
}

// renderSimpleEmoji creates a simple colored emoji representation
func (er *EmojiRenderer) renderSimpleEmoji(sequence EmojiSequence, size float64) *image.RGBA {
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

// drawSimpleEmoji draws a recognizable emoji representation
func (er *EmojiRenderer) drawSimpleEmoji(img *image.RGBA, baseColor color.RGBA, sequence EmojiSequence, size float64) {
	bounds := img.Bounds()
	center := image.Point{bounds.Dx() / 2, bounds.Dy() / 2}
	radius := int(size / 3)

	// Draw based on specific emoji if we can recognize it
	if len(sequence.Runes) > 0 {
		emoji := sequence.Runes[0]
		switch emoji {
		case 0x1F600: // 😀 grinning face
			er.drawGrinningFace(img, center, radius)
			return
		case 0x1F603: // 😃 grinning face with big eyes
			er.drawGrinningFace(img, center, radius) // Same as grinning for now
			return
		case 0x1F604: // 😄 grinning face with smiling eyes
			er.drawGrinningFace(img, center, radius) // Same as grinning for now
			return
		case 0x1F44B: // 👋 waving hand
			er.drawWavingHand(img, center, radius)
			return
		case 0x1F44D: // 👍 thumbs up
			er.drawThumbsUp(img, center, radius)
			return
		case 0x2764: // ❤ red heart
			er.drawHeart(img, center, radius)
			return
		case 0x1F31F: // 🌟 glowing star
			er.drawStar(img, center, radius)
			return
		}
	}

	// Fallback to category-based rendering
	switch sequence.Category {
	case EmojiCategorySmileys:
		er.drawGrinningFace(img, center, radius)
	case EmojiCategoryPeople:
		er.drawGenericPerson(img, center, radius, baseColor)
	case EmojiCategoryAnimals:
		er.drawGenericAnimal(img, center, radius, baseColor)
	default:
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
		"/usr/share/fonts/truetype/noto/NotoColorEmoji.ttf",
		// Android
		"/system/fonts/NotoColorEmoji.ttf",
	}

	for _, fontPath := range fontPaths {
		if er.loadEmojiFont(fontPath) {
			fmt.Printf("Loaded emoji font: %s\n", fontPath)
			break // Only load the first available font
		}
	}
}

// loadEmojiFont loads a single emoji font file
func (er *EmojiRenderer) loadEmojiFont(fontPath string) bool {
	// Check if file exists
	if _, err := os.Stat(fontPath); os.IsNotExist(err) {
		return false
	}

	// Read font file
	fontData, err := os.ReadFile(fontPath)
	if err != nil {
		return false
	}

	// Parse font using freetype
	font, err := truetype.Parse(fontData)
	if err != nil {
		// Many emoji fonts are in special formats that freetype can't parse
		// This is expected for color emoji fonts like Noto Color Emoji
		return false
	}

	// Create font face
	er.emojiFont = font
	er.emojiFontFace = truetype.NewFace(font, &truetype.Options{
		Size: er.EmojiSize,
		DPI:  72,
	})
	er.fontLoaded = true

	return true
}

// hasColorFont checks if any color fonts are loaded
func (er *EmojiRenderer) hasColorFont() bool {
	return er.fontLoaded
}

// getEmojiFont gets the loaded emoji font face
func (er *EmojiRenderer) getEmojiFont() font.Face {
	if er.fontLoaded {
		return er.emojiFontFace
	}
	return nil
}

// Specific emoji drawing methods

// drawGrinningFace draws a grinning face emoji 😀
func (er *EmojiRenderer) drawGrinningFace(img *image.RGBA, center image.Point, radius int) {
	// Yellow face
	faceColor := color.RGBA{255, 220, 100, 255}
	er.drawCircle(img, center, radius, faceColor)

	// Black eyes
	eyeColor := color.RGBA{0, 0, 0, 255}
	leftEye := image.Point{center.X - radius/3, center.Y - radius/3}
	rightEye := image.Point{center.X + radius/3, center.Y - radius/3}
	er.drawCircle(img, leftEye, radius/8, eyeColor)
	er.drawCircle(img, rightEye, radius/8, eyeColor)

	// Big smile
	er.drawSmile(img, center, radius, eyeColor)
}

// drawWavingHand draws a waving hand emoji 👋
func (er *EmojiRenderer) drawWavingHand(img *image.RGBA, center image.Point, radius int) {
	// Skin color
	handColor := color.RGBA{255, 200, 150, 255}

	// Draw palm (oval)
	er.drawOval(img, center, radius, radius*3/4, handColor)

	// Draw fingers
	fingerColor := handColor
	for i := 0; i < 4; i++ {
		fingerX := center.X - radius/2 + i*radius/4
		fingerY := center.Y - radius
		finger := image.Point{fingerX, fingerY}
		er.drawOval(img, finger, radius/6, radius/3, fingerColor)
	}
}

// drawThumbsUp draws a thumbs up emoji 👍
func (er *EmojiRenderer) drawThumbsUp(img *image.RGBA, center image.Point, radius int) {
	// Skin color
	handColor := color.RGBA{255, 200, 150, 255}

	// Draw thumb (vertical oval)
	thumbCenter := image.Point{center.X, center.Y - radius/4}
	er.drawOval(img, thumbCenter, radius/3, radius, handColor)

	// Draw fist base
	fistCenter := image.Point{center.X, center.Y + radius/2}
	er.drawOval(img, fistCenter, radius*2/3, radius/2, handColor)
}

// drawHeart draws a heart emoji ❤
func (er *EmojiRenderer) drawHeart(img *image.RGBA, center image.Point, radius int) {
	heartColor := color.RGBA{255, 50, 50, 255}

	// Draw two circles for the top of the heart
	leftTop := image.Point{center.X - radius/3, center.Y - radius/4}
	rightTop := image.Point{center.X + radius/3, center.Y - radius/4}
	er.drawCircle(img, leftTop, radius/2, heartColor)
	er.drawCircle(img, rightTop, radius/2, heartColor)

	// Draw triangle for bottom of heart
	er.drawTriangle(img, center, radius, heartColor)
}

// drawStar draws a star emoji 🌟
func (er *EmojiRenderer) drawStar(img *image.RGBA, center image.Point, radius int) {
	starColor := color.RGBA{255, 255, 100, 255}

	// Draw a simple 5-pointed star
	er.drawStarShape(img, center, radius, starColor)
}

// drawGenericPerson draws a generic person emoji
func (er *EmojiRenderer) drawGenericPerson(img *image.RGBA, center image.Point, radius int, baseColor color.RGBA) {
	// Head
	headColor := color.RGBA{255, 200, 150, 255}
	headCenter := image.Point{center.X, center.Y - radius/2}
	er.drawCircle(img, headCenter, radius/2, headColor)

	// Body
	bodyCenter := image.Point{center.X, center.Y + radius/4}
	er.drawOval(img, bodyCenter, radius/2, radius/2, baseColor)
}

// drawGenericAnimal draws a generic animal emoji
func (er *EmojiRenderer) drawGenericAnimal(img *image.RGBA, center image.Point, radius int, baseColor color.RGBA) {
	// Body
	er.drawCircle(img, center, radius, baseColor)

	// Ears
	earColor := color.RGBA{baseColor.R - 30, baseColor.G - 30, baseColor.B - 30, 255}
	leftEar := image.Point{center.X - radius/2, center.Y - radius}
	rightEar := image.Point{center.X + radius/2, center.Y - radius}
	er.drawCircle(img, leftEar, radius/4, earColor)
	er.drawCircle(img, rightEar, radius/4, earColor)

	// Eyes
	eyeColor := color.RGBA{0, 0, 0, 255}
	leftEye := image.Point{center.X - radius/4, center.Y - radius/4}
	rightEye := image.Point{center.X + radius/4, center.Y - radius/4}
	er.drawCircle(img, leftEye, radius/8, eyeColor)
	er.drawCircle(img, rightEye, radius/8, eyeColor)
}

// Helper drawing methods

// drawSmile draws a smile arc
func (er *EmojiRenderer) drawSmile(img *image.RGBA, center image.Point, radius int, c color.RGBA) {
	// Draw a simple smile as a curved line
	smileY := center.Y + radius/4
	for x := center.X - radius/2; x <= center.X+radius/2; x++ {
		// Simple parabolic curve for smile
		dx := float64(x - center.X)
		y := smileY + int(dx*dx/float64(radius*2))
		if y >= 0 && y < img.Bounds().Dy() && x >= 0 && x < img.Bounds().Dx() {
			img.SetRGBA(x, y, c)
			// Make it thicker
			if y+1 < img.Bounds().Dy() {
				img.SetRGBA(x, y+1, c)
			}
		}
	}
}

// drawOval draws a filled oval
func (er *EmojiRenderer) drawOval(img *image.RGBA, center image.Point, radiusX, radiusY int, c color.RGBA) {
	bounds := img.Bounds()
	for y := center.Y - radiusY; y <= center.Y+radiusY; y++ {
		for x := center.X - radiusX; x <= center.X+radiusX; x++ {
			if x >= bounds.Min.X && x < bounds.Max.X && y >= bounds.Min.Y && y < bounds.Max.Y {
				dx := float64(x - center.X)
				dy := float64(y - center.Y)
				// Ellipse equation: (x/a)² + (y/b)² <= 1
				if (dx*dx)/float64(radiusX*radiusX)+(dy*dy)/float64(radiusY*radiusY) <= 1 {
					img.SetRGBA(x, y, c)
				}
			}
		}
	}
}

// drawTriangle draws a filled triangle pointing down (for heart bottom)
func (er *EmojiRenderer) drawTriangle(img *image.RGBA, center image.Point, radius int, c color.RGBA) {
	bounds := img.Bounds()
	// Simple triangle pointing down
	for y := center.Y; y <= center.Y+radius; y++ {
		width := radius - (y - center.Y)
		for x := center.X - width; x <= center.X+width; x++ {
			if x >= bounds.Min.X && x < bounds.Max.X && y >= bounds.Min.Y && y < bounds.Max.Y {
				img.SetRGBA(x, y, c)
			}
		}
	}
}

// drawStarShape draws a 5-pointed star
func (er *EmojiRenderer) drawStarShape(img *image.RGBA, center image.Point, radius int, c color.RGBA) {
	// Simple star - draw as overlapping triangles
	// This is a simplified star shape
	bounds := img.Bounds()

	// Draw a diamond shape as a simple star
	for y := center.Y - radius; y <= center.Y+radius; y++ {
		for x := center.X - radius; x <= center.X+radius; x++ {
			if x >= bounds.Min.X && x < bounds.Max.X && y >= bounds.Min.Y && y < bounds.Max.Y {
				dx := abs(x - center.X)
				dy := abs(y - center.Y)
				// Diamond shape
				if dx+dy <= radius {
					img.SetRGBA(x, y, c)
				}
			}
		}
	}
}

// abs returns absolute value of an integer
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
