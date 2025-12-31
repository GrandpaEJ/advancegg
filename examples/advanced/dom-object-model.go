package main

import (
	"fmt"
	"image/color"

	"github.com/GrandpaEJ/advancegg"
)

func main() {
	fmt.Println("Creating DOM-style Object Model demo...")
	
	// Create document-based graphics demo
	createDocumentDemo()
	
	// Create CSS-style styling demo
	createStylingDemo()
	
	// Create hierarchical elements demo
	createHierarchyDemo()
	
	fmt.Println("DOM Object Model demo completed!")
}

func createDocumentDemo() {
	// Create a new document
	doc := advancegg.NewDocument()
	
	// Create elements with IDs
	title := advancegg.CreateText("title", 200, 30, "DOM Object Model Demo")
	title.AddClass("heading")
	title.SetStyle("fill", color.RGBA{0, 0, 0, 255})
	
	// Create shapes with IDs and classes
	redCircle := advancegg.CreateCircle("circle1", 100, 100, 40)
	redCircle.AddClass("shape")
	redCircle.AddClass("circle")
	redCircle.SetStyle("fill", color.RGBA{255, 100, 100, 255})
	
	blueRect := advancegg.CreateRect("rect1", 200, 60, 80, 80)
	blueRect.AddClass("shape")
	blueRect.AddClass("rectangle")
	blueRect.SetStyle("fill", color.RGBA{100, 100, 255, 255})
	
	greenLine := advancegg.CreateLine("line1", 50, 200, 350, 200)
	greenLine.AddClass("shape")
	greenLine.AddClass("line")
	greenLine.SetStyle("stroke", color.RGBA{100, 255, 100, 255})
	greenLine.SetStyle("stroke-width", 3.0)
	
	// Add elements to document
	doc.AddElement(title)
	doc.AddElement(redCircle)
	doc.AddElement(blueRect)
	doc.AddElement(greenLine)
	
	// Create context and render
	dc := advancegg.NewContext(400, 300)
	dc.SetRGB(1, 1, 1)
	dc.Clear()
	
	// Render document
	doc.Render(dc)
	
	// Save result
	advancegg.SavePNG("images/advanced/dom-document-demo.png", dc.Image())
	fmt.Println("Created dom-document-demo.png")
	
	// Demonstrate element access
	fmt.Println("\nElement Access Demo:")
	
	// Access by ID
	circle := doc.GetElementByID("circle1")
	if circle != nil {
		fmt.Printf("Found circle with ID: %s\n", circle.ID)
		fmt.Printf("Circle has classes: %v\n", circle.Classes)
	}
	
	// Access by class
	shapes := doc.GetElementsByClass("shape")
	fmt.Printf("Found %d elements with class 'shape'\n", len(shapes))
	
	circles := doc.GetElementsByClass("circle")
	fmt.Printf("Found %d elements with class 'circle'\n", len(circles))
}

func createStylingDemo() {
	// Create document with CSS-like styling
	doc := advancegg.NewDocument()
	
	// Define styles
	headingStyle := advancegg.Style{
		Fill: color.RGBA{50, 50, 50, 255},
	}
	
	primaryStyle := advancegg.Style{
		Fill:        color.RGBA{255, 100, 100, 255},
		Stroke:      color.RGBA{200, 50, 50, 255},
		StrokeWidth: 2.0,
	}
	
	secondaryStyle := advancegg.Style{
		Fill:        color.RGBA{100, 255, 100, 255},
		Stroke:      color.RGBA{50, 200, 50, 255},
		StrokeWidth: 1.5,
	}
	
	// Add styles to document
	doc.AddStyle(".heading", headingStyle)
	doc.AddStyle(".primary", primaryStyle)
	doc.AddStyle(".secondary", secondaryStyle)
	doc.AddStyle("#special", advancegg.Style{
		Fill: color.RGBA{255, 215, 0, 255}, // Gold
	})
	
	// Create elements
	title := advancegg.CreateText("title", 150, 30, "CSS-Style Styling Demo")
	title.AddClass("heading")
	
	// Primary shapes
	circle1 := advancegg.CreateCircle("c1", 80, 100, 30)
	circle1.AddClass("primary")
	
	rect1 := advancegg.CreateRect("r1", 150, 70, 60, 60)
	rect1.AddClass("primary")
	
	// Secondary shapes
	circle2 := advancegg.CreateCircle("c2", 80, 180, 25)
	circle2.AddClass("secondary")
	
	rect2 := advancegg.CreateRect("r2", 150, 155, 60, 50)
	rect2.AddClass("secondary")
	
	// Special element
	specialCircle := advancegg.CreateCircle("special", 280, 140, 35)
	specialCircle.AddClass("primary") // Will be overridden by ID style
	
	// Add elements
	doc.AddElement(title)
	doc.AddElement(circle1)
	doc.AddElement(rect1)
	doc.AddElement(circle2)
	doc.AddElement(rect2)
	doc.AddElement(specialCircle)
	
	// Create context and render
	dc := advancegg.NewContext(400, 300)
	dc.SetRGB(0.95, 0.95, 0.95)
	dc.Clear()
	
	// Render with styles applied
	doc.Render(dc)
	
	// Save result
	advancegg.SavePNG("images/advanced/dom-styling-demo.png", dc.Image())
	fmt.Println("Created dom-styling-demo.png")
}

