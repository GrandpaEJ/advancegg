package main

import (
	"fmt"
	"image/color"
	"math"

	"github.com/GrandpaEJ/advancegg"
)

func main() {
	fmt.Println("Creating Advanced Stroke examples...")

	// Basic advanced strokes
	createBasicAdvancedStrokes()

	// Dashed patterns showcase
	createDashedPatterns()

	// Gradient strokes
	createGradientStrokes()

	// Tapered strokes
	createTaperedStrokes()

	// Creative stroke effects
	createCreativeStrokeEffects()

	fmt.Println("Advanced Stroke examples completed!")
}

func createBasicAdvancedStrokes() {
	dc := advancegg.NewContext(800, 600)

	// White background
	dc.SetRGB(1, 1, 1)
	dc.Clear()

	// Title
	dc.SetRGB(0, 0, 0)
	dc.DrawString("Advanced Stroke Styles", 50, 50)

	// Basic dashed line
	dc.SetRGB(0.8, 0.2, 0.2)
	dc.DrawDashedLine(50, 100, 750, 100, []float64{10, 5})
	dc.SetRGB(0, 0, 0)
	dc.DrawString("Basic dashed line (10, 5)", 50, 90)

	// Complex dash pattern
	dc.SetRGB(0.2, 0.8, 0.2)
	dc.DrawDashedLine(50, 150, 750, 150, []float64{20, 5, 5, 5})
	dc.SetRGB(0, 0, 0)
	dc.DrawString("Complex dash pattern (20, 5, 5, 5)", 50, 140)

	// Gradient line
	gradientStops := []advancegg.StrokeGradientStop{
		{Position: 0.0, Color: color.RGBA{255, 0, 0, 255}},
		{Position: 0.5, Color: color.RGBA{255, 255, 0, 255}},
		{Position: 1.0, Color: color.RGBA{0, 0, 255, 255}},
	}
	dc.DrawGradientLine(50, 200, 750, 200, gradientStops)
	dc.SetRGB(0, 0, 0)
	dc.DrawString("Gradient stroke (red → yellow → blue)", 50, 190)

	// Tapered line
	dc.SetRGB(0.8, 0.2, 0.8)
	dc.SetLineWidth(10)
	dc.DrawTaperedLine(50, 250, 750, 250, 1.0, 0.1)
	dc.SetRGB(0, 0, 0)
	dc.DrawString("Tapered stroke (thick → thin)", 50, 240)

	// Reverse tapered line
	dc.SetRGB(0.2, 0.8, 0.8)
	dc.SetLineWidth(10)
	dc.DrawTaperedLine(50, 300, 750, 300, 0.1, 1.0)
	dc.SetRGB(0, 0, 0)
	dc.DrawString("Reverse tapered stroke (thin → thick)", 50, 290)

	// Multiple line caps demonstration
	dc.SetRGB(0, 0, 0)
	dc.DrawString("Line Cap Styles:", 50, 350)

	// Butt cap
	dc.SetRGB(0.6, 0.6, 0.6)
	dc.SetLineWidth(20)
	dc.MoveTo(100, 380)
	dc.LineTo(200, 380)
	dc.Stroke()
	dc.SetRGB(0, 0, 0)
	dc.DrawString("Butt", 120, 400)

	// Round cap
	dc.SetRGB(0.6, 0.6, 0.6)
	dc.MoveTo(300, 380)
	dc.LineTo(400, 380)
	dc.Stroke()
	dc.SetRGB(0, 0, 0)
	dc.DrawString("Round", 320, 400)

	// Square cap
	dc.SetRGB(0.6, 0.6, 0.6)
	dc.MoveTo(500, 380)
	dc.LineTo(600, 380)
	dc.Stroke()
	dc.SetRGB(0, 0, 0)
	dc.DrawString("Square", 520, 400)

	// Line join demonstration
	dc.SetRGB(0, 0, 0)
	dc.DrawString("Line Join Styles:", 50, 450)

	// Miter join
	dc.SetRGB(0.4, 0.4, 0.8)
	dc.SetLineWidth(15)
	dc.MoveTo(100, 480)
	dc.LineTo(150, 500)
	dc.LineTo(200, 480)
	dc.Stroke()
	dc.SetRGB(0, 0, 0)
	dc.DrawString("Miter", 120, 520)

	// Round join
	dc.SetRGB(0.4, 0.4, 0.8)
	dc.MoveTo(300, 480)
	dc.LineTo(350, 500)
	dc.LineTo(400, 480)
	dc.Stroke()
	dc.SetRGB(0, 0, 0)
	dc.DrawString("Round", 320, 520)

	// Bevel join
	dc.SetRGB(0.4, 0.4, 0.8)
	dc.MoveTo(500, 480)
	dc.LineTo(550, 500)
	dc.LineTo(600, 480)
	dc.Stroke()
	dc.SetRGB(0, 0, 0)
	dc.DrawString("Bevel", 520, 520)

	dc.SavePNG("images/basic-advanced-strokes.png")
	fmt.Println("Basic advanced strokes demo saved as basic-advanced-strokes.png")
}

