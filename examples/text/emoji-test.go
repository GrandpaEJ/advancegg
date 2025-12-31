package main

import (
	"fmt"
	"github.com/GrandpaEJ/advancegg/internal/core"
)

func main() {
	// Create a simple test to show the new emoji rendering
	dc := core.NewContext(400, 200)
	
	// White background
	dc.SetRGB(1, 1, 1)
	dc.Clear()
	
	// Title
	dc.SetRGB(0, 0, 0)
	dc.DrawString("Real Emoji Test", 20, 30)
	
	// Set up emoji renderer
	emojiRenderer := core.NewEmojiRenderer()
	emojiRenderer.EmojiSize = 48
	
	// Test specific emoji that we implemented
	testEmoji := []struct {
		emoji string
		name  string
		x     float64
	}{
		{"😀", "Grinning Face", 50},
		{"👋", "Waving Hand", 120},
		{"👍", "Thumbs Up", 190},
		{"❤", "Red Heart", 260},
		{"🌟", "Star", 330},
	}
	
	y := 80.0
	for _, test := range testEmoji {
		// Parse and render the emoji
		sequences := emojiRenderer.ParseEmojiSequence(test.emoji)
		if len(sequences) > 0 {
			emojiImg := emojiRenderer.RenderEmoji(sequences[0], emojiRenderer.EmojiSize)
			if emojiImg != nil {
				dc.DrawImage(emojiImg, int(test.x), int(y))
			}
		}
		
		// Draw label
		dc.SetRGB(0, 0, 0)
		dc.DrawString(test.name, test.x-10, y+60)
	}
	
	// Add comparison text
	dc.SetRGB(0.5, 0.5, 0.5)
	dc.DrawString("Before: Just colored circles", 20, 170)
	dc.DrawString("After: Recognizable emoji shapes!", 20, 190)
	
	// Save the result
	core.SavePNG("images/text/emoji-test-result.png", dc.Image())
	fmt.Println("Emoji test saved to: emoji-test-result.png")
	fmt.Println("Check the image to see the improved emoji rendering!")
}
