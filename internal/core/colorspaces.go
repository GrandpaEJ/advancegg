package core

import (
	"image/color"
	"math"
)

// Color represents a color in various color spaces
type Color struct {
	R, G, B, A float64
}

// CMYK represents a color in CMYK color space
type CMYK struct {
	C, M, Y, K float64 // Cyan, Magenta, Yellow, Key (Black)
}

// HSV represents a color in HSV color space
type HSV struct {
	H, S, V float64 // Hue (0-360), Saturation (0-1), Value (0-1)
}

// HSL represents a color in HSL color space
type HSL struct {
	H, S, L float64 // Hue (0-360), Saturation (0-1), Lightness (0-1)
}

// LAB represents a color in LAB color space
type LAB struct {
	L, A, B float64 // Lightness, A (green-red), B (blue-yellow)
}

// XYZ represents a color in XYZ color space (intermediate for LAB conversion)
type XYZ struct {
	X, Y, Z float64
}

// NewColor creates a new Color from RGBA values (0-1 range)
func NewColor(r, g, b, a float64) Color {
	return Color{R: r, G: g, B: b, A: a}
}

// NewColorFromRGBA255 creates a new Color from RGBA values (0-255 range)
func NewColorFromRGBA255(r, g, b, a uint8) Color {
	return Color{
		R: float64(r) / 255.0,
		G: float64(g) / 255.0,
		B: float64(b) / 255.0,
		A: float64(a) / 255.0,
	}
}

// RGBA returns the color as standard RGBA values (0-65535 range)
func (c Color) RGBA() (r, g, b, a uint32) {
	r = uint32(c.R * 65535)
	g = uint32(c.G * 65535)
	b = uint32(c.B * 65535)
	a = uint32(c.A * 65535)
	return
}

// RGBA255 returns the color as RGBA values (0-255 range)
func (c Color) RGBA255() (r, g, b, a uint8) {
	return uint8(c.R * 255), uint8(c.G * 255), uint8(c.B * 255), uint8(c.A * 255)
}

// ToStandardColor converts to Go's standard color.Color interface
func (c Color) ToStandardColor() color.Color {
	r, g, b, a := c.RGBA255()
	return color.RGBA{r, g, b, a}
}

// RGB to CMYK conversion
func (c Color) ToCMYK() CMYK {
	// Convert RGB to CMY
	cyan := 1.0 - c.R
	magenta := 1.0 - c.G
	yellow := 1.0 - c.B
	
	// Find the minimum CMY value (this becomes K)
	k := math.Min(cyan, math.Min(magenta, yellow))
	
	// Adjust CMY values
	if k < 1.0 {
		cyan = (cyan - k) / (1.0 - k)
		magenta = (magenta - k) / (1.0 - k)
		yellow = (yellow - k) / (1.0 - k)
	} else {
		cyan = 0
		magenta = 0
		yellow = 0
	}
	
	return CMYK{C: cyan, M: magenta, Y: yellow, K: k}
}

// CMYK to RGB conversion
func (cmyk CMYK) ToRGB() Color {
	r := (1.0 - cmyk.C) * (1.0 - cmyk.K)
	g := (1.0 - cmyk.M) * (1.0 - cmyk.K)
	b := (1.0 - cmyk.Y) * (1.0 - cmyk.K)
	
	return Color{R: r, G: g, B: b, A: 1.0}
}

// RGB to HSV conversion
func (c Color) ToHSV() HSV {
	max := math.Max(c.R, math.Max(c.G, c.B))
	min := math.Min(c.R, math.Min(c.G, c.B))
	delta := max - min
	
	var h, s, v float64
	
	// Value
	v = max
	
	// Saturation
	if max != 0 {
		s = delta / max
	} else {
		s = 0
	}
	
	// Hue
	if delta == 0 {
		h = 0
	} else if max == c.R {
		h = 60 * math.Mod((c.G-c.B)/delta, 6)
	} else if max == c.G {
		h = 60 * ((c.B-c.R)/delta + 2)
	} else {
		h = 60 * ((c.R-c.G)/delta + 4)
	}
	
	if h < 0 {
		h += 360
	}
	
	return HSV{H: h, S: s, V: v}
}

// HSV to RGB conversion
func (hsv HSV) ToRGB() Color {
	c := hsv.V * hsv.S
	x := c * (1 - math.Abs(math.Mod(hsv.H/60, 2)-1))
	m := hsv.V - c
	
	var r, g, b float64
	
	switch {
	case hsv.H >= 0 && hsv.H < 60:
		r, g, b = c, x, 0
	case hsv.H >= 60 && hsv.H < 120:
		r, g, b = x, c, 0
	case hsv.H >= 120 && hsv.H < 180:
		r, g, b = 0, c, x
	case hsv.H >= 180 && hsv.H < 240:
		r, g, b = 0, x, c
	case hsv.H >= 240 && hsv.H < 300:
		r, g, b = x, 0, c
	case hsv.H >= 300 && hsv.H < 360:
		r, g, b = c, 0, x
	}
	
	return Color{R: r + m, G: g + m, B: b + m, A: 1.0}
}

