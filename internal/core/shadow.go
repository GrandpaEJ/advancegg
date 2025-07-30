package core

import (
	"image"
	"image/color"
	"image/draw"
	"math"
)

// applyShadow applies shadow effect to the current drawing operation
func (dc *Context) applyShadow(drawFunc func()) {
	if !dc.HasShadow() {
		drawFunc()
		return
	}

	// Save current state
	originalIm := dc.im
	originalColor := dc.color

	// Create shadow image
	shadowIm := image.NewRGBA(dc.im.Bounds())
	dc.im = shadowIm
	dc.color = dc.shadowColor

	// Draw the shape for shadow
	drawFunc()

	// Apply blur if needed
	if dc.shadowBlur > 0 {
		shadowIm = dc.blurImage(shadowIm, dc.shadowBlur)
	}

	// Restore original image and color
	dc.im = originalIm
	dc.color = originalColor

	// Draw shadow with offset
	dc.drawShadowImage(shadowIm, dc.shadowOffsetX, dc.shadowOffsetY)

	// Draw the original shape
	drawFunc()
}

// drawShadowImage draws the shadow image with the specified offset
func (dc *Context) drawShadowImage(shadowIm *image.RGBA, offsetX, offsetY float64) {
	bounds := shadowIm.Bounds()
	
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			shadowPixel := shadowIm.RGBAAt(x, y)
			if shadowPixel.A > 0 {
				targetX := x + int(offsetX)
				targetY := y + int(offsetY)
				
				if targetX >= 0 && targetX < dc.width && targetY >= 0 && targetY < dc.height {
					// Alpha blend the shadow pixel
					dc.blendPixel(targetX, targetY, shadowPixel)
				}
			}
		}
	}
}

// blendPixel blends a pixel with alpha blending
func (dc *Context) blendPixel(x, y int, newPixel color.RGBA) {
	if x < 0 || x >= dc.width || y < 0 || y >= dc.height {
		return
	}

	existing := dc.im.RGBAAt(x, y)
	alpha := float64(newPixel.A) / 255.0
	invAlpha := 1.0 - alpha

	blended := color.RGBA{
		R: uint8(float64(newPixel.R)*alpha + float64(existing.R)*invAlpha),
		G: uint8(float64(newPixel.G)*alpha + float64(existing.G)*invAlpha),
		B: uint8(float64(newPixel.B)*alpha + float64(existing.B)*invAlpha),
		A: uint8(math.Min(255, float64(newPixel.A)+float64(existing.A))),
	}

	dc.im.SetRGBA(x, y, blended)
}

// blurImage applies a simple box blur to the image
func (dc *Context) blurImage(img *image.RGBA, radius float64) *image.RGBA {
	if radius <= 0 {
		return img
	}

	bounds := img.Bounds()
	blurred := image.NewRGBA(bounds)
	
	r := int(radius)
	if r < 1 {
		r = 1
	}

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			var rSum, gSum, bSum, aSum, count uint32

			// Sample surrounding pixels
			for dy := -r; dy <= r; dy++ {
				for dx := -r; dx <= r; dx++ {
					px, py := x+dx, y+dy
					if px >= bounds.Min.X && px < bounds.Max.X && py >= bounds.Min.Y && py < bounds.Max.Y {
						pixel := img.RGBAAt(px, py)
						rSum += uint32(pixel.R)
						gSum += uint32(pixel.G)
						bSum += uint32(pixel.B)
						aSum += uint32(pixel.A)
						count++
					}
				}
			}

			if count > 0 {
				blurred.SetRGBA(x, y, color.RGBA{
					R: uint8(rSum / count),
					G: uint8(gSum / count),
					B: uint8(bSum / count),
					A: uint8(aSum / count),
				})
			}
		}
	}

	return blurred
}

// Shadow-enabled drawing methods

// FillWithShadow fills the current path with shadow support
func (dc *Context) FillWithShadow() {
	dc.applyShadow(func() {
		dc.Fill()
	})
}

