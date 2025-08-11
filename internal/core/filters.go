package core

import (
	"image"
	"image/color"
	"math"
	"math/rand"
	"runtime"
)

// Filter represents an image filter function
type Filter func(img image.Image) image.Image

// ApplyFilter applies a filter to the context's current image
func (dc *Context) ApplyFilter(filter Filter) {
	dc.im = filter(dc.im).(*image.RGBA)
}

// Grayscale converts the image to grayscale
func Grayscale(img image.Image) image.Image {
	bounds := img.Bounds()
	grayImg := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			// Convert to grayscale using luminance formula
			grayValue := uint8((0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b)) / 256)
			grayImg.Set(x, y, color.RGBA{grayValue, grayValue, grayValue, uint8(a / 256)})
		}
	}
	return grayImg
}

// Invert inverts the colors of the image
func Invert(img image.Image) image.Image {
	bounds := img.Bounds()
	inverted := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			inverted.Set(x, y, color.RGBA{
				255 - uint8(r/256),
				255 - uint8(g/256),
				255 - uint8(b/256),
				uint8(a / 256),
			})
		}
	}
	return inverted
}

// Sepia applies a sepia tone effect
func Sepia(img image.Image) image.Image {
	bounds := img.Bounds()
	sepia := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			rf, gf, bf := float64(r/256), float64(g/256), float64(b/256)

			// Sepia transformation matrix
			newR := math.Min(255, 0.393*rf+0.769*gf+0.189*bf)
			newG := math.Min(255, 0.349*rf+0.686*gf+0.168*bf)
			newB := math.Min(255, 0.272*rf+0.534*gf+0.131*bf)

			sepia.Set(x, y, color.RGBA{
				uint8(newR),
				uint8(newG),
				uint8(newB),
				uint8(a / 256),
			})
		}
	}
	return sepia
}

// Brightness adjusts the brightness of the image
func Brightness(factor float64) Filter {
	return func(img image.Image) image.Image {
		bounds := img.Bounds()
		bright := image.NewRGBA(bounds)

		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				r, g, b, a := img.At(x, y).RGBA()

				newR := math.Max(0, math.Min(255, float64(r/256)*factor))
				newG := math.Max(0, math.Min(255, float64(g/256)*factor))
				newB := math.Max(0, math.Min(255, float64(b/256)*factor))

				bright.Set(x, y, color.RGBA{
					uint8(newR),
					uint8(newG),
					uint8(newB),
					uint8(a / 256),
				})
			}
		}
		return bright
	}
}

// Contrast adjusts the contrast of the image
func Contrast(factor float64) Filter {
	return func(img image.Image) image.Image {
		bounds := img.Bounds()
		contrast := image.NewRGBA(bounds)

		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				r, g, b, a := img.At(x, y).RGBA()
				rf, gf, bf := float64(r/256), float64(g/256), float64(b/256)

				// Apply contrast formula: (color - 128) * factor + 128
				newR := math.Max(0, math.Min(255, (rf-128)*factor+128))
				newG := math.Max(0, math.Min(255, (gf-128)*factor+128))
				newB := math.Max(0, math.Min(255, (bf-128)*factor+128))

				contrast.Set(x, y, color.RGBA{
					uint8(newR),
					uint8(newG),
					uint8(newB),
					uint8(a / 256),
				})
			}
		}
		return contrast
	}
}

// Blur applies a simple box blur
func Blur(radius int) Filter {
	return func(img image.Image) image.Image {
		bounds := img.Bounds()
		blurred := image.NewRGBA(bounds)

		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				var rSum, gSum, bSum, aSum, count uint32

				// Sample surrounding pixels
				for dy := -radius; dy <= radius; dy++ {
					for dx := -radius; dx <= radius; dx++ {
						px, py := x+dx, y+dy
						if px >= bounds.Min.X && px < bounds.Max.X && py >= bounds.Min.Y && py < bounds.Max.Y {
							r, g, b, a := img.At(px, py).RGBA()
							rSum += r / 256
							gSum += g / 256
							bSum += b / 256
							aSum += a / 256
							count++
						}
					}
				}

				if count > 0 {
					blurred.Set(x, y, color.RGBA{
						uint8(rSum / count),
						uint8(gSum / count),
						uint8(bSum / count),
						uint8(aSum / count),
					})
				}
			}
		}
		return blurred
	}
}

