package core

import (
	"math"
)

// Hit testing functionality for interactive graphics

// HitTestable interface for objects that can be hit tested
type HitTestable interface {
	HitTest(x, y float64) bool
	GetBounds() (float64, float64, float64, float64)
}

// HitTestManager manages hit testing for multiple objects
type HitTestManager struct {
	objects []HitTestable
}

// NewHitTestManager creates a new hit test manager
func NewHitTestManager() *HitTestManager {
	return &HitTestManager{
		objects: make([]HitTestable, 0),
	}
}

// AddObject adds an object for hit testing
func (htm *HitTestManager) AddObject(obj HitTestable) {
	htm.objects = append(htm.objects, obj)
}

// RemoveObject removes an object from hit testing
func (htm *HitTestManager) RemoveObject(obj HitTestable) {
	for i, o := range htm.objects {
		if o == obj {
			htm.objects = append(htm.objects[:i], htm.objects[i+1:]...)
			break
		}
	}
}

// HitTest performs hit testing at the given point
func (htm *HitTestManager) HitTest(x, y float64) []HitTestable {
	var hits []HitTestable
	
	for _, obj := range htm.objects {
		if obj.HitTest(x, y) {
			hits = append(hits, obj)
		}
	}
	
	return hits
}

// HitTestFirst returns the first object hit at the given point
func (htm *HitTestManager) HitTestFirst(x, y float64) HitTestable {
	for _, obj := range htm.objects {
		if obj.HitTest(x, y) {
			return obj
		}
	}
	return nil
}

// HitTestLast returns the last (topmost) object hit at the given point
func (htm *HitTestManager) HitTestLast(x, y float64) HitTestable {
	for i := len(htm.objects) - 1; i >= 0; i-- {
		if htm.objects[i].HitTest(x, y) {
			return htm.objects[i]
		}
	}
	return nil
}

// Hit testable shape implementations

// HitTestRect represents a hit testable rectangle
type HitTestRect struct {
	X, Y, Width, Height float64
	ID                  string
}

// HitTest tests if point is inside rectangle
func (r *HitTestRect) HitTest(x, y float64) bool {
	return x >= r.X && x <= r.X+r.Width && y >= r.Y && y <= r.Y+r.Height
}

// GetBounds returns the bounding box
func (r *HitTestRect) GetBounds() (float64, float64, float64, float64) {
	return r.X, r.Y, r.X + r.Width, r.Y + r.Height
}

// HitTestCircle represents a hit testable circle
type HitTestCircle struct {
	X, Y, Radius float64
	ID           string
}

// HitTest tests if point is inside circle
func (c *HitTestCircle) HitTest(x, y float64) bool {
	dx := x - c.X
	dy := y - c.Y
	return dx*dx+dy*dy <= c.Radius*c.Radius
}

// GetBounds returns the bounding box
func (c *HitTestCircle) GetBounds() (float64, float64, float64, float64) {
	return c.X - c.Radius, c.Y - c.Radius, c.X + c.Radius, c.Y + c.Radius
}

// HitTestEllipse represents a hit testable ellipse
type HitTestEllipse struct {
	X, Y, RadiusX, RadiusY float64
	ID                     string
}

// HitTest tests if point is inside ellipse
func (e *HitTestEllipse) HitTest(x, y float64) bool {
	dx := x - e.X
	dy := y - e.Y
	return (dx*dx)/(e.RadiusX*e.RadiusX)+(dy*dy)/(e.RadiusY*e.RadiusY) <= 1
}

// GetBounds returns the bounding box
func (e *HitTestEllipse) GetBounds() (float64, float64, float64, float64) {
	return e.X - e.RadiusX, e.Y - e.RadiusY, e.X + e.RadiusX, e.Y + e.RadiusY
}

// HitTestPolygon represents a hit testable polygon
type HitTestPolygon struct {
	Points []Point
	ID     string
}

// HitTest tests if point is inside polygon using ray casting algorithm
func (p *HitTestPolygon) HitTest(x, y float64) bool {
	if len(p.Points) < 3 {
		return false
	}
	
	inside := false
	j := len(p.Points) - 1
	
	for i := 0; i < len(p.Points); i++ {
		xi, yi := p.Points[i].X, p.Points[i].Y
		xj, yj := p.Points[j].X, p.Points[j].Y
		
		if ((yi > y) != (yj > y)) && (x < (xj-xi)*(y-yi)/(yj-yi)+xi) {
			inside = !inside
		}
		j = i
	}
	
	return inside
}

// GetBounds returns the bounding box
func (p *HitTestPolygon) GetBounds() (float64, float64, float64, float64) {
	if len(p.Points) == 0 {
		return 0, 0, 0, 0
	}
	
	minX, minY := p.Points[0].X, p.Points[0].Y
	maxX, maxY := p.Points[0].X, p.Points[0].Y
	
	for _, point := range p.Points {
		if point.X < minX {
			minX = point.X
		}
		if point.X > maxX {
			maxX = point.X
		}
		if point.Y < minY {
			minY = point.Y
		}
		if point.Y > maxY {
			maxY = point.Y
		}
	}
	
	return minX, minY, maxX, maxY
}

