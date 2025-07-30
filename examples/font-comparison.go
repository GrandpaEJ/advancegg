package main

import (
	"fmt"
	"os"
	
	"github.com/GrandpaEJ/advancegg"
)

func main() {
	dc := advancegg.NewContext(1000, 700)
	
	// White background
	dc.SetRGB(1, 1, 1)
	dc.Clear()
	
	// Title
	dc.SetRGB(0, 0, 0)
	dc.DrawString("TTF vs OTF Font Comparison", 50, 50)
	
	// Find available fonts
	fonts := findFonts()
	if len(fonts) == 0 {
		dc.DrawString("No fonts found for comparison", 50, 100)
		dc.SavePNG("font-comparison.png")
		return
	}
	
	y := 100.0
	
	// Display font information
	for i, fontInfo := range fonts {
		if i >= 5 { // Limit to 5 fonts to fit on screen
			break
		}
		
		// Load the font
		err := dc.LoadFontFace(fontInfo.path, 24)
		if err != nil {
			continue
		}
		
		// Font info header
		dc.SetRGB(0.2, 0.2, 0.2)
		dc.LoadFontFace(fontInfo.path, 16)
		header := fmt.Sprintf("%s - %s (%s)", fontInfo.name, fontInfo.format, fontInfo.path)
		dc.DrawString(header, 50, y)
		y += 25
		
		// Sample text with the font
		dc.SetRGB(0, 0, 0)
		dc.LoadFontFace(fontInfo.path, 20)
		sampleText := "The quick brown fox jumps over the lazy dog 1234567890"
		dc.DrawString(sampleText, 70, y)
		y += 30
		
		// Show different sizes
		sizes := []float64{12, 16, 24}
		for _, size := range sizes {
			dc.LoadFontFace(fontInfo.path, size)
			sizeText := fmt.Sprintf("Size %.0f: ABCDEFG abcdefg", size)
			dc.DrawString(sizeText, 90, y)
			y += size + 5
		}
		
		y += 20
	}
	
	// Add technical information
	y += 20
	dc.LoadFontFace(fonts[0].path, 16)
	dc.SetRGB(0.3, 0.3, 0.3)
	
	info := []string{
		"Font Format Differences:",
		"• TTF (TrueType): Quadratic Bézier curves, simpler structure",
		"• OTF (OpenType): Cubic Bézier curves, advanced typography features",
		"• Both formats are supported by AdvanceGG using the same API",
		"",
		"Advanced Features (when available):",
		"• Ligatures, kerning, and advanced typography",
		"• Multiple language support",
		"• Variable font weights and styles",
	}
	
	for _, line := range info {
		dc.DrawString(line, 50, y)
		y += 20
	}
	
	// Save the image
	dc.SavePNG("font-comparison.png")
	fmt.Println("Font comparison saved as font-comparison.png")
}

type FontInfo struct {
	name   string
	path   string
	format string
}

func findFonts() []FontInfo {
	var fonts []FontInfo
	
	// Common font paths with their likely names
	candidates := []FontInfo{
		{"Arial", "/System/Fonts/Arial.ttf", ""},
		{"Helvetica", "/System/Fonts/Helvetica.ttc", ""},
		{"Times", "/System/Library/Fonts/Times.ttc", ""},
		{"Avenir", "/System/Library/Fonts/Avenir.ttc", ""},
		{"DejaVu Sans", "/usr/share/fonts/truetype/dejavu/DejaVuSans.ttf", ""},
		{"Liberation Sans", "/usr/share/fonts/truetype/liberation/LiberationSans-Regular.ttf", ""},
		{"Arial (Windows)", "/Windows/Fonts/arial.ttf", ""},
		{"Calibri (Windows)", "/Windows/Fonts/calibri.ttf", ""},
		{"Arial (Windows Alt)", "C:/Windows/Fonts/arial.ttf", ""},
		{"Calibri (Windows Alt)", "C:/Windows/Fonts/calibri.ttf", ""},
	}
	
	for _, candidate := range candidates {
		if _, err := os.Stat(candidate.path); err == nil {
			// Detect format
			format, err := advancegg.GetFontFormat(candidate.path)
			if err == nil {
				candidate.format = format
				fonts = append(fonts, candidate)
			}
		}
	}
	
	return fonts
}
