package main

import (
	"fmt"
	
	"github.com/GrandpaEJ/advancegg"
)

func main() {
	dc := advancegg.NewContext(800, 600)
	
	// White background
	dc.SetRGB(1, 1, 1)
	dc.Clear()
	
	// Title
	dc.SetRGB(0, 0, 0)
	dc.DrawString("Text Metrics Demonstration", 50, 50)
	
	// Demonstrate various text metrics
	demonstrateBasicMetrics(dc)
	demonstrateDetailedMetrics(dc)
	demonstrateTextWrapping(dc)
	demonstrateFontMetrics(dc)
	
	// Save the demonstration
	dc.SavePNG("images/text/text-metrics.png")
	fmt.Println("Text metrics demo saved as text-metrics.png")
	
	fmt.Println("Text metrics examples completed!")
}

func demonstrateBasicMetrics(dc *advancegg.Context) {
	dc.SetRGB(0, 0, 0)
	dc.DrawString("Basic Text Metrics:", 50, 100)
	
	testText := "Hello, World!"
	
	// Get basic measurements
	width, height := dc.MeasureString(testText)
	textWidth := dc.GetTextWidth(testText)
	textHeight := dc.GetTextHeight()
	lineHeight := dc.GetLineHeight()
	baseline := dc.GetBaseline()
	
	// Draw the text
	x, y := 50.0, 150.0
	dc.SetRGB(0, 0, 1)
	dc.DrawString(testText, x, y)
	
	// Draw measurement indicators
	dc.SetRGB(1, 0, 0)
	dc.SetLineWidth(1)
	
	// Width indicator
	dc.DrawLine(x, y+10, x+width, y+10)
	dc.Stroke()
	
	// Height indicator
	dc.DrawLine(x-10, y-baseline, x-10, y-baseline+height)
	dc.Stroke()
	
	// Baseline
	dc.SetRGB(0, 1, 0)
	dc.DrawLine(x, y, x+width, y)
	dc.Stroke()
	
	// Display measurements
	dc.SetRGB(0, 0, 0)
	dc.DrawString(fmt.Sprintf("Width: %.1f", width), x+width+20, y-20)
	dc.DrawString(fmt.Sprintf("Height: %.1f", height), x+width+20, y)
	dc.DrawString(fmt.Sprintf("Text Width: %.1f", textWidth), x+width+20, y+20)
	dc.DrawString(fmt.Sprintf("Text Height: %.1f", textHeight), x+width+20, y+40)
	dc.DrawString(fmt.Sprintf("Line Height: %.1f", lineHeight), x+width+20, y+60)
	dc.DrawString(fmt.Sprintf("Baseline: %.1f", baseline), x+width+20, y+80)
}

func demonstrateDetailedMetrics(dc *advancegg.Context) {
	dc.SetRGB(0, 0, 0)
	dc.DrawString("Detailed Text Metrics:", 50, 250)
	
	testText := "Typography"
	
	// Get detailed metrics
	metrics := dc.MeasureTextMetrics(testText)
	
	// Draw the text
	x, y := 50.0, 300.0
	dc.SetRGB(0, 0, 1)
	dc.DrawString(testText, x, y)
	
	// Draw bounding box
	dc.SetRGB(1, 0, 0)
	dc.SetLineWidth(1)
	dc.DrawRectangle(x+metrics.ActualBoundingBoxLeft, y-metrics.ActualBoundingBoxAscent, 
		metrics.ActualBoundingBoxRight-metrics.ActualBoundingBoxLeft, 
		metrics.ActualBoundingBoxAscent+metrics.ActualBoundingBoxDescent)
	dc.Stroke()
	
	// Draw font bounding box
	dc.SetRGB(0, 1, 0)
	dc.DrawRectangle(x, y-metrics.FontBoundingBoxAscent, 
		metrics.Width, metrics.FontBoundingBoxAscent+metrics.FontBoundingBoxDescent)
	dc.Stroke()
	
	// Display detailed measurements
	dc.SetRGB(0, 0, 0)
	infoX := x + metrics.Width + 20
	dc.DrawString(fmt.Sprintf("Width: %.1f", metrics.Width), infoX, y-60)
	dc.DrawString(fmt.Sprintf("Height: %.1f", metrics.Height), infoX, y-40)
	dc.DrawString(fmt.Sprintf("Ascent: %.1f", metrics.FontBoundingBoxAscent), infoX, y-20)
	dc.DrawString(fmt.Sprintf("Descent: %.1f", metrics.FontBoundingBoxDescent), infoX, y)
	dc.DrawString(fmt.Sprintf("Em Height Ascent: %.1f", metrics.EmHeightAscent), infoX, y+20)
	dc.DrawString(fmt.Sprintf("Em Height Descent: %.1f", metrics.EmHeightDescent), infoX, y+40)
}

