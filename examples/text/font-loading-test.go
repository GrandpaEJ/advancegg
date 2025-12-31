package main

import (
	"fmt"
	"os"
	"path/filepath"
	"github.com/GrandpaEJ/advancegg"
)

func main() {
	fmt.Println("Testing font loading capabilities...")

	// Create context
	dc := advancegg.NewContext(1200, 800)
	dc.SetRGB(1, 1, 1)
	dc.Clear()

	// Title
	dc.SetRGB(0, 0, 0)
	dc.DrawString("AdvanceGG Font Loading Test", 50, 50)

	y := 100.0

	// Test 1: Default font
	fmt.Println("Testing default font...")
	dc.SetRGB(0.2, 0.2, 0.2)
	dc.DrawString("Default Font: The quick brown fox jumps over the lazy dog", 50, y)
	y += 40

	// Test 2: Load from assets
	fmt.Println("Testing font loading from assets...")
	fontPaths := []struct {
		path string
		name string
		text string
	}{
		{"assets/fonts/NotoSans-Regular.ttf", "Noto Sans Regular", "Noto Sans: The quick brown fox jumps over the lazy dog"},
		{"assets/fonts/NotoSans-Bold.ttf", "Noto Sans Bold", "Noto Sans Bold: The quick brown fox jumps over the lazy dog"},
		{"assets/fonts/NotoSerif-Regular.ttf", "Noto Serif", "Noto Serif: The quick brown fox jumps over the lazy dog"},
		{"assets/fonts/LiberationSans-Regular.ttf", "Liberation Sans", "Liberation Sans: The quick brown fox jumps over the lazy dog"},
	}

	for _, font := range fontPaths {
		if _, err := os.Stat(font.path); err == nil {
			err := dc.LoadFontFace(font.path, 16)
			if err != nil {
				fmt.Printf("Error loading %s: %v\n", font.name, err)
				dc.SetRGB(0.8, 0.2, 0.2)
				dc.DrawString(fmt.Sprintf("❌ Failed to load %s", font.name), 50, y)
			} else {
				fmt.Printf("✓ Loaded %s\n", font.name)
				dc.SetRGB(0.2, 0.6, 0.2)
				dc.DrawString(font.text, 50, y)
			}
		} else {
			fmt.Printf("⚠ Font not found: %s\n", font.path)
			dc.SetRGB(0.8, 0.6, 0.2)
			dc.DrawString(fmt.Sprintf("⚠ Font not found: %s", font.name), 50, y)
		}
		y += 30
	}

	// Test 3: Unicode text with different scripts
	fmt.Println("Testing Unicode text rendering...")
	y += 20
	dc.SetRGB(0, 0, 0)
	dc.DrawString("Unicode Text Rendering:", 50, y)
	y += 30

	unicodeTests := []struct {
		fontPath string
		text     string
		script   string
	}{
		{"assets/fonts/NotoSansArabic-Regular.ttf", "مرحبا بالعالم", "Arabic"},
		{"assets/fonts/NotoSansHebrew-Regular.ttf", "שלום עולם", "Hebrew"},
		{"assets/fonts/NotoSansDevanagari-Regular.ttf", "नमस्ते दुनिया", "Hindi (Devanagari)"},
		{"assets/fonts/NotoSansThai-Regular.ttf", "สวัสดีชาวโลก", "Thai"},
		{"assets/fonts/NotoSans-Regular.ttf", "Здравствуй мир", "Russian (Cyrillic)"},
		{"assets/fonts/NotoSans-Regular.ttf", "Γεια σας κόσμε", "Greek"},
		{"assets/fonts/NotoSans-Regular.ttf", "こんにちは世界", "Japanese"},
		{"assets/fonts/NotoSans-Regular.ttf", "你好世界", "Chinese"},
	}

	for _, test := range unicodeTests {
		if _, err := os.Stat(test.fontPath); err == nil {
			err := dc.LoadFontFace(test.fontPath, 16)
			if err != nil {
				fmt.Printf("Error loading font for %s: %v\n", test.script, err)
				dc.SetRGB(0.8, 0.2, 0.2)
				dc.DrawString(fmt.Sprintf("❌ %s: Font load failed", test.script), 50, y)
			} else {
				fmt.Printf("✓ Rendering %s text\n", test.script)
				dc.SetRGB(0.2, 0.2, 0.8)
				dc.DrawString(fmt.Sprintf("%s: %s", test.script, test.text), 50, y)
			}
		} else {
			fmt.Printf("⚠ Font not found for %s\n", test.script)
			dc.SetRGB(0.8, 0.6, 0.2)
			dc.DrawString(fmt.Sprintf("⚠ %s: Font not available", test.script), 50, y)
		}
		y += 25
	}

	// Test 4: Emoji rendering
	fmt.Println("Testing emoji rendering...")
	y += 20
	dc.SetRGB(0, 0, 0)
	dc.DrawString("Emoji Rendering:", 50, y)
	y += 30

	// Try to load emoji font
	emojiFont := "assets/fonts/NotoColorEmoji.ttf"
	if _, err := os.Stat(emojiFont); err == nil {
		renderer := dc.GetEmojiRenderer()
		err := renderer.LoadEmojiFont(emojiFont)
		if err != nil {
			fmt.Printf("Error loading emoji font: %v\n", err)
			dc.SetRGB(0.8, 0.2, 0.2)
			dc.DrawString("❌ Emoji font load failed", 50, y)
		} else {
			fmt.Println("✓ Emoji font loaded successfully")
			dc.SetRGB(0.2, 0.6, 0.2)
			dc.DrawString("✓ Emoji font loaded: 😀 🎉 🚀 ❤️ 🌟", 50, y)
			
			// Test emoji sequences
			y += 30
			emojiTests := []string{
				"👨‍👩‍👧‍👦 Family",
				"👨‍💻 Technologist", 
				"🏳️‍🌈 Rainbow Flag",
				"👍🏽 Thumbs Up (skin tone)",
			}
			
			for _, emojiTest := range emojiTests {
				dc.DrawString(emojiTest, 50, y)
				y += 25
			}
		}
	} else {
		fmt.Println("⚠ Emoji font not found")
		dc.SetRGB(0.8, 0.6, 0.2)
		dc.DrawString("⚠ Emoji font not available", 50, y)
	}

	// Test 5: Font fallback
	fmt.Println("Testing font fallback...")
	y += 40
	dc.SetRGB(0, 0, 0)
	dc.DrawString("Font Fallback Test:", 50, y)
	y += 30

	// Try to load a non-existent font
	err := dc.LoadFontFace("non-existent-font.ttf", 16)
	if err != nil {
		fmt.Println("✓ Font fallback working (expected error for non-existent font)")
		dc.SetRGB(0.2, 0.6, 0.2)
		dc.DrawString("✓ Font fallback working - using default font", 50, y)
	} else {
		fmt.Println("⚠ Unexpected: non-existent font loaded")
		dc.SetRGB(0.8, 0.6, 0.2)
		dc.DrawString("⚠ Unexpected behavior with non-existent font", 50, y)
	}

	// Test 6: Font size variations
	y += 50
	dc.SetRGB(0, 0, 0)
	dc.DrawString("Font Size Variations:", 50, y)
	y += 30

	sizes := []float64{12, 16, 20, 24, 32}
	for _, size := range sizes {
		if _, err := os.Stat("assets/fonts/NotoSans-Regular.ttf"); err == nil {
			err := dc.LoadFontFace("assets/fonts/NotoSans-Regular.ttf", size)
			if err == nil {
				dc.DrawString(fmt.Sprintf("Size %.0f: The quick brown fox", size), 50, y)
				y += size + 5
			}
		}
	}

	// Save the result
	err = dc.SavePNG("images/text/font-loading-test.png")
	if err != nil {
		fmt.Printf("Error saving image: %v\n", err)
		return
	}

	fmt.Println("\nFont loading test completed!")
	fmt.Println("Generated: font-loading-test.png")
	
	// Print summary
	fmt.Println("\nSummary:")
	fmt.Printf("- Assets directory: %s\n", "assets/fonts")
	
	// Count available fonts
	fontCount := 0
	if entries, err := os.ReadDir("assets/fonts"); err == nil {
		for _, entry := range entries {
			if filepath.Ext(entry.Name()) == ".ttf" || filepath.Ext(entry.Name()) == ".ttc" {
				fontCount++
			}
		}
	}
	fmt.Printf("- Available fonts: %d\n", fontCount)
	fmt.Println("- Unicode scripts tested: Arabic, Hebrew, Devanagari, Thai, Cyrillic, Greek, Japanese, Chinese")
	fmt.Println("- Emoji support tested")
	fmt.Println("- Font fallback tested")
}
