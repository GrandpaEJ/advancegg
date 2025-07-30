# AdvanceGG - Advanced 2D Graphics Library for Go

[![Go Version](https://img.shields.io/badge/Go-1.18+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Documentation](https://img.shields.io/badge/docs-online-blue.svg)](index.html)

AdvanceGG is a powerful, production-ready 2D graphics library for Go that provides comprehensive drawing capabilities, advanced text rendering, image processing, and much more. Built for performance and ease of use, it's perfect for creating visualizations, games, charts, and any application requiring high-quality 2D graphics.

**Originally forked from [GG](https://github.com/fogleman/gg) and significantly enhanced with professional-grade features.**

## ✨ Key Features

### 🎨 **Professional Graphics**
- **Layer System** - Multi-layered drawing with blend modes and opacity
- **Advanced Strokes** - Dashed patterns, gradient strokes, tapered ends
- **Path2D Support** - Complex vector paths and Bézier curves
- **High-Quality Rendering** - Anti-aliased graphics with sub-pixel precision

### 📝 **Advanced Typography**
- **Unicode Shaping** - Full support for international text and complex scripts
- **Emoji Rendering** - Color emoji with fallback mechanisms
- **Text-on-Path** - Text following curves, circles, and custom paths
- **Font Management** - TTF/OTF loading with advanced metrics

### 🖼️ **Image Processing**
- **15+ Filters** - Blur, sharpen, edge detection, sepia, and more
- **Color Management** - ICC color profiles and accurate color conversion
- **Multiple Formats** - PNG, JPEG, GIF, WebP, TIFF, BMP support
- **Non-destructive Editing** - Reversible filters and transformations

### ⚡ **Performance & Developer Experience**
- **SIMD Optimizations** - CPU vector instructions for image processing
- **Memory Pooling** - Reduced GC pressure and improved performance
- **Debug Mode** - Visual debugging tools and performance profiling
- **Comprehensive Examples** - 50+ examples covering all features

## 🚀 Installation

```bash
go get github.com/GrandpaEJ/advancegg
```

**Requirements:** Go 1.18 or later

## 📖 Documentation

**[📚 Complete Documentation & Examples](index.html)** - Interactive web documentation with live examples

- **[🚀 Getting Started](docs/getting-started.html)** - Quick start guide and tutorials
- **[📋 API Reference](docs/api/)** - Complete API documentation with examples
- **[🎨 Examples Gallery](docs/examples/)** - Practical examples and use cases
- **[💡 Interactive Demo](index.html)** - Live documentation website

### Quick Links
- [Basic Drawing](docs/examples/index.html#basic-drawing) - Shapes, paths, and curves
- [Text & Typography](docs/examples/index.html#text-typography) - Advanced text rendering
- [Image Processing](docs/examples/index.html#image-processing) - Filters and effects
- [Layer System](docs/examples/index.html#layer-system) - Multi-layer compositing
- [Data Visualization](docs/examples/index.html#data-visualization) - Charts and graphs

## 🎯 Quick Start

### Hello, Circle!

```go
package main

import "github.com/GrandpaEJ/advancegg"

func main() {
    // Create a new 800x600 canvas
    dc := advancegg.NewContext(800, 600)

    // Set background color
    dc.SetRGB(0.1, 0.1, 0.3)
    dc.Clear()

    // Draw a red circle with gradient
    dc.SetRGB(1, 0, 0)
    dc.DrawCircle(400, 300, 100)
    dc.Fill()

    // Add text with shadow
    dc.SetRGB(1, 1, 1)
    dc.DrawString("Hello AdvanceGG!", 300, 350)

    // Save as PNG
    dc.SavePNG("hello.png")
}
```

### Advanced Example - Data Visualization

```go
// Create a beautiful bar chart
dc := advancegg.NewContext(800, 600)

// Gradient background
gradient := advancegg.NewLinearGradient(0, 0, 0, 600)
gradient.AddColorStop(0, color.RGBA{240, 248, 255, 255})
gradient.AddColorStop(1, color.RGBA{230, 240, 250, 255})
dc.SetFillStyle(gradient)
dc.DrawRectangle(0, 0, 800, 600)
dc.Fill()

// Draw bars with different colors
data := []float64{85, 92, 78, 96, 88}
colors := []color.Color{
    color.RGBA{255, 99, 132, 255},
    color.RGBA{54, 162, 235, 255},
    color.RGBA{255, 205, 86, 255},
    color.RGBA{75, 192, 192, 255},
    color.RGBA{153, 102, 255, 255},
}

for i, value := range data {
    x := 100 + float64(i)*120
    height := value * 4

    dc.SetColor(colors[i])
    dc.DrawRectangle(x, 500-height, 80, height)
    dc.Fill()

    // Add value labels
    dc.SetRGB(0, 0, 0)
    dc.DrawStringAnchored(fmt.Sprintf("%.0f", value), x+40, 480-height, 0.5, 0.5)
}

dc.SavePNG("chart.png")
```

## 🎨 Examples Gallery

Explore our comprehensive [examples gallery](docs/examples/) with over 50 practical examples:

- **[Basic Drawing](examples/basic-shapes.go)** - Shapes, paths, and curves
- **[Text Effects](examples/text-effects.go)** - Typography and text-on-path
- **[Image Filters](examples/image-filters.go)** - Professional image processing
- **[Layer System](examples/layer-system.go)** - Multi-layer compositing
- **[Data Visualization](examples/data-visualization.go)** - Charts and graphs
- **[Game Graphics](examples/game-graphics.go)** - Sprites and animations

![AdvanceGG Examples](images/examples-showcase.png)

## Creating Contexts

There are a few ways of creating a context.

```go
NewContext(width, height int) *Context
NewContextForImage(im image.Image) *Context
NewContextForRGBA(im *image.RGBA) *Context
```

## Drawing Functions

Ever used a graphics library that didn't have functions for drawing rectangles
or circles? What a pain!

```go
DrawPoint(x, y, r float64)
DrawLine(x1, y1, x2, y2 float64)
DrawRectangle(x, y, w, h float64)
DrawRoundedRectangle(x, y, w, h, r float64)
DrawCircle(x, y, r float64)
DrawArc(x, y, r, angle1, angle2 float64)
DrawEllipse(x, y, rx, ry float64)
DrawEllipticalArc(x, y, rx, ry, angle1, angle2 float64)
DrawRegularPolygon(n int, x, y, r, rotation float64)
DrawImage(im image.Image, x, y int)
DrawImageAnchored(im image.Image, x, y int, ax, ay float64)
SetPixel(x, y int)

MoveTo(x, y float64)
LineTo(x, y float64)
QuadraticTo(x1, y1, x2, y2 float64)
CubicTo(x1, y1, x2, y2, x3, y3 float64)
ClosePath()
ClearPath()
NewSubPath()

Clear()
Stroke()
Fill()
StrokePreserve()
FillPreserve()
```

It is often desired to center an image at a point. Use `DrawImageAnchored` with `ax` and `ay` set to 0.5 to do this. Use 0 to left or top align. Use 1 to right or bottom align. `DrawStringAnchored` does the same for text, so you don't need to call `MeasureString` yourself.

## Text Functions

It will even do word wrap for you!

```go
DrawString(s string, x, y float64)
DrawStringAnchored(s string, x, y, ax, ay float64)
DrawStringWrapped(s string, x, y, ax, ay, width, lineSpacing float64, align Align)
MeasureString(s string) (w, h float64)
MeasureMultilineString(s string, lineSpacing float64) (w, h float64)
WordWrap(s string, w float64) []string
SetFontFace(fontFace font.Face)
LoadFontFace(path string, points float64) error
```

## Color Functions

Colors can be set in several different ways for your convenience.

```go
SetRGB(r, g, b float64)
SetRGBA(r, g, b, a float64)
SetRGB255(r, g, b int)
SetRGBA255(r, g, b, a int)
SetColor(c color.Color)
SetHexColor(x string)
```

## Stroke & Fill Options

```go
SetLineWidth(lineWidth float64)
SetLineCap(lineCap LineCap)
SetLineJoin(lineJoin LineJoin)
SetDash(dashes ...float64)
SetDashOffset(offset float64)
SetFillRule(fillRule FillRule)
```

## Gradients & Patterns

`gg` supports linear, radial and conic gradients and surface patterns. You can also implement your own patterns.

```go
SetFillStyle(pattern Pattern)
SetStrokeStyle(pattern Pattern)
NewSolidPattern(color color.Color)
NewLinearGradient(x0, y0, x1, y1 float64)
NewRadialGradient(x0, y0, r0, x1, y1, r1 float64)
NewConicGradient(cx, cy, deg float64)
NewSurfacePattern(im image.Image, op RepeatOp)
```

## Transformation Functions

```go
Identity()
Translate(x, y float64)
Scale(x, y float64)
Rotate(angle float64)
Shear(x, y float64)
ScaleAbout(sx, sy, x, y float64)
RotateAbout(angle, x, y float64)
ShearAbout(sx, sy, x, y float64)
TransformPoint(x, y float64) (tx, ty float64)
InvertY()
```

It is often desired to rotate or scale about a point that is not the origin. The functions `RotateAbout`, `ScaleAbout`, `ShearAbout` are provided as a convenience.

`InvertY` is provided in case Y should increase from bottom to top vs. the default top to bottom.

## Stack Functions

Save and restore the state of the context. These can be nested.

```go
Push()
Pop()
```

## Clipping Functions

Use clipping regions to restrict drawing operations to an area that you
defined using paths.

```go
Clip()
ClipPreserve()
ResetClip()
AsMask() *image.Alpha
SetMask(mask *image.Alpha)
InvertMask()
```

## Helper Functions

Sometimes you just don't want to write these yourself.

```go
Radians(degrees float64) float64
Degrees(radians float64) float64
LoadImage(path string) (image.Image, error)
LoadPNG(path string) (image.Image, error)
SavePNG(path string, im image.Image) error
```

![Separator](http://i.imgur.com/fsUvnPB.png)

## Another Example

See the output of this example below.

```go
package main

import "github.com/GrandpaEJ/advancegg"

func main() {
	const S = 1024
	dc := advancegg.NewContext(S, S)
	dc.SetRGBA(0, 0, 0, 0.1)
	for i := 0; i < 360; i += 15 {
		dc.Push()
		dc.RotateAbout(advancegg.Radians(float64(i)), S/2, S/2)
		dc.DrawEllipse(S/2, S/2, S*7/16, S/8)
		dc.Fill()
		dc.Pop()
	}
	dc.SavePNG("out.png")
}
```

![Ellipses](http://i.imgur.com/J9CBZef.png)

## Project Structure

```
AdvanceGG/
├── main.go               # Main library entry point (re-exports from internal/core)
├── cli.go                # CLI utilities and command-line interface
├── internal/             # Internal packages
│   └── core/            # Core graphics library implementation
├── examples/             # Example programs
│   └── images/          # Image assets for examples
├── docs/                # Documentation
│   ├── getting-started.md # Getting started guide
│   ├── api-reference.md   # Complete API documentation
│   ├── examples.md        # Example gallery
│   └── contributing.md    # Contributing guidelines
├── cmd/                 # Command-line tools (for future use)
├── src/                 # Additional source files (for future use)
├── CHANGELOG.md         # Version history
├── README.md           # This file
├── LICENSE.md          # License information
└── go.mod             # Go module definition
```

## Features

- **Simple API**: Easy-to-use functions for common graphics operations
- **Rich Drawing Functions**: Support for shapes, paths, text, and images
- **Advanced Features**: Gradients, patterns, transformations, and clipping
- **Pure Go**: No external dependencies beyond the Go standard library and golang.org/x packages
- **High Performance**: Optimized for speed and memory efficiency
- **Comprehensive Documentation**: Complete guides, API reference, and examples

## Quick Start

1. Install the library:
   ```bash
   go get github.com/GrandpaEJ/advancegg
   ```

2. Create your first graphic:
   ```go
   package main

   import "github.com/GrandpaEJ/advancegg"

   func main() {
       dc := advancegg.NewContext(800, 600)
       dc.SetRGB(1, 1, 1) // White background
       dc.Clear()

       dc.SetRGB(0, 0, 1) // Blue color
       dc.DrawCircle(400, 300, 100)
       dc.Fill()

       dc.SavePNG("my_first_graphic.png")
   }
   ```

3. Run your program:
   ```bash
   go run main.go
   ```

## Contributing

We welcome contributions! Please see our [Contributing Guide](docs/contributing.md) for details on how to get started.

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details.
