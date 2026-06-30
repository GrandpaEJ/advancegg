package main

import (
	"log"

	"github.com/GrandpaEJ/advancegg"
)

func main() {
	const W = 1400
	const H = 1200
	dc := advancegg.NewContext(W, H)

	dc.SetRGB(1, 1, 1)
	dc.Clear()

	fonts := []struct {
		name string
		path string
	}{
		{"Akaash (Comic)", "assets/fonts/comic/AkaashNormal.ttf"},
		{"Likhan (Handwriting)", "assets/fonts/comic/LikhanNormal.ttf"},
		{"Jamrul (Display)", "assets/fonts/comic/JamrulNormal.ttf"},
	}

	// Title
	dc.SetRGB(0, 0, 0)
	dc.DrawString("Bengali Comic Font GPOS Test — 3 Fonts", 30, 30)
	dc.DrawString("Matras should attach to base chars. Conjuncts should form correctly.", 30, 55)

	y := 90.0
	colW := 440.0

	for fi, font := range fonts {
		x := 30 + float64(fi)*colW

		// Font label
		dc.SetRGB(0, 0, 0)
		dc.DrawString(font.name, x, y)

		err := dc.LoadScriptFont(advancegg.ScriptBengali, font.path, 28)
		if err != nil {
			log.Printf("✗ Failed to load %s: %v", font.name, err)
			dc.SetRGB(1, 0, 0)
			dc.DrawString("LOAD FAILED", x, y+20)
			continue
		}
		log.Printf("✓ Loaded %s", font.name)

		dc.SetRGB(0, 0, 0)

		// Load same Latin font
		dc.LoadFontFace("assets/fonts/NotoSans-Regular.ttf", 28)

		// ── All matras with ক ──
		ry := y + 30
		dc.DrawString("All matras with ক:", x, ry)
		ry += 22
		dc.DrawString("ক কা কি কী কু কূ কৃ কে কৈ কো কৌ", x, ry)

		// ── More consonant tables ──
		ry += 24
		dc.DrawString("ত তা তি তী তু তূ তৃ তে তৈ তো তৌ", x, ry)
		ry += 20
		dc.DrawString("ন না নি নী নু নূ নৃ নে নৈ নো নৌ", x, ry)
		ry += 20
		dc.DrawString("ব বা বি বী বু বূ বৃ বে বৈ বো বৌ", x, ry)
		ry += 20
		dc.DrawString("ম মা মি মী মু মূ মৃ মে মৈ মো মৌ", x, ry)
		ry += 20
		dc.DrawString("স সা সি সী সু সূ সৃ সে সৈ সো সৌ", x, ry)

		// ── Conjuncts ──
		ry += 28
		dc.DrawString("Conjuncts:", x, ry)
		ry += 22
		dc.DrawString("ক্ষ ত্র জ্ঞ স্ত স্থ স্ক স্প স্ট্র", x, ry)
		ry += 20
		dc.DrawString("দ্ধ দ্ব ন্দ ন্ধ ন্ত ম্প ম্ব ম্ভ", x, ry)
		ry += 20
		dc.DrawString("ষ্ট ষ্ঠ ষ্প ষ্ণ র্থ র্ম হ্ম হ্ন", x, ry)
		ry += 20
		dc.DrawString("র্ক র্গ র্পর্শ্ব ক্র গ্র প্র ব্র", x, ry)
		ry += 20
		dc.DrawString("ভ্র শ্র প্ল ক্ল গ্ল ব্ল স্ন", x, ry)

		// ── Words & sentences ──
		ry += 28
		dc.DrawString("Words & sentences:", x, ry)
		ry += 22
		dc.DrawString("বাংলাদেশ ভাষা বিশ্ব শিক্ষা", x, ry)
		ry += 20
		dc.DrawString("সূর্য চন্দ্র গ্রহ তারা আকাশ", x, ry)
		ry += 20
		dc.DrawString("গ্রন্থ হৃদয় বর্ষা সপ্তাহ", x, ry)
		ry += 24
		dc.DrawString("বাংলা ভাষা একটি সুন্দর ভাষা।", x, ry)
		ry += 20
		dc.DrawString("আমার বিদ্যালয়ে প্রতিদিন পাঠ হয়।", x, ry)
	}

	outputPath := "images/text/bengali_comic_test.png"
	err := dc.SavePNG(outputPath)
	if err != nil {
		log.Fatalf("Failed to save image: %v", err)
	}

	log.Printf("✓ Created Bengali comic font test: %s", outputPath)
}
