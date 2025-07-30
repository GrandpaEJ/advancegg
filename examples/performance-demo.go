package main

import (
	"fmt"
	"runtime"
	"time"
	
	"github.com/GrandpaEJ/advancegg"
	"github.com/GrandpaEJ/advancegg/internal/core"
)

func main() {
	fmt.Println("AdvanceGG Performance Demonstration")
	fmt.Println("===================================")
	
	// Enable debug mode for performance tracking
	core.SetDebugMode(core.DebugModePerformance)
	
	// Run performance tests
	runBasicPerformanceTest()
	runMemoryPerformanceTest()
	runBatchPerformanceTest()
	runCachePerformanceTest()
	runSIMDPerformanceTest()
	
	fmt.Println("\nPerformance demonstration completed!")
}

func runBasicPerformanceTest() {
	fmt.Println("\n1. Basic Drawing Performance Test")
	fmt.Println("---------------------------------")
	
	// Test different canvas sizes
	sizes := []struct {
		name   string
		width  int
		height int
	}{
		{"Small", 400, 300},
		{"Medium", 800, 600},
		{"Large", 1920, 1080},
	}
	
	for _, size := range sizes {
		start := time.Now()
		
		dc := advancegg.NewContext(size.width, size.height)
		
		// Draw many shapes
		numShapes := 1000
		for i := 0; i < numShapes; i++ {
			dc.SetRGB(float64(i%255)/255, 0.5, 0.8)
			dc.DrawCircle(float64(i%(size.width-20))+10, 
				float64((i*7)%(size.height-20))+10, 5)
			dc.Fill()
		}
		
		duration := time.Since(start)
		fmt.Printf("  %s (%dx%d): %d shapes in %v (%.2f shapes/ms)\n", 
			size.name, size.width, size.height, numShapes, duration, 
			float64(numShapes)/float64(duration.Nanoseconds()/1000000))
	}
}

func runMemoryPerformanceTest() {
	fmt.Println("\n2. Memory Performance Test")
	fmt.Println("--------------------------")
	
	// Test with and without memory pooling
	fmt.Println("  Testing without memory pooling:")
	testMemoryUsage(false)
	
	fmt.Println("  Testing with memory pooling:")
	testMemoryUsage(true)
	
	// Show memory pool statistics
	stats := core.GetMemoryStats()
	fmt.Printf("  Memory Pool Stats: Image pools: %d, Byte slice pools: %d\n", 
		stats.ImagePoolSize, stats.ByteSlicePoolSize)
}

func testMemoryUsage(usePooling bool) {
	runtime.GC() // Force garbage collection
	var m1, m2 runtime.MemStats
	runtime.ReadMemStats(&m1)
	
	start := time.Now()
	
	for i := 0; i < 100; i++ {
		var dc *advancegg.Context
		
		if usePooling {
			// Use pooled context (if available)
			dc = advancegg.NewContext(400, 300)
		} else {
			dc = advancegg.NewContext(400, 300)
		}
		
		// Perform drawing operations
		for j := 0; j < 50; j++ {
			dc.SetRGB(float64(j)/50, 0.5, 0.8)
			dc.DrawCircle(float64(j*8), 150, 10)
			dc.Fill()
		}
		
		// Get image data (causes allocation)
		imageData := dc.GetImageData()
		_ = imageData.Clone()
	}
	
	duration := time.Since(start)
	runtime.ReadMemStats(&m2)
	
	fmt.Printf("    Time: %v, Memory allocated: %d KB\n", 
		duration, (m2.TotalAlloc-m1.TotalAlloc)/1024)
}

func runBatchPerformanceTest() {
	fmt.Println("\n3. Batch Operations Performance Test")
	fmt.Println("------------------------------------")
	
	dc := advancegg.NewContext(800, 600)
	numShapes := 1000
	
	// Test individual operations
	start := time.Now()
	for i := 0; i < numShapes; i++ {
		dc.SetRGBA255(i%255, 100, 200, 255)
		dc.DrawCircle(float64(i%780)+10, float64((i*3)%580)+10, 5)
		dc.Fill()
	}
	individualTime := time.Since(start)
	
	// Clear canvas
	dc.SetRGB(1, 1, 1)
	dc.Clear()
	
	// Test batch operations
	circles := make([]core.BatchCircle, numShapes)
	for i := range circles {
		circles[i] = core.BatchCircle{
			X: float64(i%780) + 10, Y: float64((i*3)%580) + 10, Radius: 5,
			Color: advancegg.NewColorFromRGBA255(uint8(i%255), 100, 200, 255).ToStandardColor(),
			Fill: true,
		}
	}
	
	start = time.Now()
	dc.BatchCircles(circles)
	batchTime := time.Since(start)
	
	fmt.Printf("  Individual operations: %v\n", individualTime)
	fmt.Printf("  Batch operations: %v\n", batchTime)
	fmt.Printf("  Speedup: %.2fx\n", float64(individualTime.Nanoseconds())/float64(batchTime.Nanoseconds()))
}

