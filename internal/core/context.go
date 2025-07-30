// Package advancegg provides a simple API for rendering 2D graphics in pure Go.
package core

import (
	"errors"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"io"
	"math"
	"os"
	"strings"

	"github.com/golang/freetype/raster"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/draw"
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/f64"
)

type LineCap int

const (
	LineCapRound LineCap = iota
	LineCapButt
	LineCapSquare
)

type LineJoin int

const (
	LineJoinRound LineJoin = iota
	LineJoinBevel
)

type FillRule int

const (
	FillRuleWinding FillRule = iota
	FillRuleEvenOdd
)

type Align int

const (
	AlignLeft Align = iota
	AlignCenter
	AlignRight
)

var (
	defaultFillStyle   = NewSolidPattern(color.White)
	defaultStrokeStyle = NewSolidPattern(color.Black)
)

type Context struct {
	width         int
	height        int
	rasterizer    *raster.Rasterizer
	im            *image.RGBA
	mask          *image.Alpha
	color         color.Color
	fillPattern   Pattern
	strokePattern Pattern
	strokePath    raster.Path
	fillPath      raster.Path
	start         Point
	current       Point
	hasCurrent    bool
	dashes        []float64
	dashOffset    float64
	lineWidth     float64
	lineCap       LineCap
	lineJoin      LineJoin
	fillRule      FillRule
	fontFace      font.Face
	fontHeight    float64
	matrix        Matrix
	stack         []*Context
	// Shadow properties
	shadowColor   color.Color
	shadowOffsetX float64
	shadowOffsetY float64
	shadowBlur    float64
	// Layer system
	layerManager   *LayerManager
	useLayerSystem bool
	// Non-destructive editing
	editStack *EditStack
	// Guides and alignment
	guideManager *GuideManager
	// Unicode shaping
	textShaper *TextShaper
	// Emoji rendering
	emojiRenderer *EmojiRenderer
}

// NewContext creates a new image.RGBA with the specified width and height
// and prepares a context for rendering onto that image.
func NewContext(width, height int) *Context {
	return NewContextForRGBA(image.NewRGBA(image.Rect(0, 0, width, height)))
}

// NewContextForImage copies the specified image into a new image.RGBA
// and prepares a context for rendering onto that image.
func NewContextForImage(im image.Image) *Context {
	return NewContextForRGBA(imageToRGBA(im))
}

// NewContextForRGBA prepares a context for rendering onto the specified image.
// No copy is made.
func NewContextForRGBA(im *image.RGBA) *Context {
	w := im.Bounds().Size().X
	h := im.Bounds().Size().Y
	return &Context{
		width:         w,
		height:        h,
		rasterizer:    raster.NewRasterizer(w, h),
		im:            im,
		color:         color.Transparent,
		fillPattern:   defaultFillStyle,
		strokePattern: defaultStrokeStyle,
		lineWidth:     1,
		fillRule:      FillRuleWinding,
		fontFace:      basicfont.Face7x13,
		fontHeight:    13,
		matrix:        Identity(),
	}
}

// GetCurrentPoint will return the current point and if there is a current point.
// The point will have been transformed by the context's transformation matrix.
func (dc *Context) GetCurrentPoint() (Point, bool) {
	if dc.hasCurrent {
		return dc.current, true
	}
	return Point{}, false
}

// Image returns the image that has been drawn by this context.
func (dc *Context) Image() image.Image {
	return dc.im
}

// Width returns the width of the image in pixels.
func (dc *Context) Width() int {
	return dc.width
}

// Height returns the height of the image in pixels.
func (dc *Context) Height() int {
	return dc.height
}

// SavePNG encodes the image as a PNG and writes it to disk.
func (dc *Context) SavePNG(path string) error {
	return SavePNG(path, dc.im)
}

// SaveJPEG saves the current image as a JPEG file.
func (dc *Context) SaveJPEG(path string, quality int) error {
	return SaveJPEG(path, dc.im, quality)
}

// SaveGIF saves the current image as a GIF file.
func (dc *Context) SaveGIF(path string) error {
	return SaveGIF(path, dc.im)
}

// SaveBMP saves the current image as a BMP file.
func (dc *Context) SaveBMP(path string) error {
	return SaveBMP(path, dc.im)
}

// SaveTIFF saves the current image as a TIFF file.
func (dc *Context) SaveTIFF(path string) error {
	return SaveTIFF(path, dc.im)
}

// ImageData methods

// GetImageData returns the current image as ImageData for pixel manipulation
func (dc *Context) GetImageData() *ImageData {
	return NewImageDataFromImage(dc.im)
}

// GetImageDataRegion returns a region of the current image as ImageData
func (dc *Context) GetImageDataRegion(x, y, width, height int) *ImageData {
	imageData := NewImageDataFromImage(dc.im)
	return imageData.GetSubImageData(x, y, width, height)
}

// PutImageData replaces the current image with ImageData
func (dc *Context) PutImageData(imageData *ImageData) {
	dc.im = imageData.ToImage()
}

// PutImageDataAt places ImageData at the specified coordinates
func (dc *Context) PutImageDataAt(imageData *ImageData, x, y int) {
	currentData := NewImageDataFromImage(dc.im)
	currentData.CopyFrom(imageData, 0, 0, imageData.Width, imageData.Height, x, y)
	dc.im = currentData.ToImage()
}

// CreateImageData creates a new ImageData with the specified dimensions
func (dc *Context) CreateImageData(width, height int) *ImageData {
	return NewImageData(width, height)
}

// SaveJPG encodes the image as a JPG and writes it to disk.
func (dc *Context) SaveJPG(path string, quality int) error {
	return SaveJPG(path, dc.im, quality)
}

// EncodePNG encodes the image as a PNG and writes it to the provided io.Writer.
func (dc *Context) EncodePNG(w io.Writer) error {
	return png.Encode(w, dc.im)
}

