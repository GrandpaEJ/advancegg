package main

import (
	"fmt"
	"image/color"
	
	"github.com/GrandpaEJ/advancegg"
)

func main() {
	fmt.Println("Creating ICC Color Profile examples...")
	
	// Color space comparison
	createColorSpaceComparison()
	
	// Color profile conversion
	createProfileConversion()
	
	// Print-ready workflow
	createPrintWorkflow()
	
	fmt.Println("ICC Color Profile examples completed!")
}

func createColorSpaceComparison() {
	dc := advancegg.NewContext(1200, 800)
	
	// White background
	dc.SetRGB(1, 1, 1)
	dc.Clear()
	
	// Title
	dc.SetRGB(0, 0, 0)
	dc.DrawString("Color Space Comparison: sRGB vs Adobe RGB", 50, 50)
	
	// Create color profiles
	srgbProfile := advancegg.CreateSRGBProfile()
	adobeProfile := advancegg.CreateAdobeRGBProfile()
	
	// Set initial profile to sRGB
	dc.SetColorProfile(srgbProfile)
	
	// Draw color swatches in sRGB
	dc.SetRGB(0, 0, 0)
	dc.DrawString("sRGB Color Space:", 50, 100)
	
	colors := []struct {
		name string
		r, g, b float64
	}{
		{"Pure Red", 1, 0, 0},
		{"Pure Green", 0, 1, 0},
		{"Pure Blue", 0, 0, 1},
		{"Cyan", 0, 1, 1},
		{"Magenta", 1, 0, 1},
		{"Yellow", 1, 1, 0},
		{"Orange", 1, 0.5, 0},
		{"Purple", 0.5, 0, 1},
	}
	
	x := 50.0
	for _, col := range colors {
		// Draw sRGB version
		dc.SetRGB(col.r, col.g, col.b)
		dc.DrawRectangle(x, 130, 80, 60)
		dc.Fill()
		
		// Label
		dc.SetRGB(0, 0, 0)
		dc.DrawString(col.name, x, 210)
		
		x += 100
	}
	
	// Draw same colors in Adobe RGB space
	dc.SetRGB(0, 0, 0)
	dc.DrawString("Adobe RGB Color Space:", 50, 250)
	
	// Set profile to Adobe RGB
	dc.SetColorProfile(adobeProfile)
	
	x = 50.0
	for _, col := range colors {
		// Draw Adobe RGB version
		dc.SetRGB(col.r, col.g, col.b)
		dc.DrawRectangle(x, 280, 80, 60)
		dc.Fill()
		
		// Label
		dc.SetRGB(0, 0, 0)
		dc.DrawString(col.name, x, 360)
		
		x += 100
	}
	
	// Color gamut visualization
	dc.SetRGB(0, 0, 0)
	dc.DrawString("Color Gamut Visualization:", 50, 420)
	
	// Draw simplified gamut triangles
	drawGamutTriangle(dc, 100, 500, 150, "sRGB", color.RGBA{255, 100, 100, 100})
	drawGamutTriangle(dc, 300, 500, 180, "Adobe RGB", color.RGBA{100, 255, 100, 100})
	
	// Profile information
	dc.SetRGB(0, 0, 0)
	dc.DrawString("Profile Information:", 50, 650)
	
	srgbInfo := fmt.Sprintf("sRGB: White Point D65 (%.3f, %.3f, %.3f)", 
		srgbProfile.WhitePoint.X, srgbProfile.WhitePoint.Y, srgbProfile.WhitePoint.Z)
	dc.DrawString(srgbInfo, 70, 680)
	
	adobeInfo := fmt.Sprintf("Adobe RGB: White Point D65 (%.3f, %.3f, %.3f)", 
		adobeProfile.WhitePoint.X, adobeProfile.WhitePoint.Y, adobeProfile.WhitePoint.Z)
	dc.DrawString(adobeInfo, 70, 700)
	
	dc.DrawString("Note: Wider gamut allows more saturated colors", 70, 730)
	
	dc.SavePNG("images/color-space-comparison.png")
	fmt.Println("Color space comparison saved as color-space-comparison.png")
}