// HitTestLine represents a hit testable line with thickness
type HitTestLine struct {
	X1, Y1, X2, Y2 float64
	Thickness      float64
	ID             string
}

// HitTest tests if point is near the line within thickness
func (l *HitTestLine) HitTest(x, y float64) bool {
	// Calculate distance from point to line segment
	distance := l.distanceToLineSegment(x, y)
	return distance <= l.Thickness/2
}

// distanceToLineSegment calculates distance from point to line segment
func (l *HitTestLine) distanceToLineSegment(x, y float64) float64 {
	dx := l.X2 - l.X1
	dy := l.Y2 - l.Y1
	
	if dx == 0 && dy == 0 {
		// Line is a point
		return math.Sqrt((x-l.X1)*(x-l.X1) + (y-l.Y1)*(y-l.Y1))
	}
	
	// Calculate parameter t for closest point on line
	t := ((x-l.X1)*dx + (y-l.Y1)*dy) / (dx*dx + dy*dy)
	
	// Clamp t to line segment
	if t < 0 {
		t = 0
	} else if t > 1 {
		t = 1
	}
	
	// Calculate closest point on line segment
	closestX := l.X1 + t*dx
	closestY := l.Y1 + t*dy
	
	// Return distance to closest point
	return math.Sqrt((x-closestX)*(x-closestX) + (y-closestY)*(y-closestY))
}

// GetBounds returns the bounding box
func (l *HitTestLine) GetBounds() (float64, float64, float64, float64) {
	minX := l.X1
	if l.X2 < minX {
		minX = l.X2
	}
	maxX := l.X1
	if l.X2 > maxX {
		maxX = l.X2
	}
	minY := l.Y1
	if l.Y2 < minY {
		minY = l.Y2
	}
	maxY := l.Y1
	if l.Y2 > maxY {
		maxY = l.Y2
	}
	
	// Expand bounds by thickness
	halfThickness := l.Thickness / 2
	return minX - halfThickness, minY - halfThickness, 
		   maxX + halfThickness, maxY + halfThickness
}

// HitTestPath represents a hit testable path
type HitTestPath struct {
	Points []Point
	Closed bool
	ID     string
}

// HitTest tests if point is inside path (if closed) or near path (if open)
func (p *HitTestPath) HitTest(x, y float64) bool {
	if p.Closed {
		// Use polygon hit test for closed paths
		polygon := &HitTestPolygon{Points: p.Points, ID: p.ID}
		return polygon.HitTest(x, y)
	} else {
		// Test distance to path segments for open paths
		tolerance := 5.0 // pixels
		
		for i := 0; i < len(p.Points)-1; i++ {
			line := &HitTestLine{
				X1: p.Points[i].X, Y1: p.Points[i].Y,
				X2: p.Points[i+1].X, Y2: p.Points[i+1].Y,
				Thickness: tolerance * 2,
			}
			if line.HitTest(x, y) {
				return true
			}
		}
		return false
	}
}

// GetBounds returns the bounding box
func (p *HitTestPath) GetBounds() (float64, float64, float64, float64) {
	if len(p.Points) == 0 {
		return 0, 0, 0, 0
	}
	
	minX, minY := p.Points[0].X, p.Points[0].Y
	maxX, maxY := p.Points[0].X, p.Points[0].Y
	
	for _, point := range p.Points {
		if point.X < minX {
			minX = point.X
		}
		if point.X > maxX {
			maxX = point.X
		}
		if point.Y < minY {
			minY = point.Y
		}
		if point.Y > maxY {
			maxY = point.Y
		}
	}
	
	return minX, minY, maxX, maxY
}

// Convenience functions

// CreateHitTestRect creates a hit testable rectangle
func CreateHitTestRect(id string, x, y, width, height float64) *HitTestRect {
	return &HitTestRect{X: x, Y: y, Width: width, Height: height, ID: id}
}

// CreateHitTestCircle creates a hit testable circle
func CreateHitTestCircle(id string, x, y, radius float64) *HitTestCircle {
	return &HitTestCircle{X: x, Y: y, Radius: radius, ID: id}
}

// CreateHitTestEllipse creates a hit testable ellipse
func CreateHitTestEllipse(id string, x, y, radiusX, radiusY float64) *HitTestEllipse {
	return &HitTestEllipse{X: x, Y: y, RadiusX: radiusX, RadiusY: radiusY, ID: id}
}

// CreateHitTestPolygon creates a hit testable polygon
func CreateHitTestPolygon(id string, points []Point) *HitTestPolygon {
	return &HitTestPolygon{Points: points, ID: id}
}

// CreateHitTestLine creates a hit testable line
func CreateHitTestLine(id string, x1, y1, x2, y2, thickness float64) *HitTestLine {
	return &HitTestLine{X1: x1, Y1: y1, X2: x2, Y2: y2, Thickness: thickness, ID: id}
}

// CreateHitTestPath creates a hit testable path
func CreateHitTestPath(id string, points []Point, closed bool) *HitTestPath {
	return &HitTestPath{Points: points, Closed: closed, ID: id}
}
