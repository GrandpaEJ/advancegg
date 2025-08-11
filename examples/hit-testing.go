package main

import (
	"fmt"
	"image/color"
	"math/rand"

	"github.com/GrandpaEJ/advancegg"
)

func main() {
	fmt.Println("Creating Hit Testing demo...")
	
	// Create basic hit testing demo
	createBasicHitTestDemo()
	
	// Create interactive hit testing demo
	createInteractiveHitTestDemo()
	
	// Create complex shapes hit testing demo
	createComplexShapesDemo()
	
	fmt.Println("Hit Testing demo completed!")
}

func createBasicHitTestDemo() {
	// Create hit test manager
	htm := advancegg.NewHitTestManager()
	
	// Create hit testable shapes
	rect1 := advancegg.CreateHitTestRect("rect1", 50, 50, 100, 80)
	circle1 := advancegg.CreateHitTestCircle("circle1", 200, 100, 40)
	ellipse1 := advancegg.CreateHitTestEllipse("ellipse1", 350, 100, 60, 30)
	line1 := advancegg.CreateHitTestLine("line1", 50, 200, 350, 250, 10)
	
	// Add to hit test manager
	htm.AddObject(rect1)
	htm.AddObject(circle1)
	htm.AddObject(ellipse1)
	htm.AddObject(line1)
	
	// Create visualization
	dc := advancegg.NewContext(400, 300)
	dc.SetRGB(1, 1, 1)
	dc.Clear()
	
	// Draw title
	dc.SetRGB(0, 0, 0)
	dc.DrawString("Hit Testing Demo - Click Points", 100, 20)
	
	// Draw shapes
	dc.SetRGB(0.8, 0.8, 1.0)
	dc.DrawRectangle(rect1.X, rect1.Y, rect1.Width, rect1.Height)
	dc.Fill()
	
	dc.SetRGB(1.0, 0.8, 0.8)
	dc.DrawCircle(circle1.X, circle1.Y, circle1.Radius)
	dc.Fill()
	
	dc.SetRGB(0.8, 1.0, 0.8)
	dc.DrawEllipse(ellipse1.X, ellipse1.Y, ellipse1.RadiusX, ellipse1.RadiusY)
	dc.Fill()
	
	dc.SetRGB(1.0, 0.8, 1.0)
	dc.SetLineWidth(line1.Thickness)
	dc.MoveTo(line1.X1, line1.Y1)
	dc.LineTo(line1.X2, line1.Y2)
	dc.Stroke()
	
	// Test points and visualize hits
	testPoints := []struct {
		x, y float64
		name string
	}{
		{75, 75, "P1"},     // Inside rect
		{200, 100, "P2"},   // Center of circle
		{350, 100, "P3"},   // Center of ellipse
		{200, 225, "P4"},   // On line
		{300, 50, "P5"},    // Empty space
		{150, 150, "P6"},   // Empty space
	}
	
	for _, point := range testPoints {
		hits := htm.HitTest(point.x, point.y)
		
		// Draw test point
		if len(hits) > 0 {
			dc.SetRGB(1, 0, 0) // Red for hits
		} else {
			dc.SetRGB(0, 0, 1) // Blue for misses
		}
		dc.DrawCircle(point.x, point.y, 3)
		dc.Fill()
		
		// Draw point label
		dc.SetRGB(0, 0, 0)
		dc.DrawString(point.name, point.x+5, point.y-5)
		
		// Print hit results
		if len(hits) > 0 {
			fmt.Printf("Point %s (%.0f, %.0f) hits: ", point.name, point.x, point.y)
			for i, hit := range hits {
				if i > 0 {
					fmt.Print(", ")
				}
				switch h := hit.(type) {
				case *advancegg.HitTestRect:
					fmt.Print(h.ID)
				case *advancegg.HitTestCircle:
					fmt.Print(h.ID)
				case *advancegg.HitTestEllipse:
					fmt.Print(h.ID)
				case *advancegg.HitTestLine:
					fmt.Print(h.ID)
				}
			}
			fmt.Println()
		} else {
			fmt.Printf("Point %s (%.0f, %.0f) hits: none\n", point.name, point.x, point.y)
		}
	}
	
	// Add legend
	dc.SetRGB(0, 0, 0)
	dc.DrawString("Red dots = Hit, Blue dots = Miss", 50, 280)
	
	// Save result
	advancegg.SavePNG("hit-testing-basic.png", dc.Image())
	fmt.Println("Created hit-testing-basic.png")
}

