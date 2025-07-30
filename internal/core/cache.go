package core

import (
	"crypto/md5"
	"fmt"
	"image"
	"sync"
	"time"
)

// Caching system for rendered elements and frequently used resources

// CacheEntry represents a cached item
type CacheEntry struct {
	Data      interface{}
	CreatedAt time.Time
	AccessAt  time.Time
	AccessCount int64
	Size      int64
}

// Cache represents a generic cache with LRU eviction
type Cache struct {
	items    map[string]*CacheEntry
	maxSize  int64
	maxItems int
	mutex    sync.RWMutex
	stats    CacheStats
}

// CacheStats holds cache performance statistics
type CacheStats struct {
	Hits        int64
	Misses      int64
	Evictions   int64
	TotalSize   int64
	ItemCount   int
}

// NewCache creates a new cache with specified limits
func NewCache(maxSize int64, maxItems int) *Cache {
	return &Cache{
		items:    make(map[string]*CacheEntry),
		maxSize:  maxSize,
		maxItems: maxItems,
	}
}

// Get retrieves an item from the cache
func (c *Cache) Get(key string) (interface{}, bool) {
	c.mutex.RLock()
	entry, exists := c.items[key]
	c.mutex.RUnlock()
	
	if !exists {
		c.mutex.Lock()
		c.stats.Misses++
		c.mutex.Unlock()
		return nil, false
	}
	
	c.mutex.Lock()
	entry.AccessAt = time.Now()
	entry.AccessCount++
	c.stats.Hits++
	c.mutex.Unlock()
	
	return entry.Data, true
}

// Set stores an item in the cache
func (c *Cache) Set(key string, data interface{}, size int64) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	
	// Check if item already exists
	if existing, exists := c.items[key]; exists {
		c.stats.TotalSize -= existing.Size
		c.stats.ItemCount--
	}
	
	// Create new entry
	entry := &CacheEntry{
		Data:        data,
		CreatedAt:   time.Now(),
		AccessAt:    time.Now(),
		AccessCount: 1,
		Size:        size,
	}
	
	c.items[key] = entry
	c.stats.TotalSize += size
	c.stats.ItemCount++
	
	// Evict if necessary
	c.evictIfNeeded()
}

// evictIfNeeded removes items if cache limits are exceeded
func (c *Cache) evictIfNeeded() {
	for c.stats.TotalSize > c.maxSize || c.stats.ItemCount > c.maxItems {
		c.evictLRU()
	}
}

// evictLRU removes the least recently used item
func (c *Cache) evictLRU() {
	var oldestKey string
	var oldestTime time.Time = time.Now()
	
	for key, entry := range c.items {
		if entry.AccessAt.Before(oldestTime) {
			oldestTime = entry.AccessAt
			oldestKey = key
		}
	}
	
	if oldestKey != "" {
		entry := c.items[oldestKey]
		delete(c.items, oldestKey)
		c.stats.TotalSize -= entry.Size
		c.stats.ItemCount--
		c.stats.Evictions++
	}
}

// Clear removes all items from the cache
func (c *Cache) Clear() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	
	c.items = make(map[string]*CacheEntry)
	c.stats = CacheStats{}
}

// GetStats returns cache statistics
func (c *Cache) GetStats() CacheStats {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.stats
}

// Specialized caches

// ImageCache caches rendered images
type ImageCache struct {
	*Cache
}

// GlobalImageCache is the global image cache
var GlobalImageCache = NewImageCache(100*1024*1024, 1000) // 100MB, 1000 items

// NewImageCache creates a new image cache
func NewImageCache(maxSize int64, maxItems int) *ImageCache {
	return &ImageCache{
		Cache: NewCache(maxSize, maxItems),
	}
}

// GetImage retrieves a cached image
func (ic *ImageCache) GetImage(key string) (*image.RGBA, bool) {
	data, exists := ic.Get(key)
	if !exists {
		return nil, false
	}
	return data.(*image.RGBA), true
}

// SetImage stores an image in the cache
func (ic *ImageCache) SetImage(key string, img *image.RGBA) {
	bounds := img.Bounds()
	size := int64(bounds.Max.X * bounds.Max.Y * 4) // 4 bytes per pixel (RGBA)
	ic.Set(key, img, size)
}

// FontCache caches font faces and metrics
type FontCache struct {
	*Cache
}

// GlobalFontCache is the global font cache
var GlobalFontCache = NewFontCache(50*1024*1024, 500) // 50MB, 500 items

// NewFontCache creates a new font cache
func NewFontCache(maxSize int64, maxItems int) *FontCache {
	return &FontCache{
		Cache: NewCache(maxSize, maxItems),
	}
}

// PathCache caches rendered paths
type PathCache struct {
	*Cache
}

// GlobalPathCache is the global path cache
var GlobalPathCache = NewPathCache(25*1024*1024, 250) // 25MB, 250 items