// Sharpen applies a sharpening filter
func Sharpen(img image.Image) image.Image {
	bounds := img.Bounds()
	sharpened := image.NewRGBA(bounds)

	// Sharpening kernel
	kernel := [][]float64{
		{0, -1, 0},
		{-1, 5, -1},
		{0, -1, 0},
	}

	for y := bounds.Min.Y + 1; y < bounds.Max.Y-1; y++ {
		for x := bounds.Min.X + 1; x < bounds.Max.X-1; x++ {
			var rSum, gSum, bSum float64

			for ky := 0; ky < 3; ky++ {
				for kx := 0; kx < 3; kx++ {
					px, py := x+kx-1, y+ky-1
					r, g, b, _ := img.At(px, py).RGBA()
					weight := kernel[ky][kx]
					rSum += float64(r/256) * weight
					gSum += float64(g/256) * weight
					bSum += float64(b/256) * weight
				}
			}

			_, _, _, a := img.At(x, y).RGBA()
			sharpened.Set(x, y, color.RGBA{
				uint8(math.Max(0, math.Min(255, rSum))),
				uint8(math.Max(0, math.Min(255, gSum))),
				uint8(math.Max(0, math.Min(255, bSum))),
				uint8(a / 256),
			})
		}
	}
	return sharpened
}

// Threshold applies a threshold effect
func Threshold(threshold uint8) Filter {
	return func(img image.Image) image.Image {
		bounds := img.Bounds()
		thresholded := image.NewRGBA(bounds)

		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				r, g, b, a := img.At(x, y).RGBA()
				// Convert to grayscale first
				gray := uint8((0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b)) / 256)

				var newColor uint8
				if gray > threshold {
					newColor = 255
				} else {
					newColor = 0
				}

				thresholded.Set(x, y, color.RGBA{newColor, newColor, newColor, uint8(a / 256)})
			}
		}
		return thresholded
	}
}

// Pixelate creates a pixelated effect
func Pixelate(blockSize int) Filter {
	return func(img image.Image) image.Image {
		bounds := img.Bounds()
		pixelated := image.NewRGBA(bounds)

		for y := bounds.Min.Y; y < bounds.Max.Y; y += blockSize {
			for x := bounds.Min.X; x < bounds.Max.X; x += blockSize {
				// Sample the center pixel of the block
				centerX := x + blockSize/2
				centerY := y + blockSize/2
				if centerX >= bounds.Max.X {
					centerX = bounds.Max.X - 1
				}
				if centerY >= bounds.Max.Y {
					centerY = bounds.Max.Y - 1
				}

				centerColor := img.At(centerX, centerY)

				// Fill the entire block with this color
				for by := y; by < y+blockSize && by < bounds.Max.Y; by++ {
					for bx := x; bx < x+blockSize && bx < bounds.Max.X; bx++ {
						pixelated.Set(bx, by, centerColor)
					}
				}
			}
		}
		return pixelated
	}
}

// Noise adds random noise to the image
func Noise(intensity float64) Filter {
	return func(img image.Image) image.Image {
		bounds := img.Bounds()
		noisy := image.NewRGBA(bounds)

		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				r, g, b, a := img.At(x, y).RGBA()

				// Add random noise
				noise := (rand.Float64() - 0.5) * intensity * 255
				newR := math.Max(0, math.Min(255, float64(r/256)+noise))
				newG := math.Max(0, math.Min(255, float64(g/256)+noise))
				newB := math.Max(0, math.Min(255, float64(b/256)+noise))

				noisy.Set(x, y, color.RGBA{
					uint8(newR),
					uint8(newG),
					uint8(newB),
					uint8(a / 256),
				})
			}
		}
		return noisy
	}
}

// EdgeDetection applies edge detection using Sobel operator
func EdgeDetection(img image.Image) image.Image {
	bounds := img.Bounds()
	edges := image.NewRGBA(bounds)

	// Sobel kernels
	sobelX := [][]float64{
		{-1, 0, 1},
		{-2, 0, 2},
		{-1, 0, 1},
	}
	sobelY := [][]float64{
		{-1, -2, -1},
		{0, 0, 0},
		{1, 2, 1},
	}

	for y := bounds.Min.Y + 1; y < bounds.Max.Y-1; y++ {
		for x := bounds.Min.X + 1; x < bounds.Max.X-1; x++ {
			var gx, gy float64

			for ky := 0; ky < 3; ky++ {
				for kx := 0; kx < 3; kx++ {
					px, py := x+kx-1, y+ky-1
					r, g, b, _ := img.At(px, py).RGBA()
					gray := 0.299*float64(r/256) + 0.587*float64(g/256) + 0.114*float64(b/256)

					gx += gray * sobelX[ky][kx]
					gy += gray * sobelY[ky][kx]
				}
			}

			magnitude := math.Sqrt(gx*gx + gy*gy)
			edgeValue := uint8(math.Min(255, magnitude))

			_, _, _, a := img.At(x, y).RGBA()
			edges.Set(x, y, color.RGBA{edgeValue, edgeValue, edgeValue, uint8(a / 256)})
		}
	}
	return edges
}

