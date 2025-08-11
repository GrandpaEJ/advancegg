package main

import (
	"fmt"
	"math"

	"github.com/GrandpaEJ/advancegg"
	"github.com/GrandpaEJ/advancegg/internal/advance"
	"github.com/GrandpaEJ/advancegg/internal/core"
)

func main() {
	fmt.Println("Testing text-on-path rendering...")

	// Create context
	dc := advancegg.NewContext(1200, 800)
	dc.SetRGB(1, 1, 1)
	dc.Clear()

	// Title
	dc.SetRGB(0, 0, 0)
	dc.DrawString("Text-on-Path Rendering Test", 50, 50)

	// Test 1: Text on circular path
	fmt.Println("Testing text on circular path...")
	testCircularPath(dc, 300, 150, 80, "Circular Path Text")

	// Test 2: Text on sine wave
	fmt.Println("Testing text on sine wave...")
	testSineWavePath(dc, 100, 300, 400, 50, "Sine Wave Text")

	// Test 3: Text on Bezier curve
	fmt.Println("Testing text on Bezier curve...")
	testBezierPath(dc, 700, 150, "Bezier Curve Text")

	// Test 4: Text on line
	fmt.Println("Testing text on line...")
	testLineText(dc, 100, 500)

	// Test 5: Text on arc
	fmt.Println("Testing text on arc...")
	testArcText(dc, 700, 500)

	// Save the result
	err := dc.SavePNG("text-on-path-test.png")
	if err != nil {
		fmt.Printf("Error saving image: %v\n", err)
		return
	}

	fmt.Println("Text-on-path test completed!")
	fmt.Println("Generated: text-on-path-test.png")
}

func testCircularPath(dc *advancegg.Context, centerX, centerY, radius float64, text string) {
	// Create circular path
	path := core.NewPath2D()
	path.Arc(centerX, centerY, radius, 0, 2*math.Pi, false)

	// Draw the path for reference (using existing circle drawing)
	dc.SetRGB(0.8, 0.8, 0.8)
	dc.SetLineWidth(1)
	dc.DrawCircle(centerX, centerY, radius)
	dc.Stroke()

	// Draw text on path using existing function
	dc.SetRGB(0, 0, 0.8)
	advance.DrawTextOnCircle(dc, text, centerX, centerY, radius)

	// Add label
	dc.SetRGB(0.5, 0.5, 0.5)
	dc.DrawString("Circular Path", centerX-40, centerY+radius+20)
}

func testSineWavePath(dc *advancegg.Context, startX, startY, width, amplitude float64, text string) {
	// Draw sine wave using existing function
	dc.SetRGB(0.8, 0.8, 0.8)
	dc.SetLineWidth(1)

	// Draw reference sine wave manually
	dc.MoveTo(startX, startY)
	steps := 100
	for i := 1; i <= steps; i++ {
		x := startX + (width * float64(i) / float64(steps))
		y := startY + amplitude*math.Sin(2*math.Pi*float64(i)/float64(steps)*2)
		dc.LineTo(x, y)
	}
	dc.Stroke()

	// Draw text on wave using existing function
	dc.SetRGB(0.8, 0, 0)
	advance.DrawTextOnWave(dc, text, startX, startY, startX+width, amplitude, 2)

	// Add label
	dc.SetRGB(0.5, 0.5, 0.5)
	dc.DrawString("Sine Wave", startX, startY+amplitude+30)
}

func testBezierPath(dc *advancegg.Context, startX, startY float64, text string) {
	// Draw reference curve manually
	dc.SetRGB(0.8, 0.8, 0.8)
	dc.SetLineWidth(1)
	dc.MoveTo(startX, startY)
	// Approximate Bezier with line segments
	steps := 20
	for i := 1; i <= steps; i++ {
		t := float64(i) / float64(steps)
		x := startX + t*200 + 50*math.Sin(t*math.Pi)
		y := startY + 30*math.Sin(t*math.Pi*2)
		dc.LineTo(x, y)
	}
	dc.Stroke()

	// Draw text on Bezier using existing function
	dc.SetRGB(0, 0.8, 0)
	advance.DrawTextOnBezier(dc, text, startX, startY, startX+100, startY-25, startX+200, startY)

	// Add label
	dc.SetRGB(0.5, 0.5, 0.5)
	dc.DrawString("Bezier Curve", startX+50, startY+40)
}

func testLineText(dc *advancegg.Context, startX, startY float64) {
	// Draw reference line
	dc.SetRGB(0.8, 0.8, 0.8)
	dc.SetLineWidth(1)
	dc.MoveTo(startX, startY)
	dc.LineTo(startX+300, startY)
	dc.Stroke()

	// Draw text using spiral with minimal turns (approximates line)
	dc.SetRGB(0, 0, 0.8)
	advance.DrawTextOnSpiral(dc, "Spiral Text", startX+150, startY, 150, 150, 1)

	// Add label
	dc.SetRGB(0.5, 0.5, 0.5)
	dc.DrawString("Spiral Text", startX, startY+30)
}

func testArcText(dc *advancegg.Context, centerX, centerY float64) {
	radius := 80.0

	// Draw reference arc
	dc.SetRGB(0.8, 0.8, 0.8)
	dc.SetLineWidth(1)
	dc.DrawArc(centerX, centerY, radius, 0, math.Pi)
	dc.Stroke()

	// Draw text on arc
	dc.SetRGB(0.8, 0, 0.8)
	advance.DrawTextOnArc(dc, "Text on Arc", centerX, centerY, radius, 0, math.Pi)

	// Add label
	dc.SetRGB(0.5, 0.5, 0.5)
	dc.DrawString("Arc Text", centerX-30, centerY+radius+20)
}
