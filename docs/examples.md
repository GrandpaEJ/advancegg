# Examples

This document contains various examples demonstrating the capabilities of AdvanceGG.

## Basic Examples

### Hello Circle

The simplest example - drawing a black circle:

```go
package main

import "github.com/GrandpaEJ/advancegg"

func main() {
    dc := advancegg.NewContext(1000, 1000)
    dc.DrawCircle(500, 500, 400)
    dc.SetRGB(0, 0, 0)
    dc.Fill()
    dc.SavePNG("circle.png")
}
```

### Colorful Rectangle

Drawing a red rectangle with transparency:

```go
package main

import "github.com/GrandpaEJ/advancegg"

func main() {
    dc := advancegg.NewContext(800, 600)
    
    // Set background to white
    dc.SetRGB(1, 1, 1)
    dc.Clear()
    
    // Draw a semi-transparent red rectangle
    dc.SetRGBA(1, 0, 0, 0.5)
    dc.DrawRectangle(200, 150, 400, 300)
    dc.Fill()
    
    dc.SavePNG("rectangle.png")
}
```

## Intermediate Examples

### Multiple Shapes

Combining different shapes:

```go
package main

import "github.com/GrandpaEJ/advancegg"

func main() {
    dc := advancegg.NewContext(800, 600)
    
    // White background
    dc.SetRGB(1, 1, 1)
    dc.Clear()
    
    // Blue circle
    dc.SetRGB(0, 0, 1)
    dc.DrawCircle(200, 200, 100)
    dc.Fill()
    
    // Green rectangle
    dc.SetRGB(0, 1, 0)
    dc.DrawRectangle(400, 150, 200, 100)
    dc.Fill()
    
    // Red line
    dc.SetRGB(1, 0, 0)
    dc.SetLineWidth(5)
    dc.DrawLine(100, 400, 700, 400)
    dc.Stroke()
    
    dc.SavePNG("shapes.png")
}
```

### Text Rendering

Drawing text with different styles:

```go
package main

import "github.com/GrandpaEJ/advancegg"

func main() {
    dc := advancegg.NewContext(800, 600)

    // White background
    dc.SetRGB(1, 1, 1)
    dc.Clear()

    // Black text
    dc.SetRGB(0, 0, 0)
    dc.DrawString("Hello, AdvanceGG!", 50, 100)

    // Centered text
    dc.DrawStringAnchored("Centered Text", 400, 300, 0.5, 0.5)

    dc.SavePNG("text.png")
}
```

### Font Loading (TTF and OTF Support)

AdvanceGG supports both TTF and OTF font formats:

```go
package main

import "github.com/GrandpaEJ/advancegg"

func main() {
    dc := advancegg.NewContext(600, 400)
    dc.SetRGB(1, 1, 1)
    dc.Clear()

    // Load TTF font
    err := dc.LoadTTFFace("/path/to/font.ttf", 24)
    if err == nil {
        dc.SetRGB(1, 0, 0)
        dc.DrawString("TTF Font Text", 50, 100)
    }

    // Load OTF font
    err = dc.LoadOTFFace("/path/to/font.otf", 24)
    if err == nil {
        dc.SetRGB(0, 0, 1)
        dc.DrawString("OTF Font Text", 50, 150)
    }

    // Auto-detect format
    err = dc.LoadFontFace("/path/to/font.ttf", 20)
    if err == nil {
        dc.SetRGB(0, 1, 0)
        dc.DrawString("Auto-detected Font", 50, 200)
    }

    dc.SavePNG("fonts.png")
}
```

## Advanced Examples

### Gradients

Creating linear gradients:

```go
package main

import "github.com/GrandpaEJ/advancegg"

func main() {
    dc := advancegg.NewContext(800, 600)
    
    // Create a linear gradient
    gradient := advancegg.NewLinearGradient(0, 0, 800, 0)
    gradient.AddColorStop(0, advancegg.Color{R: 1, G: 0, B: 0, A: 1}) // Red
    gradient.AddColorStop(1, advancegg.Color{R: 0, G: 0, B: 1, A: 1}) // Blue
    
    dc.SetFillStyle(gradient)
    dc.DrawRectangle(0, 0, 800, 600)
    dc.Fill()
    
    dc.SavePNG("gradient.png")
}
```

### Transformations

Using transformations to create interesting effects:

```go
package main

import (
    "math"
    "github.com/GrandpaEJ/advancegg"
)

func main() {
    dc := advancegg.NewContext(800, 800)
    
    // White background
    dc.SetRGB(1, 1, 1)
    dc.Clear()
    
    // Draw rotating rectangles
    dc.SetRGB(0, 0, 0)
    for i := 0; i < 12; i++ {
        dc.Push()
        dc.RotateAbout(float64(i)*math.Pi/6, 400, 400)
        dc.DrawRectangle(350, 390, 100, 20)
        dc.Fill()
        dc.Pop()
    }
    
    dc.SavePNG("rotation.png")
}
```

## Running Examples

All examples are available in the `examples/` directory. To run an example:

```bash
cd examples
go run circle.go
```

This will generate the corresponding PNG file in the examples directory.

## Path2D Examples

Path2D provides advanced path manipulation capabilities for creating reusable path objects.

### Basic Path2D Usage

```go
package main

import "github.com/GrandpaEJ/advancegg"

func main() {
    dc := advancegg.NewContext(400, 300)
    dc.SetRGB(1, 1, 1)
    dc.Clear()

    // Create a reusable path
    path := advancegg.NewPath2D()
    path.MoveTo(50, 50)
    path.LineTo(150, 50)
    path.LineTo(100, 150)
    path.ClosePath()

    // Use the path multiple times
    dc.SetRGB(1, 0, 0)
    dc.FillPath2D(path)

    dc.Translate(200, 0)
    dc.SetRGB(0, 0, 1)
    dc.StrokePath2D(path)

    dc.SavePNG("path2d-basic.png")
}
```

### Path2D with Curves and Arcs

```go
// Create a path with various curve types
curvePath := advancegg.NewPath2D()
curvePath.MoveTo(100, 200)
curvePath.QuadraticCurveTo(200, 100, 300, 200)
curvePath.BezierCurveTo(350, 250, 400, 150, 450, 200)
curvePath.Arc(500, 200, 50, 0, math.Pi, false)

dc.SetRGB(0, 1, 0)
dc.SetLineWidth(3)
dc.StrokePath2D(curvePath)
```

### Combining Multiple Paths

```go
// Create multiple paths and combine them
combinedPath := advancegg.NewPath2D()

// Add a circle
circlePath := advancegg.NewPath2D()
circlePath.Arc(200, 200, 50, 0, 2*math.Pi, false)
combinedPath.AddPath(circlePath)

// Add a rectangle
rectPath := advancegg.NewPath2D()
rectPath.Rect(150, 150, 100, 100)
combinedPath.AddPath(rectPath)

// Draw the combined path
dc.SetRGB(0.5, 0, 0.5)
dc.FillPath2D(combinedPath)
```

## More Examples

For more complex examples, check out the files in the `examples/` directory:

- `path2d-basic.go` - Basic Path2D usage
- `path2d-reuse.go` - Reusing paths for multiple drawings
- `path2d-advanced.go` - Advanced Path2D techniques with clipping
- `beziers.go` - Bézier curve demonstrations
- `gradient-*.go` - Various gradient examples
- `pattern-fill.go` - Pattern filling examples
- `text.go` - Advanced text rendering
- `stars.go` - Complex shape drawing
- And many more!
