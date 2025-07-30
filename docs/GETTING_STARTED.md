# Getting Started with AdvanceGG

Welcome to AdvanceGG, a high-performance 2D graphics library for Go! This guide will help you get up and running quickly.

## Installation

```bash
go get github.com/GrandpaEJ/advancegg
```

## Your First Drawing

Let's create a simple "Hello World" example:

```go
package main

import "github.com/GrandpaEJ/advancegg"

func main() {
    // Create a new drawing context
    dc := advancegg.NewContext(800, 600)
    
    // Set background color to white
    dc.SetRGB(1, 1, 1)
    dc.Clear()
    
    // Set drawing color to black
    dc.SetRGB(0, 0, 0)
    
    // Draw some text
    dc.DrawString("Hello, AdvanceGG!", 50, 50)
    
    // Draw a circle
    dc.SetRGB(1, 0, 0) // Red
    dc.DrawCircle(400, 300, 50)
    dc.Fill()
    
    // Save the image
    dc.SavePNG("hello.png")
}
```

## Basic Concepts

### Context
The `Context` is the main drawing surface. All drawing operations are performed on a context.

```go
dc := advancegg.NewContext(width, height)
```

### Colors
Colors can be set in multiple ways:

```go
// RGB values (0.0 to 1.0)
dc.SetRGB(1, 0, 0) // Red

// RGBA with alpha
dc.SetRGBA(1, 0, 0, 0.5) // Semi-transparent red

// RGB values (0 to 255)
dc.SetRGBA255(255, 0, 0, 255)

// Hex colors
dc.SetHexColor("#FF0000")

// Color spaces
dc.SetHSV(0, 1, 1) // Red in HSV
dc.SetCMYK(0, 1, 1, 0) // Red in CMYK
```

### Basic Shapes

```go
// Rectangle
dc.DrawRectangle(x, y, width, height)
dc.Fill() // or dc.Stroke()

// Circle
dc.DrawCircle(centerX, centerY, radius)
dc.Fill()

// Line
dc.DrawLine(x1, y1, x2, y2)
dc.Stroke()

// Rounded rectangle
dc.DrawRoundedRectangle(x, y, width, height, radius)
dc.Fill()
```

### Text Drawing

```go
// Simple text
dc.DrawString("Hello", x, y)

// Anchored text (centered, etc.)
dc.DrawStringAnchored("Centered", x, y, 0.5, 0.5)

// Wrapped text
dc.DrawStringWrapped("Long text...", x, y, ax, ay, width, lineSpacing, align)
```

### Paths
For complex shapes, use paths:

```go
// Start a new path
dc.MoveTo(100, 100)
dc.LineTo(200, 100)
dc.LineTo(150, 200)
dc.ClosePath()
dc.Fill()

// Or use Path2D for reusable paths
path := advancegg.NewPath2D()
path.MoveTo(100, 100)
path.LineTo(200, 100)
path.LineTo(150, 200)
path.ClosePath()

dc.DrawPath(path)
dc.Fill()
```

## Saving Images

AdvanceGG supports multiple image formats:

```go
dc.SavePNG("output.png")
dc.SaveJPEG("output.jpg", 90) // 90% quality
dc.SaveGIF("output.gif")
dc.SaveBMP("output.bmp")
dc.SaveTIFF("output.tiff")
dc.SaveWebP("output.webp", 90)
```

## Next Steps

- Check out the [examples](../examples/) directory for more complex examples
- Read the [API Reference](API_REFERENCE.md) for detailed documentation
- Explore [Advanced Features](ADVANCED_FEATURES.md) for performance optimization
- See [Tutorials](TUTORIALS.md) for step-by-step guides

## Common Patterns

### Drawing with Transparency
```go
dc.SetRGBA(1, 0, 0, 0.5) // 50% transparent red
dc.DrawCircle(100, 100, 50)
dc.Fill()
```

### Using Transformations
```go
dc.Push() // Save current state
dc.Translate(100, 100)
dc.Rotate(math.Pi / 4) // 45 degrees
dc.DrawRectangle(-25, -25, 50, 50)
dc.Fill()
dc.Pop() // Restore state
```

### Gradients and Patterns
```go
// Linear gradient
gradient := advancegg.CreateLinearGradient(0, 0, 200, 0,
    color.RGBA{255, 0, 0, 255},   // Red
    color.RGBA{0, 0, 255, 255},   // Blue
)
dc.SetFillPattern(gradient)
dc.DrawRectangle(0, 0, 200, 100)
dc.Fill()
```

### Image Processing
```go
// Apply filters
dc.ApplyFilter(advancegg.Blur(5))
dc.ApplyFilter(advancegg.Brightness(1.2))
dc.ApplyFilter(advancegg.Contrast(1.1))

// Pixel manipulation
imageData := dc.GetImageData()
// Modify pixels...
dc.PutImageData(imageData)
```

## Performance Tips

1. **Use batch operations** for drawing many similar shapes
2. **Enable caching** for repeated elements
3. **Use memory pooling** for frequent allocations
4. **Consider SIMD optimizations** for image processing
5. **Profile your code** to identify bottlenecks

## Troubleshooting

### Common Issues

**Q: My text isn't showing up**
A: Make sure you've set a color before drawing text, and check that the coordinates are within the canvas bounds.

**Q: Colors look wrong**
A: Remember that RGB values are 0.0-1.0, not 0-255. Use `SetRGBA255()` for 0-255 values.

**Q: Performance is slow**
A: Consider using batch operations, enabling caching, or using SIMD optimizations for image processing.

**Q: Memory usage is high**
A: Enable memory pooling and clear unused resources regularly.

## Getting Help

- Check the [FAQ](FAQ.md)
- Browse [examples](../examples/)
- Read the [API documentation](API_REFERENCE.md)
- Report issues on GitHub

Happy drawing with AdvanceGG! 🎨
