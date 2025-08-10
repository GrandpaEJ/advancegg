//go:build svg
// +build svg

package advance

import (
	"encoding/xml"
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"
)

// SVG support for vector graphics import/export

// SVGDocument represents a complete SVG document
type SVGDocument struct {
	XMLName  xml.Name     `xml:"svg"`
	Width    string       `xml:"width,attr,omitempty"`
	Height   string       `xml:"height,attr,omitempty"`
	ViewBox  string       `xml:"viewBox,attr,omitempty"`
	Xmlns    string       `xml:"xmlns,attr,omitempty"`
	Elements []SVGElement `xml:",any"`
}

// SVGElement represents any SVG element
type SVGElement interface {
	Draw(ctx *Context)
	GetBounds() (float64, float64, float64, float64)
	Transform(matrix Matrix)
}

// SVGPath represents an SVG path element
type SVGPath struct {
	XMLName     xml.Name `xml:"path"`
	D           string   `xml:"d,attr"`
	Fill        string   `xml:"fill,attr,omitempty"`
	Stroke      string   `xml:"stroke,attr,omitempty"`
	StrokeWidth string   `xml:"stroke-width,attr,omitempty"`
	Transform   string   `xml:"transform,attr,omitempty"`
	path        *Path
}

// SVGRect represents an SVG rectangle element
type SVGRect struct {
	XMLName   xml.Name `xml:"rect"`
	X         string   `xml:"x,attr,omitempty"`
	Y         string   `xml:"y,attr,omitempty"`
	Width     string   `xml:"width,attr,omitempty"`
	Height    string   `xml:"height,attr,omitempty"`
	Fill      string   `xml:"fill,attr,omitempty"`
	Stroke    string   `xml:"stroke,attr,omitempty"`
	Transform string   `xml:"transform,attr,omitempty"`
}

// SVGCircle represents an SVG circle element
type SVGCircle struct {
	XMLName   xml.Name `xml:"circle"`
	CX        string   `xml:"cx,attr,omitempty"`
	CY        string   `xml:"cy,attr,omitempty"`
	R         string   `xml:"r,attr,omitempty"`
	Fill      string   `xml:"fill,attr,omitempty"`
	Stroke    string   `xml:"stroke,attr,omitempty"`
	Transform string   `xml:"transform,attr,omitempty"`
}

// SVGEllipse represents an SVG ellipse element
type SVGEllipse struct {
	XMLName   xml.Name `xml:"ellipse"`
	CX        string   `xml:"cx,attr,omitempty"`
	CY        string   `xml:"cy,attr,omitempty"`
	RX        string   `xml:"rx,attr,omitempty"`
	RY        string   `xml:"ry,attr,omitempty"`
	Fill      string   `xml:"fill,attr,omitempty"`
	Stroke    string   `xml:"stroke,attr,omitempty"`
	Transform string   `xml:"transform,attr,omitempty"`
}

// SVGLine represents an SVG line element
type SVGLine struct {
	XMLName     xml.Name `xml:"line"`
	X1          string   `xml:"x1,attr,omitempty"`
	Y1          string   `xml:"y1,attr,omitempty"`
	X2          string   `xml:"x2,attr,omitempty"`
	Y2          string   `xml:"y2,attr,omitempty"`
	Stroke      string   `xml:"stroke,attr,omitempty"`
	StrokeWidth string   `xml:"stroke-width,attr,omitempty"`
	Transform   string   `xml:"transform,attr,omitempty"`
}

// SVGGroup represents an SVG group element
type SVGGroup struct {
	XMLName   xml.Name     `xml:"g"`
	Transform string       `xml:"transform,attr,omitempty"`
	Elements  []SVGElement `xml:",any"`
}

// NewSVGDocument creates a new SVG document
func NewSVGDocument(width, height float64) *SVGDocument {
	return &SVGDocument{
		Width:    fmt.Sprintf("%.2f", width),
		Height:   fmt.Sprintf("%.2f", height),
		ViewBox:  fmt.Sprintf("0 0 %.2f %.2f", width, height),
		Xmlns:    "http://www.w3.org/2000/svg",
		Elements: make([]SVGElement, 0),
	}
}

