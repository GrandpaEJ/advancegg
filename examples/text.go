package main

import "github.com/GrandpaEJ/advancegg"

func main() {
	const S = 1024
	dc := advancegg.NewContext(S, S)
	dc.SetRGB(1, 1, 1)
	dc.Clear()
	dc.SetRGB(0, 0, 0)
	// Try to load a system font, fallback to default if not found
	fontPaths := []string{
		"/usr/share/fonts/truetype/dejavu/DejaVuSans.ttf",
		"/usr/share/fonts/truetype/liberation/LiberationSans-Regular.ttf",
		"/usr/share/fonts/TTF/DejaVuSans.ttf",
		"/System/Library/Fonts/Arial.ttf", // macOS
		"/Library/Fonts/Arial.ttf",        // macOS
		"C:/Windows/Fonts/arial.ttf",      // Windows
	}

	var err error
	for _, fontPath := range fontPaths {
		err = dc.LoadFontFace(fontPath, 96)
		if err == nil {
			break
		}
	}
	if err != nil {
		// Use default font if no system font found - just continue with default
		// The default font will be used automatically
	}
	dc.DrawStringAnchored("Hello, world!", S/2, S/2, 0.5, 0.5)
	dc.SavePNG("out.png")
}
