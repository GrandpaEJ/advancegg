package main

import (
	"log"

	"github.com/GrandpaEJ/advancegg"
)

func main() {
	const W = 800
	const H = 600
	dc := advancegg.NewContext(W, H)
	dc.SetRGB(1, 1, 1) // White background
	dc.Clear()
	dc.SetRGB(0, 0, 0) // Black text

	// Helper to draw text
	y := 50.0
	drawLang := func(lang, fontPath, text string) {
		err := dc.LoadFontFace(fontPath, 40)
		if err != nil {
			log.Printf("Failed to load font for %s: %v", lang, err)
			return
		}

		// Draw language name (in English, if needed, but we'll specific font)
		// For simplicity, we just draw the text.

		dc.DrawString(text, 50, y)
		y += 80
	}

	// 1. Bengali
	// Font: NotoSansBengali-Regular.ttf
	drawLang("Bengali", "/usr/share/fonts/truetype/noto/NotoSansBengali-Regular.ttf", "আমার সোনার বাংলা")

	// 2. Hindi (Devanagari)
	// Font: NotoSansDevanagari-Regular.ttf
	drawLang("Hindi", "/usr/share/fonts/truetype/noto/NotoSansDevanagari-Regular.ttf", "नमस्ते दुनिया")

	// 3. Arabic (RTL)
	// Font: NotoSansArabic-Regular.ttf
	// Arabic requires shaping to join letters. HarfBuzz should handle this.
	// "مرحبا بالعالم" (Hello World)
	drawLang("Arabic", "/usr/share/fonts/truetype/noto/NotoSansArabic-Regular.ttf", "مرحبا بالعالم")

	// 4. Tamil
	drawLang("Tamil", "/usr/share/fonts/truetype/noto/NotoSansTamil-Bold.ttf", "வணக்கம்")

	// 5. Telugu
	drawLang("Telugu", "/usr/share/fonts/truetype/noto/NotoSansTelugu-Regular.ttf", "నమస్కారం")

	// 6. Thai
	drawLang("Thai", "/usr/share/fonts/truetype/noto/NotoSansTaiTham-Regular.ttf", "สวัสดี")

	// 7. Khmer
	drawLang("Khmer", "/usr/share/fonts/truetype/noto/NotoSansKhmer-Bold.ttf", "សួស្តី")

	// 8. Myanmar (Burmese)
	drawLang("Myanmar", "/usr/share/fonts/truetype/noto/NotoSansMyanmarUI-Regular.ttf", "မင်္ဂလာပါ")

	// 9. Try CJK (Chinese)
	// Font: /usr/share/fonts/opentype/noto/NotoSansCJK-Regular.ttc
	// Note: loading TTC might fail if our loader or underlying lib doesn't support it.
	// Let's try.
	err := dc.LoadFontFace("/usr/share/fonts/opentype/noto/NotoSansCJK-Regular.ttc", 40)
	if err == nil {
		dc.DrawString("你好世界", 50, y)
		y += 80
	} else {
		log.Printf("Skipping CJK (TTC load failed?): %v", err)
	}

	dc.SavePNG("images/text/multilang_demo.png")
}
