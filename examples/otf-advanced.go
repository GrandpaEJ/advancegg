package main

import (
	"fmt"
	"log"
	"os"
	
	"github.com/GrandpaEJ/advancegg"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

func main() {
	dc := advancegg.NewContext(800, 600)
	
	// White background
	dc.SetRGB(1, 1, 1)
	dc.Clear()
	
	// Try to load an OTF font with advanced options
	y := 50.0
	
	// Title
	dc.SetRGB(0, 0, 0)
	dc.DrawString("Advanced OTF Font Features", 50, y)
	y += 60
	
	// Demonstrate loading with custom options
	fontPath := findAvailableFont()
	if fontPath == "" {
		dc.DrawString("No suitable font found for demo", 50, y)
		dc.SavePNG("otf-advanced.png")
		return
	}
	
	// Load font with different hinting options
	hintingOptions := []struct {
		name    string
		hinting font.Hinting
		color   [3]float64
	}{
		{"No Hinting", font.HintingNone, [3]float64{1, 0, 0}},
		{"Vertical Hinting", font.HintingVertical, [3]float64{0, 1, 0}},
		{"Full Hinting", font.HintingFull, [3]float64{0, 0, 1}},
	}
	
	for _, option := range hintingOptions {
		// Load font with custom options
		fontBytes, err := os.ReadFile(fontPath)
		if err != nil {
			log.Printf("Error reading font: %v", err)
			continue
		}
		
		face, err := advancegg.ParseFontFaceWithOptions(fontBytes, &truetype.Options{
			Size:    24,
			DPI:     72,
			Hinting: option.hinting,
		})
		if err != nil {
			log.Printf("Error parsing font: %v", err)
			continue
		}
		
		dc.SetFontFace(face)
		dc.SetRGB(option.color[0], option.color[1], option.color[2])
		
		text := fmt.Sprintf("%s: The quick brown fox jumps", option.name)
		dc.DrawString(text, 50, y)
		y += 35
	}
	
	y += 20
	
	// Demonstrate different font sizes
	dc.SetRGB(0, 0, 0)
	dc.DrawString("Font Size Variations:", 50, y)
	y += 40
	
	sizes := []float64{12, 16, 20, 24, 32, 48}
	for _, size := range sizes {
		err := dc.LoadFontFace(fontPath, size)
		if err != nil {
			continue
		}
		
		text := fmt.Sprintf("Size %.0f: Sample Text", size)
		dc.DrawString(text, 50, y)
		y += size + 5
	}
	
	// Demonstrate font metrics
	y += 20
	dc.LoadFontFace(fontPath, 20)
	dc.DrawString("Font Metrics Demo:", 50, y)
	y += 30
	
	testText := "Measuring this text"
	w, h := dc.MeasureString(testText)
	
	// Draw the text
	dc.SetRGB(0, 0, 1)
	dc.DrawString(testText, 50, y)
	
	// Draw bounding box
	dc.SetRGB(1, 0, 0)
	dc.SetLineWidth(1)
	dc.DrawRectangle(50, y-h, w, h)
	dc.Stroke()
	
	// Add metrics info
	dc.SetRGB(0, 0, 0)
	dc.LoadFontFace(fontPath, 14)
	metricsText := fmt.Sprintf("Width: %.1f, Height: %.1f", w, h)
	dc.DrawString(metricsText, 50, y+25)
	
	// Save the image
	dc.SavePNG("otf-advanced.png")
	fmt.Println("Advanced OTF demo saved as otf-advanced.png")
}

func findAvailableFont() string {
	// List of common font paths to try
	fontPaths := []string{
		// macOS
		"/System/Fonts/Arial.ttf",
		"/System/Fonts/Helvetica.ttc",
		"/System/Library/Fonts/Avenir.ttc",
		"/System/Library/Fonts/Times.ttc",
		
		// Linux
		"/usr/share/fonts/truetype/dejavu/DejaVuSans.ttf",
		"/usr/share/fonts/truetype/liberation/LiberationSans-Regular.ttf",
		"/usr/share/fonts/TTF/arial.ttf",
		
		// Windows
		"/Windows/Fonts/arial.ttf",
		"/Windows/Fonts/calibri.ttf",
		"/Windows/Fonts/times.ttf",
		
		// Alternative Windows paths
		"C:/Windows/Fonts/arial.ttf",
		"C:/Windows/Fonts/calibri.ttf",
	}
	
	for _, path := range fontPaths {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}
	
	return ""
}
