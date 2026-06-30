# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [v1.7.1] - Bilingual Text Rendering

### Features

- **Bilingual Text Support** — বাংলা-English mixed text এখন একই line-এ ঠিকমতো render হয়
  - `segmentBidiRuns()` rewrite করা হয়েছে text-কে script boundary-তে split করার জন্য
  - নতুন `scriptForRune()` helper — প্রতিটি character-এর script detect করে
  - Mixed text এখন আলাদা sub-runs এ ভাগ হয় (যেমন "আমি" → Bengali run, "today" → Latin run)
  - প্রতিটি run নিজের script font ব্যবহার করে render হয়

- **Indic Danda Support** — Bengali/Hindi danda (।) এখন সঠিকভাবে render হয়
  - `scriptForRune()`-এ danda (U+0964, U+0965) কে Bengali script হিসেবে treat করা হয়
  - Devanagari font (NotoSansDevanagari) যোগ করা হয়েছে danda rendering-এর জন্য
  - `drawShapedString()`-এ glyph fallback logic যোগ করা হয়েছে

- **Emoji Rendering Fix (Partial)** — `drawShapedString()`-এ emoji detection যোগ করা হয়েছে
  - Emoji characters এখন emoji renderer ব্যবহার করে render হয়
  - Note: Some edge cases still need fixing (see Issue #15)

### Implementation

- `internal/core/unicode.go`:
  - নতুন `scriptForRune()` function যোগ করা হয়েছে
  - `segmentBidiRuns()` rewrite করা হয়েছে per-script segmentation-এর জন্য
  - নতুন `DetectScript()` public wrapper যোগ করা হয়েছে

- `internal/core/context.go`:
  - `drawShapedString()` update করা হয়েছে emoji detection সহ
  - Glyph fallback logic যোগ করা হয়েছে script font-এর জন্য

- `assets/fonts/`:
  - `NotoSansDevanagari-Regular.ttf` যোগ করা হয়েছে danda support-এর জন্য

- `examples/text/bengali-bilingual-demo.go`:
  - নতুন bilingual demo example যোগ করা হয়েছে
  - বাংলা-English mixed text examples সহ 4টি card

### Examples

- **Bilingual Demo**: `examples/text/bengali-bilingual-demo.go` — বাংলা-English mixed text rendering demo

## [v1.7.0] - Indic Script Rendering & Emoji Support

### Features

- **Bengali GPOS Rendering** — বাংলা (ও অন্যান্য Indic) script-এর জন্য complete glyph positioning
  - GPOS mark-to-base attachment: matras (া, ি, ী, ু, ূ, ৃ, ে, ৈ, ো, ৌ) এখন base consonant-এর সাথে attach হয়, আলাদা character না হয়ে
  - Pre-base matras (ে, ৈ) GSUB reordering-এর মাধ্যমে base consonant-এর আগে ঠিকঠাক render হয়
  - Below-base matras (ু, ূ, ৃ) GPOS offsets সহ সঠিক position-এ বসে
  - Above-base marks (ং, ঁ, ়) GPOS Y-offsets দিয়ে properly positioned
  - 75+ conjuncts (যুক্তবর্ণ) GSUB-এর মাধ্যমে সঠিকভাবে formed

- **Script Font API** — `LoadScriptFont()` যোগ করা হয়েছে: per-script font loading (যেমন বাংলা text-এর জন্য NotoSansBengali, Latin-এর জন্য main font)

- **Font Size Fix** — Shaping এখন `fontHeight` অনুযায়ী correct font size use করে, আগের hardcoded `Size=32`-এর পরিবর্তে — এতে glyph-গুলোর মাঝে excessive spacing দূর হয়েছে

- **Emoji Rendering** — Full emoji text rendering support, automatic detection সহ
  - Automatic emoji detection and rendering in `DrawString` methods
  - Color emoji font support (Noto Color Emoji, Apple Color Emoji, Segoe UI Emoji)
  - ZWJ (Zero Width Joiner) sequence handling: complex emojis (family, profession, couples) ঠিকভাবে render হয়
  - Skin tone modifier support
  - Emoji fallback rendering, category-based colors সহ
  - Performance-optimized emoji caching system

- **API**: `SetEnableAutoEmoji()` এবং `GetEnableAutoEmoji()` যোগ করা হয়েছে — automatic emoji rendering control-এর জন্য
- **API**: `SetEmojiSize()` এবং `GetEmojiSize()` যোগ করা হয়েছে — emoji size control-এর জন্য
- **API**: `ScriptBengali` public ScriptType constant হিসেবে যোগ করা হয়েছে
- **Examples**: যোগ করা হয়েছে `examples/text/bengali-demo.go` — Bengali speech bubble demo
- **Examples**: যোগ করা হয়েছে `examples/text/bengali-gpos-test.go` — comprehensive GPOS test (all matras, 28 consonant tables, 75+ conjuncts, 30 words, 15 sentences)

### Implementation

- Fixed `shapeRun()` in `unicode.go`: এখন HarfBuzz GPOS-এর `XOffset`/`YOffset`/`Advance` use করে, আগের ভুল scaling-এর বদলে
- `TextShaper`-এ `fontSize` field যোগ করা হয়েছে; shaping-এর `Size` এখন `fontHeight` থেকে derive হয়
- GPOS data `go-text/typesetting`-এর shaping output-এ আগে থেকেই ছিল — previous code শুধু `XOffset`/`YOffset` ignore করত
- `drawGlyphWithFont()` যোগ করা হয়েছে — per-script truetype font থেকে glyph render করার জন্য
- Created `emoji_integration.go`: text/emoji segmentation and mixed rendering
- Created `emoji_api.go`: public API methods for emoji control
- `DrawStringAnchored` update করা হয়েছে `drawMixedString` use করার জন্য — seamless text/emoji rendering
- নতুন context-এ emoji rendering by default enabled

## [v1.6.0] - Fallback Fixes & Benchmarking Tools

### Features
- **Benchmarks**: Added automated visualization script for benchmark results (`benchmarks/plot`)
- **API**: Exported batch operation types for external usage

### Bug Fixes
- **Core**: Fixed text shaping fallback to use Go font metrics when HarfBuzz is unavailable (Issue #13)
- **Linting**: Resolved various static analysis warnings in core and benchmark modules

## [v1.5.0] - Blend Modes & Architecture Refactor

### 🌟 **Major Features**

#### **Porter-Duff Blend Modes**
- **Added** support for 12 Porter-Duff compositing operators (`Clear`, `Source`, `Dest`, `SrcOver`, `DstOver`, `SrcIn`, `DstIn`, `SrcOut`, `DstOut`, `SrcAtop`, `DstAtop`, `Xor`)
- **Added** new blend modes: `Add` (Linear Dodge) and `plus` operators
- **Refactored** `compositeLayer` engine to distinguishing between color production and alpha composition
- **Fixed** alpha blending math for handling semi-transparent source/destinations correctly

### 🏗 **Architectural Improvements**

#### **Project Restructuring**
- **Reorganized** `examples/` directory into categorized submodules (`basics/`, `shapes/`, `text/`, `filters/`, etc.)
- **Reorganized** image assets into structured `images/` subdirectories matching example categories
- **Cleaned** project root by removing loose artifacts (`*.png`, `*.svg`, `*.gif`)
- **Added** robust recursive example runner script (`scripts/run_examples.go`) for regression testing

#### **Quality Assurance**
- **Added** comprehensive unit tests for all blend modes in `internal/core/layers_blend_test.go`
- **Added** visual verification generator `examples/layers/blend_modes_gen.go` including all 29 blend modes
- **Verified** regression suite across 70+ examples

### 🐛 **Bug Fixes**
- **Fixed** `alphaBlend` function to correctly handle transparent destination backgrounds (critical for non-rectangular layers)
- **Fixed** various example file paths to support new directory structure

### 🎨 **Major New Features Added**

#### **Layer Compositing System**
- **Added** comprehensive multi-layer compositing system with Photoshop-style functionality
- **Added** 13 professional blend modes: Normal, Multiply, Screen, Overlay, Soft Light, Hard Light, Color Dodge, Color Burn, Darken, Lighten, Difference, Exclusion
- **Added** layer opacity and visibility controls with `LayerManager` API
- **Added** automatic layer compositing pipeline for professional graphics

#### **DOM-Style Object Model**
- **Added** modern web-like API with tree structure for shapes
- **Added** element IDs and CSS-like class system with style inheritance
- **Added** CSS-style selectors (ID, class, type) and dynamic style application
- **Added** element hierarchy and manipulation for complex graphics

#### **Hit Testing System**
- **Added** comprehensive point-in-path detection for interactive graphics
- **Added** support for rectangles, circles, ellipses, polygons, lines, and paths
- **Added** spatial indexing for O(log n) performance optimization
- **Added** `HitTestManager` for managing multiple hit testable objects

#### **Animation Framework**
- **Added** frame-based animation system with GIF export capabilities
- **Added** 9 easing functions: Linear, Ease In/Out, Cubic variations, Bounce, Elastic
- **Added** animation sequences and predefined animations with property helpers
- **Added** `Animator` class with automatic frame management

#### **Advanced Stroke Styles**
- **Added** comprehensive dashed pattern support and gradient strokes
- **Added** tapered stroke effects for calligraphy and custom line cap styles
- **Added** stroke width variation along paths for artistic effects

#### **SVG Import/Export**
- **Added** basic SVG document parsing and rendering support
- **Added** support for rectangles, circles, lines, and paths with color parsing
- **Added** viewBox and coordinate transformation support

#### **Performance Optimizations**
- **Added** spatial indexing, render caching, and parallel processing
- **Added** memory optimization techniques and performance monitoring
- **Added** batch rendering optimizations for complex scenes

### 📚 **New Examples and Documentation**
- **Added** `examples/layer-compositing.go` - Layer and blend mode examples
- **Added** `examples/dom-object-model.go` - DOM-style graphics programming
- **Added** `examples/hit-testing.go` - Interactive graphics and UI examples
- **Added** `examples/animation-demo.go` - Animation and easing functions
- **Added** `examples/advanced-strokes.go` - Advanced stroke styling
- **Added** comprehensive HTML documentation with interactive examples

### ⚡ **Performance Improvements**
- **Improved** rendering pipeline efficiency by 40%
- **Reduced** hit testing complexity from O(n) to O(log n) with spatial indexing
- **Enhanced** memory usage with object pooling and smart garbage collection

### 🎯 **Generated Visual Demos**
- `blend-modes-demo.png`, `dom-styling-demo.png`, `hit-testing-basic.png`
- `bouncing-ball.gif`, `easing-demo.gif`, `complex-sequence.gif`
- Multiple performance and feature demonstration images

## [v1.0.0] - Initial Release

### Added
- Initial release of AdvanceGG graphics library
- Core drawing context with support for 2D graphics rendering
- Basic shape drawing functions:
  - Points, lines, rectangles, circles, ellipses
  - Rounded rectangles and regular polygons
  - Arcs and elliptical arcs
- Path drawing capabilities:
  - MoveTo, LineTo for basic paths
  - QuadraticTo and CubicTo for Bézier curves
  - Path management (ClosePath, ClearPath, NewSubPath)
- Text rendering functions:
  - DrawString and DrawStringAnchored
  - DrawStringWrapped with word wrapping support
  - MeasureString and MeasureMultilineString
  - Font loading and management
- Color management:
  - RGB and RGBA color setting
  - Hex color support
  - Color utilities
- Stroke and fill operations:
  - Fill and Stroke with preserve options
  - Line width, cap, and join settings
  - Dash patterns and offsets
  - Fill rule configuration
- Gradient and pattern support:
  - Linear, radial, and conic gradients
  - Surface patterns with repeat operations
  - Custom pattern interface
- Transformation functions:
  - Translation, scaling, rotation, and shearing
  - Transform about specific points
  - Matrix transformations
  - Coordinate system inversion
- State management:
  - Push and Pop for context state
  - Identity matrix reset
- Clipping operations:
  - Clip and ClipPreserve
  - Mask support with alpha channels
  - Clipping region management
- Image operations:
  - DrawImage and DrawImageAnchored
  - Image loading (PNG, JPEG)
  - Context creation from existing images
- Utility functions:
  - Angle conversion (radians/degrees)
  - Image I/O helpers
  - Mathematical utilities
- Comprehensive example collection:
  - Basic shapes and drawing
  - Text rendering examples
  - Gradient and pattern examples
  - Transformation demonstrations
  - Complex graphics compositions
- Complete documentation:
  - Getting started guide
  - Comprehensive API reference
  - Example gallery with code
  - Contributing guidelines
- Proper Go module structure:
  - Organized package layout
  - Clean import paths
  - Dependency management
- Test coverage for core functionality

### Technical Details
- Pure Go implementation with minimal external dependencies
- Optimized for performance and memory efficiency
- Cross-platform compatibility
- Thread-safe operations where applicable
- Comprehensive error handling

### Dependencies
- `github.com/golang/freetype` for font rendering
- `golang.org/x/image` for extended image support

### Package Structure
- `pkg/advancegg/` - Main library package
- `examples/` - Example programs and demonstrations
- `docs/` - Comprehensive documentation
- `examples/images/` - Image assets for examples

### Breaking Changes
- N/A (Initial release)

### Migration Guide
- N/A (Initial release)

---

## Release Notes

This is the initial release of AdvanceGG, a powerful 2D graphics library for Go. The library provides a simple and intuitive API for rendering graphics, with support for shapes, text, gradients, transformations, and more.

### Key Features
- **Simple API**: Easy-to-use functions for common graphics operations
- **Rich Drawing Functions**: Comprehensive set of drawing primitives
- **Advanced Features**: Gradients, patterns, transformations, and clipping
- **Pure Go**: Minimal external dependencies
- **High Performance**: Optimized for speed and memory efficiency
- **Comprehensive Documentation**: Complete guides and examples

### Getting Started
```bash
go get github.com/GrandpaEJ/advancegg
```

See the [Getting Started Guide](docs/getting-started.md) for detailed installation and usage instructions.
