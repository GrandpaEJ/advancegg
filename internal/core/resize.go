package core

import (
	"image"
	"image/color"
	"math"
)

// ResizeAlgorithm specifies the interpolation algorithm to use
// when resizing images.
type ResizeAlgorithm int

const (
	// ResizeNearestNeighbor uses nearest-neighbor sampling (fastest, lower quality)
	ResizeNearestNeighbor ResizeAlgorithm = iota
	// ResizeBilinear uses bilinear interpolation (good balance of quality/speed)
	ResizeBilinear
	// ResizeBicubic uses bicubic interpolation (higher quality, slower)
	ResizeBicubic
	// ResizeLanczos uses Lanczos resampling (highest quality, slowest)
	ResizeLanczos
)

// ResizeImage resizes an image to the specified width and height using
// a default algorithm (bilinear) that balances quality and performance.
func ResizeImage(img image.Image, newWidth, newHeight int) image.Image {
	return ResizeImageWithAlgorithm(img, newWidth, newHeight, ResizeBilinear)
}

// ScaleImage scales an image by a floating-point factor. For example, factor=2.0
// doubles the dimensions; factor=0.5 halves them.
func ScaleImage(img image.Image, factor float64) image.Image {
	if img == nil || factor <= 0 {
		return img
	}
	b := img.Bounds()
	ow, oh := b.Dx(), b.Dy()
	if ow == 0 || oh == 0 {
		return img
	}
	nw := int(math.Max(1, math.Round(float64(ow)*factor)))
	nh := int(math.Max(1, math.Round(float64(oh)*factor)))
	return ResizeImage(img, nw, nh)
}

// ResizeImageFit resizes an image to fit entirely within the given bounding box
// while maintaining aspect ratio. The returned image will be <= maxWidth x maxHeight.
func ResizeImageFit(img image.Image, maxWidth, maxHeight int) image.Image {
	if img == nil || maxWidth <= 0 || maxHeight <= 0 {
		return img
	}
	b := img.Bounds()
	ow, oh := b.Dx(), b.Dy()
	if ow == 0 || oh == 0 {
		return img
	}
	s := math.Min(float64(maxWidth)/float64(ow), float64(maxHeight)/float64(oh))
	if s <= 0 {
		s = 1
	}
	nw := int(math.Max(1, math.Floor(float64(ow)*s)))
	nh := int(math.Max(1, math.Floor(float64(oh)*s)))
	return ResizeImage(img, nw, nh)
}

// ResizeImageFill scales the image to completely fill the target width/height
// while maintaining aspect ratio, then center-crops to the exact size. This may crop edges.
func ResizeImageFill(img image.Image, width, height int) image.Image {
	if img == nil || width <= 0 || height <= 0 {
		return img
	}
	b := img.Bounds()
	ow, oh := b.Dx(), b.Dy()
	if ow == 0 || oh == 0 {
		return img
	}
	// Scale up so the smaller dimension fills the target
	s := math.Max(float64(width)/float64(ow), float64(height)/float64(oh))
	nw := int(math.Max(1, math.Ceil(float64(ow)*s)))
	nh := int(math.Max(1, math.Ceil(float64(oh)*s)))
	scaled := ResizeImage(img, nw, nh)
	// Center-crop to desired output size
	return cropCenter(imageToRGBA(scaled), width, height)
}

// ResizeImageWithAlgorithm resizes an image using the specified algorithm.
func ResizeImageWithAlgorithm(img image.Image, newWidth, newHeight int, algo ResizeAlgorithm) image.Image {
	if img == nil || newWidth <= 0 || newHeight <= 0 {
		return img
	}
	// Convert to RGBA for fast pixel access
	rgba := imageToRGBA(img)
	switch algo {
	case ResizeNearestNeighbor:
		return SIMDResize(rgba, newWidth, newHeight)
	case ResizeBilinear:
		return resizeBilinearRGBA(rgba, newWidth, newHeight)
	case ResizeBicubic:
		// TODO: implement true bicubic; fallback to bilinear for now
		return resizeBilinearRGBA(rgba, newWidth, newHeight)
	case ResizeLanczos:
		// TODO: implement true Lanczos; fallback to bilinear for now
		return resizeBilinearRGBA(rgba, newWidth, newHeight)
	default:
		return resizeBilinearRGBA(rgba, newWidth, newHeight)
	}
}