// RGB to HSL conversion
func (c Color) ToHSL() HSL {
	max := math.Max(c.R, math.Max(c.G, c.B))
	min := math.Min(c.R, math.Min(c.G, c.B))
	delta := max - min
	
	var h, s, l float64
	
	// Lightness
	l = (max + min) / 2
	
	// Saturation
	if delta == 0 {
		s = 0
	} else if l < 0.5 {
		s = delta / (max + min)
	} else {
		s = delta / (2 - max - min)
	}
	
	// Hue (same as HSV)
	if delta == 0 {
		h = 0
	} else if max == c.R {
		h = 60 * math.Mod((c.G-c.B)/delta, 6)
	} else if max == c.G {
		h = 60 * ((c.B-c.R)/delta + 2)
	} else {
		h = 60 * ((c.R-c.G)/delta + 4)
	}
	
	if h < 0 {
		h += 360
	}
	
	return HSL{H: h, S: s, L: l}
}

// HSL to RGB conversion
func (hsl HSL) ToRGB() Color {
	c := (1 - math.Abs(2*hsl.L-1)) * hsl.S
	x := c * (1 - math.Abs(math.Mod(hsl.H/60, 2)-1))
	m := hsl.L - c/2
	
	var r, g, b float64
	
	switch {
	case hsl.H >= 0 && hsl.H < 60:
		r, g, b = c, x, 0
	case hsl.H >= 60 && hsl.H < 120:
		r, g, b = x, c, 0
	case hsl.H >= 120 && hsl.H < 180:
		r, g, b = 0, c, x
	case hsl.H >= 180 && hsl.H < 240:
		r, g, b = 0, x, c
	case hsl.H >= 240 && hsl.H < 300:
		r, g, b = x, 0, c
	case hsl.H >= 300 && hsl.H < 360:
		r, g, b = c, 0, x
	}
	
	return Color{R: r + m, G: g + m, B: b + m, A: 1.0}
}

// RGB to XYZ conversion (intermediate step for LAB)
func (c Color) ToXYZ() XYZ {
	// Convert to linear RGB
	var r, g, b float64
	
	if c.R > 0.04045 {
		r = math.Pow((c.R+0.055)/1.055, 2.4)
	} else {
		r = c.R / 12.92
	}
	
	if c.G > 0.04045 {
		g = math.Pow((c.G+0.055)/1.055, 2.4)
	} else {
		g = c.G / 12.92
	}
	
	if c.B > 0.04045 {
		b = math.Pow((c.B+0.055)/1.055, 2.4)
	} else {
		b = c.B / 12.92
	}
	
	// Convert to XYZ using sRGB matrix
	x := r*0.4124564 + g*0.3575761 + b*0.1804375
	y := r*0.2126729 + g*0.7151522 + b*0.0721750
	z := r*0.0193339 + g*0.1191920 + b*0.9503041
	
	return XYZ{X: x * 100, Y: y * 100, Z: z * 100}
}

// XYZ to RGB conversion
func (xyz XYZ) ToRGB() Color {
	// Normalize
	x := xyz.X / 100
	y := xyz.Y / 100
	z := xyz.Z / 100
	
	// Convert to linear RGB
	r := x*3.2404542 + y*-1.5371385 + z*-0.4985314
	g := x*-0.9692660 + y*1.8760108 + z*0.0415560
	b := x*0.0556434 + y*-0.2040259 + z*1.0572252
	
	// Convert to sRGB
	if r > 0.0031308 {
		r = 1.055*math.Pow(r, 1/2.4) - 0.055
	} else {
		r = 12.92 * r
	}
	
	if g > 0.0031308 {
		g = 1.055*math.Pow(g, 1/2.4) - 0.055
	} else {
		g = 12.92 * g
	}
	
	if b > 0.0031308 {
		b = 1.055*math.Pow(b, 1/2.4) - 0.055
	} else {
		b = 12.92 * b
	}
	
	// Clamp values
	r = math.Max(0, math.Min(1, r))
	g = math.Max(0, math.Min(1, g))
	b = math.Max(0, math.Min(1, b))
	
	return Color{R: r, G: g, B: b, A: 1.0}
}

// RGB to LAB conversion
func (c Color) ToLAB() LAB {
	xyz := c.ToXYZ()
	return xyz.ToLAB()
}

// XYZ to LAB conversion
func (xyz XYZ) ToLAB() LAB {
	// Reference white D65
	xn, yn, zn := 95.047, 100.000, 108.883
	
	x := xyz.X / xn
	y := xyz.Y / yn
	z := xyz.Z / zn
	
	// Apply the LAB transformation
	fx := labF(x)
	fy := labF(y)
	fz := labF(z)
	
	l := 116*fy - 16
	a := 500 * (fx - fy)
	b := 200 * (fy - fz)
	
	return LAB{L: l, A: a, B: b}
}

// LAB to XYZ conversion
func (lab LAB) ToXYZ() XYZ {
	// Reference white D65
	xn, yn, zn := 95.047, 100.000, 108.883
	
	fy := (lab.L + 16) / 116
	fx := lab.A/500 + fy
	fz := fy - lab.B/200
	
	x := labFInv(fx) * xn
	y := labFInv(fy) * yn
	z := labFInv(fz) * zn
	
	return XYZ{X: x, Y: y, Z: z}
}

// LAB to RGB conversion
func (lab LAB) ToRGB() Color {
	xyz := lab.ToXYZ()
	return xyz.ToRGB()
}

// Helper function for LAB conversion
func labF(t float64) float64 {
	if t > 0.008856 {
		return math.Pow(t, 1.0/3.0)
	}
	return (7.787*t + 16.0/116.0)
}

// Inverse helper function for LAB conversion
func labFInv(t float64) float64 {
	if t > 0.206893 {
		return math.Pow(t, 3.0)
	}
	return (t - 16.0/116.0) / 7.787
}
