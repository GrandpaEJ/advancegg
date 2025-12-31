package main

import (
	"fmt"
	"image"
	"time"

	"github.com/GrandpaEJ/advancegg"
	"github.com/GrandpaEJ/advancegg/internal/core"
)

func main() {
	fmt.Println("Testing filter performance optimizations...")

	// Create test image
	dc := advancegg.NewContext(800, 600)
	createTestImage(dc)

	// Test standard vs optimized filters
	testFilterPerformance(dc)

	// Test batch filtering
	testBatchFiltering(dc)

	// Test parallel filtering
	testParallelFiltering(dc)

	fmt.Println("Filter performance test completed!")
}

func createTestImage(dc *advancegg.Context) {
	// Create a complex test image with various elements
	dc.SetRGB(0.2, 0.3, 0.8)
	dc.Clear()

	// Add some geometric shapes
	for i := 0; i < 20; i++ {
		dc.SetRGB(float64(i)/20, 0.5, 0.8)
		dc.DrawCircle(float64(i*40), 150, 20)
		dc.Fill()
	}

	// Add some rectangles
	for i := 0; i < 15; i++ {
		dc.SetRGB(0.8, float64(i)/15, 0.3)
		dc.DrawRectangle(float64(i*50), 300, 40, 30)
		dc.Fill()
	}

	// Add some lines
	dc.SetRGB(1, 1, 1)
	dc.SetLineWidth(2)
	for i := 0; i < 10; i++ {
		dc.MoveTo(0, float64(i*60))
		dc.LineTo(800, float64(i*60+30))
		dc.Stroke()
	}
}

func testFilterPerformance(dc *advancegg.Context) {
	fmt.Println("\n=== Filter Performance Comparison ===")

	// Get the current image
	img := dc.Image()

	// Test Grayscale
	fmt.Println("Testing Grayscale filters...")
	start := time.Now()
	_ = core.Grayscale(img)
	standardTime := time.Since(start)

	start = time.Now()
	_ = core.FastGrayscale(img)
	optimizedTime := time.Since(start)

	fmt.Printf("Grayscale - Standard: %v, Optimized: %v (%.2fx faster)\n",
		standardTime, optimizedTime, float64(standardTime)/float64(optimizedTime))

	// Test Brightness
	fmt.Println("Testing Brightness filters...")
	start = time.Now()
	_ = core.Brightness(1.5)(img)
	standardTime = time.Since(start)

	start = time.Now()
	_ = core.FastBrightness(1.5)(img)
	optimizedTime = time.Since(start)

	fmt.Printf("Brightness - Standard: %v, Optimized: %v (%.2fx faster)\n",
		standardTime, optimizedTime, float64(standardTime)/float64(optimizedTime))

	// Test Contrast
	fmt.Println("Testing Contrast filters...")
	start = time.Now()
	_ = core.Contrast(1.3)(img)
	standardTime = time.Since(start)

	start = time.Now()
	_ = core.FastContrast(1.3)(img)
	optimizedTime = time.Since(start)

	fmt.Printf("Contrast - Standard: %v, Optimized: %v (%.2fx faster)\n",
		standardTime, optimizedTime, float64(standardTime)/float64(optimizedTime))

	// Test Blur
	fmt.Println("Testing Blur filters...")
	start = time.Now()
	_ = core.Blur(3)(img)
	standardTime = time.Since(start)

	start = time.Now()
	_ = core.FastBlur(3)(img)
	optimizedTime = time.Since(start)

	fmt.Printf("Blur - Standard: %v, Optimized: %v (%.2fx faster)\n",
		standardTime, optimizedTime, float64(standardTime)/float64(optimizedTime))

	// Test Edge Detection
	fmt.Println("Testing Edge Detection filters...")
	start = time.Now()
	_ = core.EdgeDetection(img)
	standardTime = time.Since(start)

	start = time.Now()
	_ = core.FastEdgeDetection()(img)
	optimizedTime = time.Since(start)

	fmt.Printf("Edge Detection - Standard: %v, Optimized: %v (%.2fx faster)\n",
		standardTime, optimizedTime, float64(standardTime)/float64(optimizedTime))
}

