package main

import (
	"github.com/GrandpaEJ/advancegg"
)

func main() {
	// Create a canvas
	dc := advancegg.NewContext(800, 600)

	// Set background
	dc.SetRGB(1, 1, 1)
	dc.Clear()

	// Load a font
	if err := dc.LoadFontFace("assets/fonts/Roboto-Regular.ttf", 48); err != nil {
		// Try alternative font paths
		_ = dc.LoadFontFace("/usr/share/fonts/truetype/dejavu/DejaVuSans.ttf", 48)
	}

	// Draw text with emojis - emoji rendering is enabled by default
	dc.SetRGB(0, 0, 0)
	dc.DrawStringAnchored("Hello World 😀", 400, 100, 0.5, 0.5)
	dc.DrawStringAnchored("I ❤️ Go Programming! 🚀", 400, 200, 0.5, 0.5)
	dc.DrawStringAnchored("Emoji: 👋 👍 🌟 🎉", 400, 300, 0.5, 0.5)

	// Draw text with skin tone modifiers
	dc.DrawStringAnchored("Skin tones: 👋🏻 👋🏽 👋🏿", 400, 400, 0.5, 0.5)

	// Draw text without emoji (disable auto-emoji)
	dc.SetEnableAutoEmoji(false)
	dc.DrawStringAnchored("No emoji: 😀 ❤️ 🚀", 400, 500, 0.5, 0.5)

	// Save the result
	dc.SavePNG("emoji-rendering.png")
	println("Saved emoji-rendering.png")
}