// Emboss applies an emboss effect
func Emboss(img image.Image) image.Image {
	bounds := img.Bounds()
	embossed := image.NewRGBA(bounds)

	// Emboss kernel
	kernel := [][]float64{
		{-2, -1, 0},
		{-1, 1, 1},
		{0, 1, 2},
	}

	for y := bounds.Min.Y + 1; y < bounds.Max.Y-1; y++ {
		for x := bounds.Min.X + 1; x < bounds.Max.X-1; x++ {
			var rSum, gSum, bSum float64

			for ky := 0; ky < 3; ky++ {
				for kx := 0; kx < 3; kx++ {
					px, py := x+kx-1, y+ky-1
					r, g, b, _ := img.At(px, py).RGBA()
					weight := kernel[ky][kx]
					rSum += float64(r/256) * weight
					gSum += float64(g/256) * weight
					bSum += float64(b/256) * weight
				}
			}

			// Add 128 to center the values
			rSum += 128
			gSum += 128
			bSum += 128

			_, _, _, a := img.At(x, y).RGBA()
			embossed.Set(x, y, color.RGBA{
				uint8(math.Max(0, math.Min(255, rSum))),
				uint8(math.Max(0, math.Min(255, gSum))),
				uint8(math.Max(0, math.Min(255, bSum))),
				uint8(a / 256),
			})
		}
	}
	return embossed
}

// Performance-optimized filters using SIMD and parallel processing

// FastGrayscale is an optimized version of Grayscale using SIMD
func FastGrayscale(img image.Image) image.Image {
	if rgba, ok := img.(*image.RGBA); ok {
		return SIMDColorTransform(rgba, func(r, g, b, a uint8) (uint8, uint8, uint8, uint8) {
			// Optimized grayscale conversion using integer arithmetic
			gray := uint8((uint32(r)*77 + uint32(g)*150 + uint32(b)*29) >> 8)
			return gray, gray, gray, a
		})
	}

	// Fallback to standard implementation
	return Grayscale(img)
}

// FastBrightness is an optimized brightness filter
func FastBrightness(factor float64) Filter {
	return func(img image.Image) image.Image {
		if rgba, ok := img.(*image.RGBA); ok {
			// Pre-calculate lookup table for performance
			lut := make([]uint8, 256)
			for i := 0; i < 256; i++ {
				val := float64(i) * factor
				if val > 255 {
					lut[i] = 255
				} else if val < 0 {
					lut[i] = 0
				} else {
					lut[i] = uint8(val)
				}
			}

			return SIMDColorTransform(rgba, func(r, g, b, a uint8) (uint8, uint8, uint8, uint8) {
				return lut[r], lut[g], lut[b], a
			})
		}

		// Fallback to standard implementation
		return Brightness(factor)(img)
	}
}

// FastContrast is an optimized contrast filter
func FastContrast(factor float64) Filter {
	return func(img image.Image) image.Image {
		if rgba, ok := img.(*image.RGBA); ok {
			// Pre-calculate lookup table
			lut := make([]uint8, 256)
			for i := 0; i < 256; i++ {
				val := (float64(i)-128)*factor + 128
				if val > 255 {
					lut[i] = 255
				} else if val < 0 {
					lut[i] = 0
				} else {
					lut[i] = uint8(val)
				}
			}

			return SIMDColorTransform(rgba, func(r, g, b, a uint8) (uint8, uint8, uint8, uint8) {
				return lut[r], lut[g], lut[b], a
			})
		}

		return Contrast(factor)(img)
	}
}

// FastBlur is an optimized blur filter using separable convolution
func FastBlur(radius int) Filter {
	return func(img image.Image) image.Image {
		if rgba, ok := img.(*image.RGBA); ok {
			return SIMDBlur(rgba, radius)
		}

		return Blur(radius)(img)
	}
}

// FastSharpen is an optimized sharpen filter
func FastSharpen(amount float64) Filter {
	return func(img image.Image) image.Image {
		if rgba, ok := img.(*image.RGBA); ok {
			// Optimized unsharp mask
			kernel := [][]float64{
				{0, -amount, 0},
				{-amount, 1 + 4*amount, -amount},
				{0, -amount, 0},
			}
			return SIMDConvolution(rgba, kernel)
		}

		return Sharpen(img)
	}
}

