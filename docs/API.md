# AdvanceGG API Reference

## Core Context API

### Context Creation
```go
// Create a new drawing context
dc := advancegg.NewContext(width, height)

// Create with custom DPI
dc := advancegg.NewContextWithDPI(width, height, dpi)
```

### Drawing Operations

#### Basic Shapes
```go
// Rectangles
dc.DrawRectangle(x, y, width, height)
dc.DrawRoundedRectangle(x, y, width, height, radius)

// Circles and Ellipses
dc.DrawCircle(x, y, radius)
dc.DrawEllipse(x, y, rx, ry)

// Lines
dc.DrawLine(x1, y1, x2, y2)
dc.MoveTo(x, y)
dc.LineTo(x, y)
```

#### Path Operations
```go
// Start a new path
dc.NewPath()

// Move without drawing
dc.MoveTo(x, y)

// Draw lines
dc.LineTo(x, y)

// Curves
dc.QuadraticCurveTo(cpx, cpy, x, y)
dc.BezierCurveTo(cp1x, cp1y, cp2x, cp2y, x, y)

// Arcs
dc.Arc(x, y, radius, startAngle, endAngle, anticlockwise)
dc.ArcTo(x1, y1, x2, y2, radius)

// Close path
dc.ClosePath()
```

#### Fill and Stroke
```go
// Fill operations
dc.Fill()
dc.FillPreserve()

// Stroke operations
dc.Stroke()
dc.StrokePreserve()

// Combined
dc.FillAndStroke()
```

### Colors and Styles

#### Color Setting
```go
// RGB (0.0 to 1.0)
dc.SetRGB(r, g, b)

// RGBA with alpha
dc.SetRGBA(r, g, b, a)

// Hex colors
dc.SetHexColor("#FF5733")

// Named colors
dc.SetColor(color.RGBA{255, 0, 0, 255})
```

#### Gradients
```go
// Linear gradient
gradient := dc.NewLinearGradient(x1, y1, x2, y2)
gradient.AddColorStop(0, color.RGBA{255, 0, 0, 255})
gradient.AddColorStop(1, color.RGBA{0, 0, 255, 255})
dc.SetFillStyle(gradient)

// Radial gradient
radial := dc.NewRadialGradient(x1, y1, r1, x2, y2, r2)
radial.AddColorStop(0, color.RGBA{255, 255, 255, 255})
radial.AddColorStop(1, color.RGBA{0, 0, 0, 255})
dc.SetFillStyle(radial)
```

#### Line Styles
```go
// Line width
dc.SetLineWidth(width)

// Line cap styles
dc.SetLineCap(advancegg.LineCapRound)
dc.SetLineCap(advancegg.LineCapButt)
dc.SetLineCap(advancegg.LineCapSquare)

// Line join styles
dc.SetLineJoin(advancegg.LineJoinRound)
dc.SetLineJoin(advancegg.LineJoinBevel)
dc.SetLineJoin(advancegg.LineJoinMiter)

// Dashed lines
dc.SetDash([]float64{5, 5}) // 5 pixels on, 5 pixels off
dc.SetDashOffset(offset)
```

### Text Rendering

#### Basic Text
```go
// Draw text
dc.DrawString("Hello World", x, y)

// Anchored text (alignment)
dc.DrawStringAnchored("Centered", x, y, 0.5, 0.5)

// Wrapped text
dc.DrawStringWrapped("Long text...", x, y, ax, ay, width, lineSpacing, align)
```

#### Font Management
```go
// Load font from file
dc.LoadFontFace("path/to/font.ttf", size)

// Set font size
dc.SetFontSize(size)

// Text measurement
width := dc.MeasureString("text")
width, height := dc.MeasureStringWrapped("text", width)
```

#### Advanced Text Features
```go
// Text on paths
advance.DrawTextOnCircle(dc, "Circular Text", centerX, centerY, radius)
advance.DrawTextOnWave(dc, "Wave Text", startX, startY, endX, amplitude, frequency)
advance.DrawTextOnSpiral(dc, "Spiral Text", centerX, centerY, startRadius, endRadius, turns)

// Text with custom alignment
textOnPath := advance.NewSimpleTextOnPath("Custom Text")
textOnPath.Alignment = advance.SimpleAlignCenter
textOnPath.Spacing = 1.2
textOnPath.Offset = 10
```

### Transformations

#### Matrix Operations
```go
// Translation
dc.Translate(dx, dy)

// Rotation (in radians)
dc.Rotate(angle)

// Scaling
dc.Scale(sx, sy)

// Skewing
dc.Shear(sx, sy)

// Custom matrix
dc.SetMatrix(matrix)
dc.MultiplyMatrix(matrix)
```

#### Transformation Stack
```go
// Save current state
dc.Push()

// Apply transformations
dc.Translate(100, 100)
dc.Rotate(math.Pi / 4)

// Restore previous state
dc.Pop()
```

### Image Operations

#### Image Loading and Drawing
```go
// Load image
img, err := advancegg.LoadImage("image.png")
if err != nil {
    log.Fatal(err)
}

// Draw image
dc.DrawImage(img, x, y)

// Draw scaled image
dc.DrawImageAnchored(img, x, y, ax, ay)

// Get current image
currentImg := dc.Image()
```

#### Image Data Operations
```go
// Get image data
imageData := dc.GetImageData(x, y, width, height)

// Put image data
dc.PutImageData(imageData, x, y)

// Create image data from image
imageData := core.NewImageDataFromImage(img)
```

