package main

import (
	"fmt"
	"image/color"
	"math"

	"github.com/GrandpaEJ/advancegg"
)

func main() {
	fmt.Println("Creating Layer Compositing examples...")
	
	// Create comprehensive blend mode demo
	createBlendModeDemo()
	
	// Create layer effects demo
	createLayerEffectsDemo()
	
	// Create complex compositing demo
	createComplexCompositingDemo()
	
	fmt.Println("Layer Compositing examples completed!")
}

func createBlendModeDemo() {
	dc := advancegg.NewContext(1200, 900)
	
	// White background
	dc.SetRGB(1, 1, 1)
	dc.Clear()
	
	// Title
	dc.SetRGB(0, 0, 0)
	dc.DrawString("Layer Blend Modes Demo", 20, 30)
	
	// Enable layer system
	dc.EnableLayers()
	lm := dc.GetLayerManager()
	
	// Create base layer with gradient background
	baseLayer := lm.GetActiveLayer()
	baseLayer.Name = "Base"
	
	// Draw gradient background
	for y := 0; y < 900; y++ {
		for x := 0; x < 1200; x++ {
			r := float64(x) / 1200.0
			g := float64(y) / 900.0
			b := 0.5
			baseLayer.Image.SetRGBA(x, y, color.RGBA{
				uint8(r * 255),
				uint8(g * 255),
				uint8(b * 255),
				255,
			})
		}
	}
	
	// Define blend modes to demonstrate
	blendModes := []struct {
		mode advancegg.BlendMode
		name string
		x, y int
	}{
		{advancegg.BlendModeNormal, "Normal", 50, 80},
		{advancegg.BlendModeMultiply, "Multiply", 250, 80},
		{advancegg.BlendModeScreen, "Screen", 450, 80},
		{advancegg.BlendModeOverlay, "Overlay", 650, 80},
		{advancegg.BlendModeSoftLight, "Soft Light", 850, 80},
		{advancegg.BlendModeHardLight, "Hard Light", 1050, 80},
		{advancegg.BlendModeColorDodge, "Color Dodge", 50, 280},
		{advancegg.BlendModeColorBurn, "Color Burn", 250, 280},
		{advancegg.BlendModeDarken, "Darken", 450, 280},
		{advancegg.BlendModeLighten, "Lighten", 650, 280},
		{advancegg.BlendModeDifference, "Difference", 850, 280},
		{advancegg.BlendModeExclusion, "Exclusion", 1050, 280},
	}
	
	// Create layers for each blend mode
	for i, bm := range blendModes {
		layer := lm.AddLayer(bm.name)
		layer.SetBlendMode(bm.mode)
		layer.SetOpacity(0.8)
		
		// Set active layer for drawing
		lm.SetActiveLayer(i + 1)
		
		// Draw overlapping shapes
		dc.SetRGB(1, 0.2, 0.2) // Red circle
		dc.DrawCircle(float64(bm.x+60), float64(bm.y+60), 40)
		dc.Fill()
		
		dc.SetRGB(0.2, 1, 0.2) // Green circle
		dc.DrawCircle(float64(bm.x+80), float64(bm.y+80), 40)
		dc.Fill()
		
		dc.SetRGB(0.2, 0.2, 1) // Blue circle
		dc.DrawCircle(float64(bm.x+70), float64(bm.y+100), 40)
		dc.Fill()
		
		// Add label
		dc.SetRGB(0, 0, 0)
		dc.DrawString(bm.name, float64(bm.x), float64(bm.y+160))
	}
	
	// Composite and save
	result := lm.Composite()
	advancegg.SavePNG("blend-modes-demo.png", result)
	fmt.Println("Blend modes demo saved as blend-modes-demo.png")
}

func createLayerEffectsDemo() {
	dc := advancegg.NewContext(800, 600)
	
	// Dark background
	dc.SetRGB(0.1, 0.1, 0.2)
	dc.Clear()
	
	// Enable layers
	dc.EnableLayers()
	lm := dc.GetLayerManager()
	
	// Background layer
	bgLayer := lm.GetActiveLayer()
	bgLayer.Name = "Background"
	bgLayer.Fill(color.RGBA{20, 20, 40, 255})
	
	// Create text layer with glow effect
	textLayer := lm.AddLayer("Text")
	textLayer.SetBlendMode(advancegg.BlendModeScreen)
	lm.SetActiveLayer(1)
	
	// Draw text
	dc.SetRGB(1, 1, 0.8)
	dc.DrawString("LAYER EFFECTS", 200, 200)
	
	// Create glow layer
	glowLayer := lm.AddLayer("Glow")
	glowLayer.SetBlendMode(advancegg.BlendModeScreen)
	glowLayer.SetOpacity(0.6)
	lm.SetActiveLayer(2)
	
	// Draw glow effect (multiple blurred copies)
	for i := 0; i < 5; i++ {
		offset := float64(i * 2)
		dc.SetRGBA(1, 1, 0.2, 0.3)
		dc.DrawString("LAYER EFFECTS", 200-offset, 200-offset)
		dc.DrawString("LAYER EFFECTS", 200+offset, 200+offset)
	}
	
	// Create shadow layer
	shadowLayer := lm.AddLayer("Shadow")
	shadowLayer.SetBlendMode(advancegg.BlendModeMultiply)
	shadowLayer.SetOpacity(0.7)
	lm.SetActiveLayer(3)
	
	// Draw shadow
	dc.SetRGB(0, 0, 0)
	dc.DrawString("LAYER EFFECTS", 205, 205)
	
	// Composite and save
	result := lm.Composite()
	advancegg.SavePNG("layer-effects-demo.png", result)
	fmt.Println("Layer effects demo saved as layer-effects-demo.png")
}

