# AdvanceGG Examples

This directory contains comprehensive examples demonstrating all features of AdvanceGG. Each example is self-contained and includes detailed comments explaining the concepts.

## 🚀 Getting Started

### Quick Setup
```bash
# Clone the repository
git clone https://github.com/GrandpaEJ/advancegg.git
cd advancegg/examples

# Run any example
go run basic-drawing.go
go run text-rendering.go
go run filter-performance-test.go
```

## 📚 Example Categories

### 🎯 Basic Drawing
Learn the fundamentals of 2D graphics with AdvanceGG.

| Example | Description | Key Features |
|---------|-------------|--------------|
| **[basic-drawing.go](basic-drawing.go)** | Essential shapes and colors | Rectangles, circles, lines, colors |
| **[paths-and-curves.go](paths-and-curves.go)** | Complex path operations | Bézier curves, arcs, path manipulation |
| **[gradients-patterns.go](gradients-patterns.go)** | Advanced fill styles | Linear/radial gradients, patterns |
| **[transformations.go](transformations.go)** | Matrix transformations | Translate, rotate, scale, skew |

### 📝 Typography & Text
Master text rendering and advanced typography features.

| Example | Description | Key Features |
|---------|-------------|--------------|
| **[text-rendering.go](text-rendering.go)** | Basic text operations | Font loading, text measurement, alignment |
| **[text-on-path-test.go](text-on-path-test.go)** | Advanced text layouts | Text on circles, waves, custom paths |
| **[font-loading-test.go](font-loading-test.go)** | Font management | Unicode support, emoji rendering, fallbacks |
| **[unicode-text.go](unicode-text.go)** | International text | Arabic, Chinese, Japanese, emoji |

### 🖼️ Image Processing
Explore high-performance image manipulation and filters.

| Example | Description | Key Features |
|---------|-------------|--------------|
| **[filter-performance-test.go](filter-performance-test.go)** | Optimized image filters | SIMD acceleration, parallel processing |
| **[image-manipulation.go](image-manipulation.go)** | Basic image operations | Loading, scaling, compositing |
| **[custom-filters.go](custom-filters.go)** | Creating custom filters | SIMD transforms, filter chains |
| **[batch-processing.go](batch-processing.go)** | Efficient batch operations | Multiple images, automated workflows |

### 📊 Data Visualization
Create professional charts and dashboards.

| Example | Description | Key Features |
|---------|-------------|--------------|
| **[bar-charts.go](bar-charts.go)** | Professional bar charts | Gradients, shadows, labels |
| **[line-graphs.go](line-graphs.go)** | Interactive line graphs | Multiple series, legends, axes |
| **[pie-charts.go](pie-charts.go)** | Beautiful pie charts | 3D effects, animations, tooltips |
| **[dashboard.go](dashboard.go)** | Complete dashboard | Multiple chart types, layouts |

### 🎮 Game Graphics
Game development assets and effects.

| Example | Description | Key Features |
|---------|-------------|--------------|
| **[sprite-rendering.go](sprite-rendering.go)** | Game sprite systems | Sprite sheets, animations, batching |
| **[particle-effects.go](particle-effects.go)** | Particle systems | Fire, smoke, explosions, physics |
| **[ui-elements.go](ui-elements.go)** | Game UI components | Buttons, health bars, menus |
| **[tilemap-rendering.go](tilemap-rendering.go)** | Tile-based graphics | Map rendering, collision detection |

### 🎨 Creative & Artistic
Explore creative coding and generative art.

| Example | Description | Key Features |
|---------|-------------|--------------|
| **[generative-art.go](generative-art.go)** | Algorithmic art creation | Fractals, patterns, randomization |
| **[creative-effects.go](creative-effects.go)** | Artistic visual effects | Distortions, color manipulations |
| **[mandala-generator.go](mandala-generator.go)** | Geometric pattern creation | Symmetry, mathematical art |
| **[abstract-compositions.go](abstract-compositions.go)** | Abstract art generation | Color theory, composition rules |

### 🔬 Scientific Visualization
Technical and scientific graphics.

| Example | Description | Key Features |
|---------|-------------|--------------|
| **[mathematical-plots.go](mathematical-plots.go)** | Function plotting | 2D/3D plots, mathematical functions |
| **[scientific-charts.go](scientific-charts.go)** | Research visualizations | Error bars, statistical plots |
| **[heatmaps.go](heatmaps.go)** | Data density visualization | Color mapping, interpolation |
| **[network-graphs.go](network-graphs.go)** | Graph visualization | Node-link diagrams, force layouts |

## 🎯 Feature Demonstrations

### Performance Features
- **[filter-performance-test.go](filter-performance-test.go)** - SIMD optimizations showcase
- **[parallel-processing.go](parallel-processing.go)** - Multi-core rendering
- **[memory-optimization.go](memory-optimization.go)** - Efficient memory usage

### Advanced Graphics
- **[layer-compositing.go](layer-compositing.go)** - Multi-layer blend modes
- **[clipping-masking.go](clipping-masking.go)** - Advanced masking techniques
- **[vector-graphics.go](vector-graphics.go)** - Scalable vector operations

### Text & Typography
- **[advanced-typography.go](advanced-typography.go)** - Professional text layout
- **[emoji-rendering.go](emoji-rendering.go)** - Color emoji support
- **[text-effects.go](text-effects.go)** - Shadows, outlines, gradients

