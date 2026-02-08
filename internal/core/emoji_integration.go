package core

import (
	"image"

	"golang.org/x/image/font"
)

// TextSegment represents a segment of text (either regular text or emoji)
type TextSegment struct {
	Text     string
	IsEmoji  bool
	Sequence EmojiSequence
	Width    float64
}

// splitTextAndEmoji splits text into alternating text and emoji segments
func splitTextAndEmoji(text string) []TextSegment {
	var segments []TextSegment
	runes := []rune(text)
	i := 0

	for i < len(runes) {
		if IsEmoji(runes[i]) {
			// Extract emoji sequence
			renderer := NewEmojiRenderer()
			sequence := renderer.extractEmojiSequence(runes, i)
			segments = append(segments, TextSegment{
				Text:     sequence.Text,
				IsEmoji:  true,
				Sequence: sequence,
			})
			i += len(sequence.Runes)
		} else {
			// Extract regular text until next emoji
			start := i
			for i < len(runes) && !IsEmoji(runes[i]) {
				i++
			}
			segments = append(segments, TextSegment{
				Text:    string(runes[start:i]),
				IsEmoji: false,
			})
		}
	}

	return segments
}

// measureEmojiWidth calculates the width of an emoji at a given size
func measureEmojiWidth(sequence EmojiSequence, size float64) float64 {
	// Emoji width is typically equal to its height (square aspect ratio)
	return size
}

// alignEmojiBaseline aligns emoji image to text baseline
func alignEmojiBaseline(emojiImg *image.RGBA, textBaseline, emojiSize float64) (offsetX, offsetY float64) {
	// Emoji should be centered vertically around the text baseline
	// Typically emoji sits slightly above baseline
	offsetX = 0
	offsetY = -emojiSize * 0.85 // Adjust so emoji aligns with text
	return
}

// measureMixedString measures string width including emojis
func (dc *Context) measureMixedString(text string, emojiSize float64) (width, height float64) {
	segments := splitTextAndEmoji(text)
	totalWidth := 0.0
	maxHeight := dc.fontHeight

	for _, seg := range segments {
		if seg.IsEmoji {
			w := measureEmojiWidth(seg.Sequence, emojiSize)
			totalWidth += w
			if emojiSize > maxHeight {
				maxHeight = emojiSize
			}
		} else {
			// Measure text segment
			d := &font.Drawer{Face: dc.fontFace}
			a := d.MeasureString(seg.Text)
			totalWidth += float64(a >> 6)
		}
	}

	return totalWidth, maxHeight
}

// drawMixedString draws text with emoji support
func (dc *Context) drawMixedString(dst *image.RGBA, text string, x, y float64) {
	if !dc.enableAutoEmoji {
		// Fall back to regular text rendering
		dc.drawString(dst, text, x, y)
		return
	}

	segments := splitTextAndEmoji(text)
	currentX := x
	renderer := dc.GetEmojiRenderer()
	emojiSize := dc.fontHeight // Use font height as emoji size

	for _, seg := range segments {
		if seg.IsEmoji {
			// Render emoji
			emojiImg := renderer.RenderEmoji(seg.Sequence, emojiSize)
			if emojiImg != nil {
				_, offsetY := alignEmojiBaseline(emojiImg, y, emojiSize)
				dc.drawImageAt(dst, emojiImg, int(currentX), int(y+offsetY))
			}
			currentX += measureEmojiWidth(seg.Sequence, emojiSize)
		} else {
			// Render text
			dc.drawString(dst, seg.Text, currentX, y)
			// Measure text width for next position
			d := &font.Drawer{Face: dc.fontFace}
			a := d.MeasureString(seg.Text)
			currentX += float64(a >> 6)
		}
	}
}

// drawImageAt draws an image at a specific position (helper)
func (dc *Context) drawImageAt(dst *image.RGBA, src image.Image, x, y int) {
	bounds := src.Bounds()
	for dy := 0; dy < bounds.Dy(); dy++ {
		for dx := 0; dx < bounds.Dx(); dx++ {
			px := x + dx
			py := y + dy
			if px >= 0 && px < dst.Bounds().Dx() && py >= 0 && py < dst.Bounds().Dy() {
				srcColor := src.At(bounds.Min.X+dx, bounds.Min.Y+dy)
				dst.Set(px, py, srcColor)
			}
		}
	}
}