// AddElement adds an element to the SVG document
func (doc *SVGDocument) AddElement(element SVGElement) {
	doc.Elements = append(doc.Elements, element)
}

// WriteTo writes the SVG document to a writer
func (doc *SVGDocument) WriteTo(w io.Writer) error {
	encoder := xml.NewEncoder(w)
	encoder.Indent("", "  ")

	// Write XML header
	if _, err := w.Write([]byte(`<?xml version="1.0" encoding="UTF-8"?>` + "\n")); err != nil {
		return err
	}

	return encoder.Encode(doc)
}

// ParseSVG parses an SVG document from a reader
func ParseSVG(r io.Reader) (*SVGDocument, error) {
	var doc SVGDocument
	decoder := xml.NewDecoder(r)

	err := decoder.Decode(&doc)
	if err != nil {
		return nil, err
	}

	return &doc, nil
}

// Draw methods for SVG elements

// Draw renders the SVG path
func (p *SVGPath) Draw(ctx *Context) {
	if p.path == nil {
		p.path = parseSVGPath(p.D)
	}

	// Apply transform if present
	if p.Transform != "" {
		matrix := parseSVGTransform(p.Transform)
		p.path.Transform(matrix)
	}

	// Set fill color
	if p.Fill != "" && p.Fill != "none" {
		color := parseSVGColor(p.Fill)
		ctx.SetColor(color)
		ctx.DrawPath(p.path)
		ctx.Fill()
	}

	// Set stroke
	if p.Stroke != "" && p.Stroke != "none" {
		color := parseSVGColor(p.Stroke)
		ctx.SetColor(color)
		if p.StrokeWidth != "" {
			if width, err := strconv.ParseFloat(p.StrokeWidth, 64); err == nil {
				ctx.SetLineWidth(width)
			}
		}
		ctx.DrawPath(p.path)
		ctx.Stroke()
	}
}

// Draw renders the SVG rectangle
func (r *SVGRect) Draw(ctx *Context) {
	x := parseFloat(r.X)
	y := parseFloat(r.Y)
	width := parseFloat(r.Width)
	height := parseFloat(r.Height)

	// Apply transform if present
	if r.Transform != "" {
		matrix := parseSVGTransform(r.Transform)
		x, y = matrix.Transform(x, y)
	}

	// Set fill color
	if r.Fill != "" && r.Fill != "none" {
		color := parseSVGColor(r.Fill)
		ctx.SetColor(color)
		ctx.DrawRectangle(x, y, width, height)
	}

	// Set stroke
	if r.Stroke != "" && r.Stroke != "none" {
		color := parseSVGColor(r.Stroke)
		ctx.SetColor(color)
		ctx.DrawRectangle(x, y, width, height)
		ctx.Stroke()
	}
}

// Draw renders the SVG circle
func (c *SVGCircle) Draw(ctx *Context) {
	cx := parseFloat(c.CX)
	cy := parseFloat(c.CY)
	r := parseFloat(c.R)

	// Apply transform if present
	if c.Transform != "" {
		matrix := parseSVGTransform(c.Transform)
		cx, cy = matrix.Transform(cx, cy)
	}

	// Set fill color
	if c.Fill != "" && c.Fill != "none" {
		color := parseSVGColor(c.Fill)
		ctx.SetColor(color)
		ctx.DrawCircle(cx, cy, r)
	}

	// Set stroke
	if c.Stroke != "" && c.Stroke != "none" {
		color := parseSVGColor(c.Stroke)
		ctx.SetColor(color)
		ctx.DrawCircle(cx, cy, r)
		ctx.Stroke()
	}
}

// Draw renders the SVG ellipse
func (e *SVGEllipse) Draw(ctx *Context) {
	cx := parseFloat(e.CX)
	cy := parseFloat(e.CY)
	rx := parseFloat(e.RX)
	ry := parseFloat(e.RY)

	// Apply transform if present
	if e.Transform != "" {
		matrix := parseSVGTransform(e.Transform)
		cx, cy = matrix.Transform(cx, cy)
	}

	// Set fill color
	if e.Fill != "" && e.Fill != "none" {
		color := parseSVGColor(e.Fill)
		ctx.SetColor(color)
		ctx.DrawEllipse(cx, cy, rx, ry)
	}

	// Set stroke
	if e.Stroke != "" && e.Stroke != "none" {
		color := parseSVGColor(e.Stroke)
		ctx.SetColor(color)
		ctx.DrawEllipse(cx, cy, rx, ry)
		ctx.Stroke()
	}
}

