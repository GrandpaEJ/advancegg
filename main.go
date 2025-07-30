// Package advancegg provides a simple API for rendering 2D graphics in pure Go.
// This is the main entry point that re-exports all functionality from internal/core.
package advancegg

// Re-export all public types and functions from the core package
import (
	"github.com/GrandpaEJ/advancegg/internal/core"
)

// Context represents a 2D graphics context
type Context = core.Context

// Point represents a 2D point
type Point = core.Point

// Matrix represents a 2D transformation matrix
type Matrix = core.Matrix

// Path2D represents a 2D path that can be reused and manipulated
type Path2D = core.Path2D

// Color space types
type Color = core.Color
type CMYK = core.CMYK
type HSV = core.HSV
type HSL = core.HSL
type LAB = core.LAB
type XYZ = core.XYZ

// ImageData type for pixel manipulation
type ImageData = core.ImageData

// TextMetrics type for advanced text measurement
type TextMetrics = core.TextMetrics

// Pattern interface for fill and stroke patterns
type Pattern = core.Pattern

// Gradient interface for gradient patterns
type Gradient = core.Gradient

// Line cap styles
type LineCap = core.LineCap

const (
	LineCapRound  = core.LineCapRound
	LineCapButt   = core.LineCapButt
	LineCapSquare = core.LineCapSquare
)

// Line join styles
type LineJoin = core.LineJoin

const (
	LineJoinRound = core.LineJoinRound
	LineJoinBevel = core.LineJoinBevel
)

// Fill rules
type FillRule = core.FillRule

const (
	FillRuleWinding = core.FillRuleWinding
	FillRuleEvenOdd = core.FillRuleEvenOdd
)

// Alignment options
type Align = core.Align

const (
	AlignLeft   = core.AlignLeft
	AlignCenter = core.AlignCenter
	AlignRight  = core.AlignRight
)

// Repeat operations for patterns
type RepeatOp = core.RepeatOp

const (
	RepeatBoth = core.RepeatBoth
	RepeatX    = core.RepeatX
	RepeatY    = core.RepeatY
	RepeatNone = core.RepeatNone
)

// Context creation functions
var (
	NewContext         = core.NewContext
	NewContextForImage = core.NewContextForImage
	NewContextForRGBA  = core.NewContextForRGBA
)

// Path2D creation functions
var (
	NewPath2D         = core.NewPath2D
	NewPath2DFromPath = core.NewPath2DFromPath
)

// Pattern creation functions
var (
	NewSolidPattern   = core.NewSolidPattern
	NewLinearGradient = core.NewLinearGradient
	NewRadialGradient = core.NewRadialGradient
	NewConicGradient  = core.NewConicGradient
	NewSurfacePattern = core.NewSurfacePattern
)

// Utility functions
var (
	Radians   = core.Radians
	Degrees   = core.Degrees
	LoadImage = core.LoadImage
	LoadPNG   = core.LoadPNG
	LoadJPG   = core.LoadJPG
	LoadGIF   = core.LoadGIF
	LoadBMP   = core.LoadBMP
	LoadTIFF  = core.LoadTIFF
	LoadWebP  = core.LoadWebP
	SavePNG   = core.SavePNG
	SaveJPEG  = core.SaveJPEG
	SaveGIF   = core.SaveGIF
	SaveBMP   = core.SaveBMP
	SaveTIFF  = core.SaveTIFF
)

// Font loading functions
var (
	LoadFontFace             = core.LoadFontFace
	LoadTTFFace              = core.LoadTTFFace
	LoadOTFFace              = core.LoadOTFFace
	ParseFontFace            = core.ParseFontFace
	ParseFontFaceWithOptions = core.ParseFontFaceWithOptions
	GetFontFormat            = core.GetFontFormat
)

// Bezier curve functions
var (
	QuadraticBezier = core.QuadraticBezier
	CubicBezier     = core.CubicBezier
)

// Matrix functions
var (
	Identity  = core.Identity
	Translate = core.Translate
	Scale     = core.Scale
	Rotate    = core.Rotate
	Shear     = core.Shear
)

// Image filter functions
var (
	Grayscale     = core.Grayscale
	Invert        = core.Invert
	Sepia         = core.Sepia
	Brightness    = core.Brightness
	Contrast      = core.Contrast
	Blur          = core.Blur
	Sharpen       = core.Sharpen
	Threshold     = core.Threshold
	Pixelate      = core.Pixelate
	Noise         = core.Noise
	EdgeDetection = core.EdgeDetection
	Emboss        = core.Emboss
	Posterize     = core.Posterize
	Vignette      = core.Vignette
)