func createDashedPatterns() {
	dc := advancegg.NewContext(1000, 700)

	// Light background
	dc.SetRGB(0.95, 0.95, 0.98)
	dc.Clear()

	// Title
	dc.SetRGB(0, 0, 0)
	dc.DrawString("Dashed Pattern Showcase", 50, 50)

	// Various dash patterns
	patterns := []struct {
		name    string
		pattern []float64
		color   [3]float64
	}{
		{"Solid", []float64{}, [3]float64{0, 0, 0}},
		{"Basic Dash", []float64{10, 10}, [3]float64{0.8, 0.2, 0.2}},
		{"Long Dash", []float64{20, 10}, [3]float64{0.2, 0.8, 0.2}},
		{"Dot", []float64{2, 8}, [3]float64{0.2, 0.2, 0.8}},
		{"Dash-Dot", []float64{15, 5, 3, 5}, [3]float64{0.8, 0.5, 0.2}},
		{"Dash-Dot-Dot", []float64{15, 5, 3, 5, 3, 5}, [3]float64{0.8, 0.2, 0.8}},
		{"Railroad", []float64{10, 10, 10, 10, 50, 10}, [3]float64{0.5, 0.5, 0.5}},
		{"Morse Code", []float64{3, 3, 3, 3, 15, 3}, [3]float64{0.2, 0.8, 0.8}},
	}

	y := 100.0
	for _, pattern := range patterns {
		dc.SetRGB(pattern.color[0], pattern.color[1], pattern.color[2])
		dc.SetLineWidth(4)

		if len(pattern.pattern) == 0 {
			// Solid line
			dc.MoveTo(200, y)
			dc.LineTo(900, y)
			dc.Stroke()
		} else {
			dc.DrawDashedLine(200, y, 900, y, pattern.pattern)
		}

		// Label
		dc.SetRGB(0, 0, 0)
		dc.DrawString(pattern.name, 50, y+5)

		y += 60
	}

	// Curved dashed lines
	dc.SetRGB(0, 0, 0)
	dc.DrawString("Dashed Curves:", 50, 550)

	// Sine wave with dashes
	dc.SetRGB(0.8, 0.2, 0.4)
	dc.SetLineWidth(3)
	points := make([][2]float64, 0)
	for x := 100; x <= 900; x += 2 {
		y := 600 + 30*math.Sin(float64(x-100)*0.02)
		points = append(points, [2]float64{float64(x), y})
	}

	// Draw dashed curve by sampling points
	dashPattern := []float64{8, 4}
	drawDashedCurve(dc, points, dashPattern)

	dc.SavePNG("images/dashed-patterns.png")
	fmt.Println("Dashed patterns demo saved as dashed-patterns.png")
}

func drawDashedCurve(dc *advancegg.Context, points [][2]float64, pattern []float64) {
	if len(points) < 2 || len(pattern) == 0 {
		return
	}

	currentDistance := 0.0
	patternIndex := 0
	drawing := true

	for i := 1; i < len(points); i++ {
		x1, y1 := points[i-1][0], points[i-1][1]
		x2, y2 := points[i][0], points[i][1]

		segmentLength := math.Sqrt((x2-x1)*(x2-x1) + (y2-y1)*(y2-y1))

		if drawing {
			dc.MoveTo(x1, y1)
			dc.LineTo(x2, y2)
			dc.Stroke()
		}

		currentDistance += segmentLength
		patternLength := pattern[patternIndex%len(pattern)]

		if currentDistance >= patternLength {
			currentDistance -= patternLength
			patternIndex++
			drawing = !drawing
		}
	}
}

