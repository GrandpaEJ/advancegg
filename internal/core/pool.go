package core

import (
	"image"
	"image/color"
	"sync"
)

// Memory pools to reduce GC pressure and improve performance

// ImagePool manages a pool of RGBA images to reduce allocations
type ImagePool struct {
	pools map[string]*sync.Pool // Key: "widthxheight"
	mutex sync.RWMutex
}

// GlobalImagePool is the global image pool instance
var GlobalImagePool = NewImagePool()

// NewImagePool creates a new image pool
func NewImagePool() *ImagePool {
	return &ImagePool{
		pools: make(map[string]*sync.Pool),
	}
}

// Get retrieves an image from the pool or creates a new one
func (p *ImagePool) Get(width, height int) *image.RGBA {
	key := getPoolKey(width, height)

	p.mutex.RLock()
	pool, exists := p.pools[key]
	p.mutex.RUnlock()

	if !exists {
		p.mutex.Lock()
		// Double-check after acquiring write lock
		if pool, exists = p.pools[key]; !exists {
			pool = &sync.Pool{
				New: func() interface{} {
					return image.NewRGBA(image.Rect(0, 0, width, height))
				},
			}
			p.pools[key] = pool
		}
		p.mutex.Unlock()
	}

	img := pool.Get().(*image.RGBA)

	// Clear the image
	bounds := img.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			img.SetRGBA(x, y, color.RGBA{0, 0, 0, 0})
		}
	}

	return img
}

// Put returns an image to the pool
func (p *ImagePool) Put(img *image.RGBA) {
	if img == nil {
		return
	}

	bounds := img.Bounds()
	width := bounds.Max.X - bounds.Min.X
	height := bounds.Max.Y - bounds.Min.Y
	key := getPoolKey(width, height)

	p.mutex.RLock()
	pool, exists := p.pools[key]
	p.mutex.RUnlock()

	if exists {
		pool.Put(img)
	}
}

// getPoolKey generates a key for the pool map
func getPoolKey(width, height int) string {
	return string(rune(width)) + "x" + string(rune(height))
}

// ByteSlicePool manages pools of byte slices for various operations
type ByteSlicePool struct {
	pools map[int]*sync.Pool // Key: slice length
	mutex sync.RWMutex
}

// GlobalByteSlicePool is the global byte slice pool
var GlobalByteSlicePool = NewByteSlicePool()

// NewByteSlicePool creates a new byte slice pool
func NewByteSlicePool() *ByteSlicePool {
	return &ByteSlicePool{
		pools: make(map[int]*sync.Pool),
	}
}

// Get retrieves a byte slice from the pool
func (p *ByteSlicePool) Get(size int) []byte {
	p.mutex.RLock()
	pool, exists := p.pools[size]
	p.mutex.RUnlock()

	if !exists {
		p.mutex.Lock()
		if pool, exists = p.pools[size]; !exists {
			pool = &sync.Pool{
				New: func() interface{} {
					return make([]byte, size)
				},
			}
			p.pools[size] = pool
		}
		p.mutex.Unlock()
	}

	slice := pool.Get().([]byte)

	// Clear the slice
	for i := range slice {
		slice[i] = 0
	}

	return slice
}

// Put returns a byte slice to the pool
func (p *ByteSlicePool) Put(slice []byte) {
	if slice == nil {
		return
	}

	size := len(slice)
	p.mutex.RLock()
	pool, exists := p.pools[size]
	p.mutex.RUnlock()

	if exists {
		pool.Put(slice)
	}
}

// ContextPool manages a pool of Context objects
type ContextPool struct {
	pool sync.Pool
}

// GlobalContextPool is the global context pool
var GlobalContextPool = NewContextPool()

// NewContextPool creates a new context pool
func NewContextPool() *ContextPool {
	return &ContextPool{
		pool: sync.Pool{
			New: func() interface{} {
				return &Context{}
			},
		},
	}
}

// Get retrieves a context from the pool
func (p *ContextPool) Get() *Context {
	ctx := p.pool.Get().(*Context)
	// Reset context to default state
	ctx.reset()
	return ctx
}

// Put returns a context to the pool
func (p *ContextPool) Put(ctx *Context) {
	if ctx != nil {
		p.pool.Put(ctx)
	}
}

