package main

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/GrandpaEJ/advancegg"
)

func main() {
	fmt.Println("Creating game graphics examples...")

	// Create various game graphics
	createSpaceScene()
	createPlatformerLevel()
	createParticleSystem()
	createUIElements()
	createSpriteAnimations()

	fmt.Println("Game graphics examples completed!")
}

func createSpaceScene() {
	dc := advancegg.NewContext(1200, 800)

	// Space background
	createStarField(dc)

	// Planets
	createPlanet(dc, 200, 150, 60, [3]float64{0.8, 0.4, 0.2})  // Mars-like
	createPlanet(dc, 800, 300, 80, [3]float64{0.2, 0.4, 0.8})  // Earth-like
	createPlanet(dc, 1000, 600, 40, [3]float64{0.9, 0.9, 0.7}) // Moon-like

	// Spaceship
	createSpaceship(dc, 400, 400)

	// Asteroids
	for i := 0; i < 15; i++ {
		x := rand.Float64() * 1200
		y := rand.Float64() * 800
		size := 5 + rand.Float64()*15
		createAsteroid(dc, x, y, size)
	}

	// Nebula effect
	createNebula(dc, 600, 200, 300)

	// UI overlay
	createSpaceUI(dc)

	dc.SavePNG("images/games/space-scene.png")
	fmt.Println("Space scene saved as space-scene.png")
}

func createStarField(dc *advancegg.Context) {
	// Dark space background
	dc.SetRGB(0.02, 0.02, 0.1)
	dc.Clear()

	// Stars
	for i := 0; i < 500; i++ {
		x := rand.Float64() * 1200
		y := rand.Float64() * 800
		brightness := rand.Float64()
		size := 1 + rand.Float64()*2

		dc.SetRGB(brightness, brightness, brightness)
		dc.DrawCircle(x, y, size)
		dc.Fill()
	}
}

func createPlanet(dc *advancegg.Context, x, y, radius float64, color [3]float64) {
	// Planet shadow/lighting effect
	for i := 0; i < int(radius); i++ {
		r := float64(i)
		alpha := 1.0 - (r / radius)
		lightFactor := 0.3 + 0.7*math.Cos(float64(i)/radius*math.Pi/2)

		dc.SetRGBA(color[0]*lightFactor, color[1]*lightFactor, color[2]*lightFactor, alpha)
		dc.DrawCircle(x, y, radius-r)
		dc.Fill()
	}

	// Planet rings (for some planets)
	if radius > 50 {
		dc.SetRGBA(0.8, 0.8, 0.6, 0.3)
		dc.SetLineWidth(3)
		dc.DrawEllipse(x, y, radius*1.5, radius*0.3)
		dc.Stroke()
		dc.DrawEllipse(x, y, radius*1.8, radius*0.4)
		dc.Stroke()
	}
}

func createSpaceship(dc *advancegg.Context, x, y float64) {
	// Spaceship body
	dc.SetRGB(0.7, 0.7, 0.8)

	// Main hull
	dc.MoveTo(x, y-20)
	dc.LineTo(x-15, y+20)
	dc.LineTo(x-5, y+15)
	dc.LineTo(x+5, y+15)
	dc.LineTo(x+15, y+20)
	dc.ClosePath()
	dc.Fill()

	// Cockpit
	dc.SetRGB(0.3, 0.5, 0.8)
	dc.DrawCircle(x, y-5, 8)
	dc.Fill()

	// Engine glow
	dc.SetRGB(0.2, 0.6, 1.0)
	dc.DrawCircle(x-8, y+18, 3)
	dc.Fill()
	dc.DrawCircle(x+8, y+18, 3)
	dc.Fill()

	// Engine trail
	dc.SetRGBA(0.2, 0.6, 1.0, 0.5)
	dc.DrawLine(x-8, y+18, x-8, y+35)
	dc.DrawLine(x+8, y+18, x+8, y+35)
	dc.SetLineWidth(4)
	dc.Stroke()
}

