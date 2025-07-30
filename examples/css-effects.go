package main

import (
	"fmt"
	"image"
	"image/color"
	"math"

	"github.com/GrandpaEJ/advancegg"
	"github.com/GrandpaEJ/advancegg/internal/advance"
)

func main() {
	// Create a sample image to apply effects to
	dc := advancegg.NewContext(400, 300)
	createSampleImage(dc)
	dc.SavePNG("css-original.png")

	// Demonstrate CSS-like filters
	demonstrateCSSFilters(dc)

	// Demonstrate patterns
	demonstratePatterns()

	fmt.Println("CSS-like effects examples completed!")
}

func createSampleImage(dc *advancegg.Context) {
	// Create a colorful test image
	dc.SetRGB(0.3, 0.6, 0.9)
	dc.Clear()

	// Add some shapes
	dc.SetRGB(1, 0.2, 0.2)
	dc.DrawCircle(100, 75, 40)
	dc.Fill()

	dc.SetRGB(0.2, 1, 0.2)
	dc.DrawRectangle(200, 50, 80, 50)
	dc.Fill()

	dc.SetRGB(1, 1, 0.2)
	dc.DrawCircle(300, 75, 30)
	dc.Fill()

	// Add some text
	dc.SetRGB(1, 1, 1)
	dc.DrawString("CSS Effects Demo", 50, 200)
	dc.DrawString("Filters & Patterns", 50, 230)
}

func demonstrateCSSFilters(dc *advancegg.Context) {
	fmt.Println("Demonstrating CSS-like filters...")

	originalImage := dc.Image()

	// Individual filters
	filters := []struct {
		name   string
		filter advance.CSSFilter
	}{
		{"brightness", advance.BrightnessFilter{Amount: 1.5}},
		{"contrast", advance.ContrastFilter{Amount: 1.5}},
		{"saturate", advance.SaturateFilter{Amount: 2.0}},
		{"hue-rotate", advance.HueRotateFilter{Degrees: 90}},
		{"invert", advance.InvertFilter{Amount: 1.0}},
		{"opacity", advance.OpacityFilter{Amount: 0.7}},
		{"blur", advance.BlurFilter{Radius: 3}},
		{"sepia", advance.SepiaFilter{Amount: 1.0}},
	}

	for _, f := range filters {
		filtered := f.filter.Apply(originalImage)

		newDC := advancegg.NewContextForRGBA(filtered.(*image.RGBA))
		newDC.SavePNG(fmt.Sprintf("css-%s.png", f.name))
		fmt.Printf("Applied %s filter, saved as css-%s.png\n", f.name, f.name)
	}

	// Preset filter combinations
	presets := []struct {
		name   string
		filter *advance.FilterChain
	}{
		{"instagram", advance.Instagram()},
		{"vintage", advance.Vintage()},
		{"dramatic", advance.Dramatic()},
		{"blackandwhite", advance.BlackAndWhite()},
		{"warm", advance.Warm()},
		{"cool", advance.Cool()},
	}

	for _, preset := range presets {
		filtered := preset.filter.Apply(originalImage)

		newDC := advancegg.NewContextForRGBA(filtered.(*image.RGBA))
		newDC.SavePNG(fmt.Sprintf("css-%s.png", preset.name))
		fmt.Printf("Applied %s preset, saved as css-%s.png\n", preset.name, preset.name)
	}

	// Custom filter chain
	customChain := advance.NewFilterChain().
		Add(advance.ContrastFilter{Amount: 1.2}).
		Add(advance.SaturateFilter{Amount: 1.3}).
		Add(advance.HueRotateFilter{Degrees: 30}).
		Add(advance.BrightnessFilter{Amount: 1.1})

	customFiltered := customChain.Apply(originalImage)
	customDC := advancegg.NewContextForRGBA(customFiltered.(*image.RGBA))
	customDC.SavePNG("css-custom-chain.png")
	fmt.Println("Applied custom filter chain, saved as css-custom-chain.png")
}