func demonstrateTextWrapping(dc *advancegg.Context) {
	dc.SetRGB(0, 0, 0)
	dc.DrawString("Text Wrapping:", 50, 400)
	
	longText := "This is a long text that will be wrapped to fit within a specified width. The text wrapping function automatically breaks lines at word boundaries."
	maxWidth := 300.0
	
	// Wrap the text
	lines := dc.WrapText(longText, maxWidth)
	
	// Draw the wrapped text
	x, y := 50.0, 430.0
	lineHeight := dc.GetLineHeight()
	
	// Draw bounding box
	dc.SetRGB(0.8, 0.8, 0.8)
	dc.DrawRectangle(x, y, maxWidth, float64(len(lines))*lineHeight)
	dc.Fill()
	
	// Draw the text lines
	dc.SetRGB(0, 0, 0)
	for i, line := range lines {
		dc.DrawString(line, x+5, y+float64(i+1)*lineHeight)
	}
	
	// Display wrapping info
	dc.DrawString(fmt.Sprintf("Original text length: %d chars", len(longText)), x+maxWidth+20, y+20)
	dc.DrawString(fmt.Sprintf("Wrapped into %d lines", len(lines)), x+maxWidth+20, y+40)
	dc.DrawString(fmt.Sprintf("Max width: %.1f", maxWidth), x+maxWidth+20, y+60)
}

func demonstrateFontMetrics(dc *advancegg.Context) {
	dc.SetRGB(0, 0, 0)
	dc.DrawString("Font Metrics:", 400, 100)
	
	// Get font metrics
	ascent, descent, lineGap := dc.GetFontMetrics()
	
	// Draw baseline and metrics
	x, y := 400.0, 150.0
	testText := "Baseline Example"
	
	dc.SetRGB(0, 0, 1)
	dc.DrawString(testText, x, y)
	
	// Draw baseline
	dc.SetRGB(0, 1, 0)
	dc.SetLineWidth(2)
	dc.DrawLine(x, y, x+200, y)
	dc.Stroke()
	
	// Draw ascent line
	dc.SetRGB(1, 0, 0)
	dc.SetLineWidth(1)
	dc.DrawLine(x, y-ascent, x+200, y-ascent)
	dc.Stroke()
	
	// Draw descent line
	dc.SetRGB(1, 0, 0)
	dc.DrawLine(x, y+descent, x+200, y+descent)
	dc.Stroke()
	
	// Draw line gap
	dc.SetRGB(0.5, 0.5, 0.5)
	dc.DrawLine(x, y+descent+lineGap, x+200, y+descent+lineGap)
	dc.Stroke()
	
	// Labels
	dc.SetRGB(0, 0, 0)
	dc.DrawString("Baseline", x+210, y+5)
	dc.DrawString(fmt.Sprintf("Ascent: %.1f", ascent), x+210, y-ascent+5)
	dc.DrawString(fmt.Sprintf("Descent: %.1f", descent), x+210, y+descent+5)
	dc.DrawString(fmt.Sprintf("Line Gap: %.1f", lineGap), x+210, y+descent+lineGap+5)
	
	// Multiple lines example
	dc.SetRGB(0, 0, 0)
	dc.DrawString("Multi-line Example:", 400, 250)
	
	multilineText := "Line 1\nLine 2\nLine 3"
	lineSpacing := 1.2
	
	x, y = 400.0, 280.0
	width, height := dc.MeasureMultilineString(multilineText, lineSpacing)
	
	// Draw bounding box
	dc.SetRGB(0.9, 0.9, 0.9)
	dc.DrawRectangle(x, y-ascent, width, height)
	dc.Fill()
	
	// Draw the multiline text
	dc.SetRGB(0, 0, 0)
	dc.DrawStringWrapped(multilineText, x, y, 0, 0, width, lineSpacing, advancegg.AlignLeft)
	
	// Display multiline metrics
	dc.DrawString(fmt.Sprintf("Multiline Width: %.1f", width), x+width+20, y)
	dc.DrawString(fmt.Sprintf("Multiline Height: %.1f", height), x+width+20, y+20)
	dc.DrawString(fmt.Sprintf("Line Spacing: %.1f", lineSpacing), x+width+20, y+40)
}