// EncodeJPG encodes the image as a JPG and writes it to the provided io.Writer
// in JPEG 4:2:0 baseline format with the given options.
// Default parameters are used if a nil *jpeg.Options is passed.
func (dc *Context) EncodeJPG(w io.Writer, o *jpeg.Options) error {
	return jpeg.Encode(w, dc.im, o)
}

// SetDash sets the current dash pattern to use. Call with zero arguments to
// disable dashes. The values specify the lengths of each dash, with
// alternating on and off lengths.
func (dc *Context) SetDash(dashes ...float64) {
	dc.dashes = dashes
}

// SetDashOffset sets the initial offset into the dash pattern to use when
// stroking dashed paths.
func (dc *Context) SetDashOffset(offset float64) {
	dc.dashOffset = offset
}

func (dc *Context) SetLineWidth(lineWidth float64) {
	dc.lineWidth = lineWidth
}

func (dc *Context) SetLineCap(lineCap LineCap) {
	dc.lineCap = lineCap
}

func (dc *Context) SetLineCapRound() {
	dc.lineCap = LineCapRound
}

func (dc *Context) SetLineCapButt() {
	dc.lineCap = LineCapButt
}

func (dc *Context) SetLineCapSquare() {
	dc.lineCap = LineCapSquare
}

func (dc *Context) SetLineJoin(lineJoin LineJoin) {
	dc.lineJoin = lineJoin
}

func (dc *Context) SetLineJoinRound() {
	dc.lineJoin = LineJoinRound
}

func (dc *Context) SetLineJoinBevel() {
	dc.lineJoin = LineJoinBevel
}

func (dc *Context) SetFillRule(fillRule FillRule) {
	dc.fillRule = fillRule
}

func (dc *Context) SetFillRuleWinding() {
	dc.fillRule = FillRuleWinding
}

func (dc *Context) SetFillRuleEvenOdd() {
	dc.fillRule = FillRuleEvenOdd
}

// Color Setters

func (dc *Context) setFillAndStrokeColor(c color.Color) {
	dc.color = c
	dc.fillPattern = NewSolidPattern(c)
	dc.strokePattern = NewSolidPattern(c)
}

// SetFillStyle sets current fill style
func (dc *Context) SetFillStyle(pattern Pattern) {
	// if pattern is SolidPattern, also change dc.color(for dc.Clear, dc.drawString)
	if fillStyle, ok := pattern.(*solidPattern); ok {
		dc.color = fillStyle.color
	}
	dc.fillPattern = pattern
}

// SetStrokeStyle sets current stroke style
func (dc *Context) SetStrokeStyle(pattern Pattern) {
	dc.strokePattern = pattern
}

// SetColor sets the current color(for both fill and stroke).
func (dc *Context) SetColor(c color.Color) {
	dc.setFillAndStrokeColor(c)
}

// SetHexColor sets the current color using a hex string. The leading pound
// sign (#) is optional. Both 3- and 6-digit variations are supported. 8 digits
// may be provided to set the alpha value as well.
func (dc *Context) SetHexColor(x string) {
	r, g, b, a := parseHexColor(x)
	dc.SetRGBA255(r, g, b, a)
}

// SetCMYK sets the current color using CMYK values (0-1 range)
func (dc *Context) SetCMYK(c, m, y, k float64) {
	cmyk := CMYK{C: c, M: m, Y: y, K: k}
	rgb := cmyk.ToRGB()
	dc.SetRGBA(rgb.R, rgb.G, rgb.B, 1.0)
}

// SetHSV sets the current color using HSV values
func (dc *Context) SetHSV(h, s, v float64) {
	hsv := HSV{H: h, S: s, V: v}
	rgb := hsv.ToRGB()
	dc.SetRGBA(rgb.R, rgb.G, rgb.B, 1.0)
}

// SetHSL sets the current color using HSL values
func (dc *Context) SetHSL(h, s, l float64) {
	hsl := HSL{H: h, S: s, L: l}
	rgb := hsl.ToRGB()
	dc.SetRGBA(rgb.R, rgb.G, rgb.B, 1.0)
}

// SetLAB sets the current color using LAB values
func (dc *Context) SetLAB(l, a, b float64) {
	lab := LAB{L: l, A: a, B: b}
	rgb := lab.ToRGB()
	dc.SetRGBA(rgb.R, rgb.G, rgb.B, 1.0)
}

// Shadow methods

// SetShadow sets the shadow properties
func (dc *Context) SetShadow(offsetX, offsetY, blur float64, shadowColor color.Color) {
	dc.shadowOffsetX = offsetX
	dc.shadowOffsetY = offsetY
	dc.shadowBlur = blur
	dc.shadowColor = shadowColor
}

// SetShadowRGBA sets the shadow with RGBA color
func (dc *Context) SetShadowRGBA(offsetX, offsetY, blur, r, g, b, a float64) {
	dc.SetShadow(offsetX, offsetY, blur, color.RGBA{
		uint8(r * 255),
		uint8(g * 255),
		uint8(b * 255),
		uint8(a * 255),
	})
}

// ClearShadow removes the shadow effect
func (dc *Context) ClearShadow() {
	dc.shadowOffsetX = 0
	dc.shadowOffsetY = 0
	dc.shadowBlur = 0
	dc.shadowColor = nil
}

// HasShadow returns true if shadow is enabled
func (dc *Context) HasShadow() bool {
	return dc.shadowColor != nil && (dc.shadowOffsetX != 0 || dc.shadowOffsetY != 0 || dc.shadowBlur > 0)
}

// Layer system methods

// EnableLayers enables the layer system
func (dc *Context) EnableLayers() {
	if dc.layerManager == nil {
		dc.layerManager = NewLayerManager(dc.width, dc.height)
	}
	dc.useLayerSystem = true
}

// DisableLayers disables the layer system
func (dc *Context) DisableLayers() {
	dc.useLayerSystem = false
}

// GetLayerManager returns the layer manager
func (dc *Context) GetLayerManager() *LayerManager {
	return dc.layerManager
}