func createGradientStrokes() {
	dc := advancegg.NewContext(800, 600)

	// Dark background for better gradient visibility
	dc.SetRGB(0.1, 0.1, 0.2)
	dc.Clear()

	// Title
	dc.SetRGB(1, 1, 1)
	dc.DrawString("Gradient Stroke Effects", 50, 50)

	// Rainbow gradient
	rainbowStops := []advancegg.StrokeGradientStop{
		{Position: 0.0, Color: color.RGBA{255, 0, 0, 255}},     // Red
		{Position: 0.17, Color: color.RGBA{255, 165, 0, 255}},  // Orange
		{Position: 0.33, Color: color.RGBA{255, 255, 0, 255}},  // Yellow
		{Position: 0.5, Color: color.RGBA{0, 255, 0, 255}},     // Green
		{Position: 0.67, Color: color.RGBA{0, 0, 255, 255}},    // Blue
		{Position: 0.83, Color: color.RGBA{75, 0, 130, 255}},   // Indigo
		{Position: 1.0, Color: color.RGBA{238, 130, 238, 255}}, // Violet
	}

	dc.SetLineWidth(8)
	dc.DrawGradientLine(50, 100, 750, 100, rainbowStops)
	dc.SetRGB(1, 1, 1)
	dc.DrawString("Rainbow gradient", 50, 90)

	// Fire gradient
	fireStops := []advancegg.StrokeGradientStop{
		{Position: 0.0, Color: color.RGBA{255, 255, 0, 255}}, // Yellow
		{Position: 0.5, Color: color.RGBA{255, 100, 0, 255}}, // Orange
		{Position: 1.0, Color: color.RGBA{200, 0, 0, 255}},   // Dark red
	}

	dc.DrawGradientLine(50, 150, 750, 150, fireStops)
	dc.SetRGB(1, 1, 1)
	dc.DrawString("Fire gradient", 50, 140)

	// Ocean gradient
	oceanStops := []advancegg.StrokeGradientStop{
		{Position: 0.0, Color: color.RGBA{0, 255, 255, 255}}, // Cyan
		{Position: 0.5, Color: color.RGBA{0, 100, 255, 255}}, // Blue
		{Position: 1.0, Color: color.RGBA{0, 0, 100, 255}},   // Dark blue
	}

	dc.DrawGradientLine(50, 200, 750, 200, oceanStops)
	dc.SetRGB(1, 1, 1)
	dc.DrawString("Ocean gradient", 50, 190)

	// Metallic gradient
	metallicStops := []advancegg.StrokeGradientStop{
		{Position: 0.0, Color: color.RGBA{200, 200, 200, 255}}, // Light gray
		{Position: 0.3, Color: color.RGBA{255, 255, 255, 255}}, // White
		{Position: 0.7, Color: color.RGBA{150, 150, 150, 255}}, // Medium gray
		{Position: 1.0, Color: color.RGBA{100, 100, 100, 255}}, // Dark gray
	}

	dc.DrawGradientLine(50, 250, 750, 250, metallicStops)
	dc.SetRGB(1, 1, 1)
	dc.DrawString("Metallic gradient", 50, 240)

	// Gradient shapes
	dc.SetRGB(1, 1, 1)
	dc.DrawString("Gradient Shapes:", 50, 320)

	// Gradient circle
	dc.SetLineWidth(12)
	for angle := 0.0; angle < 2*math.Pi; angle += 0.1 {
		x1 := 200 + 80*math.Cos(angle)
		y1 := 400 + 80*math.Sin(angle)
		x2 := 200 + 80*math.Cos(angle+0.1)
		y2 := 400 + 80*math.Sin(angle+0.1)

		// Create gradient based on angle
		t := angle / (2 * math.Pi)
		gradientStops := []advancegg.StrokeGradientStop{
			{Position: 0.0, Color: color.RGBA{uint8(255 * t), uint8(255 * (1 - t)), 128, 255}},
			{Position: 1.0, Color: color.RGBA{uint8(255 * (1 - t)), uint8(255 * t), 128, 255}},
		}

		dc.DrawGradientLine(x1, y1, x2, y2, gradientStops)
	}

	// Gradient spiral
	dc.SetLineWidth(6)
	for t := 0.0; t < 4*math.Pi; t += 0.1 {
		radius := 20 + t*5
		x1 := 500 + radius*math.Cos(t)
		y1 := 400 + radius*math.Sin(t)
		x2 := 500 + (radius+2)*math.Cos(t+0.1)
		y2 := 400 + (radius+2)*math.Sin(t+0.1)

		// Create gradient based on spiral position
		spiralT := t / (4 * math.Pi)
		gradientStops := []advancegg.StrokeGradientStop{
			{Position: 0.0, Color: color.RGBA{255, uint8(255 * spiralT), uint8(255 * (1 - spiralT)), 255}},
			{Position: 1.0, Color: color.RGBA{uint8(255 * (1 - spiralT)), 255, uint8(255 * spiralT), 255}},
		}

		dc.DrawGradientLine(x1, y1, x2, y2, gradientStops)
	}

	dc.SavePNG("images/gradient-strokes.png")
	fmt.Println("Gradient strokes demo saved as gradient-strokes.png")
}

