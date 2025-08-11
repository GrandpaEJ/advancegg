package core

import (
	"image"
	"image/color"
	"image/gif"
	"math"
	"os"
	"time"
)

// Animation support for smooth transitions and effects

// Animator manages frame-based animations
type Animator struct {
	width, height int
	fps           float64
	duration      time.Duration
	frames        []*image.RGBA
	currentFrame  int
}

// AnimationFrame represents a single frame in an animation
type AnimationFrame struct {
	Image *image.RGBA
	Delay time.Duration
}

// Easing functions for smooth animations
type EasingFunc func(t float64) float64

// NewAnimator creates a new animator
func NewAnimator(width, height int, fps float64, duration time.Duration) *Animator {
	frameCount := int(fps * duration.Seconds())
	return &Animator{
		width:    width,
		height:   height,
		fps:      fps,
		duration: duration,
		frames:   make([]*image.RGBA, 0, frameCount),
	}
}

// AddFrame adds a frame to the animation
func (a *Animator) AddFrame(frame *image.RGBA) {
	a.frames = append(a.frames, frame)
}

// GetFrameCount returns the number of frames
func (a *Animator) GetFrameCount() int {
	return len(a.frames)
}

// GetFrame returns a specific frame
func (a *Animator) GetFrame(index int) *image.RGBA {
	if index >= 0 && index < len(a.frames) {
		return a.frames[index]
	}
	return nil
}

// SaveGIF saves the animation as a GIF
func (a *Animator) SaveGIF(filename string) error {
	if len(a.frames) == 0 {
		return nil
	}

	// Convert frames to paletted images for GIF
	images := make([]*image.Paletted, len(a.frames))
	delays := make([]int, len(a.frames))

	frameDelay := int(100 / a.fps) // Delay in 1/100th of a second

	for i, frame := range a.frames {
		// Create palette
		palette := make(color.Palette, 256)
		for j := 0; j < 256; j++ {
			palette[j] = color.RGBA{uint8(j), uint8(j), uint8(j), 255}
		}

		// Convert to paletted image
		bounds := frame.Bounds()
		paletted := image.NewPaletted(bounds, palette)

		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				paletted.Set(x, y, frame.At(x, y))
			}
		}

		images[i] = paletted
		delays[i] = frameDelay
	}

	// Create GIF
	anim := &gif.GIF{
		Image: images,
		Delay: delays,
	}

	// Save GIF file
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	return gif.EncodeAll(file, anim)
}

// Animation helpers

// AnimateProperty animates a single property over time
func AnimateProperty(start, end, t float64, easing EasingFunc) float64 {
	if easing != nil {
		t = easing(t)
	}
	return start + (end-start)*t
}

// AnimateColor animates between two colors
func AnimateColor(start, end color.RGBA, t float64, easing EasingFunc) color.RGBA {
	if easing != nil {
		t = easing(t)
	}

	return color.RGBA{
		R: uint8(float64(start.R) + (float64(end.R)-float64(start.R))*t),
		G: uint8(float64(start.G) + (float64(end.G)-float64(start.G))*t),
		B: uint8(float64(start.B) + (float64(end.B)-float64(start.B))*t),
		A: uint8(float64(start.A) + (float64(end.A)-float64(start.A))*t),
	}
}

// AnimatePoint animates between two points
func AnimatePoint(start, end Point, t float64, easing EasingFunc) Point {
	if easing != nil {
		t = easing(t)
	}

	return Point{
		X: start.X + (end.X-start.X)*t,
		Y: start.Y + (end.Y-start.Y)*t,
	}
}

// Easing functions

// Linear easing (no easing)
func EaseLinear(t float64) float64 {
	return t
}

// Ease in (slow start)
func EaseIn(t float64) float64 {
	return t * t
}

// Ease out (slow end)
func EaseOut(t float64) float64 {
	return 1 - (1-t)*(1-t)
}

// Ease in-out (slow start and end)
func EaseInOut(t float64) float64 {
	if t < 0.5 {
		return 2 * t * t
	}
	return 1 - 2*(1-t)*(1-t)
}

// Ease in cubic
func EaseInCubic(t float64) float64 {
	return t * t * t
}

// Ease out cubic
func EaseOutCubic(t float64) float64 {
	return 1 - (1-t)*(1-t)*(1-t)
}

// Ease in-out cubic
func EaseInOutCubic(t float64) float64 {
	if t < 0.5 {
		return 4 * t * t * t
	}
	return 1 - 4*(1-t)*(1-t)*(1-t)
}

// Bounce easing
func EaseBounce(t float64) float64 {
	if t < 1/2.75 {
		return 7.5625 * t * t
	} else if t < 2/2.75 {
		t -= 1.5 / 2.75
		return 7.5625*t*t + 0.75
	} else if t < 2.5/2.75 {
		t -= 2.25 / 2.75
		return 7.5625*t*t + 0.9375
	} else {
		t -= 2.625 / 2.75
		return 7.5625*t*t + 0.984375
	}
}

// Elastic easing
func EaseElastic(t float64) float64 {
	if t == 0 || t == 1 {
		return t
	}
	return -math.Pow(2, 10*(t-1)) * math.Sin((t-1.1)*5*math.Pi)
}

// Animation sequence builder

// AnimationSequence represents a sequence of animations
type AnimationSequence struct {
	steps []AnimationStep
}

// AnimationStep represents a single step in an animation sequence
type AnimationStep struct {
	duration time.Duration
	render   func(t float64, ctx *Context)
	easing   EasingFunc
}

// NewAnimationSequence creates a new animation sequence
func NewAnimationSequence() *AnimationSequence {
	return &AnimationSequence{
		steps: make([]AnimationStep, 0),
	}
}

