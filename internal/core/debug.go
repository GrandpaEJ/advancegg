package core

import (
	"fmt"
	"image/color"
	"time"
)

// Debug mode functionality for visual debugging and performance analysis

// DebugMode represents the current debug mode
type DebugMode int

const (
	DebugModeOff DebugMode = iota
	DebugModeBasic
	DebugModeVerbose
	DebugModePerformance
)

// DebugConfig holds debug configuration
type DebugConfig struct {
	Mode               DebugMode
	ShowBounds         bool
	ShowGrid           bool
	ShowCoordinates    bool
	ShowPerformance    bool
	ShowMemoryUsage    bool
	LogOperations      bool
	HighlightSlowOps   bool
	GridSize           int
	BoundsColor        color.Color
	GridColor          color.Color
	CoordinateColor    color.Color
	PerformanceOverlay bool
}

// DefaultDebugConfig returns the default debug configuration
func DefaultDebugConfig() DebugConfig {
	return DebugConfig{
		Mode:               DebugModeOff,
		ShowBounds:         false,
		ShowGrid:           false,
		ShowCoordinates:    false,
		ShowPerformance:    false,
		ShowMemoryUsage:    false,
		LogOperations:      false,
		HighlightSlowOps:   false,
		GridSize:           50,
		BoundsColor:        color.RGBA{255, 0, 0, 128},    // Semi-transparent red
		GridColor:          color.RGBA{128, 128, 128, 64}, // Light gray
		CoordinateColor:    color.RGBA{0, 255, 0, 255},    // Green
		PerformanceOverlay: false,
	}
}

// GlobalDebugConfig is the global debug configuration
var GlobalDebugConfig = DefaultDebugConfig()

// DebugInfo holds debug information for operations
type DebugInfo struct {
	Operation   string
	StartTime   time.Time
	EndTime     time.Time
	Duration    time.Duration
	MemoryUsed  int64
	Parameters  map[string]interface{}
	BoundingBox Rectangle
}

// DebugContext wraps a Context with debug functionality
type DebugContext struct {
	*Context
	debugInfo []DebugInfo
	currentOp *DebugInfo
}

// NewDebugContext creates a new debug context
func NewDebugContext(width, height int) *DebugContext {
	return &DebugContext{
		Context:   NewContext(width, height),
		debugInfo: make([]DebugInfo, 0),
	}
}

// SetDebugMode sets the global debug mode
func SetDebugMode(mode DebugMode) {
	GlobalDebugConfig.Mode = mode

	// Configure debug settings based on mode
	switch mode {
	case DebugModeOff:
		GlobalDebugConfig.ShowBounds = false
		GlobalDebugConfig.ShowGrid = false
		GlobalDebugConfig.LogOperations = false
	case DebugModeBasic:
		GlobalDebugConfig.ShowBounds = true
		GlobalDebugConfig.LogOperations = true
	case DebugModeVerbose:
		GlobalDebugConfig.ShowBounds = true
		GlobalDebugConfig.ShowGrid = true
		GlobalDebugConfig.ShowCoordinates = true
		GlobalDebugConfig.LogOperations = true
	case DebugModePerformance:
		GlobalDebugConfig.ShowPerformance = true
		GlobalDebugConfig.ShowMemoryUsage = true
		GlobalDebugConfig.HighlightSlowOps = true
		GlobalDebugConfig.LogOperations = true
	}
}

// SetDebugConfig sets the global debug configuration
func SetDebugConfig(config DebugConfig) {
	GlobalDebugConfig = config
}

// Debug operation tracking

// startOperation begins tracking a debug operation
func (dc *DebugContext) startOperation(operation string, params map[string]interface{}) {
	if GlobalDebugConfig.Mode == DebugModeOff {
		return
	}

	dc.currentOp = &DebugInfo{
		Operation:  operation,
		StartTime:  time.Now(),
		Parameters: params,
	}

	if GlobalDebugConfig.LogOperations {
		fmt.Printf("[DEBUG] Starting operation: %s\n", operation)
	}
}