func demonstratePatterns() {
	fmt.Println("Demonstrating patterns...")

	// Create pattern examples
	patterns := []struct {
		name    string
		pattern advance.Pattern
	}{
		{"linear-gradient", advance.CreateLinearGradient(400, 300,
			color.RGBA{255, 0, 0, 255},   // Red
			color.RGBA{255, 255, 0, 255}, // Yellow
			color.RGBA{0, 255, 0, 255},   // Green
			color.RGBA{0, 255, 255, 255}, // Cyan
			color.RGBA{0, 0, 255, 255},   // Blue
		)},
		{"radial-gradient", advance.CreateRadialGradient(200, 150, 150,
			color.RGBA{255, 255, 255, 255}, // White center
			color.RGBA{255, 0, 0, 255},     // Red
			color.RGBA{0, 0, 0, 255},       // Black edge
		)},
		{"checkerboard", advance.CreateCheckerboard(20)},
		{"stripes", advance.CreateStripes(15)},
		{"polka-dots", advance.CreatePolkaDots(40, 15)},
		{"noise", advance.NoisePattern{
			Scale:     0.1,
			BaseColor: color.RGBA{100, 150, 200, 255},
			Intensity: 0.5,
		}},
		{"waves", advance.WavePattern{
			Wavelength: 50,
			Amplitude:  20,
			Angle:      0,
			Color1:     color.RGBA{0, 100, 200, 255},
			Color2:     color.RGBA{200, 100, 0, 255},
		}},
	}

	for _, p := range patterns {
		dc := advancegg.NewContext(400, 300)

		// Fill with pattern
		advance.PatternFill(dc.Image().(*image.RGBA), p.pattern)

		dc.SavePNG(fmt.Sprintf("pattern-%s.png", p.name))
		fmt.Printf("Created %s pattern, saved as pattern-%s.png\n", p.name, p.name)
	}

	// Create a composite showing all patterns
	createPatternComposite()
}

func createPatternComposite() {
	// Create a grid showing multiple patterns
	dc := advancegg.NewContext(800, 600)
	dc.SetRGB(0.9, 0.9, 0.9)
	dc.Clear()

	// Title
	dc.SetRGB(0, 0, 0)
	dc.DrawString("Pattern Gallery", 50, 50)

	patterns := []struct {
		name    string
		pattern advance.Pattern
		x, y    float64
	}{
		{"Linear Gradient", advance.CreateLinearGradient(150, 100,
			color.RGBA{255, 0, 0, 255},
			color.RGBA{0, 0, 255, 255},
		), 50, 80},
		{"Radial Gradient", advance.CreateRadialGradient(75, 50, 60,
			color.RGBA{255, 255, 0, 255},
			color.RGBA{255, 0, 0, 255},
		), 250, 80},
		{"Checkerboard", advance.CreateCheckerboard(10), 450, 80},
		{"Stripes", advance.CreateStripes(8), 650, 80},
		{"Polka Dots", advance.CreatePolkaDots(25, 8), 50, 220},
		{"Noise", advance.NoisePattern{
			Scale:     0.05,
			BaseColor: color.RGBA{150, 100, 200, 255},
			Intensity: 0.3,
		}, 250, 220},
		{"Waves", advance.WavePattern{
			Wavelength: 30,
			Amplitude:  15,
			Angle:      math.Pi / 6,
			Color1:     color.RGBA{0, 150, 255, 255},
			Color2:     color.RGBA{255, 150, 0, 255},
		}, 450, 220},
	}

	for _, p := range patterns {
		// Create small pattern sample
		sampleDC := advancegg.NewContext(150, 100)
		advance.PatternFill(sampleDC.Image().(*image.RGBA), p.pattern)

		// Draw the pattern sample
		dc.DrawImage(sampleDC.Image(), int(p.x), int(p.y))

		// Add label
		dc.SetRGB(0, 0, 0)
		dc.DrawString(p.name, p.x, p.y+120)
	}

	dc.SavePNG("pattern-gallery.png")
	fmt.Println("Pattern gallery saved as pattern-gallery.png")
}

// Helper function to create context from RGBA image
func NewContextForRGBA(img *image.RGBA) *advancegg.Context {
	bounds := img.Bounds()
	dc := advancegg.NewContext(bounds.Max.X, bounds.Max.Y)

	// Copy the image data
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := img.RGBAAt(x, y)
			dc.SetRGBA255(int(c.R), int(c.G), int(c.B), int(c.A))
			dc.SetPixel(x, y)
		}
	}

	return dc
}