## 🏃‍♂️ Quick Start Examples

### Hello World
```go
package main

import "github.com/GrandpaEJ/advancegg"

func main() {
    dc := advancegg.NewContext(800, 600)
    dc.SetRGB(1, 1, 1)
    dc.Clear()
    
    dc.SetRGB(0, 0, 0)
    dc.DrawString("Hello, AdvanceGG!", 400, 300)
    
    dc.SavePNG("hello.png")
}
```

### Gradient Circle
```go
package main

import "github.com/GrandpaEJ/advancegg"

func main() {
    dc := advancegg.NewContext(400, 400)
    
    // Gradient background
    gradient := dc.NewRadialGradient(200, 200, 0, 200, 200, 150)
    gradient.AddColorStop(0, color.RGBA{255, 100, 100, 255})
    gradient.AddColorStop(1, color.RGBA{100, 100, 255, 255})
    
    dc.SetFillStyle(gradient)
    dc.DrawCircle(200, 200, 150)
    dc.Fill()
    
    dc.SavePNG("gradient-circle.png")
}
```

### High-Performance Filter
```go
package main

import (
    "github.com/GrandpaEJ/advancegg"
    "github.com/GrandpaEJ/advancegg/internal/core"
)

func main() {
    dc := advancegg.NewContext(800, 600)
    // ... create your image ...
    
    img := dc.Image()
    
    // Apply optimized filters
    filtered := core.BatchFilter(
        core.FastGrayscale,
        core.FastBrightness(1.2),
        core.FastBlur(3),
    )(img)
    
    // Save result
    imageData := core.NewImageDataFromImage(filtered)
    dc.PutImageData(imageData)
    dc.SavePNG("filtered.png")
}
```

## 📖 Learning Path

### Beginner (Start Here)
1. **[basic-drawing.go](basic-drawing.go)** - Learn fundamental shapes
2. **[text-rendering.go](text-rendering.go)** - Add text to your graphics
3. **[gradients-patterns.go](gradients-patterns.go)** - Beautiful fill styles
4. **[transformations.go](transformations.go)** - Move and rotate objects

### Intermediate
1. **[paths-and-curves.go](paths-and-curves.go)** - Complex path operations
2. **[image-manipulation.go](image-manipulation.go)** - Work with images
3. **[bar-charts.go](bar-charts.go)** - Create data visualizations
4. **[text-on-path-test.go](text-on-path-test.go)** - Advanced typography

### Advanced
1. **[filter-performance-test.go](filter-performance-test.go)** - High-performance processing
2. **[layer-compositing.go](layer-compositing.go)** - Multi-layer graphics
3. **[generative-art.go](generative-art.go)** - Creative programming
4. **[custom-filters.go](custom-filters.go)** - Build your own filters

## 🔧 Running Examples

### Individual Examples
```bash
# Run a specific example
go run basic-drawing.go

# With custom output
go run text-rendering.go -output custom-name.png

# With verbose logging
go run filter-performance-test.go -verbose
```

### Batch Generation
```bash
# Generate all examples
go run generate-all.go

# Generate specific category
go run generate-all.go -category visualization

# Generate with custom settings
go run generate-all.go -size 1920x1080 -quality high
```

## 📊 Performance Benchmarks

Many examples include performance measurements:

```bash
# Run performance tests
go run filter-performance-test.go
go run parallel-processing.go
go run memory-optimization.go

# Generate performance report
go run benchmark-all.go > performance-report.txt
```

## 🎨 Output Gallery

All examples generate high-quality output images. Check the generated files:

- **PNG files** - Lossless graphics, perfect for screenshots
- **JPEG files** - Compressed photos, smaller file sizes
- **PDF files** - Vector graphics, scalable output
- **SVG files** - Web-compatible vector graphics

## 🤝 Contributing Examples

We welcome new examples! Here's how to contribute:

### Example Template
```go
package main

import (
    "fmt"
    "github.com/GrandpaEJ/advancegg"
)

func main() {
    fmt.Println("Example: [Your Example Name]")
    
    // Create context
    dc := advancegg.NewContext(800, 600)
    
    // Your code here
    
    // Save result
    dc.SavePNG("your-example.png")
    fmt.Println("Generated: your-example.png")
}
```

### Guidelines
1. **Clear comments** - Explain what each section does
2. **Self-contained** - Include all necessary code
3. **Error handling** - Handle potential errors gracefully
4. **Performance notes** - Mention any performance considerations
5. **Output description** - Describe what the example generates

### Submission Process
1. Create your example file
2. Test it thoroughly
3. Add it to this README
4. Submit a pull request

## 📚 Additional Resources

- **[API Documentation](../docs/API.md)** - Complete API reference
- **[Performance Guide](../docs/PERFORMANCE.md)** - Optimization tips
- **[Troubleshooting](../docs/TROUBLESHOOTING.md)** - Common issues and solutions
- **[Contributing Guide](../CONTRIBUTING.md)** - How to contribute

## 🎯 Next Steps

After exploring these examples:

1. **Build your own project** using the techniques learned
2. **Contribute examples** to help other developers
3. **Join the community** on GitHub Discussions
4. **Share your creations** with the AdvanceGG community

Happy coding with AdvanceGG! 🚀