// AddLayer adds a new layer
func (dc *Context) AddLayer(name string) *Layer {
	if dc.layerManager == nil {
		dc.EnableLayers()
	}
	return dc.layerManager.AddLayer(name)
}

// SetActiveLayer sets the active layer
func (dc *Context) SetActiveLayer(index int) bool {
	if dc.layerManager == nil {
		return false
	}
	return dc.layerManager.SetActiveLayer(index)
}

// SetActiveLayerByName sets the active layer by name
func (dc *Context) SetActiveLayerByName(name string) bool {
	if dc.layerManager == nil {
		return false
	}
	return dc.layerManager.SetActiveLayerByName(name)
}

// GetActiveLayer returns the active layer
func (dc *Context) GetActiveLayer() *Layer {
	if dc.layerManager == nil {
		return nil
	}
	return dc.layerManager.GetActiveLayer()
}

// CompositeToImage renders all layers to the main image
func (dc *Context) CompositeToImage() {
	if dc.layerManager == nil || !dc.useLayerSystem {
		return
	}

	composited := dc.layerManager.Composite()
	dc.im = composited
}

// SetRGBA255 sets the current color. r, g, b, a values should be between 0 and
// 255, inclusive.
func (dc *Context) SetRGBA255(r, g, b, a int) {
	dc.color = color.NRGBA{uint8(r), uint8(g), uint8(b), uint8(a)}
	dc.setFillAndStrokeColor(dc.color)
}

// SetRGB255 sets the current color. r, g, b values should be between 0 and 255,
// inclusive. Alpha will be set to 255 (fully opaque).
func (dc *Context) SetRGB255(r, g, b int) {
	dc.SetRGBA255(r, g, b, 255)
}

// SetRGBA sets the current color. r, g, b, a values should be between 0 and 1,
// inclusive.
func (dc *Context) SetRGBA(r, g, b, a float64) {
	dc.color = color.NRGBA{
		uint8(r * 255),
		uint8(g * 255),
		uint8(b * 255),
		uint8(a * 255),
	}
	dc.setFillAndStrokeColor(dc.color)
}

// SetRGB sets the current color. r, g, b values should be between 0 and 1,
// inclusive. Alpha will be set to 1 (fully opaque).
func (dc *Context) SetRGB(r, g, b float64) {
	dc.SetRGBA(r, g, b, 1)
}

// Path Manipulation

// MoveTo starts a new subpath within the current path starting at the
// specified point.
func (dc *Context) MoveTo(x, y float64) {
	if dc.hasCurrent {
		dc.fillPath.Add1(dc.start.Fixed())
	}
	x, y = dc.TransformPoint(x, y)
	p := Point{x, y}
	dc.strokePath.Start(p.Fixed())
	dc.fillPath.Start(p.Fixed())
	dc.start = p
	dc.current = p
	dc.hasCurrent = true
}

// LineTo adds a line segment to the current path starting at the current
// point. If there is no current point, it is equivalent to MoveTo(x, y)
func (dc *Context) LineTo(x, y float64) {
	if !dc.hasCurrent {
		dc.MoveTo(x, y)
	} else {
		x, y = dc.TransformPoint(x, y)
		p := Point{x, y}
		dc.strokePath.Add1(p.Fixed())
		dc.fillPath.Add1(p.Fixed())
		dc.current = p
	}
}

// QuadraticTo adds a quadratic bezier curve to the current path starting at
// the current point. If there is no current point, it first performs
// MoveTo(x1, y1)
func (dc *Context) QuadraticTo(x1, y1, x2, y2 float64) {
	if !dc.hasCurrent {
		dc.MoveTo(x1, y1)
	}
	x1, y1 = dc.TransformPoint(x1, y1)
	x2, y2 = dc.TransformPoint(x2, y2)
	p1 := Point{x1, y1}
	p2 := Point{x2, y2}
	dc.strokePath.Add2(p1.Fixed(), p2.Fixed())
	dc.fillPath.Add2(p1.Fixed(), p2.Fixed())
	dc.current = p2
}

// CubicTo adds a cubic bezier curve to the current path starting at the
// current point. If there is no current point, it first performs
// MoveTo(x1, y1). Because freetype/raster does not support cubic beziers,
// this is emulated with many small line segments.
func (dc *Context) CubicTo(x1, y1, x2, y2, x3, y3 float64) {
	if !dc.hasCurrent {
		dc.MoveTo(x1, y1)
	}
	x0, y0 := dc.current.X, dc.current.Y
	x1, y1 = dc.TransformPoint(x1, y1)
	x2, y2 = dc.TransformPoint(x2, y2)
	x3, y3 = dc.TransformPoint(x3, y3)
	points := CubicBezier(x0, y0, x1, y1, x2, y2, x3, y3)
	previous := dc.current.Fixed()
	for _, p := range points[1:] {
		f := p.Fixed()
		if f == previous {
			// TODO: this fixes some rendering issues but not all
			continue
		}
		previous = f
		dc.strokePath.Add1(f)
		dc.fillPath.Add1(f)
		dc.current = p
	}
}

// ClosePath adds a line segment from the current point to the beginning
// of the current subpath. If there is no current point, this is a no-op.
func (dc *Context) ClosePath() {
	if dc.hasCurrent {
		dc.strokePath.Add1(dc.start.Fixed())
		dc.fillPath.Add1(dc.start.Fixed())
		dc.current = dc.start
	}
}

// ClearPath clears the current path. There is no current point after this
// operation.
func (dc *Context) ClearPath() {
	dc.strokePath.Clear()
	dc.fillPath.Clear()
	dc.hasCurrent = false
}

// NewSubPath starts a new subpath within the current path. There is no current
// point after this operation.
func (dc *Context) NewSubPath() {
	if dc.hasCurrent {
		dc.fillPath.Add1(dc.start.Fixed())
	}
	dc.hasCurrent = false
}

// Path2D methods