func createComplexCompositingDemo() {
	dc := advancegg.NewContext(1000, 800)
	
	// Enable layers
	dc.EnableLayers()
	lm := dc.GetLayerManager()
	
	// Sky background
	skyLayer := lm.GetActiveLayer()
	skyLayer.Name = "Sky"
	drawSkyGradient(dc, skyLayer)
	
	// Mountains layer
	mountainLayer := lm.AddLayer("Mountains")
	mountainLayer.SetBlendMode(advancegg.BlendModeMultiply)
	mountainLayer.SetOpacity(0.8)
	lm.SetActiveLayer(1)
	drawMountains(dc)
	
	// Clouds layer
	cloudLayer := lm.AddLayer("Clouds")
	cloudLayer.SetBlendMode(advancegg.BlendModeScreen)
	cloudLayer.SetOpacity(0.7)
	lm.SetActiveLayer(2)
	drawClouds(dc)
	
	// Sun layer
	sunLayer := lm.AddLayer("Sun")
	sunLayer.SetBlendMode(advancegg.BlendModeColorDodge)
	sunLayer.SetOpacity(0.9)
	lm.SetActiveLayer(3)
	drawSun(dc)
	
	// Water layer
	waterLayer := lm.AddLayer("Water")
	waterLayer.SetBlendMode(advancegg.BlendModeOverlay)
	waterLayer.SetOpacity(0.6)
	lm.SetActiveLayer(4)
	drawWater(dc)
	
	// Reflection layer
	reflectionLayer := lm.AddLayer("Reflection")
	reflectionLayer.SetBlendMode(advancegg.BlendModeSoftLight)
	reflectionLayer.SetOpacity(0.4)
	lm.SetActiveLayer(5)
	drawReflection(dc)
	
	// Composite and save
	result := lm.Composite()
	advancegg.SavePNG("complex-compositing-demo.png", result)
	fmt.Println("Complex compositing demo saved as complex-compositing-demo.png")
}

// Helper functions for drawing scene elements

func drawSkyGradient(dc *advancegg.Context, layer *advancegg.Layer) {
	// Draw sky gradient from blue to orange
	for y := 0; y < 800; y++ {
		t := float64(y) / 800.0
		r := 0.2 + t*0.6  // Blue to orange
		g := 0.4 + t*0.4
		b := 0.8 - t*0.6
		
		for x := 0; x < 1000; x++ {
			layer.Image.SetRGBA(x, y, color.RGBA{
				uint8(r * 255),
				uint8(g * 255),
				uint8(b * 255),
				255,
			})
		}
	}
}

func drawMountains(dc *advancegg.Context) {
	dc.SetRGB(0.3, 0.2, 0.4)
	
	// Draw mountain silhouettes
	for i := 0; i < 5; i++ {
		baseX := float64(i * 200)
		height := 200 + float64(i*50)
		
		dc.MoveTo(baseX, 800)
		dc.LineTo(baseX+100, 800-height)
		dc.LineTo(baseX+200, 800)
		dc.ClosePath()
		dc.Fill()
	}
}

func drawClouds(dc *advancegg.Context) {
	dc.SetRGBA(1, 1, 1, 0.8)
	
	// Draw fluffy clouds
	for i := 0; i < 8; i++ {
		x := float64(i*130 + 50)
		y := float64(100 + i*20)
		
		// Multiple overlapping circles for cloud effect
		for j := 0; j < 5; j++ {
			radius := 20 + float64(j*5)
			offsetX := float64(j * 15)
			dc.DrawCircle(x+offsetX, y, radius)
			dc.Fill()
		}
	}
}

func drawSun(dc *advancegg.Context) {
	dc.SetRGB(1, 0.9, 0.3)
	
	// Sun disc
	dc.DrawCircle(800, 150, 60)
	dc.Fill()
	
	// Sun rays
	for i := 0; i < 12; i++ {
		angle := float64(i) * 2 * math.Pi / 12
		x1 := 800 + 80*math.Cos(angle)
		y1 := 150 + 80*math.Sin(angle)
		x2 := 800 + 120*math.Cos(angle)
		y2 := 150 + 120*math.Sin(angle)
		
		dc.SetLineWidth(4)
		dc.MoveTo(x1, y1)
		dc.LineTo(x2, y2)
		dc.Stroke()
	}
}

func drawWater(dc *advancegg.Context) {
	dc.SetRGB(0.2, 0.4, 0.8)
	
	// Water surface
	dc.DrawRectangle(0, 600, 1000, 200)
	dc.Fill()
	
	// Water ripples
	for i := 0; i < 20; i++ {
		x := float64(i * 50)
		y := 650 + 20*math.Sin(float64(i)*0.5)
		
		dc.SetRGBA(0.3, 0.5, 0.9, 0.5)
		dc.DrawCircle(x, y, 15)
		dc.Fill()
	}
}

func drawReflection(dc *advancegg.Context) {
	// Simplified reflection of mountains and sun
	dc.SetRGBA(0.5, 0.5, 0.5, 0.3)
	
	// Reflected mountains (inverted)
	for i := 0; i < 3; i++ {
		baseX := float64(i * 200)
		height := 100 + float64(i*30)
		
		dc.MoveTo(baseX, 600)
		dc.LineTo(baseX+100, 600+height)
		dc.LineTo(baseX+200, 600)
		dc.ClosePath()
		dc.Fill()
	}
	
	// Reflected sun
	dc.SetRGBA(1, 0.9, 0.3, 0.4)
	dc.DrawCircle(800, 650, 40)
	dc.Fill()
}
