package main

import "github.com/GrandpaEJ/advancegg"

func main() {
	dc := advancegg.NewContext(1000, 1000)
	dc.DrawCircle(500, 500, 400)
	dc.SetRGB(0, 0, 0)
	dc.Fill()
	dc.SavePNG("out.png")
}