// FastEdgeDetection is an optimized edge detection filter
func FastEdgeDetection() Filter {
	return func(img image.Image) image.Image {
		if rgba, ok := img.(*image.RGBA); ok {
			// Sobel operator
			sobelX := [][]float64{
				{-1, 0, 1},
				{-2, 0, 2},
				{-1, 0, 1},
			}

			sobelY := [][]float64{
				{-1, -2, -1},
				{0, 0, 0},
				{1, 2, 1},
			}

			// Apply both kernels and combine
			edgesX := SIMDConvolution(rgba, sobelX)
			edgesY := SIMDConvolution(rgba, sobelY)

			return combineEdges(edgesX, edgesY)
		}

		return EdgeDetection(img)
	}
}

// combineEdges combines horizontal and vertical edge detection results
func combineEdges(edgesX, edgesY *image.RGBA) *image.RGBA {
	bounds := edgesX.Bounds()
	result := image.NewRGBA(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			px := edgesX.RGBAAt(x, y)
			py := edgesY.RGBAAt(x, y)

			// Combine using magnitude
			gx := float64(px.R)
			gy := float64(py.R)
			magnitude := uint8(math.Min(255, math.Sqrt(gx*gx+gy*gy)))

			result.SetRGBA(x, y, color.RGBA{magnitude, magnitude, magnitude, px.A})
		}
	}

	return result
}

// BatchFilter applies multiple filters efficiently
func BatchFilter(filters ...Filter) Filter {
	return func(img image.Image) image.Image {
		result := img
		for _, filter := range filters {
			result = filter(result)
		}
		return result
	}
}

// ParallelFilter applies a filter using parallel processing
func ParallelFilter(filter Filter, numWorkers int) Filter {
	return func(img image.Image) image.Image {
		if rgba, ok := img.(*image.RGBA); ok {
			bounds := rgba.Bounds()
			width := bounds.Max.X - bounds.Min.X
			height := bounds.Max.Y - bounds.Min.Y
			result := image.NewRGBA(bounds)

			if numWorkers <= 0 {
				numWorkers = runtime.NumCPU()
			}

			rowsPerWorker := height / numWorkers
			done := make(chan bool, numWorkers)

			for worker := 0; worker < numWorkers; worker++ {
				startY := worker * rowsPerWorker
				endY := startY + rowsPerWorker
				if worker == numWorkers-1 {
					endY = height
				}

				go func(startY, endY int) {
					// Create sub-image for this worker
					subBounds := image.Rect(bounds.Min.X, bounds.Min.Y+startY, bounds.Max.X, bounds.Min.Y+endY)
					subImg := rgba.SubImage(subBounds)

					// Apply filter to sub-image
					filtered := filter(subImg)

					// Copy result back
					if filteredRGBA, ok := filtered.(*image.RGBA); ok {
						for y := startY; y < endY; y++ {
							for x := 0; x < width; x++ {
								pixel := filteredRGBA.RGBAAt(x, y-startY)
								result.SetRGBA(x, y, pixel)
							}
						}
					}

					done <- true
				}(startY, endY)
			}

			// Wait for all workers
			for i := 0; i < numWorkers; i++ {
				<-done
			}

			return result
		}

		return filter(img)
	}
}

// Posterize reduces the number of colors
func Posterize(levels int) Filter {
	return func(img image.Image) image.Image {
		bounds := img.Bounds()
		posterized := image.NewRGBA(bounds)

		step := 255.0 / float64(levels-1)

		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				r, g, b, a := img.At(x, y).RGBA()

				newR := math.Round(float64(r/256)/step) * step
				newG := math.Round(float64(g/256)/step) * step
				newB := math.Round(float64(b/256)/step) * step

				posterized.Set(x, y, color.RGBA{
					uint8(math.Min(255, newR)),
					uint8(math.Min(255, newG)),
					uint8(math.Min(255, newB)),
					uint8(a / 256),
				})
			}
		}
		return posterized
	}
}

// Vignette applies a vignette effect
func Vignette(strength float64) Filter {
	return func(img image.Image) image.Image {
		bounds := img.Bounds()
		vignetted := image.NewRGBA(bounds)

		centerX := float64(bounds.Max.X-bounds.Min.X) / 2
		centerY := float64(bounds.Max.Y-bounds.Min.Y) / 2
		maxDistance := math.Sqrt(centerX*centerX + centerY*centerY)

		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				r, g, b, a := img.At(x, y).RGBA()

				// Calculate distance from center
				dx := float64(x) - centerX
				dy := float64(y) - centerY
				distance := math.Sqrt(dx*dx + dy*dy)

				// Calculate vignette factor
				factor := 1.0 - (distance/maxDistance)*strength
				factor = math.Max(0, factor)

				vignetted.Set(x, y, color.RGBA{
					uint8(float64(r/256) * factor),
					uint8(float64(g/256) * factor),
					uint8(float64(b/256) * factor),
					uint8(a / 256),
				})
			}
		}
		return vignetted
	}
}
