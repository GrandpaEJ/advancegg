package main

import (
	"log"

	"github.com/GrandpaEJ/advancegg"
)

func main() {
	const W = 800
	const H = 600
	dc := advancegg.NewContext(W, H)

	// White background
	dc.SetRGB(1, 1, 1)
	dc.Clear()

	// Load main font (for English/Latin text)
	err := dc.LoadFontFace("assets/fonts/NotoSans-Regular.ttf", 32)
	if err != nil {
		log.Fatalf("Failed to load main font: %v", err)
	}

	// Load Bengali script font using the new script font API
	err = dc.LoadScriptFont(advancegg.ScriptBengali, "assets/fonts/NotoSansBengali-Regular.ttf", 32)
	if err != nil {
		log.Fatalf("Failed to load Bengali font: %v", err)
	}
	log.Println("✓ Bengali script font loaded successfully")

	// Draw title
	dc.SetRGB(0, 0, 0)
	dc.DrawString("Bengali Text Rendering Test", 50, 50)

	// Draw speech bubble 1 (top)
	drawSpeechBubble(dc, 100, 100, 300, 120, "উফ... আমি কোনোমতে বেঁচে ফিরিলাম।")

	// Draw speech bubble 2 (middle)
	drawSpeechBubble(dc, 400, 250, 350, 100, "আমার কাজে প্রথম দিনে দেরি হয়ে যেতে পারত।")

	// Draw speech bubble 3 (bottom)
	drawSpeechBubble(dc, 150, 400, 280, 110, "বাংলা ভাষা সুন্দর।")

	// Draw some English text to show mixed rendering
	dc.SetRGB(0, 0, 0.5)
	dc.DrawString("English text works too!", 50, 550)

	// Save the image
	outputPath := "images/text/bengali_demo.png"
	err = dc.SavePNG(outputPath)
	if err != nil {
		log.Fatalf("Failed to save image: %v", err)
	}

	log.Printf("✓ Created Bengali demo image: %s", outputPath)
}

// drawSpeechBubble draws a rounded rectangle with text inside
func drawSpeechBubble(dc *advancegg.Context, x, y, w, h float64, text string) {
	// Draw bubble background
	dc.SetRGB(0.95, 0.95, 0.95)
	dc.DrawRoundedRectangle(x, y, w, h, 15)
	dc.Fill()

	// Draw bubble border
	dc.SetRGB(0.3, 0.3, 0.3)
	dc.SetLineWidth(2)
	dc.DrawRoundedRectangle(x, y, w, h, 15)
	dc.Stroke()

	// Draw text centered in bubble
	dc.SetRGB(0, 0, 0)
	textX := x + w/2
	textY := y + h/2
	dc.DrawStringAnchored(text, textX, textY, 0.5, 0.5)
}
