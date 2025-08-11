# AdvanceGG Font Assets

This directory contains font files used for testing and examples in AdvanceGG.

## Font Files

### Core Fonts
- `NotoSans-Regular.ttf` - Main sans-serif font for Latin text
- `NotoSerif-Regular.ttf` - Main serif font for Latin text
- `NotoColorEmoji.ttf` - Color emoji font

### Script-Specific Fonts
- `NotoSansArabic-Regular.ttf` - Arabic script support
- `NotoSansHebrew-Regular.ttf` - Hebrew script support
- `NotoSansDevanagari-Regular.ttf` - Devanagari script support (Hindi, Sanskrit)
- `NotoSansThai-Regular.ttf` - Thai script support
- `NotoSansCJK-Regular.ttf` - Chinese, Japanese, Korean support

## Usage

These fonts are automatically detected and used by AdvanceGG's font loading system. The library will fall back to system fonts if these are not available.

### In Examples

```go
// Load a specific font
err := dc.LoadFontFace("assets/fonts/NotoSans-Regular.ttf", 16)
if err != nil {
    // Fallback to system font
    dc.LoadFontFace("/usr/share/fonts/truetype/noto/NotoSans-Regular.ttf", 16)
}

// Use emoji font
renderer := dc.GetEmojiRenderer()
renderer.LoadEmojiFont("assets/fonts/NotoColorEmoji.ttf")
```

### Font Loading Priority

1. Local assets/fonts/ directory
2. System font directories (/usr/share/fonts, /System/Library/Fonts, etc.)
3. Built-in fallback fonts

## License

These fonts are from the Noto font family by Google, licensed under the SIL Open Font License 1.1.

## Installation

To set up fonts for development:

```bash
# Copy system fonts to local assets (Linux)
cp /usr/share/fonts/truetype/noto/NotoSans-Regular.ttf assets/fonts/
cp /usr/share/fonts/truetype/noto/NotoColorEmoji.ttf assets/fonts/

# Or use the setup script
go run scripts/setup-fonts.go
```

## Testing

Run font tests with:

```bash
go run examples/unicode-emoji.go
go run examples/font-loading-test.go
```

This will test:
- Font loading from different sources
- Unicode text rendering
- Emoji rendering with color fonts
- Script-specific text shaping
- Fallback font behavior
