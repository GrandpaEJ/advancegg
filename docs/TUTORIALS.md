# AdvanceGG Tutorials

This document contains step-by-step tutorials for common use cases with AdvanceGG.

## Tutorial 1: Creating a Simple Chart

Learn how to create a basic bar chart from scratch.

### Step 1: Set up the canvas
```go
dc := advancegg.NewContext(800, 600)
dc.SetRGB(1, 1, 1) // White background
dc.Clear()
```

### Step 2: Prepare your data
```go
data := []struct {
    label string
    value float64
    color [3]float64
}{
    {"Jan", 120, [3]float64{0.8, 0.2, 0.2}},
    {"Feb", 150, [3]float64{0.2, 0.8, 0.2}},
    {"Mar", 180, [3]float64{0.2, 0.2, 0.8}},
}
```

### Step 3: Draw the axes
```go
chartX, chartY := 100.0, 100.0
chartWidth, chartHeight := 600.0, 400.0

dc.SetRGB(0, 0, 0)
dc.SetLineWidth(2)
dc.DrawLine(chartX, chartY+chartHeight, chartX+chartWidth, chartY+chartHeight) // X-axis
dc.DrawLine(chartX, chartY, chartX, chartY+chartHeight) // Y-axis
dc.Stroke()
```

### Step 4: Draw the bars
```go
maxValue := 200.0 // Find your actual max
barWidth := chartWidth / float64(len(data)) * 0.8

for i, item := range data {
    barHeight := (item.value / maxValue) * chartHeight
    barX := chartX + float64(i)*(chartWidth/float64(len(data))) + 10
    barY := chartY + chartHeight - barHeight
    
    dc.SetRGB(item.color[0], item.color[1], item.color[2])
    dc.DrawRectangle(barX, barY, barWidth, barHeight)
    dc.Fill()
}
```

## Tutorial 2: Image Processing Pipeline

Learn how to apply multiple filters to an image.

### Step 1: Load an image
```go
dc := advancegg.NewContext(800, 600)
// Create or load your image content here
```

### Step 2: Apply filters in sequence
```go
// Method 1: Individual filters
dc.ApplyFilter(advancegg.Blur(3))
dc.ApplyFilter(advancegg.Brightness(1.2))
dc.ApplyFilter(advancegg.Contrast(1.1))

// Method 2: Filter chain (more efficient)
chain := advancegg.NewFilterChain().
    Add(advancegg.BlurFilter{Radius: 3}).
    Add(advancegg.BrightnessFilter{Amount: 1.2}).
    Add(advancegg.ContrastFilter{Amount: 1.1})

processedImage := chain.Apply(dc.Image())
```

### Step 3: Save the result
```go
dc.SavePNG("processed-image.png")
```

## Tutorial 3: Creating Animations

Learn how to create frame-by-frame animations.

### Step 1: Set up animation parameters
```go
frameCount := 60
frameWidth, frameHeight := 400, 400
```

### Step 2: Generate frames
```go
for frame := 0; frame < frameCount; frame++ {
    dc := advancegg.NewContext(frameWidth, frameHeight)
    
    // Clear background
    dc.SetRGB(0.1, 0.1, 0.2)
    dc.Clear()
    
    // Animate a rotating circle
    angle := float64(frame) * 2 * math.Pi / float64(frameCount)
    x := 200 + 100*math.Cos(angle)
    y := 200 + 100*math.Sin(angle)
    
    dc.SetRGB(1, 0.5, 0)
    dc.DrawCircle(x, y, 20)
    dc.Fill()
    
    // Save frame
    filename := fmt.Sprintf("frame_%03d.png", frame)
    dc.SavePNG(filename)
}
```

### Step 3: Combine frames (external tool needed)
```bash
# Use ffmpeg to create video
ffmpeg -r 30 -i frame_%03d.png -c:v libx264 animation.mp4

# Or create GIF
ffmpeg -r 10 -i frame_%03d.png animation.gif
```

## Tutorial 4: Custom Patterns and Gradients

Learn how to create complex fill patterns.

### Step 1: Create a linear gradient
```go
gradient := advancegg.CreateLinearGradient(0, 0, 200, 0,
    color.RGBA{255, 0, 0, 255},   // Red
    color.RGBA{255, 255, 0, 255}, // Yellow
    color.RGBA{0, 255, 0, 255},   // Green
)
```

### Step 2: Create a radial gradient
```go
radialGradient := advancegg.CreateRadialGradient(100, 100, 80,
    color.RGBA{255, 255, 255, 255}, // White center
    color.RGBA{0, 0, 255, 255},     // Blue edge
)
```

### Step 3: Create custom patterns
```go
// Checkerboard pattern
checkerboard := advancegg.CreateCheckerboard(20)

// Polka dots
polkaDots := advancegg.CreatePolkaDots(40, 15)

// Custom noise pattern
noise := advancegg.NoisePattern{
    Scale:     0.1,
    BaseColor: color.RGBA{100, 150, 200, 255},
    Intensity: 0.5,
}
```