// Draw renders the SVG line
func (l *SVGLine) Draw(ctx *Context) {
	x1 := parseFloat(l.X1)
	y1 := parseFloat(l.Y1)
	x2 := parseFloat(l.X2)
	y2 := parseFloat(l.Y2)

	// Apply transform if present
	if l.Transform != "" {
		matrix := parseSVGTransform(l.Transform)
		x1, y1 = matrix.Transform(x1, y1)
		x2, y2 = matrix.Transform(x2, y2)
	}

	// Set stroke
	if l.Stroke != "" && l.Stroke != "none" {
		color := parseSVGColor(l.Stroke)
		ctx.SetColor(color)
		if l.StrokeWidth != "" {
			if width, err := strconv.ParseFloat(l.StrokeWidth, 64); err == nil {
				ctx.SetLineWidth(width)
			}
		}
		ctx.DrawLine(x1, y1, x2, y2)
	}
}

// Draw renders the SVG group
func (g *SVGGroup) Draw(ctx *Context) {
	// Apply transform if present
	if g.Transform != "" {
		matrix := parseSVGTransform(g.Transform)
		// Apply transform to context
		ctx.transform = ctx.transform.Multiply(matrix)
	}

	// Draw all elements in the group
	for _, element := range g.Elements {
		element.Draw(ctx)
	}
}

// GetBounds methods
func (p *SVGPath) GetBounds() (float64, float64, float64, float64) {
	if p.path == nil {
		p.path = parseSVGPath(p.D)
	}
	return p.path.GetBounds()
}

func (r *SVGRect) GetBounds() (float64, float64, float64, float64) {
	x := parseFloat(r.X)
	y := parseFloat(r.Y)
	width := parseFloat(r.Width)
	height := parseFloat(r.Height)
	return x, y, x + width, y + height
}

func (c *SVGCircle) GetBounds() (float64, float64, float64, float64) {
	cx := parseFloat(c.CX)
	cy := parseFloat(c.CY)
	r := parseFloat(c.R)
	return cx - r, cy - r, cx + r, cy + r
}

func (e *SVGEllipse) GetBounds() (float64, float64, float64, float64) {
	cx := parseFloat(e.CX)
	cy := parseFloat(e.CY)
	rx := parseFloat(e.RX)
	ry := parseFloat(e.RY)
	return cx - rx, cy - ry, cx + rx, cy + ry
}

func (l *SVGLine) GetBounds() (float64, float64, float64, float64) {
	x1 := parseFloat(l.X1)
	y1 := parseFloat(l.Y1)
	x2 := parseFloat(l.X2)
	y2 := parseFloat(l.Y2)
	return math.Min(x1, x2), math.Min(y1, y2), math.Max(x1, x2), math.Max(y1, y2)
}

func (g *SVGGroup) GetBounds() (float64, float64, float64, float64) {
	if len(g.Elements) == 0 {
		return 0, 0, 0, 0
	}

	minX, minY, maxX, maxY := g.Elements[0].GetBounds()
	for i := 1; i < len(g.Elements); i++ {
		x1, y1, x2, y2 := g.Elements[i].GetBounds()
		if x1 < minX {
			minX = x1
		}
		if y1 < minY {
			minY = y1
		}
		if x2 > maxX {
			maxX = x2
		}
		if y2 > maxY {
			maxY = y2
		}
	}

	return minX, minY, maxX, maxY
}

// Transform methods
func (p *SVGPath) Transform(matrix Matrix) {
	if p.path == nil {
		p.path = parseSVGPath(p.D)
	}
	p.path.Transform(matrix)
}

func (r *SVGRect) Transform(matrix Matrix) {
	// Transform would update the rect coordinates
}

func (c *SVGCircle) Transform(matrix Matrix) {
	// Transform would update the circle coordinates
}