// DrawPath2D draws a Path2D object to the context
func (dc *Context) DrawPath2D(path2d *Path2D) {
	if path2d == nil || path2d.IsEmpty() {
		return
	}

	// Convert Path2D to the context's path format
	path2dPath := path2d.GetPath()
	for i := 0; i < len(path2dPath); {
		switch path2dPath[i] {
		case 0: // MoveTo
			if i+2 < len(path2dPath) {
				x := float64(path2dPath[i+1]) / 64
				y := float64(path2dPath[i+2]) / 64
				dc.MoveTo(x, y)
				i += 3
			} else {
				i++
			}
		case 1: // LineTo
			if i+2 < len(path2dPath) {
				x := float64(path2dPath[i+1]) / 64
				y := float64(path2dPath[i+2]) / 64
				dc.LineTo(x, y)
				i += 3
			} else {
				i++
			}
		case 2: // QuadTo
			if i+4 < len(path2dPath) {
				cpx := float64(path2dPath[i+1]) / 64
				cpy := float64(path2dPath[i+2]) / 64
				x := float64(path2dPath[i+3]) / 64
				y := float64(path2dPath[i+4]) / 64
				dc.QuadraticTo(cpx, cpy, x, y)
				i += 5
			} else {
				i++
			}
		case 3: // CubeTo
			if i+6 < len(path2dPath) {
				cp1x := float64(path2dPath[i+1]) / 64
				cp1y := float64(path2dPath[i+2]) / 64
				cp2x := float64(path2dPath[i+3]) / 64
				cp2y := float64(path2dPath[i+4]) / 64
				x := float64(path2dPath[i+5]) / 64
				y := float64(path2dPath[i+6]) / 64
				dc.CubicTo(cp1x, cp1y, cp2x, cp2y, x, y)
				i += 7
			} else {
				i++
			}
		default:
			i++
		}
	}
}

// FillPath2D fills a Path2D object
func (dc *Context) FillPath2D(path2d *Path2D) {
	if path2d == nil || path2d.IsEmpty() {
		return
	}
	dc.Push()
	dc.ClearPath()
	dc.DrawPath2D(path2d)
	dc.Fill()
	dc.Pop()
}

// StrokePath2D strokes a Path2D object
func (dc *Context) StrokePath2D(path2d *Path2D) {
	if path2d == nil || path2d.IsEmpty() {
		return
	}
	dc.Push()
	dc.ClearPath()
	dc.DrawPath2D(path2d)
	dc.Stroke()
	dc.Pop()
}

// ClipPath2D sets a Path2D object as the clipping region
func (dc *Context) ClipPath2D(path2d *Path2D) {
	if path2d == nil || path2d.IsEmpty() {
		return
	}
	dc.Push()
	dc.ClearPath()
	dc.DrawPath2D(path2d)
	dc.Clip()
	dc.Pop()
}

// IsPointInPath2D tests if a point is inside a Path2D object
func (dc *Context) IsPointInPath2D(path2d *Path2D, x, y float64) bool {
	if path2d == nil || path2d.IsEmpty() {
		return false
	}

	// Save current state
	dc.Push()
	dc.ClearPath()
	dc.DrawPath2D(path2d)

	// Use a simple point-in-polygon test (simplified implementation)
	result := dc.isPointInCurrentPath(x, y)

	dc.Pop()
	return result
}

// Helper method for point-in-path testing (simplified implementation)
func (dc *Context) isPointInCurrentPath(x, y float64) bool {
	// This is a placeholder implementation
	// A full implementation would use proper point-in-polygon algorithms
	// like ray casting or winding number calculation
	return false
}

// Path Drawing

func (dc *Context) capper() raster.Capper {
	switch dc.lineCap {
	case LineCapButt:
		return raster.ButtCapper
	case LineCapRound:
		return raster.RoundCapper
	case LineCapSquare:
		return raster.SquareCapper
	}
	return nil
}

func (dc *Context) joiner() raster.Joiner {
	switch dc.lineJoin {
	case LineJoinBevel:
		return raster.BevelJoiner
	case LineJoinRound:
		return raster.RoundJoiner
	}
	return nil
}

func (dc *Context) stroke(painter raster.Painter) {
	path := dc.strokePath
	if len(dc.dashes) > 0 {
		path = dashed(path, dc.dashes, dc.dashOffset)
	} else {
		// TODO: this is a temporary workaround to remove tiny segments
		// that result in rendering issues
		path = rasterPath(flattenPath(path))
	}
	r := dc.rasterizer
	r.UseNonZeroWinding = true
	r.Clear()
	r.AddStroke(path, fix(dc.lineWidth), dc.capper(), dc.joiner())
	r.Rasterize(painter)
}

func (dc *Context) fill(painter raster.Painter) {
	path := dc.fillPath
	if dc.hasCurrent {
		path = make(raster.Path, len(dc.fillPath))
		copy(path, dc.fillPath)
		path.Add1(dc.start.Fixed())
	}
	r := dc.rasterizer
	r.UseNonZeroWinding = dc.fillRule == FillRuleWinding
	r.Clear()
	r.AddPath(path)
	r.Rasterize(painter)
}

// StrokePreserve strokes the current path with the current color, line width,
// line cap, line join and dash settings. The path is preserved after this
// operation.
func (dc *Context) StrokePreserve() {
	var painter raster.Painter
	if dc.mask == nil {
		if pattern, ok := dc.strokePattern.(*solidPattern); ok {
			// with a nil mask and a solid color pattern, we can be more efficient
			// TODO: refactor so we don't have to do this type assertion stuff?
			p := raster.NewRGBAPainter(dc.im)
			p.SetColor(pattern.color)
			painter = p
		}
	}
	if painter == nil {
		painter = newPatternPainter(dc.im, dc.mask, dc.strokePattern)
	}
	dc.stroke(painter)
}

// Stroke strokes the current path with the current color, line width,
// line cap, line join and dash settings. The path is cleared after this
// operation.
func (dc *Context) Stroke() {
	dc.StrokePreserve()
	dc.ClearPath()
}

