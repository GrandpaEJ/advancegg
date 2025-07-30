package core

import (
	"image"
)

// Non-destructive editing system for reversible filters and transformations

// EditOperation represents a single editing operation
type EditOperation interface {
	Apply(img *image.RGBA) *image.RGBA
	GetType() string
	GetParameters() map[string]interface{}
	SetParameters(params map[string]interface{})
	Clone() EditOperation
}

// EditStack manages a stack of non-destructive operations
type EditStack struct {
	Operations   []EditOperation
	BaseImage    *image.RGBA
	CachedResult *image.RGBA
	CacheDirty   bool
}

// NewEditStack creates a new edit stack
func NewEditStack(baseImage *image.RGBA) *EditStack {
	return &EditStack{
		Operations: make([]EditOperation, 0),
		BaseImage:  baseImage,
		CacheDirty: true,
	}
}

// AddOperation adds an operation to the stack
func (es *EditStack) AddOperation(op EditOperation) {
	es.Operations = append(es.Operations, op)
	es.CacheDirty = true
}

// InsertOperation inserts an operation at the specified index
func (es *EditStack) InsertOperation(index int, op EditOperation) {
	if index < 0 || index > len(es.Operations) {
		es.AddOperation(op)
		return
	}

	es.Operations = append(es.Operations[:index], append([]EditOperation{op}, es.Operations[index:]...)...)
	es.CacheDirty = true
}

// RemoveOperation removes an operation by index
func (es *EditStack) RemoveOperation(index int) bool {
	if index < 0 || index >= len(es.Operations) {
		return false
	}

	es.Operations = append(es.Operations[:index], es.Operations[index+1:]...)
	es.CacheDirty = true
	return true
}

// MoveOperation moves an operation from one index to another
func (es *EditStack) MoveOperation(from, to int) bool {
	if from < 0 || from >= len(es.Operations) || to < 0 || to >= len(es.Operations) {
		return false
	}

	op := es.Operations[from]
	es.Operations = append(es.Operations[:from], es.Operations[from+1:]...)

	if to > from {
		to--
	}
	es.Operations = append(es.Operations[:to], append([]EditOperation{op}, es.Operations[to:]...)...)
	es.CacheDirty = true
	return true
}

// UpdateOperation updates an operation's parameters
func (es *EditStack) UpdateOperation(index int, params map[string]interface{}) bool {
	if index < 0 || index >= len(es.Operations) {
		return false
	}

	es.Operations[index].SetParameters(params)
	es.CacheDirty = true
	return true
}

// GetResult returns the final processed image
func (es *EditStack) GetResult() *image.RGBA {
	if !es.CacheDirty && es.CachedResult != nil {
		return es.CachedResult
	}

	result := cloneImage(es.BaseImage)

	for _, op := range es.Operations {
		result = op.Apply(result)
	}

	es.CachedResult = result
	es.CacheDirty = false
	return result
}

// GetPreview returns a preview up to the specified operation index
func (es *EditStack) GetPreview(upToIndex int) *image.RGBA {
	if upToIndex < 0 {
		return cloneImage(es.BaseImage)
	}

	if upToIndex >= len(es.Operations) {
		return es.GetResult()
	}

	result := cloneImage(es.BaseImage)

	for i := 0; i <= upToIndex; i++ {
		result = es.Operations[i].Apply(result)
	}

	return result
}

// Clear removes all operations
func (es *EditStack) Clear() {
	es.Operations = es.Operations[:0]
	es.CacheDirty = true
}

// Clone creates a copy of the edit stack
func (es *EditStack) Clone() *EditStack {
	clone := &EditStack{
		Operations: make([]EditOperation, len(es.Operations)),
		BaseImage:  cloneImage(es.BaseImage),
		CacheDirty: true,
	}

	for i, op := range es.Operations {
		clone.Operations[i] = op.Clone()
	}

	return clone
}

// Specific edit operations

// BrightnessOperation adjusts image brightness
type BrightnessOperation struct {
	Amount float64
}

func (op *BrightnessOperation) Apply(img *image.RGBA) *image.RGBA {
	return SIMDColorTransform(img, func(r, g, b, a uint8) (uint8, uint8, uint8, uint8) {
		newR := clampUint8(float64(r) * op.Amount)
		newG := clampUint8(float64(g) * op.Amount)
		newB := clampUint8(float64(b) * op.Amount)
		return newR, newG, newB, a
	})
}

