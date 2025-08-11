package main

import "github.com/GrandpaEJ/advancegg"

func main() {
	const S = 1024
	dc := advancegg.NewContext(S, S)
	dc.SetRGB(1, 1, 1)
	dc.Clear()
	// Try to load a bold system font, fallback to default if not found
	fontPaths := []string{
		"/usr/share/fonts/truetype/dejavu/DejaVuSans-Bold.ttf",
		"/usr/share/fonts/truetype/liberation/LiberationSans-Bold.ttf",
		"/usr/share/fonts/TTF/DejaVuSans-Bold.ttf",
		"/System/Library/Fonts/Impact.ttf",                // macOS
		"/Library/Fonts/Impact.ttf",                       // macOS
		"C:/Windows/Fonts/impact.ttf",                     // Windows
		"/usr/share/fonts/truetype/dejavu/DejaVuSans.ttf", // fallback to regular
	}

	var err error
	for _, fontPath := range fontPaths {
		err = dc.LoadFontFace(fontPath, 96)
		if err == nil {
			break
		}
	}
	// Continue with default font if none found
	dc.SetRGB(0, 0, 0)
	s := "ONE DOES NOT SIMPLY"
	n := 6 // "stroke" size
	for dy := -n; dy <= n; dy++ {
		for dx := -n; dx <= n; dx++ {
			if dx*dx+dy*dy >= n*n {
				// give it rounded corners
				continue
			}
			x := S/2 + float64(dx)
			y := S/2 + float64(dy)
			dc.DrawStringAnchored(s, x, y, 0.5, 0.5)
		}
	}
	dc.SetRGB(1, 1, 1)
	dc.DrawStringAnchored(s, S/2, S/2, 0.5, 0.5)
	dc.SavePNG("out.png")
}
