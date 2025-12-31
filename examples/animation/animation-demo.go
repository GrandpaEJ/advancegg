package main

import (
	"fmt"
	"image"
	"image/color"
	"time"

	"github.com/GrandpaEJ/advancegg"
)

func main() {
	fmt.Println("Creating Animation demo...")

	// Create basic animation demo
	createBasicAnimationDemo()

	// Create easing functions demo
	createEasingDemo()

	// Create complex animation sequence demo
	createComplexSequenceDemo()

	fmt.Println("Animation demo completed!")
}

func createBasicAnimationDemo() {
	fmt.Println("Creating basic animation...")

	// Create a simple bouncing ball animation (optimized for faster generation)
	width, height := 400, 300
	fps := 10.0                 // Reduced FPS for faster generation
	duration := 2 * time.Second // Shorter duration

	animator := advancegg.NewAnimator(width, height, fps, duration)
	frameCount := int(fps * duration.Seconds())

	for i := 0; i < frameCount; i++ {
		t := float64(i) / float64(frameCount-1)

		// Create frame
		ctx := advancegg.NewContext(width, height)
		ctx.SetRGB(0.9, 0.9, 1.0)
		ctx.Clear()

		// Draw title
		ctx.SetRGB(0, 0, 0)
		ctx.DrawString("Bouncing Ball Animation", 120, 30)

		// Animate ball position with bounce
		x := advancegg.AnimateProperty(50, 350, t, advancegg.EaseLinear)
		y := 150 + 80*advancegg.EaseBounce(t)

		// Draw ball
		ctx.SetRGB(1, 0.3, 0.3)
		ctx.DrawCircle(x, y, 15)
		ctx.Fill()

		// Draw ground
		ctx.SetRGB(0.3, 0.7, 0.3)
		ctx.DrawRectangle(0, 250, float64(width), 50)
		ctx.Fill()

		// Add frame to animator
		animator.AddFrame(ctx.Image().(*image.RGBA))
	}

	// Save as GIF
	err := animator.SaveGIF("images/animation/bouncing-ball.gif")
	if err != nil {
		fmt.Printf("Error saving GIF: %v\n", err)
	} else {
		fmt.Println("Created bouncing-ball.gif")
	}
}

func createEasingDemo() {
	fmt.Println("Creating easing functions demo...")

	// Create animation showing different easing functions (optimized)
	width, height := 600, 400
	fps := 10.0                 // Reduced FPS
	duration := 2 * time.Second // Shorter duration

	animator := advancegg.NewAnimator(width, height, fps, duration)
	frameCount := int(fps * duration.Seconds())

	// Define easing functions to demonstrate
	easings := []struct {
		name  string
		func_ advancegg.EasingFunc
		color color.RGBA
	}{
		{"Linear", advancegg.EaseLinear, color.RGBA{255, 100, 100, 255}},
		{"Ease In", advancegg.EaseIn, color.RGBA{100, 255, 100, 255}},
		{"Ease Out", advancegg.EaseOut, color.RGBA{100, 100, 255, 255}},
		{"Ease In-Out", advancegg.EaseInOut, color.RGBA{255, 255, 100, 255}},
		{"Bounce", advancegg.EaseBounce, color.RGBA{255, 100, 255, 255}},
		{"Elastic", advancegg.EaseElastic, color.RGBA{100, 255, 255, 255}},
	}

	for i := 0; i < frameCount; i++ {
		t := float64(i) / float64(frameCount-1)

		// Create frame
		ctx := advancegg.NewContext(width, height)
		ctx.SetRGB(1, 1, 1)
		ctx.Clear()

		// Draw title
		ctx.SetRGB(0, 0, 0)
		ctx.DrawString("Easing Functions Comparison", 200, 30)

		// Draw each easing function
		for j, easing := range easings {
			y := 70 + float64(j)*50

			// Draw label
			ctx.SetRGB(0, 0, 0)
			ctx.DrawString(easing.name, 20, y+5)

			// Draw progress bar background
			ctx.SetRGB(0.9, 0.9, 0.9)
			ctx.DrawRectangle(120, y-10, 400, 20)
			ctx.Fill()

			// Calculate eased position
			easedT := easing.func_(t)
			x := 120 + 400*easedT

			// Draw progress
			ctx.SetColor(easing.color)
			ctx.DrawRectangle(120, y-10, 400*easedT, 20)
			ctx.Fill()

			// Draw moving circle
			ctx.SetColor(easing.color)
			ctx.DrawCircle(x, y, 8)
			ctx.Fill()
		}

		// Draw time indicator
		ctx.SetRGB(0, 0, 0)
		timeText := fmt.Sprintf("Time: %.1f%%", t*100)
		ctx.DrawString(timeText, 450, 350)

		// Add frame to animator
		animator.AddFrame(ctx.Image().(*image.RGBA))
	}

	// Save as GIF
	err := animator.SaveGIF("images/animation/easing-demo.gif")
	if err != nil {
		fmt.Printf("Error saving GIF: %v\n", err)
	} else {
		fmt.Println("Created easing-demo.gif")
	}
}

