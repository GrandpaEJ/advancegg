package core

import (
	"image/color"
	"testing"
)

func TestGradientWithAlpha(t *testing.T) {
	dc := NewContext(200, 100)

	grad := NewLinearGradient(0, 0, 180, 0)
	grad.AddColorStop(0, color.RGBA{255, 0, 0, 200})
	grad.AddColorStop(1, color.RGBA{0, 0, 255, 100})

	dc.SetFillStyle(grad)
	dc.DrawRectangle(20, 20, 160, 60)
	dc.Fill()

	dc.SavePNG("out.png")
}