// reset resets a context to its default state
func (ctx *Context) reset() {
	ctx.width = 0
	ctx.height = 0
	ctx.im = nil
	ctx.mask = nil
	ctx.color = nil
	ctx.fillPattern = nil
	ctx.strokePattern = nil
	ctx.strokePath = nil
	ctx.fillPath = nil
	ctx.start = Point{}
	ctx.current = Point{}
	ctx.hasCurrent = false
	ctx.dashes = nil
	ctx.dashOffset = 0
	ctx.lineWidth = 1
	ctx.lineCap = LineCapRound
	ctx.lineJoin = LineJoinRound
	ctx.fillRule = FillRuleWinding
	ctx.fontFace = nil
	ctx.fontHeight = 0
	ctx.matrix = Identity()
	ctx.stack = nil
	ctx.shadowColor = nil
	ctx.shadowOffsetX = 0
	ctx.shadowOffsetY = 0
	ctx.shadowBlur = 0
}

// PathPool manages a pool of Path2D objects
type PathPool struct {
	pool sync.Pool
}

// GlobalPathPool is the global path pool
var GlobalPathPool = NewPathPool()

// NewPathPool creates a new path pool
func NewPathPool() *PathPool {
	return &PathPool{
		pool: sync.Pool{
			New: func() interface{} {
				return &Path2D{}
			},
		},
	}
}

// Get retrieves a path from the pool
func (p *PathPool) Get() *Path2D {
	path := p.pool.Get().(*Path2D)
	// Reset the path (Path2D doesn't have Clear method, create new one)
	*path = Path2D{}
	return path
}

// Put returns a path to the pool
func (p *PathPool) Put(path *Path2D) {
	if path != nil {
		p.pool.Put(path)
	}
}

// PooledContext creates a new context using pooled resources
func PooledContext(width, height int) *Context {
	ctx := GlobalContextPool.Get()
	ctx.width = width
	ctx.height = height
	ctx.im = GlobalImagePool.Get(width, height)
	// Note: rasterizer initialization would need to be implemented
	ctx.color = color.Black
	ctx.lineWidth = 1
	ctx.lineCap = LineCapRound
	ctx.lineJoin = LineJoinRound
	ctx.fillRule = FillRuleWinding
	ctx.matrix = Identity()
	return ctx
}

// ReleaseContext returns a context and its resources to pools
func ReleaseContext(ctx *Context) {
	if ctx == nil {
		return
	}

	// Return image to pool
	if ctx.im != nil {
		GlobalImagePool.Put(ctx.im)
	}

	// Return context to pool
	GlobalContextPool.Put(ctx)
}

// PooledPath2D creates a new Path2D using pooled resources
func PooledPath2D() *Path2D {
	return GlobalPathPool.Get()
}

// ReleasePath2D returns a Path2D to the pool
func ReleasePath2D(path *Path2D) {
	GlobalPathPool.Put(path)
}

// MemoryStats provides memory usage statistics
type MemoryStats struct {
	ImagePoolSize     int
	ByteSlicePoolSize int
	ContextPoolSize   int
	PathPoolSize      int
}

// GetMemoryStats returns current memory pool statistics
func GetMemoryStats() MemoryStats {
	stats := MemoryStats{}

	// Count image pools
	GlobalImagePool.mutex.RLock()
	stats.ImagePoolSize = len(GlobalImagePool.pools)
	GlobalImagePool.mutex.RUnlock()

	// Count byte slice pools
	GlobalByteSlicePool.mutex.RLock()
	stats.ByteSlicePoolSize = len(GlobalByteSlicePool.pools)
	GlobalByteSlicePool.mutex.RUnlock()

	// Note: sync.Pool doesn't expose size information
	// In a real implementation, you might track this separately

	return stats
}

// ClearPools clears all memory pools (useful for testing)
func ClearPools() {
	GlobalImagePool = NewImagePool()
	GlobalByteSlicePool = NewByteSlicePool()
	GlobalContextPool = NewContextPool()
	GlobalPathPool = NewPathPool()
}

// PoolConfig holds configuration for memory pools
type PoolConfig struct {
	MaxImagePoolSize     int
	MaxByteSlicePoolSize int
	EnablePooling        bool
}

// DefaultPoolConfig returns the default pool configuration
func DefaultPoolConfig() PoolConfig {
	return PoolConfig{
		MaxImagePoolSize:     100,
		MaxByteSlicePoolSize: 50,
		EnablePooling:        true,
	}
}

// SetPoolConfig sets the global pool configuration
func SetPoolConfig(config PoolConfig) {
	// In a real implementation, this would configure the pools
	// For now, we'll just store the config
	_ = config
}