### High-Performance Filters

#### SIMD-Optimized Filters
```go
// Fast grayscale conversion
filtered := core.FastGrayscale(img)

// Fast brightness adjustment
filtered = core.FastBrightness(1.5)(img)

// Fast contrast adjustment
filtered = core.FastContrast(1.3)(img)

// Fast blur
filtered = core.FastBlur(radius)(img)

// Fast edge detection
filtered = core.FastEdgeDetection()(img)

// Fast sharpen
filtered = core.FastSharpen(amount)(img)
```

#### Batch and Parallel Processing
```go
// Batch multiple filters
batchFilter := core.BatchFilter(
    core.FastGrayscale,
    core.FastBrightness(1.2),
    core.FastContrast(1.1),
    core.FastBlur(3),
)
result := batchFilter(img)

// Parallel processing
parallelFilter := core.ParallelFilter(core.FastBlur(5), 4) // 4 workers
result = parallelFilter(img)
```

#### Standard Filters
```go
// Basic filters
filtered := core.Grayscale(img)
filtered = core.Invert(img)
filtered = core.Sepia(img)

// Adjustments
filtered = core.Brightness(factor)(img)
filtered = core.Contrast(factor)(img)
filtered = core.Gamma(gamma)(img)

// Effects
filtered = core.Blur(radius)(img)
filtered = core.Sharpen(img)
filtered = core.EdgeDetection(img)
filtered = core.Emboss(img)

// Color effects
filtered = core.Hue(degrees)(img)
filtered = core.Saturation(factor)(img)
filtered = core.Threshold(value)(img)
```

### Clipping and Masking

#### Clipping Paths
```go
// Set clipping region
dc.DrawCircle(centerX, centerY, radius)
dc.Clip()

// Draw within clipped region
dc.SetRGB(1, 0, 0)
dc.DrawRectangle(0, 0, width, height)
dc.Fill()

// Reset clipping
dc.ResetClip()
```

#### Masking
```go
// Create mask
mask := advancegg.NewContext(width, height)
mask.DrawCircle(centerX, centerY, radius)
mask.Fill()

// Apply mask
dc.SetMask(mask.Image())
```

### Export and Save

#### Image Formats
```go
// PNG (lossless)
dc.SavePNG("output.png")

// JPEG with quality
dc.SaveJPEG("output.jpg", 95) // Quality 0-100

// PDF
dc.SavePDF("output.pdf")

// Custom writer
var buf bytes.Buffer
dc.EncodePNG(&buf)
```

### Advanced Features

#### Emoji Rendering
```go
// Get emoji renderer
renderer := dc.GetEmojiRenderer()

// Load emoji font
err := renderer.LoadEmojiFont("path/to/emoji.ttf")

// Render emoji text
dc.DrawString("Hello 👋 World 🌍", x, y)
```

#### Performance Monitoring
```go
// Enable performance tracking
dc.SetPerformanceTracking(true)

// Get performance metrics
metrics := dc.GetPerformanceMetrics()
fmt.Printf("Render time: %v\n", metrics.RenderTime)
fmt.Printf("Memory usage: %d bytes\n", metrics.MemoryUsage)
```

## Constants and Enums

### Line Caps
```go
advancegg.LineCapButt    // Flat line endings
advancegg.LineCapRound   // Rounded line endings
advancegg.LineCapSquare  // Square line endings
```

### Line Joins
```go
advancegg.LineJoinMiter  // Sharp corners
advancegg.LineJoinRound  // Rounded corners
advancegg.LineJoinBevel  // Beveled corners
```

### Text Alignment
```go
advance.SimpleAlignStart   // Left alignment
advance.SimpleAlignCenter  // Center alignment
advance.SimpleAlignEnd     // Right alignment
```

### Blend Modes
```go
advancegg.BlendModeNormal      // Normal blending
advancegg.BlendModeMultiply    // Multiply blending
advancegg.BlendModeScreen      // Screen blending
advancegg.BlendModeOverlay     // Overlay blending
// ... and more
```

## Error Handling

Most operations in AdvanceGG are designed to be safe and will not panic. However, some operations may return errors:

```go
// Font loading
err := dc.LoadFontFace("font.ttf", 16)
if err != nil {
    log.Printf("Failed to load font: %v", err)
    // Use default font
}

// Image loading
img, err := advancegg.LoadImage("image.png")
if err != nil {
    log.Printf("Failed to load image: %v", err)
    return
}

// File saving
err = dc.SavePNG("output.png")
if err != nil {
    log.Printf("Failed to save PNG: %v", err)
}
```

## Best Practices

### Performance
1. **Reuse contexts** when possible to avoid allocation overhead
2. **Use SIMD filters** for image processing operations
3. **Batch operations** when applying multiple filters
4. **Pre-load fonts** at application startup
5. **Use appropriate image formats** (PNG for graphics, JPEG for photos)

### Memory Management
1. **Clear large images** when no longer needed
2. **Use memory pooling** for frequent operations
3. **Avoid creating many small contexts**
4. **Monitor memory usage** in long-running applications

### Quality
1. **Use high DPI** for crisp output on modern displays
2. **Enable anti-aliasing** for smooth graphics
3. **Choose appropriate line widths** for your output resolution
4. **Test on different platforms** to ensure consistency