func createInteractiveHitTestDemo() {
	// Create hit test manager
	htm := advancegg.NewHitTestManager()
	
	// Create overlapping shapes to test layering
	shapes := []advancegg.HitTestable{
		advancegg.CreateHitTestRect("bg-rect", 50, 50, 200, 150),
		advancegg.CreateHitTestCircle("circle1", 100, 100, 30),
		advancegg.CreateHitTestCircle("circle2", 150, 120, 25),
		advancegg.CreateHitTestRect("small-rect", 120, 80, 40, 30),
	}
	
	for _, shape := range shapes {
		htm.AddObject(shape)
	}
	
	// Create visualization
	dc := advancegg.NewContext(400, 300)
	dc.SetRGB(1, 1, 1)
	dc.Clear()
	
	// Draw title
	dc.SetRGB(0, 0, 0)
	dc.DrawString("Interactive Hit Testing - Overlapping Shapes", 50, 20)
	
	// Draw shapes with different colors
	colors := []color.RGBA{
		{200, 200, 255, 255}, // Light blue
		{255, 200, 200, 255}, // Light red
		{200, 255, 200, 255}, // Light green
		{255, 255, 200, 255}, // Light yellow
	}
	
	for i, shape := range shapes {
		dc.SetColor(colors[i])
		
		switch s := shape.(type) {
		case *advancegg.HitTestRect:
			dc.DrawRectangle(s.X, s.Y, s.Width, s.Height)
			dc.Fill()
		case *advancegg.HitTestCircle:
			dc.DrawCircle(s.X, s.Y, s.Radius)
			dc.Fill()
		}
	}
	
	// Test hit testing with different methods
	testX, testY := 125.0, 100.0
	
	// Test all hits
	allHits := htm.HitTest(testX, testY)
	fmt.Printf("\nAll hits at (%.0f, %.0f): %d objects\n", testX, testY, len(allHits))
	
	// Test first hit
	firstHit := htm.HitTestFirst(testX, testY)
	if firstHit != nil {
		fmt.Printf("First hit: %s\n", getShapeID(firstHit))
	}
	
	// Test last (topmost) hit
	lastHit := htm.HitTestLast(testX, testY)
	if lastHit != nil {
		fmt.Printf("Last (topmost) hit: %s\n", getShapeID(lastHit))
	}
	
	// Draw test point
	dc.SetRGB(1, 0, 0)
	dc.DrawCircle(testX, testY, 4)
	dc.Fill()
	dc.SetRGB(0, 0, 0)
	dc.DrawString("Test Point", testX+10, testY)
	
	// Save result
	advancegg.SavePNG("hit-testing-interactive.png", dc.Image())
	fmt.Println("Created hit-testing-interactive.png")
}

