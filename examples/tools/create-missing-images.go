package main

import (
	"fmt"
	"image/color"
	"math"
	"math/rand"

	"github.com/GrandpaEJ/advancegg"
)

func main() {
	fmt.Println("Creating missing images for README...")

	createDataVisualizationDemo()
	createGameGraphicsDemo()
	createFilterShowcase()
	createAdvancedStrokeEffects()

	fmt.Println("All missing images created successfully!")
}

func createDataVisualizationDemo() {
	dc := advancegg.NewContext(800, 600)

	// Professional gradient background
	for y := 0; y < 600; y++ {
		t := float64(y) / 600.0
		dc.SetRGB(0.95+t*0.05, 0.96+t*0.04, 0.98+t*0.02)
		dc.DrawLine(0, float64(y), 800, float64(y))
		dc.Stroke()
	}

	// Title
	dc.SetRGB(0.1, 0.1, 0.2)
	dc.DrawString("Data Visualization Dashboard", 50, 50)

	// Bar Chart
	data := []float64{85, 92, 78, 96, 88, 94}
	colors := []color.Color{
		color.RGBA{59, 130, 246, 255}, // Blue
		color.RGBA{16, 185, 129, 255}, // Green
		color.RGBA{245, 158, 11, 255}, // Yellow
		color.RGBA{239, 68, 68, 255},  // Red
		color.RGBA{139, 92, 246, 255}, // Purple
		color.RGBA{6, 182, 212, 255},  // Cyan
	}

	// Draw bars with shadows
	for i, value := range data {
		x := 80 + float64(i)*100
		height := value * 3

		// Shadow
		dc.SetRGBA(0, 0, 0, 0.1)
		dc.DrawRoundedRectangle(x+3, 450-height+3, 60, height, 5)
		dc.Fill()

		// Main bar
		dc.SetColor(colors[i])
		dc.DrawRoundedRectangle(x, 450-height, 60, height, 5)
		dc.Fill()

		// Value labels
		dc.SetRGB(0.2, 0.2, 0.2)
		dc.DrawStringAnchored(fmt.Sprintf("%.0f%%", value), x+30, 430-height, 0.5, 0.5)

		// Month labels
		months := []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun"}
		dc.SetRGB(0.4, 0.4, 0.4)
		dc.DrawStringAnchored(months[i], x+30, 470, 0.5, 0.5)
	}

	// Line Chart
	dc.SetRGB(0.1, 0.1, 0.2)
	dc.DrawString("Trend Analysis", 450, 100)

	// Line chart data
	lineData := []float64{65, 72, 68, 85, 78, 92, 88, 95}

	// Draw line
	dc.SetRGB(0.2, 0.6, 1)
	dc.SetLineWidth(3)
	for i, value := range lineData {
		x := 480 + float64(i)*35
		y := 300 - value*1.5

		if i == 0 {
			dc.MoveTo(x, y)
		} else {
			dc.LineTo(x, y)
		}
	}
	dc.Stroke()

	// Draw points
	for i, value := range lineData {
		x := 480 + float64(i)*35
		y := 300 - value*1.5

		dc.SetRGB(1, 1, 1)
		dc.DrawCircle(x, y, 6)
		dc.Fill()

		dc.SetRGB(0.2, 0.6, 1)
		dc.DrawCircle(x, y, 4)
		dc.Fill()
	}

	// Pie Chart
	dc.SetRGB(0.1, 0.1, 0.2)
	dc.DrawString("Distribution", 50, 350)

	pieData := []float64{30, 25, 20, 15, 10}
	pieColors := []color.Color{
		color.RGBA{255, 99, 132, 255},
		color.RGBA{54, 162, 235, 255},
		color.RGBA{255, 205, 86, 255},
		color.RGBA{75, 192, 192, 255},
		color.RGBA{153, 102, 255, 255},
	}

	centerX, centerY := 150.0, 520.0
	radius := 60.0
	startAngle := 0.0

	for i, value := range pieData {
		angle := value / 100.0 * 2 * math.Pi

		dc.SetColor(pieColors[i])
		dc.MoveTo(centerX, centerY)
		dc.DrawArc(centerX, centerY, radius, startAngle, startAngle+angle)
		dc.ClosePath()
		dc.Fill()

		startAngle += angle
	}

	dc.SavePNG("images/tools/images/data-visualization-demo.png")
	fmt.Println("Created data-visualization-demo.png")
}

