package benchmarks

import (
	"image/color"
	"testing"

	"github.com/GrandpaEJ/advancegg"
)

// Comprehensive benchmarking suite for AdvanceGG

// Basic drawing operations benchmarks

func BenchmarkNewContext(b *testing.B) {
	for i := 0; i < b.N; i++ {
		dc := advancegg.NewContext(800, 600)
		_ = dc
	}
}

func BenchmarkDrawCircle(b *testing.B) {
	dc := advancegg.NewContext(800, 600)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		dc.DrawCircle(400, 300, 50)
		dc.Fill()
	}
}

func BenchmarkDrawRectangle(b *testing.B) {
	dc := advancegg.NewContext(800, 600)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		dc.DrawRectangle(100, 100, 200, 150)
		dc.Fill()
	}
}

func BenchmarkDrawLine(b *testing.B) {
	dc := advancegg.NewContext(800, 600)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		dc.DrawLine(0, 0, 800, 600)
		dc.Stroke()
	}
}

func BenchmarkDrawText(b *testing.B) {
	dc := advancegg.NewContext(800, 600)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		dc.DrawString("Hello, World!", 100, 100)
	}
}

// Complex drawing benchmarks

func BenchmarkComplexScene(b *testing.B) {
	for i := 0; i < b.N; i++ {
		dc := advancegg.NewContext(800, 600)

		// Background
		dc.SetRGB(0.2, 0.3, 0.8)
		dc.Clear()

		// Draw multiple shapes
		for j := 0; j < 50; j++ {
			dc.SetRGB(float64(j%255)/255, 0.5, 0.8)
			dc.DrawCircle(float64(j*15), float64(j*10), 20)
			dc.Fill()
		}

		// Draw rectangles
		for j := 0; j < 30; j++ {
			dc.SetRGB(0.8, float64(j%255)/255, 0.3)
			dc.DrawRectangle(float64(j*25), float64(j*15), 40, 30)
			dc.Fill()
		}

		// Draw text
		dc.SetRGB(1, 1, 1)
		dc.DrawString("Complex Scene Benchmark", 50, 50)
	}
}

// Image processing benchmarks

func BenchmarkImageFilters(b *testing.B) {
	dc := advancegg.NewContext(400, 300)

	// Create test image
	dc.SetRGB(0.5, 0.7, 0.9)
	dc.Clear()
	for i := 0; i < 20; i++ {
		dc.SetRGB(float64(i)/20, 0.5, 0.8)
		dc.DrawCircle(float64(i*20), 150, 15)
		dc.Fill()
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// Apply various filters
		dc.ApplyFilter(advancegg.Blur(3))
		dc.ApplyFilter(advancegg.Brightness(1.2))
		dc.ApplyFilter(advancegg.Contrast(1.1))
	}
}

func BenchmarkImageResize(b *testing.B) {
	dc := advancegg.NewContext(800, 600)
	dc.SetRGB(0.5, 0.7, 0.9)
	dc.Clear()

	imageData := dc.GetImageData()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		resized := imageData.Resize(400, 300)
		_ = resized
	}
}

// Color space benchmarks

func BenchmarkColorSpaceConversions(b *testing.B) {
	color := advancegg.NewColor(0.8, 0.6, 0.4, 1.0)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		cmyk := color.ToCMYK()
		hsv := color.ToHSV()
		lab := color.ToLAB()

		// Convert back
		_ = cmyk.ToRGB()
		_ = hsv.ToRGB()
		_ = lab.ToRGB()
	}
}

// Memory allocation benchmarks

func BenchmarkMemoryAllocation(b *testing.B) {
	for i := 0; i < b.N; i++ {
		dc := advancegg.NewContext(1920, 1080)

		// Perform operations that allocate memory
		for j := 0; j < 100; j++ {
			dc.DrawCircle(float64(j*10), float64(j*10), 5)
			dc.Fill()
		}

		// Force some allocations
		imageData := dc.GetImageData()
		_ = imageData.Clone()
	}
}

// Batch operations benchmarks

func BenchmarkBatchOperations(b *testing.B) {
	dc := advancegg.NewContext(800, 600)

	// Prepare batch data
	circles := make([]advancegg.BatchCircle, 100)
	for i := range circles {
		circles[i] = advancegg.BatchCircle{
			X: float64(i * 8), Y: float64(i * 6), Radius: 10,
			Color: color.RGBA{uint8(i), 100, 200, 255}, Fill: true,
		}
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		dc.BatchCircles(circles)
	}
}