// Color space functions
var (
	NewColor            = core.NewColor
	NewColorFromRGBA255 = core.NewColorFromRGBA255
)

// ImageData functions
var (
	NewImageData          = core.NewImageData
	NewImageDataFromImage = core.NewImageDataFromImage
)

// Layer system exports
type Layer = core.Layer
type LayerManager = core.LayerManager
type BlendMode = core.BlendMode

const (
	BlendModeNormal     = core.BlendModeNormal
	BlendModeMultiply   = core.BlendModeMultiply
	BlendModeScreen     = core.BlendModeScreen
	BlendModeOverlay    = core.BlendModeOverlay
	BlendModeSoftLight  = core.BlendModeSoftLight
	BlendModeHardLight  = core.BlendModeHardLight
	BlendModeColorDodge = core.BlendModeColorDodge
	BlendModeColorBurn  = core.BlendModeColorBurn
	BlendModeDarken     = core.BlendModeDarken
	BlendModeLighten    = core.BlendModeLighten
	BlendModeDifference = core.BlendModeDifference
	BlendModeExclusion  = core.BlendModeExclusion
	BlendModeHue        = core.BlendModeHue
	BlendModeSaturation = core.BlendModeSaturation
	BlendModeColor      = core.BlendModeColor
	BlendModeLuminosity = core.BlendModeLuminosity
)

var (
	NewLayerManager = core.NewLayerManager
)

// Non-destructive editing exports
type EditOperation = core.EditOperation
type EditStack = core.EditStack
type BrightnessOperation = core.BrightnessOperation
type ContrastOperation = core.ContrastOperation
type BlurOperation = core.BlurOperation
type CropOperation = core.CropOperation

var (
	NewEditStack = core.NewEditStack
)

// Guides and alignment exports
type Guide = core.Guide
type GuideManager = core.GuideManager
type GuideType = core.GuideType
type GuideOrientation = core.GuideOrientation
type AlignmentTarget = core.AlignmentTarget
type Margins = core.Margins

const (
	GuideTypeManual   = core.GuideTypeManual
	GuideTypeGrid     = core.GuideTypeGrid
	GuideTypeCenter   = core.GuideTypeCenter
	GuideTypeBaseline = core.GuideTypeBaseline
	GuideTypeMargin   = core.GuideTypeMargin

	GuideHorizontal = core.GuideHorizontal
	GuideVertical   = core.GuideVertical
)

var (
	NewGuideManager = core.NewGuideManager

	// Alignment functions
	AlignTargetsLeft             = core.AlignTargetsLeft
	AlignTargetsRight            = core.AlignTargetsRight
	AlignTargetsTop              = core.AlignTargetsTop
	AlignTargetsBottom           = core.AlignTargetsBottom
	AlignTargetsCenterHorizontal = core.AlignTargetsCenterHorizontal
	AlignTargetsCenterVertical   = core.AlignTargetsCenterVertical
	DistributeHorizontally       = core.DistributeHorizontally
	DistributeVertically         = core.DistributeVertically
)

// Unicode shaping exports
type TextShaper = core.TextShaper
type ShapedText = core.ShapedText
type ShapedGlyph = core.ShapedGlyph
type TextDirection = core.TextDirection
type ScriptType = core.ScriptType

const (
	TextDirectionLTR = core.TextDirectionLTR
	TextDirectionRTL = core.TextDirectionRTL
	TextDirectionTTB = core.TextDirectionTTB
	TextDirectionBTT = core.TextDirectionBTT

	ScriptLatin      = core.ScriptLatin
	ScriptArabic     = core.ScriptArabic
	ScriptHebrew     = core.ScriptHebrew
	ScriptDevanagari = core.ScriptDevanagari
	ScriptThai       = core.ScriptThai
	ScriptChinese    = core.ScriptChinese
	ScriptJapanese   = core.ScriptJapanese
	ScriptKorean     = core.ScriptKorean
	ScriptCyrillic   = core.ScriptCyrillic
	ScriptGreek      = core.ScriptGreek
)

var (
	NewTextShaper = core.NewTextShaper
)

// Emoji rendering exports
type EmojiRenderer = core.EmojiRenderer
type EmojiInfo = core.EmojiInfo
type EmojiSequence = core.EmojiSequence

var (
	NewEmojiRenderer = core.NewEmojiRenderer
	IsEmoji          = core.IsEmoji
)
