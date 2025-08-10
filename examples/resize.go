package main

import (
	"fmt"

	"github.com/GrandpaEJ/advancegg/internal/core"
)

func main() {
	// Create or load a sample image; we'll draw a simple gradient to have content
	dc := core.NewContext(640, 360)
	for y := 0; y < 360; y++ {
		// vertical gradient
		t := float64(y) / 359.0
		dc.SetRGB(0.2*(1-t)+0.0*t, 0.5*(1-t)+0.8*t, 0.9*(1-t)+0.3*t)
		dc.DrawLine(0, float64(y), 640, float64(y))
		dc.Stroke()
	}
	// add a circle to see shapes after resizing
	dc.SetRGB(1, 1, 1)
	dc.DrawCircle(320, 180, 70)
	dc.Fill()
	core.SavePNG("resize_input.png", dc.Image())

	img := dc.Image()

	// 1) Simple resize to specific dimensions
	resized := core.ResizeImage(img, 800, 600)
	core.SavePNG("resize_800x600_bilinear.png", resized)

	// 2) Fit within 300x300, maintaining aspect ratio (no crop)
	fit := core.ResizeImageFit(img, 300, 300)
	core.SavePNG("resize_fit_300x300.png", fit)

	// 3) Fill to 300x300, maintaining aspect ratio with center crop
	fill := core.ResizeImageFill(img, 300, 300)
	core.SavePNG("resize_fill_300x300.png", fill)

	// 4) Scale by a factor (2x)
	scaled := core.ScaleImage(img, 2.0)
	core.SavePNG("resize_scale_2x.png", scaled)

	// 5) Explicit algorithm: nearest-neighbor (fast)
	nearest := core.ResizeImageWithAlgorithm(img, 300, 180, core.ResizeNearestNeighbor)
	core.SavePNG("resize_300x180_nearest.png", nearest)

	fmt.Println("Resize example images saved:")
	fmt.Println(" - resize_input.png")
	fmt.Println(" - resize_800x600_bilinear.png")
	fmt.Println(" - resize_fit_300x300.png")
	fmt.Println(" - resize_fill_300x300.png")
	fmt.Println(" - resize_scale_2x.png")
	fmt.Println(" - resize_300x180_nearest.png")
}
