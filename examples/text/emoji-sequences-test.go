package main

import (
	"fmt"
	"github.com/GrandpaEJ/advancegg"
)

func main() {
	fmt.Println("Testing enhanced emoji sequence handling...")

	// Create context
	dc := advancegg.NewContext(800, 600)
	dc.SetRGB(1, 1, 1)
	dc.Clear()

	// Test various emoji sequences
	testEmojis := []struct {
		emoji       string
		description string
		x, y        float64
	}{
		{"👨‍👩‍👧‍👦", "Family: Man, Woman, Girl, Boy", 100, 100},
		{"👨‍💻", "Man Technologist", 300, 100},
		{"👩‍⚕️", "Woman Health Worker", 500, 100},
		{"👨‍❤️‍👨", "Couple: Man, Man", 100, 250},
		{"👩‍❤️‍👩", "Couple: Woman, Woman", 300, 250},
		{"🧑‍🤝‍🧑", "People Holding Hands", 500, 250},
		{"👨🏻", "Man: Light Skin Tone", 100, 400},
		{"👩🏾", "Woman: Medium-Dark Skin Tone", 300, 400},
		{"🏳️‍🌈", "Rainbow Flag", 500, 400},
		{"👍🏽", "Thumbs Up: Medium Skin Tone", 700, 400},
	}

	// Set up emoji renderer
	renderer := dc.GetEmojiRenderer()
	renderer.EmojiSize = 64

	// Draw title
	dc.SetRGB(0, 0, 0)
	dc.DrawString("Enhanced Emoji Sequence Handling Test", 50, 50)

	// Test each emoji sequence
	for _, test := range testEmojis {
		fmt.Printf("Testing: %s (%s)\n", test.emoji, test.description)
		
		// Parse the emoji sequence
		sequences := renderer.ParseEmojiSequence(test.emoji)
		
		if len(sequences) > 0 {
			sequence := sequences[0]
			fmt.Printf("  - Parsed %d runes\n", len(sequence.Runes))
			fmt.Printf("  - Is ZWJ sequence: %v\n", sequence.IsZWJ)
			fmt.Printf("  - Has modifier: %v\n", sequence.HasModifier)
			if sequence.HasModifier {
				fmt.Printf("  - Skin tone: %s\n", sequence.SkinTone)
			}
			fmt.Printf("  - Category: %s\n", sequence.Category)
			
			// Render the emoji
			emojiImg := renderer.RenderEmoji(sequence, renderer.EmojiSize)
			if emojiImg != nil {
				dc.DrawImage(emojiImg, int(test.x), int(test.y))
			}
		}
		
		// Draw description
		dc.SetRGB(0.3, 0.3, 0.3)
		dc.DrawString(test.description, test.x, test.y+80)
		dc.SetRGB(0, 0, 0)
	}

	// Save the result
	err := dc.SavePNG("images/text/emoji-sequences-test.png")
	if err != nil {
		fmt.Printf("Error saving image: %v\n", err)
		return
	}

	fmt.Println("Enhanced emoji sequence test completed!")
	fmt.Println("Generated: emoji-sequences-test.png")
}
