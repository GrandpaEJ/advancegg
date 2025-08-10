package advance

import (
	"image"
	"image/color"
	"testing"
)

// Helper function to create a test image
func createTestImage(width, height int, c color.RGBA) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.SetRGBA(x, y, c)
		}
	}
	return img
}

func TestFilterChain_Creation(t *testing.T) {
	chain := NewFilterChain()
	if chain == nil {
		t.Fatal("NewFilterChain() returned nil")
	}

	if len(chain.filters) != 0 {
		t.Errorf("New filter chain should be empty, got %d filters", len(chain.filters))
	}
}

func TestFilterChain_Add(t *testing.T) {
	chain := NewFilterChain()
	brightness := BrightnessFilter{Amount: 1.5}

	chain.Add(brightness)

	if len(chain.filters) != 1 {
		t.Errorf("Expected 1 filter after adding, got %d", len(chain.filters))
	}
}

func TestFilterChain_Apply(t *testing.T) {
	// Create a test image
	img := createTestImage(10, 10, color.RGBA{128, 128, 128, 255})

	// Create filter chain
	chain := NewFilterChain()
	chain.Add(BrightnessFilter{Amount: 2.0}) // Double brightness

	// Apply filters
	result := chain.Apply(img)

	if result == nil {
		t.Fatal("Filter chain returned nil")
	}

	// Check that the result is different from the original
	resultRGBA := result.(*image.RGBA)
	originalColor := img.RGBAAt(0, 0)
	resultColor := resultRGBA.RGBAAt(0, 0)

	if originalColor == resultColor {
		t.Error("Filter should have changed the image")
	}
}

func TestBrightnessFilter(t *testing.T) {
	img := createTestImage(10, 10, color.RGBA{100, 100, 100, 255})
	filter := BrightnessFilter{Amount: 2.0}

	result := filter.Apply(img)
	resultRGBA := result.(*image.RGBA)
	resultColor := resultRGBA.RGBAAt(0, 0)

	// Brightness should increase the RGB values
	if resultColor.R <= 100 || resultColor.G <= 100 || resultColor.B <= 100 {
		t.Errorf("Brightness filter should increase RGB values, got %+v", resultColor)
	}
}

func TestContrastFilter(t *testing.T) {
	img := createTestImage(10, 10, color.RGBA{100, 100, 100, 255}) // Use non-middle gray
	filter := ContrastFilter{Amount: 2.0}

	result := filter.Apply(img)
	resultRGBA := result.(*image.RGBA)
	resultColor := resultRGBA.RGBAAt(0, 0)

	// Contrast should change the color
	originalColor := img.RGBAAt(0, 0)
	if resultColor.R == originalColor.R && resultColor.G == originalColor.G && resultColor.B == originalColor.B {
		t.Error("Contrast filter should change the color")
	}
}

func TestSaturateFilter(t *testing.T) {
	// Create a colored image
	img := createTestImage(10, 10, color.RGBA{200, 100, 100, 255})
	filter := SaturateFilter{Amount: 0.0} // Desaturate completely

	result := filter.Apply(img)
	resultRGBA := result.(*image.RGBA)
	resultColor := resultRGBA.RGBAAt(0, 0)

	// Complete desaturation should make R, G, B equal (grayscale)
	if resultColor.R != resultColor.G || resultColor.G != resultColor.B {
		t.Errorf("Complete desaturation should produce grayscale, got %+v", resultColor)
	}
}

func TestHueRotateFilter(t *testing.T) {
	img := createTestImage(10, 10, color.RGBA{255, 0, 0, 255}) // Pure red
	filter := HueRotateFilter{Degrees: 120}                    // Rotate to green

	result := filter.Apply(img)
	resultRGBA := result.(*image.RGBA)
	resultColor := resultRGBA.RGBAAt(0, 0)

	// Hue rotation should change the color
	originalColor := img.RGBAAt(0, 0)
	if resultColor == originalColor {
		t.Error("Hue rotation should change the color")
	}
}

func TestInvertFilter(t *testing.T) {
	img := createTestImage(10, 10, color.RGBA{100, 150, 200, 255})
	filter := InvertFilter{Amount: 1.0} // Full inversion

	result := filter.Apply(img)
	resultRGBA := result.(*image.RGBA)
	resultColor := resultRGBA.RGBAAt(0, 0)

	// Invert filter should change the color significantly
	originalColor := img.RGBAAt(0, 0)
	if resultColor.R == originalColor.R && resultColor.G == originalColor.G && resultColor.B == originalColor.B {
		t.Error("Invert filter should change the color")
	}

	// For full inversion, the result should be significantly different
	// We'll just check that the color changed significantly
	rDiff := int(resultColor.R) - int(originalColor.R)
	if rDiff < 50 && rDiff > -50 {
		t.Errorf("Invert filter should change color significantly: original R=%d, result R=%d",
			originalColor.R, resultColor.R)
	}
}

