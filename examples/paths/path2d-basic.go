package main

import "github.com/GrandpaEJ/advancegg"

func main() {
	dc := advancegg.NewContext(800, 600)
	
	// White background
	dc.SetRGB(1, 1, 1)
	dc.Clear()
	
	// Create a Path2D object
	path := advancegg.NewPath2D()
	
	// Draw a triangle
	path.MoveTo(100, 100)
	path.LineTo(200, 100)
	path.LineTo(150, 50)
	path.ClosePath()
	
	// Fill the path in blue
	dc.SetRGB(0, 0, 1)
	dc.FillPath2D(path)
	
	// Create another path for a rectangle
	rectPath := advancegg.NewPath2D()
	rectPath.Rect(300, 100, 100, 80)
	
	// Stroke the rectangle in red
	dc.SetRGB(1, 0, 0)
	dc.SetLineWidth(3)
	dc.StrokePath2D(rectPath)
	
	// Create a path with curves
	curvePath := advancegg.NewPath2D()
	curvePath.MoveTo(100, 300)
	curvePath.QuadraticCurveTo(200, 200, 300, 300)
	curvePath.BezierCurveTo(350, 350, 400, 250, 450, 300)
	
	// Stroke the curve in green
	dc.SetRGB(0, 1, 0)
	dc.SetLineWidth(2)
	dc.StrokePath2D(curvePath)
	
	// Create a path with an arc
	arcPath := advancegg.NewPath2D()
	arcPath.Arc(600, 200, 50, 0, 3.14159, false) // Half circle
	
	// Fill the arc in purple
	dc.SetRGB(0.5, 0, 0.5)
	dc.FillPath2D(arcPath)
	
	// Create an ellipse path
	ellipsePath := advancegg.NewPath2D()
	ellipsePath.Ellipse(600, 400, 80, 40, 0, 0, 6.28318, false) // Full ellipse
	
	// Stroke the ellipse in orange
	dc.SetRGB(1, 0.5, 0)
	dc.SetLineWidth(2)
	dc.StrokePath2D(ellipsePath)
	
	// Save the image
	dc.SavePNG("images/paths/path2d-basic.png")
}
