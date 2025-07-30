package core

import (
	"image"
	"image/color"
)

// ImageData represents pixel data that can be manipulated directly
type ImageData struct {
	Data   []uint8 // RGBA data: [R, G, B, A, R, G, B, A, ...]
	Width  int
	Height int
}

// NewImageData creates a new ImageData with the specified dimensions
func NewImageData(width, height int) *ImageData {
	return &ImageData{
		Data:   make([]uint8, width*height*4),
		Width:  width,
		Height: height,
	}
}

// NewImageDataFromImage creates ImageData from an existing image
func NewImageDataFromImage(img image.Image) *ImageData {
	bounds := img.Bounds()
	width := bounds.Max.X - bounds.Min.X
	height := bounds.Max.Y - bounds.Min.Y
	
	data := make([]uint8, width*height*4)
	
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r, g, b, a := img.At(bounds.Min.X+x, bounds.Min.Y+y).RGBA()
			
			index := (y*width + x) * 4
			data[index] = uint8(r >> 8)     // R
			data[index+1] = uint8(g >> 8)   // G
			data[index+2] = uint8(b >> 8)   // B
			data[index+3] = uint8(a >> 8)   // A
		}
	}
	
	return &ImageData{
		Data:   data,
		Width:  width,
		Height: height,
	}
}

// GetPixel returns the RGBA values at the specified coordinates
func (id *ImageData) GetPixel(x, y int) (r, g, b, a uint8) {
	if x < 0 || x >= id.Width || y < 0 || y >= id.Height {
		return 0, 0, 0, 0
	}
	
	index := (y*id.Width + x) * 4
	return id.Data[index], id.Data[index+1], id.Data[index+2], id.Data[index+3]
}

// SetPixel sets the RGBA values at the specified coordinates
func (id *ImageData) SetPixel(x, y int, r, g, b, a uint8) {
	if x < 0 || x >= id.Width || y < 0 || y >= id.Height {
		return
	}
	
	index := (y*id.Width + x) * 4
	id.Data[index] = r
	id.Data[index+1] = g
	id.Data[index+2] = b
	id.Data[index+3] = a
}

// GetPixelColor returns the color at the specified coordinates
func (id *ImageData) GetPixelColor(x, y int) color.RGBA {
	r, g, b, a := id.GetPixel(x, y)
	return color.RGBA{r, g, b, a}
}

// SetPixelColor sets the color at the specified coordinates
func (id *ImageData) SetPixelColor(x, y int, c color.Color) {
	r, g, b, a := c.RGBA()
	id.SetPixel(x, y, uint8(r>>8), uint8(g>>8), uint8(b>>8), uint8(a>>8))
}

// ToImage converts ImageData to a standard Go image
func (id *ImageData) ToImage() *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, id.Width, id.Height))
	
	for y := 0; y < id.Height; y++ {
		for x := 0; x < id.Width; x++ {
			r, g, b, a := id.GetPixel(x, y)
			img.Set(x, y, color.RGBA{r, g, b, a})
		}
	}
	
	return img
}

// Clone creates a copy of the ImageData
func (id *ImageData) Clone() *ImageData {
	newData := make([]uint8, len(id.Data))
	copy(newData, id.Data)
	
	return &ImageData{
		Data:   newData,
		Width:  id.Width,
		Height: id.Height,
	}
}

// Fill fills the entire ImageData with the specified color
func (id *ImageData) Fill(r, g, b, a uint8) {
	for i := 0; i < len(id.Data); i += 4 {
		id.Data[i] = r
		id.Data[i+1] = g
		id.Data[i+2] = b
		id.Data[i+3] = a
	}
}

// FillRect fills a rectangular region with the specified color
func (id *ImageData) FillRect(x, y, width, height int, r, g, b, a uint8) {
	for dy := 0; dy < height; dy++ {
		for dx := 0; dx < width; dx++ {
			id.SetPixel(x+dx, y+dy, r, g, b, a)
		}
	}
}