func createGameGraphicsDemo() {
	dc := advancegg.NewContext(800, 600)

	// Game-style dark background
	dc.SetRGB(0.05, 0.05, 0.1)
	dc.Clear()

	// Title
	dc.SetRGB(1, 1, 0.8)
	dc.DrawString("Game Graphics Showcase", 50, 50)

	// Draw a simple character sprite
	// Head
	dc.SetRGB(1, 0.8, 0.6)
	dc.DrawCircle(150, 150, 30)
	dc.Fill()

	// Eyes
	dc.SetRGB(0, 0, 0)
	dc.DrawCircle(140, 140, 5)
	dc.Fill()
	dc.DrawCircle(160, 140, 5)
	dc.Fill()

	// Body
	dc.SetRGB(0.2, 0.4, 0.8)
	dc.DrawRectangle(130, 180, 40, 60)
	dc.Fill()

	// Arms
	dc.SetRGB(1, 0.8, 0.6)
	dc.DrawRectangle(110, 190, 15, 40)
	dc.Fill()
	dc.DrawRectangle(175, 190, 15, 40)
	dc.Fill()

	// Legs
	dc.SetRGB(0.1, 0.2, 0.6)
	dc.DrawRectangle(135, 240, 12, 40)
	dc.Fill()
	dc.DrawRectangle(153, 240, 12, 40)
	dc.Fill()

	// UI Elements
	// Health bar
	dc.SetRGB(0.8, 0.8, 0.8)
	dc.DrawRoundedRectangle(50, 350, 200, 20, 10)
	dc.Fill()

	dc.SetRGB(0.2, 0.8, 0.2)
	dc.DrawRoundedRectangle(52, 352, 150, 16, 8)
	dc.Fill()

	dc.SetRGB(1, 1, 1)
	dc.DrawString("Health: 75%", 60, 365)

	// Score display
	dc.SetRGB(1, 1, 0)
	dc.DrawString("Score: 12,450", 50, 400)

	// Particle effects (stars)
	for i := 0; i < 20; i++ {
		x := 400 + rand.Float64()*300
		y := 100 + rand.Float64()*200
		size := 2 + rand.Float64()*4

		dc.SetRGBA(1, 1, 0.8, 0.7)
		dc.DrawCircle(x, y, size)
		dc.Fill()
	}

	// Game objects (coins)
	for i := 0; i < 5; i++ {
		x := 450 + float64(i)*60
		y := 350.0

		// Coin shadow
		dc.SetRGBA(0, 0, 0, 0.3)
		dc.DrawCircle(x+2, y+2, 15)
		dc.Fill()

		// Coin
		dc.SetRGB(1, 0.8, 0)
		dc.DrawCircle(x, y, 15)
		dc.Fill()

		// Coin detail
		dc.SetRGB(0.8, 0.6, 0)
		dc.DrawCircle(x, y, 10)
		dc.Stroke()
	}

	// Platform
	dc.SetRGB(0.3, 0.2, 0.1)
	dc.DrawRoundedRectangle(400, 450, 300, 40, 20)
	dc.Fill()

	// Grass texture on platform
	dc.SetRGB(0.2, 0.6, 0.2)
	for i := 0; i < 30; i++ {
		x := 410 + float64(i)*9
		dc.DrawLine(x, 450, x, 445)
		dc.Stroke()
	}

	dc.SavePNG("images/tools/images/game-graphics-demo.png")
	fmt.Println("Created game-graphics-demo.png")
}

