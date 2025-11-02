package main

import (
	"fmt"
	"image/color"

	"github.com/GrandpaEJ/advancegg"
)

func main() {
	fmt.Println("Creating basic operations examples...")

	// Demonstrate new image creation functions
	createNewImageExamples()

	// Demonstrate paste operations
	pasteImageExamples()

	// Demonstrate text box and alignment
	textBoxExamples()

	// Demonstrate enhanced geometric operations
	geometricExamples()

	fmt.Println("Basic operations examples completed!")
}

func createNewImageExamples() {
	// Create new images with different backgrounds
	rgbImg := advancegg.CreateNewRGB(400, 300, 255, 200, 150)
	rgbImg.SavePNG("images/new-rgb-image.png")
	fmt.Println("Created RGB image: new-rgb-image.png")

	rgbaImg := advancegg.CreateNewRGBA(400, 300, 100, 150, 200, 180)
	rgbaImg.SavePNG("images/new-rgba-image.png")
	fmt.Println("Created RGBA image: new-rgba-image.png")

	grayImg := advancegg.CreateNewGrayscale(400, 300, 128)
	grayImg.SavePNG("images/new-grayscale-image.png")
	fmt.Println("Created grayscale image: new-grayscale-image.png")

	transparentImg := advancegg.CreateNewTransparent(400, 300)
	transparentImg.SavePNG("images/new-transparent-image.png")
	fmt.Println("Created transparent image: new-transparent-image.png")
}

func pasteImageExamples() {
	// Create base image
	base := advancegg.CreateNewRGB(600, 400, 240, 240, 240)

	// Create smaller image to paste
	smallImg := advancegg.CreateNewRGB(150, 100, 255, 100, 100)
	smallImg.DrawString("PASTED", 20, 50)

	// Paste the small image onto the base
	base.PasteImage(smallImg.Image(), 50, 50)
	base.SavePNG("images/paste-basic.png")
	fmt.Println("Created basic paste example: paste-basic.png")

	// Create another example with mask
	base2 := advancegg.CreateNewRGB(600, 400, 200, 220, 240)

	// Create circular image with mask
	circleImg := advancegg.CreateNewRGB(120, 120, 100, 200, 100)
	circleImg.SetRGB(1, 0, 0)
	circleImg.DrawCircle(60, 60, 50)
	circleImg.Fill()

	// Create circular mask
	mask := advancegg.CreateNewGrayscale(120, 120, 0)
	mask.SetRGB(1, 1, 1)
	mask.DrawCircle(60, 60, 50)
	mask.Fill()

	base2.PasteImageWithMask(circleImg.Image(), mask.Image(), 200, 150)
	base2.SavePNG("images/paste-with-mask.png")
	fmt.Println("Created paste with mask example: paste-with-mask.png")
}

func textBoxExamples() {
	dc := advancegg.NewContext(800, 600)

	// White background
	dc.SetRGB(1, 1, 1)
	dc.Clear()

	// Load a font
	err := dc.LoadFontFace("fonts/arial.ttf", 24)
	if err != nil {
		fmt.Println("Font not found, using default")
	}

	// Example 1: Basic text box
	dc.SetRGB(0, 0, 0)
	text := "This is a text box example that demonstrates automatic word wrapping and text layout within a bounding box."
	dc.DrawTextBox(text, 50, 50, 300, 150, advancegg.AlignLeft)

	// Draw box outline
	dc.SetRGB(0.7, 0.7, 0.7)
	dc.SetLineWidth(1)
	dc.DrawRectangle(50, 50, 300, 150)
	dc.Stroke()

	// Example 2: Centered text box
	dc.SetRGB(0.2, 0.2, 0.8)
	dc.DrawTextBox("Centered Text\nMultiple Lines\nWord Wrapped", 400, 50, 300, 150, advancegg.AlignCenter)

	// Draw box outline
	dc.SetRGB(0.7, 0.7, 0.7)
	dc.DrawRectangle(400, 50, 300, 150)
	dc.Stroke()

	// Example 3: Right-aligned text box
	dc.SetRGB(0.8, 0.2, 0.2)
	dc.DrawTextBox("Right Aligned Text\nWith Multiple Lines\nAnd Word Wrapping", 50, 250, 300, 150, advancegg.AlignRight)

	// Draw box outline
	dc.SetRGB(0.7, 0.7, 0.7)
	dc.DrawRectangle(50, 250, 300, 150)
	dc.Stroke()

	// Example 4: Anchored text box
	dc.SetRGB(0.2, 0.8, 0.2)
	dc.DrawTextBoxAnchored("Anchored Text Box\nCenter positioned\nWith anchor point", 400, 250, 0.5, 0.5, 300, 150, advancegg.AlignCenter)

	// Draw box outline
	dc.SetRGB(0.7, 0.7, 0.7)
	dc.DrawRectangle(250, 175, 300, 150) // Adjusted for anchor
	dc.Stroke()

	dc.SavePNG("images/text-box-examples.png")
	fmt.Println("Created text box examples: text-box-examples.png")
}

