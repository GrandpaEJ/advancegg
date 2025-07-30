package main

import (
	"log"

	"github.com/GrandpaEJ/advancegg"
)

func main() {
	im, err := advancegg.LoadImage("examples/baboon.png")
	if err != nil {
		log.Fatal(err)
	}

	dc := advancegg.NewContext(512, 512)
	dc.DrawRoundedRectangle(0, 0, 512, 512, 64)
	dc.Clip()
	dc.DrawImage(im, 0, 0)
	dc.SavePNG("out.png")
}