func createFilterShowcase() {
	dc := advancegg.NewContext(800, 600)

	// White background
	dc.SetRGB(1, 1, 1)
	dc.Clear()

	// Title
	dc.SetRGB(0, 0, 0)
	dc.DrawString("Image Filter Showcase", 50, 50)

	// Create a sample image to filter
	sampleDC := advancegg.NewContext(150, 100)

	// Gradient background
	for y := 0; y < 100; y++ {
		t := float64(y) / 100.0
		sampleDC.SetRGB(0.2+t*0.6, 0.4+t*0.4, 0.8+t*0.2)
		sampleDC.DrawLine(0, float64(y), 150, float64(y))
		sampleDC.Stroke()
	}

	// Add some shapes
	sampleDC.SetRGB(1, 0.5, 0)
	sampleDC.DrawCircle(75, 50, 30)
	sampleDC.Fill()

	sampleDC.SetRGB(0.8, 0.2, 0.8)
	sampleDC.DrawRectangle(20, 20, 40, 30)
	sampleDC.Fill()

	originalImg := sampleDC.Image()

	// Display original
	dc.DrawImage(originalImg, 50, 100)
	dc.SetRGB(0, 0, 0)
	dc.DrawString("Original", 50, 220)

	// Blur effect simulation
	dc.DrawImage(originalImg, 250, 100)
	dc.SetRGBA(1, 1, 1, 0.3)
	dc.DrawRectangle(250, 100, 150, 100)
	dc.Fill()
	dc.SetRGB(0, 0, 0)
	dc.DrawString("Blur", 250, 220)

	// Grayscale effect simulation
	dc.SetRGB(0.5, 0.5, 0.5)
	dc.DrawRectangle(450, 100, 150, 100)
	dc.Fill()
	dc.SetRGB(0, 0, 0)
	dc.DrawString("Grayscale", 450, 220)

	// Sepia effect simulation
	dc.SetRGB(0.8, 0.7, 0.5)
	dc.DrawRectangle(50, 280, 150, 100)
	dc.Fill()
	dc.SetRGB(0, 0, 0)
	dc.DrawString("Sepia", 50, 400)

	// Edge detection simulation
	dc.SetRGB(1, 1, 1)
	dc.DrawRectangle(250, 280, 150, 100)
	dc.Fill()
	dc.SetRGB(0, 0, 0)
	dc.SetLineWidth(2)
	dc.DrawCircle(325, 330, 30)
	dc.Stroke()
	dc.DrawRectangle(270, 300, 40, 30)
	dc.Stroke()
	dc.DrawString("Edge Detection", 250, 400)

	// Emboss effect simulation
	dc.SetRGB(0.7, 0.7, 0.7)
	dc.DrawRectangle(450, 280, 150, 100)
	dc.Fill()
	dc.SetRGB(0.9, 0.9, 0.9)
	dc.DrawCircle(525, 330, 30)
	dc.Fill()
	dc.SetRGB(0.5, 0.5, 0.5)
	dc.DrawRectangle(470, 300, 40, 30)
	dc.Fill()
	dc.SetRGB(0, 0, 0)
	dc.DrawString("Emboss", 450, 400)

	// Filter labels
	dc.SetRGB(0.2, 0.2, 0.2)
	dc.DrawString("15+ Professional Filters Available", 50, 450)
	dc.DrawString("• Blur • Sharpen • Grayscale • Sepia • Invert", 50, 470)
	dc.DrawString("• Edge Detection • Emboss • Brightness • Contrast", 50, 490)
	dc.DrawString("• Saturation • Hue Rotate • Pixelate • Vignette", 50, 510)

	dc.SavePNG("images/tools/images/filter-showcase.png")
	fmt.Println("Created filter-showcase.png")
}

