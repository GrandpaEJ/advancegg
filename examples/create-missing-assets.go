package main

import (
	"math"
	"math/rand"
	"time"

	"github.com/GrandpaEJ/advancegg"
)

func createBaboonImage() {
	// Create a colorful test image that resembles the classic baboon test image
	dc := advancegg.NewContext(512, 512)

	// Create a gradient background
	for y := 0; y < 512; y++ {
		for x := 0; x < 512; x++ {
			// Create a complex pattern with multiple colors
			r := float64(x) / 512.0
			g := float64(y) / 512.0
			b := math.Sin(float64(x+y)/50.0)*0.5 + 0.5

			dc.SetRGB(r*0.8+0.2, g*0.6+0.3, b*0.7+0.2)
			dc.DrawRectangle(float64(x), float64(y), 1, 1)
			dc.Fill()
		}
	}

	// Add some noise and texture
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 5000; i++ {
		x := rand.Float64() * 512
		y := rand.Float64() * 512
		size := rand.Float64()*3 + 1

		dc.SetRGBA(rand.Float64(), rand.Float64(), rand.Float64(), 0.3)
		dc.DrawCircle(x, y, size)
		dc.Fill()
	}

	// Add some geometric patterns
	dc.SetRGBA(0.8, 0.2, 0.2, 0.6)
	for i := 0; i < 20; i++ {
		angle := float64(i) * 2 * math.Pi / 20
		x := 256 + 100*math.Cos(angle)
		y := 256 + 100*math.Sin(angle)
		dc.DrawCircle(x, y, 15)
		dc.Fill()
	}

	dc.SavePNG("examples/baboon.png")
}

func createGopherImage() {
	// Create a simple gopher-like mascot image
	dc := advancegg.NewContext(400, 400)

	// Background
	dc.SetRGB(0.9, 0.95, 1.0)
	dc.Clear()

	// Gopher body (simplified)
	dc.SetRGB(0.4, 0.7, 0.9) // Light blue
	dc.DrawEllipse(200, 250, 80, 100)
	dc.Fill()

	// Gopher head
	dc.SetRGB(0.4, 0.7, 0.9)
	dc.DrawCircle(200, 150, 60)
	dc.Fill()

	// Eyes
	dc.SetRGB(0, 0, 0)
	dc.DrawCircle(180, 140, 8)
	dc.Fill()
	dc.DrawCircle(220, 140, 8)
	dc.Fill()

	// Eye highlights
	dc.SetRGB(1, 1, 1)
	dc.DrawCircle(182, 138, 3)
	dc.Fill()
	dc.DrawCircle(222, 138, 3)
	dc.Fill()

	// Nose
	dc.SetRGB(0.2, 0.2, 0.2)
	dc.DrawEllipse(200, 160, 4, 6)
	dc.Fill()

	// Mouth
	dc.SetLineWidth(3)
	dc.SetRGB(0.2, 0.2, 0.2)
	dc.MoveTo(190, 175)
	dc.QuadraticTo(200, 185, 210, 175)
	dc.Stroke()

	// Ears
	dc.SetRGB(0.4, 0.7, 0.9)
	dc.DrawEllipse(170, 120, 15, 25)
	dc.Fill()
	dc.DrawEllipse(230, 120, 15, 25)
	dc.Fill()

	// Arms
	dc.SetRGB(0.4, 0.7, 0.9)
	dc.DrawEllipse(140, 220, 20, 40)
	dc.Fill()
	dc.DrawEllipse(260, 220, 20, 40)
	dc.Fill()

	// Legs
	dc.DrawEllipse(170, 320, 25, 50)
	dc.Fill()
	dc.DrawEllipse(230, 320, 25, 50)
	dc.Fill()

	// Add "Go" text
	dc.SetRGB(0.2, 0.4, 0.6)
	dc.DrawStringAnchored("Go", 200, 380, 0.5, 0.5)

	dc.SavePNG("examples/gopher.png")
}

func main() {
	println("Creating missing asset files...")

	createBaboonImage()
	println("Created examples/baboon.png")

	createGopherImage()
	println("Created examples/gopher.png")

	println("All missing assets created successfully!")
}
