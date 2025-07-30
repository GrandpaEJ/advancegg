package main

import (
	"fmt"
	"image/color"

	"github.com/GrandpaEJ/advancegg"
)

func main() {
	fmt.Println("Creating advanced features examples...")

	// Layer system example
	createLayerExample()

	// Non-destructive editing example
	createNonDestructiveExample()

	// Smart guides example
	createGuidesExample()

	fmt.Println("Advanced features examples completed!")
}

func createLayerExample() {
	dc := advancegg.NewContext(800, 600)

	// Enable layer system
	dc.EnableLayers()
	lm := dc.GetLayerManager()

	// Create background layer with gradient
	bgLayer := lm.GetActiveLayer()
	bgLayer.Name = "Background"
	bgLayer.Fill(color.RGBA{50, 50, 100, 255})

	// Add a shapes layer
	dc.AddLayer("Shapes")
	dc.SetActiveLayerByName("Shapes")

	// Draw some shapes on the shapes layer
	dc.SetRGB(1, 0.5, 0)
	dc.DrawCircle(200, 150, 80)
	dc.Fill()

	dc.SetRGB(0, 0.8, 0.3)
	dc.DrawRectangle(400, 100, 150, 100)
	dc.Fill()

	// Add a text layer with transparency
	textLayer := dc.AddLayer("Text")
	textLayer.SetOpacity(0.8)
	textLayer.SetBlendMode(advancegg.BlendModeOverlay)
	dc.SetActiveLayerByName("Text")

	dc.SetRGB(1, 1, 1)
	dc.DrawString("Layer System Demo", 50, 50)
	dc.DrawString("Multiple layers with blend modes", 50, 80)

	// Add an effects layer
	effectsLayer := dc.AddLayer("Effects")
	effectsLayer.SetOpacity(0.6)
	effectsLayer.SetBlendMode(advancegg.BlendModeMultiply)
	dc.SetActiveLayerByName("Effects")

	// Draw some overlay effects
	dc.SetRGB(0.8, 0.2, 0.8)
	for i := 0; i < 10; i++ {
		x := float64(i * 80)
		y := 300 + float64(i*10)
		dc.DrawCircle(x, y, 30)
		dc.Fill()
	}

	// Composite all layers
	dc.CompositeToImage()

	dc.SavePNG("images/layer-system-demo.png")
	fmt.Println("Layer system demo saved as layer-system-demo.png")
}

func createNonDestructiveExample() {
	dc := advancegg.NewContext(600, 400)

	// Create base image
	dc.SetRGB(0.2, 0.4, 0.8)
	dc.Clear()

	// Draw some content
	dc.SetRGB(1, 1, 0)
	dc.DrawCircle(150, 150, 80)
	dc.Fill()

	dc.SetRGB(1, 0.5, 0)
	dc.DrawRectangle(300, 100, 120, 100)
	dc.Fill()

	dc.SetRGB(0, 1, 0.5)
	dc.DrawCircle(450, 250, 60)
	dc.Fill()

	// Enable non-destructive editing
	dc.EnableNonDestructiveEditing()

	// Add brightness adjustment
	brightnessOp := &advancegg.BrightnessOperation{Amount: 1.3}
	dc.AddEditOperation(brightnessOp)

	// Add contrast adjustment
	contrastOp := &advancegg.ContrastOperation{Amount: 1.2}
	dc.AddEditOperation(contrastOp)

	// Add blur effect
	blurOp := &advancegg.BlurOperation{Radius: 3}
	dc.AddEditOperation(blurOp)

	// Apply all non-destructive edits
	dc.ApplyNonDestructiveEdits()

	dc.SavePNG("images/non-destructive-demo.png")
	fmt.Println("Non-destructive editing demo saved as non-destructive-demo.png")

	// Show edit stack info
	editStack := dc.GetEditStack()
	fmt.Printf("Edit stack has %d operations\n", len(editStack.Operations))

	// Demonstrate editing an operation
	editStack.UpdateOperation(0, map[string]interface{}{"amount": 0.8}) // Reduce brightness
	dc.ApplyNonDestructiveEdits()

	dc.SavePNG("images/non-destructive-modified.png")
	fmt.Println("Modified non-destructive demo saved as non-destructive-modified.png")
}