### Step 4: Apply patterns
```go
dc.SetFillPattern(gradient)
dc.DrawRectangle(0, 0, 200, 100)
dc.Fill()

// Or fill entire canvas with pattern
advancegg.PatternFill(dc.Image().(*image.RGBA), checkerboard)
```

## Tutorial 5: Performance Optimization

Learn how to optimize your graphics code for better performance.

### Step 1: Use batch operations
```go
// Instead of individual draws
for i := 0; i < 1000; i++ {
    dc.SetRGB(float64(i)/1000, 0.5, 0.8)
    dc.DrawCircle(float64(i%800), float64(i/800)*600, 5)
    dc.Fill()
}

// Use batch operations
circles := make([]advancegg.BatchCircle, 1000)
for i := range circles {
    circles[i] = advancegg.BatchCircle{
        X: float64(i%800), Y: float64(i/800)*600, Radius: 5,
        Color: color.RGBA{uint8(i/4), 128, 200, 255}, Fill: true,
    }
}
dc.BatchCircles(circles)
```

### Step 2: Enable caching
```go
// For repeated elements
dc.DrawCachedCircle(x, y, radius, true)
dc.DrawCachedText("Repeated text", x, y)
```

### Step 3: Use memory pooling
```go
// Use pooled contexts for temporary work
ctx := advancegg.PooledContext(400, 300)
defer advancegg.ReleaseContext(ctx)

// Use pooled paths
path := advancegg.PooledPath2D()
defer advancegg.ReleasePath2D(path)
```

### Step 4: SIMD optimizations
```go
// For image processing
img := dc.Image().(*image.RGBA)
blurred := advancegg.SIMDBlur(img, 5)
dc.PutImageData(advancegg.NewImageDataFromImage(blurred))
```

## Tutorial 6: Game Graphics

Learn how to create game-style graphics and UI elements.

### Step 1: Create a game character
```go
func drawCharacter(dc *advancegg.Context, x, y float64) {
    // Body
    dc.SetRGB(0.2, 0.4, 0.8)
    dc.DrawRectangle(x-10, y-30, 20, 30)
    dc.Fill()
    
    // Head
    dc.SetRGB(1, 0.8, 0.6)
    dc.DrawCircle(x, y-35, 8)
    dc.Fill()
    
    // Add details...
}
```

### Step 2: Create UI elements
```go
func drawHealthBar(dc *advancegg.Context, x, y, width, height, percentage float64) {
    // Background
    dc.SetRGB(0.3, 0.1, 0.1)
    dc.DrawRectangle(x, y, width, height)
    dc.Fill()
    
    // Health fill
    dc.SetRGB(0.8, 0.2, 0.2)
    dc.DrawRectangle(x+2, y+2, (width-4)*percentage, height-4)
    dc.Fill()
}
```

### Step 3: Add particle effects
```go
func drawFireEffect(dc *advancegg.Context, x, y float64, numParticles int) {
    for i := 0; i < numParticles; i++ {
        px := x + (rand.Float64()-0.5)*40
        py := y - rand.Float64()*100
        
        heat := rand.Float64()
        dc.SetRGBA(1.0, heat*0.8, heat*heat*0.3, 0.3+heat*0.7)
        dc.DrawCircle(px, py, 2+rand.Float64()*6)
        dc.Fill()
    }
}
```

## Tutorial 7: Data Visualization

Learn how to create professional data visualizations.

### Step 1: Prepare your data structure
```go
type DataPoint struct {
    X, Y  float64
    Label string
    Color color.Color
}

data := []DataPoint{
    {1, 10, "Point 1", color.RGBA{255, 0, 0, 255}},
    {2, 15, "Point 2", color.RGBA{0, 255, 0, 255}},
    // ... more data
}
```

### Step 2: Calculate scales and bounds
```go
minX, maxX := findMinMax(data, func(d DataPoint) float64 { return d.X })
minY, maxY := findMinMax(data, func(d DataPoint) float64 { return d.Y })

scaleX := chartWidth / (maxX - minX)
scaleY := chartHeight / (maxY - minY)
```

### Step 3: Draw the visualization
```go
for _, point := range data {
    screenX := chartX + (point.X-minX)*scaleX
    screenY := chartY + chartHeight - (point.Y-minY)*scaleY
    
    dc.SetColor(point.Color)
    dc.DrawCircle(screenX, screenY, 5)
    dc.Fill()
}
```

## Best Practices

1. **Always clear the background** before drawing
2. **Use appropriate coordinate systems** for your use case
3. **Handle edge cases** (empty data, zero values, etc.)
4. **Profile performance** for complex visualizations
5. **Use consistent styling** across your application
6. **Add proper error handling** for file operations
7. **Consider accessibility** (color blindness, contrast)

## Common Gotchas

- RGB values are 0.0-1.0, not 0-255 (unless using RGBA255)
- Text coordinates are at the baseline, not top-left
- Transformations affect all subsequent drawing operations
- Remember to call `Stroke()` or `Fill()` after drawing shapes
- Path operations accumulate until you start a new path

## Next Steps

- Explore the [examples](../examples/) directory
- Read the [API Reference](API_REFERENCE.md)
- Check out [Advanced Features](ADVANCED_FEATURES.md)
- Join the community discussions
