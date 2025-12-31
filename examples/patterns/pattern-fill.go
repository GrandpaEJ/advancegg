package main

import "github.com/GrandpaEJ/advancegg"

func main() {
	im, err := advancegg.LoadPNG("examples/baboon.png")
	if err != nil {
		panic(err)
	}
	pattern := advancegg.NewSurfacePattern(im, advancegg.RepeatBoth)
	dc := advancegg.NewContext(600, 600)
	dc.MoveTo(20, 20)
	dc.LineTo(590, 20)
	dc.LineTo(590, 590)
	dc.LineTo(20, 590)
	dc.ClosePath()
	dc.SetFillStyle(pattern)
	dc.Fill()
	dc.SavePNG("images/patterns/out.png")
}
