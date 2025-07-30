package core

import (
	"image"
	"image/color"
)

// Batch operations for optimizing multiple draw calls

// BatchOperation represents a single drawing operation
type BatchOperation interface {
	Execute(ctx *Context)
	GetBounds() (x, y, width, height float64)
	GetType() BatchOpType
}

// BatchOpType represents the type of batch operation
type BatchOpType int

const (
	BatchOpCircle BatchOpType = iota
	BatchOpRectangle
	BatchOpLine
	BatchOpText
	BatchOpImage
	BatchOpPath
)

// Batch manages a collection of drawing operations
type Batch struct {
	operations []BatchOperation
	bounds     Rectangle
	dirty      bool
}

// Rectangle represents a bounding rectangle
type Rectangle struct {
	X, Y, Width, Height float64
}

// NewBatch creates a new batch
func NewBatch() *Batch {
	return &Batch{
		operations: make([]BatchOperation, 0),
		dirty:      true,
	}
}

// Add adds an operation to the batch
func (b *Batch) Add(op BatchOperation) {
	b.operations = append(b.operations, op)
	b.dirty = true
}

// Execute executes all operations in the batch
func (b *Batch) Execute(ctx *Context) {
	// Sort operations by type for better performance
	b.sortOperations()
	
	// Execute operations in batches by type
	b.executeBatched(ctx)
}

// sortOperations sorts operations by type for optimal batching
func (b *Batch) sortOperations() {
	if !b.dirty {
		return
	}
	
	// Simple sort by operation type
	for i := 0; i < len(b.operations)-1; i++ {
		for j := i + 1; j < len(b.operations); j++ {
			if b.operations[i].GetType() > b.operations[j].GetType() {
				b.operations[i], b.operations[j] = b.operations[j], b.operations[i]
			}
		}
	}
	
	b.dirty = false
}

// executeBatched executes operations in optimized batches
func (b *Batch) executeBatched(ctx *Context) {
	currentType := BatchOpType(-1)
	
	for _, op := range b.operations {
		opType := op.GetType()
		
		// If operation type changed, optimize for the new type
		if opType != currentType {
			b.optimizeForType(ctx, opType)
			currentType = opType
		}
		
		op.Execute(ctx)
	}
}

// optimizeForType optimizes context settings for a specific operation type
func (b *Batch) optimizeForType(ctx *Context, opType BatchOpType) {
	switch opType {
	case BatchOpCircle, BatchOpRectangle:
		// Optimize for shape drawing
		ctx.SetLineWidth(1)
	case BatchOpLine:
		// Optimize for line drawing
		ctx.SetLineCap(LineCapRound)
	case BatchOpText:
		// Optimize for text rendering
		// Font settings would be optimized here
	case BatchOpImage:
		// Optimize for image drawing
		// Image blending settings would be optimized here
	}
}

// Clear clears all operations from the batch
func (b *Batch) Clear() {
	b.operations = b.operations[:0]
	b.dirty = true
}

// GetBounds returns the bounding rectangle of all operations
func (b *Batch) GetBounds() Rectangle {
	if len(b.operations) == 0 {
		return Rectangle{}
	}
	
	minX, minY := float64(1e9), float64(1e9)
	maxX, maxY := float64(-1e9), float64(-1e9)
	
	for _, op := range b.operations {
		x, y, w, h := op.GetBounds()
		if x < minX {
			minX = x
		}
		if y < minY {
			minY = y
		}
		if x+w > maxX {
			maxX = x + w
		}
		if y+h > maxY {
			maxY = y + h
		}
	}
	
	return Rectangle{
		X:      minX,
		Y:      minY,
		Width:  maxX - minX,
		Height: maxY - minY,
	}
}

// Specific batch operation implementations

// BatchCircle represents a circle drawing operation
type BatchCircle struct {
	X, Y, Radius float64
	Color        color.Color
	Fill         bool
}

func (op BatchCircle) Execute(ctx *Context) {
	ctx.SetColor(op.Color)
	ctx.DrawCircle(op.X, op.Y, op.Radius)
	if op.Fill {
		ctx.Fill()
	} else {
		ctx.Stroke()
	}
}

func (op BatchCircle) GetBounds() (x, y, width, height float64) {
	return op.X - op.Radius, op.Y - op.Radius, op.Radius * 2, op.Radius * 2
}

func (op BatchCircle) GetType() BatchOpType {
	return BatchOpCircle
}

// BatchRectangle represents a rectangle drawing operation
type BatchRectangle struct {
	X, Y, Width, Height float64
	Color               color.Color
	Fill                bool
}

func (op BatchRectangle) Execute(ctx *Context) {
	ctx.SetColor(op.Color)
	ctx.DrawRectangle(op.X, op.Y, op.Width, op.Height)
	if op.Fill {
		ctx.Fill()
	} else {
		ctx.Stroke()
	}
}

