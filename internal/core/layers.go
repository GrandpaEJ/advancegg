package core

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"math"
	"time"
)

// Layer system for multi-layered drawing like Photoshop

// BlendMode represents different layer blending modes
type BlendMode int

const (
	BlendModeNormal BlendMode = iota
	// Photoshop-style Blend Modes
	BlendModeMultiply
	BlendModeScreen
	BlendModeOverlay
	BlendModeSoftLight
	BlendModeHardLight
	BlendModeColorDodge
	BlendModeColorBurn
	BlendModeDarken
	BlendModeLighten
	BlendModeDifference
	BlendModeExclusion
	BlendModeHue
	BlendModeSaturation
	BlendModeColor
	BlendModeLuminosity

	// Porter-Duff Compositing Operators
	BlendModeClear
	BlendModeSource
	BlendModeDest
	BlendModeSrcOver // Alias for Normal but explicit
	BlendModeDstOver
	BlendModeSrcIn
	BlendModeDstIn
	BlendModeSrcOut
	BlendModeDstOut
	BlendModeSrcAtop
	BlendModeDstAtop
	BlendModeXor
	BlendModeAdd
)

// Layer represents a single drawing layer
type Layer struct {
	Name      string
	Image     *image.RGBA
	Opacity   float64
	BlendMode BlendMode
	Visible   bool
	Locked    bool
	Mask      *image.Alpha
	Transform Matrix
	ID        string
}

// LayerManager manages multiple layers
type LayerManager struct {
	Layers      []*Layer
	ActiveLayer int
	Width       int
	Height      int
	Background  color.Color
}

// NewLayerManager creates a new layer manager
func NewLayerManager(width, height int) *LayerManager {
	lm := &LayerManager{
		Layers:      make([]*Layer, 0),
		ActiveLayer: -1,
		Width:       width,
		Height:      height,
		Background:  color.RGBA{255, 255, 255, 255}, // White background
	}

	// Create default background layer
	lm.AddLayer("Background")
	return lm
}

// NewLayer creates a new layer
func NewLayer(name string, width, height int) *Layer {
	return &Layer{
		Name:      name,
		Image:     image.NewRGBA(image.Rect(0, 0, width, height)),
		Opacity:   1.0,
		BlendMode: BlendModeNormal,
		Visible:   true,
		Locked:    false,
		Transform: Identity(),
		ID:        generateLayerID(),
	}
}

// AddLayer adds a new layer
func (lm *LayerManager) AddLayer(name string) *Layer {
	layer := NewLayer(name, lm.Width, lm.Height)
	lm.Layers = append(lm.Layers, layer)
	lm.ActiveLayer = len(lm.Layers) - 1
	return layer
}

// InsertLayer inserts a layer at the specified index
func (lm *LayerManager) InsertLayer(index int, name string) *Layer {
	if index < 0 || index > len(lm.Layers) {
		return lm.AddLayer(name)
	}

	layer := NewLayer(name, lm.Width, lm.Height)

	// Insert at index
	lm.Layers = append(lm.Layers[:index], append([]*Layer{layer}, lm.Layers[index:]...)...)
	lm.ActiveLayer = index

	return layer
}

// RemoveLayer removes a layer by index
func (lm *LayerManager) RemoveLayer(index int) bool {
	if index < 0 || index >= len(lm.Layers) || len(lm.Layers) <= 1 {
		return false // Can't remove last layer
	}

	lm.Layers = append(lm.Layers[:index], lm.Layers[index+1:]...)

	// Adjust active layer
	if lm.ActiveLayer >= len(lm.Layers) {
		lm.ActiveLayer = len(lm.Layers) - 1
	}

	return true
}

// MoveLayer moves a layer from one index to another
func (lm *LayerManager) MoveLayer(from, to int) bool {
	if from < 0 || from >= len(lm.Layers) || to < 0 || to >= len(lm.Layers) {
		return false
	}

	layer := lm.Layers[from]

	// Remove from old position
	lm.Layers = append(lm.Layers[:from], lm.Layers[from+1:]...)

	// Insert at new position
	if to > from {
		to-- // Adjust for removal
	}
	lm.Layers = append(lm.Layers[:to], append([]*Layer{layer}, lm.Layers[to:]...)...)

	lm.ActiveLayer = to
	return true
}

// GetActiveLayer returns the currently active layer
func (lm *LayerManager) GetActiveLayer() *Layer {
	if lm.ActiveLayer < 0 || lm.ActiveLayer >= len(lm.Layers) {
		return nil
	}
	return lm.Layers[lm.ActiveLayer]
}

