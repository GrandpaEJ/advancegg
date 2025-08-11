package main

import (
	"fmt"
	"image/color"
	"math/rand"
	"runtime"
	"time"

	"github.com/GrandpaEJ/advancegg"
)

func main() {
	fmt.Println("Creating Performance Optimizations demo...")
	
	// Demonstrate basic performance optimizations
	demonstrateBasicOptimizations()
	
	// Demonstrate batch rendering
	demonstrateBatchRendering()
	
	// Demonstrate memory optimization
	demonstrateMemoryOptimization()
	
	// Demonstrate parallel processing
	demonstrateParallelProcessing()
	
	fmt.Println("Performance optimizations demo completed!")
}

func demonstrateBasicOptimizations() {
	fmt.Println("\n=== Basic Performance Optimizations ===")
	
	// Test 1: Efficient shape drawing
	fmt.Println("Testing efficient shape drawing...")
	
	start := time.Now()
	
	// Create context
	dc := advancegg.NewContext(1000, 800)
	dc.SetRGB(1, 1, 1)
	dc.Clear()
	
	// Draw many shapes efficiently
	colors := []color.RGBA{
		{255, 100, 100, 255},
		{100, 255, 100, 255},
		{100, 100, 255, 255},
		{255, 255, 100, 255},
		{255, 100, 255, 255},
		{100, 255, 255, 255},
	}
	
	// Batch similar operations
	for colorIndex, col := range colors {
		dc.SetColor(col)
		
		// Draw all circles of this color at once
		for i := 0; i < 100; i++ {
			x := float64(50 + (i%10)*90)
			y := float64(50 + (i/10)*70 + colorIndex*120)
			dc.DrawCircle(x, y, 20)
			dc.Fill()
		}
	}
	
	duration := time.Since(start)
	fmt.Printf("Drew 600 shapes in %v\n", duration)
	
	// Save result
	advancegg.SavePNG("performance-basic.png", dc.Image())
	fmt.Println("Created performance-basic.png")
}

func demonstrateBatchRendering() {
	fmt.Println("\n=== Batch Rendering ===")
	
	// Compare individual vs batch rendering
	
	// Test 1: Individual rendering
	start := time.Now()
	dc1 := advancegg.NewContext(800, 600)
	dc1.SetRGB(1, 1, 1)
	dc1.Clear()
	
	for i := 0; i < 1000; i++ {
		x := rand.Float64() * 800
		y := rand.Float64() * 600
		r := uint8(rand.Intn(256))
		g := uint8(rand.Intn(256))
		b := uint8(rand.Intn(256))
		
		dc1.SetRGB(float64(r)/255, float64(g)/255, float64(b)/255)
		dc1.DrawCircle(x, y, 5)
		dc1.Fill()
	}
	
	individualTime := time.Since(start)
	fmt.Printf("Individual rendering: %v\n", individualTime)
	
	// Test 2: Batch rendering (group by color)
	start = time.Now()
	dc2 := advancegg.NewContext(800, 600)
	dc2.SetRGB(1, 1, 1)
	dc2.Clear()
	
	// Pre-generate positions and colors
	type Circle struct {
		x, y float64
		r, g, b uint8
	}
	
	circles := make([]Circle, 1000)
	for i := range circles {
		circles[i] = Circle{
			x: rand.Float64() * 800,
			y: rand.Float64() * 600,
			r: uint8(rand.Intn(256)),
			g: uint8(rand.Intn(256)),
			b: uint8(rand.Intn(256)),
		}
	}
	
	// Group by similar colors (simplified batching)
	colorGroups := make(map[color.RGBA][]Circle)
	for _, circle := range circles {
		// Quantize colors to reduce groups
		quantizedColor := color.RGBA{
			circle.r & 0xF0, // Keep only upper 4 bits
			circle.g & 0xF0,
			circle.b & 0xF0,
			255,
		}
		colorGroups[quantizedColor] = append(colorGroups[quantizedColor], circle)
	}
	
	// Render by color groups
	for col, group := range colorGroups {
		dc2.SetColor(col)
		for _, circle := range group {
			dc2.DrawCircle(circle.x, circle.y, 5)
			dc2.Fill()
		}
	}
	
	batchTime := time.Since(start)
	fmt.Printf("Batch rendering: %v\n", batchTime)
	fmt.Printf("Speedup: %.2fx\n", float64(individualTime)/float64(batchTime))
	
	// Save results
	advancegg.SavePNG("performance-individual.png", dc1.Image())
	advancegg.SavePNG("performance-batch.png", dc2.Image())
	fmt.Println("Created performance comparison images")
}

