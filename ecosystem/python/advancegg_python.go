package main

import "C"
import (
	"image/color"
	"unsafe"
	
	"github.com/GrandpaEJ/advancegg"
)

// Global context storage
var contexts = make(map[int]*advancegg.Context)
var nextID = 1

//export create_context
func create_context(width, height int) int {
	dc := advancegg.NewContext(width, height)
	id := nextID
	contexts[id] = dc
	nextID++
	return id
}

//export destroy_context
func destroy_context(id int) {
	delete(contexts, id)
}

//export set_rgb
func set_rgb(id int, r, g, b float64) {
	if dc, ok := contexts[id]; ok {
		dc.SetRGB(r, g, b)
	}
}

//export set_rgba
func set_rgba(id int, r, g, b, a float64) {
	if dc, ok := contexts[id]; ok {
		dc.SetRGBA(r, g, b, a)
	}
}

//export clear
func clear(id int) {
	if dc, ok := contexts[id]; ok {
		dc.Clear()
	}
}

//export draw_rectangle
func draw_rectangle(id int, x, y, width, height float64) {
	if dc, ok := contexts[id]; ok {
		dc.DrawRectangle(x, y, width, height)
	}
}

//export draw_rounded_rectangle
func draw_rounded_rectangle(id int, x, y, width, height, radius float64) {
	if dc, ok := contexts[id]; ok {
		dc.DrawRoundedRectangle(x, y, width, height, radius)
	}
}

//export draw_circle
func draw_circle(id int, x, y, radius float64) {
	if dc, ok := contexts[id]; ok {
		dc.DrawCircle(x, y, radius)
	}
}

//export draw_ellipse
func draw_ellipse(id int, x, y, rx, ry float64) {
	if dc, ok := contexts[id]; ok {
		dc.DrawEllipse(x, y, rx, ry)
	}
}

//export draw_line
func draw_line(id int, x1, y1, x2, y2 float64) {
	if dc, ok := contexts[id]; ok {
		dc.DrawLine(x1, y1, x2, y2)
	}
}

//export move_to
func move_to(id int, x, y float64) {
	if dc, ok := contexts[id]; ok {
		dc.MoveTo(x, y)
	}
}

//export line_to
func line_to(id int, x, y float64) {
	if dc, ok := contexts[id]; ok {
		dc.LineTo(x, y)
	}
}

//export curve_to
func curve_to(id int, cp1x, cp1y, cp2x, cp2y, x, y float64) {
	if dc, ok := contexts[id]; ok {
		dc.CubicTo(cp1x, cp1y, cp2x, cp2y, x, y)
	}
}

//export close_path
func close_path(id int) {
	if dc, ok := contexts[id]; ok {
		dc.ClosePath()
	}
}

//export fill
func fill(id int) {
	if dc, ok := contexts[id]; ok {
		dc.Fill()
	}
}

//export stroke
func stroke(id int) {
	if dc, ok := contexts[id]; ok {
		dc.Stroke()
	}
}

//export set_line_width
func set_line_width(id int, width float64) {
	if dc, ok := contexts[id]; ok {
		dc.SetLineWidth(width)
	}
}

//export draw_string
func draw_string(id int, text *C.char, x, y float64) {
	if dc, ok := contexts[id]; ok {
		goText := C.GoString(text)
		dc.DrawString(goText, x, y)
	}
}

//export draw_string_anchored
func draw_string_anchored(id int, text *C.char, x, y, ax, ay float64) {
	if dc, ok := contexts[id]; ok {
		goText := C.GoString(text)
		dc.DrawStringAnchored(goText, x, y, ax, ay)
	}
}

//export load_font_face
func load_font_face(id int, path *C.char, size float64) {
	if dc, ok := contexts[id]; ok {
		goPath := C.GoString(path)
		dc.LoadFontFace(goPath, size)
	}
}

//export save_png
func save_png(id int, path *C.char) {
	if dc, ok := contexts[id]; ok {
		goPath := C.GoString(path)
		dc.SavePNG(goPath)
	}
}

//export save_jpeg
func save_jpeg(id int, path *C.char, quality int) {
	if dc, ok := contexts[id]; ok {
		goPath := C.GoString(path)
		dc.SaveJPEG(goPath, quality)
	}
}

