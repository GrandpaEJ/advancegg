package core

import (
	"image/color"
	"strings"
)

// DOM-style object model for shapes with IDs and styles

// Element represents a drawable element with ID and styles
type Element struct {
	ID       string
	Classes  []string
	Styles   map[string]interface{}
	Children []*Element
	Parent   *Element
	Shape    Shape
}

// Shape interface for different drawable shapes
type Shape interface {
	Draw(ctx *Context)
	GetBounds() (float64, float64, float64, float64)
	Clone() Shape
}

// Document represents a collection of elements
type Document struct {
	Root     *Element
	elements map[string]*Element // ID -> Element mapping
	styles   map[string]Style    // CSS-like styles
}

// Style represents CSS-like styling
type Style struct {
	Fill        color.Color
	Stroke      color.Color
	StrokeWidth float64
	Opacity     float64
	Transform   Matrix
	Visible     bool
}

// Shape implementations

// RectShape represents a rectangle
type RectShape struct {
	X, Y, Width, Height float64
}

// CircleShape represents a circle
type CircleShape struct {
	X, Y, Radius float64
}

// LineShape represents a line
type LineShape struct {
	X1, Y1, X2, Y2 float64
}

// PathShape represents a complex path (simplified)
type PathShape struct {
	Points []Point
}

// TextShape represents text
type TextShape struct {
	X, Y float64
	Text string
}

// NewDocument creates a new document
func NewDocument() *Document {
	return &Document{
		Root:     NewElement("root"),
		elements: make(map[string]*Element),
		styles:   make(map[string]Style),
	}
}

// NewElement creates a new element
func NewElement(id string) *Element {
	return &Element{
		ID:       id,
		Classes:  make([]string, 0),
		Styles:   make(map[string]interface{}),
		Children: make([]*Element, 0),
	}
}

// Element methods

// AddChild adds a child element
func (e *Element) AddChild(child *Element) {
	child.Parent = e
	e.Children = append(e.Children, child)
}

// RemoveChild removes a child element
func (e *Element) RemoveChild(child *Element) {
	for i, c := range e.Children {
		if c == child {
			e.Children = append(e.Children[:i], e.Children[i+1:]...)
			child.Parent = nil
			break
		}
	}
}

// AddClass adds a CSS class
func (e *Element) AddClass(class string) {
	for _, c := range e.Classes {
		if c == class {
			return // Already has class
		}
	}
	e.Classes = append(e.Classes, class)
}

// RemoveClass removes a CSS class
func (e *Element) RemoveClass(class string) {
	for i, c := range e.Classes {
		if c == class {
			e.Classes = append(e.Classes[:i], e.Classes[i+1:]...)
			break
		}
	}
}

// HasClass checks if element has a CSS class
func (e *Element) HasClass(class string) bool {
	for _, c := range e.Classes {
		if c == class {
			return true
		}
	}
	return false
}

// SetStyle sets a style property
func (e *Element) SetStyle(property string, value interface{}) {
	e.Styles[property] = value
}

// GetStyle gets a style property
func (e *Element) GetStyle(property string) interface{} {
	return e.Styles[property]
}

// Draw renders the element and its children
func (e *Element) Draw(ctx *Context) {
	if e.Shape != nil {
		// Apply styles before drawing
		e.applyStyles(ctx)
		e.Shape.Draw(ctx)
	}

	// Draw children
	for _, child := range e.Children {
		child.Draw(ctx)
	}
}

// applyStyles applies element styles to context
func (e *Element) applyStyles(ctx *Context) {
	// Apply fill color
	if fill, ok := e.Styles["fill"].(color.Color); ok {
		ctx.SetColor(fill)
	}

	// Apply stroke
	if stroke, ok := e.Styles["stroke"].(color.Color); ok {
		ctx.SetColor(stroke)
	}

	// Apply stroke width
	if width, ok := e.Styles["stroke-width"].(float64); ok {
		ctx.SetLineWidth(width)
	}

	// Apply opacity
	if opacity, ok := e.Styles["opacity"].(float64); ok {
		// TODO: Implement opacity support
		_ = opacity
	}

	// Apply transform
	if transform, ok := e.Styles["transform"].(Matrix); ok {
		ctx.Push()
		// Apply matrix transformation (simplified)
		ctx.Scale(transform.XX, transform.YY)
		ctx.Translate(transform.X0, transform.Y0)
	}
}

// Document methods

// AddElement adds an element to the document
func (d *Document) AddElement(element *Element) {
	if element.ID != "" {
		d.elements[element.ID] = element
	}
	d.Root.AddChild(element)
}

// GetElementByID gets an element by ID
func (d *Document) GetElementByID(id string) *Element {
	return d.elements[id]
}

// GetElementsByClass gets elements by class name
func (d *Document) GetElementsByClass(className string) []*Element {
	var elements []*Element
	d.walkElements(d.Root, func(e *Element) {
		if e.HasClass(className) {
			elements = append(elements, e)
		}
	})
	return elements
}

// AddStyle adds a CSS-like style rule
func (d *Document) AddStyle(selector string, style Style) {
	d.styles[selector] = style
}

// ApplyStyles applies CSS-like styles to elements
func (d *Document) ApplyStyles() {
	for selector, style := range d.styles {
		elements := d.selectElements(selector)
		for _, element := range elements {
			d.applyStyleToElement(element, style)
		}
	}
}

// selectElements selects elements based on CSS-like selector
func (d *Document) selectElements(selector string) []*Element {
	var elements []*Element

	if strings.HasPrefix(selector, "#") {
		// ID selector
		id := selector[1:]
		if element := d.GetElementByID(id); element != nil {
			elements = append(elements, element)
		}
	} else if strings.HasPrefix(selector, ".") {
		// Class selector
		className := selector[1:]
		elements = d.GetElementsByClass(className)
	} else {
		// Type selector (simplified - would need shape type info)
		// For now, just return all elements
		d.walkElements(d.Root, func(e *Element) {
			elements = append(elements, e)
		})
	}

	return elements
}

