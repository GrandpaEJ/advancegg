package advance

import (
	"image"
	"image/color"
	"math"
)

// CSSFilter represents a CSS-like filter that can be applied to images
type CSSFilter interface {
	Apply(img image.Image) image.Image
}

// FilterChain allows chaining multiple filters
type FilterChain struct {
	filters []CSSFilter
}

// NewFilterChain creates a new filter chain
func NewFilterChain() *FilterChain {
	return &FilterChain{filters: make([]CSSFilter, 0)}
}

// Add adds a filter to the chain
func (fc *FilterChain) Add(filter CSSFilter) *FilterChain {
	fc.filters = append(fc.filters, filter)
	return fc
}

// Apply applies all filters in the chain
func (fc *FilterChain) Apply(img image.Image) image.Image {
	result := img
	for _, filter := range fc.filters {
		result = filter.Apply(result)
	}
	return result
}

// Brightness filter - CSS brightness()
type BrightnessFilter struct {
	Amount float64 // 0 = black, 1 = normal, >1 = brighter
}

func (f BrightnessFilter) Apply(img image.Image) image.Image {
	bounds := img.Bounds()
	result := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()

			newR := uint8(math.Min(255, float64(r>>8)*f.Amount))
			newG := uint8(math.Min(255, float64(g>>8)*f.Amount))
			newB := uint8(math.Min(255, float64(b>>8)*f.Amount))

			result.Set(x, y, color.RGBA{newR, newG, newB, uint8(a >> 8)})
		}
	}

	return result
}

// Contrast filter - CSS contrast()
type ContrastFilter struct {
	Amount float64 // 0 = gray, 1 = normal, >1 = more contrast
}

func (f ContrastFilter) Apply(img image.Image) image.Image {
	bounds := img.Bounds()
	result := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()

			// Apply contrast formula: (color - 128) * amount + 128
			newR := uint8(math.Max(0, math.Min(255, (float64(r>>8)-128)*f.Amount+128)))
			newG := uint8(math.Max(0, math.Min(255, (float64(g>>8)-128)*f.Amount+128)))
			newB := uint8(math.Max(0, math.Min(255, (float64(b>>8)-128)*f.Amount+128)))

			result.Set(x, y, color.RGBA{newR, newG, newB, uint8(a >> 8)})
		}
	}

	return result
}

// Saturate filter - CSS saturate()
type SaturateFilter struct {
	Amount float64 // 0 = grayscale, 1 = normal, >1 = more saturated
}

func (f SaturateFilter) Apply(img image.Image) image.Image {
	bounds := img.Bounds()
	result := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			rf, gf, bf := float64(r>>8), float64(g>>8), float64(b>>8)

			// Convert to grayscale
			gray := 0.299*rf + 0.587*gf + 0.114*bf

			// Interpolate between grayscale and original
			newR := uint8(gray + (rf-gray)*f.Amount)
			newG := uint8(gray + (gf-gray)*f.Amount)
			newB := uint8(gray + (bf-gray)*f.Amount)

			result.Set(x, y, color.RGBA{newR, newG, newB, uint8(a >> 8)})
		}
	}

	return result
}

// HueRotate filter - CSS hue-rotate()
type HueRotateFilter struct {
	Degrees float64 // Rotation in degrees
}

func (f HueRotateFilter) Apply(img image.Image) image.Image {
	bounds := img.Bounds()
	result := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()

			// Convert RGB to HSV
			h, s, v := rgbToHSV(float64(r>>8), float64(g>>8), float64(b>>8))

			// Rotate hue
			h = math.Mod(h+f.Degrees, 360)
			if h < 0 {
				h += 360
			}

			// Convert back to RGB
			newR, newG, newB := hsvToRGB(h, s, v)

			result.Set(x, y, color.RGBA{
				uint8(newR), uint8(newG), uint8(newB), uint8(a >> 8),
			})
		}
	}

	return result
}

// Invert filter - CSS invert()
type InvertFilter struct {
	Amount float64 // 0 = normal, 1 = fully inverted
}

func (f InvertFilter) Apply(img image.Image) image.Image {
	bounds := img.Bounds()
	result := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()

			rf, gf, bf := float64(r>>8), float64(g>>8), float64(b>>8)

			// Interpolate between original and inverted
			invertedR := 255 - rf
			invertedG := 255 - gf
			invertedB := 255 - bf

			newR := uint8(rf + (invertedR-rf)*f.Amount)
			newG := uint8(gf + (invertedG-gf)*f.Amount)
			newB := uint8(bf + (invertedB-bf)*f.Amount)

			result.Set(x, y, color.RGBA{newR, newG, newB, uint8(a >> 8)})
		}
	}

	return result
}

// Opacity filter - CSS opacity()
type OpacityFilter struct {
	Amount float64 // 0 = transparent, 1 = opaque
}

func (f OpacityFilter) Apply(img image.Image) image.Image {
	bounds := img.Bounds()
	result := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()

			newA := uint8(float64(a>>8) * f.Amount)

			result.Set(x, y, color.RGBA{
				uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), newA,
			})
		}
	}

	return result
}

// Blur filter - CSS blur()
type BlurFilter struct {
	Radius float64 // Blur radius in pixels
}