func createAsteroid(dc *advancegg.Context, x, y, size float64) {
	// Irregular asteroid shape
	dc.SetRGB(0.4, 0.3, 0.2)

	numPoints := 8
	for i := 0; i < numPoints; i++ {
		angle := float64(i) * 2 * math.Pi / float64(numPoints)
		variation := 0.7 + rand.Float64()*0.6
		px := x + math.Cos(angle)*size*variation
		py := y + math.Sin(angle)*size*variation

		if i == 0 {
			dc.MoveTo(px, py)
		} else {
			dc.LineTo(px, py)
		}
	}
	dc.ClosePath()
	dc.Fill()

	// Crater details
	for i := 0; i < 3; i++ {
		craterX := x + (rand.Float64()-0.5)*size
		craterY := y + (rand.Float64()-0.5)*size
		craterSize := size * 0.1 * (1 + rand.Float64())

		dc.SetRGB(0.2, 0.15, 0.1)
		dc.DrawCircle(craterX, craterY, craterSize)
		dc.Fill()
	}
}

func createNebula(dc *advancegg.Context, x, y, size float64) {
	// Nebula cloud effect
	for i := 0; i < 100; i++ {
		cloudX := x + (rand.Float64()-0.5)*size*2
		cloudY := y + (rand.Float64()-0.5)*size*2
		cloudSize := size * 0.1 * (1 + rand.Float64()*2)

		distance := math.Sqrt((cloudX-x)*(cloudX-x) + (cloudY-y)*(cloudY-y))
		alpha := math.Max(0, 1-distance/size) * 0.3

		dc.SetRGBA(0.8, 0.2, 0.6, alpha)
		dc.DrawCircle(cloudX, cloudY, cloudSize)
		dc.Fill()
	}
}

func createSpaceUI(dc *advancegg.Context) {
	// HUD elements
	dc.SetRGB(0, 1, 0)
	dc.SetLineWidth(2)

	// Radar
	dc.DrawCircle(100, 100, 50)
	dc.Stroke()
	dc.DrawLine(50, 100, 150, 100)
	dc.DrawLine(100, 50, 100, 150)
	dc.Stroke()

	// Health bar
	dc.SetRGB(1, 0, 0)
	dc.DrawRectangle(50, 700, 200, 20)
	dc.Fill()
	dc.SetRGB(0, 1, 0)
	dc.DrawRectangle(50, 700, 150, 20) // 75% health
	dc.Fill()

	// Speed indicator
	dc.SetRGB(0, 1, 1)
	dc.DrawString("SPEED: 1250 km/h", 50, 750)
	dc.DrawString("FUEL: 85%", 50, 770)
}

func createPlatformerLevel() {
	dc := advancegg.NewContext(1200, 600)

	// Sky background
	createSkyGradient(dc)

	// Clouds
	createClouds(dc)

	// Platforms
	createPlatforms(dc)

	// Player character
	createPlayer(dc, 100, 400)

	// Enemies
	createEnemy(dc, 400, 450, "goomba")
	createEnemy(dc, 800, 350, "spike")

	// Collectibles
	for i := 0; i < 10; i++ {
		x := 200 + float64(i)*100
		y := 300 + math.Sin(float64(i)*0.5)*50
		createCoin(dc, x, y)
	}

	// Background elements
	createTrees(dc)
	createMountains(dc)

	dc.SavePNG("images/games/platformer-level.png")
	fmt.Println("Platformer level saved as platformer-level.png")
}

func createSkyGradient(dc *advancegg.Context) {
	// Sky gradient from light blue to darker blue
	for y := 0; y < 600; y++ {
		t := float64(y) / 600.0
		r := 0.5 + 0.3*(1-t)
		g := 0.7 + 0.2*(1-t)
		b := 1.0

		dc.SetRGB(r, g, b)
		dc.DrawLine(0, float64(y), 1200, float64(y))
		dc.Stroke()
	}
}

