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

## Path2D

Path2D provides advanced path manipulation capabilities, allowing you to create reusable path objects that can be drawn multiple times.

### NewPath2D
```go
func NewPath2D() *Path2D
```
Creates a new empty Path2D object.

### NewPath2DFromPath
```go
func NewPath2DFromPath(other *Path2D) *Path2D
```
Creates a new Path2D object by copying an existing one.

### Path2D Methods

#### MoveTo
```go
func (p *Path2D) MoveTo(x, y float64)
```
Moves the current point to the specified coordinates.

#### LineTo
```go
func (p *Path2D) LineTo(x, y float64)
```
Draws a line from the current point to the specified coordinates.

#### QuadraticCurveTo
```go
func (p *Path2D) QuadraticCurveTo(cpx, cpy, x, y float64)
```
Draws a quadratic Bézier curve.

#### BezierCurveTo
```go
func (p *Path2D) BezierCurveTo(cp1x, cp1y, cp2x, cp2y, x, y float64)
```
Draws a cubic Bézier curve.

#### Arc
```go
func (p *Path2D) Arc(x, y, radius, startAngle, endAngle float64, counterclockwise bool)
```
Adds an arc to the path.

#### ArcTo
```go
func (p *Path2D) ArcTo(x1, y1, x2, y2, radius float64)
```
Adds an arc with the given control points and radius.

#### Ellipse
```go
func (p *Path2D) Ellipse(x, y, radiusX, radiusY, rotation, startAngle, endAngle float64, counterclockwise bool)
```
Adds an ellipse to the path.

#### Rect
```go
func (p *Path2D) Rect(x, y, width, height float64)
```
Adds a rectangle to the path.

#### ClosePath
```go
func (p *Path2D) ClosePath()
```
Closes the current path by drawing a line to the start point.

#### AddPath
```go
func (p *Path2D) AddPath(other *Path2D)
```
Adds another path to this path.

### Context Methods for Path2D

#### DrawPath2D
```go
func (dc *Context) DrawPath2D(path2d *Path2D)
```
Draws a Path2D object to the context's current path.

#### FillPath2D
```go
func (dc *Context) FillPath2D(path2d *Path2D)
```
Fills a Path2D object with the current fill style.

#### StrokePath2D
```go
func (dc *Context) StrokePath2D(path2d *Path2D)
```
Strokes a Path2D object with the current stroke style.

#### ClipPath2D
```go
func (dc *Context) ClipPath2D(path2d *Path2D)
```
Sets a Path2D object as the clipping region.

#### IsPointInPath2D
```go
func (dc *Context) IsPointInPath2D(path2d *Path2D, x, y float64) bool
```
Tests if a point is inside a Path2D object.

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
