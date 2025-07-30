package main

import (
	"fmt"

	"github.com/GrandpaEJ/advancegg"
)

func main() {
	fmt.Println("Creating Unicode and Emoji examples...")

	// Unicode shaping example
	createUnicodeExample()

	// Emoji rendering example
	createEmojiExample()

	// Mixed text example
	createMixedTextExample()

	fmt.Println("Unicode and Emoji examples completed!")
}

func createUnicodeExample() {
	dc := advancegg.NewContext(800, 600)

	// White background
	dc.SetRGB(1, 1, 1)
	dc.Clear()

	// Title
	dc.SetRGB(0, 0, 0)
	dc.DrawString("Unicode Text Shaping Demo", 50, 50)

	// Set up text shaper
	shaper := advancegg.NewTextShaper()
	dc.SetTextShaper(shaper)

	// Latin text with ligatures
	dc.SetRGB(0.2, 0.2, 0.8)
	dc.DrawString("Latin text: Hello World! (with fi, fl ligatures)", 50, 100)

	// Simulate Arabic text (RTL)
	shaper.Direction = advancegg.TextDirectionRTL
	shaper.Script = advancegg.ScriptArabic
	dc.SetTextShaper(shaper)

	dc.SetRGB(0.8, 0.2, 0.2)
	dc.DrawString("Arabic (RTL): مرحبا بالعالم", 50, 150)

	// Hebrew text (RTL)
	shaper.Script = advancegg.ScriptHebrew
	dc.SetTextShaper(shaper)

	dc.SetRGB(0.2, 0.8, 0.2)
	dc.DrawString("Hebrew (RTL): שלום עולם", 50, 200)

	// Devanagari text (complex shaping)
	shaper.Direction = advancegg.TextDirectionLTR
	shaper.Script = advancegg.ScriptDevanagari
	dc.SetTextShaper(shaper)

	dc.SetRGB(0.8, 0.5, 0.2)
	dc.DrawString("Devanagari: नमस्ते दुनिया", 50, 250)

	// Thai text
	shaper.Script = advancegg.ScriptThai
	dc.SetTextShaper(shaper)

	dc.SetRGB(0.5, 0.2, 0.8)
	dc.DrawString("Thai: สวัสดีชาวโลก", 50, 300)

	// Chinese text
	shaper.Script = advancegg.ScriptChinese
	dc.SetTextShaper(shaper)

	dc.SetRGB(0.2, 0.8, 0.8)
	dc.DrawString("Chinese: 你好世界", 50, 350)

	// Japanese text
	shaper.Script = advancegg.ScriptJapanese
	dc.SetTextShaper(shaper)

	dc.SetRGB(0.8, 0.2, 0.8)
	dc.DrawString("Japanese: こんにちは世界", 50, 400)

	// Korean text
	shaper.Script = advancegg.ScriptKorean
	dc.SetTextShaper(shaper)

	dc.SetRGB(0.6, 0.8, 0.2)
	dc.DrawString("Korean: 안녕하세요 세계", 50, 450)

	// Mixed direction text
	dc.SetRGB(0, 0, 0)
	dc.DrawString("Mixed: Hello مرحبا World שלום!", 50, 500)

	// Text shaping info
	dc.SetRGB(0.5, 0.5, 0.5)
	dc.DrawString("Note: This demo shows Unicode script detection and shaping", 50, 550)
	dc.DrawString("Real implementation would use proper font shaping engines", 50, 570)

	dc.SavePNG("images/unicode-shaping-demo.png")
	fmt.Println("Unicode shaping demo saved as unicode-shaping-demo.png")
}