func drawGamutTriangle(dc *advancegg.Context, x, y, size float64, label string, fillColor color.RGBA) {
	// Draw triangle representing color gamut
	dc.SetColor(fillColor)
	
	// Triangle points (simplified representation)
	dc.MoveTo(x, y)
	dc.LineTo(x+size, y)
	dc.LineTo(x+size/2, y-size*0.866) // Equilateral triangle
	dc.ClosePath()
	dc.Fill()
	
	// Outline
	dc.SetRGB(0, 0, 0)
	dc.SetLineWidth(2)
	dc.MoveTo(x, y)
	dc.LineTo(x+size, y)
	dc.LineTo(x+size/2, y-size*0.866)
	dc.ClosePath()
	dc.Stroke()
	
	// Label
	dc.DrawString(label, x+size/2-20, y+20)
}

func createProfileConversion() {
	dc := advancegg.NewContext(1000, 600)
	
	// White background
	dc.SetRGB(1, 1, 1)
	dc.Clear()
	
	// Title
	dc.SetRGB(0, 0, 0)
	dc.DrawString("Color Profile Conversion Demo", 50, 50)
	
	// Create profiles
	srgbProfile := advancegg.CreateSRGBProfile()
	adobeProfile := advancegg.CreateAdobeRGBProfile()
	
	// Create color converter
	converter := advancegg.NewColorConverter(srgbProfile, adobeProfile)
	dc.SetColorConverter(converter)
	
	// Original image in sRGB
	dc.SetColorProfile(srgbProfile)
	dc.SetRGB(0, 0, 0)
	dc.DrawString("Original (sRGB):", 50, 100)
	
	// Draw test image
	createTestImage(dc, 50, 130, 400, 200)
	
	// Convert to Adobe RGB
	dc.SetRGB(0, 0, 0)
	dc.DrawString("Converted to Adobe RGB:", 500, 100)
	
	// Create second test image and convert it
	createTestImage(dc, 500, 130, 400, 200)
	dc.ConvertToColorSpace(adobeProfile)
	
	// Conversion information
	dc.SetRGB(0, 0, 0)
	dc.DrawString("Color Conversion Process:", 50, 380)
	dc.DrawString("1. Source RGB → XYZ (using source profile)", 70, 410)
	dc.DrawString("2. Chromatic adaptation (if white points differ)", 70, 430)
	dc.DrawString("3. XYZ → Destination RGB (using dest profile)", 70, 450)
	dc.DrawString("4. Apply tone curves and gamma correction", 70, 470)
	
	// Rendering intents
	dc.DrawString("Rendering Intents:", 50, 510)
	intents := []string{
		"Perceptual: Preserves overall appearance",
		"Relative Colorimetric: Preserves color accuracy",
		"Saturation: Preserves color saturation",
		"Absolute Colorimetric: Preserves absolute colors",
	}
	
	for i, intent := range intents {
		dc.DrawString(fmt.Sprintf("• %s", intent), 70, 530+float64(i)*20)
	}
	
	dc.SavePNG("images/profile-conversion-demo.png")
	fmt.Println("Profile conversion demo saved as profile-conversion-demo.png")
}