// applyStyleToElement applies a style to an element
func (d *Document) applyStyleToElement(element *Element, style Style) {
	if style.Fill != nil {
		element.SetStyle("fill", style.Fill)
	}
	if style.Stroke != nil {
		element.SetStyle("stroke", style.Stroke)
	}
	if style.StrokeWidth > 0 {
		element.SetStyle("stroke-width", style.StrokeWidth)
	}
	if style.Opacity > 0 {
		element.SetStyle("opacity", style.Opacity)
	}
}

// walkElements walks through all elements in the tree
func (d *Document) walkElements(element *Element, fn func(*Element)) {
	fn(element)
	for _, child := range element.Children {
		d.walkElements(child, fn)
	}
}

// Render renders the entire document
func (d *Document) Render(ctx *Context) {
	// Apply styles first
	d.ApplyStyles()

	// Render from root
	d.Root.Draw(ctx)
}

// Shape implementations

// RectShape methods
func (r *RectShape) Draw(ctx *Context) {
	ctx.DrawRectangle(r.X, r.Y, r.Width, r.Height)
	ctx.Fill()
}

func (r *RectShape) GetBounds() (float64, float64, float64, float64) {
	return r.X, r.Y, r.X + r.Width, r.Y + r.Height
}

func (r *RectShape) Clone() Shape {
	return &RectShape{X: r.X, Y: r.Y, Width: r.Width, Height: r.Height}
}

// CircleShape methods
func (c *CircleShape) Draw(ctx *Context) {
	ctx.DrawCircle(c.X, c.Y, c.Radius)
	ctx.Fill()
}

func (c *CircleShape) GetBounds() (float64, float64, float64, float64) {
	return c.X - c.Radius, c.Y - c.Radius, c.X + c.Radius, c.Y + c.Radius
}

func (c *CircleShape) Clone() Shape {
	return &CircleShape{X: c.X, Y: c.Y, Radius: c.Radius}
}

// LineShape methods
func (l *LineShape) Draw(ctx *Context) {
	ctx.MoveTo(l.X1, l.Y1)
	ctx.LineTo(l.X2, l.Y2)
	ctx.Stroke()
}

func (l *LineShape) GetBounds() (float64, float64, float64, float64) {
	minX := l.X1
	if l.X2 < minX {
		minX = l.X2
	}
	maxX := l.X1
	if l.X2 > maxX {
		maxX = l.X2
	}
	minY := l.Y1
	if l.Y2 < minY {
		minY = l.Y2
	}
	maxY := l.Y1
	if l.Y2 > maxY {
		maxY = l.Y2
	}
	return minX, minY, maxX, maxY
}

func (l *LineShape) Clone() Shape {
	return &LineShape{X1: l.X1, Y1: l.Y1, X2: l.X2, Y2: l.Y2}
}

// PathShape methods
func (p *PathShape) Draw(ctx *Context) {
	if len(p.Points) > 0 {
		// Draw path as connected lines
		ctx.MoveTo(p.Points[0].X, p.Points[0].Y)
		for i := 1; i < len(p.Points); i++ {
			ctx.LineTo(p.Points[i].X, p.Points[i].Y)
		}
		ctx.Fill()
	}
}

func (p *PathShape) GetBounds() (float64, float64, float64, float64) {
	if len(p.Points) == 0 {
		return 0, 0, 0, 0
	}

	minX, minY := p.Points[0].X, p.Points[0].Y
	maxX, maxY := p.Points[0].X, p.Points[0].Y

	for _, point := range p.Points {
		if point.X < minX {
			minX = point.X
		}
		if point.X > maxX {
			maxX = point.X
		}
		if point.Y < minY {
			minY = point.Y
		}
		if point.Y > maxY {
			maxY = point.Y
		}
	}

	return minX, minY, maxX, maxY
}

func (p *PathShape) Clone() Shape {
	newPoints := make([]Point, len(p.Points))
	copy(newPoints, p.Points)
	return &PathShape{Points: newPoints}
}

// TextShape methods
func (t *TextShape) Draw(ctx *Context) {
	ctx.DrawString(t.Text, t.X, t.Y)
}

func (t *TextShape) GetBounds() (float64, float64, float64, float64) {
	// Simplified bounds - would need font metrics for accurate calculation
	width := float64(len(t.Text)) * 8 // Approximate character width
	height := 16.0                    // Approximate character height
	return t.X, t.Y - height, t.X + width, t.Y
}

func (t *TextShape) Clone() Shape {
	return &TextShape{X: t.X, Y: t.Y, Text: t.Text}
}

// Convenience functions for creating elements

// CreateRect creates a rectangle element
func CreateRect(id string, x, y, width, height float64) *Element {
	element := NewElement(id)
	element.Shape = &RectShape{X: x, Y: y, Width: width, Height: height}
	return element
}

// CreateCircle creates a circle element
func CreateCircle(id string, x, y, radius float64) *Element {
	element := NewElement(id)
	element.Shape = &CircleShape{X: x, Y: y, Radius: radius}
	return element
}

// CreateLine creates a line element
func CreateLine(id string, x1, y1, x2, y2 float64) *Element {
	element := NewElement(id)
	element.Shape = &LineShape{X1: x1, Y1: y1, X2: x2, Y2: y2}
	return element
}

// CreateText creates a text element
func CreateText(id string, x, y float64, text string) *Element {
	element := NewElement(id)
	element.Shape = &TextShape{X: x, Y: y, Text: text}
	return element
}
