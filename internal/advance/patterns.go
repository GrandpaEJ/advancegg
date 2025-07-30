package advance

import (
	"image"
	"image/color"
	"math"
)

// Pattern represents a repeatable pattern
type Pattern interface {
	ColorAt(x, y float64) color.Color
}

// LinearGradientPattern creates a linear gradient pattern
type LinearGradientPattern struct {
	X1, Y1, X2, Y2 float64
	ColorStops     []ColorStop
}

// ColorStop represents a color stop in a gradient
type ColorStop struct {
	Position float64    // 0.0 to 1.0
	Color    color.Color
}

// ColorAt returns the color at the specified coordinates
func (p LinearGradientPattern) ColorAt(x, y float64) color.Color {
	// Calculate position along gradient line
	dx := p.X2 - p.X1
	dy := p.Y2 - p.Y1
	length := math.Sqrt(dx*dx + dy*dy)
	
	if length == 0 {
		return p.ColorStops[0].Color
	}
	
	// Project point onto gradient line
	t := ((x-p.X1)*dx + (y-p.Y1)*dy) / (length * length)
	t = math.Max(0, math.Min(1, t))
	
	// Find color stops to interpolate between
	for i := 0; i < len(p.ColorStops)-1; i++ {
		stop1 := p.ColorStops[i]
		stop2 := p.ColorStops[i+1]
		
		if t >= stop1.Position && t <= stop2.Position {
			// Interpolate between colors
			ratio := (t - stop1.Position) / (stop2.Position - stop1.Position)
			return interpolateColors(stop1.Color, stop2.Color, ratio)
		}
	}
	
	// Return last color if beyond range
	return p.ColorStops[len(p.ColorStops)-1].Color
}

// RadialGradientPattern creates a radial gradient pattern
type RadialGradientPattern struct {
	CX, CY     float64 // Center
	Radius     float64
	ColorStops []ColorStop
}

// ColorAt returns the color at the specified coordinates
func (p RadialGradientPattern) ColorAt(x, y float64) color.Color {
	// Calculate distance from center
	dx := x - p.CX
	dy := y - p.CY
	distance := math.Sqrt(dx*dx + dy*dy)
	
	// Normalize distance
	t := distance / p.Radius
	t = math.Max(0, math.Min(1, t))
	
	// Find color stops to interpolate between
	for i := 0; i < len(p.ColorStops)-1; i++ {
		stop1 := p.ColorStops[i]
		stop2 := p.ColorStops[i+1]
		
		if t >= stop1.Position && t <= stop2.Position {
			ratio := (t - stop1.Position) / (stop2.Position - stop1.Position)
			return interpolateColors(stop1.Color, stop2.Color, ratio)
		}
	}
	
	return p.ColorStops[len(p.ColorStops)-1].Color
}

// CheckerboardPattern creates a checkerboard pattern
type CheckerboardPattern struct {
	Size   float64
	Color1 color.Color
	Color2 color.Color
}

// ColorAt returns the color at the specified coordinates
func (p CheckerboardPattern) ColorAt(x, y float64) color.Color {
	cellX := int(math.Floor(x / p.Size))
	cellY := int(math.Floor(y / p.Size))
	
	if (cellX+cellY)%2 == 0 {
		return p.Color1
	}
	return p.Color2
}

// StripePattern creates a stripe pattern
type StripePattern struct {
	Width     float64
	Angle     float64 // In radians
	Color1    color.Color
	Color2    color.Color
}

// ColorAt returns the color at the specified coordinates
func (p StripePattern) ColorAt(x, y float64) color.Color {
	// Rotate coordinates
	cos := math.Cos(p.Angle)
	sin := math.Sin(p.Angle)
	rotX := x*cos - y*sin
	
	// Determine stripe
	stripe := int(math.Floor(rotX / p.Width))
	if stripe%2 == 0 {
		return p.Color1
	}
	return p.Color2
}

// PolkaDotPattern creates a polka dot pattern
type PolkaDotPattern struct {
	SpacingX, SpacingY float64
	Radius             float64
	DotColor           color.Color
	BackgroundColor    color.Color
}

// ColorAt returns the color at the specified coordinates
func (p PolkaDotPattern) ColorAt(x, y float64) color.Color {
	// Find nearest dot center
	cellX := math.Floor(x / p.SpacingX)
	cellY := math.Floor(y / p.SpacingY)
	
	centerX := (cellX + 0.5) * p.SpacingX
	centerY := (cellY + 0.5) * p.SpacingY
	
	// Calculate distance to center
	dx := x - centerX
	dy := y - centerY
	distance := math.Sqrt(dx*dx + dy*dy)
	
	if distance <= p.Radius {
		return p.DotColor
	}
	return p.BackgroundColor
}

// NoisePattern creates a noise pattern
type NoisePattern struct {
	Scale     float64
	BaseColor color.Color
	Intensity float64
}

