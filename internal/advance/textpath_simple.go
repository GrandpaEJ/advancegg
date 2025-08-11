package advance

import (
	"math"

	"github.com/GrandpaEJ/advancegg/internal/core"
	"github.com/golang/freetype/raster"
	"golang.org/x/image/math/fixed"
)

// Simplified Text-on-Path functionality

// SimpleTextOnPath handles basic text rendering along simple paths
type SimpleTextOnPath struct {
	Text      string
	FontSize  float64
	Spacing   float64
	Offset    float64
	Alignment SimpleTextAlignment
}

// SimpleTextAlignment represents text alignment along path
type SimpleTextAlignment int

const (
	SimpleAlignStart SimpleTextAlignment = iota
	SimpleAlignCenter
	SimpleAlignEnd
)

// NewSimpleTextOnPath creates a new simple text-on-path renderer
func NewSimpleTextOnPath(text string) *SimpleTextOnPath {
	return &SimpleTextOnPath{
		Text:      text,
		FontSize:  12,
		Spacing:   1.0,
		Offset:    0,
		Alignment: SimpleAlignStart,
	}
}

// DrawTextOnCircle draws text along a circular path
func DrawTextOnCircle(dc *core.Context, text string, centerX, centerY, radius float64) {
	if text == "" {
		return
	}

	textOnPath := NewSimpleTextOnPath(text)
	textOnPath.renderOnCircle(dc, centerX, centerY, radius)
}

// DrawTextOnWave draws text along a wave path
func DrawTextOnWave(dc *core.Context, text string, startX, startY, endX, amplitude, frequency float64) {
	if text == "" {
		return
	}

	textOnPath := NewSimpleTextOnPath(text)
	textOnPath.renderOnWave(dc, startX, startY, endX, amplitude, frequency)
}

// DrawTextOnSpiral draws text along a spiral path
func DrawTextOnSpiral(dc *core.Context, text string, centerX, centerY, startRadius, endRadius float64, turns int) {
	if text == "" {
		return
	}

	textOnPath := NewSimpleTextOnPath(text)
	textOnPath.renderOnSpiral(dc, centerX, centerY, startRadius, endRadius, turns)
}

// renderOnCircle renders text along a circular path
func (stp *SimpleTextOnPath) renderOnCircle(dc *core.Context, centerX, centerY, radius float64) {
	textWidth := stp.estimateTextWidth()

	// Calculate starting angle based on alignment
	startAngle := stp.Offset / radius
	switch stp.Alignment {
	case SimpleAlignCenter:
		startAngle -= (textWidth / radius) / 2
	case SimpleAlignEnd:
		startAngle -= textWidth / radius
	}

	currentAngle := startAngle

	for _, r := range stp.Text {
		// Calculate position
		x := centerX + radius*math.Cos(currentAngle)
		y := centerY + radius*math.Sin(currentAngle)

		// Calculate rotation (tangent to circle)
		rotation := currentAngle + math.Pi/2

		// Render character
		stp.renderCharacter(dc, r, x, y, rotation)

		// Advance angle
		charWidth := stp.getCharacterWidth(r)
		angleAdvance := (charWidth * stp.Spacing) / radius
		currentAngle += angleAdvance

		if currentAngle > 2*math.Pi {
			break // Don't wrap around
		}
	}
}

// renderOnWave renders text along a wave path
func (stp *SimpleTextOnPath) renderOnWave(dc *core.Context, startX, startY, endX, amplitude, frequency float64) {
	width := endX - startX
	textWidth := stp.estimateTextWidth()

	// Calculate starting position based on alignment
	startOffset := stp.Offset
	switch stp.Alignment {
	case SimpleAlignCenter:
		startOffset += (width - textWidth) / 2
	case SimpleAlignEnd:
		startOffset += width - textWidth
	}

	currentX := startX + startOffset

	for _, r := range stp.Text {
		if currentX >= endX {
			break
		}

		// Calculate wave position
		t := (currentX - startX) / width
		y := startY + amplitude*math.Sin(frequency*t*2*math.Pi)

		// Calculate tangent (derivative of sine wave)
		tangent := math.Atan(amplitude * frequency * 2 * math.Pi * math.Cos(frequency*t*2*math.Pi) / width)

		// Render character
		stp.renderCharacter(dc, r, currentX, y, tangent)

		// Advance position
		charWidth := stp.getCharacterWidth(r)
		currentX += charWidth * stp.Spacing
	}
}