// SetActiveLayer sets the active layer by index
func (lm *LayerManager) SetActiveLayer(index int) bool {
	if index < 0 || index >= len(lm.Layers) {
		return false
	}
	lm.ActiveLayer = index
	return true
}

// SetActiveLayerByName sets the active layer by name
func (lm *LayerManager) SetActiveLayerByName(name string) bool {
	for i, layer := range lm.Layers {
		if layer.Name == name {
			lm.ActiveLayer = i
			return true
		}
	}
	return false
}

// DuplicateLayer duplicates a layer
func (lm *LayerManager) DuplicateLayer(index int) *Layer {
	if index < 0 || index >= len(lm.Layers) {
		return nil
	}

	original := lm.Layers[index]
	duplicate := NewLayer(original.Name+" Copy", lm.Width, lm.Height)

	// Copy properties
	duplicate.Opacity = original.Opacity
	duplicate.BlendMode = original.BlendMode
	duplicate.Visible = original.Visible
	duplicate.Transform = original.Transform

	// Copy image data
	draw.Draw(duplicate.Image, duplicate.Image.Bounds(), original.Image, image.Point{}, draw.Src)

	// Copy mask if exists
	if original.Mask != nil {
		duplicate.Mask = image.NewAlpha(original.Mask.Bounds())
		draw.Draw(duplicate.Mask, duplicate.Mask.Bounds(), original.Mask, image.Point{}, draw.Src)
	}

	// Insert after original
	lm.Layers = append(lm.Layers[:index+1], append([]*Layer{duplicate}, lm.Layers[index+1:]...)...)
	lm.ActiveLayer = index + 1

	return duplicate
}

// Composite renders all layers into a single image
func (lm *LayerManager) Composite() *image.RGBA {
	result := image.NewRGBA(image.Rect(0, 0, lm.Width, lm.Height))

	// Fill with background color
	draw.Draw(result, result.Bounds(), &image.Uniform{lm.Background}, image.Point{}, draw.Src)

	// Composite layers from bottom to top
	for _, layer := range lm.Layers {
		if !layer.Visible {
			continue
		}

		lm.compositeLayer(result, layer)
	}

	return result
}

// compositeLayer composites a single layer onto the result
func (lm *LayerManager) compositeLayer(dst *image.RGBA, layer *Layer) {
	bounds := dst.Bounds()
	mode := layer.BlendMode

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			// Get source pixel
			srcPixel := layer.Image.RGBAAt(x, y)

			// Apply layer opacity
			if layer.Opacity < 1.0 {
				srcPixel.A = uint8(float64(srcPixel.A) * layer.Opacity)
				// Pre-multiply RGB by new Alpha if strictly needed, but BlendMode logic handles it
				// For correctness in Composition, we might need PREMULTIPLIED alpha?
				// Most Go image libraries work with NON-premultiplied RGBA struct, but math is often easier with premultiplied
				// However, existing simple alphaBlend logic below assumed non-premultiplied.
				// We will scale alpha here. RGB remains "straight" (associated with alpha).
			}

			// Apply mask if exists
			if layer.Mask != nil {
				maskAlpha := float64(layer.Mask.AlphaAt(x, y).A) / 255.0
				srcPixel.A = uint8(float64(srcPixel.A) * maskAlpha)
			}

			// Skip completely transparent pixels ONLY if we are in a mode where source 0 means no change
			// For modes like Clear, SrcIn, etc., transparent source DOES affect the output
			if srcPixel.A == 0 && isStandardBlendMode(mode) {
				continue
			}

			// Get destination pixel
			dstPixel := dst.RGBAAt(x, y)

			var finalPixel color.RGBA

			// Handle Compositing Operators (Porter-Duff) vs Standard Blend Modes
			if isCompositingOperator(mode) {
				finalPixel = applyCompositingOperator(srcPixel, dstPixel, mode)
			} else {
				// Standard Photoshop-style Blend Modes (Source Over composition)
				// 1. Blend Colors
				blended := lm.applyBlendMode(srcPixel, dstPixel, mode)
				// 2. Alpha Composite (SrcOver)
				finalPixel = alphaBlend(blended, dstPixel, float64(srcPixel.A)/255.0)
			}

			dst.SetRGBA(x, y, finalPixel)
		}
	}
}