func (e *SVGEllipse) Transform(matrix Matrix) {
	// Transform would update the ellipse coordinates
}

func (l *SVGLine) Transform(matrix Matrix) {
	// Transform would update the line coordinates
}

func (g *SVGGroup) Transform(matrix Matrix) {
	for _, element := range g.Elements {
		element.Transform(matrix)
	}
}

// Utility functions for SVG parsing

// parseSVGPath parses an SVG path data string
func parseSVGPath(d string) *Path {
	path := NewPath()
	// Simplified path parsing - in a real implementation, this would be much more comprehensive
	commands := strings.Fields(d)

	for i := 0; i < len(commands); i++ {
		cmd := commands[i]
		switch cmd {
		case "M":
			if i+2 < len(commands) {
				x := parseFloat(commands[i+1])
				y := parseFloat(commands[i+2])
				path.MoveTo(x, y)
				i += 2
			}
		case "L":
			if i+2 < len(commands) {
				x := parseFloat(commands[i+1])
				y := parseFloat(commands[i+2])
				path.LineTo(x, y)
				i += 2
			}
		case "C":
			if i+6 < len(commands) {
				cp1x := parseFloat(commands[i+1])
				cp1y := parseFloat(commands[i+2])
				cp2x := parseFloat(commands[i+3])
				cp2y := parseFloat(commands[i+4])
				x := parseFloat(commands[i+5])
				y := parseFloat(commands[i+6])
				path.CurveTo(cp1x, cp1y, cp2x, cp2y, x, y)
				i += 6
			}
		case "Z":
			path.ClosePath()
		}
	}

	return path
}

// parseSVGTransform parses an SVG transform attribute
func parseSVGTransform(transform string) Matrix {
	// Simplified transform parsing - in a real implementation, this would handle all transform types
	matrix := NewIdentityMatrix()

	if strings.Contains(transform, "translate") {
		// Parse translate transform
		start := strings.Index(transform, "(")
		end := strings.Index(transform, ")")
		if start != -1 && end != -1 {
			values := strings.Split(transform[start+1:end], ",")
			if len(values) >= 2 {
				tx := parseFloat(values[0])
				ty := parseFloat(values[1])
				matrix = Matrix{1, 0, 0, 1, tx, ty}
			}
		}
	}

	return matrix
}

// parseSVGColor parses an SVG color value
func parseSVGColor(color string) *RGBColor {
	// Simplified color parsing
	if strings.HasPrefix(color, "#") {
		// Hex color
		if len(color) == 7 {
			r, _ := strconv.ParseInt(color[1:3], 16, 0)
			g, _ := strconv.ParseInt(color[3:5], 16, 0)
			b, _ := strconv.ParseInt(color[5:7], 16, 0)
			return &RGBColor{uint8(r), uint8(g), uint8(b), 255}
		}
	}

	// Named colors
	switch color {
	case "red":
		return &RGBColor{255, 0, 0, 255}
	case "green":
		return &RGBColor{0, 255, 0, 255}
	case "blue":
		return &RGBColor{0, 0, 255, 255}
	case "black":
		return &RGBColor{0, 0, 0, 255}
	case "white":
		return &RGBColor{255, 255, 255, 255}
	default:
		return &RGBColor{0, 0, 0, 255}
	}
}

// parseFloat parses a string to float64
func parseFloat(s string) float64 {
	if s == "" {
		return 0
	}
	f, _ := strconv.ParseFloat(s, 64)
	return f
}

// RGBColor represents an RGB color
type RGBColor struct {
	R, G, B, A uint8
}

// RGBA implements the color.Color interface
func (c *RGBColor) RGBA() (r, g, b, a uint32) {
	r = uint32(c.R)
	r |= r << 8
	g = uint32(c.G)
	g |= g << 8
	b = uint32(c.B)
	b |= b << 8
	a = uint32(c.A)
	a |= a << 8
	return
}

// ExportToSVG exports the current context to SVG
func (c *Context) ExportToSVG() *SVGDocument {
	doc := NewSVGDocument(float64(c.width), float64(c.height))

	// Convert context content to SVG elements
	// This is a simplified implementation - a real version would be much more comprehensive

	return doc
}
