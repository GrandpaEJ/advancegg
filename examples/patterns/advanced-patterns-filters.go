package main

import (
	"fmt"
	"image"
	"image/color"
	"math"

	"github.com/GrandpaEJ/advancegg"
)

func main() {
	// Create a large canvas for demonstrating advanced patterns and filters
	dc := advancegg.NewContext(800, 600)

	// White background
	dc.SetRGB(1, 1, 1)
	dc.Clear()

	// Title
	dc.SetRGB(0, 0, 0)
	dc.DrawString("Advanced Patterns & Filters Demo", 20, 30)

	// Section 1: Advanced Patterns
	dc.DrawString("1. Advanced Patterns", 20, 70)

	// Linear gradient pattern
	linearGrad := advancegg.CreateLinearGradientPattern(100, 100,
		color.RGBA{255, 0, 0, 255},
		color.RGBA{0, 0, 255, 255})
	drawPatternSample(dc, linearGrad, 20, 90, 80, 60, "Linear Gradient")

	// Radial gradient pattern
	radialGrad := advancegg.CreateRadialGradientPattern(40, 40, 30,
		color.RGBA{255, 255, 0, 255},
		color.RGBA{255, 0, 255, 255})
	drawPatternSample(dc, radialGrad, 120, 90, 80, 60, "Radial Gradient")

	// Checkerboard pattern
	checkerboard := advancegg.CreateCheckerboard(8)
	drawPatternSample(dc, checkerboard, 220, 90, 80, 60, "Checkerboard")

	// Stripe pattern
	stripes := advancegg.CreateStripes(6)
	drawPatternSample(dc, stripes, 320, 90, 80, 60, "Stripes")

	// Polka dots pattern
	polkaDots := advancegg.CreatePolkaDots(15, 4)
	drawPatternSample(dc, polkaDots, 420, 90, 80, 60, "Polka Dots")

	// Section 2: Pattern Transforms
	dc.SetRGB(0, 0, 0)
	dc.DrawString("2. Pattern Transforms", 20, 200)

	// Original checkerboard
	originalPattern := advancegg.CreateCheckerboard(10)
	drawPatternSample(dc, originalPattern, 20, 220, 80, 60, "Original")

	// Translated pattern
	translatedPattern := advancegg.WithTranslation(originalPattern, 5, 5)
	drawPatternSample(dc, translatedPattern, 120, 220, 80, 60, "Translated")

	// Scaled pattern
	scaledPattern := advancegg.WithScale(originalPattern, 0.5, 0.5)
	drawPatternSample(dc, scaledPattern, 220, 220, 80, 60, "Scaled")

	// Rotated pattern
	rotatedPattern := advancegg.WithRotation(originalPattern, math.Pi/4)
	drawPatternSample(dc, rotatedPattern, 320, 220, 80, 60, "Rotated")

	// Section 3: CSS Filters
	dc.SetRGB(0, 0, 0)
	dc.DrawString("3. CSS Filters", 20, 330)

	// Create a base image for filtering
	baseImgCtx := createColorfulTestImage()
	baseImg := baseImgCtx.Image()

	// Original
	dc.DrawImage(baseImg, 20, 350)
	dc.DrawString("Original", 20, 410)

	// Brightness filter
	brightnessFilter := advancegg.BrightnessFilter{Amount: 1.5}
	brightImg := brightnessFilter.Apply(baseImg)
	dc.DrawImage(brightImg, 120, 350)
	dc.DrawString("Brightness", 120, 410)

	// Contrast filter
	contrastFilter := advancegg.ContrastFilter{Amount: 2.0}
	contrastImg := contrastFilter.Apply(baseImg)
	dc.DrawImage(contrastImg, 220, 350)
	dc.DrawString("Contrast", 220, 410)

	// Saturate filter
	saturateFilter := advancegg.SaturateFilter{Amount: 2.0}
	saturateImg := saturateFilter.Apply(baseImg)
	dc.DrawImage(saturateImg, 320, 350)
	dc.DrawString("Saturate", 320, 410)

	// Sepia filter
	sepiaFilter := advancegg.SepiaFilter{Amount: 1.0}
	sepiaImg := sepiaFilter.Apply(baseImg)
	dc.DrawImage(sepiaImg, 420, 350)
	dc.DrawString("Sepia", 420, 410)

	// Section 4: Filter Chains & Presets
	dc.SetRGB(0, 0, 0)
	dc.DrawString("4. Filter Presets", 20, 460)

	// Instagram preset
	instagramImg := advancegg.Instagram().Apply(baseImg)
	dc.DrawImage(instagramImg, 20, 480)
	dc.DrawString("Instagram", 20, 540)

	// Vintage preset
	vintageImg := advancegg.Vintage().Apply(baseImg)
	dc.DrawImage(vintageImg, 120, 480)
	dc.DrawString("Vintage", 120, 540)

	// Dramatic preset
	dramaticImg := advancegg.Dramatic().Apply(baseImg)
	dc.DrawImage(dramaticImg, 220, 480)
	dc.DrawString("Dramatic", 220, 540)

	// Black and white preset
	bwImg := advancegg.BlackAndWhite().Apply(baseImg)
	dc.DrawImage(bwImg, 320, 480)
	dc.DrawString("B&W", 320, 540)

	// Custom filter chain
	customChain := advancegg.NewFilterChain().
		Add(advancegg.BrightnessFilter{Amount: 1.2}).
		Add(advancegg.ContrastFilter{Amount: 1.3}).
		Add(advancegg.SaturateFilter{Amount: 0.8})
	customImg := customChain.Apply(baseImg)
	dc.DrawImage(customImg, 420, 480)
	dc.DrawString("Custom Chain", 420, 540)

	// Save the result
	advancegg.SavePNG("images/patterns/advanced-patterns-filters-demo.png", dc.Image())
	fmt.Println("Advanced patterns and filters demo saved to: advanced-patterns-filters-demo.png")
}

// Helper function to draw a pattern sample
func drawPatternSample(dc *advancegg.Context, pattern advancegg.AdvancedPattern, x, y, w, h float64, label string) {
	// Create a small image and fill it with the pattern
	img := advancegg.NewContext(int(w), int(h))
	if rgba, ok := img.Image().(*image.RGBA); ok {
		advancegg.PatternFill(rgba, pattern)
	}

	// Draw the pattern sample
	dc.DrawImage(img.Image(), int(x), int(y))

	// Draw label
	dc.SetRGB(0, 0, 0)
	dc.DrawString(label, x, y+h+15)
}

// Helper function to create a colorful test image
func createColorfulTestImage() *advancegg.Context {
	img := advancegg.NewContext(80, 50)

	// Create a gradient background
	for y := 0; y < 50; y++ {
		for x := 0; x < 80; x++ {
			r := float64(x) / 80.0
			g := float64(y) / 50.0
			b := 0.5
			img.SetRGB(r, g, b)
			img.SetPixel(x, y)
		}
	}

	// Add some shapes
	img.SetRGB(1, 1, 0) // Yellow
	img.DrawCircle(20, 15, 8)
	img.Fill()

	img.SetRGB(0, 1, 1) // Cyan
	img.DrawRectangle(50, 10, 20, 15)
	img.Fill()

	img.SetRGB(1, 0, 1) // Magenta
	img.MoveTo(10, 35)
	img.LineTo(30, 25)
	img.LineTo(30, 45)
	img.ClosePath()
	img.Fill()

	return img
}
