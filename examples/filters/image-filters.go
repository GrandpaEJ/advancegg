package main

import (
	"fmt"
	"image"

	"github.com/GrandpaEJ/advancegg"
)

func main() {
	// Create a sample image to apply filters to
	dc := advancegg.NewContext(800, 600)

	// Create a colorful test image
	createTestImage(dc)

	// Save original
	dc.SavePNG("images/filters/filters-original.png")
	fmt.Println("Original image saved as filters-original.png")

	// Apply various filters and save each one
	filters := []struct {
		name   string
		filter func(image.Image) image.Image
	}{
		{"grayscale", advancegg.Grayscale},
		{"invert", advancegg.Invert},
		{"sepia", advancegg.Sepia},
		{"brightness", advancegg.Brightness(1.5)},
		{"contrast", advancegg.Contrast(2.0)},
		{"blur", advancegg.Blur(3)},
		{"sharpen", advancegg.Sharpen},
		{"threshold", advancegg.Threshold(128)},
		{"pixelate", advancegg.Pixelate(8)},
		{"noise", advancegg.Noise(0.1)},
		{"edge-detection", advancegg.EdgeDetection},
		{"emboss", advancegg.Emboss},
		{"posterize", advancegg.Posterize(4)},
		{"vignette", advancegg.Vignette(0.8)},
	}

	for _, f := range filters {
		// Create a new context with the original image
		filterDC := advancegg.NewContext(800, 600)
		createTestImage(filterDC)

		// Apply the filter
		filterDC.ApplyFilter(f.filter)

		// Save the filtered image
		filename := fmt.Sprintf("filters-%s.png", f.name)
		filterDC.SavePNG(filename)
		fmt.Printf("Applied %s filter, saved as %s\n", f.name, filename)
	}

	// Create a composite image showing multiple filters
	createFilterComposite()

	fmt.Println("All filter examples completed!")
}

func createTestImage(dc *advancegg.Context) {
	// White background
	dc.SetRGB(1, 1, 1)
	dc.Clear()

	// Draw colorful shapes for testing filters

	// Red circle
	dc.SetRGB(1, 0, 0)
	dc.DrawCircle(200, 150, 80)
	dc.Fill()

	// Green rectangle
	dc.SetRGB(0, 1, 0)
	dc.DrawRectangle(350, 100, 150, 100)
	dc.Fill()

	// Blue triangle
	dc.SetRGB(0, 0, 1)
	dc.MoveTo(600, 100)
	dc.LineTo(550, 200)
	dc.LineTo(650, 200)
	dc.ClosePath()
	dc.Fill()

	// Yellow ellipse
	dc.SetRGB(1, 1, 0)
	dc.DrawEllipse(150, 350, 100, 60)
	dc.Fill()

	// Magenta rounded rectangle
	dc.SetRGB(1, 0, 1)
	dc.DrawRoundedRectangle(300, 300, 200, 100, 20)
	dc.Fill()

	// Cyan lines
	dc.SetRGB(0, 1, 1)
	dc.SetLineWidth(5)
	for i := 0; i < 10; i++ {
		dc.DrawLine(500+float64(i)*20, 350, 500+float64(i)*20, 450)
		dc.Stroke()
	}

	// Add some text
	dc.SetRGB(0, 0, 0)
	dc.DrawString("Filter Test Image", 50, 500)
	dc.DrawString("Various shapes and colors for testing", 50, 530)
}

func createFilterComposite() {
	// Create a composite image showing multiple filters side by side
	compositeDC := advancegg.NewContext(1600, 1200)
	compositeDC.SetRGB(0.9, 0.9, 0.9)
	compositeDC.Clear()

	// Create small test images for each filter
	filters := []struct {
		name   string
		filter func(image.Image) image.Image
		x, y   float64
	}{
		{"Original", nil, 50, 50},
		{"Grayscale", advancegg.Grayscale, 450, 50},
		{"Sepia", advancegg.Sepia, 850, 50},
		{"Invert", advancegg.Invert, 1250, 50},
		{"Blur", advancegg.Blur(2), 50, 350},
		{"Sharpen", advancegg.Sharpen, 450, 350},
		{"Edge Detection", advancegg.EdgeDetection, 850, 350},
		{"Emboss", advancegg.Emboss, 1250, 350},
		{"Pixelate", advancegg.Pixelate(6), 50, 650},
		{"Posterize", advancegg.Posterize(4), 450, 650},
		{"Vignette", advancegg.Vignette(0.6), 850, 650},
		{"Threshold", advancegg.Threshold(128), 1250, 650},
	}

	for _, f := range filters {
		// Create a small test image
		smallDC := advancegg.NewContext(300, 200)
		createSmallTestImage(smallDC)

		// Apply filter if not original
		if f.filter != nil {
			smallDC.ApplyFilter(f.filter)
		}

		// Draw the filtered image onto the composite
		compositeDC.DrawImageAnchored(smallDC.Image(), int(f.x+150), int(f.y+100), 0.5, 0.5)

		// Add label
		compositeDC.SetRGB(0, 0, 0)
		compositeDC.DrawStringAnchored(f.name, f.x+150, f.y+220, 0.5, 0.5)
	}

	compositeDC.SavePNG("images/filters/filters-composite.png")
	fmt.Println("Filter composite saved as filters-composite.png")
}

func createSmallTestImage(dc *advancegg.Context) {
	// Create a smaller version of the test image
	dc.SetRGB(1, 1, 1)
	dc.Clear()

	// Colorful gradient background
	for y := 0; y < 200; y++ {
		for x := 0; x < 300; x++ {
			r := float64(x) / 300.0
			g := float64(y) / 200.0
			b := 0.5
			dc.SetRGB(r, g, b)
			dc.SetPixel(x, y)
		}
	}

	// Add some shapes
	dc.SetRGB(1, 1, 1)
	dc.DrawCircle(75, 50, 25)
	dc.Fill()

	dc.SetRGB(0, 0, 0)
	dc.DrawRectangle(150, 25, 50, 50)
	dc.Fill()

	dc.SetRGB(1, 0, 0)
	dc.DrawCircle(225, 50, 20)
	dc.Fill()
}