func createHierarchyDemo() {
	// Create document with hierarchical elements
	doc := advancegg.NewDocument()
	
	// Create a group element (container)
	group1 := advancegg.NewElement("group1")
	group1.AddClass("group")
	
	// Add shapes to the group
	circle := advancegg.CreateCircle("gc1", 100, 100, 30)
	circle.AddClass("grouped-shape")
	circle.SetStyle("fill", color.RGBA{255, 150, 150, 255})
	
	rect := advancegg.CreateRect("gr1", 150, 70, 60, 60)
	rect.AddClass("grouped-shape")
	rect.SetStyle("fill", color.RGBA{150, 150, 255, 255})
	
	// Add shapes to group
	group1.AddChild(circle)
	group1.AddChild(rect)
	
	// Create another group
	group2 := advancegg.NewElement("group2")
	group2.AddClass("group")
	
	line1 := advancegg.CreateLine("gl1", 50, 200, 150, 200)
	line1.SetStyle("stroke", color.RGBA{255, 100, 100, 255})
	line1.SetStyle("stroke-width", 3.0)
	
	line2 := advancegg.CreateLine("gl2", 150, 180, 250, 220)
	line2.SetStyle("stroke", color.RGBA{100, 255, 100, 255})
	line2.SetStyle("stroke-width", 3.0)
	
	group2.AddChild(line1)
	group2.AddChild(line2)
	
	// Add title
	title := advancegg.CreateText("title", 120, 30, "Hierarchical Elements")
	title.SetStyle("fill", color.RGBA{0, 0, 0, 255})
	
	// Add elements to document
	doc.AddElement(title)
	doc.AddElement(group1)
	doc.AddElement(group2)
	
	// Create context and render
	dc := advancegg.NewContext(400, 300)
	dc.SetRGB(1, 1, 1)
	dc.Clear()
	
	// Render document
	doc.Render(dc)
	
	// Save result
	advancegg.SavePNG("images/advanced/dom-hierarchy-demo.png", dc.Image())
	fmt.Println("Created dom-hierarchy-demo.png")
	
	// Demonstrate hierarchy navigation
	fmt.Println("\nHierarchy Navigation Demo:")
	
	// Find grouped shapes
	groupedShapes := doc.GetElementsByClass("grouped-shape")
	fmt.Printf("Found %d grouped shapes\n", len(groupedShapes))
	
	for _, shape := range groupedShapes {
		if shape.Parent != nil {
			fmt.Printf("Shape %s is child of %s\n", shape.ID, shape.Parent.ID)
		}
	}
	
	// Access group and its children
	group := doc.GetElementByID("group1")
	if group != nil {
		fmt.Printf("Group %s has %d children\n", group.ID, len(group.Children))
		for _, child := range group.Children {
			fmt.Printf("  Child: %s\n", child.ID)
		}
	}
}

// Demonstrate dynamic manipulation
func demonstrateDynamicManipulation() {
	fmt.Println("\nDynamic Manipulation capabilities:")
	fmt.Println("- Add/remove elements by ID")
	fmt.Println("- Modify element styles at runtime")
	fmt.Println("- Change element hierarchy")
	fmt.Println("- Apply CSS-like selectors")
	fmt.Println("- Clone and transform elements")
	fmt.Println("- Event-driven updates (with additional implementation)")
	
	// Example of dynamic changes
	doc := advancegg.NewDocument()
	
	// Create initial element
	circle := advancegg.CreateCircle("dynamic", 100, 100, 30)
	circle.SetStyle("fill", color.RGBA{255, 0, 0, 255})
	doc.AddElement(circle)
	
	// Modify element
	circle.SetStyle("fill", color.RGBA{0, 255, 0, 255}) // Change to green
	circle.AddClass("modified")
	
	// Add new class-based style
	doc.AddStyle(".modified", advancegg.Style{
		Stroke:      color.RGBA{0, 0, 255, 255},
		StrokeWidth: 3.0,
	})
	
	fmt.Println("Dynamic manipulation example completed")
}
