package main

import (
	"fmt"

	"github.com/GrandpaEJ/advancegg"
)

func main() {
	dc := advancegg.NewContext(800, 600)

	// White background
	dc.SetRGB(1, 1, 1)
	dc.Clear()

	// Black text color
	dc.SetRGB(0, 0, 0)

	// Demonstrate font format detection
	y := 50.0

	// Title
	dc.DrawString("Font Format Support Demo", 50, y)
	y += 60

	// Try to load different font formats
	fontPaths := []string{
		"/System/Fonts/Arial.ttf",                         // Common TTF on macOS
		"/System/Fonts/Helvetica.ttc",                     // TTC (TrueType Collection)
		"/System/Library/Fonts/Avenir.ttc",                // Another common font
		"/usr/share/fonts/truetype/dejavu/DejaVuSans.ttf", // Linux TTF
		"/Windows/Fonts/arial.ttf",                        // Windows TTF
		"/Windows/Fonts/calibri.ttf",                      // Windows TTF
	}

	// Test each font path
	for _, fontPath := range fontPaths {
		// Try to detect font format
		format, err := advancegg.GetFontFormat(fontPath)
		if err != nil {
			continue // Skip if file doesn't exist
		}

		// Try to load the font
		err = dc.LoadFontFace(fontPath, 24)
		if err != nil {
			continue // Skip if can't load
		}

		// Display font info
		text := fmt.Sprintf("✓ Loaded %s font: %s", format, fontPath)
		dc.DrawString(text, 50, y)
		y += 30

		// Show sample text with this font
		dc.DrawString("Sample text: The quick brown fox jumps over the lazy dog", 70, y)
		y += 40

		break // Use the first successfully loaded font for demo
	}

	// If no system fonts found, use built-in font
	if y <= 110 {
		dc.DrawString("No system fonts found. Using built-in font.", 50, y)
		y += 40
	}

	// Demonstrate different font loading methods
	dc.DrawString("Font Loading Methods:", 50, y)
	y += 40

	methods := []string{
		"LoadFontFace(path, size) - Auto-detect TTF/OTF",
		"LoadTTFFace(path, size) - Explicit TTF loading",
		"LoadOTFFace(path, size) - Explicit OTF loading",
		"LoadFontFaceFromBytes(bytes, size) - From memory",
		"LoadFontFaceWithOptions(path, options) - Custom options",
	}

	for _, method := range methods {
		dc.DrawString("• "+method, 70, y)
		y += 25
	}

	// Add format information
	y += 20
	dc.DrawString("Supported Font Formats:", 50, y)
	y += 30

	formats := []string{
		"TTF (TrueType Font) - Standard vector font format",
		"OTF (OpenType Font) - Advanced font format with more features",
		"TTC (TrueType Collection) - Multiple fonts in one file",
	}

	for _, format := range formats {
		dc.DrawString("• "+format, 70, y)
		y += 25
	}

	// Save the image
	dc.SavePNG("images/text/font-formats.png")
	fmt.Println("Font formats demo saved as font-formats.png")
}