// endOperation finishes tracking a debug operation
func (dc *DebugContext) endOperation() {
	if GlobalDebugConfig.Mode == DebugModeOff || dc.currentOp == nil {
		return
	}

	dc.currentOp.EndTime = time.Now()
	dc.currentOp.Duration = dc.currentOp.EndTime.Sub(dc.currentOp.StartTime)

	// Log slow operations
	if GlobalDebugConfig.HighlightSlowOps && dc.currentOp.Duration > time.Millisecond*10 {
		fmt.Printf("[DEBUG] SLOW OPERATION: %s took %v\n", dc.currentOp.Operation, dc.currentOp.Duration)
	}

	if GlobalDebugConfig.LogOperations {
		fmt.Printf("[DEBUG] Completed operation: %s (took %v)\n", dc.currentOp.Operation, dc.currentOp.Duration)
	}

	dc.debugInfo = append(dc.debugInfo, *dc.currentOp)
	dc.currentOp = nil
}

// Debug drawing methods

// DebugDrawCircle draws a circle with debug information
func (dc *DebugContext) DebugDrawCircle(x, y, radius float64) {
	params := map[string]interface{}{
		"x": x, "y": y, "radius": radius,
	}
	dc.startOperation("DrawCircle", params)

	// Draw the actual circle
	dc.DrawCircle(x, y, radius)

	// Add debug visualizations
	if GlobalDebugConfig.ShowBounds {
		dc.drawBounds(x-radius, y-radius, radius*2, radius*2)
	}

	if GlobalDebugConfig.ShowCoordinates {
		dc.drawCoordinates(x, y)
	}

	dc.endOperation()
}

// DebugDrawRectangle draws a rectangle with debug information
func (dc *DebugContext) DebugDrawRectangle(x, y, width, height float64) {
	params := map[string]interface{}{
		"x": x, "y": y, "width": width, "height": height,
	}
	dc.startOperation("DrawRectangle", params)

	// Draw the actual rectangle
	dc.DrawRectangle(x, y, width, height)

	// Add debug visualizations
	if GlobalDebugConfig.ShowBounds {
		dc.drawBounds(x, y, width, height)
	}

	if GlobalDebugConfig.ShowCoordinates {
		dc.drawCoordinates(x, y)
		dc.drawCoordinates(x+width, y+height)
	}

	dc.endOperation()
}

// DebugDrawString draws text with debug information
func (dc *DebugContext) DebugDrawString(text string, x, y float64) {
	params := map[string]interface{}{
		"text": text, "x": x, "y": y,
	}
	dc.startOperation("DrawString", params)

	// Measure text for bounds
	width, height := dc.MeasureString(text)

	// Draw the actual text
	dc.DrawString(text, x, y)

	// Add debug visualizations
	if GlobalDebugConfig.ShowBounds {
		dc.drawBounds(x, y-height, width, height)
	}

	if GlobalDebugConfig.ShowCoordinates {
		dc.drawCoordinates(x, y)
	}

	dc.endOperation()
}

// Debug visualization helpers

// drawBounds draws bounding box around an area
func (dc *DebugContext) drawBounds(x, y, width, height float64) {
	originalColor := dc.color
	originalLineWidth := dc.lineWidth

	dc.SetColor(GlobalDebugConfig.BoundsColor)
	dc.SetLineWidth(1)
	dc.DrawRectangle(x, y, width, height)
	dc.Stroke()

	dc.SetColor(originalColor)
	dc.SetLineWidth(originalLineWidth)
}

// drawCoordinates draws coordinate markers
func (dc *DebugContext) drawCoordinates(x, y float64) {
	originalColor := dc.color

	dc.SetColor(GlobalDebugConfig.CoordinateColor)

	// Draw crosshair
	dc.DrawLine(x-5, y, x+5, y)
	dc.DrawLine(x, y-5, x, y+5)
	dc.Stroke()

	// Draw coordinate text
	coordText := fmt.Sprintf("(%.0f,%.0f)", x, y)
	dc.DrawString(coordText, x+8, y-8)

	dc.SetColor(originalColor)
}

