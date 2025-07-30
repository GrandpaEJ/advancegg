package main

import (
	"fmt"
	"log"
	
	"github.com/GrandpaEJ/advancegg"
)

func main() {
	// Create a sample image
	dc := advancegg.NewContext(400, 300)
	
	// Create a colorful test image
	createSampleImage(dc)
	
	// Save in different formats
	formats := []struct {
		name string
		save func() error
	}{
		{"PNG", func() error { return dc.SavePNG("sample.png") }},
		{"JPEG", func() error { return dc.SaveJPEG("sample.jpg", 90) }},
		{"GIF", func() error { return dc.SaveGIF("sample.gif") }},
		{"BMP", func() error { return dc.SaveBMP("sample.bmp") }},
		{"TIFF", func() error { return dc.SaveTIFF("sample.tiff") }},
	}
	
	fmt.Println("Saving image in multiple formats:")
	for _, format := range formats {
		err := format.save()
		if err != nil {
			log.Printf("Error saving %s: %v", format.name, err)
		} else {
			fmt.Printf("✓ Saved as %s\n", format.name)
		}
	}
	
	// Test loading different formats
	fmt.Println("\nTesting image format loading:")
	
	loadTests := []struct {
		name string
		load func() error
	}{
		{"PNG", func() error { _, err := advancegg.LoadPNG("sample.png"); return err }},
		{"JPEG", func() error { _, err := advancegg.LoadJPG("sample.jpg"); return err }},
		{"GIF", func() error { _, err := advancegg.LoadGIF("sample.gif"); return err }},
		{"BMP", func() error { _, err := advancegg.LoadBMP("sample.bmp"); return err }},
		{"TIFF", func() error { _, err := advancegg.LoadTIFF("sample.tiff"); return err }},
	}
	
	for _, test := range loadTests {
		err := test.load()
		if err != nil {
			log.Printf("Error loading %s: %v", test.name, err)
		} else {
			fmt.Printf("✓ Loaded %s successfully\n", test.name)
		}
	}
	
	// Create a composite image showing format comparison
	createFormatComparison()
	
	fmt.Println("\nImage format demo completed!")
}

func createSampleImage(dc *advancegg.Context) {
	// Gradient background
	for y := 0; y < 300; y++ {
		for x := 0; x < 400; x++ {
			r := float64(x) / 400.0
			g := float64(y) / 300.0
			b := 0.5
			dc.SetRGB(r, g, b)
			dc.SetPixel(x, y)
		}
	}
	
	// Add some shapes
	dc.SetRGB(1, 1, 1)
	dc.DrawCircle(100, 75, 40)
	dc.Fill()
	
	dc.SetRGB(0, 0, 0)
	dc.DrawRectangle(200, 50, 80, 50)
	dc.Fill()
	
	dc.SetRGB(1, 0, 0)
	dc.DrawCircle(300, 75, 30)
	dc.Fill()
	
	// Add text
	dc.SetRGB(1, 1, 1)
	dc.DrawString("Format Test", 50, 200)
	dc.DrawString("PNG • JPEG • GIF • BMP • TIFF", 50, 230)
}

func createFormatComparison() {
	// Create a comparison image showing different format characteristics
	compDC := advancegg.NewContext(800, 600)
	compDC.SetRGB(0.9, 0.9, 0.9)
	compDC.Clear()
	
	// Title
	compDC.SetRGB(0, 0, 0)
	compDC.DrawString("Image Format Comparison", 50, 50)
	
	// Format information
	formats := []struct {
		name        string
		description string
		y           float64
	}{
		{"PNG", "Lossless compression, transparency support, best for graphics", 100},
		{"JPEG", "Lossy compression, smaller files, best for photos", 140},
		{"GIF", "Limited colors (256), animation support, very small files", 180},
		{"BMP", "Uncompressed, large files, universal compatibility", 220},
		{"TIFF", "Lossless/lossy options, professional photography", 260},
	}
	
	for _, format := range formats {
		compDC.SetRGB(0, 0, 0)
		compDC.DrawString(fmt.Sprintf("%s:", format.name), 70, format.y)
		compDC.SetRGB(0.3, 0.3, 0.3)
		compDC.DrawString(format.description, 150, format.y)
	}
	
	// Add usage recommendations
	compDC.SetRGB(0, 0, 0)
	compDC.DrawString("Usage Recommendations:", 50, 350)
	
	recommendations := []string{
		"• PNG: Graphics, logos, images with transparency",
		"• JPEG: Photographs, images with many colors",
		"• GIF: Simple animations, images with few colors",
		"• BMP: When file size is not a concern",
		"• TIFF: Professional photography, archival storage",
	}
	
	for i, rec := range recommendations {
		compDC.SetRGB(0.2, 0.2, 0.2)
		compDC.DrawString(rec, 70, 380+float64(i)*25)
	}
	
	compDC.SavePNG("format-comparison.png")
	fmt.Println("Format comparison saved as format-comparison.png")
}
