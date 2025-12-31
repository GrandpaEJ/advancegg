package main

import (
	"fmt"
	"image/color"
	"path/filepath"

	"github.com/GrandpaEJ/advancegg"
)

func main() {
	modes := []struct {
		name string
		mode advancegg.BlendMode
	}{
		// Standard
		{"Normal", advancegg.BlendModeNormal},
		{"Multiply", advancegg.BlendModeMultiply},
		{"Screen", advancegg.BlendModeScreen},
		{"Overlay", advancegg.BlendModeOverlay},
		{"Darken", advancegg.BlendModeDarken},
		{"Lighten", advancegg.BlendModeLighten},
		{"ColorDodge", advancegg.BlendModeColorDodge},
		{"ColorBurn", advancegg.BlendModeColorBurn},
		{"HardLight", advancegg.BlendModeHardLight},
		{"SoftLight", advancegg.BlendModeSoftLight},
		{"Difference", advancegg.BlendModeDifference},
		{"Exclusion", advancegg.BlendModeExclusion},
		{"Hue", advancegg.BlendModeHue},
		{"Saturation", advancegg.BlendModeSaturation},
		{"Color", advancegg.BlendModeColor},
		{"Luminosity", advancegg.BlendModeLuminosity},

		// Porter-Duff
		{"Clear", advancegg.BlendModeClear},
		{"Source", advancegg.BlendModeSource},
		{"Dest", advancegg.BlendModeDest},
		{"SrcOver", advancegg.BlendModeSrcOver},
		{"DstOver", advancegg.BlendModeDstOver},
		{"SrcIn", advancegg.BlendModeSrcIn},
		{"DstIn", advancegg.BlendModeDstIn},
		{"SrcOut", advancegg.BlendModeSrcOut},
		{"DstOut", advancegg.BlendModeDstOut},
		{"SrcAtop", advancegg.BlendModeSrcAtop},
		{"DstAtop", advancegg.BlendModeDstAtop},
		{"Xor", advancegg.BlendModeXor},
		{"Add", advancegg.BlendModeAdd},
	}

	for _, m := range modes {
		// Create canvas
		width, height := 300, 300

		// Use LayerManager directly for blending features
		lm := advancegg.NewLayerManager(width, height)
		lm.Background = color.Transparent // Transparent background for Porter-Duff checks

		// 1. Destination Layer (Blue Square)
		// We make it the "base" layer on top of transparent background
		dstLayer := lm.AddLayer("Dest")

		// Draw Square on Dest Layer
		// Manually draw to the layer image using a temporary context
		dstCtx := advancegg.NewContextForRGBA(dstLayer.Image)
		dstCtx.SetRGB(0, 0, 1) // Blue
		dstCtx.DrawRectangle(50, 50, 150, 150)
		dstCtx.Fill()

		// 2. Source Layer (Red Circle)
		srcLayer := lm.AddLayer("Source")
		srcLayer.SetBlendMode(m.mode)

		srcCtx := advancegg.NewContextForRGBA(srcLayer.Image)
		srcCtx.SetRGB(1, 0, 0) // Red
		srcCtx.DrawCircle(200, 200, 80)
		srcCtx.Fill()

		// Composite
		result := lm.Composite()

		// Save
		outCtx := advancegg.NewContextForRGBA(result)
		filename := filepath.Join("images", fmt.Sprintf("blend_%s.png", m.name))
		if err := outCtx.SavePNG(filename); err != nil {
			fmt.Printf("Failed to save %s: %v\n", filename, err)
		} else {
			fmt.Printf("Generated %s\n", filename)
		}
	}
}
