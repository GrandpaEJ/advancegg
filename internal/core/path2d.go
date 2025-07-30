package core

import (
	"math"

	"github.com/golang/freetype/raster"
	"golang.org/x/image/math/fixed"
)

// Path2D represents a 2D path that can be reused and manipulated
type Path2D struct {
	path     raster.Path
	currentX float64
	currentY float64
	startX   float64
	startY   float64
	hasStart bool
}

// NewPath2D creates a new empty Path2D
func NewPath2D() *Path2D {
	return &Path2D{
		path: make(raster.Path, 0),
	}
}

// NewPath2DFromPath creates a new Path2D from an existing Path2D
func NewPath2DFromPath(other *Path2D) *Path2D {
	if other == nil {
		return NewPath2D()
	}

	newPath := &Path2D{
		path:     make(raster.Path, len(other.path)),
		currentX: other.currentX,
		currentY: other.currentY,
		startX:   other.startX,
		startY:   other.startY,
		hasStart: other.hasStart,
	}
	copy(newPath.path, other.path)
	return newPath
}

// MoveTo moves the current point to the specified coordinates
func (p *Path2D) MoveTo(x, y float64) {
	p.path = append(p.path, 0) // raster.MoveTo
	p.path = append(p.path, fixedFromFloat(x), fixedFromFloat(y))
	p.currentX = x
	p.currentY = y
	if !p.hasStart {
		p.startX = x
		p.startY = y
		p.hasStart = true
	}
}

// LineTo draws a line from the current point to the specified coordinates
func (p *Path2D) LineTo(x, y float64) {
	if !p.hasStart {
		p.MoveTo(0, 0)
	}
	p.path = append(p.path, 1) // raster.LineTo
	p.path = append(p.path, fixedFromFloat(x), fixedFromFloat(y))
	p.currentX = x
	p.currentY = y
}

// QuadraticCurveTo draws a quadratic Bézier curve
func (p *Path2D) QuadraticCurveTo(cpx, cpy, x, y float64) {
	if !p.hasStart {
		p.MoveTo(0, 0)
	}
	p.path = append(p.path, 2) // raster.QuadTo
	p.path = append(p.path, fixedFromFloat(cpx), fixedFromFloat(cpy), fixedFromFloat(x), fixedFromFloat(y))
	p.currentX = x
	p.currentY = y
}

// BezierCurveTo draws a cubic Bézier curve
func (p *Path2D) BezierCurveTo(cp1x, cp1y, cp2x, cp2y, x, y float64) {
	if !p.hasStart {
		p.MoveTo(0, 0)
	}
	p.path = append(p.path, 3) // raster.CubeTo
	p.path = append(p.path, fixedFromFloat(cp1x), fixedFromFloat(cp1y), fixedFromFloat(cp2x), fixedFromFloat(cp2y), fixedFromFloat(x), fixedFromFloat(y))
	p.currentX = x
	p.currentY = y
}

// ClosePath closes the current path by drawing a line to the start point
func (p *Path2D) ClosePath() {
	if p.hasStart {
		p.LineTo(p.startX, p.startY)
	}
}

// Arc adds an arc to the path
func (p *Path2D) Arc(x, y, radius, startAngle, endAngle float64, counterclockwise bool) {
	if radius < 0 {
		return
	}

	// Normalize angles
	if counterclockwise {
		for endAngle > startAngle {
			endAngle -= 2 * math.Pi
		}
	} else {
		for endAngle < startAngle {
			endAngle += 2 * math.Pi
		}
	}

	// Calculate start point
	startX := x + radius*math.Cos(startAngle)
	startY := y + radius*math.Sin(startAngle)

	// Move to start point if this is the first operation
	if !p.hasStart {
		p.MoveTo(startX, startY)
	} else {
		p.LineTo(startX, startY)
	}

	// Draw the arc using multiple cubic Bézier curves
	p.drawArc(x, y, radius, startAngle, endAngle, counterclockwise)
}