func createEmojiExample() {
	dc := advancegg.NewContext(800, 600)

	// White background
	dc.SetRGB(1, 1, 1)
	dc.Clear()

	// Title
	dc.SetRGB(0, 0, 0)
	dc.DrawString("Emoji Rendering Demo", 50, 50)

	// Set up emoji renderer
	emojiRenderer := advancegg.NewEmojiRenderer()
	emojiRenderer.EmojiSize = 32
	dc.SetEmojiRenderer(emojiRenderer)

	// Different emoji categories
	categories := []struct {
		name   string
		emojis string
		y      float64
	}{
		{"Smileys & Emotion:", "😀😃😄😁😆😅😂🤣", 100},
		{"People & Body:", "👋👍👎👌✌️🤞🤟🤘", 150},
		{"Animals & Nature:", "🐶🐱🐭🐹🐰🦊🐻🐼", 200},
		{"Food & Drink:", "🍎🍌🍇🍓🥝🍅🥕🌽", 250},
		{"Activities:", "⚽🏀🏈⚾🎾🏐🏉🎱", 300},
		{"Travel & Places:", "🚗🚕🚙🚌🚎🏎️🚓🚑", 350},
		{"Objects:", "📱💻⌨️🖥️🖨️📷📹📼", 400},
		{"Symbols:", "❤️💛💚💙💜🖤🤍🤎", 450},
		{"Flags:", "🏁🚩🏴🏳️🏳️‍🌈🏳️‍⚧️🏴‍☠️", 500},
	}

	for _, category := range categories {
		// Draw category label
		dc.SetRGB(0, 0, 0)
		dc.DrawString(category.name, 50, category.y)

		// Draw emojis with color rendering
		x := 200.0
		for _, emoji := range []rune(category.emojis) {
			if advancegg.IsEmoji(emoji) {
				// Create emoji sequence
				sequence := advancegg.EmojiSequence{
					Runes:    []rune{emoji},
					Text:     string(emoji),
					Category: emojiRenderer.GetEmojiCategory(emoji),
				}

				// Render emoji
				emojiImg := emojiRenderer.RenderEmoji(sequence, emojiRenderer.EmojiSize)
				if emojiImg != nil {
					dc.DrawImage(emojiImg, int(x), int(category.y-emojiRenderer.EmojiSize))
				}

				x += emojiRenderer.EmojiSize + 5
			}
		}
	}

	// Skin tone variations
	dc.SetRGB(0, 0, 0)
	dc.DrawString("Skin tone variations:", 50, 550)

	// Show skin tone modifiers (simplified representation)
	skinTones := []string{"👋", "👋🏻", "👋🏼", "👋🏽", "👋🏾", "👋🏿"}
	x := 200.0
	for _, emoji := range skinTones {
		dc.DrawStringWithEmoji(emoji, x, 550)
		x += 40
	}

	dc.SavePNG("images/emoji-rendering-demo.png")
	fmt.Println("Emoji rendering demo saved as emoji-rendering-demo.png")
}

func createMixedTextExample() {
	dc := advancegg.NewContext(800, 600)

	// Gradient background
	for y := 0; y < 600; y++ {
		t := float64(y) / 600.0
		dc.SetRGB(0.9+t*0.1, 0.95+t*0.05, 1.0)
		dc.DrawLine(0, float64(y), 800, float64(y))
		dc.Stroke()
	}

	// Title
	dc.SetRGB(0, 0, 0)
	dc.DrawString("Mixed Text with Unicode & Emoji", 50, 50)

	// Set up renderers
	emojiRenderer := advancegg.NewEmojiRenderer()
	emojiRenderer.EmojiSize = 24
	dc.SetEmojiRenderer(emojiRenderer)

	// Mixed content examples
	examples := []struct {
		text string
		y    float64
		desc string
	}{
		{"Hello 👋 World 🌍!", 120, "English with emojis"},
		{"Café ☕ with emoji", 160, "Latin with diacritics and emoji"},
		{"مرحبا 👋 بالعالم 🌍", 200, "Arabic with emojis"},
		{"שלום 👋 עולם 🌍", 240, "Hebrew with emojis"},
		{"こんにちは 👋 世界 🌍", 280, "Japanese with emojis"},
		{"안녕하세요 👋 세계 🌍", 320, "Korean with emojis"},
		{"Привет 👋 мир 🌍", 360, "Cyrillic with emojis"},
		{"Γεια σας 👋 κόσμος 🌍", 400, "Greek with emojis"},
		{"🎉 Party time! 🎊🥳🎈", 440, "Emoji-heavy text"},
		{"Code: func() 💻 { return 🚀 }", 480, "Technical text with emojis"},
	}

	for _, example := range examples {
		// Draw description
		dc.SetRGB(0.5, 0.5, 0.5)
		dc.DrawString(example.desc, 50, example.y-15)

		// Draw mixed text
		dc.SetRGB(0, 0, 0)
		dc.DrawStringWithEmoji(example.text, 50, example.y)
	}

	// Emoji info section
	dc.SetRGB(0.2, 0.2, 0.8)
	dc.DrawString("Emoji Information:", 50, 530)

	// Get info about a specific emoji
	emojiInfo := emojiRenderer.GetEmojiInfo("🎉")
	if emojiInfo != nil {
		dc.SetRGB(0, 0, 0)
		infoText := fmt.Sprintf("🎉 - %s (%s)", emojiInfo.Name, emojiInfo.Category)
		dc.DrawString(infoText, 50, 550)
	}

	dc.SavePNG("images/mixed-text-demo.png")
	fmt.Println("Mixed text demo saved as mixed-text-demo.png")
}