// renderOnSpiral renders text along a spiral path
func (stp *SimpleTextOnPath) renderOnSpiral(dc *core.Context, centerX, centerY, startRadius, endRadius float64, turns int) {
	totalAngle := float64(turns) * 2 * math.Pi
	radiusDelta := endRadius - startRadius

	currentAngle := 0.0
	charIndex := 0

	for _, r := range stp.Text {
		if currentAngle >= totalAngle {
			break
		}

		// Calculate spiral position
		t := currentAngle / totalAngle
		radius := startRadius + t*radiusDelta

		x := centerX + radius*math.Cos(currentAngle)
		y := centerY + radius*math.Sin(currentAngle)

		// Calculate rotation (tangent to spiral)
		rotation := currentAngle + math.Pi/2

		// Render character
		stp.renderCharacter(dc, r, x, y, rotation)

		// Advance angle
		charWidth := stp.getCharacterWidth(r)
		angleAdvance := (charWidth * stp.Spacing) / radius
		currentAngle += angleAdvance
		charIndex++
	}
}

// renderCharacter renders a single character at the specified position and rotation
func (stp *SimpleTextOnPath) renderCharacter(dc *core.Context, r rune, x, y, rotation float64) {
	// Save current transformation
	dc.Push()

	// Apply transformation
	dc.Translate(x, y)
	dc.Rotate(rotation)

	// Draw character
	char := string(r)
	dc.DrawString(char, 0, 0)

	// Restore transformation
	dc.Pop()
}

// estimateTextWidth estimates the total width of the text
func (stp *SimpleTextOnPath) estimateTextWidth() float64 {
	width := 0.0
	for _, r := range stp.Text {
		width += stp.getCharacterWidth(r) * stp.Spacing
	}
	return width
}

// getCharacterWidth gets the width of a character
func (stp *SimpleTextOnPath) getCharacterWidth(r rune) float64 {
	// Simplified character width calculation
	if r == ' ' {
		return stp.FontSize * 0.3
	}
	return stp.FontSize * 0.6 // Average character width
}

// Predefined path generators for simple paths

// CreateSimpleCircularPath creates parameters for a circular text path
func CreateSimpleCircularPath(centerX, centerY, radius float64) (float64, float64, float64) {
	return centerX, centerY, radius
}

// CreateSimpleWavePath creates parameters for a wave text path
func CreateSimpleWavePath(startX, startY, endX, amplitude, frequency float64) (float64, float64, float64, float64, float64) {
	return startX, startY, endX, amplitude, frequency
}

// CreateSimpleSpiralPath creates parameters for a spiral text path
func CreateSimpleSpiralPath(centerX, centerY, startRadius, endRadius float64, turns int) (float64, float64, float64, float64, int) {
	return centerX, centerY, startRadius, endRadius, turns
}

// Advanced options for simple text-on-path

// SetSimpleTextAlignment sets text alignment
func (stp *SimpleTextOnPath) SetAlignment(alignment SimpleTextAlignment) {
	stp.Alignment = alignment
}

// SetSimpleTextSpacing sets character spacing multiplier
func (stp *SimpleTextOnPath) SetSpacing(spacing float64) {
	stp.Spacing = spacing
}

// SetSimpleTextOffset sets offset along path
func (stp *SimpleTextOnPath) SetOffset(offset float64) {
	stp.Offset = offset
}

// SetSimpleTextFontSize sets font size
func (stp *SimpleTextOnPath) SetFontSize(size float64) {
	stp.FontSize = size
}

// Additional convenience methods

// DrawTextOnArc draws text along an arc
func DrawTextOnArc(dc *core.Context, text string, centerX, centerY, radius, startAngle, endAngle float64) {
	if text == "" {
		return
	}

	textOnPath := NewSimpleTextOnPath(text)
	textOnPath.renderOnArc(dc, centerX, centerY, radius, startAngle, endAngle)
}