// ArcTo adds an arc with the given control points and radius
func (p *Path2D) ArcTo(x1, y1, x2, y2, radius float64) {
	if !p.hasStart {
		p.MoveTo(0, 0)
	}

	x0, y0 := p.currentX, p.currentY

	// Calculate the arc
	if radius <= 0 || (x0 == x1 && y0 == y1) || (x1 == x2 && y1 == y2) {
		p.LineTo(x1, y1)
		return
	}

	// Vector from current point to first control point
	dx0 := x0 - x1
	dy0 := y0 - y1

	// Vector from first control point to second control point
	dx1 := x2 - x1
	dy1 := y2 - y1

	// Calculate lengths
	len0 := math.Sqrt(dx0*dx0 + dy0*dy0)
	len1 := math.Sqrt(dx1*dx1 + dy1*dy1)

	if len0 == 0 || len1 == 0 {
		p.LineTo(x1, y1)
		return
	}

	// Normalize vectors
	dx0 /= len0
	dy0 /= len0
	dx1 /= len1
	dy1 /= len1

	// Calculate angle between vectors
	dot := dx0*dx1 + dy0*dy1
	if dot >= 1 {
		p.LineTo(x1, y1)
		return
	}

	// Calculate tangent length
	halfAngle := math.Acos(dot) / 2
	tangentLength := radius / math.Tan(halfAngle)

	// Calculate tangent points
	t0 := math.Min(tangentLength, len0)
	t1 := math.Min(tangentLength, len1)

	startX := x1 + dx0*t0
	startY := y1 + dy0*t0
	endX := x1 + dx1*t1
	endY := y1 + dy1*t1

	// Line to start of arc
	p.LineTo(startX, startY)

	// Calculate center of arc
	centerX := startX + dy0*radius
	centerY := startY - dx0*radius

	// Calculate angles
	startAngle := math.Atan2(startY-centerY, startX-centerX)
	endAngle := math.Atan2(endY-centerY, endX-centerX)

	// Draw the arc
	p.drawArc(centerX, centerY, radius, startAngle, endAngle, false)
}

// Rect adds a rectangle to the path
func (p *Path2D) Rect(x, y, width, height float64) {
	p.MoveTo(x, y)
	p.LineTo(x+width, y)
	p.LineTo(x+width, y+height)
	p.LineTo(x, y+height)
	p.ClosePath()
}

// Ellipse adds an ellipse to the path
func (p *Path2D) Ellipse(x, y, radiusX, radiusY, rotation, startAngle, endAngle float64, counterclockwise bool) {
	if radiusX < 0 || radiusY < 0 {
		return
	}

	// Save current transformation state
	cos := math.Cos(rotation)
	sin := math.Sin(rotation)

	// Normalize angles
	if counterclockwise {
		for endAngle > startAngle {
			endAngle -= 2 * math.Pi
		}
	} else {
		for endAngle < startAngle {
			endAngle += 2 * math.Pi
		}
	}

	// Calculate start point
	startX := x + radiusX*math.Cos(startAngle)*cos - radiusY*math.Sin(startAngle)*sin
	startY := y + radiusX*math.Cos(startAngle)*sin + radiusY*math.Sin(startAngle)*cos

	// Move to start point
	if !p.hasStart {
		p.MoveTo(startX, startY)
	} else {
		p.LineTo(startX, startY)
	}

	// Draw the ellipse using multiple cubic Bézier curves
	p.drawEllipse(x, y, radiusX, radiusY, rotation, startAngle, endAngle, counterclockwise)
}

// AddPath adds another path to this path
func (p *Path2D) AddPath(other *Path2D) {
	if other == nil {
		return
	}
	p.path = append(p.path, other.path...)
	p.currentX = other.currentX
	p.currentY = other.currentY
}

// GetPath returns the internal raster.Path for use with Context
func (p *Path2D) GetPath() raster.Path {
	return p.path
}

// IsEmpty returns true if the path has no operations
func (p *Path2D) IsEmpty() bool {
	return len(p.path) == 0
}