func createComplexSequenceDemo() {
	fmt.Println("Creating complex animation sequence...")

	// Create a complex animation using the sequence builder
	width, height := 500, 400

	// Define the content to animate
	renderContent := func(ctx *advancegg.Context) {
		// Draw a logo-like design
		ctx.SetRGB(0.2, 0.4, 0.8)
		ctx.DrawCircle(250, 150, 60)
		ctx.Fill()

		ctx.SetRGB(0.8, 0.4, 0.2)
		ctx.DrawRectangle(200, 200, 100, 40)
		ctx.Fill()

		ctx.SetRGB(0.4, 0.8, 0.2)
		ctx.MoveTo(200, 280)
		ctx.LineTo(250, 260)
		ctx.LineTo(300, 280)
		ctx.LineTo(280, 320)
		ctx.LineTo(220, 320)
		ctx.ClosePath()
		ctx.Fill()

		// Add text
		ctx.SetRGB(0, 0, 0)
		ctx.DrawString("AdvanceGG", 210, 350)
	}

	// Create animation sequence
	seq := advancegg.NewAnimationSequence()

	// 1. Fade in
	seq.AddStep(1*time.Second, func(t float64, ctx *advancegg.Context) {
		renderContent(ctx)
		// Apply fade effect
		ctx.SetRGBA(1, 1, 1, 1-t)
		ctx.DrawRectangle(0, 0, float64(width), float64(height))
		ctx.Fill()
	}, advancegg.EaseOut)

	// 2. Scale up
	seq.AddStep(1*time.Second, func(t float64, ctx *advancegg.Context) {
		scale := advancegg.AnimateProperty(1.0, 1.2, t, advancegg.EaseInOut)

		ctx.Push()
		ctx.Translate(float64(width)/2, float64(height)/2)
		ctx.Scale(scale, scale)
		ctx.Translate(-float64(width)/2, -float64(height)/2)

		renderContent(ctx)
		ctx.Pop()
	}, advancegg.EaseInOut)

	// 3. Rotate
	seq.AddStep(2*time.Second, func(t float64, ctx *advancegg.Context) {
		angle := advancegg.AnimateProperty(0, 360, t, advancegg.EaseLinear)

		ctx.Push()
		ctx.Translate(float64(width)/2, float64(height)/2)
		ctx.Rotate(angle * 3.14159 / 180) // Convert to radians
		ctx.Translate(-float64(width)/2, -float64(height)/2)

		renderContent(ctx)
		ctx.Pop()
	}, advancegg.EaseLinear)

	// 4. Pulse
	seq.AddStep(1*time.Second, func(t float64, ctx *advancegg.Context) {
		// Create pulsing effect
		pulseT := t * 4 // 4 pulses in 1 second
		scale := 1.0 + 0.1*advancegg.EaseBounce(pulseT-float64(int(pulseT)))

		ctx.Push()
		ctx.Translate(float64(width)/2, float64(height)/2)
		ctx.Scale(scale, scale)
		ctx.Translate(-float64(width)/2, -float64(height)/2)

		renderContent(ctx)
		ctx.Pop()
	}, advancegg.EaseLinear)

	// 5. Slide out
	seq.AddStep(1*time.Second, func(t float64, ctx *advancegg.Context) {
		offset := advancegg.AnimateProperty(0, float64(width), t, advancegg.EaseIn)

		ctx.Push()
		ctx.Translate(offset, 0)

		renderContent(ctx)
		ctx.Pop()
	}, advancegg.EaseIn)

	// Render sequence to animation (optimized)
	fps := 10.0 // Reduced FPS
	animator := seq.Render(width, height, fps)

	// Save as GIF
	err := animator.SaveGIF("images/animation/complex-sequence.gif")
	if err != nil {
		fmt.Printf("Error saving GIF: %v\n", err)
	} else {
		fmt.Println("Created complex-sequence.gif")
	}

	fmt.Printf("Animation has %d frames\n", animator.GetFrameCount())
}

