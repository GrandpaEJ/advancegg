
---

## 🧠 Suggested Additions to TODO & Feature Roadmap

### 🔥 High Priority (Critical/Powerful)

* [x] **Layer System** – Multi-layered drawing (like Photoshop layers) ✅ COMPLETED
* [x] **Non-destructive Editing** – Filters/transformations as reversible states ✅ COMPLETED
* [x] **Smart Guides & Alignment** – Snap to grid, guides, center, baseline ✅ COMPLETED
* [x] **Unicode Shaping Support** – Indic, RTL, complex scripts ✅ COMPLETED
* [x] **Emoji Rendering** – Color emoji font + fallback ✅ COMPLETED
* [ ] **Color Profiles (ICC)** – Accurate color conversion for print-ready files

---

### 📋 Medium Priority (Usability & Power Tools)

* [ ] **Text-on-Path** – Draw text along curves or shapes
* [ ] **Stroke Styles** – Dashed stroke patterns, gradient stroke, taper ends
* [ ] **Layer Compositing Modes** – Multiply, Screen, Overlay (CSS-like blending)
* [ ] **SVG Parser/Loader** – Import SVG shapes into canvas
* [ ] **Filter Chains** – Chainable filter pipelines (like shaders)
* [ ] **Object Model (DOM-style)** – Tree structure for shapes with IDs & styles

---

### 📝 Low Priority (Flexibility / Cross-Platform)

* [ ] **Canvas to HTML5 Export** – Export canvas as HTML5 `<canvas>` with JS fallback
* [ ] **FFI-safe API** – Call from Rust, Zig, WASM, etc.
* [ ] **Texture Atlas Generator** – For games/UI sprites
* [ ] **Color Blindness Simulation** – Filter preview for protanopia, deuteranopia
* [ ] **Tiled Rendering Engine** – High-resolution render in chunks
* [ ] **Offline Font Subsetter** – Reduce font size for only used glyphs

---

### 🛠 Developer/Power User Features

* [ ] **Live Reload** – Auto-refresh image on code/file changes
* [ ] **Canvas Inspector** – Show bounding boxes, alignment guides (visual debug)
* [ ] **Visual Unit Grid** – px/inch/mm toggle + Rulers
* [ ] **Render Trace Log** – View ordered draw calls for debugging
* [ ] **Configurable Defaults** – Central style/theme (colors, fonts, stroke)

---

### 🌍 Ecosystem / Interop / Utility

* [ ] **Markdown to Image Renderer** – Turn `.md` or rich text into styled images
* [ ] **Chart Drawing API** – Built-in drawing for bar/line/pie charts
* [ ] **Image Metadata Reader** – Read EXIF, ICC, DPI, orientation
* [ ] **Template System** – Define reusable image templates (with variables)
* [ ] **Headless Browser Preview** – Use Chrome headless to preview/export via script
* [ ] **Graphviz-style Graph API** – Node/edge diagram support

---

## ✅ Bonus: Feature Tagging Suggestion

Consider tagging each TODO item with:

```
[core] [io] [render] [filter] [interop] [dev] [perf] [ux]
```

🔍 Example:

```md
- [ ] [filter] Blur, Sharpen, Posterize filters
- [ ] [render] Layer Compositing (Multiply, Overlay)
- [ ] [interop] Python Bindings via cgo
- [ ] [perf] SIMD Path Rasterization
```

---
