package core

import (
	"image/color"
	"testing"
)

func TestPorterDuffBlending(t *testing.T) {
	// Setup colors
	red := color.RGBA{255, 0, 0, 255}
	blue := color.RGBA{0, 0, 255, 255}
	transparent := color.RGBA{0, 0, 0, 0}

	// Ensure we test partial alpha too
	// semiRed := color.RGBA{255, 0, 0, 128}  // ~0.5 alpha
	// semiBlue := color.RGBA{0, 0, 255, 128} // ~0.5 alpha

	tests := []struct {
		name     string
		mode     BlendMode
		src      color.RGBA
		dst      color.RGBA
		expected color.RGBA
	}{
		// Clear: Result should be transparent
		{"Clear", BlendModeClear, red, blue, transparent},

		// Source: Result should be Source
		{"Source", BlendModeSource, red, blue, red},

		// Dest: Result should be Dest
		{"Dest", BlendModeDest, red, blue, blue},

		// SrcIn: Src where Dst exists (Alpha logic)
		// red (255) in blue (255) -> red
		{"SrcIn Full Overlap", BlendModeSrcIn, red, blue, red},
		// red (255) in transparent (0) -> transparent
		{"SrcIn No Overlap", BlendModeSrcIn, red, transparent, transparent},

		// SrcOut: Src where Dst does NOT exist
		// red (255) out blue (255) -> transparent
		{"SrcOut Full Overlap", BlendModeSrcOut, red, blue, transparent},
		// red (255) out transparent (0) -> red
		{"SrcOut No Overlap", BlendModeSrcOut, red, transparent, red},

		// DstIn: Dst where Src exists
		// blue (255) in red (255) -> blue
		{"DstIn Full Overlap", BlendModeDstIn, red, blue, blue},

		// DstOut: Dst where Src does NOT exist
		// blue (255) out red (255) -> transparent
		{"DstOut Full Overlap", BlendModeDstOut, red, blue, transparent},

		// SrcAtop: Src over Dst, but only where Dst exists
		// red atop blue -> red (because blue is opaque)
		{"SrcAtop Opaque", BlendModeSrcAtop, red, blue, red},
		// red atop transparent -> transparent (because dest is missing)
		{"SrcAtop Transparent", BlendModeSrcAtop, red, transparent, transparent},

		// DstAtop: Dst over Src, but only where Src exists
		// blue atop red -> blue
		{"DstAtop Opaque", BlendModeDstAtop, red, blue, blue},

		// Xor: Src where Dst missing + Dst where Src missing
		// red xor blue (both opaque) -> transparent
		{"Xor Opaque", BlendModeXor, red, blue, transparent},
		// red xor transparent -> red
		{"Xor Src Only", BlendModeXor, red, transparent, red},
		// transparent xor blue -> blue
		{"Xor Dst Only", BlendModeXor, transparent, blue, blue},

		// Add: Src + Dst
		// red + blue = magenta
		{"Add", BlendModeAdd, red, blue, color.RGBA{255, 0, 255, 255}},
	}

	// Helper to check approximate equality for colors
	// We allow +/- 1 diff for rounding errors in blending math
	colorsEqual := func(c1, c2 color.RGBA) bool {
		diff := func(a, b uint8) int {
			d := int(a) - int(b)
			if d < 0 {
				return -d
			}
			return d
		}
		return diff(c1.R, c2.R) <= 1 && diff(c1.G, c2.G) <= 1 && diff(c1.B, c2.B) <= 1 && diff(c1.A, c2.A) <= 1
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// We can't directly call internal compositeLayer easily without setup,
			// but we can call applyCompositingOperator if we export it or use it here.
			// Since applyCompositingOperator is unexported in layers.go, we should test via LayerManager.

			lm := NewLayerManager(1, 1)
			lm.Background = color.Transparent // Ensure buffer starts transparent

			// Set background to destination color
			bgLayer := lm.GetActiveLayer()
			bgLayer.Fill(tt.dst)

			// Add source layer
			l := lm.AddLayer("Src")
			l.Fill(tt.src)
			l.SetBlendMode(tt.mode)

			// Composite
			resultRaw := lm.Composite()
			result := resultRaw.RGBAAt(0, 0)

			if !colorsEqual(result, tt.expected) {
				t.Errorf("Mode %v: expected %v, got %v", tt.name, tt.expected, result)
			}
		})
	}
}

func TestSemiTransparentBlending(t *testing.T) {
	// Test SrcOver with semi-transparent colors
	src := color.RGBA{255, 0, 0, 128} // 50% Red
	dst := color.RGBA{0, 0, 255, 128} // 50% Blue

	// Math:
	// sA = 0.5, dA = 0.5
	// outA = sA + dA*(1-sA) = 0.5 + 0.5*0.5 = 0.75 (191)
	// R = (sR*sA + dR*dA*(1-sA)) / outA
	// R = (255*0.5 + 0) / 0.75 = 127.5 / 0.75 = 170
	// G = 0
	// B = (0 + 255*0.5*(1-0.5)) / 0.75 = (255*0.25) / 0.75 = 63.75 / 0.75 = 85

	expected := color.RGBA{170, 0, 85, 191}

	lm := NewLayerManager(1, 1)
	lm.Layers[0].Fill(dst) // Background
	// Note: LayerManager init creates a white opaque background by default logic "Background" layer?
	// NewLayerManager init sets Background color to White (Opaque).
	// But `lm.AddLayer("Background")` is a layer.
	// `lm.Composite()` fills with `lm.Background` color first.
	// So we need to ensure the base canvas is CLEARED or handled.

	lm.Background = color.Transparent // Transparent base

	lm.Layers[0].Fill(dst) // First layer is dest

	l := lm.AddLayer("Src")
	l.Fill(src)
	l.SetBlendMode(BlendModeSrcOver)

	resultRaw := lm.Composite()
	result := resultRaw.RGBAAt(0, 0)

	// Allow margin of error
	diff := func(a, b uint8) int {
		d := int(a) - int(b)
		if d < 0 {
			return -d
		}
		return d
	}

	if diff(result.R, expected.R) > 2 || diff(result.B, expected.B) > 2 || diff(result.A, expected.A) > 2 {
		t.Errorf("SemiTransparent SrcOver: expected ~%v, got %v", expected, result)
	}
}