func createAdvancedTextExample() {
	dc := advancegg.NewContext(1000, 700)

	// Dark background
	dc.SetRGB(0.1, 0.1, 0.2)
	dc.Clear()

	// Title
	dc.SetRGB(1, 1, 1)
	dc.DrawString("Advanced Text Rendering Features", 50, 50)

	// Set up advanced text shaper
	shaper := advancegg.NewTextShaper()
	shaper.EnableLigatures = true
	shaper.EnableKerning = true
	dc.SetTextShaper(shaper)

	// Ligature examples
	dc.SetRGB(0.8, 0.8, 0.2)
	dc.DrawString("Ligatures: fi fl ff ffi ffl", 50, 120)

	// Kerning examples
	dc.SetRGB(0.2, 0.8, 0.8)
	dc.DrawString("Kerning: AV AW AY AT VA WA YA TA", 50, 160)

	// Complex script examples
	dc.SetRGB(0.8, 0.2, 0.8)
	dc.DrawString("Complex Scripts:", 50, 220)

	// Arabic contextual forms
	shaper.Script = advancegg.ScriptArabic
	shaper.Direction = advancegg.TextDirectionRTL
	dc.SetTextShaper(shaper)

	dc.SetRGB(0.8, 0.4, 0.2)
	dc.DrawString("Arabic contextual: بسم الله الرحمن الرحيم", 50, 260)

	// Devanagari conjuncts
	shaper.Script = advancegg.ScriptDevanagari
	shaper.Direction = advancegg.TextDirectionLTR
	dc.SetTextShaper(shaper)

	dc.SetRGB(0.2, 0.8, 0.4)
	dc.DrawString("Devanagari conjuncts: क्ष त्र ज्ञ श्र", 50, 300)

	// Emoji with skin tones
	emojiRenderer := advancegg.NewEmojiRenderer()
	emojiRenderer.EmojiSize = 32
	dc.SetEmojiRenderer(emojiRenderer)

	dc.SetRGB(1, 1, 1)
	dc.DrawString("Emoji skin tones:", 50, 360)

	skinToneEmojis := []string{"👶", "👶🏻", "👶🏼", "👶🏽", "👶🏾", "👶🏿"}
	x := 50.0
	for _, emoji := range skinToneEmojis {
		dc.DrawStringWithEmoji(emoji, x, 400)
		x += 60
	}

	// ZWJ sequences (simplified)
	dc.SetRGB(1, 1, 1)
	dc.DrawString("ZWJ sequences:", 50, 460)
	dc.DrawStringWithEmoji("👨‍👩‍👧‍👦 👨‍💻 👩‍🚀 🏳️‍🌈", 50, 500)

	// Text direction indicators
	dc.SetRGB(0.6, 0.6, 0.6)
	dc.DrawString("← RTL direction", 700, 260)
	dc.DrawString("LTR direction →", 700, 300)

	// Performance info
	dc.SetRGB(0.5, 0.5, 0.5)
	dc.DrawString("Features demonstrated:", 50, 580)
	dc.DrawString("• Unicode script detection and shaping", 70, 600)
	dc.DrawString("• Bidirectional text (BiDi) algorithm", 70, 620)
	dc.DrawString("• Color emoji rendering with fallbacks", 70, 640)
	dc.DrawString("• Ligatures and kerning", 70, 660)

	dc.SavePNG("images/advanced-text-demo.png")
	fmt.Println("Advanced text demo saved as advanced-text-demo.png")
}
