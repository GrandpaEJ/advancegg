package core

import (
	"fmt"
	"image"
	"image/color"
	"runtime"
	"sync"
	"time"
)

// Performance optimizations for graphics rendering

// SpatialIndex provides spatial indexing for fast object lookup
type SpatialIndex struct {
	bounds   Bounds
	objects  []SpatialObject
	children []*SpatialIndex
	maxDepth int
	depth    int
	maxItems int
}

// SpatialObject represents an object that can be spatially indexed
type SpatialObject interface {
	GetBounds() (float64, float64, float64, float64)
	GetID() string
}

// Bounds represents a bounding rectangle
type Bounds struct {
	X, Y, Width, Height float64
}

// NewSpatialIndex creates a new spatial index
func NewSpatialIndex(bounds Bounds, maxDepth, maxItems int) *SpatialIndex {
	return &SpatialIndex{
		bounds:   bounds,
		objects:  make([]SpatialObject, 0),
		children: nil,
		maxDepth: maxDepth,
		depth:    0,
		maxItems: maxItems,
	}
}

// Insert adds an object to the spatial index
func (si *SpatialIndex) Insert(obj SpatialObject) {
	if si.children != nil {
		// If subdivided, insert into appropriate child
		index := si.getIndex(obj)
		if index != -1 {
			si.children[index].Insert(obj)
			return
		}
	}

	si.objects = append(si.objects, obj)

	// Check if we need to subdivide
	if len(si.objects) > si.maxItems && si.depth < si.maxDepth {
		si.subdivide()

		// Move objects to children
		i := 0
		for i < len(si.objects) {
			index := si.getIndex(si.objects[i])
			if index != -1 {
				si.children[index].Insert(si.objects[i])
				// Remove from current level
				si.objects = append(si.objects[:i], si.objects[i+1:]...)
			} else {
				i++
			}
		}
	}
}

// Query returns objects that intersect with the given bounds
func (si *SpatialIndex) Query(bounds Bounds) []SpatialObject {
	var result []SpatialObject

	// Check if bounds intersect with this node
	if !si.intersects(bounds) {
		return result
	}

	// Add objects from this node
	for _, obj := range si.objects {
		objBounds := si.getObjectBounds(obj)
		if si.rectanglesIntersect(bounds, objBounds) {
			result = append(result, obj)
		}
	}

	// Query children if subdivided
	if si.children != nil {
		for _, child := range si.children {
			result = append(result, child.Query(bounds)...)
		}
	}

	return result
}

// subdivide splits the node into four quadrants
func (si *SpatialIndex) subdivide() {
	halfWidth := si.bounds.Width / 2
	halfHeight := si.bounds.Height / 2

	si.children = make([]*SpatialIndex, 4)

	// Top-left
	si.children[0] = &SpatialIndex{
		bounds:   Bounds{si.bounds.X, si.bounds.Y, halfWidth, halfHeight},
		objects:  make([]SpatialObject, 0),
		maxDepth: si.maxDepth,
		depth:    si.depth + 1,
		maxItems: si.maxItems,
	}

	// Top-right
	si.children[1] = &SpatialIndex{
		bounds:   Bounds{si.bounds.X + halfWidth, si.bounds.Y, halfWidth, halfHeight},
		objects:  make([]SpatialObject, 0),
		maxDepth: si.maxDepth,
		depth:    si.depth + 1,
		maxItems: si.maxItems,
	}

	// Bottom-left
	si.children[2] = &SpatialIndex{
		bounds:   Bounds{si.bounds.X, si.bounds.Y + halfHeight, halfWidth, halfHeight},
		objects:  make([]SpatialObject, 0),
		maxDepth: si.maxDepth,
		depth:    si.depth + 1,
		maxItems: si.maxItems,
	}

	// Bottom-right
	si.children[3] = &SpatialIndex{
		bounds:   Bounds{si.bounds.X + halfWidth, si.bounds.Y + halfHeight, halfWidth, halfHeight},
		objects:  make([]SpatialObject, 0),
		maxDepth: si.maxDepth,
		depth:    si.depth + 1,
		maxItems: si.maxItems,
	}
}

