package core

import (
	"image/color"
	"math"
)

// Advanced stroke styles including dashed patterns, gradient strokes, and tapered ends

// StrokeStyle represents advanced stroke styling options
type StrokeStyle struct {
	Width       float64
	Color       color.Color
	DashPattern []float64
	DashOffset  float64
	LineCap     StrokeLineCap
	LineJoin    StrokeLineJoin
	MiterLimit  float64
	Gradient    *StrokeGradient
	Taper       *StrokeTaper
}

// StrokeLineCap represents line cap styles
type StrokeLineCap int

const (
	StrokeLineCapButt   StrokeLineCap = iota // Square end
	StrokeLineCapRound                       // Rounded end
	StrokeLineCapSquare                      // Extended square end
)

// StrokeLineJoin represents line join styles
type StrokeLineJoin int

const (
	StrokeLineJoinMiter StrokeLineJoin = iota // Sharp corner
	StrokeLineJoinRound                       // Rounded corner
	StrokeLineJoinBevel                       // Beveled corner
)

// StrokeGradient represents gradient stroke coloring
type StrokeGradient struct {
	Type      StrokeGradientType
	Colors    []StrokeGradientStop
	StartX    float64
	StartY    float64
	EndX      float64
	EndY      float64
	CenterX   float64 // For radial gradients
	CenterY   float64
	Radius    float64
}

// StrokeGradientType represents gradient types
type StrokeGradientType int

const (
	StrokeGradientLinear StrokeGradientType = iota
	StrokeGradientRadial
	StrokeGradientConic
)

// StrokeGradientStop represents a color stop in a gradient
type StrokeGradientStop struct {
	Position float64 // 0.0 to 1.0
	Color    color.Color
}

// StrokeTaper represents tapered stroke ends
type StrokeTaper struct {
	StartWidth float64 // Width at start (0.0 to 1.0 multiplier)
	EndWidth   float64 // Width at end (0.0 to 1.0 multiplier)
	Type       StrokeTaperType
}

// StrokeTaperType represents taper styles
type StrokeTaperType int

const (
	StrokeTaperLinear     StrokeTaperType = iota // Linear taper
	StrokeTaperExponential                       // Exponential taper
	StrokeTaperSinusoidal                        // Smooth sinusoidal taper
)

// NewStrokeStyle creates a new stroke style with defaults
func NewStrokeStyle() *StrokeStyle {
	return &StrokeStyle{
		Width:      1.0,
		Color:      color.RGBA{0, 0, 0, 255},
		LineCap:    StrokeLineCapRound,
		LineJoin:   StrokeLineJoinRound,
		MiterLimit: 10.0,
	}
}

// SetDashPattern sets a dash pattern for the stroke
func (ss *StrokeStyle) SetDashPattern(pattern []float64, offset float64) {
	ss.DashPattern = make([]float64, len(pattern))
	copy(ss.DashPattern, pattern)
	ss.DashOffset = offset
}

// SetLinearGradient sets a linear gradient for the stroke
func (ss *StrokeStyle) SetLinearGradient(startX, startY, endX, endY float64, stops []StrokeGradientStop) {
	ss.Gradient = &StrokeGradient{
		Type:   StrokeGradientLinear,
		Colors: make([]StrokeGradientStop, len(stops)),
		StartX: startX,
		StartY: startY,
		EndX:   endX,
		EndY:   endY,
	}
	copy(ss.Gradient.Colors, stops)
}

// SetRadialGradient sets a radial gradient for the stroke
func (ss *StrokeStyle) SetRadialGradient(centerX, centerY, radius float64, stops []StrokeGradientStop) {
	ss.Gradient = &StrokeGradient{
		Type:    StrokeGradientRadial,
		Colors:  make([]StrokeGradientStop, len(stops)),
		CenterX: centerX,
		CenterY: centerY,
		Radius:  radius,
	}
	copy(ss.Gradient.Colors, stops)
}

// SetTaper sets stroke tapering
func (ss *StrokeStyle) SetTaper(startWidth, endWidth float64, taperType StrokeTaperType) {
	ss.Taper = &StrokeTaper{
		StartWidth: startWidth,
		EndWidth:   endWidth,
		Type:       taperType,
	}
}

// Context integration