func createClouds(dc *advancegg.Context) {
	clouds := []struct{ x, y, size float64 }{
		{200, 100, 40},
		{500, 80, 60},
		{900, 120, 50},
		{1100, 90, 35},
	}

	for _, cloud := range clouds {
		dc.SetRGBA(1, 1, 1, 0.8)

		// Cloud made of overlapping circles
		dc.DrawCircle(cloud.x, cloud.y, cloud.size)
		dc.Fill()
		dc.DrawCircle(cloud.x-cloud.size*0.6, cloud.y, cloud.size*0.8)
		dc.Fill()
		dc.DrawCircle(cloud.x+cloud.size*0.6, cloud.y, cloud.size*0.8)
		dc.Fill()
		dc.DrawCircle(cloud.x, cloud.y-cloud.size*0.4, cloud.size*0.6)
		dc.Fill()
	}
}

func createPlatforms(dc *advancegg.Context) {
	platforms := []struct{ x, y, width, height float64 }{
		{0, 500, 300, 100},   // Ground
		{400, 450, 200, 20},  // Platform 1
		{700, 350, 150, 20},  // Platform 2
		{950, 400, 250, 100}, // Ground 2
		{300, 300, 100, 20},  // High platform
	}

	for _, platform := range platforms {
		// Platform shadow
		dc.SetRGBA(0, 0, 0, 0.3)
		dc.DrawRectangle(platform.x+3, platform.y+3, platform.width, platform.height)
		dc.Fill()

		// Platform
		dc.SetRGB(0.4, 0.8, 0.2) // Green grass
		dc.DrawRectangle(platform.x, platform.y, platform.width, platform.height)
		dc.Fill()

		// Platform edge
		dc.SetRGB(0.3, 0.6, 0.1)
		dc.DrawRectangle(platform.x, platform.y, platform.width, 5)
		dc.Fill()
	}
}

func createPlayer(dc *advancegg.Context, x, y float64) {
	// Player character (simple)

	// Body
	dc.SetRGB(0.2, 0.4, 0.8)
	dc.DrawRectangle(x-10, y-30, 20, 30)
	dc.Fill()

	// Head
	dc.SetRGB(1, 0.8, 0.6)
	dc.DrawCircle(x, y-35, 8)
	dc.Fill()

	// Hat
	dc.SetRGB(0.8, 0.2, 0.2)
	dc.DrawRectangle(x-8, y-45, 16, 8)
	dc.Fill()

	// Arms
	dc.SetRGB(1, 0.8, 0.6)
	dc.DrawCircle(x-12, y-20, 4)
	dc.Fill()
	dc.DrawCircle(x+12, y-20, 4)
	dc.Fill()

	// Legs
	dc.SetRGB(0.2, 0.4, 0.8)
	dc.DrawRectangle(x-8, y, 6, 15)
	dc.Fill()
	dc.DrawRectangle(x+2, y, 6, 15)
	dc.Fill()
}

func createEnemy(dc *advancegg.Context, x, y float64, enemyType string) {
	switch enemyType {
	case "goomba":
		// Brown mushroom enemy
		dc.SetRGB(0.6, 0.4, 0.2)
		dc.DrawCircle(x, y-10, 15)
		dc.Fill()

		// Eyes
		dc.SetRGB(0, 0, 0)
		dc.DrawCircle(x-5, y-12, 2)
		dc.Fill()
		dc.DrawCircle(x+5, y-12, 2)
		dc.Fill()

		// Feet
		dc.SetRGB(0.4, 0.2, 0.1)
		dc.DrawRectangle(x-10, y, 8, 5)
		dc.Fill()
		dc.DrawRectangle(x+2, y, 8, 5)
		dc.Fill()

	case "spike":
		// Spiky enemy
		dc.SetRGB(0.8, 0.2, 0.2)
		dc.DrawCircle(x, y-8, 12)
		dc.Fill()

		// Spikes
		for i := 0; i < 8; i++ {
			angle := float64(i) * math.Pi / 4
			spikeTipX := x + math.Cos(angle)*18
			spikeTipY := y - 8 + math.Sin(angle)*18

			dc.MoveTo(x, y-8)
			dc.LineTo(spikeTipX, spikeTipY)
			dc.LineTo(x+math.Cos(angle+0.3)*12, y-8+math.Sin(angle+0.3)*12)
			dc.ClosePath()
			dc.Fill()
		}
	}
}

