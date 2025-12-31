package main

import (
	"math/rand"

	"github.com/GrandpaEJ/advancegg"
)

func CreatePoints(n int) []advancegg.Point {
	points := make([]advancegg.Point, n)
	for i := 0; i < n; i++ {
		x := 0.5 + rand.NormFloat64()*0.1
		y := x + rand.NormFloat64()*0.1
		points[i] = advancegg.Point{X: x, Y: y}
	}
	return points
}

func main() {
	const S = 1024
	const P = 64
	dc := advancegg.NewContext(S, S)
	dc.InvertY()
	dc.SetRGB(1, 1, 1)
	dc.Clear()
	points := CreatePoints(1000)
	dc.Translate(P, P)
	dc.Scale(S-P*2, S-P*2)
	// draw minor grid
	for i := 1; i <= 10; i++ {
		x := float64(i) / 10
		dc.MoveTo(x, 0)
		dc.LineTo(x, 1)
		dc.MoveTo(0, x)
		dc.LineTo(1, x)
	}
	dc.SetRGBA(0, 0, 0, 0.25)
	dc.SetLineWidth(1)
	dc.Stroke()
	// draw axes
	dc.MoveTo(0, 0)
	dc.LineTo(1, 0)
	dc.MoveTo(0, 0)
	dc.LineTo(0, 1)
	dc.SetRGB(0, 0, 0)
	dc.SetLineWidth(4)
	dc.Stroke()
	// draw points
	dc.SetRGBA(0, 0, 1, 0.5)
	for _, p := range points {
		dc.DrawCircle(p.X, p.Y, 3.0/S)
		dc.Fill()
	}
	// draw text
	dc.Identity()
	dc.SetRGB(0, 0, 0)

	// Try to load bold font for title
	boldFontPaths := []string{
		"/usr/share/fonts/truetype/dejavu/DejaVuSans-Bold.ttf",
		"/usr/share/fonts/truetype/liberation/LiberationSans-Bold.ttf",
		"/System/Library/Fonts/Arial Bold.ttf", // macOS
		"/Library/Fonts/Arial Bold.ttf",        // macOS
		"C:/Windows/Fonts/arialbd.ttf",         // Windows
	}

	for _, fontPath := range boldFontPaths {
		if dc.LoadFontFace(fontPath, 24) == nil {
			break
		}
	}
	dc.DrawStringAnchored("Chart Title", S/2, P/2, 0.5, 0.5)

	// Try to load regular font for axis
	regularFontPaths := []string{
		"/usr/share/fonts/truetype/dejavu/DejaVuSans.ttf",
		"/usr/share/fonts/truetype/liberation/LiberationSans-Regular.ttf",
		"/System/Library/Fonts/Arial.ttf", // macOS
		"/Library/Fonts/Arial.ttf",        // macOS
		"C:/Windows/Fonts/arial.ttf",      // Windows
	}

	for _, fontPath := range regularFontPaths {
		if dc.LoadFontFace(fontPath, 18) == nil {
			break
		}
	}
	dc.DrawStringAnchored("X Axis Title", S/2, S-P/2, 0.5, 0.5)
	dc.SavePNG("images/misc/out.png")
}
