# AdvanceGG Python Integration

Use AdvanceGG from Python through Go's C-shared library interface.

## Installation

### Prerequisites

- Python 3.7+
- Go 1.18+
- GCC or compatible C compiler

### Build the Shared Library

```bash
# Clone AdvanceGG
git clone https://github.com/GrandpaEJ/advancegg.git
cd advancegg

# Build Python shared library
cd ecosystem/python
go build -buildmode=c-shared -o advancegg.so advancegg_python.go

# Install Python package
pip install -e .
```

## Quick Start

```python
import advancegg

# Create a new canvas
canvas = advancegg.Canvas(800, 600)

# Set background
canvas.set_rgb(0.1, 0.1, 0.3)
canvas.clear()

# Draw a red circle
canvas.set_rgb(1, 0, 0)
canvas.draw_circle(400, 300, 100)
canvas.fill()

# Add text
canvas.set_rgb(1, 1, 1)
canvas.draw_string("Hello from Python!", 300, 350)

# Save as PNG
canvas.save_png("output.png")
```

## Features

### Basic Drawing

```python
import advancegg

canvas = advancegg.Canvas(800, 600)

# Shapes
canvas.draw_rectangle(100, 100, 200, 150)
canvas.draw_circle(400, 300, 80)
canvas.draw_ellipse(600, 200, 100, 60)

# Lines and paths
canvas.move_to(100, 400)
canvas.line_to(200, 350)
canvas.line_to(300, 400)
canvas.close_path()
canvas.stroke()
```

### Advanced Features

```python
# Gradients
gradient = canvas.create_linear_gradient(0, 0, 200, 0)
gradient.add_color_stop(0, (1, 0, 0, 1))  # Red
gradient.add_color_stop(1, (0, 0, 1, 1))  # Blue
canvas.set_fill_style(gradient)

# Text effects
canvas.load_font("arial.ttf", 24)
canvas.draw_text_on_circle("Circular Text", 400, 300, 150)

# Image processing
image = canvas.load_image("input.jpg")
blurred = canvas.apply_blur(image, 5.0)
canvas.draw_image(blurred, 0, 0)

# Layers
layer_manager = advancegg.LayerManager(800, 600)
bg_layer = layer_manager.add_layer("background")
fg_layer = layer_manager.add_layer("foreground")

bg_layer.set_rgb(0.2, 0.2, 0.4)
bg_layer.clear()

fg_layer.set_rgb(1, 1, 0)
fg_layer.draw_circle(400, 300, 100)
fg_layer.fill()

result = layer_manager.flatten()
result.save_png("layered.png")
```

### Data Visualization

```python
import advancegg
import numpy as np

def create_chart(data, labels):
    canvas = advancegg.Canvas(800, 600)
    
    # Background gradient
    bg = canvas.create_linear_gradient(0, 0, 0, 600)
    bg.add_color_stop(0, (0.95, 0.95, 0.98, 1))
    bg.add_color_stop(1, (0.9, 0.9, 0.95, 1))
    canvas.set_fill_style(bg)
    canvas.draw_rectangle(0, 0, 800, 600)
    canvas.fill()
    
    # Draw bars
    colors = [(0.2, 0.4, 0.8, 1), (0.8, 0.2, 0.4, 1), (0.2, 0.8, 0.2, 1)]
    bar_width = 60
    spacing = 100
    
    for i, (value, label) in enumerate(zip(data, labels)):
        x = 100 + i * spacing
        height = value * 4
        
        # Bar with gradient
        bar_gradient = canvas.create_linear_gradient(x, 500-height, x, 500)
        color = colors[i % len(colors)]
        bar_gradient.add_color_stop(0, color)
        bar_gradient.add_color_stop(1, (color[0]*0.7, color[1]*0.7, color[2]*0.7, 1))
        
        canvas.set_fill_style(bar_gradient)
        canvas.draw_rounded_rectangle(x, 500-height, bar_width, height, 8)
        canvas.fill()
        
        # Value label
        canvas.set_rgb(0.2, 0.2, 0.2)
        canvas.draw_string_anchored(f"{value:.1f}", x + bar_width/2, 480-height, 0.5, 0.5)
        
        # Category label
        canvas.draw_string_anchored(label, x + bar_width/2, 520, 0.5, 0.5)
    
    return canvas

# Usage
data = [85.5, 92.3, 78.1, 96.7, 88.9]
labels = ["Q1", "Q2", "Q3", "Q4", "Q5"]
chart = create_chart(data, labels)
chart.save_png("python_chart.png")
```

### Scientific Plotting