func geometricExamples() {
	dc := advancegg.NewContext(1000, 800)

	// Light background
	dc.SetRGB(0.95, 0.95, 0.95)
	dc.Clear()

	// Load font for labels
	err := dc.LoadFontFace("fonts/arial.ttf", 16)
	if err != nil {
		fmt.Println("Font not found, using default")
	}

	// Example 1: Circles on existing image
	dc.SetRGB(0, 0, 0)
	dc.DrawString("Circle Operations:", 50, 30)

	// Filled circle
	dc.SetRGB(1, 0, 0)
	dc.DrawFilledCircle(100, 80, 40)

	// Stroked circle
	dc.SetRGB(0, 1, 0)
	dc.SetLineWidth(3)
	dc.DrawStrokedCircle(200, 80, 40)

	// Circle with border
	dc.DrawCircleWithBorder(300, 80, 40, color.RGBA{0, 0, 255, 255}, color.RGBA{0, 0, 0, 255}, 4)

	// Example 2: Ellipses
	dc.SetRGB(0, 0, 0)
	dc.DrawString("Ellipse Operations:", 50, 160)

	dc.SetRGB(1, 0.5, 0)
	dc.DrawEllipseOnImage(100, 200, 60, 40)

	dc.SetRGB(0, 0.5, 1)
	dc.DrawEllipseOnImage(200, 200, 40, 60)

	// Example 3: Rectangles
	dc.SetRGB(0, 0, 0)
	dc.DrawString("Rectangle Operations:", 50, 280)

	dc.SetRGB(1, 0, 1)
	dc.DrawRectangleOnImage(100, 320, 80, 50)

	dc.SetRGB(0, 1, 1)
	dc.DrawRoundedRectangleOnImage(200, 320, 80, 50, 10)

	// Example 4: Polygons and shapes
	dc.SetRGB(0, 0, 0)
	dc.DrawString("Polygon & Shape Operations:", 50, 400)

	// Regular polygon
	dc.SetRGB(1, 0.5, 0.5)
	dc.DrawPolygon(100, 450, 40, 6) // Hexagon

	// Star
	dc.SetRGB(0.5, 0.5, 1)
	dc.DrawStar(200, 450, 45, 20, 5) // 5-pointed star

	// Pie slice
	dc.SetRGB(0.5, 1, 0.5)
	dc.DrawPieSlice(300, 450, 40, 0, 3.14159/2) // Quarter circle

	// Donut
	dc.SetRGB(1, 0.5, 1)
	dc.DrawDonut(400, 450, 40, 20)

	// Example 5: Arcs
	dc.SetRGB(0, 0, 0)
	dc.DrawString("Arc Operations:", 50, 520)

	dc.SetRGB(1, 0.8, 0)
	dc.SetLineWidth(4)
	dc.DrawArc(100, 570, 40, 0, 3.14159) // Semicircle arc

	dc.SetRGB(0, 0.8, 1)
	dc.DrawArc(200, 570, 40, 3.14159/4, 3.14159*1.5) // 3/4 arc

	dc.SavePNG("images/geometric-operations.png")
	fmt.Println("Created geometric operations examples: geometric-operations.png")
}
