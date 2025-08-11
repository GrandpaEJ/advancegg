package main

import "github.com/GrandpaEJ/advancegg"

func main() {
	const S = 4096 * 2
	const T = 16 * 2
	const F = 28
	dc := advancegg.NewContext(S, S)
	dc.SetRGB(1, 1, 1)
	dc.Clear()
	dc.SetRGB(0, 0, 0)
	// Try to load a Unicode-capable font
	fontPaths := []string{
		"/usr/share/fonts/truetype/dejavu/DejaVuSans.ttf",
		"/usr/share/fonts/truetype/liberation/LiberationSans-Regular.ttf",
		"/usr/share/fonts/truetype/noto/NotoSans-Regular.ttf",
		"/System/Library/Fonts/Arial.ttf", // macOS
		"/Library/Fonts/Arial.ttf",        // macOS
		"C:/Windows/Fonts/arial.ttf",      // Windows
	}

	for _, fontPath := range fontPaths {
		if dc.LoadFontFace(fontPath, F) == nil {
			break
		}
	}
	// Continue with default font if none found
	for r := 0; r < 256; r++ {
		for c := 0; c < 256; c++ {
			i := r*256 + c
			x := float64(c*T) + T/2
			y := float64(r*T) + T/2
			dc.DrawStringAnchored(string(rune(i)), x, y, 0.5, 0.5)
		}
	}
	dc.SavePNG("out.png")
}
