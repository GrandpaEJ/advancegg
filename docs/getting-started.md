# Getting Started with AdvanceGG

## Installation

To install AdvanceGG, use `go get`:

```bash
go get github.com/GrandpaEJ/advancegg
```

## Basic Usage

Here's a simple example that creates a black circle:

```go
package main

import "github.com/GrandpaEJ/advancegg"

func main() {
    // Create a new context with 1000x1000 pixels
    dc := advancegg.NewContext(1000, 1000)
    
    // Draw a circle at the center
    dc.DrawCircle(500, 500, 400)
    
    // Set the color to black
    dc.SetRGB(0, 0, 0)
    
    // Fill the circle
    dc.Fill()
    
    // Save as PNG
    dc.SavePNG("circle.png")
}
```

## Creating Contexts

There are several ways to create a drawing context:

```go
// Create a new context with specified dimensions
dc := advancegg.NewContext(width, height)

// Create a context from an existing image
dc := advancegg.NewContextForImage(img)

// Create a context from an RGBA image
dc := advancegg.NewContextForRGBA(rgbaImg)
```

## Basic Drawing

### Shapes

```go
// Rectangle
dc.DrawRectangle(x, y, width, height)

// Circle
dc.DrawCircle(centerX, centerY, radius)

// Line
dc.DrawLine(x1, y1, x2, y2)
```

### Colors

```go
// RGB values (0.0 to 1.0)
dc.SetRGB(r, g, b)

// RGBA values (0.0 to 1.0)
dc.SetRGBA(r, g, b, a)

// RGB values (0 to 255)
dc.SetRGB255(r, g, b)

// Hex color
dc.SetHexColor("#FF0000")
```

### Rendering

```go
// Fill the current path
dc.Fill()

// Stroke the current path
dc.Stroke()

// Clear the entire context
dc.Clear()
```

## Next Steps

- Check out the [API Reference](api-reference.md) for complete function documentation
- Browse the [Examples](examples.md) for more complex usage patterns
- Learn about advanced features like gradients, transformations, and text rendering