// DrawDebugGrid draws a debug grid overlay
func (dc *DebugContext) DrawDebugGrid() {
	if !GlobalDebugConfig.ShowGrid {
		return
	}

	originalColor := dc.color
	originalLineWidth := dc.lineWidth

	dc.SetColor(GlobalDebugConfig.GridColor)
	dc.SetLineWidth(1)

	// Draw vertical lines
	for x := 0; x < dc.width; x += GlobalDebugConfig.GridSize {
		dc.DrawLine(float64(x), 0, float64(x), float64(dc.height))
		dc.Stroke()
	}

	// Draw horizontal lines
	for y := 0; y < dc.height; y += GlobalDebugConfig.GridSize {
		dc.DrawLine(0, float64(y), float64(dc.width), float64(y))
		dc.Stroke()
	}

	dc.SetColor(originalColor)
	dc.SetLineWidth(originalLineWidth)
}

// DrawPerformanceOverlay draws performance information
func (dc *DebugContext) DrawPerformanceOverlay() {
	if !GlobalDebugConfig.ShowPerformance {
		return
	}

	originalColor := dc.color

	// Calculate performance stats
	totalOps := len(dc.debugInfo)
	var totalTime time.Duration
	var slowOps int

	for _, info := range dc.debugInfo {
		totalTime += info.Duration
		if info.Duration > time.Millisecond*5 {
			slowOps++
		}
	}

	// Draw performance overlay
	dc.SetRGBA(0, 0, 0, 0.7) // Semi-transparent black background
	dc.DrawRectangle(10, 10, 200, 100)
	dc.Fill()

	dc.SetColor(color.RGBA{255, 255, 255, 255}) // White text

	y := 30.0
	dc.DrawString(fmt.Sprintf("Operations: %d", totalOps), 20, y)
	y += 20
	dc.DrawString(fmt.Sprintf("Total Time: %v", totalTime), 20, y)
	y += 20
	dc.DrawString(fmt.Sprintf("Slow Ops: %d", slowOps), 20, y)
	y += 20
	if totalOps > 0 {
		avgTime := totalTime / time.Duration(totalOps)
		dc.DrawString(fmt.Sprintf("Avg Time: %v", avgTime), 20, y)
	}

	dc.SetColor(originalColor)
}

// GetDebugInfo returns debug information for all operations
func (dc *DebugContext) GetDebugInfo() []DebugInfo {
	return dc.debugInfo
}

// ClearDebugInfo clears all debug information
func (dc *DebugContext) ClearDebugInfo() {
	dc.debugInfo = dc.debugInfo[:0]
}

// PrintDebugSummary prints a summary of debug information
func (dc *DebugContext) PrintDebugSummary() {
	fmt.Println("\n=== DEBUG SUMMARY ===")
	fmt.Printf("Total Operations: %d\n", len(dc.debugInfo))

	if len(dc.debugInfo) == 0 {
		return
	}

	var totalTime time.Duration
	opCounts := make(map[string]int)
	opTimes := make(map[string]time.Duration)

	for _, info := range dc.debugInfo {
		totalTime += info.Duration
		opCounts[info.Operation]++
		opTimes[info.Operation] += info.Duration
	}

	fmt.Printf("Total Time: %v\n", totalTime)
	fmt.Printf("Average Time: %v\n", totalTime/time.Duration(len(dc.debugInfo)))

	fmt.Println("\nOperation Breakdown:")
	for op, count := range opCounts {
		avgTime := opTimes[op] / time.Duration(count)
		fmt.Printf("  %s: %d calls, %v total, %v avg\n", op, count, opTimes[op], avgTime)
	}

	fmt.Println("=====================")
}

// Debug utilities

// IsDebugMode returns true if debug mode is enabled
func IsDebugMode() bool {
	return GlobalDebugConfig.Mode != DebugModeOff
}

// DebugLog logs a debug message if debug mode is enabled
func DebugLog(format string, args ...interface{}) {
	if GlobalDebugConfig.LogOperations {
		fmt.Printf("[DEBUG] "+format+"\n", args...)
	}
}

// DebugAssert checks a condition and logs an error if it fails
func DebugAssert(condition bool, message string) {
	if !condition && IsDebugMode() {
		fmt.Printf("[DEBUG ASSERT FAILED] %s\n", message)
	}
}