// getIndex determines which child quadrant an object belongs to
func (si *SpatialIndex) getIndex(obj SpatialObject) int {
	objBounds := si.getObjectBounds(obj)

	verticalMidpoint := si.bounds.X + si.bounds.Width/2
	horizontalMidpoint := si.bounds.Y + si.bounds.Height/2

	// Object can completely fit within the top quadrants
	topQuadrant := objBounds.Y < horizontalMidpoint && objBounds.Y+objBounds.Height < horizontalMidpoint

	// Object can completely fit within the bottom quadrants
	bottomQuadrant := objBounds.Y > horizontalMidpoint

	// Object can completely fit within the left quadrants
	if objBounds.X < verticalMidpoint && objBounds.X+objBounds.Width < verticalMidpoint {
		if topQuadrant {
			return 0 // Top-left
		} else if bottomQuadrant {
			return 2 // Bottom-left
		}
	}

	// Object can completely fit within the right quadrants
	if objBounds.X > verticalMidpoint {
		if topQuadrant {
			return 1 // Top-right
		} else if bottomQuadrant {
			return 3 // Bottom-right
		}
	}

	return -1 // Object doesn't fit completely in any quadrant
}

// Helper methods
func (si *SpatialIndex) getObjectBounds(obj SpatialObject) Bounds {
	x1, y1, x2, y2 := obj.GetBounds()
	return Bounds{x1, y1, x2 - x1, y2 - y1}
}

func (si *SpatialIndex) intersects(bounds Bounds) bool {
	return si.rectanglesIntersect(si.bounds, bounds)
}

func (si *SpatialIndex) rectanglesIntersect(a, b Bounds) bool {
	return a.X < b.X+b.Width && a.X+a.Width > b.X && a.Y < b.Y+b.Height && a.Y+a.Height > b.Y
}

// RenderCache provides caching for rendered elements
type RenderCache struct {
	cache       map[string]*image.RGBA
	mutex       sync.RWMutex
	maxSize     int
	currentSize int
}

// NewRenderCache creates a new render cache
func NewRenderCache(maxSize int) *RenderCache {
	return &RenderCache{
		cache:   make(map[string]*image.RGBA),
		maxSize: maxSize,
	}
}

// Get retrieves a cached image
func (rc *RenderCache) Get(key string) (*image.RGBA, bool) {
	rc.mutex.RLock()
	defer rc.mutex.RUnlock()

	img, exists := rc.cache[key]
	return img, exists
}

// Set stores an image in the cache
func (rc *RenderCache) Set(key string, img *image.RGBA) {
	rc.mutex.Lock()
	defer rc.mutex.Unlock()

	// Check if we need to evict items
	if rc.currentSize >= rc.maxSize {
		rc.evictLRU()
	}

	rc.cache[key] = img
	rc.currentSize++
}

// evictLRU removes the least recently used item (simplified)
func (rc *RenderCache) evictLRU() {
	// Simple eviction - remove first item found
	for key := range rc.cache {
		delete(rc.cache, key)
		rc.currentSize--
		break
	}
}

// Clear clears the cache
func (rc *RenderCache) Clear() {
	rc.mutex.Lock()
	defer rc.mutex.Unlock()

	rc.cache = make(map[string]*image.RGBA)
	rc.currentSize = 0
}

// ParallelRenderer provides parallel rendering capabilities
type ParallelRenderer struct {
	workers int
	jobs    chan RenderJob
	results chan RenderResult
	wg      sync.WaitGroup
}

// RenderJob represents a rendering job
type RenderJob struct {
	ID     string
	Bounds Bounds
	Render func() *image.RGBA
}

// RenderResult represents a rendering result
type RenderResult struct {
	ID    string
	Image *image.RGBA
	Error error
}

// NewParallelRenderer creates a new parallel renderer
func NewParallelRenderer(workers int) *ParallelRenderer {
	if workers <= 0 {
		workers = runtime.NumCPU()
	}

	return &ParallelRenderer{
		workers: workers,
		jobs:    make(chan RenderJob, workers*2),
		results: make(chan RenderResult, workers*2),
	}
}

// Start starts the parallel renderer workers
func (pr *ParallelRenderer) Start() {
	for i := 0; i < pr.workers; i++ {
		go pr.worker()
	}
}

// Stop stops the parallel renderer
func (pr *ParallelRenderer) Stop() {
	close(pr.jobs)
	pr.wg.Wait()
	close(pr.results)
}

// Submit submits a rendering job
func (pr *ParallelRenderer) Submit(job RenderJob) {
	pr.jobs <- job
}

// GetResult gets a rendering result
func (pr *ParallelRenderer) GetResult() RenderResult {
	return <-pr.results
}

// worker processes rendering jobs
func (pr *ParallelRenderer) worker() {
	pr.wg.Add(1)
	defer pr.wg.Done()

	for job := range pr.jobs {
		result := RenderResult{ID: job.ID}

		start := time.Now()
		result.Image = job.Render()
		duration := time.Since(start)

		// Log performance if needed
		_ = duration

		pr.results <- result
	}
}