func createGuidesExample() {
	dc := advancegg.NewContext(800, 600)

	// White background
	dc.SetRGB(1, 1, 1)
	dc.Clear()

	// Enable guides
	dc.EnableGuides()
	gm := dc.GetGuideManager()

	// Set up grid
	gm.GridSize = 40
	gm.GridVisible = true
	gm.SnapToGrid = true

	// Add some manual guides
	dc.AddGuide(200, advancegg.GuideVertical)
	dc.AddGuide(400, advancegg.GuideVertical)
	dc.AddGuide(150, advancegg.GuideHorizontal)
	dc.AddGuide(300, advancegg.GuideHorizontal)

	// Enable center guides
	gm.ShowCenterLines = true
	gm.GenerateCenterGuides(800, 600)

	// Enable margin guides
	gm.Margins = advancegg.Margins{Top: 50, Right: 50, Bottom: 50, Left: 50}
	gm.GenerateMarginGuides(800, 600)

	// Draw grid
	if gm.GridVisible {
		dc.SetRGBA(0.8, 0.8, 0.8, 0.3)
		dc.SetLineWidth(1)

		// Vertical grid lines
		for x := 0.0; x <= 800; x += gm.GridSize {
			dc.DrawLine(x, 0, x, 600)
			dc.Stroke()
		}

		// Horizontal grid lines
		for y := 0.0; y <= 600; y += gm.GridSize {
			dc.DrawLine(0, y, 800, y)
			dc.Stroke()
		}
	}

	// Draw guides
	for _, guide := range gm.Guides {
		if !guide.Visible {
			continue
		}

		r, g, b, a := guide.Color.RGBA()
		dc.SetRGBA(float64(r)/65535, float64(g)/65535, float64(b)/65535, float64(a)/65535)
		dc.SetLineWidth(2)

		switch guide.Orientation {
		case advancegg.GuideVertical:
			dc.DrawLine(guide.Position, 0, guide.Position, 600)
		case advancegg.GuideHorizontal:
			dc.DrawLine(0, guide.Position, 800, guide.Position)
		}
		dc.Stroke()
	}

	// Draw some shapes that snap to guides
	shapes := []struct {
		x, y, w, h float64
		color      [3]float64
	}{
		{100, 100, 80, 60, [3]float64{1, 0.3, 0.3}},
		{300, 200, 100, 80, [3]float64{0.3, 1, 0.3}},
		{500, 150, 120, 100, [3]float64{0.3, 0.3, 1}},
		{150, 350, 90, 70, [3]float64{1, 1, 0.3}},
	}

	for _, shape := range shapes {
		// Snap the shape position
		snappedX, snappedY := gm.SnapPoint(shape.x, shape.y)
		snappedX, snappedY, snappedW, snappedH := gm.SnapRectangle(snappedX, snappedY, shape.w, shape.h)

		dc.SetRGB(shape.color[0], shape.color[1], shape.color[2])
		dc.DrawRectangle(snappedX, snappedY, snappedW, snappedH)
		dc.Fill()

		// Draw original position with outline to show snapping
		dc.SetRGBA(0, 0, 0, 0.5)
		dc.SetLineWidth(1)
		dc.DrawRectangle(shape.x, shape.y, shape.w, shape.h)
		dc.Stroke()
	}

	// Demonstrate alignment
	targets := []advancegg.AlignmentTarget{
		{X: 50, Y: 450, Width: 60, Height: 40, ID: "rect1"},
		{X: 150, Y: 470, Width: 80, Height: 30, ID: "rect2"},
		{X: 280, Y: 460, Width: 70, Height: 50, ID: "rect3"},
	}

	// Draw original positions
	dc.SetRGBA(1, 0, 0, 0.3)
	for _, target := range targets {
		dc.DrawRectangle(target.X, target.Y, target.Width, target.Height)
		dc.Fill()
	}

	// Align to top
	alignedTargets := advancegg.AlignTargetsTop(targets)

	// Draw aligned positions
	dc.SetRGBA(0, 1, 0, 0.7)
	for _, target := range alignedTargets {
		dc.DrawRectangle(target.X, target.Y, target.Width, target.Height)
		dc.Fill()
	}

	// Add labels
	dc.SetRGB(0, 0, 0)
	dc.DrawString("Smart Guides & Alignment Demo", 50, 30)
	dc.DrawString("Grid: 40px, Guides: Blue lines, Center: Red lines, Margins: Green lines", 50, 550)
	dc.DrawString("Red shapes: Original positions, Green shapes: Aligned to top", 50, 570)

	dc.SavePNG("images/guides-alignment-demo.png")
	fmt.Println("Guides and alignment demo saved as guides-alignment-demo.png")
}