func isStandardBlendMode(mode BlendMode) bool {
	return mode <= BlendModeLuminosity
}

func isCompositingOperator(mode BlendMode) bool {
	return mode >= BlendModeClear
}

// applyCompositingOperator implements Porter-Duff compositing rules
func applyCompositingOperator(src, dst color.RGBA, mode BlendMode) color.RGBA {
	sR, sG, sB, sA := float64(src.R), float64(src.G), float64(src.B), float64(src.A)/255.0
	dR, dG, dB, dA := float64(dst.R), float64(dst.G), float64(dst.B), float64(dst.A)/255.0

	// Helper for byte conversion
	toByte := func(v float64) uint8 {
		if v < 0 {
			return 0
		}
		if v > 255 {
			return 255
		}
		return uint8(v)
	}

	switch mode {
	case BlendModeClear:
		return color.RGBA{0, 0, 0, 0}

	case BlendModeSource:
		return src

	case BlendModeDest:
		return dst

	case BlendModeSrcOver, BlendModeNormal: // Use standard formula
		// out = src + dst(1-sa)
		outA := sA + dA*(1.0-sA)
		if outA == 0 {
			return color.RGBA{0, 0, 0, 0}
		}
		return color.RGBA{
			R: toByte((sR*sA + dR*dA*(1.0-sA)) / outA),
			G: toByte((sG*sA + dG*dA*(1.0-sA)) / outA),
			B: toByte((sB*sA + dB*dA*(1.0-sA)) / outA),
			A: toByte(outA * 255.0),
		}

	case BlendModeDstOver:
		// out = dst + src(1-da)
		outA := dA + sA*(1.0-dA)
		if outA == 0 {
			return color.RGBA{0, 0, 0, 0}
		}
		return color.RGBA{
			R: toByte((dR*dA + sR*sA*(1.0-dA)) / outA),
			G: toByte((dG*dA + sG*sA*(1.0-dA)) / outA),
			B: toByte((dB*dA + sB*sA*(1.0-dA)) / outA),
			A: toByte(outA * 255.0),
		}

	case BlendModeSrcIn:
		// out = src * da
		return color.RGBA{
			R: toByte(sR * dA),
			G: toByte(sG * dA),
			B: toByte(sB * dA),
			A: toByte(sA * 255.0 * dA),
		}

	case BlendModeDstIn:
		// out = dst * sa
		return color.RGBA{
			R: toByte(dR * sA),
			G: toByte(dG * sA),
			B: toByte(dB * sA),
			A: toByte(dA * 255.0 * sA),
		}

	case BlendModeSrcOut:
		// out = src * (1-da)
		f := 1.0 - dA
		return color.RGBA{
			R: toByte(sR * f),
			G: toByte(sG * f),
			B: toByte(sB * f),
			A: toByte(sA * 255.0 * f),
		}

	case BlendModeDstOut:
		// out = dst * (1-sa)
		f := 1.0 - sA
		return color.RGBA{
			R: toByte(dR * f),
			G: toByte(dG * f),
			B: toByte(dB * f),
			A: toByte(dA * 255.0 * f),
		}

	case BlendModeSrcAtop:
		// out = src*da + dst*(1-sa)
		// r = s*da + d*(1-sa)
		// a = s*da + d*(1-sa) = da(s + 1 - sa)? No.
		// A = sa*da + da*(1-sa) = da(sa + 1 - sa) = da. Correct.
		// Result Alpha is always dA.
		return color.RGBA{
			R: toByte(sR*dA + dR*(1.0-sA)),
			G: toByte(sG*dA + dG*(1.0-sA)),
			B: toByte(sB*dA + dB*(1.0-sA)),
			A: toByte(dA * 255.0),
		}

	case BlendModeDstAtop:
		// out = dst*sa + src*(1-da)
		// Result Alpha is always sA
		return color.RGBA{
			R: toByte(dR*sA + sR*(1.0-dA)),
			G: toByte(dG*sA + sG*(1.0-dA)),
			B: toByte(dB*sA + sB*(1.0-dA)),
			A: toByte(sA * 255.0),
		}

	case BlendModeXor:
		// out = src*(1-da) + dst*(1-sa)
		outA := sA*(1.0-dA) + dA*(1.0-sA)
		if outA == 0 {
			return color.RGBA{0, 0, 0, 0}
		}
		// Colors roughly premultiplied in formula:
		// R = (sR*sA)*(1-dA) + (dR*dA)*(1-sA)
		// We have sR, dR as STRAIGHT colors.
		// Formula using straight colors:
		// C = (Cs * As * (1-Ad) + Cd * Ad * (1-As)) / Aout

		r := (sR*sA*(1.0-dA) + dR*dA*(1.0-sA)) / outA
		g := (sG*sA*(1.0-dA) + dG*dA*(1.0-sA)) / outA
		b := (sB*sA*(1.0-dA) + dB*dA*(1.0-sA)) / outA

		return color.RGBA{
			R: toByte(r),
			G: toByte(g),
			B: toByte(b),
			A: toByte(outA * 255.0),
		}

	case BlendModeAdd:
		// Simple addition limited to 255
		// R = min(255, sR + dR)
		// A = min(255, sA + dA)
		return color.RGBA{
			R: toByte(math.Min(255, sR+dR)),
			G: toByte(math.Min(255, sG+dG)),
			B: toByte(math.Min(255, sB+dB)),
			A: toByte(math.Min(255, float64(src.A)+float64(dst.A))),
		}
	}

	return src // Fallback
}