// AddStep adds a step to the animation sequence
func (seq *AnimationSequence) AddStep(duration time.Duration, render func(t float64, ctx *Context), easing EasingFunc) *AnimationSequence {
	seq.steps = append(seq.steps, AnimationStep{
		duration: duration,
		render:   render,
		easing:   easing,
	})
	return seq
}

// Render renders the animation sequence to frames
func (seq *AnimationSequence) Render(width, height int, fps float64) *Animator {
	totalDuration := time.Duration(0)
	for _, step := range seq.steps {
		totalDuration += step.duration
	}

	animator := NewAnimator(width, height, fps, totalDuration)
	frameCount := int(fps * totalDuration.Seconds())
	frameDuration := totalDuration / time.Duration(frameCount)

	for i := 0; i < frameCount; i++ {
		currentTime := time.Duration(i) * frameDuration

		// Find which step we're in
		stepTime := time.Duration(0)
		for _, step := range seq.steps {
			if currentTime < stepTime+step.duration {
				// We're in this step
				t := float64(currentTime-stepTime) / float64(step.duration)
				if step.easing != nil {
					t = step.easing(t)
				}

				// Render frame
				ctx := NewContext(width, height)
				ctx.SetRGB(1, 1, 1)
				ctx.Clear()

				step.render(t, ctx)

				animator.AddFrame(ctx.Image().(*image.RGBA))
				break
			}
			stepTime += step.duration
		}
	}

	return animator
}

// Predefined animations

// FadeIn creates a fade-in animation
func FadeIn(duration time.Duration, render func(ctx *Context)) *AnimationSequence {
	seq := NewAnimationSequence()
	seq.AddStep(duration, func(t float64, ctx *Context) {
		render(ctx)
		// Apply fade effect by drawing a semi-transparent overlay
		ctx.SetRGBA(1, 1, 1, 1-t)
		ctx.DrawRectangle(0, 0, float64(ctx.Width()), float64(ctx.Height()))
		ctx.Fill()
	}, EaseOut)
	return seq
}

// FadeOut creates a fade-out animation
func FadeOut(duration time.Duration, render func(ctx *Context)) *AnimationSequence {
	seq := NewAnimationSequence()
	seq.AddStep(duration, func(t float64, ctx *Context) {
		render(ctx)
		// Apply fade effect by drawing a semi-transparent overlay
		ctx.SetRGBA(1, 1, 1, t)
		ctx.DrawRectangle(0, 0, float64(ctx.Width()), float64(ctx.Height()))
		ctx.Fill()
	}, EaseIn)
	return seq
}

// SlideIn creates a slide-in animation
func SlideIn(duration time.Duration, direction string, render func(ctx *Context)) *AnimationSequence {
	seq := NewAnimationSequence()
	seq.AddStep(duration, func(t float64, ctx *Context) {
		ctx.Push()

		switch direction {
		case "left":
			ctx.Translate(-float64(ctx.Width())*(1-t), 0)
		case "right":
			ctx.Translate(float64(ctx.Width())*(1-t), 0)
		case "up":
			ctx.Translate(0, -float64(ctx.Height())*(1-t))
		case "down":
			ctx.Translate(0, float64(ctx.Height())*(1-t))
		}

		render(ctx)
		ctx.Pop()
	}, EaseOut)
	return seq
}

// AnimScale creates a scale animation
func AnimScale(duration time.Duration, fromScale, toScale float64, render func(ctx *Context)) *AnimationSequence {
	seq := NewAnimationSequence()
	seq.AddStep(duration, func(t float64, ctx *Context) {
		scale := AnimateProperty(fromScale, toScale, t, EaseInOut)

		ctx.Push()
		ctx.Translate(float64(ctx.Width())/2, float64(ctx.Height())/2)
		ctx.Scale(scale, scale)
		ctx.Translate(-float64(ctx.Width())/2, -float64(ctx.Height())/2)

		render(ctx)
		ctx.Pop()
	}, EaseInOut)
	return seq
}

// AnimRotate creates a rotation animation
func AnimRotate(duration time.Duration, fromAngle, toAngle float64, render func(ctx *Context)) *AnimationSequence {
	seq := NewAnimationSequence()
	seq.AddStep(duration, func(t float64, ctx *Context) {
		angle := AnimateProperty(fromAngle, toAngle, t, EaseLinear)

		ctx.Push()
		ctx.Translate(float64(ctx.Width())/2, float64(ctx.Height())/2)
		ctx.Rotate(angle)
		ctx.Translate(-float64(ctx.Width())/2, -float64(ctx.Height())/2)

		render(ctx)
		ctx.Pop()
	}, EaseLinear)
	return seq
}

// Pulse creates a pulsing animation
func Pulse(duration time.Duration, minScale, maxScale float64, render func(ctx *Context)) *AnimationSequence {
	seq := NewAnimationSequence()
	seq.AddStep(duration/2, func(t float64, ctx *Context) {
		scale := AnimateProperty(minScale, maxScale, t, EaseInOut)

		ctx.Push()
		ctx.Translate(float64(ctx.Width())/2, float64(ctx.Height())/2)
		ctx.Scale(scale, scale)
		ctx.Translate(-float64(ctx.Width())/2, -float64(ctx.Height())/2)

		render(ctx)
		ctx.Pop()
	}, EaseInOut)

	seq.AddStep(duration/2, func(t float64, ctx *Context) {
		scale := AnimateProperty(maxScale, minScale, t, EaseInOut)

		ctx.Push()
		ctx.Translate(float64(ctx.Width())/2, float64(ctx.Height())/2)
		ctx.Scale(scale, scale)
		ctx.Translate(-float64(ctx.Width())/2, -float64(ctx.Height())/2)

		render(ctx)
		ctx.Pop()
	}, EaseInOut)

	return seq
}