// resizeBilinearRGBA performs bilinear interpolation on an RGBA image.
func resizeBilinearRGBA(src *image.RGBA, newW, newH int) *image.RGBA {
	if newW <= 0 || newH <= 0 {
		return src
	}
	b := src.Bounds()
	ow := b.Dx()
	oh := b.Dy()
	if ow == 0 || oh == 0 {
		return src
	}
	dst := image.NewRGBA(image.Rect(0, 0, newW, newH))

	xRatio := float64(ow-1) / float64(newW)
	yRatio := float64(oh-1) / float64(newH)
	for y := 0; y < newH; y++ {
		fy := yRatio * float64(y)
		sy := int(math.Floor(fy))
		dy := fy - float64(sy)
		if sy >= oh-1 {
			sy = oh - 2
			dy = 1
		}
		for x := 0; x < newW; x++ {
			fx := xRatio * float64(x)
			sx := int(math.Floor(fx))
			dx := fx - float64(sx)
			if sx >= ow-1 {
				sx = ow - 2
				dx = 1
			}

			c00 := src.RGBAAt(b.Min.X+sx, b.Min.Y+sy)
			c10 := src.RGBAAt(b.Min.X+sx+1, b.Min.Y+sy)
			c01 := src.RGBAAt(b.Min.X+sx, b.Min.Y+sy+1)
			c11 := src.RGBAAt(b.Min.X+sx+1, b.Min.Y+sy+1)

			w00 := (1 - dx) * (1 - dy)
			w10 := dx * (1 - dy)
			w01 := (1 - dx) * dy
			w11 := dx * dy

			r := float64(c00.R)*w00 + float64(c10.R)*w10 + float64(c01.R)*w01 + float64(c11.R)*w11
			g := float64(c00.G)*w00 + float64(c10.G)*w10 + float64(c01.G)*w01 + float64(c11.G)*w11
			b := float64(c00.B)*w00 + float64(c10.B)*w10 + float64(c01.B)*w01 + float64(c11.B)*w11
			a := float64(c00.A)*w00 + float64(c10.A)*w10 + float64(c01.A)*w01 + float64(c11.A)*w11

			dst.SetRGBA(x, y, colorFromFloats(r, g, b, a))
		}
	}
	return dst
}

// cropCenter crops the image to target width/height centered.
func cropCenter(src *image.RGBA, width, height int) *image.RGBA {
	b := src.Bounds()
	ow, oh := b.Dx(), b.Dy()
	if width >= ow && height >= oh {
		return src
	}
	x := (ow - width) / 2
	y := (oh - height) / 2
	if x < 0 {
		x = 0
	}
	if y < 0 {
		y = 0
	}
	rect := image.Rect(0, 0, width, height)
	dst := image.NewRGBA(rect)
	for j := 0; j < height; j++ {
		for i := 0; i < width; i++ {
			dst.SetRGBA(i, j, src.RGBAAt(b.Min.X+x+i, b.Min.Y+y+j))
		}
	}
	return dst
}

func colorFromFloats(r, g, b, a float64) (c color.RGBA) {
	c.R = uint8(math.Max(0, math.Min(255, math.Round(r))))
	c.G = uint8(math.Max(0, math.Min(255, math.Round(g))))
	c.B = uint8(math.Max(0, math.Min(255, math.Round(b))))
	c.A = uint8(math.Max(0, math.Min(255, math.Round(a))))
	return
}