// FillPreserve fills the current path with the current color. Open subpaths
// are implicity closed. The path is preserved after this operation.
func (dc *Context) FillPreserve() {
	var painter raster.Painter
	if dc.mask == nil {
		if pattern, ok := dc.fillPattern.(*solidPattern); ok {
			// with a nil mask and a solid color pattern, we can be more efficient
			// TODO: refactor so we don't have to do this type assertion stuff?
			p := raster.NewRGBAPainter(dc.im)
			p.SetColor(pattern.color)
			painter = p
		}
	}
	if painter == nil {
		painter = newPatternPainter(dc.im, dc.mask, dc.fillPattern)
	}
	dc.fill(painter)
}

// Fill fills the current path with the current color. Open subpaths
// are implicity closed. The path is cleared after this operation.
func (dc *Context) Fill() {
	dc.FillPreserve()
	dc.ClearPath()
}

// ClipPreserve updates the clipping region by intersecting the current
// clipping region with the current path as it would be filled by dc.Fill().
// The path is preserved after this operation.
func (dc *Context) ClipPreserve() {
	clip := image.NewAlpha(image.Rect(0, 0, dc.width, dc.height))
	painter := raster.NewAlphaOverPainter(clip)
	dc.fill(painter)
	if dc.mask == nil {
		dc.mask = clip
	} else {
		mask := image.NewAlpha(image.Rect(0, 0, dc.width, dc.height))
		draw.DrawMask(mask, mask.Bounds(), clip, image.ZP, dc.mask, image.ZP, draw.Over)
		dc.mask = mask
	}
}

// SetMask allows you to directly set the *image.Alpha to be used as a clipping
// mask. It must be the same size as the context, else an error is returned
// and the mask is unchanged.
func (dc *Context) SetMask(mask *image.Alpha) error {
	if mask.Bounds().Size() != dc.im.Bounds().Size() {
		return errors.New("mask size must match context size")
	}
	dc.mask = mask
	return nil
}

// AsMask returns an *image.Alpha representing the alpha channel of this
// context. This can be useful for advanced clipping operations where you first
// render the mask geometry and then use it as a mask.
func (dc *Context) AsMask() *image.Alpha {
	mask := image.NewAlpha(dc.im.Bounds())
	draw.Draw(mask, dc.im.Bounds(), dc.im, image.ZP, draw.Src)
	return mask
}

// InvertMask inverts the alpha values in the current clipping mask such that
// a fully transparent region becomes fully opaque and vice versa.
func (dc *Context) InvertMask() {
	if dc.mask == nil {
		dc.mask = image.NewAlpha(dc.im.Bounds())
	} else {
		for i, a := range dc.mask.Pix {
			dc.mask.Pix[i] = 255 - a
		}
	}
}

// Clip updates the clipping region by intersecting the current
// clipping region with the current path as it would be filled by dc.Fill().
// The path is cleared after this operation.
func (dc *Context) Clip() {
	dc.ClipPreserve()
	dc.ClearPath()
}

// ResetClip clears the clipping region.
func (dc *Context) ResetClip() {
	dc.mask = nil
}

// Convenient Drawing Functions

// Clear fills the entire image with the current color.
func (dc *Context) Clear() {
	src := image.NewUniform(dc.color)
	draw.Draw(dc.im, dc.im.Bounds(), src, image.ZP, draw.Src)
}

// SetPixel sets the color of the specified pixel using the current color.
func (dc *Context) SetPixel(x, y int) {
	dc.im.Set(x, y, dc.color)
}

// DrawPoint is like DrawCircle but ensures that a circle of the specified
// size is drawn regardless of the current transformation matrix. The position
// is still transformed, but not the shape of the point.
func (dc *Context) DrawPoint(x, y, r float64) {
	dc.Push()
	tx, ty := dc.TransformPoint(x, y)
	dc.Identity()
	dc.DrawCircle(tx, ty, r)
	dc.Pop()
}

func (dc *Context) DrawLine(x1, y1, x2, y2 float64) {
	dc.MoveTo(x1, y1)
	dc.LineTo(x2, y2)
}

func (dc *Context) DrawRectangle(x, y, w, h float64) {
	dc.NewSubPath()
	dc.MoveTo(x, y)
	dc.LineTo(x+w, y)
	dc.LineTo(x+w, y+h)
	dc.LineTo(x, y+h)
	dc.ClosePath()
}

func (dc *Context) DrawRoundedRectangle(x, y, w, h, r float64) {
	x0, x1, x2, x3 := x, x+r, x+w-r, x+w
	y0, y1, y2, y3 := y, y+r, y+h-r, y+h
	dc.NewSubPath()
	dc.MoveTo(x1, y0)
	dc.LineTo(x2, y0)
	dc.DrawArc(x2, y1, r, Radians(270), Radians(360))
	dc.LineTo(x3, y2)
	dc.DrawArc(x2, y2, r, Radians(0), Radians(90))
	dc.LineTo(x1, y3)
	dc.DrawArc(x1, y2, r, Radians(90), Radians(180))
	dc.LineTo(x0, y1)
	dc.DrawArc(x1, y1, r, Radians(180), Radians(270))
	dc.ClosePath()
}

func (dc *Context) DrawEllipticalArc(x, y, rx, ry, angle1, angle2 float64) {
	const n = 16
	for i := 0; i < n; i++ {
		p1 := float64(i+0) / n
		p2 := float64(i+1) / n
		a1 := angle1 + (angle2-angle1)*p1
		a2 := angle1 + (angle2-angle1)*p2
		x0 := x + rx*math.Cos(a1)
		y0 := y + ry*math.Sin(a1)
		x1 := x + rx*math.Cos((a1+a2)/2)
		y1 := y + ry*math.Sin((a1+a2)/2)
		x2 := x + rx*math.Cos(a2)
		y2 := y + ry*math.Sin(a2)
		cx := 2*x1 - x0/2 - x2/2
		cy := 2*y1 - y0/2 - y2/2
		if i == 0 {
			if dc.hasCurrent {
				dc.LineTo(x0, y0)
			} else {
				dc.MoveTo(x0, y0)
			}
		}
		dc.QuadraticTo(cx, cy, x2, y2)
	}
}

