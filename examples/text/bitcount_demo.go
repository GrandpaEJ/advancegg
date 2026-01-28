package main

import (
	"log"

	"github.com/GrandpaEJ/advancegg"
)

func main() {
	// 1. Create a new context
	width, height := 800, 200
	dc := advancegg.NewContext(width, height)

	// 2. Set white background
	dc.SetRGB(1, 1, 1)
	dc.Clear()

	// 3. Load the downloaded Bitcount Single font
	fontPath := "assets/fonts/BitcountSingle-Regular.ttf"
	if err := dc.LoadFontFace(fontPath, 64); err != nil {
		log.Fatalf("Failed to load font: %v", err)
	}

	// 4. Set text color (black)
	dc.SetRGB(0, 0, 0)

	// 5. Draw text
	text := "Bitcount Single Test"
	// Center the text simply for this test
	textWidth, textHeight := dc.MeasureString(text)
	x := float64(width)/2 - textWidth/2
	y := float64(height)/2 + textHeight/4 // Roughly centered vertically

	dc.DrawString(text, x, y)

	// 6. Save the result
	outputPath := "bitcount_test.png"
	if err := dc.SavePNG(outputPath); err != nil {
		log.Fatalf("Failed to save image: %v", err)
	}

	log.Printf("Successfully saved %s", outputPath)
}
