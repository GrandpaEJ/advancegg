package main

import (
	"fmt"

	"github.com/GrandpaEJ/advancegg"
)

func main() {
	dc := advancegg.NewContext(800, 600)

	// Light background
	dc.SetRGB(0.95, 0.95, 0.95)
	dc.Clear()

	// Title
	dc.SetRGB(0, 0, 0)
	dc.DrawString("Shadow Effects Demonstration", 50, 50)

	// Demonstrate various shadow effects
	demonstrateBasicShadows(dc)
	demonstrateTextShadows(dc)
	demonstrateShapeShadows(dc)
	demonstrateAdvancedShadows(dc)

	// Save the demonstration
	dc.SavePNG("shadow-effects.png")
	fmt.Println("Shadow effects demo saved as shadow-effects.png")

	fmt.Println("Shadow effects examples completed!")
}

func demonstrateBasicShadows(dc *advancegg.Context) {
	dc.SetRGB(0, 0, 0)
	dc.DrawString("Basic Shadows:", 50, 100)

	// Simple drop shadow
	dc.SetShadowRGBA(5, 5, 0, 0, 0, 0, 0.3) // Black shadow with 30% opacity
	dc.SetRGB(0, 0.5, 1)                    // Blue color
	dc.DrawCircleWithShadow(150, 150, 30)

	// Shadow with blur
	dc.SetShadowRGBA(3, 3, 5, 0, 0, 0, 0.4) // Blurred shadow
	dc.SetRGB(1, 0.5, 0)                    // Orange color
	dc.DrawRectangleWithShadow(250, 120, 60, 60)

	// Colored shadow
	dc.SetShadowRGBA(4, 4, 2, 1, 0, 0, 0.5) // Red shadow
	dc.SetRGB(1, 1, 0)                      // Yellow color
	dc.DrawRoundedRectangleWithShadow(370, 120, 60, 60, 10)

	// Large blur shadow
	dc.SetShadowRGBA(2, 2, 10, 0, 0, 0, 0.2) // Large blur
	dc.SetRGB(0, 1, 0)                       // Green color
	dc.DrawEllipseWithShadow(520, 150, 40, 25)

	dc.ClearShadow() // Clear shadow for next section
}

func demonstrateTextShadows(dc *advancegg.Context) {
	dc.SetRGB(0, 0, 0)
	dc.DrawString("Text Shadows:", 50, 220)

	// Simple text shadow
	dc.SetShadowRGBA(2, 2, 0, 0, 0, 0, 0.5)
	dc.SetRGB(0, 0, 1)
	dc.DrawStringWithShadow("Drop Shadow", 50, 250)

	// Blurred text shadow
	dc.SetShadowRGBA(1, 1, 3, 0, 0, 0, 0.6)
	dc.SetRGB(1, 0, 0)
	dc.DrawStringWithShadow("Blurred Shadow", 200, 250)

	// Colored text shadow
	dc.SetShadowRGBA(3, 3, 1, 0, 0, 1, 0.4) // Blue shadow
	dc.SetRGB(1, 1, 1)                      // White text
	dc.DrawStringWithShadow("Colored Shadow", 350, 250)

	// Multiple shadow effect (simulated)
	dc.SetShadowRGBA(4, 4, 2, 1, 0, 0, 0.3) // Red shadow
	dc.SetRGB(1, 1, 0)
	dc.DrawStringWithShadow("Layered", 520, 250)

	dc.SetShadowRGBA(2, 2, 1, 0, 0, 1, 0.3) // Blue shadow
	dc.DrawStringWithShadow("Layered", 520, 250)

	dc.ClearShadow()
}

func demonstrateShapeShadows(dc *advancegg.Context) {
	dc.SetRGB(0, 0, 0)
	dc.DrawString("Shape Shadows:", 50, 300)

	// Card-like shadow
	dc.SetShadowRGBA(0, 4, 8, 0, 0, 0, 0.15)
	dc.SetRGB(1, 1, 1)
	dc.DrawRoundedRectangleWithShadow(50, 330, 120, 80, 8)

	// Button-like shadow
	dc.SetShadowRGBA(0, 2, 4, 0, 0, 0, 0.2)
	dc.SetRGB(0.2, 0.6, 1)
	dc.DrawRoundedRectangleWithShadow(200, 340, 100, 40, 20)

	// Floating element shadow
	dc.SetShadowRGBA(0, 8, 16, 0, 0, 0, 0.1)
	dc.SetRGB(0.9, 0.9, 0.9)
	dc.DrawRoundedRectangleWithShadow(330, 325, 100, 70, 5)

	// Inner shadow effect (simulated with dark border)
	dc.ClearShadow()
	dc.SetRGB(0.8, 0.8, 0.8)
	dc.DrawRoundedRectangleWithShadow(460, 330, 80, 80, 10)

	dc.ClearShadow()
}