func runCachePerformanceTest() {
	fmt.Println("\n4. Cache Performance Test")
	fmt.Println("-------------------------")
	
	dc := advancegg.NewContext(400, 300)
	
	// Test without caching
	start := time.Now()
	for i := 0; i < 100; i++ {
		dc.SetRGB(0.8, 0.6, 0.4)
		dc.DrawCircle(200, 150, 50)
		dc.Fill()
	}
	noCacheTime := time.Since(start)
	
	// Test with caching (if available)
	start = time.Now()
	for i := 0; i < 100; i++ {
		dc.DrawCachedCircle(200, 150, 50, true)
	}
	cacheTime := time.Since(start)
	
	fmt.Printf("  Without caching: %v\n", noCacheTime)
	fmt.Printf("  With caching: %v\n", cacheTime)
	if cacheTime.Nanoseconds() > 0 {
		fmt.Printf("  Speedup: %.2fx\n", float64(noCacheTime.Nanoseconds())/float64(cacheTime.Nanoseconds()))
	}
	
	// Show cache statistics
	cacheStats := core.GetCacheStats()
	for name, stats := range cacheStats {
		fmt.Printf("  %s cache: %d hits, %d misses, %.2f%% hit rate\n", 
			name, stats.Hits, stats.Misses, 
			float64(stats.Hits)/float64(stats.Hits+stats.Misses)*100)
	}
}

func runSIMDPerformanceTest() {
	fmt.Println("\n5. SIMD Performance Test")
	fmt.Println("------------------------")
	
	// Create test image
	dc := advancegg.NewContext(800, 600)
	dc.SetRGB(0.5, 0.7, 0.9)
	dc.Clear()
	
	// Add some content
	for i := 0; i < 50; i++ {
		dc.SetRGB(float64(i)/50, 0.5, 0.8)
		dc.DrawCircle(float64(i*15), 300, 20)
		dc.Fill()
	}
	
	img := dc.Image().(*advancegg.ImageRGBA)
	
	// Test regular blur
	start := time.Now()
	regularBlur := advancegg.Blur(5).Apply(img)
	regularTime := time.Since(start)
	
	// Test SIMD blur (if available)
	start = time.Now()
	simdBlur := core.SIMDBlur(img, 5)
	simdTime := time.Since(start)
	
	fmt.Printf("  Regular blur: %v\n", regularTime)
	fmt.Printf("  SIMD blur: %v\n", simdTime)
	if simdTime.Nanoseconds() > 0 {
		fmt.Printf("  Speedup: %.2fx\n", float64(regularTime.Nanoseconds())/float64(simdTime.Nanoseconds()))
	}
	
	// Test color transformations
	transform := func(r, g, b, a uint8) (uint8, uint8, uint8, uint8) {
		// Simple brightness adjustment
		factor := 1.2
		return uint8(float64(r)*factor), uint8(float64(g)*factor), uint8(float64(b)*factor), a
	}
	
	start = time.Now()
	simdTransform := core.SIMDColorTransform(img, transform)
	transformTime := time.Since(start)
	
	fmt.Printf("  SIMD color transform: %v\n", transformTime)
	
	// Save results for comparison
	resultDC := advancegg.NewContextForRGBA(regularBlur)
	resultDC.SavePNG("performance-regular-blur.png")
	
	resultDC2 := advancegg.NewContextForRGBA(simdBlur)
	resultDC2.SavePNG("performance-simd-blur.png")
	
	resultDC3 := advancegg.NewContextForRGBA(simdTransform)
	resultDC3.SavePNG("performance-simd-transform.png")
	
	fmt.Println("  Performance test images saved")
}

// Helper function to create context from RGBA image
func NewContextForRGBA(img *advancegg.ImageRGBA) *advancegg.Context {
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