// applyBlendMode applies the specified blend mode
func (lm *LayerManager) applyBlendMode(src, dst color.RGBA, mode BlendMode) color.RGBA {
	switch mode {
	case BlendModeNormal:
		return src
	case BlendModeMultiply:
		return color.RGBA{
			R: uint8(int(src.R) * int(dst.R) / 255),
			G: uint8(int(src.G) * int(dst.G) / 255),
			B: uint8(int(src.B) * int(dst.B) / 255),
			A: src.A,
		}
	case BlendModeScreen:
		return color.RGBA{
			R: uint8(255 - (255-int(src.R))*(255-int(dst.R))/255),
			G: uint8(255 - (255-int(src.G))*(255-int(dst.G))/255),
			B: uint8(255 - (255-int(src.B))*(255-int(dst.B))/255),
			A: src.A,
		}
	case BlendModeOverlay:
		return color.RGBA{
			R: overlayBlend(src.R, dst.R),
			G: overlayBlend(src.G, dst.G),
			B: overlayBlend(src.B, dst.B),
			A: src.A,
		}
	case BlendModeSoftLight:
		return color.RGBA{
			R: softLightBlend(src.R, dst.R),
			G: softLightBlend(src.G, dst.G),
			B: softLightBlend(src.B, dst.B),
			A: src.A,
		}
	case BlendModeHardLight:
		return color.RGBA{
			R: hardLightBlend(src.R, dst.R),
			G: hardLightBlend(src.G, dst.G),
			B: hardLightBlend(src.B, dst.B),
			A: src.A,
		}
	case BlendModeColorDodge:
		return color.RGBA{
			R: colorDodgeBlend(src.R, dst.R),
			G: colorDodgeBlend(src.G, dst.G),
			B: colorDodgeBlend(src.B, dst.B),
			A: src.A,
		}
	case BlendModeColorBurn:
		return color.RGBA{
			R: colorBurnBlend(src.R, dst.R),
			G: colorBurnBlend(src.G, dst.G),
			B: colorBurnBlend(src.B, dst.B),
			A: src.A,
		}
	case BlendModeDarken:
		return color.RGBA{
			R: minUint8(src.R, dst.R),
			G: minUint8(src.G, dst.G),
			B: minUint8(src.B, dst.B),
			A: src.A,
		}
	case BlendModeLighten:
		return color.RGBA{
			R: maxUint8(src.R, dst.R),
			G: maxUint8(src.G, dst.G),
			B: maxUint8(src.B, dst.B),
			A: src.A,
		}
	case BlendModeDifference:
		return color.RGBA{
			R: uint8(absInt(int(src.R) - int(dst.R))),
			G: uint8(absInt(int(src.G) - int(dst.G))),
			B: uint8(absInt(int(src.B) - int(dst.B))),
			A: src.A,
		}
	case BlendModeExclusion:
		return color.RGBA{
			R: uint8(int(src.R) + int(dst.R) - 2*int(src.R)*int(dst.R)/255),
			G: uint8(int(src.G) + int(dst.G) - 2*int(src.G)*int(dst.G)/255),
			B: uint8(int(src.B) + int(dst.B) - 2*int(src.B)*int(dst.B)/255),
			A: src.A,
		}
	default:
		return src
	}
}

// Helper functions for blend modes
func overlayBlend(src, dst uint8) uint8 {
	if dst < 128 {
		return uint8(2 * int(src) * int(dst) / 255)
	}
	return uint8(255 - 2*(255-int(src))*(255-int(dst))/255)
}