// CopyFrom copies pixel data from another ImageData
func (id *ImageData) CopyFrom(src *ImageData, srcX, srcY, srcWidth, srcHeight, dstX, dstY int) {
	for dy := 0; dy < srcHeight; dy++ {
		for dx := 0; dx < srcWidth; dx++ {
			if srcX+dx >= 0 && srcX+dx < src.Width && srcY+dy >= 0 && srcY+dy < src.Height {
				r, g, b, a := src.GetPixel(srcX+dx, srcY+dy)
				id.SetPixel(dstX+dx, dstY+dy, r, g, b, a)
			}
		}
	}
}

// ApplyKernel applies a convolution kernel to the ImageData
func (id *ImageData) ApplyKernel(kernel [][]float64) *ImageData {
	result := id.Clone()
	kernelSize := len(kernel)
	offset := kernelSize / 2
	
	for y := offset; y < id.Height-offset; y++ {
		for x := offset; x < id.Width-offset; x++ {
			var rSum, gSum, bSum float64
			
			for ky := 0; ky < kernelSize; ky++ {
				for kx := 0; kx < kernelSize; kx++ {
					px := x + kx - offset
					py := y + ky - offset
					
					r, g, b, _ := id.GetPixel(px, py)
					weight := kernel[ky][kx]
					
					rSum += float64(r) * weight
					gSum += float64(g) * weight
					bSum += float64(b) * weight
				}
			}
			
			// Clamp values
			rSum = clamp(rSum, 0, 255)
			gSum = clamp(gSum, 0, 255)
			bSum = clamp(bSum, 0, 255)
			
			_, _, _, a := id.GetPixel(x, y)
			result.SetPixel(x, y, uint8(rSum), uint8(gSum), uint8(bSum), a)
		}
	}
	
	return result
}

// GetSubImageData extracts a rectangular region as new ImageData
func (id *ImageData) GetSubImageData(x, y, width, height int) *ImageData {
	subData := NewImageData(width, height)
	
	for dy := 0; dy < height; dy++ {
		for dx := 0; dx < width; dx++ {
			if x+dx >= 0 && x+dx < id.Width && y+dy >= 0 && y+dy < id.Height {
				r, g, b, a := id.GetPixel(x+dx, y+dy)
				subData.SetPixel(dx, dy, r, g, b, a)
			}
		}
	}
	
	return subData
}

// Resize creates a new ImageData with different dimensions using nearest neighbor
func (id *ImageData) Resize(newWidth, newHeight int) *ImageData {
	result := NewImageData(newWidth, newHeight)
	
	xRatio := float64(id.Width) / float64(newWidth)
	yRatio := float64(id.Height) / float64(newHeight)
	
	for y := 0; y < newHeight; y++ {
		for x := 0; x < newWidth; x++ {
			srcX := int(float64(x) * xRatio)
			srcY := int(float64(y) * yRatio)
			
			r, g, b, a := id.GetPixel(srcX, srcY)
			result.SetPixel(x, y, r, g, b, a)
		}
	}
	
	return result
}

// FlipHorizontal flips the ImageData horizontally
func (id *ImageData) FlipHorizontal() *ImageData {
	result := NewImageData(id.Width, id.Height)
	
	for y := 0; y < id.Height; y++ {
		for x := 0; x < id.Width; x++ {
			r, g, b, a := id.GetPixel(x, y)
			result.SetPixel(id.Width-1-x, y, r, g, b, a)
		}
	}
	
	return result
}

// FlipVertical flips the ImageData vertically
func (id *ImageData) FlipVertical() *ImageData {
	result := NewImageData(id.Width, id.Height)
	
	for y := 0; y < id.Height; y++ {
		for x := 0; x < id.Width; x++ {
			r, g, b, a := id.GetPixel(x, y)
			result.SetPixel(x, id.Height-1-y, r, g, b, a)
		}
	}
	
	return result
}

// Rotate90 rotates the ImageData 90 degrees clockwise
func (id *ImageData) Rotate90() *ImageData {
	result := NewImageData(id.Height, id.Width)
	
	for y := 0; y < id.Height; y++ {
		for x := 0; x < id.Width; x++ {
			r, g, b, a := id.GetPixel(x, y)
			result.SetPixel(id.Height-1-y, x, r, g, b, a)
		}
	}
	
	return result
}

// Helper function to clamp values
func clamp(value, min, max float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}