// Demonstrate color animation
func demonstrateColorAnimation() {
	fmt.Println("Creating color animation demo...")

	width, height := 400, 300
	fps := 10.0                 // Reduced FPS
	duration := 2 * time.Second // Shorter duration

	animator := advancegg.NewAnimator(width, height, fps, duration)
	frameCount := int(fps * duration.Seconds())

	startColor := color.RGBA{255, 100, 100, 255} // Red
	endColor := color.RGBA{100, 100, 255, 255}   // Blue

	for i := 0; i < frameCount; i++ {
		t := float64(i) / float64(frameCount-1)

		// Create frame
		ctx := advancegg.NewContext(width, height)
		ctx.SetRGB(1, 1, 1)
		ctx.Clear()

		// Animate color
		currentColor := advancegg.AnimateColor(startColor, endColor, t, advancegg.EaseInOut)

		// Draw animated circle
		ctx.SetColor(currentColor)
		ctx.DrawCircle(200, 150, 80)
		ctx.Fill()

		// Draw title
		ctx.SetRGB(0, 0, 0)
		ctx.DrawString("Color Animation", 150, 50)

		// Add frame to animator
		animator.AddFrame(ctx.Image().(*image.RGBA))
	}

	// Save as GIF
	err := animator.SaveGIF("images/animation/color-animation.gif")
	if err != nil {
		fmt.Printf("Error saving GIF: %v\n", err)
	} else {
		fmt.Println("Created color-animation.gif")
	}
}

// Demonstrate point animation
func demonstratePointAnimation() {
	fmt.Println("Creating point animation demo...")

	width, height := 400, 300
	fps := 10.0                 // Reduced FPS
	duration := 1 * time.Second // Shorter duration

	animator := advancegg.NewAnimator(width, height, fps, duration)
	frameCount := int(fps * duration.Seconds())

	// Define path points
	points := []advancegg.Point{
		{X: 50, Y: 150},
		{X: 150, Y: 50},
		{X: 250, Y: 150},
		{X: 350, Y: 50},
		{X: 350, Y: 250},
		{X: 50, Y: 250},
	}

	for i := 0; i < frameCount; i++ {
		t := float64(i) / float64(frameCount-1)

		// Create frame
		ctx := advancegg.NewContext(width, height)
		ctx.SetRGB(1, 1, 1)
		ctx.Clear()

		// Draw path
		ctx.SetRGB(0.8, 0.8, 0.8)
		ctx.SetLineWidth(2)
		ctx.MoveTo(points[0].X, points[0].Y)
		for j := 1; j < len(points); j++ {
			ctx.LineTo(points[j].X, points[j].Y)
		}
		ctx.Stroke()

		// Animate point along path
		pathT := t * float64(len(points)-1)
		segmentIndex := int(pathT)
		segmentT := pathT - float64(segmentIndex)

		if segmentIndex < len(points)-1 {
			currentPoint := advancegg.AnimatePoint(points[segmentIndex], points[segmentIndex+1], segmentT, advancegg.EaseInOut)

			// Draw moving circle
			ctx.SetRGB(1, 0.3, 0.3)
			ctx.DrawCircle(currentPoint.X, currentPoint.Y, 10)
			ctx.Fill()
		}

		// Draw title
		ctx.SetRGB(0, 0, 0)
		ctx.DrawString("Point Animation", 150, 30)

		// Add frame to animator
		animator.AddFrame(ctx.Image().(*image.RGBA))
	}

	// Save as GIF
	err := animator.SaveGIF("images/animation/point-animation.gif")
	if err != nil {
		fmt.Printf("Error saving GIF: %v\n", err)
	} else {
		fmt.Println("Created point-animation.gif")
	}
}