```python
import advancegg
import math

def plot_function(func, x_range, y_range, canvas_size=(800, 600)):
    canvas = advancegg.Canvas(*canvas_size)
    
    # White background
    canvas.set_rgb(1, 1, 1)
    canvas.clear()
    
    # Draw axes
    canvas.set_rgb(0, 0, 0)
    canvas.set_line_width(2)
    
    # X-axis
    canvas.draw_line(50, canvas_size[1]-50, canvas_size[0]-50, canvas_size[1]-50)
    # Y-axis
    canvas.draw_line(50, 50, 50, canvas_size[1]-50)
    canvas.stroke()
    
    # Plot function
    canvas.set_rgb(0.2, 0.4, 0.8)
    canvas.set_line_width(3)
    
    x_min, x_max = x_range
    y_min, y_max = y_range
    
    first_point = True
    for i in range(canvas_size[0] - 100):
        x = x_min + (x_max - x_min) * i / (canvas_size[0] - 100)
        y = func(x)
        
        # Convert to canvas coordinates
        canvas_x = 50 + i
        canvas_y = canvas_size[1] - 50 - (y - y_min) / (y_max - y_min) * (canvas_size[1] - 100)
        
        if first_point:
            canvas.move_to(canvas_x, canvas_y)
            first_point = False
        else:
            canvas.line_to(canvas_x, canvas_y)
    
    canvas.stroke()
    return canvas

# Plot sine wave
sine_plot = plot_function(math.sin, (-2*math.pi, 2*math.pi), (-1.5, 1.5))
sine_plot.save_png("sine_wave.png")

# Plot polynomial
poly_plot = plot_function(lambda x: x**3 - 2*x**2 + x, (-2, 3), (-5, 5))
poly_plot.save_png("polynomial.png")
```

## API Reference

### Canvas Class

```python
class Canvas:
    def __init__(self, width: int, height: int)
    
    # Basic drawing
    def draw_rectangle(self, x: float, y: float, width: float, height: float)
    def draw_circle(self, x: float, y: float, radius: float)
    def draw_line(self, x1: float, y1: float, x2: float, y2: float)
    
    # Path operations
    def move_to(self, x: float, y: float)
    def line_to(self, x: float, y: float)
    def curve_to(self, cp1x: float, cp1y: float, cp2x: float, cp2y: float, x: float, y: float)
    def close_path(self)
    
    # Fill and stroke
    def fill(self)
    def stroke(self)
    def clear(self)
    
    # Colors
    def set_rgb(self, r: float, g: float, b: float)
    def set_rgba(self, r: float, g: float, b: float, a: float)
    def set_hex_color(self, hex_color: str)
    
    # Text
    def load_font(self, path: str, size: float)
    def draw_string(self, text: str, x: float, y: float)
    def draw_string_anchored(self, text: str, x: float, y: float, ax: float, ay: float)
    
    # Images
    def load_image(self, path: str) -> Image
    def draw_image(self, image: Image, x: float, y: float)
    def save_png(self, path: str)
    def save_jpeg(self, path: str, quality: int = 95)
    
    # Advanced features
    def create_linear_gradient(self, x1: float, y1: float, x2: float, y2: float) -> Gradient
    def apply_blur(self, image: Image, radius: float) -> Image
    def draw_text_on_circle(self, text: str, x: float, y: float, radius: float)
```

### LayerManager Class

```python
class LayerManager:
    def __init__(self, width: int, height: int)
    def add_layer(self, name: str) -> Canvas
    def remove_layer(self, name: str)
    def set_layer_opacity(self, name: str, opacity: float)
    def set_layer_blend_mode(self, name: str, mode: str)
    def flatten(self) -> Canvas
```

## Examples

See the `examples/` directory for complete Python examples:

- `basic_drawing.py` - Basic shapes and colors
- `data_visualization.py` - Charts and graphs
- `image_processing.py` - Filters and effects
- `scientific_plotting.py` - Mathematical visualizations
- `game_graphics.py` - Sprites and animations
- `web_graphics.py` - Server-side image generation

## Performance Tips

1. **Reuse Canvas objects** when possible
2. **Use layers** for complex compositions
3. **Enable caching** for repeated operations
4. **Batch drawing operations** for better performance
5. **Use appropriate image formats** (PNG for graphics, JPEG for photos)

## Troubleshooting

### Common Issues

**ImportError: cannot import name 'advancegg'**
- Ensure the shared library is built and in the correct location
- Check that all dependencies are installed

**Font loading errors**
- Use absolute paths for font files
- Ensure font files are accessible and valid TTF/OTF format

**Memory issues with large images**
- Enable memory pooling: `canvas.set_memory_pooling(True)`
- Process images in smaller chunks for very large files

## Contributing

Contributions to the Python integration are welcome! Please see the main AdvanceGG repository for contribution guidelines.

## License

Same as AdvanceGG - MIT License