func createTaperedStrokes() {
	dc := advancegg.NewContext(800, 600)

	// White background
	dc.SetRGB(1, 1, 1)
	dc.Clear()

	// Title
	dc.SetRGB(0, 0, 0)
	dc.DrawString("Tapered Stroke Effects", 50, 50)

	// Linear taper examples
	dc.SetRGB(0, 0, 0)
	dc.DrawString("Linear Tapers:", 50, 90)

	// Thick to thin
	dc.SetRGB(0.8, 0.2, 0.2)
	dc.SetLineWidth(20)
	dc.DrawTaperedLine(100, 120, 700, 120, 1.0, 0.1)
	dc.SetRGB(0, 0, 0)
	dc.DrawString("Thick → Thin", 50, 130)

	// Thin to thick
	dc.SetRGB(0.2, 0.8, 0.2)
	dc.SetLineWidth(20)
	dc.DrawTaperedLine(100, 170, 700, 170, 0.1, 1.0)
	dc.SetRGB(0, 0, 0)
	dc.DrawString("Thin → Thick", 50, 180)

	// Both ends thin
	dc.SetRGB(0.2, 0.2, 0.8)
	dc.SetLineWidth(20)
	dc.DrawTaperedLine(100, 220, 700, 220, 0.2, 0.2)
	dc.SetRGB(0, 0, 0)
	dc.DrawString("Thin → Thick → Thin (simulated)", 50, 230)

	// Calligraphy-style strokes
	dc.SetRGB(0, 0, 0)
	dc.DrawString("Calligraphy Style:", 50, 280)

	// Brush stroke simulation
	brushStrokes := []struct {
		x1, y1, x2, y2 float64
		start, end     float64
		color          [3]float64
	}{
		{100, 320, 300, 300, 0.8, 0.2, [3]float64{0.2, 0.2, 0.2}},
		{320, 300, 500, 330, 0.2, 0.9, [3]float64{0.2, 0.2, 0.2}},
		{520, 330, 700, 310, 0.9, 0.1, [3]float64{0.2, 0.2, 0.2}},
	}

	dc.SetLineWidth(25)
	for _, stroke := range brushStrokes {
		dc.SetRGB(stroke.color[0], stroke.color[1], stroke.color[2])
		dc.DrawTaperedLine(stroke.x1, stroke.y1, stroke.x2, stroke.y2, stroke.start, stroke.end)
	}

	// Artistic flourishes
	dc.SetRGB(0, 0, 0)
	dc.DrawString("Artistic Flourishes:", 50, 380)

	// Curved tapered lines
	dc.SetRGB(0.6, 0.2, 0.8)
	dc.SetLineWidth(15)

	// Create curved paths with tapers
	curves := []struct {
		points [][2]float64
		taper  [2]float64
	}{
		{
			points: [][2]float64{{100, 420}, {200, 400}, {300, 440}, {400, 410}},
			taper:  [2]float64{1.0, 0.2},
		},
		{
			points: [][2]float64{{450, 410}, {550, 450}, {650, 400}, {750, 430}},
			taper:  [2]float64{0.3, 1.0},
		},
	}

	for _, curve := range curves {
		for i := 1; i < len(curve.points); i++ {
			t := float64(i-1) / float64(len(curve.points)-2)
			taperStart := curve.taper[0] + t*(curve.taper[1]-curve.taper[0])
			taperEnd := curve.taper[0] + (t+0.1)*(curve.taper[1]-curve.taper[0])

			dc.DrawTaperedLine(
				curve.points[i-1][0], curve.points[i-1][1],
				curve.points[i][0], curve.points[i][1],
				taperStart, taperEnd,
			)
		}
	}

	// Feather effect
	dc.SetRGB(0, 0, 0)
	dc.DrawString("Feather Effect:", 50, 480)

	dc.SetRGB(0.4, 0.6, 0.8)
	dc.SetLineWidth(12)

	// Main feather shaft
	dc.DrawTaperedLine(400, 520, 400, 580, 1.0, 0.3)

	// Feather barbs
	for i := 0; i < 10; i++ {
		y := 525 + float64(i)*5
		leftX := 400 - float64(i)*3
		rightX := 400 + float64(i)*3

		dc.SetLineWidth(3)
		dc.DrawTaperedLine(400, y, leftX, y-5, 0.8, 0.1)
		dc.DrawTaperedLine(400, y, rightX, y-5, 0.8, 0.1)
	}

	dc.SavePNG("images/tapered-strokes.png")
	fmt.Println("Tapered strokes demo saved as tapered-strokes.png")
}

