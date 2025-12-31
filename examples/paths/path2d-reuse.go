package main

import (
	"math"

	"github.com/GrandpaEJ/advancegg"
)

func main() {
	dc := advancegg.NewContext(800, 600)

	// White background
	dc.SetRGB(1, 1, 1)
	dc.Clear()

	// Create a star path that we'll reuse
	starPath := createStarPath()

	// Draw the star multiple times in different colors and positions
	positions := []struct{ x, y float64 }{
		{150, 150},
		{400, 150},
		{650, 150},
		{275, 350},
		{525, 350},
	}

	colors := []struct{ r, g, b float64 }{
		{1, 0, 0}, // Red
		{0, 1, 0}, // Green
		{0, 0, 1}, // Blue
		{1, 1, 0}, // Yellow
		{1, 0, 1}, // Magenta
	}

	for i, pos := range positions {
		dc.Push()
		dc.Translate(pos.x, pos.y)
		dc.SetRGB(colors[i].r, colors[i].g, colors[i].b)
		dc.FillPath2D(starPath)
		dc.Pop()
	}

	// Create a path by combining multiple paths
	combinedPath := advancegg.NewPath2D()

	// Add a circle
	circlePath := advancegg.NewPath2D()
	circlePath.Arc(400, 450, 30, 0, 6.28318, false)
	combinedPath.AddPath(circlePath)

	// Add a square
	squarePath := advancegg.NewPath2D()
	squarePath.Rect(350, 400, 60, 60)
	combinedPath.AddPath(squarePath)

	// Stroke the combined path
	dc.SetRGB(0, 0, 0)
	dc.SetLineWidth(2)
	dc.StrokePath2D(combinedPath)

	// Save the image
	dc.SavePNG("images/paths/path2d-reuse.png")
}

func createStarPath() *advancegg.Path2D {
	path := advancegg.NewPath2D()

	// Create a 5-pointed star
	centerX, centerY := 0.0, 0.0
	outerRadius := 50.0
	innerRadius := 20.0

	for i := 0; i < 10; i++ {
		angle := float64(i) * 3.14159 / 5
		var radius float64
		if i%2 == 0 {
			radius = outerRadius
		} else {
			radius = innerRadius
		}

		x := centerX + radius*math.Cos(angle-3.14159/2)
		y := centerY + radius*math.Sin(angle-3.14159/2)

		if i == 0 {
			path.MoveTo(x, y)
		} else {
			path.LineTo(x, y)
		}
	}
	path.ClosePath()

	return path
}