func createCoin(dc *advancegg.Context, x, y float64) {
	// Spinning coin effect
	dc.SetRGB(1, 0.8, 0)
	dc.DrawCircle(x, y, 8)
	dc.Fill()

	// Inner detail
	dc.SetRGB(1, 1, 0.5)
	dc.DrawCircle(x, y, 5)
	dc.Fill()

	// Shine effect
	dc.SetRGB(1, 1, 1)
	dc.DrawCircle(x-2, y-2, 2)
	dc.Fill()
}

func createTrees(dc *advancegg.Context) {
	trees := []struct{ x, y float64 }{
		{50, 500},
		{1150, 400},
		{600, 600},
	}

	for _, tree := range trees {
		// Trunk
		dc.SetRGB(0.4, 0.2, 0.1)
		dc.DrawRectangle(tree.x-5, tree.y-50, 10, 50)
		dc.Fill()

		// Leaves
		dc.SetRGB(0.2, 0.6, 0.1)
		dc.DrawCircle(tree.x, tree.y-60, 25)
		dc.Fill()
		dc.DrawCircle(tree.x-15, tree.y-45, 20)
		dc.Fill()
		dc.DrawCircle(tree.x+15, tree.y-45, 20)
		dc.Fill()
	}
}

func createMountains(dc *advancegg.Context) {
	// Background mountains
	dc.SetRGBA(0.3, 0.3, 0.5, 0.6)

	// Mountain 1
	dc.MoveTo(0, 300)
	dc.LineTo(200, 150)
	dc.LineTo(400, 250)
	dc.LineTo(600, 300)
	dc.LineTo(0, 300)
	dc.Fill()

	// Mountain 2
	dc.MoveTo(500, 300)
	dc.LineTo(700, 100)
	dc.LineTo(900, 200)
	dc.LineTo(1200, 300)
	dc.LineTo(500, 300)
	dc.Fill()
}

func createParticleSystem() {
	dc := advancegg.NewContext(800, 600)

	// Dark background
	dc.SetRGB(0.1, 0.1, 0.2)
	dc.Clear()

	// Fire particles
	createFireEffect(dc, 200, 500, 100)

	// Magic sparkles
	createSparkleEffect(dc, 400, 300, 50)

	// Smoke
	createSmokeEffect(dc, 600, 400, 80)

	dc.SavePNG("images/games/particle-system.png")
	fmt.Println("Particle system saved as particle-system.png")
}

func createFireEffect(dc *advancegg.Context, x, y float64, numParticles int) {
	for i := 0; i < numParticles; i++ {
		// Random position around fire source
		px := x + (rand.Float64()-0.5)*40
		py := y - rand.Float64()*100

		// Fire colors (red to yellow to white)
		heat := rand.Float64()
		r := 1.0
		g := heat * 0.8
		b := heat * heat * 0.3
		alpha := 0.3 + heat*0.7

		size := 2 + rand.Float64()*6

		dc.SetRGBA(r, g, b, alpha)
		dc.DrawCircle(px, py, size)
		dc.Fill()
	}
}