// PerformanceMonitor tracks rendering performance
type PerformanceMonitor struct {
	frameCount    int
	totalTime     time.Duration
	lastFrameTime time.Time
	fps           float64
	mutex         sync.Mutex
}

// NewPerformanceMonitor creates a new performance monitor
func NewPerformanceMonitor() *PerformanceMonitor {
	return &PerformanceMonitor{
		lastFrameTime: time.Now(),
	}
}

// StartFrame marks the start of a frame
func (pm *PerformanceMonitor) StartFrame() {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	now := time.Now()
	if pm.frameCount > 0 {
		frameDuration := now.Sub(pm.lastFrameTime)
		pm.totalTime += frameDuration
		pm.fps = float64(pm.frameCount) / pm.totalTime.Seconds()
	}

	pm.lastFrameTime = now
	pm.frameCount++
}

// GetFPS returns the current FPS
func (pm *PerformanceMonitor) GetFPS() float64 {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()
	return pm.fps
}

// GetFrameCount returns the total frame count
func (pm *PerformanceMonitor) GetFrameCount() int {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()
	return pm.frameCount
}

// Reset resets the performance monitor
func (pm *PerformanceMonitor) Reset() {
	pm.mutex.Lock()
	defer pm.mutex.Unlock()

	pm.frameCount = 0
	pm.totalTime = 0
	pm.lastFrameTime = time.Now()
	pm.fps = 0
}

// OptimizedContext provides an optimized rendering context
type OptimizedContext struct {
	*Context
	spatialIndex *SpatialIndex
	renderCache  *RenderCache
	perfMonitor  *PerformanceMonitor
	dirtyRegions []Bounds
}

// NewOptimizedContext creates a new optimized context
func NewOptimizedContext(width, height int) *OptimizedContext {
	bounds := Bounds{0, 0, float64(width), float64(height)}

	return &OptimizedContext{
		Context:      NewContext(width, height),
		spatialIndex: NewSpatialIndex(bounds, 6, 10),
		renderCache:  NewRenderCache(100),
		perfMonitor:  NewPerformanceMonitor(),
		dirtyRegions: make([]Bounds, 0),
	}
}

// AddDirtyRegion marks a region as needing redraw
func (oc *OptimizedContext) AddDirtyRegion(bounds Bounds) {
	oc.dirtyRegions = append(oc.dirtyRegions, bounds)
}

// ClearDirtyRegions clears all dirty regions
func (oc *OptimizedContext) ClearDirtyRegions() {
	oc.dirtyRegions = oc.dirtyRegions[:0]
}

// GetDirtyRegions returns the current dirty regions
func (oc *OptimizedContext) GetDirtyRegions() []Bounds {
	return oc.dirtyRegions
}

// GetSpatialIndex returns the spatial index
func (oc *OptimizedContext) GetSpatialIndex() *SpatialIndex {
	return oc.spatialIndex
}

// GetRenderCache returns the render cache
func (oc *OptimizedContext) GetRenderCache() *RenderCache {
	return oc.renderCache
}

// GetPerformanceMonitor returns the performance monitor
func (oc *OptimizedContext) GetPerformanceMonitor() *PerformanceMonitor {
	return oc.perfMonitor
}

// Optimized drawing methods

// DrawOptimizedRectangle draws a rectangle with caching
func (oc *OptimizedContext) DrawOptimizedRectangle(id string, x, y, width, height float64, fillColor color.Color) {
	cacheKey := fmt.Sprintf("rect_%s_%.0f_%.0f_%.0f_%.0f_%v", id, x, y, width, height, fillColor)

	if cached, exists := oc.renderCache.Get(cacheKey); exists {
		// Use cached version
		oc.DrawImage(cached, int(x), int(y))
		return
	}

	// Render and cache
	tempCtx := NewContext(int(width), int(height))
	tempCtx.SetColor(fillColor)
	tempCtx.DrawRectangle(0, 0, width, height)
	tempCtx.Fill()

	oc.renderCache.Set(cacheKey, tempCtx.Image().(*image.RGBA))
	oc.DrawImage(tempCtx.Image(), int(x), int(y))
}

// Frustum culling for efficient rendering
func (oc *OptimizedContext) IsVisible(bounds Bounds) bool {
	viewBounds := Bounds{0, 0, float64(oc.Width()), float64(oc.Height())}
	return oc.spatialIndex.rectanglesIntersect(bounds, viewBounds)
}
