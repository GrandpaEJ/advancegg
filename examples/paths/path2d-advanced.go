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
	
	// Create a complex path with multiple shapes
	complexPath := advancegg.NewPath2D()
	
	// Add a rounded rectangle using arcs
	createRoundedRect(complexPath, 50, 50, 200, 100, 20)
	
	// Add a heart shape
	createHeart(complexPath, 400, 100, 50)
	
	// Add a gear shape
	createGear(complexPath, 650, 100, 40, 8)
	
	// Fill the complex path with gradient effect (simulated with multiple colors)
	dc.SetRGB(0.2, 0.4, 0.8)
	dc.FillPath2D(complexPath)
	
	// Create a path for clipping
	clipPath := advancegg.NewPath2D()
	clipPath.Arc(400, 400, 150, 0, 2*math.Pi, false)
	
	// Use the clip path
	dc.ClipPath2D(clipPath)
	
	// Draw something that will be clipped
	for i := 0; i < 20; i++ {
		linePath := advancegg.NewPath2D()
		x := 250 + float64(i)*10
		linePath.MoveTo(x, 250)
		linePath.LineTo(x, 550)
		
		// Alternate colors
		if i%2 == 0 {
			dc.SetRGB(1, 0, 0)
		} else {
			dc.SetRGB(0, 0, 1)
		}
		dc.SetLineWidth(8)
		dc.StrokePath2D(linePath)
	}
	
	// Save the image
	dc.SavePNG("images/paths/path2d-advanced.png")
}

func createRoundedRect(path *advancegg.Path2D, x, y, width, height, radius float64) {
	path.MoveTo(x+radius, y)
	path.LineTo(x+width-radius, y)
	path.Arc(x+width-radius, y+radius, radius, -math.Pi/2, 0, false)
	path.LineTo(x+width, y+height-radius)
	path.Arc(x+width-radius, y+height-radius, radius, 0, math.Pi/2, false)
	path.LineTo(x+radius, y+height)
	path.Arc(x+radius, y+height-radius, radius, math.Pi/2, math.Pi, false)
	path.LineTo(x, y+radius)
	path.Arc(x+radius, y+radius, radius, math.Pi, 3*math.Pi/2, false)
	path.ClosePath()
}

func createHeart(path *advancegg.Path2D, x, y, size float64) {
	// Heart shape using curves
	path.MoveTo(x, y+size/4)
	path.BezierCurveTo(x, y-size/2, x-size, y-size/2, x-size/2, y)
	path.BezierCurveTo(x-size, y+size/2, x, y+size, x, y+size)
	path.BezierCurveTo(x, y+size, x+size, y+size/2, x+size/2, y)
	path.BezierCurveTo(x+size, y-size/2, x, y-size/2, x, y+size/4)
	path.ClosePath()
}

func createGear(path *advancegg.Path2D, x, y, radius float64, teeth int) {
	innerRadius := radius * 0.7
	toothHeight := radius * 0.3
	
	for i := 0; i < teeth*2; i++ {
		angle := float64(i) * math.Pi / float64(teeth)
		var r float64
		if i%2 == 0 {
			r = radius + toothHeight
		} else {
			r = radius
		}
		
		px := x + r*math.Cos(angle)
		py := y + r*math.Sin(angle)
		
		if i == 0 {
			path.MoveTo(px, py)
		} else {
			path.LineTo(px, py)
		}
	}
	path.ClosePath()
	
	// Add inner circle (hole)
	path.Arc(x, y, innerRadius, 0, 2*math.Pi, false)
}
