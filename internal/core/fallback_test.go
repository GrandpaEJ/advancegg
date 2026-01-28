package core

import (
	"testing"

	"golang.org/x/image/font/basicfont"
)

func TestFallbackShapingWithGoFont(t *testing.T) {
	// Initialize Shaper
	ts := NewTextShaper()

	// Set Go font face ONLY (hbFont remains nil)
	// basicfont.Face7x13 has a fixed width of 7 pixels
	ts.SetGoFontFace(basicfont.Face7x13)

	// Verify hbFont is nil (internal check, but effectively we haven't called SetFont/SetFontBytes)
	if ts.HasFont() {
		t.Fatal("Expected no HarfBuzz font to be loaded")
	}

	text := "Hello"
	shaped := ts.ShapeText(text)

	if len(shaped.Glyphs) != 5 {
		t.Fatalf("Expected 5 glyphs, got %d", len(shaped.Glyphs))
	}

	// Check advances
	// Before fix: 12.0
	// After fix: 7.0 (from Face7x13)
	expectedAdvance := 7.0

	for i, glyph := range shaped.Glyphs {
		if glyph.AdvanceX != expectedAdvance {
			t.Errorf("Glyph %d: expected advance %.2f, got %.2f", i, expectedAdvance, glyph.AdvanceX)
		}
	}
}