// renderOnArc renders text along an arc
func (stp *SimpleTextOnPath) renderOnArc(dc *core.Context, centerX, centerY, radius, startAngle, endAngle float64) {
	arcLength := math.Abs(endAngle-startAngle) * radius
	textWidth := stp.estimateTextWidth()

	// Calculate starting angle based on alignment
	currentAngle := startAngle + stp.Offset/radius
	switch stp.Alignment {
	case SimpleAlignCenter:
		currentAngle += (arcLength - textWidth) / (2 * radius)
	case SimpleAlignEnd:
		currentAngle += (arcLength - textWidth) / radius
	}

	direction := 1.0
	if endAngle < startAngle {
		direction = -1.0
	}

	for _, r := range stp.Text {
		if (direction > 0 && currentAngle >= endAngle) || (direction < 0 && currentAngle <= endAngle) {
			break
		}

		// Calculate position
		x := centerX + radius*math.Cos(currentAngle)
		y := centerY + radius*math.Sin(currentAngle)

		// Calculate rotation (tangent to arc)
		rotation := currentAngle + math.Pi/2*direction

		// Render character
		stp.renderCharacter(dc, r, x, y, rotation)

		// Advance angle
		charWidth := stp.getCharacterWidth(r)
		angleAdvance := (charWidth * stp.Spacing) / radius * direction
		currentAngle += angleAdvance
	}
}

// DrawTextOnBezier draws text along a simple quadratic Bezier curve
func DrawTextOnBezier(dc *core.Context, text string, startX, startY, controlX, controlY, endX, endY float64) {
	if text == "" {
		return
	}

	textOnPath := NewSimpleTextOnPath(text)
	textOnPath.renderOnBezier(dc, startX, startY, controlX, controlY, endX, endY)
}

// renderOnBezier renders text along a quadratic Bezier curve
func (stp *SimpleTextOnPath) renderOnBezier(dc *core.Context, startX, startY, controlX, controlY, endX, endY float64) {
	// Estimate curve length
	curveLength := stp.estimateBezierLength(startX, startY, controlX, controlY, endX, endY)
	textWidth := stp.estimateTextWidth()

	// Calculate starting parameter based on alignment
	startT := stp.Offset / curveLength
	switch stp.Alignment {
	case SimpleAlignCenter:
		startT += (curveLength - textWidth) / (2 * curveLength)
	case SimpleAlignEnd:
		startT += (curveLength - textWidth) / curveLength
	}

	currentT := startT

	for _, r := range stp.Text {
		if currentT >= 1.0 {
			break
		}

		// Calculate position on Bezier curve
		x, y := stp.evaluateQuadBezier(startX, startY, controlX, controlY, endX, endY, currentT)

		// Calculate tangent
		tangent := stp.getBezierTangent(startX, startY, controlX, controlY, endX, endY, currentT)

		// Render character
		stp.renderCharacter(dc, r, x, y, tangent)

		// Advance parameter
		charWidth := stp.getCharacterWidth(r)
		tAdvance := (charWidth * stp.Spacing) / curveLength
		currentT += tAdvance
	}
}

// evaluateQuadBezier evaluates a quadratic Bezier curve at parameter t
func (stp *SimpleTextOnPath) evaluateQuadBezier(x0, y0, x1, y1, x2, y2, t float64) (float64, float64) {
	u := 1 - t
	x := u*u*x0 + 2*u*t*x1 + t*t*x2
	y := u*u*y0 + 2*u*t*y1 + t*t*y2
	return x, y
}

// getBezierTangent gets the tangent angle at parameter t
func (stp *SimpleTextOnPath) getBezierTangent(x0, y0, x1, y1, x2, y2, t float64) float64 {
	// Derivative of quadratic Bezier
	u := 1 - t
	dx := 2*u*(x1-x0) + 2*t*(x2-x1)
	dy := 2*u*(y1-y0) + 2*t*(y2-y1)
	return math.Atan2(dy, dx)
}

// estimateBezierLength estimates the length of a quadratic Bezier curve
func (stp *SimpleTextOnPath) estimateBezierLength(x0, y0, x1, y1, x2, y2 float64) float64 {
	// Simple approximation using chord and control polygon
	chord := math.Sqrt((x2-x0)*(x2-x0) + (y2-y0)*(y2-y0))
	poly := math.Sqrt((x1-x0)*(x1-x0)+(y1-y0)*(y1-y0)) +
		math.Sqrt((x2-x1)*(x2-x1)+(y2-y1)*(y2-y1))
	return (chord + poly) / 2
}

// Advanced text-on-path features

// DrawTextOnPath draws text along a custom path using flattened points
func DrawTextOnPath(dc *core.Context, text string, path *core.Path2D) {
	if text == "" || path == nil || path.IsEmpty() {
		return
	}

	textOnPath := NewSimpleTextOnPath(text)
	textOnPath.renderOnFlattenedPath(dc, path)
}

