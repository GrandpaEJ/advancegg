package core

// SetEnableAutoEmoji enables or disables automatic emoji rendering in text
func (dc *Context) SetEnableAutoEmoji(enable bool) {
	dc.enableAutoEmoji = enable
}

// GetEnableAutoEmoji returns whether automatic emoji rendering is enabled
func (dc *Context) GetEnableAutoEmoji() bool {
	return dc.enableAutoEmoji
}

// SetEmojiSize sets the size for emoji rendering (defaults to font height)
func (dc *Context) SetEmojiSize(size float64) {
	renderer := dc.GetEmojiRenderer()
	renderer.EmojiSize = size
}

// GetEmojiSize returns the current emoji rendering size
func (dc *Context) GetEmojiSize() float64 {
	renderer := dc.GetEmojiRenderer()
	return renderer.EmojiSize
}