func (op BatchRectangle) GetBounds() (x, y, width, height float64) {
	return op.X, op.Y, op.Width, op.Height
}

func (op BatchRectangle) GetType() BatchOpType {
	return BatchOpRectangle
}

// BatchLine represents a line drawing operation
type BatchLine struct {
	X1, Y1, X2, Y2 float64
	Color          color.Color
	Width          float64
}

func (op BatchLine) Execute(ctx *Context) {
	ctx.SetColor(op.Color)
	ctx.SetLineWidth(op.Width)
	ctx.DrawLine(op.X1, op.Y1, op.X2, op.Y2)
	ctx.Stroke()
}

func (op BatchLine) GetBounds() (x, y, width, height float64) {
	minX := op.X1
	if op.X2 < minX {
		minX = op.X2
	}
	minY := op.Y1
	if op.Y2 < minY {
		minY = op.Y2
	}
	maxX := op.X1
	if op.X2 > maxX {
		maxX = op.X2
	}
	maxY := op.Y1
	if op.Y2 > maxY {
		maxY = op.Y2
	}
	return minX, minY, maxX - minX, maxY - minY
}

func (op BatchLine) GetType() BatchOpType {
	return BatchOpLine
}

// BatchText represents a text drawing operation
type BatchText struct {
	Text  string
	X, Y  float64
	Color color.Color
}

func (op BatchText) Execute(ctx *Context) {
	ctx.SetColor(op.Color)
	ctx.DrawString(op.Text, op.X, op.Y)
}

func (op BatchText) GetBounds() (x, y, width, height float64) {
	// Simplified bounds calculation
	// In a real implementation, this would use actual text metrics
	return op.X, op.Y - 20, float64(len(op.Text)) * 10, 20
}

func (op BatchText) GetType() BatchOpType {
	return BatchOpText
}

// BatchImage represents an image drawing operation
type BatchImage struct {
	Image image.Image
	X, Y  int
}

func (op BatchImage) Execute(ctx *Context) {
	ctx.DrawImage(op.Image, op.X, op.Y)
}

func (op BatchImage) GetBounds() (x, y, width, height float64) {
	bounds := op.Image.Bounds()
	return float64(op.X), float64(op.Y), float64(bounds.Max.X), float64(bounds.Max.Y)
}

func (op BatchImage) GetType() BatchOpType {
	return BatchOpImage
}

// Context extensions for batch operations

// BeginBatch starts a new batch operation
func (dc *Context) BeginBatch() *Batch {
	return NewBatch()
}

// ExecuteBatch executes a batch of operations
func (dc *Context) ExecuteBatch(batch *Batch) {
	batch.Execute(dc)
}

// BatchCircles draws multiple circles efficiently
func (dc *Context) BatchCircles(circles []BatchCircle) {
	batch := NewBatch()
	for _, circle := range circles {
		batch.Add(circle)
	}
	dc.ExecuteBatch(batch)
}

// BatchRectangles draws multiple rectangles efficiently
func (dc *Context) BatchRectangles(rectangles []BatchRectangle) {
	batch := NewBatch()
	for _, rect := range rectangles {
		batch.Add(rect)
	}
	dc.ExecuteBatch(batch)
}

// BatchLines draws multiple lines efficiently
func (dc *Context) BatchLines(lines []BatchLine) {
	batch := NewBatch()
	for _, line := range lines {
		batch.Add(line)
	}
	dc.ExecuteBatch(batch)
}

// BatchTexts draws multiple text strings efficiently
func (dc *Context) BatchTexts(texts []BatchText) {
	batch := NewBatch()
	for _, text := range texts {
		batch.Add(text)
	}
	dc.ExecuteBatch(batch)
}

// BatchImages draws multiple images efficiently
func (dc *Context) BatchImages(images []BatchImage) {
	batch := NewBatch()
	for _, img := range images {
		batch.Add(img)
	}
	dc.ExecuteBatch(batch)
}

// Advanced batching features

// ConditionalBatch executes operations only if they're within the viewport
type ConditionalBatch struct {
	*Batch
	viewport Rectangle
}

// NewConditionalBatch creates a batch that only executes operations within the viewport
func NewConditionalBatch(viewport Rectangle) *ConditionalBatch {
	return &ConditionalBatch{
		Batch:    NewBatch(),
		viewport: viewport,
	}
}

// Execute executes only operations that intersect with the viewport
func (cb *ConditionalBatch) Execute(ctx *Context) {
	for _, op := range cb.operations {
		x, y, w, h := op.GetBounds()
		opBounds := Rectangle{X: x, Y: y, Width: w, Height: h}
		
		if cb.intersects(opBounds, cb.viewport) {
			op.Execute(ctx)
		}
	}
}

// intersects checks if two rectangles intersect
func (cb *ConditionalBatch) intersects(a, b Rectangle) bool {
	return a.X < b.X+b.Width && a.X+a.Width > b.X &&
		a.Y < b.Y+b.Height && a.Y+a.Height > b.Y
}