// SetAdvancedStroke sets an advanced stroke style for the context
func (dc *Context) SetAdvancedStroke(style *StrokeStyle) {
	dc.advancedStroke = style
}

// GetAdvancedStroke returns the current advanced stroke style
func (dc *Context) GetAdvancedStroke() *StrokeStyle {
	return dc.advancedStroke
}

// DrawDashedLine draws a dashed line
func (dc *Context) DrawDashedLine(x1, y1, x2, y2 float64, pattern []float64) {
	style := NewStrokeStyle()
	style.Width = dc.lineWidth
	style.Color = dc.color
	style.SetDashPattern(pattern, 0)
	
	dc.drawDashedLine(x1, y1, x2, y2, style)
}

// DrawGradientLine draws a line with gradient stroke
func (dc *Context) DrawGradientLine(x1, y1, x2, y2 float64, stops []StrokeGradientStop) {
	style := NewStrokeStyle()
	style.Width = dc.lineWidth
	style.SetLinearGradient(x1, y1, x2, y2, stops)
	
	dc.drawGradientLine(x1, y1, x2, y2, style)
}

// DrawTaperedLine draws a line with tapered ends
func (dc *Context) DrawTaperedLine(x1, y1, x2, y2 float64, startWidth, endWidth float64) {
	style := NewStrokeStyle()
	style.Width = dc.lineWidth
	style.Color = dc.color
	style.SetTaper(startWidth, endWidth, StrokeTaperLinear)
	
	dc.drawTaperedLine(x1, y1, x2, y2, style)
}

// Internal rendering methods

// drawDashedLine draws a dashed line
func (dc *Context) drawDashedLine(x1, y1, x2, y2 float64, style *StrokeStyle) {
	if len(style.DashPattern) == 0 {
		dc.MoveTo(x1, y1)
		dc.LineTo(x2, y2)
		dc.Stroke()
		return
	}
	
	// Calculate line length and direction
	dx := x2 - x1
	dy := y2 - y1
	length := math.Sqrt(dx*dx + dy*dy)
	
	if length == 0 {
		return
	}
	
	// Normalize direction
	dirX := dx / length
	dirY := dy / length
	
	// Draw dashed segments
	currentDistance := style.DashOffset
	patternIndex := 0
	drawing := true // Start drawing (not in gap)
	
	for currentDistance < length {
		patternLength := style.DashPattern[patternIndex%len(style.DashPattern)]
		segmentEnd := math.Min(currentDistance+patternLength, length)
		
		if drawing {
			// Draw segment
			startX := x1 + currentDistance*dirX
			startY := y1 + currentDistance*dirY
			endX := x1 + segmentEnd*dirX
			endY := y1 + segmentEnd*dirY
			
			dc.MoveTo(startX, startY)
			dc.LineTo(endX, endY)
			dc.Stroke()
		}
		
		currentDistance = segmentEnd
		patternIndex++
		drawing = !drawing
	}
}

// drawGradientLine draws a line with gradient stroke
func (dc *Context) drawGradientLine(x1, y1, x2, y2 float64, style *StrokeStyle) {
	if style.Gradient == nil || len(style.Gradient.Colors) < 2 {
		dc.MoveTo(x1, y1)
		dc.LineTo(x2, y2)
		dc.Stroke()
		return
	}
	
	// Draw line as multiple segments with interpolated colors
	segments := 20
	for i := 0; i < segments; i++ {
		t1 := float64(i) / float64(segments)
		t2 := float64(i+1) / float64(segments)
		
		x1Seg := x1 + t1*(x2-x1)
		y1Seg := y1 + t1*(y2-y1)
		x2Seg := x1 + t2*(x2-x1)
		y2Seg := y1 + t2*(y2-y1)
		
		// Interpolate color
		segmentColor := dc.interpolateGradientColor(t1, style.Gradient.Colors)
		dc.SetColor(segmentColor)
		
		dc.MoveTo(x1Seg, y1Seg)
		dc.LineTo(x2Seg, y2Seg)
		dc.Stroke()
	}
}

