package main

import (
	"fmt"
	"math"
	
	"github.com/GrandpaEJ/advancegg"
)

func main() {
	fmt.Println("Creating Text-on-Path examples...")
	
	// Basic text-on-path examples
	createBasicTextOnPath()
	
	// Advanced text-on-path examples
	createAdvancedTextOnPath()
	
	// Creative text effects
	createCreativeTextEffects()
	
	fmt.Println("Text-on-Path examples completed!")
}

func createBasicTextOnPath() {
	dc := advancegg.NewContext(800, 600)
	
	// White background
	dc.SetRGB(1, 1, 1)
	dc.Clear()
	
	// Title
	dc.SetRGB(0, 0, 0)
	dc.DrawString("Text-on-Path Examples", 50, 50)
	
	// Circular text
	dc.SetRGB(0.8, 0.2, 0.2)
	advancegg.DrawTextOnCircle(dc, "CIRCULAR TEXT AROUND A CIRCLE", 200, 150, 80)
	
	// Draw the circle outline for reference
	dc.SetRGBA(0.8, 0.2, 0.2, 0.3)
	dc.SetLineWidth(1)
	dc.DrawCircle(200, 150, 80)
	dc.Stroke()
	
	// Wave text
	dc.SetRGB(0.2, 0.8, 0.2)
	advancegg.DrawTextOnWave(dc, "This text follows a wave pattern!", 50, 300, 750, 30, 3)
	
	// Draw the wave outline for reference
	dc.SetRGBA(0.2, 0.8, 0.2, 0.3)
	dc.SetLineWidth(1)
	dc.MoveTo(50, 300)
	for x := 51; x <= 750; x++ {
		t := float64(x-50) / (750 - 50)
		y := 300 + 30*math.Sin(3*t*2*math.Pi)
		dc.LineTo(float64(x), y)
	}
	dc.Stroke()
	
	// Spiral text
	dc.SetRGB(0.2, 0.2, 0.8)
	advancegg.DrawTextOnSpiral(dc, "Spiral text goes round and round", 600, 150, 20, 80, 2)
	
	// Arc text
	dc.SetRGB(0.8, 0.5, 0.2)
	advancegg.DrawTextOnArc(dc, "Arc text follows a curve", 400, 450, 100, math.Pi, 2*math.Pi)
	
	// Draw the arc outline for reference
	dc.SetRGBA(0.8, 0.5, 0.2, 0.3)
	dc.SetLineWidth(1)
	dc.DrawArc(400, 450, 100, math.Pi, 2*math.Pi)
	dc.Stroke()
	
	// Bezier curve text
	dc.SetRGB(0.8, 0.2, 0.8)
	advancegg.DrawTextOnBezier(dc, "Bezier curve text", 100, 500, 300, 400, 500, 500)
	
	// Draw the Bezier curve for reference
	dc.SetRGBA(0.8, 0.2, 0.8, 0.3)
	dc.SetLineWidth(1)
	dc.MoveTo(100, 500)
	dc.QuadraticTo(300, 400, 500, 500)
	dc.Stroke()
	
	dc.SavePNG("images/text/images/basic-text-on-path.png")
	fmt.Println("Basic text-on-path demo saved as basic-text-on-path.png")
}

