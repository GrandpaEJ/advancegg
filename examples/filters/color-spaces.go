package main

import (
	"fmt"
	"math"

	"github.com/GrandpaEJ/advancegg"
)

func main() {
	dc := advancegg.NewContext(1200, 800)

	// White background
	dc.SetRGB(1, 1, 1)
	dc.Clear()

	// Title
	dc.SetRGB(0, 0, 0)
	dc.DrawString("Color Space Demonstrations", 50, 50)

	// Demonstrate different color spaces
	demonstrateRGBtoCMYK(dc)
	demonstrateHSVColorWheel(dc)
	demonstrateHSLGradients(dc)
	demonstrateLABColorSpace(dc)

	// Save the demonstration
	dc.SavePNG("images/filters/color-spaces.png")
	fmt.Println("Color spaces demo saved as color-spaces.png")

	// Create individual color space examples
	createCMYKExample()
	createHSVExample()
	createLABExample()

	fmt.Println("All color space examples completed!")
}

func demonstrateRGBtoCMYK(dc *advancegg.Context) {
	// RGB to CMYK conversion demonstration
	dc.SetRGB(0, 0, 0)
	dc.DrawString("RGB to CMYK Conversion:", 50, 100)

	colors := []struct {
		name    string
		r, g, b float64
		x, y    float64
	}{
		{"Red", 1, 0, 0, 50, 130},
		{"Green", 0, 1, 0, 200, 130},
		{"Blue", 0, 0, 1, 350, 130},
		{"Yellow", 1, 1, 0, 500, 130},
		{"Magenta", 1, 0, 1, 650, 130},
		{"Cyan", 0, 1, 1, 800, 130},
	}

	for _, color := range colors {
		// Draw RGB color
		dc.SetRGB(color.r, color.g, color.b)
		dc.DrawCircle(color.x, color.y, 20)
		dc.Fill()

		// Convert to CMYK and display values
		rgbColor := advancegg.NewColor(color.r, color.g, color.b, 1.0)
		cmyk := rgbColor.ToCMYK()

		dc.SetRGB(0, 0, 0)
		dc.DrawString(color.name, color.x-15, color.y+35)
		dc.DrawString(fmt.Sprintf("C:%.2f M:%.2f", cmyk.C, cmyk.M), color.x-25, color.y+50)
		dc.DrawString(fmt.Sprintf("Y:%.2f K:%.2f", cmyk.Y, cmyk.K), color.x-25, color.y+65)
	}
}

func demonstrateHSVColorWheel(dc *advancegg.Context) {
	// HSV color wheel
	dc.SetRGB(0, 0, 0)
	dc.DrawString("HSV Color Wheel:", 50, 220)

	centerX, centerY := 150.0, 300.0
	radius := 60.0

	// Draw HSV color wheel
	for angle := 0.0; angle < 360; angle += 2 {
		for r := 0.0; r < radius; r += 2 {
			h := angle
			s := r / radius
			v := 1.0

			dc.SetHSV(h, s, v)

			x := centerX + r*math.Cos(angle*math.Pi/180)
			y := centerY + r*math.Sin(angle*math.Pi/180)

			dc.SetPixel(int(x), int(y))
		}
	}

	// HSV sliders demonstration
	dc.SetRGB(0, 0, 0)
	dc.DrawString("HSV Sliders:", 300, 250)

	// Hue slider
	for x := 0; x < 200; x++ {
		h := float64(x) * 360.0 / 200.0
		dc.SetHSV(h, 1.0, 1.0)
		dc.DrawLine(300+float64(x), 270, 300+float64(x), 285)
		dc.Stroke()
	}

	// Saturation slider
	for x := 0; x < 200; x++ {
		s := float64(x) / 200.0
		dc.SetHSV(240, s, 1.0) // Blue hue
		dc.DrawLine(300+float64(x), 300, 300+float64(x), 315)
		dc.Stroke()
	}

	// Value slider
	for x := 0; x < 200; x++ {
		v := float64(x) / 200.0
		dc.SetHSV(240, 1.0, v) // Blue hue, full saturation
		dc.DrawLine(300+float64(x), 330, 300+float64(x), 345)
		dc.Stroke()
	}

	dc.SetRGB(0, 0, 0)
	dc.DrawString("H", 280, 280)
	dc.DrawString("S", 280, 310)
	dc.DrawString("V", 280, 340)
}

func demonstrateHSLGradients(dc *advancegg.Context) {
	// HSL gradients
	dc.SetRGB(0, 0, 0)
	dc.DrawString("HSL Gradients:", 550, 250)

	// Lightness gradient
	for x := 0; x < 200; x++ {
		l := float64(x) / 200.0
		dc.SetHSL(240, 1.0, l) // Blue hue, full saturation
		dc.DrawLine(550+float64(x), 270, 550+float64(x), 285)
		dc.Stroke()
	}

	// Saturation gradient (HSL)
	for x := 0; x < 200; x++ {
		s := float64(x) / 200.0
		dc.SetHSL(240, s, 0.5) // Blue hue, 50% lightness
		dc.DrawLine(550+float64(x), 300, 550+float64(x), 315)
		dc.Stroke()
	}

	dc.SetRGB(0, 0, 0)
	dc.DrawString("L", 530, 280)
	dc.DrawString("S", 530, 310)
}

