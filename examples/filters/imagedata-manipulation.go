package main

import (
	"fmt"
	"math"
	"math/rand"
	
	"github.com/GrandpaEJ/advancegg"
)

func main() {
	// Create a sample image
	dc := advancegg.NewContext(400, 300)
	createSampleImage(dc)
	dc.SavePNG("images/filters/imagedata-original.png")
	
	// Get ImageData for manipulation
	imageData := dc.GetImageData()
	
	// Demonstrate various ImageData operations
	demonstratePixelManipulation(imageData)
	demonstrateImageTransforms(imageData)
	demonstrateKernelOperations(imageData)
	demonstrateImageEffects(imageData)
	
	fmt.Println("ImageData manipulation examples completed!")
}

func createSampleImage(dc *advancegg.Context) {
	// Create a colorful test image
	dc.SetRGB(0.2, 0.3, 0.8)
	dc.Clear()
	
	// Add some shapes
	dc.SetRGB(1, 0, 0)
	dc.DrawCircle(100, 75, 40)
	dc.Fill()
	
	dc.SetRGB(0, 1, 0)
	dc.DrawRectangle(200, 50, 80, 50)
	dc.Fill()
	
	dc.SetRGB(1, 1, 0)
	dc.DrawCircle(300, 75, 30)
	dc.Fill()
	
	// Add gradient
	for x := 0; x < 400; x++ {
		for y := 200; y < 250; y++ {
			r := float64(x) / 400.0
			g := 0.5
			b := float64(y-200) / 50.0
			dc.SetRGB(r, g, b)
			dc.SetPixel(x, y)
		}
	}
	
	dc.SetRGB(1, 1, 1)
	dc.DrawString("ImageData Test", 50, 280)
}

func demonstratePixelManipulation(imageData *advancegg.ImageData) {
	fmt.Println("Demonstrating pixel manipulation...")
	
	// Clone the original for manipulation
	manipulated := imageData.Clone()
	
	// Direct pixel manipulation - create a checkerboard pattern overlay
	for y := 0; y < manipulated.Height; y += 10 {
		for x := 0; x < manipulated.Width; x += 10 {
			if (x/10+y/10)%2 == 0 {
				// Get original pixel
				r, g, b, a := manipulated.GetPixel(x, y)
				// Blend with white
				newR := uint8((int(r) + 255) / 2)
				newG := uint8((int(g) + 255) / 2)
				newB := uint8((int(b) + 255) / 2)
				
				// Fill 10x10 block
				for dy := 0; dy < 10 && y+dy < manipulated.Height; dy++ {
					for dx := 0; dx < 10 && x+dx < manipulated.Width; dx++ {
						manipulated.SetPixel(x+dx, y+dy, newR, newG, newB, a)
					}
				}
			}
		}
	}
	
	// Save the manipulated image
	dc := advancegg.NewContext(manipulated.Width, manipulated.Height)
	dc.PutImageData(manipulated)
	dc.SavePNG("images/filters/imagedata-checkerboard.png")
	fmt.Println("Saved checkerboard overlay as imagedata-checkerboard.png")
}

func demonstrateImageTransforms(imageData *advancegg.ImageData) {
	fmt.Println("Demonstrating image transforms...")
	
	// Flip horizontal
	flippedH := imageData.FlipHorizontal()
	dc1 := advancegg.NewContext(flippedH.Width, flippedH.Height)
	dc1.PutImageData(flippedH)
	dc1.SavePNG("images/filters/imagedata-flip-horizontal.png")
	
	// Flip vertical
	flippedV := imageData.FlipVertical()
	dc2 := advancegg.NewContext(flippedV.Width, flippedV.Height)
	dc2.PutImageData(flippedV)
	dc2.SavePNG("images/filters/imagedata-flip-vertical.png")
	
	// Rotate 90 degrees
	rotated := imageData.Rotate90()
	dc3 := advancegg.NewContext(rotated.Width, rotated.Height)
	dc3.PutImageData(rotated)
	dc3.SavePNG("images/filters/imagedata-rotate90.png")
	
	// Resize (scale down)
	resized := imageData.Resize(200, 150)
	dc4 := advancegg.NewContext(resized.Width, resized.Height)
	dc4.PutImageData(resized)
	dc4.SavePNG("images/filters/imagedata-resized.png")
	
	fmt.Println("Saved transform examples")
}

func demonstrateKernelOperations(imageData *advancegg.ImageData) {
	fmt.Println("Demonstrating kernel operations...")
	
	// Blur kernel
	blurKernel := [][]float64{
		{1.0/9, 1.0/9, 1.0/9},
		{1.0/9, 1.0/9, 1.0/9},
		{1.0/9, 1.0/9, 1.0/9},
	}
	
	blurred := imageData.ApplyKernel(blurKernel)
	dc1 := advancegg.NewContext(blurred.Width, blurred.Height)
	dc1.PutImageData(blurred)
	dc1.SavePNG("images/filters/imagedata-blur-kernel.png")
	
	// Sharpen kernel
	sharpenKernel := [][]float64{
		{0, -1, 0},
		{-1, 5, -1},
		{0, -1, 0},
	}
	
	sharpened := imageData.ApplyKernel(sharpenKernel)
	dc2 := advancegg.NewContext(sharpened.Width, sharpened.Height)
	dc2.PutImageData(sharpened)
	dc2.SavePNG("images/filters/imagedata-sharpen-kernel.png")
	
	// Edge detection kernel
	edgeKernel := [][]float64{
		{-1, -1, -1},
		{-1, 8, -1},
		{-1, -1, -1},
	}
	
	edges := imageData.ApplyKernel(edgeKernel)
	dc3 := advancegg.NewContext(edges.Width, edges.Height)
	dc3.PutImageData(edges)
	dc3.SavePNG("images/filters/imagedata-edge-kernel.png")
	
	fmt.Println("Saved kernel operation examples")
}