func createAdvancedTextOnPath() {
	dc := advancegg.NewContext(1000, 700)
	
	// Gradient background
	for y := 0; y < 700; y++ {
		t := float64(y) / 700.0
		dc.SetRGB(0.9+t*0.1, 0.95+t*0.05, 1.0)
		dc.DrawLine(0, float64(y), 1000, float64(y))
		dc.Stroke()
	}
	
	// Title
	dc.SetRGB(0, 0, 0)
	dc.DrawString("Advanced Text-on-Path Effects", 50, 50)
	
	// Multiple concentric circles
	dc.SetRGB(0.8, 0.2, 0.2)
	advancegg.DrawTextOnCircle(dc, "OUTER CIRCLE TEXT", 250, 200, 120)
	
	dc.SetRGB(0.2, 0.8, 0.2)
	advancegg.DrawTextOnCircle(dc, "MIDDLE CIRCLE", 250, 200, 90)
	
	dc.SetRGB(0.2, 0.2, 0.8)
	advancegg.DrawTextOnCircle(dc, "INNER", 250, 200, 60)
	
	// Complex wave patterns
	dc.SetRGB(0.8, 0.5, 0.2)
	advancegg.DrawTextOnWave(dc, "High frequency wave text", 50, 400, 950, 15, 8)
	
	dc.SetRGB(0.5, 0.8, 0.2)
	advancegg.DrawTextOnWave(dc, "Low frequency wave text", 50, 450, 950, 40, 2)
	
	// Multiple spirals
	dc.SetRGB(0.8, 0.2, 0.8)
	advancegg.DrawTextOnSpiral(dc, "Expanding spiral text gets bigger", 750, 200, 10, 100, 3)
	
	dc.SetRGB(0.2, 0.8, 0.8)
	advancegg.DrawTextOnSpiral(dc, "Contracting spiral", 750, 450, 80, 20, 2)
	
	// Multiple arcs forming a pattern
	angles := []float64{0, math.Pi/3, 2*math.Pi/3, math.Pi, 4*math.Pi/3, 5*math.Pi/3}
	colors := [][3]float64{
		{1, 0, 0}, {1, 0.5, 0}, {1, 1, 0},
		{0, 1, 0}, {0, 0, 1}, {0.5, 0, 1},
	}
	
	for i, angle := range angles {
		dc.SetRGB(colors[i][0], colors[i][1], colors[i][2])
		startAngle := angle
		endAngle := angle + math.Pi/2
		advancegg.DrawTextOnArc(dc, "Arc "+fmt.Sprintf("%d", i+1), 500, 550, 80, startAngle, endAngle)
	}
	
	// Complex Bezier curves
	dc.SetRGB(0, 0, 0)
	advancegg.DrawTextOnBezier(dc, "S-curve text", 100, 600, 200, 500, 300, 600)
	advancegg.DrawTextOnBezier(dc, "Reverse S-curve", 400, 600, 500, 650, 600, 600)
	
	dc.SavePNG("images/text/images/advanced-text-on-path.png")
	fmt.Println("Advanced text-on-path demo saved as advanced-text-on-path.png")
}

func createCreativeTextEffects() {
	dc := advancegg.NewContext(1200, 800)
	
	// Dark background for dramatic effect
	dc.SetRGB(0.1, 0.1, 0.2)
	dc.Clear()
	
	// Title
	dc.SetRGB(1, 1, 1)
	dc.DrawString("Creative Text Effects with Paths", 50, 50)
	
	// Logo-style circular text
	dc.SetRGB(1, 0.8, 0.2)
	advancegg.DrawTextOnCircle(dc, "★ ADVANCEGG GRAPHICS LIBRARY ★", 300, 200, 150)
	
	// Inner circle with smaller text
	dc.SetRGB(0.8, 0.8, 1)
	advancegg.DrawTextOnCircle(dc, "High Performance • Easy to Use", 300, 200, 100)
	
	// Decorative elements
	dc.SetRGB(1, 0.8, 0.2)
	dc.DrawCircle(300, 200, 10)
	dc.Fill()
	
	// Wave banner
	dc.SetRGB(0.2, 1, 0.8)
	advancegg.DrawTextOnWave(dc, "WAVE BANNER TEXT WITH STYLE", 50, 400, 1150, 50, 4)
	
	// Double wave effect
	dc.SetRGB(1, 0.4, 0.8)
	advancegg.DrawTextOnWave(dc, "Double wave effect", 50, 500, 1150, 30, 6)
	dc.SetRGB(0.8, 0.6, 1)
	advancegg.DrawTextOnWave(dc, "creates depth", 50, 530, 1150, 25, 6)
	
	// Spiral galaxy effect
	dc.SetRGB(1, 1, 0.5)
	advancegg.DrawTextOnSpiral(dc, "Spiral galaxy text effect creates amazing visual impact", 900, 200, 5, 120, 4)
	
	// Multiple arc rainbow
	rainbowColors := [][3]float64{
		{1, 0, 0},     // Red
		{1, 0.5, 0},   // Orange
		{1, 1, 0},     // Yellow
		{0, 1, 0},     // Green
		{0, 0, 1},     // Blue
		{0.3, 0, 0.7}, // Indigo
		{0.5, 0, 1},   // Violet
	}
	
	for i, color := range rainbowColors {
		dc.SetRGB(color[0], color[1], color[2])
		radius := 60 + float64(i)*15
		advancegg.DrawTextOnArc(dc, "RAINBOW", 600, 650, radius, math.Pi, 2*math.Pi)
	}
	
	// Artistic Bezier text
	dc.SetRGB(0.8, 1, 0.8)
	advancegg.DrawTextOnBezier(dc, "Artistic flowing text", 100, 650, 300, 550, 500, 650)
	
	dc.SetRGB(1, 0.8, 0.8)
	advancegg.DrawTextOnBezier(dc, "follows curves", 100, 700, 200, 750, 400, 700)
	
	dc.SetRGB(0.8, 0.8, 1)
	advancegg.DrawTextOnBezier(dc, "naturally", 100, 750, 150, 700, 300, 750)
	
	// Add some decorative stars
	stars := []struct{ x, y float64 }{
		{100, 100}, {200, 120}, {1100, 100}, {1050, 150},
		{50, 300}, {1150, 350}, {150, 600}, {1000, 580},
	}
	
	for _, star := range stars {
		dc.SetRGB(1, 1, 0.8)
		drawStar(dc, star.x, star.y, 8, 5)
	}
	
	// Information text
	dc.SetRGB(0.7, 0.7, 0.7)
	dc.DrawString("Text-on-Path enables creative typography and logo design", 50, 780)
	
	dc.SavePNG("images/text/images/creative-text-effects.png")
	fmt.Println("Creative text effects demo saved as creative-text-effects.png")
}