func createAdvancedCompositeExample() {
	dc := advancegg.NewContext(1000, 700)

	// Enable all advanced features
	dc.EnableLayers()
	dc.EnableNonDestructiveEditing()
	dc.EnableGuides()

	// Set up guides
	gm := dc.GetGuideManager()
	gm.GridSize = 50
	gm.SnapToGrid = true
	gm.ShowCenterLines = true
	gm.GenerateCenterGuides(1000, 700)

	// Create multiple layers with different content

	// Background layer
	bgLayer := dc.GetLayerManager().GetActiveLayer()
	bgLayer.Name = "Background"
	bgLayer.Fill(color.RGBA{20, 30, 50, 255})

	// Shapes layer
	dc.AddLayer("Geometric Shapes")
	dc.SetActiveLayerByName("Geometric Shapes")

	// Draw snapped geometric shapes
	shapes := []struct {
		x, y, size float64
		color      color.RGBA
		shape      string
	}{
		{200, 150, 80, color.RGBA{255, 100, 100, 200}, "circle"},
		{400, 200, 100, color.RGBA{100, 255, 100, 200}, "square"},
		{600, 150, 90, color.RGBA{100, 100, 255, 200}, "circle"},
		{300, 350, 120, color.RGBA{255, 255, 100, 200}, "square"},
	}

	for _, s := range shapes {
		snappedX, snappedY := gm.SnapPoint(s.x, s.y)
		dc.SetColor(s.color)

		if s.shape == "circle" {
			dc.DrawCircle(snappedX, snappedY, s.size/2)
		} else {
			dc.DrawRectangle(snappedX-s.size/2, snappedY-s.size/2, s.size, s.size)
		}
		dc.Fill()
	}

	// Text layer with effects
	textLayer := dc.AddLayer("Text Effects")
	textLayer.SetOpacity(0.9)
	textLayer.SetBlendMode(advancegg.BlendModeOverlay)
	dc.SetActiveLayerByName("Text Effects")

	dc.SetRGB(1, 1, 1)
	dc.DrawString("Advanced AdvanceGG Features", 50, 50)
	dc.DrawString("Layers + Non-destructive + Guides", 50, 80)

	// Apply non-destructive effects to the entire composition
	brightnessOp := &advancegg.BrightnessOperation{Amount: 1.1}
	dc.AddEditOperation(brightnessOp)

	contrastOp := &advancegg.ContrastOperation{Amount: 1.15}
	dc.AddEditOperation(contrastOp)

	// Composite and apply effects
	dc.CompositeToImage()
	dc.ApplyNonDestructiveEdits()

	dc.SavePNG("images/advanced-composite-demo.png")
	fmt.Println("Advanced composite demo saved as advanced-composite-demo.png")
}