func createComplexShapesDemo() {
	// Create hit test manager
	htm := advancegg.NewHitTestManager()
	
	// Create polygon (triangle)
	trianglePoints := []advancegg.Point{
		{X: 100, Y: 50},
		{X: 50, Y: 150},
		{X: 150, Y: 150},
	}
	triangle := advancegg.CreateHitTestPolygon("triangle", trianglePoints)
	htm.AddObject(triangle)
	
	// Create star polygon
	starPoints := []advancegg.Point{
		{X: 250, Y: 50},  // Top
		{X: 260, Y: 80},  // Inner right
		{X: 290, Y: 80},  // Outer right
		{X: 270, Y: 100}, // Inner bottom right
		{X: 280, Y: 130}, // Bottom right
		{X: 250, Y: 110}, // Inner bottom
		{X: 220, Y: 130}, // Bottom left
		{X: 230, Y: 100}, // Inner bottom left
		{X: 210, Y: 80},  // Outer left
		{X: 240, Y: 80},  // Inner left
	}
	star := advancegg.CreateHitTestPolygon("star", starPoints)
	htm.AddObject(star)
	
	// Create path (zigzag)
	pathPoints := []advancegg.Point{
		{X: 50, Y: 200},
		{X: 100, Y: 180},
		{X: 150, Y: 220},
		{X: 200, Y: 180},
		{X: 250, Y: 220},
		{X: 300, Y: 180},
	}
	path := advancegg.CreateHitTestPath("zigzag", pathPoints, false)
	htm.AddObject(path)
	
	// Create visualization
	dc := advancegg.NewContext(400, 300)
	dc.SetRGB(1, 1, 1)
	dc.Clear()
	
	// Draw title
	dc.SetRGB(0, 0, 0)
	dc.DrawString("Complex Shapes Hit Testing", 100, 20)
	
	// Draw triangle
	dc.SetRGB(0.8, 0.8, 1.0)
	dc.MoveTo(trianglePoints[0].X, trianglePoints[0].Y)
	for i := 1; i < len(trianglePoints); i++ {
		dc.LineTo(trianglePoints[i].X, trianglePoints[i].Y)
	}
	dc.ClosePath()
	dc.Fill()
	
	// Draw star
	dc.SetRGB(1.0, 0.8, 0.8)
	dc.MoveTo(starPoints[0].X, starPoints[0].Y)
	for i := 1; i < len(starPoints); i++ {
		dc.LineTo(starPoints[i].X, starPoints[i].Y)
	}
	dc.ClosePath()
	dc.Fill()
	
	// Draw path
	dc.SetRGB(0.8, 1.0, 0.8)
	dc.SetLineWidth(8)
	dc.MoveTo(pathPoints[0].X, pathPoints[0].Y)
	for i := 1; i < len(pathPoints); i++ {
		dc.LineTo(pathPoints[i].X, pathPoints[i].Y)
	}
	dc.Stroke()
	
	// Test multiple points
	testPoints := []struct {
		x, y float64
		name string
	}{
		{100, 100, "T1"}, // Inside triangle
		{250, 90, "S1"},  // Inside star
		{150, 200, "P1"}, // On path
		{80, 80, "M1"},   // Miss
		{270, 120, "M2"}, // Miss
	}
	
	for _, point := range testPoints {
		hits := htm.HitTest(point.x, point.y)
		
		// Draw test point
		if len(hits) > 0 {
			dc.SetRGB(1, 0, 0) // Red for hits
		} else {
			dc.SetRGB(0, 0, 1) // Blue for misses
		}
		dc.DrawCircle(point.x, point.y, 3)
		dc.Fill()
		
		// Draw point label
		dc.SetRGB(0, 0, 0)
		dc.DrawString(point.name, point.x+5, point.y-5)
		
		// Print results
		if len(hits) > 0 {
			fmt.Printf("Point %s hits: %s\n", point.name, getShapeID(hits[0]))
		} else {
			fmt.Printf("Point %s hits: none\n", point.name)
		}
	}
	
	// Save result
	advancegg.SavePNG("hit-testing-complex.png", dc.Image())
	fmt.Println("Created hit-testing-complex.png")
}

// Helper function to get shape ID
func getShapeID(shape advancegg.HitTestable) string {
	switch s := shape.(type) {
	case *advancegg.HitTestRect:
		return s.ID
	case *advancegg.HitTestCircle:
		return s.ID
	case *advancegg.HitTestEllipse:
		return s.ID
	case *advancegg.HitTestPolygon:
		return s.ID
	case *advancegg.HitTestLine:
		return s.ID
	case *advancegg.HitTestPath:
		return s.ID
	default:
		return "unknown"
	}
}

// Demonstrate performance testing
func demonstratePerformance() {
	fmt.Println("\nPerformance Testing:")
	
	// Create many objects for performance testing
	htm := advancegg.NewHitTestManager()
	
	// Add many random shapes
	for i := 0; i < 1000; i++ {
		x := rand.Float64() * 800
		y := rand.Float64() * 600
		
		switch i % 3 {
		case 0:
			rect := advancegg.CreateHitTestRect(fmt.Sprintf("rect%d", i), x, y, 20, 20)
			htm.AddObject(rect)
		case 1:
			circle := advancegg.CreateHitTestCircle(fmt.Sprintf("circle%d", i), x, y, 10)
			htm.AddObject(circle)
		case 2:
			line := advancegg.CreateHitTestLine(fmt.Sprintf("line%d", i), x, y, x+20, y+20, 2)
			htm.AddObject(line)
		}
	}
	
	// Perform many hit tests
	hitCount := 0
	testCount := 10000
	
	for i := 0; i < testCount; i++ {
		x := rand.Float64() * 800
		y := rand.Float64() * 600
		hits := htm.HitTest(x, y)
		hitCount += len(hits)
	}
	
	fmt.Printf("Performed %d hit tests on 1000 objects\n", testCount)
	fmt.Printf("Total hits: %d\n", hitCount)
	fmt.Printf("Average hits per test: %.2f\n", float64(hitCount)/float64(testCount))
}