func (dc *Context) DrawEllipse(x, y, rx, ry float64) {
	dc.NewSubPath()
	dc.DrawEllipticalArc(x, y, rx, ry, 0, 2*math.Pi)
	dc.ClosePath()
}

func (dc *Context) DrawArc(x, y, r, angle1, angle2 float64) {
	dc.DrawEllipticalArc(x, y, r, r, angle1, angle2)
}

func (dc *Context) DrawCircle(x, y, r float64) {
	dc.NewSubPath()
	dc.DrawEllipticalArc(x, y, r, r, 0, 2*math.Pi)
	dc.ClosePath()
}

func (dc *Context) DrawRegularPolygon(n int, x, y, r, rotation float64) {
	angle := 2 * math.Pi / float64(n)
	rotation -= math.Pi / 2
	if n%2 == 0 {
		rotation += angle / 2
	}
	dc.NewSubPath()
	for i := 0; i < n; i++ {
		a := rotation + angle*float64(i)
		dc.LineTo(x+r*math.Cos(a), y+r*math.Sin(a))
	}
	dc.ClosePath()
}

// DrawImage draws the specified image at the specified point.
func (dc *Context) DrawImage(im image.Image, x, y int) {
	dc.DrawImageAnchored(im, x, y, 0, 0)
}

// DrawImageAnchored draws the specified image at the specified anchor point.
// The anchor point is x - w * ax, y - h * ay, where w, h is the size of the
// image. Use ax=0.5, ay=0.5 to center the image at the specified point.
func (dc *Context) DrawImageAnchored(im image.Image, x, y int, ax, ay float64) {
	s := im.Bounds().Size()
	x -= int(ax * float64(s.X))
	y -= int(ay * float64(s.Y))
	transformer := draw.BiLinear
	fx, fy := float64(x), float64(y)
	m := dc.matrix.Translate(fx, fy)
	s2d := f64.Aff3{m.XX, m.XY, m.X0, m.YX, m.YY, m.Y0}
	if dc.mask == nil {
		transformer.Transform(dc.im, s2d, im, im.Bounds(), draw.Over, nil)
	} else {
		transformer.Transform(dc.im, s2d, im, im.Bounds(), draw.Over, &draw.Options{
			DstMask:  dc.mask,
			DstMaskP: image.ZP,
		})
	}
}

// Text Functions

func (dc *Context) SetFontFace(fontFace font.Face) {
	dc.fontFace = fontFace
	dc.fontHeight = float64(fontFace.Metrics().Height) / 64
}

func (dc *Context) LoadFontFace(path string, points float64) error {
	face, err := LoadFontFace(path, points)
	if err == nil {
		dc.fontFace = face
		dc.fontHeight = points * 72 / 96
	}
	return err
}

// LoadTTFFace loads a TTF font file and sets it as the current font face.
func (dc *Context) LoadTTFFace(path string, points float64) error {
	face, err := LoadTTFFace(path, points)
	if err == nil {
		dc.fontFace = face
		dc.fontHeight = points * 72 / 96
	}
	return err
}

// LoadOTFFace loads an OTF font file and sets it as the current font face.
func (dc *Context) LoadOTFFace(path string, points float64) error {
	face, err := LoadOTFFace(path, points)
	if err == nil {
		dc.fontFace = face
		dc.fontHeight = points * 72 / 96
	}
	return err
}

// LoadFontFaceFromBytes loads a font from byte data and sets it as the current font face.
// Supports both TTF and OTF formats.
func (dc *Context) LoadFontFaceFromBytes(fontBytes []byte, points float64) error {
	face, err := ParseFontFace(fontBytes, points)
	if err == nil {
		dc.fontFace = face
		dc.fontHeight = points * 72 / 96
	}
	return err
}

// LoadFontFaceWithOptions loads a font with custom truetype options.
func (dc *Context) LoadFontFaceWithOptions(path string, options *truetype.Options) error {
	fontBytes, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	face, err := ParseFontFaceWithOptions(fontBytes, options)
	if err == nil {
		dc.fontFace = face
		dc.fontHeight = options.Size * 72 / 96
	}
	return err
}

func (dc *Context) FontHeight() float64 {
	return dc.fontHeight
}

func (dc *Context) drawString(im *image.RGBA, s string, x, y float64) {
	d := &font.Drawer{
		Dst:  im,
		Src:  image.NewUniform(dc.color),
		Face: dc.fontFace,
		Dot:  fixp(x, y),
	}
	// based on Drawer.DrawString() in golang.org/x/image/font/font.go
	prevC := rune(-1)
	for _, c := range s {
		if prevC >= 0 {
			d.Dot.X += d.Face.Kern(prevC, c)
		}
		dr, mask, maskp, advance, ok := d.Face.Glyph(d.Dot, c)
		if !ok {
			// TODO: is falling back on the U+FFFD glyph the responsibility of
			// the Drawer or the Face?
			// TODO: set prevC = '\ufffd'?
			continue
		}
		sr := dr.Sub(dr.Min)
		transformer := draw.BiLinear
		fx, fy := float64(dr.Min.X), float64(dr.Min.Y)
		m := dc.matrix.Translate(fx, fy)
		s2d := f64.Aff3{m.XX, m.XY, m.X0, m.YX, m.YY, m.Y0}
		transformer.Transform(d.Dst, s2d, d.Src, sr, draw.Over, &draw.Options{
			SrcMask:  mask,
			SrcMaskP: maskp,
		})
		d.Dot.X += advance
		prevC = c
	}
}

// DrawString draws the specified text at the specified point.
func (dc *Context) DrawString(s string, x, y float64) {
	dc.DrawStringAnchored(s, x, y, 0, 0)
}