// ColorAt returns the color at the specified coordinates
func (p NoisePattern) ColorAt(x, y float64) color.Color {
	// Simple noise function (pseudo-random)
	noise := p.simpleNoise(x*p.Scale, y*p.Scale)
	
	r, g, b, a := p.BaseColor.RGBA()
	
	// Apply noise
	factor := 1.0 + (noise-0.5)*p.Intensity
	
	newR := uint8(math.Max(0, math.Min(255, float64(r>>8)*factor)))
	newG := uint8(math.Max(0, math.Min(255, float64(g>>8)*factor)))
	newB := uint8(math.Max(0, math.Min(255, float64(b>>8)*factor)))
	
	return color.RGBA{newR, newG, newB, uint8(a >> 8)}
}

// Simple noise function
func (p NoisePattern) simpleNoise(x, y float64) float64 {
	// Simple pseudo-random noise
	n := math.Sin(x*12.9898+y*78.233) * 43758.5453
	return n - math.Floor(n)
}

// WavePattern creates a wave pattern
type WavePattern struct {
	Wavelength  float64
	Amplitude   float64
	Angle       float64
	Color1      color.Color
	Color2      color.Color
}

// ColorAt returns the color at the specified coordinates
func (p WavePattern) ColorAt(x, y float64) color.Color {
	// Rotate coordinates
	cos := math.Cos(p.Angle)
	sin := math.Sin(p.Angle)
	rotX := x*cos - y*sin
	
	// Calculate wave
	wave := math.Sin(2*math.Pi*rotX/p.Wavelength) * p.Amplitude
	
	// Interpolate between colors based on wave value
	t := (wave + p.Amplitude) / (2 * p.Amplitude) // Normalize to 0-1
	return interpolateColors(p.Color1, p.Color2, t)
}

// Helper function to interpolate between two colors
func interpolateColors(c1, c2 color.Color, t float64) color.Color {
	r1, g1, b1, a1 := c1.RGBA()
	r2, g2, b2, a2 := c2.RGBA()
	
	t = math.Max(0, math.Min(1, t))
	
	r := uint8((float64(r1>>8)*(1-t) + float64(r2>>8)*t))
	g := uint8((float64(g1>>8)*(1-t) + float64(g2>>8)*t))
	b := uint8((float64(b1>>8)*(1-t) + float64(b2>>8)*t))
	a := uint8((float64(a1>>8)*(1-t) + float64(a2>>8)*t))
	
	return color.RGBA{r, g, b, a}
}

// PatternFill fills an image with a pattern
func PatternFill(img *image.RGBA, pattern Pattern) {
	bounds := img.Bounds()
	
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c := pattern.ColorAt(float64(x), float64(y))
			img.Set(x, y, c)
		}
	}
}

// Convenience functions for creating common patterns

// CreateLinearGradient creates a linear gradient from top to bottom
func CreateLinearGradient(width, height float64, colors ...color.Color) LinearGradientPattern {
	stops := make([]ColorStop, len(colors))
	for i, c := range colors {
		stops[i] = ColorStop{
			Position: float64(i) / float64(len(colors)-1),
			Color:    c,
		}
	}
	
	return LinearGradientPattern{
		X1: 0, Y1: 0,
		X2: 0, Y2: height,
		ColorStops: stops,
	}
}

// CreateRadialGradient creates a radial gradient from center
func CreateRadialGradient(cx, cy, radius float64, colors ...color.Color) RadialGradientPattern {
	stops := make([]ColorStop, len(colors))
	for i, c := range colors {
		stops[i] = ColorStop{
			Position: float64(i) / float64(len(colors)-1),
			Color:    c,
		}
	}
	
	return RadialGradientPattern{
		CX: cx, CY: cy,
		Radius:     radius,
		ColorStops: stops,
	}
}

// CreateCheckerboard creates a simple checkerboard pattern
func CreateCheckerboard(size float64) CheckerboardPattern {
	return CheckerboardPattern{
		Size:   size,
		Color1: color.RGBA{255, 255, 255, 255}, // White
		Color2: color.RGBA{0, 0, 0, 255},       // Black
	}
}

// CreateStripes creates a diagonal stripe pattern
func CreateStripes(width float64) StripePattern {
	return StripePattern{
		Width:  width,
		Angle:  math.Pi / 4, // 45 degrees
		Color1: color.RGBA{255, 255, 255, 255},
		Color2: color.RGBA{0, 0, 0, 255},
	}
}

// CreatePolkaDots creates a polka dot pattern
func CreatePolkaDots(spacing, radius float64) PolkaDotPattern {
	return PolkaDotPattern{
		SpacingX:        spacing,
		SpacingY:        spacing,
		Radius:          radius,
		DotColor:        color.RGBA{255, 0, 0, 255}, // Red dots
		BackgroundColor: color.RGBA{255, 255, 255, 255}, // White background
	}
}
