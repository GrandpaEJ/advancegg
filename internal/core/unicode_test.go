package core

import (
	"testing"
)

func TestTextShaper_Creation(t *testing.T) {
	shaper := NewTextShaper()
	if shaper == nil {
		t.Fatal("NewTextShaper() returned nil")
	}

	if shaper.Direction != TextDirectionLTR {
		t.Errorf("Expected default direction LTR, got %v", shaper.Direction)
	}

	if shaper.Script != ScriptLatin {
		t.Errorf("Expected default script Latin, got %v", shaper.Script)
	}

	if shaper.Language != "en" {
		t.Errorf("Expected default language 'en', got %s", shaper.Language)
	}
}

func TestTextShaper_ScriptDetection(t *testing.T) {
	shaper := NewTextShaper()

	tests := []struct {
		text     string
		expected ScriptType
	}{
		{"Hello World", ScriptLatin},
		{"مرحبا", ScriptArabic},
		{"שלום", ScriptHebrew},
		{"नमस्ते", ScriptDevanagari},
		{"สวัสดี", ScriptThai},
		{"你好", ScriptChinese},
		{"Привет", ScriptCyrillic},
		{"Γεια", ScriptGreek},
	}

	for _, test := range tests {
		detected := shaper.detectScript(test.text)
		if detected != test.expected {
			t.Errorf("detectScript(%q) = %v, expected %v", test.text, detected, test.expected)
		}
	}
}

func TestTextShaper_BidiRunSegmentation(t *testing.T) {
	shaper := NewTextShaper()

	// Test simple LTR text
	runs := shaper.segmentBidiRuns("Hello World")
	if len(runs) != 1 {
		t.Errorf("Expected 1 run for LTR text, got %d", len(runs))
	}
	if runs[0].Direction != TextDirectionLTR {
		t.Errorf("Expected LTR direction, got %v", runs[0].Direction)
	}

	// Test simple RTL text
	runs = shaper.segmentBidiRuns("مرحبا")
	if len(runs) != 1 {
		t.Errorf("Expected 1 run for RTL text, got %d", len(runs))
	}
	if runs[0].Direction != TextDirectionRTL {
		t.Errorf("Expected RTL direction, got %v", runs[0].Direction)
	}
}

func TestTextShaper_ShapeText(t *testing.T) {
	shaper := NewTextShaper()

	// Test basic shaping (fallback mode)
	shaped := shaper.ShapeText("Hello")
	if shaped == nil {
		t.Fatal("ShapeText returned nil")
	}

	if len(shaped.Glyphs) != 5 {
		t.Errorf("Expected 5 glyphs for 'Hello', got %d", len(shaped.Glyphs))
	}

	// Check that glyphs have reasonable properties
	for i, glyph := range shaped.Glyphs {
		if glyph.AdvanceX <= 0 {
			t.Errorf("Glyph %d has non-positive advance: %f", i, glyph.AdvanceX)
		}
		if glyph.Cluster != i {
			t.Errorf("Glyph %d has wrong cluster: %d", i, glyph.Cluster)
		}
	}
}

func TestTextShaper_EmptyText(t *testing.T) {
	shaper := NewTextShaper()

	shaped := shaper.ShapeText("")
	if shaped == nil {
		t.Fatal("ShapeText returned nil for empty string")
	}

	if len(shaped.Glyphs) != 0 {
		t.Errorf("Expected 0 glyphs for empty string, got %d", len(shaped.Glyphs))
	}
}

func TestIsEmoji(t *testing.T) {
	tests := []struct {
		r        rune
		expected bool
	}{
		{'A', false},
		{'1', false},
		{' ', false},
		{0x1F600, true}, // 😀
		{0x1F1FA, true}, // Regional indicator U
		{0x1F1F8, true}, // Regional indicator S
		{0x200D, true},  // Zero width joiner
		{0x2600, true},  // ☀ (sun)
		{0x1F3FB, true}, // Light skin tone
	}

	for _, test := range tests {
		result := IsEmoji(test.r)
		if result != test.expected {
			t.Errorf("IsEmoji(%U) = %v, expected %v", test.r, result, test.expected)
		}
	}
}

func BenchmarkTextShaper_ShapeText(b *testing.B) {
	shaper := NewTextShaper()
	text := "Hello World! This is a test of text shaping performance."

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		shaped := shaper.ShapeText(text)
		_ = shaped
	}
}

func BenchmarkScriptDetection(b *testing.B) {
	shaper := NewTextShaper()
	texts := []string{
		"Hello World",
		"مرحبا بالعالم",
		"שלום עולם",
		"नमस्ते दुनिया",
		"สวัสดีชาวโลก",
		"你好世界",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, text := range texts {
			_ = shaper.detectScript(text)
		}
	}
}