func createAdvancedStrokeEffects() {
	dc := advancegg.NewContext(800, 600)

	// Light background
	dc.SetRGB(0.98, 0.98, 1)
	dc.Clear()

	// Title
	dc.SetRGB(0, 0, 0)
	dc.DrawString("Advanced Stroke Effects", 50, 50)

	// Dashed line examples
	dc.SetRGB(0.8, 0.2, 0.2)
	dc.SetLineWidth(4)
	dc.DrawDashedLine(50, 120, 350, 120, []float64{10, 5})
	dc.SetRGB(0, 0, 0)
	dc.DrawString("Dashed Strokes", 50, 110)

	// Complex dash pattern
	dc.SetRGB(0.2, 0.8, 0.2)
	dc.DrawDashedLine(50, 160, 350, 160, []float64{20, 5, 5, 5})
	dc.SetRGB(0, 0, 0)
	dc.DrawString("Complex Patterns", 50, 150)

	// Gradient stroke simulation
	dc.SetRGB(0, 0, 0)
	dc.DrawString("Gradient Strokes", 50, 190)

	// Simulate gradient by drawing segments
	for i := 0; i < 20; i++ {
		t := float64(i) / 19.0
		r := 1.0 * (1 - t)
		b := 1.0 * t
		dc.SetRGB(r, 0, b)
		dc.SetLineWidth(6)
		x1 := 50 + float64(i)*15
		x2 := 50 + float64(i+1)*15
		dc.DrawLine(x1, 200, x2, 200)
		dc.Stroke()
	}

	// Tapered stroke simulation
	dc.SetRGB(0, 0, 0)
	dc.DrawString("Tapered Strokes", 50, 230)

	// Simulate taper by varying line width
	for i := 0; i < 30; i++ {
		t := float64(i) / 29.0
		width := 10*(1-t) + 1
		dc.SetRGB(0.8, 0.2, 0.8)
		dc.SetLineWidth(width)
		x1 := 50 + float64(i)*10
		x2 := 50 + float64(i+1)*10
		dc.DrawLine(x1, 250, x2, 250)
		dc.Stroke()
	}

	// Artistic brush strokes
	dc.SetRGB(0, 0, 0)
	dc.DrawString("Artistic Effects", 450, 110)

	// Calligraphy-style strokes
	dc.SetRGB(0.2, 0.2, 0.2)
	dc.SetLineWidth(15)
	dc.DrawTaperedLine(480, 140, 600, 120, 1.0, 0.2)

	dc.SetLineWidth(12)
	dc.DrawTaperedLine(520, 160, 650, 180, 0.3, 1.0)

	dc.SetLineWidth(18)
	dc.DrawTaperedLine(460, 200, 580, 190, 0.8, 0.1)

	// Neon effect simulation
	dc.SetRGB(0, 0, 0)
	dc.DrawString("Neon Effects", 450, 250)

	// Outer glow
	dc.SetRGBA(0, 1, 1, 0.3)
	dc.SetLineWidth(20)
	dc.DrawLine(480, 280, 650, 280)
	dc.Stroke()

	// Inner glow
	dc.SetRGBA(0, 1, 1, 0.6)
	dc.SetLineWidth(12)
	dc.DrawLine(480, 280, 650, 280)
	dc.Stroke()

	// Core
	dc.SetRGB(1, 1, 1)
	dc.SetLineWidth(4)
	dc.DrawLine(480, 280, 650, 280)
	dc.Stroke()

	// Feature list
	dc.SetRGB(0.2, 0.2, 0.2)
	dc.DrawString("Advanced Stroke Features:", 50, 350)
	dc.DrawString("• Dashed and dotted patterns", 70, 380)
	dc.DrawString("• Gradient stroke colors", 70, 400)
	dc.DrawString("• Tapered calligraphy effects", 70, 420)
	dc.DrawString("• Custom line caps and joins", 70, 440)
	dc.DrawString("• Artistic brush simulations", 70, 460)
	dc.DrawString("• Neon and glow effects", 70, 480)

	dc.SavePNG("images/tools/images/advanced-stroke-effects.png")
	fmt.Println("Created advanced-stroke-effects.png")
}