func createSparkleEffect(dc *advancegg.Context, x, y float64, numParticles int) {
	for i := 0; i < numParticles; i++ {
		// Random position in circle
		angle := rand.Float64() * 2 * math.Pi
		distance := rand.Float64() * 80
		px := x + math.Cos(angle)*distance
		py := y + math.Sin(angle)*distance

		// Sparkle colors
		colors := [][3]float64{
			{1, 1, 0.5}, // Gold
			{0.5, 1, 1}, // Cyan
			{1, 0.5, 1}, // Magenta
			{1, 1, 1},   // White
		}

		color := colors[rand.Intn(len(colors))]
		size := 1 + rand.Float64()*3

		dc.SetRGBA(color[0], color[1], color[2], 0.8)

		// Star shape
		for j := 0; j < 4; j++ {
			starAngle := float64(j) * math.Pi / 2
			x1 := px + math.Cos(starAngle)*size
			y1 := py + math.Sin(starAngle)*size
			x2 := px + math.Cos(starAngle+math.Pi)*size
			y2 := py + math.Sin(starAngle+math.Pi)*size

			dc.DrawLine(x1, y1, x2, y2)
			dc.Stroke()
		}
	}
}

func createSmokeEffect(dc *advancegg.Context, x, y float64, numParticles int) {
	for i := 0; i < numParticles; i++ {
		// Smoke rises and spreads
		px := x + (rand.Float64()-0.5)*60
		py := y - rand.Float64()*120

		// Smoke gets lighter as it rises
		height := (y - py) / 120
		alpha := 0.1 + (1-height)*0.4
		gray := 0.3 + height*0.4

		size := 5 + height*15

		dc.SetRGBA(gray, gray, gray, alpha)
		dc.DrawCircle(px, py, size)
		dc.Fill()
	}
}

func createUIElements() {
	dc := advancegg.NewContext(800, 600)

	// Background
	dc.SetRGB(0.2, 0.2, 0.3)
	dc.Clear()

	// Title
	dc.SetRGB(1, 1, 1)
	dc.DrawString("Game UI Elements", 50, 50)

	// Health bar
	createHealthBar(dc, 50, 100, 200, 20, 0.75)

	// Mana bar
	createManaBar(dc, 50, 140, 200, 20, 0.6)

	// Buttons
	createButton(dc, 50, 200, 120, 40, "START", true)
	createButton(dc, 200, 200, 120, 40, "OPTIONS", false)
	createButton(dc, 350, 200, 120, 40, "QUIT", false)

	// Inventory slots
	createInventoryGrid(dc, 50, 280, 8, 4)

	// Minimap
	createMinimap(dc, 500, 100, 200, 150)

	// Chat box
	createChatBox(dc, 500, 300, 250, 200)

	dc.SavePNG("images/games/ui-elements.png")
	fmt.Println("UI elements saved as ui-elements.png")
}

func createHealthBar(dc *advancegg.Context, x, y, width, height, percentage float64) {
	// Background
	dc.SetRGB(0.3, 0.1, 0.1)
	dc.DrawRectangle(x, y, width, height)
	dc.Fill()

	// Health fill
	dc.SetRGB(0.8, 0.2, 0.2)
	dc.DrawRectangle(x+2, y+2, (width-4)*percentage, height-4)
	dc.Fill()

	// Border
	dc.SetRGB(0.6, 0.6, 0.6)
	dc.SetLineWidth(2)
	dc.DrawRectangle(x, y, width, height)
	dc.Stroke()

	// Text
	dc.SetRGB(1, 1, 1)
	dc.DrawString(fmt.Sprintf("HP: %.0f%%", percentage*100), x+5, y+15)
}

func createManaBar(dc *advancegg.Context, x, y, width, height, percentage float64) {
	// Background
	dc.SetRGB(0.1, 0.1, 0.3)
	dc.DrawRectangle(x, y, width, height)
	dc.Fill()

	// Mana fill
	dc.SetRGB(0.2, 0.2, 0.8)
	dc.DrawRectangle(x+2, y+2, (width-4)*percentage, height-4)
	dc.Fill()

	// Border
	dc.SetRGB(0.6, 0.6, 0.6)
	dc.SetLineWidth(2)
	dc.DrawRectangle(x, y, width, height)
	dc.Stroke()

	// Text
	dc.SetRGB(1, 1, 1)
	dc.DrawString(fmt.Sprintf("MP: %.0f%%", percentage*100), x+5, y+15)
}