// renderOnFlattenedPath renders text along a flattened path
func (stp *SimpleTextOnPath) renderOnFlattenedPath(dc *core.Context, path *core.Path2D) {
	// Flatten the path to get points
	rasterPath := path.GetPath()
	flatPaths := flattenPath(rasterPath)

	if len(flatPaths) == 0 {
		return
	}

	// Use the first path (for simplicity)
	points := flatPaths[0]
	if len(points) < 2 {
		return
	}

	// Calculate total path length
	totalLength := stp.calculateFlatPathLength(points)
	textWidth := stp.estimateTextWidth()

	// Calculate starting position based on alignment
	startDistance := stp.Offset
	switch stp.Alignment {
	case SimpleAlignCenter:
		startDistance += (totalLength - textWidth) / 2
	case SimpleAlignEnd:
		startDistance += totalLength - textWidth
	}

	currentDistance := startDistance

	for _, r := range stp.Text {
		if currentDistance >= totalLength {
			break
		}

		// Find position and tangent at current distance
		x, y, tangent := stp.getPositionAtDistanceFlat(points, currentDistance)

		// Render character
		stp.renderCharacter(dc, r, x, y, tangent)

		// Advance distance
		charWidth := stp.getCharacterWidth(r)
		currentDistance += charWidth * stp.Spacing
	}
}

// calculateFlatPathLength calculates the total length of a flattened path
func (stp *SimpleTextOnPath) calculateFlatPathLength(points []core.Point) float64 {
	if len(points) < 2 {
		return 0
	}

	totalLength := 0.0
	for i := 1; i < len(points); i++ {
		dx := points[i].X - points[i-1].X
		dy := points[i].Y - points[i-1].Y
		totalLength += math.Sqrt(dx*dx + dy*dy)
	}

	return totalLength
}

// getPositionAtDistanceFlat finds position and tangent at distance along flattened path
func (stp *SimpleTextOnPath) getPositionAtDistanceFlat(points []core.Point, targetDistance float64) (float64, float64, float64) {
	if len(points) < 2 {
		return points[0].X, points[0].Y, 0
	}

	currentDistance := 0.0

	for i := 1; i < len(points); i++ {
		dx := points[i].X - points[i-1].X
		dy := points[i].Y - points[i-1].Y
		segmentLength := math.Sqrt(dx*dx + dy*dy)

		if currentDistance+segmentLength >= targetDistance {
			// Target is on this segment
			if segmentLength == 0 {
				return points[i-1].X, points[i-1].Y, 0
			}

			t := (targetDistance - currentDistance) / segmentLength
			x := points[i-1].X + t*dx
			y := points[i-1].Y + t*dy
			tangent := math.Atan2(dy, dx)
			return x, y, tangent
		}

		currentDistance += segmentLength
	}

	// If we reach here, return the last point
	lastPoint := points[len(points)-1]
	return lastPoint.X, lastPoint.Y, 0
}

// Helper function to access flattenPath from core package
func flattenPath(p raster.Path) [][]core.Point {
	// This function is defined in internal/core/path.go
	// We need to make it accessible or reimplement it
	var result [][]core.Point
	var path []core.Point
	var cx, cy float64

	for i := 0; i < len(p); {
		switch p[i] {
		case 0: // MoveTo
			if len(path) > 0 {
				result = append(result, path)
				path = nil
			}
			x := unfix(p[i+1])
			y := unfix(p[i+2])
			path = append(path, core.Point{X: x, Y: y})
			cx, cy = x, y
			i += 4
		case 1: // LineTo
			x := unfix(p[i+1])
			y := unfix(p[i+2])
			path = append(path, core.Point{X: x, Y: y})
			cx, cy = x, y
			i += 4
		case 2: // QuadTo
			x1 := unfix(p[i+1])
			y1 := unfix(p[i+2])
			x2 := unfix(p[i+3])
			y2 := unfix(p[i+4])
			points := core.QuadraticBezier(cx, cy, x1, y1, x2, y2)
			path = append(path, points...)
			cx, cy = x2, y2
			i += 6
		case 3: // CubeTo
			x1 := unfix(p[i+1])
			y1 := unfix(p[i+2])
			x2 := unfix(p[i+3])
			y2 := unfix(p[i+4])
			x3 := unfix(p[i+5])
			y3 := unfix(p[i+6])
			points := core.CubicBezier(cx, cy, x1, y1, x2, y2, x3, y3)
			path = append(path, points...)
			cx, cy = x3, y3
			i += 8
		default:
			// Skip unknown commands
			i++
		}
	}
	if len(path) > 0 {
		result = append(result, path)
	}
	return result
}

// Helper function to convert fixed.Int26_6 to float64
func unfix(x fixed.Int26_6) float64 {
	return float64(x) / 64
}