// NewPathCache creates a new path cache
func NewPathCache(maxSize int64, maxItems int) *PathCache {
	return &PathCache{
		Cache: NewCache(maxSize, maxItems),
	}
}

// Cache key generation utilities

// GenerateImageKey generates a cache key for an image operation
func GenerateImageKey(operation string, params ...interface{}) string {
	hash := md5.New()
	hash.Write([]byte(operation))
	for _, param := range params {
		hash.Write([]byte(fmt.Sprintf("%v", param)))
	}
	return fmt.Sprintf("img_%x", hash.Sum(nil))
}

// GenerateFontKey generates a cache key for a font operation
func GenerateFontKey(fontPath string, size float64, params ...interface{}) string {
	hash := md5.New()
	hash.Write([]byte(fontPath))
	hash.Write([]byte(fmt.Sprintf("%.2f", size)))
	for _, param := range params {
		hash.Write([]byte(fmt.Sprintf("%v", param)))
	}
	return fmt.Sprintf("font_%x", hash.Sum(nil))
}

// GeneratePathKey generates a cache key for a path operation
func GeneratePathKey(pathData string, transform Matrix) string {
	hash := md5.New()
	hash.Write([]byte(pathData))
	hash.Write([]byte(fmt.Sprintf("%.6f,%.6f,%.6f,%.6f,%.6f,%.6f", 
		transform.XX, transform.XY, transform.X0,
		transform.YX, transform.YY, transform.Y0)))
	return fmt.Sprintf("path_%x", hash.Sum(nil))
}

// Context extensions for caching

// DrawCachedCircle draws a circle using cache if available
func (dc *Context) DrawCachedCircle(x, y, radius float64, fill bool) {
	key := GenerateImageKey("circle", x, y, radius, fill, dc.color)
	
	if cached, exists := GlobalImageCache.GetImage(key); exists {
		dc.DrawImage(cached, 0, 0)
		return
	}
	
	// Create temporary context for rendering
	tempDC := NewContext(int(radius*2+10), int(radius*2+10))
	tempDC.SetColor(dc.color)
	tempDC.DrawCircle(radius+5, radius+5, radius)
	if fill {
		tempDC.Fill()
	} else {
		tempDC.Stroke()
	}
	
	// Cache the result
	GlobalImageCache.SetImage(key, tempDC.im)
	
	// Draw the cached image
	dc.DrawImage(tempDC.im, int(x-radius-5), int(y-radius-5))
}

// DrawCachedText draws text using cache if available
func (dc *Context) DrawCachedText(text string, x, y float64) {
	if dc.fontFace == nil {
		dc.DrawString(text, x, y)
		return
	}
	
	key := GenerateFontKey("text", dc.fontHeight, text, dc.color)
	
	if cached, exists := GlobalImageCache.GetImage(key); exists {
		dc.DrawImage(cached, int(x), int(y))
		return
	}
	
	// Measure text to create appropriately sized context
	width, height := dc.MeasureString(text)
	
	// Create temporary context for rendering
	tempDC := NewContext(int(width+10), int(height+10))
	tempDC.SetFontFace(dc.fontFace)
	tempDC.SetColor(dc.color)
	tempDC.DrawString(text, 5, height+5)
	
	// Cache the result
	GlobalImageCache.SetImage(key, tempDC.im)
	
	// Draw the cached image
	dc.DrawImage(tempDC.im, int(x-5), int(y-height-5))
}

// CacheConfig holds cache configuration
type CacheConfig struct {
	ImageCacheSize int64
	FontCacheSize  int64
	PathCacheSize  int64
	EnableCaching  bool
}

// DefaultCacheConfig returns the default cache configuration
func DefaultCacheConfig() CacheConfig {
	return CacheConfig{
		ImageCacheSize: 100 * 1024 * 1024, // 100MB
		FontCacheSize:  50 * 1024 * 1024,  // 50MB
		PathCacheSize:  25 * 1024 * 1024,  // 25MB
		EnableCaching:  true,
	}
}

// SetCacheConfig configures the global caches
func SetCacheConfig(config CacheConfig) {
	if config.EnableCaching {
		GlobalImageCache = NewImageCache(config.ImageCacheSize, 1000)
		GlobalFontCache = NewFontCache(config.FontCacheSize, 500)
		GlobalPathCache = NewPathCache(config.PathCacheSize, 250)
	} else {
		GlobalImageCache.Clear()
		GlobalFontCache.Clear()
		GlobalPathCache.Clear()
	}
}

// GetCacheStats returns statistics for all caches
func GetCacheStats() map[string]CacheStats {
	return map[string]CacheStats{
		"image": GlobalImageCache.GetStats(),
		"font":  GlobalFontCache.GetStats(),
		"path":  GlobalPathCache.GetStats(),
	}
}

// ClearAllCaches clears all global caches
func ClearAllCaches() {
	GlobalImageCache.Clear()
	GlobalFontCache.Clear()
	GlobalPathCache.Clear()
}