func createTestImage(dc *advancegg.Context, x, y, width, height float64) {
	// Create a test image with various colors
	cellWidth := width / 8
	cellHeight := height / 4
	
	testColors := [][]color.RGBA{
		{
			{255, 0, 0, 255},   // Red
			{255, 128, 0, 255}, // Orange
			{255, 255, 0, 255}, // Yellow
			{128, 255, 0, 255}, // Yellow-green
			{0, 255, 0, 255},   // Green
			{0, 255, 128, 255}, // Green-cyan
			{0, 255, 255, 255}, // Cyan
			{0, 128, 255, 255}, // Cyan-blue
		},
		{
			{0, 0, 255, 255},   // Blue
			{128, 0, 255, 255}, // Blue-magenta
			{255, 0, 255, 255}, // Magenta
			{255, 0, 128, 255}, // Magenta-red
			{128, 128, 128, 255}, // Gray
			{192, 192, 192, 255}, // Light gray
			{64, 64, 64, 255},    // Dark gray
			{0, 0, 0, 255},       // Black
		},
		{
			{255, 192, 192, 255}, // Light red
			{192, 255, 192, 255}, // Light green
			{192, 192, 255, 255}, // Light blue
			{255, 255, 192, 255}, // Light yellow
			{255, 192, 255, 255}, // Light magenta
			{192, 255, 255, 255}, // Light cyan
			{255, 255, 255, 255}, // White
			{128, 64, 0, 255},    // Brown
		},
		{
			{64, 0, 0, 255},     // Dark red
			{0, 64, 0, 255},     // Dark green
			{0, 0, 64, 255},     // Dark blue
			{64, 64, 0, 255},    // Dark yellow
			{64, 0, 64, 255},    // Dark magenta
			{0, 64, 64, 255},    // Dark cyan
			{255, 128, 64, 255}, // Skin tone
			{128, 64, 32, 255},  // Wood
		},
	}
	
	for row := 0; row < 4; row++ {
		for col := 0; col < 8; col++ {
			cellX := x + float64(col)*cellWidth
			cellY := y + float64(row)*cellHeight
			
			dc.SetColor(testColors[row][col])
			dc.DrawRectangle(cellX, cellY, cellWidth, cellHeight)
			dc.Fill()
		}
	}
}

func createPrintWorkflow() {
	dc := advancegg.NewContext(1000, 700)
	
	// White background
	dc.SetRGB(1, 1, 1)
	dc.Clear()
	
	// Title
	dc.SetRGB(0, 0, 0)
	dc.DrawString("Print-Ready Color Workflow", 50, 50)
	
	// Workflow steps
	steps := []struct {
		title string
		desc  string
		y     float64
	}{
		{"1. Design in RGB", "Create artwork in sRGB or Adobe RGB", 100},
		{"2. Soft Proofing", "Preview how colors will look when printed", 180},
		{"3. Profile Conversion", "Convert to printer's ICC profile", 260},
		{"4. Gamut Mapping", "Handle out-of-gamut colors", 340},
		{"5. Print Output", "Final print with accurate colors", 420},
	}
	
	for _, step := range steps {
		// Step box
		dc.SetRGB(0.9, 0.9, 1)
		dc.DrawRectangle(50, step.y, 900, 60)
		dc.Fill()
		
		dc.SetRGB(0, 0, 0.8)
		dc.SetLineWidth(2)
		dc.DrawRectangle(50, step.y, 900, 60)
		dc.Stroke()
		
		// Step text
		dc.SetRGB(0, 0, 0)
		dc.DrawString(step.title, 70, step.y+25)
		dc.DrawString(step.desc, 70, step.y+45)
	}
	
	// Color management benefits
	dc.SetRGB(0, 0, 0)
	dc.DrawString("Benefits of ICC Color Management:", 50, 520)
	
	benefits := []string{
		"• Consistent colors across devices",
		"• Predictable print output",
		"• Reduced waste from color mismatches",
		"• Professional color accuracy",
		"• Support for wide-gamut displays and printers",
	}
	
	for i, benefit := range benefits {
		dc.DrawString(benefit, 70, 550+float64(i)*25)
	}
	
	// Profile examples
	dc.SetRGB(0, 0, 0)
	dc.DrawString("Common ICC Profiles:", 500, 520)
	
	profiles := []string{
		"• sRGB: Standard web/monitor profile",
		"• Adobe RGB: Photography workflow",
		"• ProPhoto RGB: Wide gamut editing",
		"• CMYK: Printing press profiles",
		"• Grayscale: Black & white printing",
	}
	
	for i, profile := range profiles {
		dc.DrawString(profile, 520, 550+float64(i)*25)
	}
	
	dc.SavePNG("images/print-workflow-demo.png")
	fmt.Println("Print workflow demo saved as print-workflow-demo.png")
}