// StrokeWithShadow strokes the current path with shadow support
func (dc *Context) StrokeWithShadow() {
	dc.applyShadow(func() {
		dc.Stroke()
	})
}

// DrawStringWithShadow draws text with shadow support
func (dc *Context) DrawStringWithShadow(s string, x, y float64) {
	dc.applyShadow(func() {
		dc.DrawString(s, x, y)
	})
}

// DrawStringAnchoredWithShadow draws anchored text with shadow support
func (dc *Context) DrawStringAnchoredWithShadow(s string, x, y, ax, ay float64) {
	dc.applyShadow(func() {
		dc.DrawStringAnchored(s, x, y, ax, ay)
	})
}

// DrawCircleWithShadow draws a circle with shadow support
func (dc *Context) DrawCircleWithShadow(x, y, r float64) {
	dc.applyShadow(func() {
		dc.DrawCircle(x, y, r)
		dc.Fill()
	})
}

// DrawRectangleWithShadow draws a rectangle with shadow support
func (dc *Context) DrawRectangleWithShadow(x, y, w, h float64) {
	dc.applyShadow(func() {
		dc.DrawRectangle(x, y, w, h)
		dc.Fill()
	})
}

// DrawRoundedRectangleWithShadow draws a rounded rectangle with shadow support
func (dc *Context) DrawRoundedRectangleWithShadow(x, y, w, h, r float64) {
	dc.applyShadow(func() {
		dc.DrawRoundedRectangle(x, y, w, h, r)
		dc.Fill()
	})
}

// DrawEllipseWithShadow draws an ellipse with shadow support
func (dc *Context) DrawEllipseWithShadow(x, y, rx, ry float64) {
	dc.applyShadow(func() {
		dc.DrawEllipse(x, y, rx, ry)
		dc.Fill()
	})
}

// DrawImageWithShadow draws an image with shadow support
func (dc *Context) DrawImageWithShadow(im image.Image, x, y int) {
	dc.applyShadow(func() {
		dc.DrawImage(im, x, y)
	})
}

// DrawImageAnchoredWithShadow draws an anchored image with shadow support
func (dc *Context) DrawImageAnchoredWithShadow(im image.Image, x, y int, ax, ay float64) {
	dc.applyShadow(func() {
		dc.DrawImageAnchored(im, x, y, ax, ay)
	})
}

// Helper method to copy context state for shadow rendering
func (dc *Context) copyForShadow() *Context {
	shadow := &Context{
		width:         dc.width,
		height:        dc.height,
		rasterizer:    dc.rasterizer,
		im:            image.NewRGBA(dc.im.Bounds()),
		mask:          dc.mask,
		color:         dc.shadowColor,
		fillPattern:   dc.fillPattern,
		strokePattern: dc.strokePattern,
		strokePath:    dc.strokePath,
		fillPath:      dc.fillPath,
		start:         dc.start,
		current:       dc.current,
		hasCurrent:    dc.hasCurrent,
		dashes:        dc.dashes,
		dashOffset:    dc.dashOffset,
		lineWidth:     dc.lineWidth,
		lineCap:       dc.lineCap,
		lineJoin:      dc.lineJoin,
		fillRule:      dc.fillRule,
		fontFace:      dc.fontFace,
		fontHeight:    dc.fontHeight,
		matrix:        dc.matrix,
	}
	
	// Clear the shadow image
	draw.Draw(shadow.im, shadow.im.Bounds(), &image.Uniform{color.Transparent}, image.Point{}, draw.Src)
	
	return shadow
}

// Gaussian blur implementation for better shadow quality
func (dc *Context) gaussianBlur(img *image.RGBA, radius float64) *image.RGBA {
	if radius <= 0 {
		return img
	}

	// For simplicity, use box blur approximation
	// A true Gaussian blur would require more complex implementation
	return dc.blurImage(img, radius)
}