func createCreativeStrokeEffects() {
	dc := advancegg.NewContext(1200, 800)

	// Gradient background
	for y := 0; y < 800; y++ {
		t := float64(y) / 800.0
		dc.SetRGB(0.1+t*0.1, 0.1+t*0.1, 0.2+t*0.1)
		dc.DrawLine(0, float64(y), 1200, float64(y))
		dc.Stroke()
	}

	// Title
	dc.SetRGB(1, 1, 1)
	dc.DrawString("Creative Advanced Stroke Effects", 50, 50)

	// Lightning effect with tapered dashed lines
	dc.SetRGB(1, 1, 0.8)
	dc.DrawString("Lightning Effect:", 50, 100)

	lightningPath := [][2]float64{
		{100, 150}, {120, 200}, {110, 250}, {130, 300}, {115, 350},
	}

	dc.SetLineWidth(8)
	for i := 1; i < len(lightningPath); i++ {
		// Main bolt
		dc.SetRGB(1, 1, 0.8)
		dc.DrawTaperedLine(
			lightningPath[i-1][0], lightningPath[i-1][1],
			lightningPath[i][0], lightningPath[i][1],
			1.0, 0.6,
		)

		// Glow effect
		dc.SetRGB(0.8, 0.8, 1)
		dc.SetLineWidth(15)
		dc.DrawTaperedLine(
			lightningPath[i-1][0], lightningPath[i-1][1],
			lightningPath[i][0], lightningPath[i][1],
			0.3, 0.1,
		)
		dc.SetLineWidth(8)
	}

	// Neon sign effect
	dc.SetRGB(1, 1, 1)
	dc.DrawString("Neon Sign Effect:", 300, 100)

	// Neon "OPEN" sign
	dc.SetLineWidth(12)

	// O
	neonColor := []advancegg.StrokeGradientStop{
		{Position: 0.0, Color: color.RGBA{255, 0, 100, 255}},
		{Position: 0.5, Color: color.RGBA{255, 100, 200, 255}},
		{Position: 1.0, Color: color.RGBA{255, 0, 100, 255}},
	}

	for angle := 0.0; angle < 2*math.Pi; angle += 0.2 {
		x1 := 350 + 30*math.Cos(angle)
		y1 := 200 + 30*math.Sin(angle)
		x2 := 350 + 30*math.Cos(angle+0.2)
		y2 := 200 + 30*math.Sin(angle+0.2)
		dc.DrawGradientLine(x1, y1, x2, y2, neonColor)
	}

	// Brush painting effect
	dc.SetRGB(1, 1, 1)
	dc.DrawString("Brush Painting:", 600, 100)

	// Paint brush strokes
	paintStrokes := []struct {
		path  [][2]float64
		color [3]float64
		width float64
	}{
		{
			path:  [][2]float64{{650, 150}, {700, 180}, {750, 160}, {800, 190}},
			color: [3]float64{0.8, 0.2, 0.2},
			width: 20,
		},
		{
			path:  [][2]float64{{660, 200}, {710, 220}, {760, 200}, {810, 230}},
			color: [3]float64{0.2, 0.8, 0.2},
			width: 15,
		},
		{
			path:  [][2]float64{{670, 250}, {720, 270}, {770, 250}, {820, 280}},
			color: [3]float64{0.2, 0.2, 0.8},
			width: 18,
		},
	}

	for _, stroke := range paintStrokes {
		dc.SetRGB(stroke.color[0], stroke.color[1], stroke.color[2])
		dc.SetLineWidth(stroke.width)

		for i := 1; i < len(stroke.path); i++ {
			// Vary taper for organic brush effect
			taper1 := 0.7 + 0.3*math.Sin(float64(i)*0.5)
			taper2 := 0.7 + 0.3*math.Sin(float64(i+1)*0.5)

			dc.DrawTaperedLine(
				stroke.path[i-1][0], stroke.path[i-1][1],
				stroke.path[i][0], stroke.path[i][1],
				taper1, taper2,
			)
		}
	}

	// Circuit board effect
	dc.SetRGB(1, 1, 1)
	dc.DrawString("Circuit Board Effect:", 50, 400)

	// Circuit traces
	dc.SetRGB(0, 1, 0.5)
	dc.SetLineWidth(3)

	circuitPaths := []struct {
		start, end [2]float64
		pattern    []float64
	}{
		{[2]float64{100, 450}, [2]float64{300, 450}, []float64{5, 2}},
		{[2]float64{300, 450}, [2]float64{300, 550}, []float64{3, 3}},
		{[2]float64{300, 550}, [2]float64{500, 550}, []float64{8, 4}},
		{[2]float64{150, 500}, [2]float64{450, 500}, []float64{2, 2}},
	}

	for _, path := range circuitPaths {
		dc.DrawDashedLine(path.start[0], path.start[1], path.end[0], path.end[1], path.pattern)
	}

	// Circuit nodes
	nodes := [][2]float64{
		{100, 450}, {300, 450}, {300, 550}, {500, 550}, {150, 500}, {450, 500},
	}

	for _, node := range nodes {
		dc.SetRGB(0, 1, 1)
		dc.DrawCircle(node[0], node[1], 8)
		dc.Fill()
	}

	// Calligraphy showcase
	dc.SetRGB(1, 1, 1)
	dc.DrawString("Calligraphy Showcase:", 600, 400)

	// Elegant script letters
	dc.SetRGB(0.2, 0.2, 0.2)
	dc.SetLineWidth(25)

	// Letter "A"
	dc.DrawTaperedLine(650, 550, 680, 450, 0.2, 1.0)
	dc.DrawTaperedLine(680, 450, 710, 550, 1.0, 0.2)
	dc.DrawTaperedLine(665, 500, 695, 500, 0.6, 0.6)

	// Letter "G"
	for angle := math.Pi; angle < 2*math.Pi; angle += 0.1 {
		x1 := 780 + 30*math.Cos(angle)
		y1 := 500 + 30*math.Sin(angle)
		x2 := 780 + 30*math.Cos(angle+0.1)
		y2 := 500 + 30*math.Sin(angle+0.1)

		taper := 0.5 + 0.5*math.Sin(angle*2)
		dc.DrawTaperedLine(x1, y1, x2, y2, taper, taper)
	}

	// Information
	dc.SetRGB(0.7, 0.7, 0.7)
	dc.DrawString("Advanced stroke styles enable professional graphics and artistic effects", 50, 750)

	dc.SavePNG("images/creative-stroke-effects.png")
	fmt.Println("Creative stroke effects demo saved as creative-stroke-effects.png")
}