//export set_hex_color
func set_hex_color(id int, hexColor *C.char) {
	if dc, ok := contexts[id]; ok {
		goHex := C.GoString(hexColor)
		dc.SetHexColor(goHex)
	}
}

//export draw_dashed_line
func draw_dashed_line(id int, x1, y1, x2, y2 float64, pattern *float64, patternLen int) {
	if dc, ok := contexts[id]; ok {
		// Convert C array to Go slice
		dashPattern := (*[1000]float64)(unsafe.Pointer(pattern))[:patternLen:patternLen]
		dc.DrawDashedLine(x1, y1, x2, y2, dashPattern)
	}
}

//export apply_blur
func apply_blur(id int, radius float64) int {
	if dc, ok := contexts[id]; ok {
		img := dc.Image()
		blurred := advancegg.ApplyBlur(img, radius)
		
		// Create new context with blurred image
		newDC := advancegg.NewContextForImage(blurred)
		newID := nextID
		contexts[newID] = newDC
		nextID++
		return newID
	}
	return -1
}

//export apply_grayscale
func apply_grayscale(id int) int {
	if dc, ok := contexts[id]; ok {
		img := dc.Image()
		grayscale := advancegg.ApplyGrayscale(img)
		
		// Create new context with grayscale image
		newDC := advancegg.NewContextForImage(grayscale)
		newID := nextID
		contexts[newID] = newDC
		nextID++
		return newID
	}
	return -1
}

//export draw_text_on_circle
func draw_text_on_circle(id int, text *C.char, x, y, radius float64) {
	if dc, ok := contexts[id]; ok {
		goText := C.GoString(text)
		advancegg.DrawTextOnCircle(dc, goText, x, y, radius)
	}
}

// Layer management
var layerManagers = make(map[int]*advancegg.LayerManager)
var nextLayerManagerID = 1

//export create_layer_manager
func create_layer_manager(width, height int) int {
	lm := advancegg.NewLayerManager(width, height)
	id := nextLayerManagerID
	layerManagers[id] = lm
	nextLayerManagerID++
	return id
}

//export destroy_layer_manager
func destroy_layer_manager(id int) {
	delete(layerManagers, id)
}

//export add_layer
func add_layer(id int, name *C.char) int {
	if lm, ok := layerManagers[id]; ok {
		goName := C.GoString(name)
		layer := lm.AddLayer(goName)
		
		// Store layer context
		layerID := nextID
		contexts[layerID] = layer
		nextID++
		return layerID
	}
	return -1
}

//export set_layer_opacity
func set_layer_opacity(id int, name *C.char, opacity float64) {
	if lm, ok := layerManagers[id]; ok {
		goName := C.GoString(name)
		lm.SetLayerOpacity(goName, opacity)
	}
}

//export flatten_layers
func flatten_layers(id int) int {
	if lm, ok := layerManagers[id]; ok {
		result := lm.Flatten()
		
		// Store result context
		resultID := nextID
		contexts[resultID] = result
		nextID++
		return resultID
	}
	return -1
}

// Gradient support
var gradients = make(map[int]*advancegg.Gradient)
var nextGradientID = 1

//export create_linear_gradient
func create_linear_gradient(x1, y1, x2, y2 float64) int {
	gradient := advancegg.NewLinearGradient(x1, y1, x2, y2)
	id := nextGradientID
	gradients[id] = gradient
	nextGradientID++
	return id
}

//export add_color_stop
func add_color_stop(id int, position, r, g, b, a float64) {
	if gradient, ok := gradients[id]; ok {
		c := color.RGBA{
			R: uint8(r * 255),
			G: uint8(g * 255),
			B: uint8(b * 255),
			A: uint8(a * 255),
		}
		gradient.AddColorStop(position, c)
	}
}

//export set_fill_style_gradient
func set_fill_style_gradient(contextID, gradientID int) {
	if dc, ok := contexts[contextID]; ok {
		if gradient, ok := gradients[gradientID]; ok {
			dc.SetFillStyle(gradient)
		}
	}
}

//export destroy_gradient
func destroy_gradient(id int) {
	delete(gradients, id)
}

func main() {
	// Required for C shared library
}