func TestOpacityFilter(t *testing.T) {
	img := createTestImage(10, 10, color.RGBA{255, 255, 255, 255})
	filter := OpacityFilter{Amount: 0.5} // Half opacity

	result := filter.Apply(img)
	resultRGBA := result.(*image.RGBA)
	resultColor := resultRGBA.RGBAAt(0, 0)

	// Alpha should be reduced
	if resultColor.A >= 255 {
		t.Errorf("Opacity filter should reduce alpha, got %d", resultColor.A)
	}
}

func TestSepiaFilter(t *testing.T) {
	img := createTestImage(10, 10, color.RGBA{100, 200, 50, 255})
	filter := SepiaFilter{Amount: 1.0} // Full sepia

	result := filter.Apply(img)
	resultRGBA := result.(*image.RGBA)
	resultColor := resultRGBA.RGBAAt(0, 0)

	// Sepia should give a brownish tint
	originalColor := img.RGBAAt(0, 0)
	if resultColor == originalColor {
		t.Error("Sepia filter should change the color")
	}
}

func TestBlurFilter(t *testing.T) {
	// Create an image with a sharp edge
	img := image.NewRGBA(image.Rect(0, 0, 10, 10))
	for y := 0; y < 10; y++ {
		for x := 0; x < 10; x++ {
			if x < 5 {
				img.SetRGBA(x, y, color.RGBA{0, 0, 0, 255})
			} else {
				img.SetRGBA(x, y, color.RGBA{255, 255, 255, 255})
			}
		}
	}

	filter := BlurFilter{Radius: 1.0}
	result := filter.Apply(img)

	if result == nil {
		t.Fatal("Blur filter returned nil")
	}

	// The result should be different from the original (blur should soften edges)
	resultRGBA := result.(*image.RGBA)

	// Check that the sharp edge has been softened
	// The middle area (around x=5) should have intermediate values after blur
	middleColor := resultRGBA.RGBAAt(5, 5)

	// After blur, the middle should not be pure black or pure white
	if (middleColor.R == 0 && middleColor.G == 0 && middleColor.B == 0) ||
		(middleColor.R == 255 && middleColor.G == 255 && middleColor.B == 255) {
		t.Error("Blur filter should create intermediate values at edges")
	}
}

func TestPresetFilters(t *testing.T) {
	img := createTestImage(10, 10, color.RGBA{128, 128, 128, 255})

	presets := []struct {
		name   string
		filter *FilterChain
	}{
		{"Instagram", Instagram()},
		{"Vintage", Vintage()},
		{"Dramatic", Dramatic()},
		{"BlackAndWhite", BlackAndWhite()},
		{"Warm", Warm()},
		{"Cool", Cool()},
	}

	for _, preset := range presets {
		t.Run(preset.name, func(t *testing.T) {
			result := preset.filter.Apply(img)
			if result == nil {
				t.Errorf("%s filter returned nil", preset.name)
			}

			// Check that the filter changed the image
			resultRGBA := result.(*image.RGBA)
			originalColor := img.RGBAAt(0, 0)
			resultColor := resultRGBA.RGBAAt(0, 0)

			if originalColor == resultColor {
				t.Errorf("%s filter should have changed the image", preset.name)
			}
		})
	}
}

func BenchmarkBrightnessFilter(b *testing.B) {
	img := createTestImage(100, 100, color.RGBA{128, 128, 128, 255})
	filter := BrightnessFilter{Amount: 1.5}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = filter.Apply(img)
	}
}

func BenchmarkFilterChain(b *testing.B) {
	img := createTestImage(100, 100, color.RGBA{128, 128, 128, 255})
	chain := NewFilterChain().
		Add(BrightnessFilter{Amount: 1.2}).
		Add(ContrastFilter{Amount: 1.1}).
		Add(SaturateFilter{Amount: 1.3})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = chain.Apply(img)
	}
}

func BenchmarkInstagramFilter(b *testing.B) {
	img := createTestImage(100, 100, color.RGBA{128, 128, 128, 255})
	filter := Instagram()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = filter.Apply(img)
	}
}
