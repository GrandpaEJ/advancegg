# API Reference

## Context Creation

### NewContext
```go
func NewContext(width, height int) *Context
```
Creates a new drawing context with the specified dimensions.

### NewContextForImage
```go
func NewContextForImage(im image.Image) *Context
```
Creates a new context for an existing image.

### NewContextForRGBA
```go
func NewContextForRGBA(im *image.RGBA) *Context
```
Creates a new context for an RGBA image.

## Drawing Functions

### Basic Shapes

#### DrawPoint
```go
func (dc *Context) DrawPoint(x, y, r float64)
```
Draws a point (small circle) at the specified coordinates.

#### DrawLine
```go
func (dc *Context) DrawLine(x1, y1, x2, y2 float64)
```
Draws a line between two points.

#### DrawRectangle
```go
func (dc *Context) DrawRectangle(x, y, w, h float64)
```
Draws a rectangle.

#### DrawRoundedRectangle
```go
func (dc *Context) DrawRoundedRectangle(x, y, w, h, r float64)
```
Draws a rectangle with rounded corners.

#### DrawCircle
```go
func (dc *Context) DrawCircle(x, y, r float64)
```
Draws a circle.

#### DrawEllipse
```go
func (dc *Context) DrawEllipse(x, y, rx, ry float64)
```
Draws an ellipse.

### Path Functions

#### MoveTo
```go
func (dc *Context) MoveTo(x, y float64)
```
Moves the current point to the specified coordinates.

#### LineTo
```go
func (dc *Context) LineTo(x, y float64)
```
Draws a line from the current point to the specified coordinates.

#### QuadraticTo
```go
func (dc *Context) QuadraticTo(x1, y1, x2, y2 float64)
```
Draws a quadratic Bézier curve.

#### CubicTo
```go
func (dc *Context) CubicTo(x1, y1, x2, y2, x3, y3 float64)
```
Draws a cubic Bézier curve.

## Color Functions

### SetRGB
```go
func (dc *Context) SetRGB(r, g, b float64)
```
Sets the color using RGB values (0.0 to 1.0).

### SetRGBA
```go
func (dc *Context) SetRGBA(r, g, b, a float64)
```
Sets the color using RGBA values (0.0 to 1.0).

### SetRGB255
```go
func (dc *Context) SetRGB255(r, g, b int)
```
Sets the color using RGB values (0 to 255).

### SetHexColor
```go
func (dc *Context) SetHexColor(x string)
```
Sets the color using a hex color string (e.g., "#FF0000").

## Text Functions

### DrawString
```go
func (dc *Context) DrawString(s string, x, y float64)
```
Draws text at the specified coordinates.

### DrawStringAnchored
```go
func (dc *Context) DrawStringAnchored(s string, x, y, ax, ay float64)
```
Draws text with specified anchor point.

### MeasureString
```go
func (dc *Context) MeasureString(s string) (w, h float64)
```
Measures the dimensions of text.

## Rendering Functions

### Fill
```go
func (dc *Context) Fill()
```
Fills the current path with the current color.

### Stroke
```go
func (dc *Context) Stroke()
```
Strokes the current path with the current color.

### Clear
```go
func (dc *Context) Clear()
```
Clears the entire context.

## Transformation Functions

### Translate
```go
func (dc *Context) Translate(x, y float64)
```
Translates the coordinate system.

### Scale
```go
func (dc *Context) Scale(x, y float64)
```
Scales the coordinate system.

### Rotate
```go
func (dc *Context) Rotate(angle float64)
```
Rotates the coordinate system.

## File I/O

### SavePNG
```go
func (dc *Context) SavePNG(path string) error
```
Saves the context as a PNG file.

### LoadImage
```go
func LoadImage(path string) (image.Image, error)
```
Loads an image from a file.
