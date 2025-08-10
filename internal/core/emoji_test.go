package core

import (
	"testing"
)

func TestEmojiRenderer_Creation(t *testing.T) {
	renderer := NewEmojiRenderer()
	if renderer == nil {
		t.Fatal("NewEmojiRenderer() returned nil")
	}
	
	if len(renderer.ColorFonts) == 0 {
		t.Error("Expected some default color fonts")
	}
	
	if renderer.EmojiSize <= 0 {
		t.Errorf("Expected positive emoji size, got %f", renderer.EmojiSize)
	}
	
	if renderer.Cache == nil {
		t.Error("Expected cache to be initialized")
	}
}

func TestEmojiRenderer_IsEmoji(t *testing.T) {
	tests := []struct {
		r        rune
		expected bool
		desc     string
	}{
		{0x1F600, true, "grinning face"},
		{0x1F1FA, true, "regional indicator U"},
		{0x1F1F8, true, "regional indicator S"},
		{0x200D, true, "zero width joiner"},
		{0x1F3FB, true, "light skin tone"},
		{0x2600, true, "black sun with rays"},
		{0x2764, true, "heavy black heart"},
		{'A', false, "latin letter A"},
		{'1', false, "digit 1"},
		{' ', false, "space"},
		{0x0041, false, "latin A"},
	}
	
	for _, test := range tests {
		result := IsEmoji(test.r)
		if result != test.expected {
			t.Errorf("IsEmoji(%U) [%s] = %v, expected %v", test.r, test.desc, result, test.expected)
		}
	}
}

func TestEmojiRenderer_ParseEmojiSequence(t *testing.T) {
	renderer := NewEmojiRenderer()
	
	// Test simple emoji
	sequences := renderer.ParseEmojiSequence("😀")
	if len(sequences) != 1 {
		t.Errorf("Expected 1 sequence for simple emoji, got %d", len(sequences))
	}
	if len(sequences) > 0 && len(sequences[0].Runes) != 1 {
		t.Errorf("Expected 1 rune in sequence, got %d", len(sequences[0].Runes))
	}
	
	// Test text with no emoji
	sequences = renderer.ParseEmojiSequence("Hello")
	if len(sequences) != 0 {
		t.Errorf("Expected 0 sequences for text without emoji, got %d", len(sequences))
	}
	
	// Test mixed text and emoji
	sequences = renderer.ParseEmojiSequence("Hello 😀 World")
	if len(sequences) != 1 {
		t.Errorf("Expected 1 sequence for mixed text, got %d", len(sequences))
	}
}

func TestEmojiRenderer_SkinToneModifiers(t *testing.T) {
	renderer := NewEmojiRenderer()
	
	tests := []struct {
		r        rune
		expected bool
		name     string
	}{
		{0x1F3FB, true, "light skin tone"},
		{0x1F3FC, true, "medium-light skin tone"},
		{0x1F3FD, true, "medium skin tone"},
		{0x1F3FE, true, "medium-dark skin tone"},
		{0x1F3FF, true, "dark skin tone"},
		{0x1F600, false, "grinning face"},
		{'A', false, "latin A"},
	}
	
	for _, test := range tests {
		result := renderer.isSkinToneModifier(test.r)
		if result != test.expected {
			t.Errorf("isSkinToneModifier(%U) [%s] = %v, expected %v", test.r, test.name, result, test.expected)
		}
	}
}

func TestEmojiRenderer_GetSkinToneName(t *testing.T) {
	renderer := NewEmojiRenderer()
	
	tests := []struct {
		r        rune
		expected string
	}{
		{0x1F3FB, "light"},
		{0x1F3FC, "medium-light"},
		{0x1F3FD, "medium"},
		{0x1F3FE, "medium-dark"},
		{0x1F3FF, "dark"},
		{0x1F600, ""}, // not a skin tone
	}
	
	for _, test := range tests {
		result := renderer.getSkinToneName(test.r)
		if result != test.expected {
			t.Errorf("getSkinToneName(%U) = %q, expected %q", test.r, result, test.expected)
		}
	}
}

func TestEmojiRenderer_GetEmojiCategory(t *testing.T) {
	renderer := NewEmojiRenderer()
	
	tests := []struct {
		r        rune
		expected string
		desc     string
	}{
		{0x1F600, EmojiCategorySmileys, "grinning face"},
		{0x1F466, EmojiCategoryPeople, "boy"},
		{0x1F400, EmojiCategoryAnimals, "rat"},
		{0x1F32D, EmojiCategoryFood, "hot dog"},
		{0x1F3A0, EmojiCategoryActivities, "carousel horse"},
		{0x1F680, EmojiCategoryTravel, "rocket"},
		{0x1F4A0, EmojiCategoryObjects, "diamond with a dot"},
		{0x1F500, EmojiCategorySymbols, "twisted rightwards arrows"},
		{0x1F1FA, EmojiCategoryFlags, "regional indicator U"},
	}
	
	for _, test := range tests {
		result := renderer.GetEmojiCategory(test.r)
		if result != test.expected {
			t.Errorf("GetEmojiCategory(%U) [%s] = %q, expected %q", test.r, test.desc, result, test.expected)
		}
	}
}

func TestEmojiRenderer_RenderEmoji(t *testing.T) {
	renderer := NewEmojiRenderer()
	
	// Create a simple emoji sequence
	sequence := EmojiSequence{
		Runes:    []rune{0x1F600}, // 😀
		Text:     "😀",
		Category: EmojiCategorySmileys,
	}
	
	// Test rendering
	img := renderer.RenderEmoji(sequence, 32)
	if img == nil {
		t.Fatal("RenderEmoji returned nil")
	}
	
	bounds := img.Bounds()
	if bounds.Dx() != 32 || bounds.Dy() != 32 {
		t.Errorf("Expected 32x32 image, got %dx%d", bounds.Dx(), bounds.Dy())
	}
}

func TestEmojiRenderer_Cache(t *testing.T) {
	renderer := NewEmojiRenderer()
	
	sequence := EmojiSequence{
		Runes:    []rune{0x1F600},
		Text:     "😀",
		Category: EmojiCategorySmileys,
	}
	
	// First render
	img1 := renderer.RenderEmoji(sequence, 32)
	if img1 == nil {
		t.Fatal("First render returned nil")
	}
	
	// Second render should use cache
	img2 := renderer.RenderEmoji(sequence, 32)
	if img2 == nil {
		t.Fatal("Second render returned nil")
	}
	
	// Should be the same image from cache
	if img1 != img2 {
		t.Error("Expected cached image to be returned")
	}
}

func BenchmarkEmojiRenderer_ParseSequence(b *testing.B) {
	renderer := NewEmojiRenderer()
	text := "Hello 😀 World 👋 Test 🌍"
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sequences := renderer.ParseEmojiSequence(text)
		_ = sequences
	}
}

func BenchmarkEmojiRenderer_RenderEmoji(b *testing.B) {
	renderer := NewEmojiRenderer()
	sequence := EmojiSequence{
		Runes:    []rune{0x1F600},
		Text:     "😀",
		Category: EmojiCategorySmileys,
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		img := renderer.RenderEmoji(sequence, 32)
		_ = img
	}
}

func BenchmarkIsEmoji(b *testing.B) {
	runes := []rune{
		'A', '1', ' ',
		0x1F600, 0x1F1FA, 0x200D,
		0x2600, 0x1F3FB,
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, r := range runes {
			_ = IsEmoji(r)
		}
	}
}
