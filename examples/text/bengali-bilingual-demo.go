package main

import (
	"log"

	"github.com/GrandpaEJ/advancegg"
)

func main() {
	const W = 900
	const H = 700
	dc := advancegg.NewContext(W, H)

	dc.SetRGB(1, 1, 1)
	dc.Clear()

	err := dc.LoadFontFace("assets/fonts/NotoSans-Regular.ttf", 28)
	if err != nil {
		log.Fatalf("Failed to load main font: %v", err)
	}

	err = dc.LoadScriptFont(advancegg.ScriptBengali, "assets/fonts/NotoSansBengali-Regular.ttf", 28)
	if err != nil {
		log.Fatalf("Failed to load Bengali font: %v", err)
	}

	err = dc.LoadScriptFont(advancegg.ScriptDevanagari, "assets/fonts/NotoSansDevanagari-Regular.ttf", 28)
	if err != nil {
		log.Fatalf("Failed to load Devanagari font: %v", err)
	}
	log.Println("✓ Bilingual demo — fonts loaded")

	// Title — mixed Bangla + English
	dc.SetRGB(0, 0, 0)
	dc.DrawStringAnchored("বাংলা-English Bilingual Demo — Mixing Two Languages", W/2, 40, 0.5, 0.5)

	// Subtitle
	dc.SetRGB(0.3, 0.3, 0.3)
	dc.DrawStringAnchored("প্রতিটি line-এ Bangla আর English naturally mix করা হয়েছে", W/2, 75, 0.5, 0.5)

	// Section 1: Daily conversations with mixed language
	drawCard(dc, 30, 110, 410, 250, "Daily Conversation (দৈনিক কথাবার্তা)", []string{
		"আমি today market-এ যাব।",
		"She told me যে সে আসবে না।",
		"এই book-টা খুব interesting।",
		"Please wait, আমি call করছি।",
		"আজকের weather খুব nice।",
	})

	// Section 2: Technical / Office mixed language
	drawCard(dc, 460, 110, 410, 250, "Office / Tech Talk", []string{
		"আমি team-কে feature propose করেছি।",
		"This API টা fix করতে হবে।",
		"Server-টা আবার down হয়েছে।",
		"Code review-এ suggestion পেয়েছি।",
		"আমার report-টা almost ready।",
	})

	// Section 3: Sentences with technical terms about this library
	drawCard(dc, 30, 390, 410, 250, "Library Features (Bilingual)", []string{
		"AdvanceGG এখন Bengali GPOS support করে।",
		"Matras (া, ি, ী) base consonant-এ attach হয়।",
		"75+ conjuncts (যুক্তবর্ণ) GSUB-এ formed হয়।",
		"Per-script font loading এ বাংলা আলাদা font-এ render হয়।",
		"Emoji-ও automatic detect হয়! 🎉",
	})

	// Section 4: Mixed proverbs & fun
	drawCard(dc, 460, 390, 410, 250, "Proverbs & Fun", []string{
		"Time and tide কারও জন্য wait করে না!",
		"All that glitters is not gold!",
		"যত বড়ই problem হোক, solution আছে।",
		"Learning a new language is tough!",
		"বাংলা + English = ❤️",
	})

	// Footer
	dc.SetRGB(0.5, 0.5, 0.5)
	dc.DrawStringAnchored("AdvanceGG — Bilingual (Bangla + English) Text Rendering Demo", W/2, H-30, 0.5, 0.5)

	outputPath := "images/text/bengali_bilingual_demo.png"
	err = dc.SavePNG(outputPath)
	if err != nil {
		log.Fatalf("Failed to save image: %v", err)
	}

	log.Printf("✓ Created bilingual demo image: %s", outputPath)
}

func drawCard(dc *advancegg.Context, x, y, w, h float64, title string, lines []string) {
	// Card background
	dc.SetRGB(0.97, 0.97, 0.97)
	dc.DrawRoundedRectangle(x, y, w, h, 10)
	dc.Fill()

	// Card border
	dc.SetRGB(0.7, 0.7, 0.7)
	dc.SetLineWidth(1)
	dc.DrawRoundedRectangle(x, y, w, h, 10)
	dc.Stroke()

	// Title
	dc.SetRGB(0, 0, 0)
	dc.DrawString(title, x+15, y+25)

	// Lines
	dc.SetRGB(0.15, 0.15, 0.15)
	lineH := 35.0
	for i, line := range lines {
		dc.DrawString("• "+line, x+15, y+55+float64(i)*lineH)
	}
}