// drawTaperedLine draws a line with tapered ends
func (dc *Context) drawTaperedLine(x1, y1, x2, y2 float64, style *StrokeStyle) {
	if style.Taper == nil {
		dc.MoveTo(x1, y1)
		dc.LineTo(x2, y2)
		dc.Stroke()
		return
	}
	
	// Draw line as multiple segments with varying width
	segments := 20
	baseWidth := style.Width
	
	for i := 0; i < segments; i++ {
		t := float64(i) / float64(segments)
		
		// Calculate tapered width
		var widthMultiplier float64
		switch style.Taper.Type {
		case StrokeTaperLinear:
			widthMultiplier = style.Taper.StartWidth + t*(style.Taper.EndWidth-style.Taper.StartWidth)
		case StrokeTaperExponential:
			widthMultiplier = style.Taper.StartWidth * math.Pow(style.Taper.EndWidth/style.Taper.StartWidth, t)
		case StrokeTaperSinusoidal:
			widthMultiplier = style.Taper.StartWidth + (style.Taper.EndWidth-style.Taper.StartWidth)*math.Sin(t*math.Pi/2)
		default:
			widthMultiplier = 1.0
		}
		
		width := baseWidth * widthMultiplier
		if width < 0.1 {
			width = 0.1 // Minimum width
		}
		
		t1 := float64(i) / float64(segments)
		t2 := float64(i+1) / float64(segments)
		
		x1Seg := x1 + t1*(x2-x1)
		y1Seg := y1 + t1*(y2-y1)
		x2Seg := x1 + t2*(x2-x1)
		y2Seg := y1 + t2*(y2-y1)
		
		dc.SetLineWidth(width)
		dc.MoveTo(x1Seg, y1Seg)
		dc.LineTo(x2Seg, y2Seg)
		dc.Stroke()
	}
	
	// Restore original width
	dc.SetLineWidth(baseWidth)
}

// interpolateGradientColor interpolates color from gradient stops
func (dc *Context) interpolateGradientColor(t float64, stops []StrokeGradientStop) color.Color {
	if len(stops) == 0 {
		return color.RGBA{0, 0, 0, 255}
	}
	
	if len(stops) == 1 {
		return stops[0].Color
	}
	
	// Clamp t to [0, 1]
	if t <= 0 {
		return stops[0].Color
	}
	if t >= 1 {
		return stops[len(stops)-1].Color
	}
	
	// Find the two stops to interpolate between
	for i := 0; i < len(stops)-1; i++ {
		if t >= stops[i].Position && t <= stops[i+1].Position {
			// Interpolate between stops[i] and stops[i+1]
			localT := (t - stops[i].Position) / (stops[i+1].Position - stops[i].Position)
			return dc.interpolateColors(stops[i].Color, stops[i+1].Color, localT)
		}
	}
	
	return stops[len(stops)-1].Color
}

// interpolateColors interpolates between two colors
func (dc *Context) interpolateColors(c1, c2 color.Color, t float64) color.Color {
	r1, g1, b1, a1 := c1.RGBA()
	r2, g2, b2, a2 := c2.RGBA()
	
	r := uint8((float64(r1>>8)*(1-t) + float64(r2>>8)*t))
	g := uint8((float64(g1>>8)*(1-t) + float64(g2>>8)*t))
	b := uint8((float64(b1>>8)*(1-t) + float64(b2>>8)*t))
	a := uint8((float64(a1>>8)*(1-t) + float64(a2>>8)*t))
	
	return color.RGBA{r, g, b, a}
}

// Convenience functions for common stroke styles

// CreateDashedStroke creates a dashed stroke style
func CreateDashedStroke(width float64, color color.Color, pattern []float64) *StrokeStyle {
	style := NewStrokeStyle()
	style.Width = width
	style.Color = color
	style.SetDashPattern(pattern, 0)
	return style
}

// CreateGradientStroke creates a gradient stroke style
func CreateGradientStroke(width float64, startX, startY, endX, endY float64, stops []StrokeGradientStop) *StrokeStyle {
	style := NewStrokeStyle()
	style.Width = width
	style.SetLinearGradient(startX, startY, endX, endY, stops)
	return style
}

// CreateTaperedStroke creates a tapered stroke style
func CreateTaperedStroke(width float64, color color.Color, startWidth, endWidth float64) *StrokeStyle {
	style := NewStrokeStyle()
	style.Width = width
	style.Color = color
	style.SetTaper(startWidth, endWidth, StrokeTaperLinear)
	return style
}