func demonstrateAdvancedShadows(dc *advancegg.Context) {
	dc.SetRGB(0, 0, 0)
	dc.DrawString("Advanced Shadow Effects:", 50, 450)

	// Long shadow effect
	dc.SetShadowRGBA(15, 15, 0, 0, 0, 0, 0.2)
	dc.SetRGB(1, 0.3, 0.3)
	dc.DrawCircleWithShadow(120, 500, 25)

	// Soft glow effect
	dc.SetShadowRGBA(0, 0, 15, 0.3, 0.7, 1, 0.6) // Blue glow
	dc.SetRGB(1, 1, 1)
	dc.DrawCircleWithShadow(250, 500, 20)

	// Neon effect (multiple glows)
	dc.SetShadowRGBA(0, 0, 8, 1, 0, 1, 0.8) // Magenta glow
	dc.SetRGB(1, 1, 1)
	dc.DrawStringWithShadow("NEON", 320, 510)

	dc.SetShadowRGBA(0, 0, 4, 1, 0, 1, 1) // Stronger inner glow
	dc.DrawStringWithShadow("NEON", 320, 510)

	// Embossed effect (light and dark shadows)
	dc.SetShadowRGBA(-1, -1, 0, 1, 1, 1, 0.8) // Light shadow (top-left)
	dc.SetRGB(0.7, 0.7, 0.7)
	dc.DrawRoundedRectangleWithShadow(450, 480, 80, 40, 5)

	dc.SetShadowRGBA(1, 1, 0, 0, 0, 0, 0.3) // Dark shadow (bottom-right)
	dc.DrawRoundedRectangleWithShadow(450, 480, 80, 40, 5)

	// Create a complex shadow composition
	createShadowComposition(dc)

	dc.ClearShadow()
}

func createShadowComposition(dc *advancegg.Context) {
	// Create a card stack effect
	x, y := 580.0, 460.0

	// Bottom card (largest shadow)
	dc.SetShadowRGBA(2, 6, 12, 0, 0, 0, 0.1)
	dc.SetRGB(0.9, 0.9, 0.9)
	dc.DrawRoundedRectangleWithShadow(x+8, y+8, 100, 60, 5)

	// Middle card
	dc.SetShadowRGBA(1, 4, 8, 0, 0, 0, 0.15)
	dc.SetRGB(0.95, 0.95, 0.95)
	dc.DrawRoundedRectangleWithShadow(x+4, y+4, 100, 60, 5)

	// Top card
	dc.SetShadowRGBA(0, 2, 4, 0, 0, 0, 0.2)
	dc.SetRGB(1, 1, 1)
	dc.DrawRoundedRectangleWithShadow(x, y, 100, 60, 5)

	// Add some content to the top card
	dc.ClearShadow()
	dc.SetRGB(0.3, 0.3, 0.3)
	dc.DrawString("Card", x+35, y+35)
}

// Additional helper functions for more complex shadow effects

func createDropShadowExample(dc *advancegg.Context) {
	// Create a realistic drop shadow example
	dc.SetRGB(1, 1, 1)
	dc.Clear()

	// Material Design inspired shadows
	elevations := []struct {
		name   string
		offset float64
		blur   float64
		alpha  float64
		x, y   float64
	}{
		{"Elevation 1", 1, 3, 0.12, 100, 100},
		{"Elevation 2", 2, 6, 0.16, 250, 100},
		{"Elevation 4", 4, 12, 0.19, 400, 100},
		{"Elevation 8", 8, 24, 0.25, 550, 100},
	}

	for _, elev := range elevations {
		dc.SetShadowRGBA(0, elev.offset, elev.blur, 0, 0, 0, elev.alpha)
		dc.SetRGB(1, 1, 1)
		dc.DrawRoundedRectangleWithShadow(elev.x, elev.y, 80, 80, 4)

		dc.ClearShadow()
		dc.SetRGB(0, 0, 0)
		dc.DrawString(elev.name, elev.x+5, elev.y+100)
	}

	dc.SavePNG("material-shadows.png")
	fmt.Println("Material Design shadows saved as material-shadows.png")
}