func createAdvancedColorExample() {
	dc := advancegg.NewContext(1200, 800)
	
	// Gradient background
	for y := 0; y < 800; y++ {
		t := float64(y) / 800.0
		dc.SetRGB(0.95+t*0.05, 0.95+t*0.05, 1.0)
		dc.DrawLine(0, float64(y), 1200, float64(y))
		dc.Stroke()
	}
	
	// Title
	dc.SetRGB(0, 0, 0)
	dc.DrawString("Advanced Color Management Features", 50, 50)
	
	// Create multiple profiles for comparison
	srgbProfile := advancegg.CreateSRGBProfile()
	adobeProfile := advancegg.CreateAdobeRGBProfile()
	
	// Demonstrate different rendering intents
	intents := []advancegg.RenderingIntent{
		advancegg.IntentPerceptual,
		advancegg.IntentRelativeColorimetric,
		advancegg.IntentSaturation,
		advancegg.IntentAbsoluteColorimetric,
	}
	
	intentNames := []string{
		"Perceptual",
		"Relative Colorimetric",
		"Saturation",
		"Absolute Colorimetric",
	}
	
	// Draw color patches with different intents
	dc.SetRGB(0, 0, 0)
	dc.DrawString("Rendering Intent Comparison:", 50, 100)
	
	for i, intent := range intents {
		converter := advancegg.NewColorConverter(srgbProfile, adobeProfile)
		converter.Intent = intent
		
		x := 50 + float64(i)*280
		
		// Intent label
		dc.SetRGB(0, 0, 0)
		dc.DrawString(intentNames[i], x, 130)
		
		// Draw test colors
		testColors := []color.RGBA{
			{255, 0, 0, 255},   // Saturated red
			{0, 255, 0, 255},   // Saturated green
			{0, 0, 255, 255},   // Saturated blue
			{255, 255, 0, 255}, // Yellow
		}
		
		for j, testColor := range testColors {
			// Original color
			dc.SetColor(testColor)
			dc.DrawRectangle(x, 150+float64(j)*40, 30, 30)
			dc.Fill()
			
			// Converted color
			convertedColor := converter.ConvertColor(testColor)
			dc.SetColor(convertedColor)
			dc.DrawRectangle(x+40, 150+float64(j)*40, 30, 30)
			dc.Fill()
			
			// Arrow
			dc.SetRGB(0, 0, 0)
			dc.DrawString("→", x+32, 170+float64(j)*40)
		}
	}
	
	// White point adaptation example
	dc.SetRGB(0, 0, 0)
	dc.DrawString("White Point Adaptation:", 50, 350)
	
	// Show different illuminants
	illuminants := []struct {
		name string
		wp   advancegg.XYZColor
	}{
		{"D50", advancegg.XYZColor{X: 0.9642, Y: 1.0000, Z: 0.8249}},
		{"D65", advancegg.XYZColor{X: 0.9505, Y: 1.0000, Z: 1.0890}},
	}
	
	for i, illum := range illuminants {
		x := 50 + float64(i)*200
		dc.DrawString(fmt.Sprintf("%s: (%.3f, %.3f, %.3f)", 
			illum.name, illum.wp.X, illum.wp.Y, illum.wp.Z), x, 380)
	}
	
	// Gamut mapping visualization
	dc.SetRGB(0, 0, 0)
	dc.DrawString("Gamut Mapping Strategies:", 50, 450)
	
	strategies := []string{
		"Clipping: Hard limit to gamut boundary",
		"Compression: Scale all colors proportionally", 
		"Perceptual: Preserve relationships between colors",
		"Saturation: Maintain color vividness",
	}
	
	for i, strategy := range strategies {
		dc.DrawString(fmt.Sprintf("• %s", strategy), 70, 480+float64(i)*25)
	}
	
	// Color accuracy metrics
	dc.SetRGB(0, 0, 0)
	dc.DrawString("Color Accuracy Metrics:", 50, 600)
	
	metrics := []string{
		"ΔE76: CIE 1976 color difference formula",
		"ΔE94: Improved perceptual uniformity",
		"ΔE2000: Most accurate modern formula",
		"ΔE < 1: Not perceptible to human eye",
		"ΔE 1-3: Perceptible but acceptable",
		"ΔE > 3: Clearly visible difference",
	}
	
	for i, metric := range metrics {
		dc.DrawString(fmt.Sprintf("• %s", metric), 70, 630+float64(i)*20)
	}
	
	dc.SavePNG("images/advanced-color-demo.png")
	fmt.Println("Advanced color demo saved as advanced-color-demo.png")
}