func demonstrateMemoryOptimization() {
	fmt.Println("\n=== Memory Optimization ===")
	
	// Show memory usage before
	var m1 runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&m1)
	fmt.Printf("Memory before: %d KB\n", m1.Alloc/1024)
	
	// Create many contexts (memory intensive)
	contexts := make([]*advancegg.Context, 100)
	for i := range contexts {
		contexts[i] = advancegg.NewContext(200, 200)
		contexts[i].SetRGB(rand.Float64(), rand.Float64(), rand.Float64())
		contexts[i].Clear()
		
		// Draw something simple
		contexts[i].SetRGB(1, 1, 1)
		contexts[i].DrawCircle(100, 100, 50)
		contexts[i].Fill()
	}
	
	// Show memory usage after
	var m2 runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&m2)
	fmt.Printf("Memory after creating 100 contexts: %d KB\n", m2.Alloc/1024)
	fmt.Printf("Memory increase: %d KB\n", (m2.Alloc-m1.Alloc)/1024)
	
	// Demonstrate memory cleanup
	contexts = nil // Release references
	runtime.GC()
	
	var m3 runtime.MemStats
	runtime.ReadMemStats(&m3)
	fmt.Printf("Memory after cleanup: %d KB\n", m3.Alloc/1024)
	fmt.Printf("Memory freed: %d KB\n", (m2.Alloc-m3.Alloc)/1024)
}

func demonstrateParallelProcessing() {
	fmt.Println("\n=== Parallel Processing ===")
	
	// Test 1: Sequential processing
	start := time.Now()
	
	results1 := make([]*advancegg.Context, 8)
	for i := range results1 {
		results1[i] = createComplexImage(i)
	}
	
	sequentialTime := time.Since(start)
	fmt.Printf("Sequential processing: %v\n", sequentialTime)
	
	// Test 2: Parallel processing
	start = time.Now()
	
	results2 := make([]*advancegg.Context, 8)
	jobs := make(chan int, 8)
	results := make(chan struct {
		index int
		ctx   *advancegg.Context
	}, 8)
	
	// Start workers
	numWorkers := runtime.NumCPU()
	for w := 0; w < numWorkers; w++ {
		go func() {
			for index := range jobs {
				ctx := createComplexImage(index)
				results <- struct {
					index int
					ctx   *advancegg.Context
				}{index, ctx}
			}
		}()
	}
	
	// Submit jobs
	for i := 0; i < 8; i++ {
		jobs <- i
	}
	close(jobs)
	
	// Collect results
	for i := 0; i < 8; i++ {
		result := <-results
		results2[result.index] = result.ctx
	}
	
	parallelTime := time.Since(start)
	fmt.Printf("Parallel processing: %v\n", parallelTime)
	fmt.Printf("Speedup: %.2fx\n", float64(sequentialTime)/float64(parallelTime))
	
	// Save one result from each method
	advancegg.SavePNG("performance-sequential.png", results1[0].Image())
	advancegg.SavePNG("performance-parallel.png", results2[0].Image())
	fmt.Println("Created parallel processing comparison images")
}

func createComplexImage(seed int) *advancegg.Context {
	// Create a complex image for performance testing
	rand.Seed(int64(seed))
	
	dc := advancegg.NewContext(400, 400)
	dc.SetRGB(0.1, 0.1, 0.2)
	dc.Clear()
	
	// Draw many overlapping shapes
	for i := 0; i < 500; i++ {
		x := rand.Float64() * 400
		y := rand.Float64() * 400
		radius := rand.Float64()*20 + 5
		
		r := rand.Float64()
		g := rand.Float64()
		b := rand.Float64()
		a := rand.Float64()*0.5 + 0.3
		
		dc.SetRGBA(r, g, b, a)
		dc.DrawCircle(x, y, radius)
		dc.Fill()
	}
	
	// Add some complex paths
	for i := 0; i < 50; i++ {
		dc.SetRGBA(rand.Float64(), rand.Float64(), rand.Float64(), 0.7)
		dc.SetLineWidth(rand.Float64()*5 + 1)
		
		// Draw a random path
		dc.MoveTo(rand.Float64()*400, rand.Float64()*400)
		for j := 0; j < 10; j++ {
			dc.LineTo(rand.Float64()*400, rand.Float64()*400)
		}
		dc.Stroke()
	}
	
	return dc
}
