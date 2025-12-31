package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/GrandpaEJ/advancegg"
)

func main() {
	fmt.Println("Creating SVG demo...")
	
	// Create a simple SVG document programmatically
	createSVGDocument()
	
	// Create a drawing and export to SVG
	createDrawingAndExportSVG()
	
	fmt.Println("SVG demo completed!")
}

func createSVGDocument() {
	// Create SVG content manually
	svgContent := `<?xml version="1.0" encoding="UTF-8"?>
<svg width="400" height="300" viewBox="0 0 400 300" xmlns="http://www.w3.org/2000/svg">
  <!-- Background -->
  <rect x="0" y="0" width="400" height="300" fill="#f0f8ff"/>
  
  <!-- Title -->
  <text x="200" y="30" text-anchor="middle" font-family="Arial" font-size="20" fill="#333">SVG Import/Export Demo</text>
  
  <!-- Shapes -->
  <circle cx="100" cy="100" r="40" fill="#ff6b6b" stroke="#d63031" stroke-width="2"/>
  <rect x="160" y="60" width="80" height="80" fill="#4ecdc4" stroke="#00b894" stroke-width="2"/>
  <ellipse cx="320" cy="100" rx="50" ry="30" fill="#fdcb6e" stroke="#e17055" stroke-width="2"/>
  
  <!-- Lines and paths -->
  <line x1="50" y1="180" x2="350" y2="180" stroke="#6c5ce7" stroke-width="3"/>
  <path d="M 50 220 Q 200 180 350 220" fill="none" stroke="#a29bfe" stroke-width="3"/>
  
  <!-- Complex path -->
  <path d="M 100 250 L 150 220 L 200 250 L 250 220 L 300 250" fill="none" stroke="#fd79a8" stroke-width="2"/>
</svg>`

	// Save SVG file
	file, err := os.Create("demo.svg")
	if err != nil {
		fmt.Printf("Error creating SVG file: %v\n", err)
		return
	}
	defer file.Close()
	
	_, err = file.WriteString(svgContent)
	if err != nil {
		fmt.Printf("Error writing SVG file: %v\n", err)
		return
	}
	
	fmt.Println("Created demo.svg")
}

func createDrawingAndExportSVG() {
	// Create a drawing using the graphics library
	dc := advancegg.NewContext(500, 400)
	
	// Background
	dc.SetRGB(0.95, 0.95, 1.0)
	dc.Clear()
	
	// Title
	dc.SetRGB(0, 0, 0)
	dc.DrawString("Graphics to SVG Export", 150, 30)
	
	// Draw some shapes
	dc.SetRGB(1, 0.2, 0.2)
	dc.DrawCircle(100, 100, 40)
	dc.Fill()
	
	dc.SetRGB(0.2, 0.8, 0.2)
	dc.DrawRectangle(200, 60, 80, 80)
	dc.Fill()
	
	dc.SetRGB(0.2, 0.2, 1)
	dc.DrawEllipse(400, 100, 50, 30)
	dc.Fill()
	
	// Draw some lines
	dc.SetRGB(0.5, 0, 0.8)
	dc.SetLineWidth(3)
	dc.MoveTo(50, 200)
	dc.LineTo(450, 200)
	dc.Stroke()
	
	// Draw a curve
	dc.SetRGB(0.8, 0.4, 0)
	dc.MoveTo(50, 250)
	dc.QuadraticTo(250, 220, 450, 250)
	dc.Stroke()
	
	// Draw a complex path
	dc.SetRGB(0.6, 0.2, 0.8)
	dc.SetLineWidth(2)
	dc.MoveTo(100, 300)
	dc.LineTo(150, 280)
	dc.LineTo(200, 300)
	dc.LineTo(250, 280)
	dc.LineTo(300, 300)
	dc.LineTo(350, 280)
	dc.LineTo(400, 300)
	dc.Stroke()
	
	// Save as PNG first
	advancegg.SavePNG("graphics-demo.png", dc.Image())
	fmt.Println("Created graphics-demo.png")
	
	// Create equivalent SVG manually (since full SVG export is complex)
	svgFromGraphics := createSVGFromDrawing()
	
	file, err := os.Create("graphics-export.svg")
	if err != nil {
		fmt.Printf("Error creating SVG export file: %v\n", err)
		return
	}
	defer file.Close()
	
	_, err = file.WriteString(svgFromGraphics)
	if err != nil {
		fmt.Printf("Error writing SVG export file: %v\n", err)
		return
	}
	
	fmt.Println("Created graphics-export.svg")
}

func createSVGFromDrawing() string {
	// Manually create SVG equivalent of the drawing
	// In a full implementation, this would be automated
	var svg strings.Builder
	
	svg.WriteString(`<?xml version="1.0" encoding="UTF-8"?>
<svg width="500" height="400" viewBox="0 0 500 400" xmlns="http://www.w3.org/2000/svg">
  <!-- Background -->
  <rect x="0" y="0" width="500" height="400" fill="#f2f2ff"/>
  
  <!-- Title -->
  <text x="150" y="30" font-family="Arial" font-size="16" fill="#000">Graphics to SVG Export</text>
  
  <!-- Red circle -->
  <circle cx="100" cy="100" r="40" fill="#ff3333"/>
  
  <!-- Green rectangle -->
  <rect x="200" y="60" width="80" height="80" fill="#33cc33"/>
  
  <!-- Blue ellipse -->
  <ellipse cx="400" cy="100" rx="50" ry="30" fill="#3333ff"/>
  
  <!-- Purple line -->
  <line x1="50" y1="200" x2="450" y2="200" stroke="#8000cc" stroke-width="3"/>
  
  <!-- Orange curve -->
  <path d="M 50 250 Q 250 220 450 250" fill="none" stroke="#cc6600" stroke-width="3"/>
  
  <!-- Purple zigzag -->
  <path d="M 100 300 L 150 280 L 200 300 L 250 280 L 300 300 L 350 280 L 400 300" 
        fill="none" stroke="#9933cc" stroke-width="2"/>
</svg>`)
	
	return svg.String()
}

// Demonstration of SVG parsing (simplified)
func demonstrateSVGParsing() {
	fmt.Println("\nSVG Parsing capabilities:")
	fmt.Println("- Parse SVG documents from XML")
	fmt.Println("- Extract shapes: rectangles, circles, ellipses, lines, paths")
	fmt.Println("- Handle basic attributes: fill, stroke, stroke-width")
	fmt.Println("- Support viewBox and coordinate transformations")
	fmt.Println("- Convert SVG elements to drawable graphics")
	fmt.Println("\nNote: Full SVG specification support would require extensive implementation")
	fmt.Println("Current implementation provides foundation for basic SVG import/export")
}