// GetCurrentPoint returns the current point of the path
func (p *Path2D) GetCurrentPoint() (x, y float64) {
	return p.currentX, p.currentY
}

// Helper function to convert float64 to fixed.Int26_6
func fixedFromFloat(x float64) fixed.Int26_6 {
	return fixed.Int26_6(x * 64)
}

// Helper function to draw an arc using cubic Bézier curves
func (p *Path2D) drawArc(cx, cy, radius, startAngle, endAngle float64, counterclockwise bool) {
	// Implementation of arc drawing using cubic Bézier approximation
	// This is a simplified version - a full implementation would use multiple curves for better accuracy

	angleDiff := endAngle - startAngle
	if counterclockwise {
		angleDiff = -angleDiff
	}

	// Use multiple curves for better approximation
	segments := int(math.Ceil(math.Abs(angleDiff) / (math.Pi / 2)))
	if segments == 0 {
		segments = 1
	}

	segmentAngle := angleDiff / float64(segments)

	for i := 0; i < segments; i++ {
		a1 := startAngle + float64(i)*segmentAngle
		a2 := startAngle + float64(i+1)*segmentAngle

		// Calculate control points for cubic Bézier approximation of arc segment
		alpha := math.Sin(a2-a1) * (math.Sqrt(4+3*math.Tan((a2-a1)/2)*math.Tan((a2-a1)/2)) - 1) / 3

		x1 := cx + radius*math.Cos(a1)
		y1 := cy + radius*math.Sin(a1)
		x2 := cx + radius*math.Cos(a2)
		y2 := cy + radius*math.Sin(a2)

		cp1x := x1 - alpha*radius*math.Sin(a1)
		cp1y := y1 + alpha*radius*math.Cos(a1)
		cp2x := x2 + alpha*radius*math.Sin(a2)
		cp2y := y2 - alpha*radius*math.Cos(a2)

		p.BezierCurveTo(cp1x, cp1y, cp2x, cp2y, x2, y2)
	}
}

// Helper function to draw an ellipse using cubic Bézier curves
func (p *Path2D) drawEllipse(cx, cy, rx, ry, rotation, startAngle, endAngle float64, counterclockwise bool) {
	// Similar to drawArc but for ellipses
	// This is a simplified implementation
	angleDiff := endAngle - startAngle
	if counterclockwise {
		angleDiff = -angleDiff
	}

	segments := int(math.Ceil(math.Abs(angleDiff) / (math.Pi / 2)))
	if segments == 0 {
		segments = 1
	}

	segmentAngle := angleDiff / float64(segments)
	cos := math.Cos(rotation)
	sin := math.Sin(rotation)

	for i := 0; i < segments; i++ {
		a1 := startAngle + float64(i)*segmentAngle
		a2 := startAngle + float64(i+1)*segmentAngle

		// Calculate points on ellipse
		ex1 := rx * math.Cos(a1)
		ey1 := ry * math.Sin(a1)
		ex2 := rx * math.Cos(a2)
		ey2 := ry * math.Sin(a2)

		// Apply rotation
		x1 := cx + ex1*cos - ey1*sin
		y1 := cy + ex1*sin + ey1*cos
		x2 := cx + ex2*cos - ey2*sin
		y2 := cy + ex2*sin + ey2*cos

		// Calculate control points (simplified)
		alpha := math.Tan((a2-a1)/2) * 4 / 3

		cp1x := x1 - alpha*ry*math.Sin(a1)*cos - alpha*rx*math.Cos(a1)*sin
		cp1y := y1 - alpha*ry*math.Sin(a1)*sin + alpha*rx*math.Cos(a1)*cos
		cp2x := x2 + alpha*ry*math.Sin(a2)*cos + alpha*rx*math.Cos(a2)*sin
		cp2y := y2 + alpha*ry*math.Sin(a2)*sin - alpha*rx*math.Cos(a2)*cos

		p.BezierCurveTo(cp1x, cp1y, cp2x, cp2y, x2, y2)
	}
}