func createButton(dc *advancegg.Context, x, y, width, height float64, text string, pressed bool) {
	offset := 0.0
	if pressed {
		offset = 2
	}

	// Button shadow
	if !pressed {
		dc.SetRGB(0.1, 0.1, 0.1)
		dc.DrawRoundedRectangle(x+3, y+3, width, height, 5)
		dc.Fill()
	}

	// Button background
	if pressed {
		dc.SetRGB(0.4, 0.4, 0.5)
	} else {
		dc.SetRGB(0.5, 0.5, 0.6)
	}
	dc.DrawRoundedRectangle(x+offset, y+offset, width, height, 5)
	dc.Fill()

	// Button border
	dc.SetRGB(0.7, 0.7, 0.8)
	dc.SetLineWidth(2)
	dc.DrawRoundedRectangle(x+offset, y+offset, width, height, 5)
	dc.Stroke()

	// Button text
	dc.SetRGB(1, 1, 1)
	dc.DrawString(text, x+width/2-20+offset, y+height/2+5+offset)
}

func createInventoryGrid(dc *advancegg.Context, x, y float64, cols, rows int) {
	slotSize := 40.0

	for row := 0; row < rows; row++ {
		for col := 0; col < cols; col++ {
			slotX := x + float64(col)*(slotSize+2)
			slotY := y + float64(row)*(slotSize+2)

			// Slot background
			dc.SetRGB(0.3, 0.3, 0.4)
			dc.DrawRectangle(slotX, slotY, slotSize, slotSize)
			dc.Fill()

			// Slot border
			dc.SetRGB(0.5, 0.5, 0.6)
			dc.SetLineWidth(1)
			dc.DrawRectangle(slotX, slotY, slotSize, slotSize)
			dc.Stroke()

			// Random item (sometimes)
			if rand.Float64() < 0.3 {
				itemColors := [][3]float64{
					{0.8, 0.2, 0.2}, // Red potion
					{0.2, 0.8, 0.2}, // Green potion
					{0.2, 0.2, 0.8}, // Blue potion
					{0.8, 0.8, 0.2}, // Gold coin
				}

				color := itemColors[rand.Intn(len(itemColors))]
				dc.SetRGB(color[0], color[1], color[2])
				dc.DrawCircle(slotX+slotSize/2, slotY+slotSize/2, slotSize/3)
				dc.Fill()
			}
		}
	}
}

func createMinimap(dc *advancegg.Context, x, y, width, height float64) {
	// Minimap background
	dc.SetRGBA(0, 0, 0, 0.7)
	dc.DrawRectangle(x, y, width, height)
	dc.Fill()

	// Minimap border
	dc.SetRGB(0.6, 0.6, 0.6)
	dc.SetLineWidth(2)
	dc.DrawRectangle(x, y, width, height)
	dc.Stroke()

	// Terrain features
	dc.SetRGB(0.2, 0.6, 0.2) // Green for land
	dc.DrawRectangle(x+10, y+10, width-20, height-20)
	dc.Fill()

	// Water
	dc.SetRGB(0.2, 0.2, 0.8)
	dc.DrawCircle(x+width*0.3, y+height*0.7, 20)
	dc.Fill()

	// Mountains
	dc.SetRGB(0.5, 0.4, 0.3)
	dc.MoveTo(x+width*0.7, y+height*0.8)
	dc.LineTo(x+width*0.8, y+height*0.4)
	dc.LineTo(x+width*0.9, y+height*0.8)
	dc.ClosePath()
	dc.Fill()

	// Player position
	dc.SetRGB(1, 0, 0)
	dc.DrawCircle(x+width*0.5, y+height*0.6, 3)
	dc.Fill()

	// Title
	dc.SetRGB(1, 1, 1)
	dc.DrawString("Map", x+5, y+15)
}