func demonstrateLABColorSpace(dc *advancegg.Context) {
	// LAB color space demonstration
	dc.SetRGB(0, 0, 0)
	dc.DrawString("LAB Color Space:", 50, 420)

	// A* axis (green to red)
	dc.DrawString("A* axis (Green to Red):", 50, 450)
	for x := 0; x < 200; x++ {
		a := (float64(x)/200.0)*200 - 100 // -100 to +100
		dc.SetLAB(50, a, 0)               // L=50, B=0
		dc.DrawLine(50+float64(x), 470, 50+float64(x), 485)
		dc.Stroke()
	}

	// B* axis (blue to yellow)
	dc.SetRGB(0, 0, 0)
	dc.DrawString("B* axis (Blue to Yellow):", 50, 510)
	for x := 0; x < 200; x++ {
		b := (float64(x)/200.0)*200 - 100 // -100 to +100
		dc.SetLAB(50, 0, b)               // L=50, A=0
		dc.DrawLine(50+float64(x), 530, 50+float64(x), 545)
		dc.Stroke()
	}

	// L* axis (lightness)
	dc.SetRGB(0, 0, 0)
	dc.DrawString("L* axis (Lightness):", 50, 570)
	for x := 0; x < 200; x++ {
		l := (float64(x) / 200.0) * 100 // 0 to 100
		dc.SetLAB(l, 0, 0)              // A=0, B=0 (neutral)
		dc.DrawLine(50+float64(x), 590, 50+float64(x), 605)
		dc.Stroke()
	}
}

func createCMYKExample() {
	dc := advancegg.NewContext(600, 400)
	dc.SetRGB(1, 1, 1)
	dc.Clear()

	dc.SetRGB(0, 0, 0)
	dc.DrawString("CMYK Color Model", 50, 50)

	// CMYK color separations
	separations := []struct {
		name       string
		c, m, y, k float64
		x, yPos    float64
	}{
		{"Cyan", 1, 0, 0, 0, 100, 100},
		{"Magenta", 0, 1, 0, 0, 200, 100},
		{"Yellow", 0, 0, 1, 0, 300, 100},
		{"Black", 0, 0, 0, 1, 400, 100},
		{"Red", 0, 1, 1, 0, 100, 200},
		{"Green", 1, 0, 1, 0, 200, 200},
		{"Blue", 1, 1, 0, 0, 300, 200},
		{"Gray", 0, 0, 0, 0.5, 400, 200},
	}

	for _, sep := range separations {
		dc.SetCMYK(sep.c, sep.m, sep.y, sep.k)
		dc.DrawCircle(sep.x, sep.yPos, 30)
		dc.Fill()

		dc.SetRGB(0, 0, 0)
		dc.DrawString(sep.name, sep.x-20, sep.yPos+50)
	}

	dc.SavePNG("images/filters/cmyk-example.png")
	fmt.Println("CMYK example saved as cmyk-example.png")
}

func createHSVExample() {
	dc := advancegg.NewContext(600, 400)
	dc.SetRGB(1, 1, 1)
	dc.Clear()

	dc.SetRGB(0, 0, 0)
	dc.DrawString("HSV Color Model", 50, 50)

	// Create HSV color variations
	baseHue := 240.0 // Blue

	// Saturation variations
	dc.DrawString("Saturation Variations:", 50, 100)
	for i := 0; i < 10; i++ {
		s := float64(i) / 9.0
		dc.SetHSV(baseHue, s, 1.0)
		dc.DrawCircle(80+float64(i)*40, 130, 15)
		dc.Fill()
	}

	// Value variations
	dc.SetRGB(0, 0, 0)
	dc.DrawString("Value Variations:", 50, 200)
	for i := 0; i < 10; i++ {
		v := float64(i) / 9.0
		dc.SetHSV(baseHue, 1.0, v)
		dc.DrawCircle(80+float64(i)*40, 230, 15)
		dc.Fill()
	}

	// Hue variations
	dc.SetRGB(0, 0, 0)
	dc.DrawString("Hue Variations:", 50, 300)
	for i := 0; i < 12; i++ {
		h := float64(i) * 30.0 // Every 30 degrees
		dc.SetHSV(h, 1.0, 1.0)
		dc.DrawCircle(80+float64(i)*35, 330, 15)
		dc.Fill()
	}

	dc.SavePNG("images/filters/hsv-example.png")
	fmt.Println("HSV example saved as hsv-example.png")
}

func createLABExample() {
	dc := advancegg.NewContext(600, 400)
	dc.SetRGB(1, 1, 1)
	dc.Clear()

	dc.SetRGB(0, 0, 0)
	dc.DrawString("LAB Color Space", 50, 50)

	// LAB color grid
	dc.DrawString("LAB Color Grid (L=50):", 50, 100)

	gridSize := 20
	for i := 0; i < gridSize; i++ {
		for j := 0; j < gridSize; j++ {
			a := (float64(i)/float64(gridSize-1))*200 - 100 // -100 to +100
			b := (float64(j)/float64(gridSize-1))*200 - 100 // -100 to +100

			dc.SetLAB(50, a, b)
			dc.DrawRectangle(50+float64(i)*15, 130+float64(j)*10, 14, 9)
			dc.Fill()
		}
	}

	dc.SavePNG("images/filters/lab-example.png")
	fmt.Println("LAB example saved as lab-example.png")
}