// DrawStringAnchored draws the specified text at the specified anchor point.
// The anchor point is x - w * ax, y - h * ay, where w, h is the size of the
// text. Use ax=0.5, ay=0.5 to center the text at the specified point.
func (dc *Context) DrawStringAnchored(s string, x, y, ax, ay float64) {
	w, h := dc.MeasureString(s)
	x -= ax * w
	y += ay * h
	if dc.mask == nil {
		dc.drawString(dc.im, s, x, y)
	} else {
		im := image.NewRGBA(image.Rect(0, 0, dc.width, dc.height))
		dc.drawString(im, s, x, y)
		draw.DrawMask(dc.im, dc.im.Bounds(), im, image.ZP, dc.mask, image.ZP, draw.Over)
	}
}

// DrawStringWrapped word-wraps the specified string to the given max width
// and then draws it at the specified anchor point using the given line
// spacing and text alignment.
func (dc *Context) DrawStringWrapped(s string, x, y, ax, ay, width, lineSpacing float64, align Align) {
	lines := dc.WordWrap(s, width)

	// sync h formula with MeasureMultilineString
	h := float64(len(lines)) * dc.fontHeight * lineSpacing
	h -= (lineSpacing - 1) * dc.fontHeight

	x -= ax * width
	y -= ay * h
	switch align {
	case AlignLeft:
		ax = 0
	case AlignCenter:
		ax = 0.5
		x += width / 2
	case AlignRight:
		ax = 1
		x += width
	}
	ay = 1
	for _, line := range lines {
		dc.DrawStringAnchored(line, x, y, ax, ay)
		y += dc.fontHeight * lineSpacing
	}
}

func (dc *Context) MeasureMultilineString(s string, lineSpacing float64) (width, height float64) {
	lines := strings.Split(s, "\n")

	// sync h formula with DrawStringWrapped
	height = float64(len(lines)) * dc.fontHeight * lineSpacing
	height -= (lineSpacing - 1) * dc.fontHeight

	d := &font.Drawer{
		Face: dc.fontFace,
	}

	// max width from lines
	for _, line := range lines {
		adv := d.MeasureString(line)
		currentWidth := float64(adv >> 6) // from gg.Context.MeasureString
		if currentWidth > width {
			width = currentWidth
		}
	}

	return width, height
}

// MeasureString returns the rendered width and height of the specified text
// given the current font face.
func (dc *Context) MeasureString(s string) (w, h float64) {
	d := &font.Drawer{
		Face: dc.fontFace,
	}
	a := d.MeasureString(s)
	return float64(a >> 6), dc.fontHeight
}

// WordWrap wraps the specified string to the given max width and current
// font face.
func (dc *Context) WordWrap(s string, w float64) []string {
	return wordWrap(dc, s, w)
}

// TextMetrics represents detailed text measurement information
type TextMetrics struct {
	Width                    float64 // Total width of the text
	Height                   float64 // Total height of the text
	ActualBoundingBoxLeft    float64 // Distance from text baseline to left edge of bounding box
	ActualBoundingBoxRight   float64 // Distance from text baseline to right edge of bounding box
	ActualBoundingBoxAscent  float64 // Distance from text baseline to top of bounding box
	ActualBoundingBoxDescent float64 // Distance from text baseline to bottom of bounding box
	FontBoundingBoxAscent    float64 // Distance from baseline to top of font bounding box
	FontBoundingBoxDescent   float64 // Distance from baseline to bottom of font bounding box
	EmHeightAscent           float64 // Distance from baseline to top of em square
	EmHeightDescent          float64 // Distance from baseline to bottom of em square
	HangingBaseline          float64 // Distance from alphabetic baseline to hanging baseline
	AlphabeticBaseline       float64 // Distance from alphabetic baseline to alphabetic baseline (0)
	IdeographicBaseline      float64 // Distance from alphabetic baseline to ideographic baseline
}

// MeasureTextMetrics returns detailed metrics for the specified text
func (dc *Context) MeasureTextMetrics(s string) TextMetrics {
	if dc.fontFace == nil {
		return TextMetrics{}
	}

	d := &font.Drawer{
		Face: dc.fontFace,
	}

	// Get basic measurements
	advance := d.MeasureString(s)
	width := float64(advance >> 6)

	// Get font metrics
	metrics := dc.fontFace.Metrics()
	ascent := float64(metrics.Ascent >> 6)
	descent := float64(metrics.Descent >> 6)
	height := ascent + descent

	return TextMetrics{
		Width:                    width,
		Height:                   height,
		ActualBoundingBoxLeft:    0, // Simplified
		ActualBoundingBoxRight:   width,
		ActualBoundingBoxAscent:  ascent,
		ActualBoundingBoxDescent: descent,
		FontBoundingBoxAscent:    ascent,
		FontBoundingBoxDescent:   descent,
		EmHeightAscent:           ascent,
		EmHeightDescent:          descent,
		HangingBaseline:          ascent * 0.8, // Approximation
		AlphabeticBaseline:       0,
		IdeographicBaseline:      -descent * 0.5, // Approximation
	}
}

// GetFontMetrics returns metrics for the current font
func (dc *Context) GetFontMetrics() (ascent, descent, lineGap float64) {
	if dc.fontFace == nil {
		return 0, 0, 0
	}

	metrics := dc.fontFace.Metrics()
	ascent = float64(metrics.Ascent >> 6)
	descent = float64(metrics.Descent >> 6)
	lineGap = float64(metrics.Height>>6) - ascent - descent

	return ascent, descent, lineGap
}

// GetTextWidth returns the width of the specified text
func (dc *Context) GetTextWidth(s string) float64 {
	width, _ := dc.MeasureString(s)
	return width
}

// GetTextHeight returns the height of the current font
func (dc *Context) GetTextHeight() float64 {
	return dc.fontHeight
}

