package advance

import (
	"image/color"
	"math"
	"testing"
)

func TestLinearGradientPattern(t *testing.T) {
	// Create a simple linear gradient from red to blue
	pattern := CreateLinearGradient(100, 100, color.RGBA{255, 0, 0, 255}, color.RGBA{0, 0, 255, 255})

	// Test start point (should be red)
	c := pattern.ColorAt(0, 0)
	r, g, b, a := c.RGBA()
	if r>>8 != 255 || g>>8 != 0 || b>>8 != 0 || a>>8 != 255 {
		t.Errorf("Expected red at start, got RGBA(%d, %d, %d, %d)", r>>8, g>>8, b>>8, a>>8)
	}

	// Test end point (should be blue)
	c = pattern.ColorAt(100, 100)
	r, g, b, a = c.RGBA()
	if r>>8 != 0 || g>>8 != 0 || b>>8 != 255 || a>>8 != 255 {
		t.Errorf("Expected blue at end, got RGBA(%d, %d, %d, %d)", r>>8, g>>8, b>>8, a>>8)
	}
}

func TestRadialGradientPattern(t *testing.T) {
	// Create a radial gradient from center
	pattern := CreateRadialGradient(50, 50, 50, color.RGBA{255, 255, 255, 255}, color.RGBA{0, 0, 0, 255})

	// Test center (should be white)
	c := pattern.ColorAt(50, 50)
	r, g, b, a := c.RGBA()
	if r>>8 != 255 || g>>8 != 255 || b>>8 != 255 || a>>8 != 255 {
		t.Errorf("Expected white at center, got RGBA(%d, %d, %d, %d)", r>>8, g>>8, b>>8, a>>8)
	}

	// Test edge (should be black)
	c = pattern.ColorAt(100, 50)
	r, g, b, a = c.RGBA()
	if r>>8 != 0 || g>>8 != 0 || b>>8 != 0 || a>>8 != 255 {
		t.Errorf("Expected black at edge, got RGBA(%d, %d, %d, %d)", r>>8, g>>8, b>>8, a>>8)
	}
}

func TestCheckerboardPattern(t *testing.T) {
	pattern := CreateCheckerboard(10)

	// Test different cells
	c1 := pattern.ColorAt(5, 5)   // Should be white
	c2 := pattern.ColorAt(15, 5)  // Should be black
	c3 := pattern.ColorAt(5, 15)  // Should be black
	c4 := pattern.ColorAt(15, 15) // Should be white

	// Check that adjacent cells have different colors
	if c1 == c2 {
		t.Error("Adjacent checkerboard cells should have different colors")
	}
	if c1 == c3 {
		t.Error("Adjacent checkerboard cells should have different colors")
	}
	if c1 != c4 {
		t.Error("Diagonal checkerboard cells should have same colors")
	}
}

func TestStripePattern(t *testing.T) {
	pattern := CreateStripes(10)

	// Test that stripes alternate
	c1 := pattern.ColorAt(5, 0)
	c2 := pattern.ColorAt(15, 0)

	if c1 == c2 {
		t.Error("Adjacent stripes should have different colors")
	}
}

func TestPolkaDotPattern(t *testing.T) {
	pattern := CreatePolkaDots(20, 5)

	// Test center of dot (should be dot color)
	c1 := pattern.ColorAt(10, 10)

	// Test far from dot (should be background color)
	c2 := pattern.ColorAt(0, 0)

	if c1 == c2 {
		t.Error("Dot center and background should have different colors")
	}
}

func TestPatternTransform_Identity(t *testing.T) {
	transform := NewPatternTransform()

	// Identity transform should not change coordinates
	x, y := 10.0, 20.0
	tx, ty := transform.transformPoint(x, y)

	if tx != x || ty != y {
		t.Errorf("Identity transform should not change coordinates: got (%f, %f), expected (%f, %f)", tx, ty, x, y)
	}
}

func TestPatternTransform_Translate(t *testing.T) {
	transform := NewPatternTransform().Translate(5, 10)

	x, y := 0.0, 0.0
	tx, ty := transform.transformPoint(x, y)

	if tx != 5.0 || ty != 10.0 {
		t.Errorf("Translation failed: got (%f, %f), expected (5, 10)", tx, ty)
	}
}

func TestPatternTransform_Scale(t *testing.T) {
	transform := NewPatternTransform().Scale(2, 3)

	x, y := 1.0, 1.0
	tx, ty := transform.transformPoint(x, y)

	if tx != 2.0 || ty != 3.0 {
		t.Errorf("Scaling failed: got (%f, %f), expected (2, 3)", tx, ty)
	}
}

func TestPatternTransform_Rotate(t *testing.T) {
	transform := NewPatternTransform().Rotate(math.Pi / 2) // 90 degrees

	x, y := 1.0, 0.0
	tx, ty := transform.transformPoint(x, y)

	// After 90-degree rotation, (1,0) should become approximately (0,1) or (0,-1) depending on rotation direction
	if math.Abs(tx) > 1e-10 || (math.Abs(ty-1.0) > 1e-10 && math.Abs(ty+1.0) > 1e-10) {
		t.Errorf("Rotation failed: got (%f, %f), expected (0, ±1)", tx, ty)
	}
}

func TestTransformablePattern(t *testing.T) {
	// Create a simple pattern
	basePattern := CreateCheckerboard(10)

	// Create transformable pattern with translation
	transformedPattern := WithTranslation(basePattern, 5, 5)

	// The transformed pattern should give different results
	c2 := transformedPattern.ColorAt(0, 0)

	// Due to the translation, the color at (0,0) in the transformed pattern
	// should be the same as the color at (5,5) in the base pattern
	c3 := basePattern.ColorAt(5, 5)

	if c2 != c3 {
		t.Error("Transformed pattern should apply translation correctly")
	}
}

func TestNewTransformablePattern(t *testing.T) {
	basePattern := CreateCheckerboard(10)
	transformable := NewTransformablePattern(basePattern)

	if transformable.Pattern != basePattern {
		t.Error("NewTransformablePattern should wrap the base pattern")
	}

	// Should start with identity transform
	c1 := basePattern.ColorAt(5, 5)
	c2 := transformable.ColorAt(5, 5)

	if c1 != c2 {
		t.Error("Identity transform should not change pattern colors")
	}
}

func BenchmarkLinearGradient(b *testing.B) {
	pattern := CreateLinearGradient(100, 100, color.RGBA{255, 0, 0, 255}, color.RGBA{0, 0, 255, 255})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = pattern.ColorAt(50, 50)
	}
}

func BenchmarkCheckerboard(b *testing.B) {
	pattern := CreateCheckerboard(10)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = pattern.ColorAt(float64(i%100), float64(i%100))
	}
}

func BenchmarkTransformablePattern(b *testing.B) {
	basePattern := CreateCheckerboard(10)
	transformedPattern := WithRotation(basePattern, math.Pi/4)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = transformedPattern.ColorAt(float64(i%100), float64(i%100))
	}
}
