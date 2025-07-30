package main

import "github.com/GrandpaEJ/advancegg"

func main() {
	const NX = 4
	const NY = 3
	im, err := advancegg.LoadPNG("examples/gopher.png")
	if err != nil {
		panic(err)
	}
	w := im.Bounds().Size().X
	h := im.Bounds().Size().Y
	dc := advancegg.NewContext(w*NX, h*NY)
	for y := 0; y < NY; y++ {
		for x := 0; x < NX; x++ {
			dc.DrawImage(im, x*w, y*h)
		}
	}
	dc.SavePNG("out.png")
}