func softLightBlend(src, dst uint8) uint8 {
	s := float64(src) / 255.0
	d := float64(dst) / 255.0

	var result float64
	if s <= 0.5 {
		result = d - (1-2*s)*d*(1-d)
	} else {
		var g float64
		if d <= 0.25 {
			g = ((16*d-12)*d + 4) * d
		} else {
			g = math.Sqrt(d)
		}
		result = d + (2*s-1)*(g-d)
	}

	return uint8(clampFloat(result*255, 0, 255))
}

func hardLightBlend(src, dst uint8) uint8 {
	if src < 128 {
		return uint8(2 * int(src) * int(dst) / 255)
	}
	return uint8(255 - 2*(255-int(src))*(255-int(dst))/255)
}

func colorDodgeBlend(src, dst uint8) uint8 {
	if src == 255 {
		return 255
	}
	result := int(dst) * 255 / (255 - int(src))
	if result > 255 {
		return 255
	}
	return uint8(result)
}

func colorBurnBlend(src, dst uint8) uint8 {
	if src == 0 {
		return 0
	}
	result := 255 - (255-int(dst))*255/int(src)
	if result < 0 {
		return 0
	}
	return uint8(result)
}

func minUint8(a, b uint8) uint8 {
	if a < b {
		return a
	}
	return b
}

func maxUint8(a, b uint8) uint8 {
	if a > b {
		return a
	}
	return b
}

func absInt(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func clampFloat(value, min, max float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

func alphaBlend(src, dst color.RGBA, srcAlpha float64) color.RGBA {
	// Standard Source-Over Composition:
	// outA = srcA + dstA*(1-srcA)
	// outC = (srcC*srcA + dstC*dstA*(1-srcA)) / outA

	sA := srcAlpha
	dA := float64(dst.A) / 255.0
	outA := sA + dA*(1.0-sA)

	if outA == 0 {
		return color.RGBA{0, 0, 0, 0}
	}

	resultAlpha := uint8(outA * 255.0)

	// Colors
	sR, sG, sB := float64(src.R), float64(src.G), float64(src.B)
	dR, dG, dB := float64(dst.R), float64(dst.G), float64(dst.B)

	r := (sR*sA + dR*dA*(1.0-sA)) / outA
	g := (sG*sA + dG*dA*(1.0-sA)) / outA
	b := (sB*sA + dB*dA*(1.0-sA)) / outA

	toByte := func(v float64) uint8 {
		if v < 0 {
			return 0
		}
		if v > 255 {
			return 255
		}
		return uint8(v)
	}

	return color.RGBA{toByte(r), toByte(g), toByte(b), resultAlpha}
}

// generateLayerID generates a unique layer ID
func generateLayerID() string {
	// Simple ID generation - in production, use UUID
	return fmt.Sprintf("layer_%d", time.Now().UnixNano())
}

// Layer utility methods

// Clear clears the layer
func (l *Layer) Clear() {
	draw.Draw(l.Image, l.Image.Bounds(), &image.Uniform{color.Transparent}, image.Point{}, draw.Src)
}

// Fill fills the layer with a color
func (l *Layer) Fill(c color.Color) {
	draw.Draw(l.Image, l.Image.Bounds(), &image.Uniform{c}, image.Point{}, draw.Src)
}

// SetOpacity sets the layer opacity
func (l *Layer) SetOpacity(opacity float64) {
	if opacity < 0 {
		opacity = 0
	}
	if opacity > 1 {
		opacity = 1
	}
	l.Opacity = opacity
}

// SetBlendMode sets the layer blend mode
func (l *Layer) SetBlendMode(mode BlendMode) {
	l.BlendMode = mode
}

// SetVisible sets layer visibility
func (l *Layer) SetVisible(visible bool) {
	l.Visible = visible
}

// SetLocked sets layer lock state
func (l *Layer) SetLocked(locked bool) {
	l.Locked = locked
}

// AddMask adds a layer mask
func (l *Layer) AddMask() {
	l.Mask = image.NewAlpha(l.Image.Bounds())
	// Fill with white (fully visible)
	draw.Draw(l.Mask, l.Mask.Bounds(), &image.Uniform{color.Alpha{255}}, image.Point{}, draw.Src)
}

// RemoveMask removes the layer mask
func (l *Layer) RemoveMask() {
	l.Mask = nil
}