func (f BlurFilter) Apply(img image.Image) image.Image {
	bounds := img.Bounds()
	result := image.NewRGBA(bounds)

	r := int(f.Radius)
	if r < 1 {
		r = 1
	}

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			var rSum, gSum, bSum, aSum, count uint32

			for dy := -r; dy <= r; dy++ {
				for dx := -r; dx <= r; dx++ {
					px, py := x+dx, y+dy
					if px >= bounds.Min.X && px < bounds.Max.X && py >= bounds.Min.Y && py < bounds.Max.Y {
						pr, pg, pb, pa := img.At(px, py).RGBA()
						rSum += pr >> 8
						gSum += pg >> 8
						bSum += pb >> 8
						aSum += pa >> 8
						count++
					}
				}
			}

			if count > 0 {
				result.Set(x, y, color.RGBA{
					uint8(rSum / count),
					uint8(gSum / count),
					uint8(bSum / count),
					uint8(aSum / count),
				})
			}
		}
	}

	return result
}

// Sepia filter - CSS sepia()
type SepiaFilter struct {
	Amount float64 // 0 = normal, 1 = full sepia
}

func (f SepiaFilter) Apply(img image.Image) image.Image {
	bounds := img.Bounds()
	result := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			rf, gf, bf := float64(r>>8), float64(g>>8), float64(b>>8)

			// Sepia transformation
			sepiaR := math.Min(255, 0.393*rf+0.769*gf+0.189*bf)
			sepiaG := math.Min(255, 0.349*rf+0.686*gf+0.168*bf)
			sepiaB := math.Min(255, 0.272*rf+0.534*gf+0.131*bf)

			// Interpolate between original and sepia
			newR := uint8(rf + (sepiaR-rf)*f.Amount)
			newG := uint8(gf + (sepiaG-gf)*f.Amount)
			newB := uint8(bf + (sepiaB-bf)*f.Amount)

			result.Set(x, y, color.RGBA{newR, newG, newB, uint8(a >> 8)})
		}
	}

	return result
}

// Helper functions for color space conversions

func rgbToHSV(r, g, b float64) (h, s, v float64) {
	r, g, b = r/255, g/255, b/255
	max := math.Max(r, math.Max(g, b))
	min := math.Min(r, math.Min(g, b))
	delta := max - min

	v = max

	if max != 0 {
		s = delta / max
	} else {
		s = 0
	}

	if delta == 0 {
		h = 0
	} else if max == r {
		h = 60 * math.Mod((g-b)/delta, 6)
	} else if max == g {
		h = 60 * ((b-r)/delta + 2)
	} else {
		h = 60 * ((r-g)/delta + 4)
	}

	if h < 0 {
		h += 360
	}

	return h, s, v
}

func hsvToRGB(h, s, v float64) (r, g, b float64) {
	c := v * s
	x := c * (1 - math.Abs(math.Mod(h/60, 2)-1))
	m := v - c

	switch {
	case h >= 0 && h < 60:
		r, g, b = c, x, 0
	case h >= 60 && h < 120:
		r, g, b = x, c, 0
	case h >= 120 && h < 180:
		r, g, b = 0, c, x
	case h >= 180 && h < 240:
		r, g, b = 0, x, c
	case h >= 240 && h < 300:
		r, g, b = x, 0, c
	case h >= 300 && h < 360:
		r, g, b = c, 0, x
	}

	r = (r + m) * 255
	g = (g + m) * 255
	b = (b + m) * 255

	return r, g, b
}

// Convenience functions for creating common filter combinations

// Instagram creates an Instagram-like filter
func Instagram() *FilterChain {
	return NewFilterChain().
		Add(ContrastFilter{Amount: 1.1}).
		Add(SaturateFilter{Amount: 1.2}).
		Add(BrightnessFilter{Amount: 1.05})
}

// Vintage creates a vintage photo effect
func Vintage() *FilterChain {
	return NewFilterChain().
		Add(SepiaFilter{Amount: 0.8}).
		Add(ContrastFilter{Amount: 1.2}).
		Add(BrightnessFilter{Amount: 0.9})
}

// Dramatic creates a dramatic effect
func Dramatic() *FilterChain {
	return NewFilterChain().
		Add(ContrastFilter{Amount: 1.5}).
		Add(SaturateFilter{Amount: 1.3}).
		Add(BrightnessFilter{Amount: 0.95})
}

// BlackAndWhite creates a black and white effect
func BlackAndWhite() *FilterChain {
	return NewFilterChain().
		Add(SaturateFilter{Amount: 0}).
		Add(ContrastFilter{Amount: 1.1})
}

// Warm creates a warm tone effect
func Warm() *FilterChain {
	return NewFilterChain().
		Add(HueRotateFilter{Degrees: 15}).
		Add(SaturateFilter{Amount: 1.1}).
		Add(BrightnessFilter{Amount: 1.05})
}

// Cool creates a cool tone effect
func Cool() *FilterChain {
	return NewFilterChain().
		Add(HueRotateFilter{Degrees: -15}).
		Add(SaturateFilter{Amount: 1.1}).
		Add(BrightnessFilter{Amount: 1.05})
}