func createChatBox(dc *advancegg.Context, x, y, width, height float64) {
	// Chat background
	dc.SetRGBA(0, 0, 0, 0.8)
	dc.DrawRectangle(x, y, width, height)
	dc.Fill()

	// Chat border
	dc.SetRGB(0.5, 0.5, 0.5)
	dc.SetLineWidth(1)
	dc.DrawRectangle(x, y, width, height)
	dc.Stroke()

	// Chat messages
	messages := []struct {
		player string
		text   string
		color  [3]float64
	}{
		{"Player1", "Hello everyone!", [3]float64{0.8, 0.8, 0.8}},
		{"Player2", "Ready for the raid?", [3]float64{0.2, 0.8, 0.2}},
		{"Player3", "Let's go!", [3]float64{0.8, 0.2, 0.8}},
		{"System", "Player4 joined the game", [3]float64{1, 1, 0}},
	}

	for i, msg := range messages {
		msgY := y + 20 + float64(i)*25

		// Player name
		dc.SetRGB(msg.color[0], msg.color[1], msg.color[2])
		dc.DrawString(msg.player+":", x+5, msgY)

		// Message text
		dc.SetRGB(1, 1, 1)
		dc.DrawString(msg.text, x+80, msgY)
	}

	// Input field
	inputY := y + height - 30
	dc.SetRGB(0.2, 0.2, 0.2)
	dc.DrawRectangle(x+5, inputY, width-10, 20)
	dc.Fill()

	dc.SetRGB(0.5, 0.5, 0.5)
	dc.DrawRectangle(x+5, inputY, width-10, 20)
	dc.Stroke()

	dc.SetRGB(0.8, 0.8, 0.8)
	dc.DrawString("Type message...", x+10, inputY+15)
}

func createSpriteAnimations() {
	// Create a sprite sheet showing animation frames
	dc := advancegg.NewContext(800, 200)

	// Background
	dc.SetRGB(0.9, 0.9, 0.9)
	dc.Clear()

	// Title
	dc.SetRGB(0, 0, 0)
	dc.DrawString("Character Animation Frames", 50, 30)

	// Walking animation frames
	for frame := 0; frame < 8; frame++ {
		x := 50 + float64(frame)*90
		y := 100.0

		createAnimationFrame(dc, x, y, frame)

		// Frame number
		dc.SetRGB(0, 0, 0)
		dc.DrawString(fmt.Sprintf("Frame %d", frame+1), x-10, y+60)
	}

	dc.SavePNG("images/games/sprite-animations.png")
	fmt.Println("Sprite animations saved as sprite-animations.png")
}

func createAnimationFrame(dc *advancegg.Context, x, y float64, frame int) {
	// Character body (constant)
	dc.SetRGB(0.2, 0.4, 0.8)
	dc.DrawRectangle(x-8, y-25, 16, 25)
	dc.Fill()

	// Head (constant)
	dc.SetRGB(1, 0.8, 0.6)
	dc.DrawCircle(x, y-30, 6)
	dc.Fill()

	// Walking animation - legs
	legOffset := math.Sin(float64(frame)*math.Pi/4) * 5

	// Left leg
	dc.SetRGB(0.2, 0.4, 0.8)
	dc.DrawRectangle(x-6, y, 5, 12+legOffset)
	dc.Fill()

	// Right leg
	dc.DrawRectangle(x+1, y, 5, 12-legOffset)
	dc.Fill()

	// Arms swinging
	armOffset := math.Sin(float64(frame)*math.Pi/4) * 3

	// Left arm
	dc.SetRGB(1, 0.8, 0.6)
	dc.DrawCircle(x-10, y-18+armOffset, 3)
	dc.Fill()

	// Right arm
	dc.DrawCircle(x+10, y-18-armOffset, 3)
	dc.Fill()

	// Frame border
	dc.SetRGB(0.5, 0.5, 0.5)
	dc.SetLineWidth(1)
	dc.DrawRectangle(x-20, y-45, 40, 70)
	dc.Stroke()
}