func demonstrateImageEffects(imageData *advancegg.ImageData) {
	fmt.Println("Demonstrating image effects...")
	
	// Create a noise effect
	noisy := imageData.Clone()
	for y := 0; y < noisy.Height; y++ {
		for x := 0; x < noisy.Width; x++ {
			r, g, b, a := noisy.GetPixel(x, y)
			
			// Add random noise
			noise := int((rand.Float64() - 0.5) * 50)
			newR := clampUint8(int(r) + noise)
			newG := clampUint8(int(g) + noise)
			newB := clampUint8(int(b) + noise)
			
			noisy.SetPixel(x, y, newR, newG, newB, a)
		}
	}
	
	dc1 := advancegg.NewContext(noisy.Width, noisy.Height)
	dc1.PutImageData(noisy)
	dc1.SavePNG("images/filters/imagedata-noise.png")
	
	// Create a wave distortion effect
	waved := imageData.Clone()
	for y := 0; y < waved.Height; y++ {
		for x := 0; x < waved.Width; x++ {
			// Calculate wave displacement
			waveX := int(math.Sin(float64(y)*0.1) * 10)
			waveY := int(math.Cos(float64(x)*0.1) * 5)
			
			srcX := x - waveX
			srcY := y - waveY
			
			if srcX >= 0 && srcX < imageData.Width && srcY >= 0 && srcY < imageData.Height {
				r, g, b, a := imageData.GetPixel(srcX, srcY)
				waved.SetPixel(x, y, r, g, b, a)
			}
		}
	}
	
	dc2 := advancegg.NewContext(waved.Width, waved.Height)
	dc2.PutImageData(waved)
	dc2.SavePNG("images/filters/imagedata-wave.png")
	
	// Create a mosaic effect
	mosaic := imageData.Clone()
	blockSize := 8
	for y := 0; y < mosaic.Height; y += blockSize {
		for x := 0; x < mosaic.Width; x += blockSize {
			// Get average color of the block
			var rSum, gSum, bSum, count int
			for dy := 0; dy < blockSize && y+dy < mosaic.Height; dy++ {
				for dx := 0; dx < blockSize && x+dx < mosaic.Width; dx++ {
					r, g, b, _ := imageData.GetPixel(x+dx, y+dy)
					rSum += int(r)
					gSum += int(g)
					bSum += int(b)
					count++
				}
			}
			
			if count > 0 {
				avgR := uint8(rSum / count)
				avgG := uint8(gSum / count)
				avgB := uint8(bSum / count)
				
				// Fill the block with average color
				for dy := 0; dy < blockSize && y+dy < mosaic.Height; dy++ {
					for dx := 0; dx < blockSize && x+dx < mosaic.Width; dx++ {
						_, _, _, a := imageData.GetPixel(x+dx, y+dy)
						mosaic.SetPixel(x+dx, y+dy, avgR, avgG, avgB, a)
					}
				}
			}
		}
	}
	
	dc3 := advancegg.NewContext(mosaic.Width, mosaic.Height)
	dc3.PutImageData(mosaic)
	dc3.SavePNG("images/filters/imagedata-mosaic.png")
	
	// Create a composite showing all effects
	createEffectsComposite(imageData, noisy, waved, mosaic)
	
	fmt.Println("Saved image effects examples")
}

func createEffectsComposite(original, noisy, waved, mosaic *advancegg.ImageData) {
	// Create a 2x2 grid showing all effects
	compositeWidth := original.Width * 2
	compositeHeight := original.Height * 2
	
	composite := advancegg.NewImageData(compositeWidth, compositeHeight)
	composite.Fill(255, 255, 255, 255) // White background
	
	// Copy images to grid positions
	composite.CopyFrom(original, 0, 0, original.Width, original.Height, 0, 0)
	composite.CopyFrom(noisy, 0, 0, noisy.Width, noisy.Height, original.Width, 0)
	composite.CopyFrom(waved, 0, 0, waved.Width, waved.Height, 0, original.Height)
	composite.CopyFrom(mosaic, 0, 0, mosaic.Width, mosaic.Height, original.Width, original.Height)
	
	dc := advancegg.NewContext(composite.Width, composite.Height)
	dc.PutImageData(composite)
	dc.SavePNG("images/filters/imagedata-effects-composite.png")
	fmt.Println("Saved effects composite as imagedata-effects-composite.png")
}

func clampUint8(value int) uint8 {
	if value < 0 {
		return 0
	}
	if value > 255 {
		return 255
	}
	return uint8(value)
}