func BenchmarkIndividualOperations(b *testing.B) {
	dc := advancegg.NewContext(800, 600)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for j := 0; j < 100; j++ {
			dc.SetRGBA255(j, 100, 200, 255)
			dc.DrawCircle(float64(j*8), float64(j*6), 10)
			dc.Fill()
		}
	}
}

// Font rendering benchmarks

func BenchmarkFontRendering(b *testing.B) {
	dc := advancegg.NewContext(800, 600)

	// Try to load a font (will use default if not available)
	// dc.LoadFontFace("arial.ttf", 16)

	text := "The quick brown fox jumps over the lazy dog"
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for j := 0; j < 20; j++ {
			dc.DrawString(text, 10, float64(j*25+20))
		}
	}
}

// Path operations benchmarks

func BenchmarkPathOperations(b *testing.B) {
	dc := advancegg.NewContext(800, 600)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		path := advancegg.NewPath2D()

		// Create complex path
		path.MoveTo(100, 100)
		for j := 0; j < 50; j++ {
			path.LineTo(float64(100+j*10), float64(100+j*5))
			path.QuadraticCurveTo(float64(150+j*10), float64(120+j*5), float64(200+j*10), float64(100+j*5))
		}
		path.ClosePath()

		dc.DrawPath2D(path)
		dc.Fill()
	}
}

// Parallel processing benchmarks

func BenchmarkParallelDrawing(b *testing.B) {
	for i := 0; i < b.N; i++ {
		dc := advancegg.NewContext(1920, 1080)

		// Simulate parallel drawing operations
		// In a real scenario, this would use goroutines
		for j := 0; j < 1000; j++ {
			dc.SetRGB(float64(j%255)/255, 0.5, 0.8)
			dc.DrawCircle(float64(j%1920), float64((j/1920)*1080), 5)
			dc.Fill()
		}
	}
}

// Image format benchmarks

func BenchmarkImageSaving(b *testing.B) {
	dc := advancegg.NewContext(800, 600)

	// Create test image
	dc.SetRGB(0.5, 0.7, 0.9)
	dc.Clear()
	for i := 0; i < 50; i++ {
		dc.SetRGB(float64(i)/50, 0.5, 0.8)
		dc.DrawCircle(float64(i*15), 300, 20)
		dc.Fill()
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// Benchmark different formats
		dc.SavePNG("benchmark_test.png")
		dc.SaveJPEG("benchmark_test.jpg", 90)
	}
}

// Stress tests

func BenchmarkStressTest(b *testing.B) {
	for i := 0; i < b.N; i++ {
		dc := advancegg.NewContext(1920, 1080)

		// Stress test with many operations
		for j := 0; j < 10000; j++ {
			switch j % 4 {
			case 0:
				dc.DrawCircle(float64(j%1920), float64((j/1920)*1080), 2)
				dc.Fill()
			case 1:
				dc.DrawRectangle(float64(j%1920), float64((j/1920)*1080), 4, 4)
				dc.Fill()
			case 2:
				dc.DrawLine(float64(j%1920), float64((j/1920)*1080),
					float64((j+10)%1920), float64(((j+10)/1920)*1080))
				dc.Stroke()
			case 3:
				dc.SetPixel(j%1920, (j/1920)*1080)
			}
		}
	}
}

// Benchmark comparison helpers

func BenchmarkComparison(b *testing.B) {
	sizes := []struct {
		name   string
		width  int
		height int
	}{
		{"Small", 400, 300},
		{"Medium", 800, 600},
		{"Large", 1920, 1080},
		{"XLarge", 3840, 2160},
	}

	for _, size := range sizes {
		b.Run(size.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				dc := advancegg.NewContext(size.width, size.height)

				// Standard drawing operations
				dc.SetRGB(0.2, 0.3, 0.8)
				dc.Clear()

				numOps := (size.width * size.height) / 10000
				for j := 0; j < numOps; j++ {
					dc.SetRGB(float64(j%255)/255, 0.5, 0.8)
					dc.DrawCircle(float64(j%(size.width-20))+10,
						float64((j*7)%(size.height-20))+10, 5)
					dc.Fill()
				}
			}
		})
	}
}