// GetLineHeight returns the recommended line height for the current font
func (dc *Context) GetLineHeight() float64 {
	if dc.fontFace == nil {
		return 0
	}
	metrics := dc.fontFace.Metrics()
	return float64(metrics.Height >> 6)
}

// GetBaseline returns the baseline position for the current font
func (dc *Context) GetBaseline() float64 {
	if dc.fontFace == nil {
		return 0
	}
	metrics := dc.fontFace.Metrics()
	return float64(metrics.Ascent >> 6)
}

// WrapText wraps text to fit within the specified width and returns the lines
func (dc *Context) WrapText(text string, maxWidth float64) []string {
	if dc.fontFace == nil {
		return []string{text}
	}

	words := strings.Fields(text)
	if len(words) == 0 {
		return []string{""}
	}

	var lines []string
	var currentLine strings.Builder

	for _, word := range words {
		testLine := currentLine.String()
		if testLine != "" {
			testLine += " "
		}
		testLine += word

		width, _ := dc.MeasureString(testLine)
		if width <= maxWidth {
			if currentLine.Len() > 0 {
				currentLine.WriteString(" ")
			}
			currentLine.WriteString(word)
		} else {
			if currentLine.Len() > 0 {
				lines = append(lines, currentLine.String())
				currentLine.Reset()
			}
			currentLine.WriteString(word)
		}
	}

	if currentLine.Len() > 0 {
		lines = append(lines, currentLine.String())
	}

	return lines
}

// Transformation Matrix Operations

// Identity resets the current transformation matrix to the identity matrix.
// This results in no translating, scaling, rotating, or shearing.
func (dc *Context) Identity() {
	dc.matrix = Identity()
}

// Translate updates the current matrix with a translation.
func (dc *Context) Translate(x, y float64) {
	dc.matrix = dc.matrix.Translate(x, y)
}

// Scale updates the current matrix with a scaling factor.
// Scaling occurs about the origin.
func (dc *Context) Scale(x, y float64) {
	dc.matrix = dc.matrix.Scale(x, y)
}

// ScaleAbout updates the current matrix with a scaling factor.
// Scaling occurs about the specified point.
func (dc *Context) ScaleAbout(sx, sy, x, y float64) {
	dc.Translate(x, y)
	dc.Scale(sx, sy)
	dc.Translate(-x, -y)
}

// Rotate updates the current matrix with a anticlockwise rotation.
// Rotation occurs about the origin. Angle is specified in radians.
func (dc *Context) Rotate(angle float64) {
	dc.matrix = dc.matrix.Rotate(angle)
}

// RotateAbout updates the current matrix with a anticlockwise rotation.
// Rotation occurs about the specified point. Angle is specified in radians.
func (dc *Context) RotateAbout(angle, x, y float64) {
	dc.Translate(x, y)
	dc.Rotate(angle)
	dc.Translate(-x, -y)
}

// Shear updates the current matrix with a shearing angle.
// Shearing occurs about the origin.
func (dc *Context) Shear(x, y float64) {
	dc.matrix = dc.matrix.Shear(x, y)
}

// ShearAbout updates the current matrix with a shearing angle.
// Shearing occurs about the specified point.
func (dc *Context) ShearAbout(sx, sy, x, y float64) {
	dc.Translate(x, y)
	dc.Shear(sx, sy)
	dc.Translate(-x, -y)
}

// TransformPoint multiplies the specified point by the current matrix,
// returning a transformed position.
func (dc *Context) TransformPoint(x, y float64) (tx, ty float64) {
	return dc.matrix.TransformPoint(x, y)
}

// InvertY flips the Y axis so that Y grows from bottom to top and Y=0 is at
// the bottom of the image.
func (dc *Context) InvertY() {
	dc.Translate(0, float64(dc.height))
	dc.Scale(1, -1)
}

// Stack

// Push saves the current state of the context for later retrieval. These
// can be nested.
func (dc *Context) Push() {
	x := *dc
	dc.stack = append(dc.stack, &x)
}

// Pop restores the last saved context state from the stack.
func (dc *Context) Pop() {
	before := *dc
	s := dc.stack
	x, s := s[len(s)-1], s[:len(s)-1]
	*dc = *x
	dc.mask = before.mask
	dc.strokePath = before.strokePath
	dc.fillPath = before.fillPath
	dc.start = before.start
	dc.current = before.current
	dc.hasCurrent = before.hasCurrent
}

// Non-destructive editing methods

// EnableNonDestructiveEditing enables non-destructive editing
func (dc *Context) EnableNonDestructiveEditing() {
	if dc.editStack == nil {
		dc.editStack = NewEditStack(dc.im)
	}
}

// AddEditOperation adds a non-destructive edit operation
func (dc *Context) AddEditOperation(op EditOperation) {
	if dc.editStack == nil {
		dc.EnableNonDestructiveEditing()
	}
	dc.editStack.AddOperation(op)
}

// ApplyNonDestructiveEdits applies all non-destructive edits to the image
func (dc *Context) ApplyNonDestructiveEdits() {
	if dc.editStack != nil {
		dc.im = dc.editStack.GetResult()
	}
}

// GetEditStack returns the edit stack
func (dc *Context) GetEditStack() *EditStack {
	return dc.editStack
}

// Guide system methods

// EnableGuides enables the guide system
func (dc *Context) EnableGuides() {
	if dc.guideManager == nil {
		dc.guideManager = NewGuideManager()
	}
}

// GetGuideManager returns the guide manager
func (dc *Context) GetGuideManager() *GuideManager {
	if dc.guideManager == nil {
		dc.EnableGuides()
	}
	return dc.guideManager
}

// AddGuide adds a guide
func (dc *Context) AddGuide(position float64, orientation GuideOrientation) *Guide {
	return dc.GetGuideManager().AddGuide(position, orientation)
}

// SnapPoint snaps a point to guides or grid
func (dc *Context) SnapPoint(x, y float64) (float64, float64) {
	if dc.guideManager == nil {
		return x, y
	}
	return dc.guideManager.SnapPoint(x, y)
}