func drawStar(dc *advancegg.Context, centerX, centerY, outerRadius, innerRadius float64) {
	points := 5
	angle := 2 * math.Pi / float64(points*2)
	
	dc.MoveTo(centerX+outerRadius, centerY)
	
	for i := 1; i < points*2; i++ {
		radius := outerRadius
		if i%2 == 1 {
			radius = innerRadius
		}
		
		x := centerX + radius*math.Cos(float64(i)*angle)
		y := centerY + radius*math.Sin(float64(i)*angle)
		dc.LineTo(x, y)
	}
	
	dc.ClosePath()
	dc.Fill()
}

func createTextOnPathShowcase() {
	dc := advancegg.NewContext(1400, 1000)
	
	// Professional showcase background
	dc.SetRGB(0.95, 0.95, 0.98)
	dc.Clear()
	
	// Main title
	dc.SetRGB(0, 0, 0)
	dc.DrawString("AdvanceGG Text-on-Path Showcase", 50, 50)
	
	// Feature demonstrations
	features := []struct {
		title string
		demo  func()
		x, y  float64
	}{
		{"Circular Text", func() {
			dc.SetRGB(0.2, 0.4, 0.8)
			advancegg.DrawTextOnCircle(dc, "Perfect for logos and badges", 200, 200, 80)
		}, 50, 120},
		{"Wave Text", func() {
			dc.SetRGB(0.8, 0.2, 0.4)
			advancegg.DrawTextOnWave(dc, "Great for banners and headers", 400, 200, 800, 25, 3)
		}, 350, 120},
		{"Spiral Text", func() {
			dc.SetRGB(0.4, 0.8, 0.2)
			advancegg.DrawTextOnSpiral(dc, "Eye-catching spiral effects", 1000, 200, 20, 70, 2)
		}, 850, 120},
		{"Arc Text", func() {
			dc.SetRGB(0.8, 0.6, 0.2)
			advancegg.DrawTextOnArc(dc, "Curved text for design", 200, 500, 60, 0, math.Pi)
		}, 50, 400},
		{"Bezier Text", func() {
			dc.SetRGB(0.6, 0.2, 0.8)
			advancegg.DrawTextOnBezier(dc, "Smooth flowing curves", 600, 500, 800, 400, 1000, 500)
		}, 550, 400},
	}
	
	for _, feature := range features {
		// Feature title
		dc.SetRGB(0, 0, 0)
		dc.DrawString(feature.title, feature.x, feature.y)
		
		// Demonstrate the feature
		feature.demo()
	}
	
	// Usage examples
	dc.SetRGB(0, 0, 0)
	dc.DrawString("Perfect for:", 50, 700)
	
	examples := []string{
		"• Logo design and branding",
		"• Decorative headers and banners",
		"• Artistic typography effects",
		"• Circular badges and stamps",
		"• Creative poster design",
		"• Web graphics and UI elements",
	}
	
	for i, example := range examples {
		dc.DrawString(example, 70, 730+float64(i)*25)
	}
	
	// Technical features
	dc.SetRGB(0, 0, 0)
	dc.DrawString("Technical Features:", 500, 700)
	
	techFeatures := []string{
		"• Automatic character spacing",
		"• Multiple alignment options",
		"• Smooth curve interpolation",
		"• Rotation and positioning",
		"• High-quality rendering",
		"• Easy-to-use API",
	}
	
	for i, feature := range techFeatures {
		dc.DrawString(feature, 520, 730+float64(i)*25)
	}
	
	dc.SavePNG("images/text/images/text-on-path-showcase.png")
	fmt.Println("Text-on-path showcase saved as text-on-path-showcase.png")
}