func testBatchFiltering(dc *advancegg.Context) {
	fmt.Println("\n=== Batch Filtering Test ===")

	img := dc.Image()

	// Test individual filters
	fmt.Println("Testing individual filter application...")
	start := time.Now()
	result := img
	result = core.FastGrayscale(result)
	result = core.FastBrightness(1.2)(result)
	result = core.FastContrast(1.1)(result)
	result = core.FastBlur(2)(result)
	individualTime := time.Since(start)

	// Test batch filter
	fmt.Println("Testing batch filter application...")
	start = time.Now()
	batchFilter := core.BatchFilter(
		core.FastGrayscale,
		core.FastBrightness(1.2),
		core.FastContrast(1.1),
		core.FastBlur(2),
	)
	_ = batchFilter(img)
	batchTime := time.Since(start)

	fmt.Printf("Individual: %v, Batch: %v (%.2fx faster)\n",
		individualTime, batchTime, float64(individualTime)/float64(batchTime))

	// Save batch result
	filteredImg := batchFilter(img)
	imageData := core.NewImageDataFromImage(filteredImg)
	dc.PutImageData(imageData)
	dc.SavePNG("images/filters/filter-batch-result.png")
	fmt.Println("Saved batch filter result as filter-batch-result.png")
}

func testParallelFiltering(dc *advancegg.Context) {
	fmt.Println("\n=== Parallel Filtering Test ===")

	img := dc.Image()

	// Test standard filter
	fmt.Println("Testing standard filter...")
	start := time.Now()
	_ = core.FastBlur(5)(img)
	standardTime := time.Since(start)

	// Test parallel filter with different worker counts
	workerCounts := []int{1, 2, 4, 8}

	for _, workers := range workerCounts {
		fmt.Printf("Testing parallel filter with %d workers...\n", workers)
		start = time.Now()
		parallelFilter := core.ParallelFilter(core.FastBlur(5), workers)
		_ = parallelFilter(img)
		parallelTime := time.Since(start)

		fmt.Printf("Workers: %d, Time: %v (%.2fx vs standard)\n",
			workers, parallelTime, float64(standardTime)/float64(parallelTime))
	}

	// Test complex parallel filter chain
	fmt.Println("Testing complex parallel filter chain...")
	start = time.Now()
	complexFilter := core.ParallelFilter(
		core.BatchFilter(
			core.FastGrayscale,
			core.FastContrast(1.5),
			core.FastEdgeDetection(),
		), 4)
	result := complexFilter(img)
	complexTime := time.Since(start)

	fmt.Printf("Complex parallel filter chain: %v\n", complexTime)

	// Save parallel result
	imageData := core.NewImageDataFromImage(result)
	dc.PutImageData(imageData)
	dc.SavePNG("images/filters/filter-parallel-result.png")
	fmt.Println("Saved parallel filter result as filter-parallel-result.png")
}

func benchmarkFilter(name string, filter core.Filter, img image.Image, iterations int) {
	fmt.Printf("Benchmarking %s (%d iterations)...\n", name, iterations)

	start := time.Now()
	for i := 0; i < iterations; i++ {
		_ = filter(img)
	}
	totalTime := time.Since(start)
	avgTime := totalTime / time.Duration(iterations)

	fmt.Printf("%s - Total: %v, Average: %v\n", name, totalTime, avgTime)
}

func demonstrateFilterEffects(dc *advancegg.Context) {
	fmt.Println("\n=== Filter Effects Demonstration ===")

	// Create a grid of filter effects
	originalImg := dc.Image()

	// Create a larger canvas for the grid
	gridDC := advancegg.NewContext(1600, 1200)
	gridDC.SetRGB(0.1, 0.1, 0.1)
	gridDC.Clear()

	filters := []struct {
		name   string
		filter core.Filter
		x, y   float64
	}{
		{"Original", func(img image.Image) image.Image { return img }, 0, 0},
		{"Grayscale", core.FastGrayscale, 400, 0},
		{"Bright", core.FastBrightness(1.5), 800, 0},
		{"Contrast", core.FastContrast(1.5), 1200, 0},
		{"Blur", core.FastBlur(3), 0, 300},
		{"Sharpen", core.FastSharpen(1.0), 400, 300},
		{"Edges", core.FastEdgeDetection(), 800, 300},
		{"Threshold", core.Threshold(128), 1200, 300},
	}

	for _, f := range filters {
		// Apply filter and save individual result
		filtered := f.filter(originalImg)

		// Create individual context for this filter
		filterDC := advancegg.NewContext(400, 300)
		imageData := core.NewImageDataFromImage(filtered)
		filterDC.PutImageData(imageData)

		// Save individual filter result
		filename := fmt.Sprintf("filter-%s.png", f.name)
		filterDC.SavePNG(filename)

		// Add label to grid
		gridDC.SetRGB(1, 1, 1)
		gridDC.DrawString(f.name, f.x+10, f.y+290)
	}

	gridDC.SavePNG("images/filters/filter-effects-grid.png")
	fmt.Println("Saved filter effects grid as filter-effects-grid.png")
	fmt.Println("Saved individual filter results as filter-*.png")
}