func (op *BrightnessOperation) GetType() string {
	return "brightness"
}

func (op *BrightnessOperation) GetParameters() map[string]interface{} {
	return map[string]interface{}{"amount": op.Amount}
}

func (op *BrightnessOperation) SetParameters(params map[string]interface{}) {
	if amount, ok := params["amount"].(float64); ok {
		op.Amount = amount
	}
}

func (op *BrightnessOperation) Clone() EditOperation {
	return &BrightnessOperation{Amount: op.Amount}
}

// ContrastOperation adjusts image contrast
type ContrastOperation struct {
	Amount float64
}

func (op *ContrastOperation) Apply(img *image.RGBA) *image.RGBA {
	return SIMDColorTransform(img, func(r, g, b, a uint8) (uint8, uint8, uint8, uint8) {
		newR := clampUint8((float64(r)-128)*op.Amount + 128)
		newG := clampUint8((float64(g)-128)*op.Amount + 128)
		newB := clampUint8((float64(b)-128)*op.Amount + 128)
		return newR, newG, newB, a
	})
}

func (op *ContrastOperation) GetType() string {
	return "contrast"
}

func (op *ContrastOperation) GetParameters() map[string]interface{} {
	return map[string]interface{}{"amount": op.Amount}
}

func (op *ContrastOperation) SetParameters(params map[string]interface{}) {
	if amount, ok := params["amount"].(float64); ok {
		op.Amount = amount
	}
}

func (op *ContrastOperation) Clone() EditOperation {
	return &ContrastOperation{Amount: op.Amount}
}

// BlurOperation applies blur
type BlurOperation struct {
	Radius float64
}

func (op *BlurOperation) Apply(img *image.RGBA) *image.RGBA {
	return SIMDBlur(img, int(op.Radius))
}

func (op *BlurOperation) GetType() string {
	return "blur"
}

func (op *BlurOperation) GetParameters() map[string]interface{} {
	return map[string]interface{}{"radius": op.Radius}
}

func (op *BlurOperation) SetParameters(params map[string]interface{}) {
	if radius, ok := params["radius"].(float64); ok {
		op.Radius = radius
	}
}

func (op *BlurOperation) Clone() EditOperation {
	return &BlurOperation{Radius: op.Radius}
}

// CropOperation crops the image
type CropOperation struct {
	X, Y, Width, Height int
}

func (op *CropOperation) Apply(img *image.RGBA) *image.RGBA {
	bounds := img.Bounds()
	cropRect := image.Rect(op.X, op.Y, op.X+op.Width, op.Y+op.Height)

	// Ensure crop rect is within bounds
	cropRect = cropRect.Intersect(bounds)

	result := image.NewRGBA(image.Rect(0, 0, cropRect.Dx(), cropRect.Dy()))

	for y := cropRect.Min.Y; y < cropRect.Max.Y; y++ {
		for x := cropRect.Min.X; x < cropRect.Max.X; x++ {
			pixel := img.RGBAAt(x, y)
			result.SetRGBA(x-cropRect.Min.X, y-cropRect.Min.Y, pixel)
		}
	}

	return result
}

func (op *CropOperation) GetType() string {
	return "crop"
}

func (op *CropOperation) GetParameters() map[string]interface{} {
	return map[string]interface{}{
		"x": op.X, "y": op.Y, "width": op.Width, "height": op.Height,
	}
}

func (op *CropOperation) SetParameters(params map[string]interface{}) {
	if x, ok := params["x"].(int); ok {
		op.X = x
	}
	if y, ok := params["y"].(int); ok {
		op.Y = y
	}
	if width, ok := params["width"].(int); ok {
		op.Width = width
	}
	if height, ok := params["height"].(int); ok {
		op.Height = height
	}
}

func (op *CropOperation) Clone() EditOperation {
	return &CropOperation{X: op.X, Y: op.Y, Width: op.Width, Height: op.Height}
}

// Helper functions

func cloneImage(img *image.RGBA) *image.RGBA {
	bounds := img.Bounds()
	clone := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			clone.SetRGBA(x, y, img.RGBAAt(x, y))
		}
	}

	return clone
}

func clampUint8(value float64) uint8 {
	if value < 0 {
		return 0
	}
	if value > 255 {
		return 255
	}
	return uint8(value)
}

// Context integration methods are in context.go
