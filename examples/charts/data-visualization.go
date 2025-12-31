package main

import (
	"fmt"
	"math"
	"math/rand"

	"github.com/GrandpaEJ/advancegg"
)

func main() {
	fmt.Println("Creating data visualization examples...")

	// Create various data visualizations
	createBarChart()
	createLineChart()
	createPieChart()
	createScatterPlot()
	createHeatmap()
	createDashboard()

	fmt.Println("Data visualization examples completed!")
}

func createBarChart() {
	dc := advancegg.NewContext(800, 600)

	// White background
	dc.SetRGB(1, 1, 1)
	dc.Clear()

	// Title
	dc.SetRGB(0, 0, 0)
	dc.DrawString("Monthly Sales Data", 50, 50)

	// Sample data
	data := []struct {
		label string
		value float64
		color [3]float64
	}{
		{"Jan", 120, [3]float64{0.8, 0.2, 0.2}},
		{"Feb", 150, [3]float64{0.2, 0.8, 0.2}},
		{"Mar", 180, [3]float64{0.2, 0.2, 0.8}},
		{"Apr", 200, [3]float64{0.8, 0.8, 0.2}},
		{"May", 160, [3]float64{0.8, 0.2, 0.8}},
		{"Jun", 220, [3]float64{0.2, 0.8, 0.8}},
	}

	// Chart area
	chartX, chartY := 100.0, 100.0
	chartWidth, chartHeight := 600.0, 400.0

	// Draw axes
	dc.SetRGB(0, 0, 0)
	dc.SetLineWidth(2)
	dc.DrawLine(chartX, chartY+chartHeight, chartX+chartWidth, chartY+chartHeight) // X-axis
	dc.DrawLine(chartX, chartY, chartX, chartY+chartHeight)                        // Y-axis
	dc.Stroke()

	// Find max value for scaling
	maxValue := 0.0
	for _, item := range data {
		if item.value > maxValue {
			maxValue = item.value
		}
	}

	// Draw bars
	barWidth := chartWidth / float64(len(data)) * 0.8
	barSpacing := chartWidth / float64(len(data))

	for i, item := range data {
		barHeight := (item.value / maxValue) * chartHeight
		barX := chartX + float64(i)*barSpacing + barSpacing*0.1
		barY := chartY + chartHeight - barHeight

		// Draw bar with shadow
		dc.SetShadowRGBA(2, 2, 4, 0, 0, 0, 0.3)
		dc.SetRGB(item.color[0], item.color[1], item.color[2])
		dc.DrawRoundedRectangleWithShadow(barX, barY, barWidth, barHeight, 5)

		// Draw value on top
		dc.ClearShadow()
		dc.SetRGB(0, 0, 0)
		dc.DrawString(fmt.Sprintf("%.0f", item.value), barX+barWidth/2-10, barY-10)

		// Draw label
		dc.DrawString(item.label, barX+barWidth/2-10, chartY+chartHeight+20)
	}

	// Y-axis labels
	for i := 0; i <= 5; i++ {
		value := maxValue * float64(i) / 5
		y := chartY + chartHeight - (float64(i)/5)*chartHeight
		dc.DrawString(fmt.Sprintf("%.0f", value), chartX-30, y+5)

		// Grid lines
		dc.SetRGB(0.9, 0.9, 0.9)
		dc.SetLineWidth(1)
		dc.DrawLine(chartX, y, chartX+chartWidth, y)
		dc.Stroke()
	}

	dc.SavePNG("images/charts/bar-chart.png")
	fmt.Println("Bar chart saved as bar-chart.png")
}

func createLineChart() {
	dc := advancegg.NewContext(800, 600)

	// White background
	dc.SetRGB(1, 1, 1)
	dc.Clear()

	// Title
	dc.SetRGB(0, 0, 0)
	dc.DrawString("Stock Price Trends", 50, 50)

	// Generate sample data
	dataPoints := 50
	data1 := make([]float64, dataPoints)
	data2 := make([]float64, dataPoints)

	for i := 0; i < dataPoints; i++ {
		t := float64(i) / float64(dataPoints-1)
		data1[i] = 100 + 50*math.Sin(t*4*math.Pi) + rand.Float64()*20
		data2[i] = 120 + 30*math.Cos(t*3*math.Pi) + rand.Float64()*15
	}

	// Chart area
	chartX, chartY := 100.0, 100.0
	chartWidth, chartHeight := 600.0, 400.0

	// Draw axes
	dc.SetRGB(0, 0, 0)
	dc.SetLineWidth(2)
	dc.DrawLine(chartX, chartY+chartHeight, chartX+chartWidth, chartY+chartHeight)
	dc.DrawLine(chartX, chartY, chartX, chartY+chartHeight)
	dc.Stroke()

	// Find min/max for scaling
	minVal, maxVal := data1[0], data1[0]
	for _, val := range data1 {
		if val < minVal {
			minVal = val
		}
		if val > maxVal {
			maxVal = val
		}
	}
	for _, val := range data2 {
		if val < minVal {
			minVal = val
		}
		if val > maxVal {
			maxVal = val
		}
	}

	// Draw grid
	dc.SetRGB(0.9, 0.9, 0.9)
	dc.SetLineWidth(1)
	for i := 0; i <= 10; i++ {
		x := chartX + float64(i)*chartWidth/10
		y := chartY + float64(i)*chartHeight/10
		dc.DrawLine(x, chartY, x, chartY+chartHeight)
		dc.DrawLine(chartX, y, chartX+chartWidth, y)
		dc.Stroke()
	}

	// Draw line 1
	dc.SetRGB(0.8, 0.2, 0.2)
	dc.SetLineWidth(3)
	for i := 0; i < dataPoints-1; i++ {
		x1 := chartX + float64(i)*chartWidth/float64(dataPoints-1)
		y1 := chartY + chartHeight - (data1[i]-minVal)/(maxVal-minVal)*chartHeight
		x2 := chartX + float64(i+1)*chartWidth/float64(dataPoints-1)
		y2 := chartY + chartHeight - (data1[i+1]-minVal)/(maxVal-minVal)*chartHeight

		dc.DrawLine(x1, y1, x2, y2)
		dc.Stroke()
	}

	// Draw line 2
	dc.SetRGB(0.2, 0.2, 0.8)
	for i := 0; i < dataPoints-1; i++ {
		x1 := chartX + float64(i)*chartWidth/float64(dataPoints-1)
		y1 := chartY + chartHeight - (data2[i]-minVal)/(maxVal-minVal)*chartHeight
		x2 := chartX + float64(i+1)*chartWidth/float64(dataPoints-1)
		y2 := chartY + chartHeight - (data2[i+1]-minVal)/(maxVal-minVal)*chartHeight

		dc.DrawLine(x1, y1, x2, y2)
		dc.Stroke()
	}

	// Legend
	dc.SetRGB(0.8, 0.2, 0.2)
	dc.DrawLine(50, 550, 80, 550)
	dc.Stroke()
	dc.SetRGB(0, 0, 0)
	dc.DrawString("Stock A", 90, 555)

	dc.SetRGB(0.2, 0.2, 0.8)
	dc.DrawLine(200, 550, 230, 550)
	dc.Stroke()
	dc.SetRGB(0, 0, 0)
	dc.DrawString("Stock B", 240, 555)

	dc.SavePNG("images/charts/line-chart.png")
	fmt.Println("Line chart saved as line-chart.png")
}

func createPieChart() {
	dc := advancegg.NewContext(600, 600)

	// White background
	dc.SetRGB(1, 1, 1)
	dc.Clear()

	// Title
	dc.SetRGB(0, 0, 0)
	dc.DrawString("Market Share Distribution", 50, 50)

	// Sample data
	data := []struct {
		label string
		value float64
		color [3]float64
	}{
		{"Product A", 35, [3]float64{0.8, 0.2, 0.2}},
		{"Product B", 25, [3]float64{0.2, 0.8, 0.2}},
		{"Product C", 20, [3]float64{0.2, 0.2, 0.8}},
		{"Product D", 15, [3]float64{0.8, 0.8, 0.2}},
		{"Others", 5, [3]float64{0.6, 0.6, 0.6}},
	}

	// Calculate total
	total := 0.0
	for _, item := range data {
		total += item.value
	}

	// Chart center and radius
	centerX, centerY := 300.0, 300.0
	radius := 150.0

	// Draw pie slices
	currentAngle := -math.Pi / 2 // Start at top

	for _, item := range data {
		sliceAngle := (item.value / total) * 2 * math.Pi

		// Draw slice with shadow
		dc.SetShadowRGBA(3, 3, 6, 0, 0, 0, 0.3)
		dc.SetRGB(item.color[0], item.color[1], item.color[2])

		// Create pie slice path
		dc.MoveTo(centerX, centerY)
		dc.DrawArc(centerX, centerY, radius, currentAngle, currentAngle+sliceAngle)
		dc.ClosePath()
		dc.FillWithShadow()

		// Draw slice outline
		dc.ClearShadow()
		dc.SetRGB(1, 1, 1)
		dc.SetLineWidth(2)
		dc.Stroke()

		// Draw label
		labelAngle := currentAngle + sliceAngle/2
		labelX := centerX + math.Cos(labelAngle)*(radius+30)
		labelY := centerY + math.Sin(labelAngle)*(radius+30)

		dc.SetRGB(0, 0, 0)
		dc.DrawString(fmt.Sprintf("%s (%.1f%%)", item.label, item.value/total*100),
			labelX-30, labelY)

		currentAngle += sliceAngle
	}

	dc.SavePNG("images/charts/pie-chart.png")
	fmt.Println("Pie chart saved as pie-chart.png")
}

func createScatterPlot() {
	dc := advancegg.NewContext(800, 600)

	// White background
	dc.SetRGB(1, 1, 1)
	dc.Clear()

	// Title
	dc.SetRGB(0, 0, 0)
	dc.DrawString("Height vs Weight Correlation", 50, 50)

	// Generate sample data
	numPoints := 100
	heights := make([]float64, numPoints)
	weights := make([]float64, numPoints)

	for i := 0; i < numPoints; i++ {
		heights[i] = 150 + rand.Float64()*50                 // 150-200 cm
		weights[i] = heights[i]*0.8 + rand.Float64()*20 - 10 // Correlated with some noise
	}

	// Chart area
	chartX, chartY := 100.0, 100.0
	chartWidth, chartHeight := 600.0, 400.0

	// Draw axes
	dc.SetRGB(0, 0, 0)
	dc.SetLineWidth(2)
	dc.DrawLine(chartX, chartY+chartHeight, chartX+chartWidth, chartY+chartHeight)
	dc.DrawLine(chartX, chartY, chartX, chartY+chartHeight)
	dc.Stroke()

	// Find min/max for scaling
	minHeight, maxHeight := heights[0], heights[0]
	minWeight, maxWeight := weights[0], weights[0]

	for i := 0; i < numPoints; i++ {
		if heights[i] < minHeight {
			minHeight = heights[i]
		}
		if heights[i] > maxHeight {
			maxHeight = heights[i]
		}
		if weights[i] < minWeight {
			minWeight = weights[i]
		}
		if weights[i] > maxWeight {
			maxWeight = weights[i]
		}
	}

	// Draw points
	for i := 0; i < numPoints; i++ {
		x := chartX + (heights[i]-minHeight)/(maxHeight-minHeight)*chartWidth
		y := chartY + chartHeight - (weights[i]-minWeight)/(maxWeight-minWeight)*chartHeight

		// Color based on value
		intensity := (weights[i] - minWeight) / (maxWeight - minWeight)
		dc.SetRGB(intensity, 0.5, 1-intensity)
		dc.DrawCircle(x, y, 4)
		dc.Fill()
	}

	// Axis labels
	dc.SetRGB(0, 0, 0)
	dc.DrawString("Height (cm)", chartX+chartWidth/2-30, chartY+chartHeight+40)

	// Rotate for Y-axis label (simplified)
	dc.DrawString("Weight (kg)", 20, chartY+chartHeight/2)

	dc.SavePNG("images/charts/scatter-plot.png")
	fmt.Println("Scatter plot saved as scatter-plot.png")
}

func createHeatmap() {
	dc := advancegg.NewContext(600, 600)

	// White background
	dc.SetRGB(1, 1, 1)
	dc.Clear()

	// Title
	dc.SetRGB(0, 0, 0)
	dc.DrawString("Temperature Heatmap", 50, 50)

	// Generate sample data
	gridSize := 20
	data := make([][]float64, gridSize)
	for i := range data {
		data[i] = make([]float64, gridSize)
		for j := range data[i] {
			// Create some pattern
			x := float64(i) / float64(gridSize-1)
			y := float64(j) / float64(gridSize-1)
			data[i][j] = math.Sin(x*4*math.Pi)*math.Cos(y*3*math.Pi) +
				rand.Float64()*0.3
		}
	}

	// Find min/max for color scaling
	minVal, maxVal := data[0][0], data[0][0]
	for i := 0; i < gridSize; i++ {
		for j := 0; j < gridSize; j++ {
			if data[i][j] < minVal {
				minVal = data[i][j]
			}
			if data[i][j] > maxVal {
				maxVal = data[i][j]
			}
		}
	}

	// Chart area
	chartX, chartY := 100.0, 100.0
	cellSize := 400.0 / float64(gridSize)

	// Draw heatmap
	for i := 0; i < gridSize; i++ {
		for j := 0; j < gridSize; j++ {
			// Normalize value to 0-1
			normalized := (data[i][j] - minVal) / (maxVal - minVal)

			// Color mapping: blue (cold) to red (hot)
			r := normalized
			g := 0.5 * (1 - math.Abs(normalized-0.5))
			b := 1 - normalized

			dc.SetRGB(r, g, b)
			dc.DrawRectangle(chartX+float64(i)*cellSize, chartY+float64(j)*cellSize,
				cellSize, cellSize)
			dc.Fill()
		}
	}

	// Draw color scale
	scaleX := chartX + 420
	scaleWidth := 30.0
	scaleHeight := 400.0

	for i := 0; i < 100; i++ {
		normalized := float64(i) / 99.0
		r := normalized
		g := 0.5 * (1 - math.Abs(normalized-0.5))
		b := 1 - normalized

		dc.SetRGB(r, g, b)
		dc.DrawRectangle(scaleX, chartY+float64(i)*scaleHeight/100,
			scaleWidth, scaleHeight/100)
		dc.Fill()
	}

	// Scale labels
	dc.SetRGB(0, 0, 0)
	dc.DrawString(fmt.Sprintf("%.2f", maxVal), scaleX+35, chartY+10)
	dc.DrawString(fmt.Sprintf("%.2f", minVal), scaleX+35, chartY+scaleHeight)

	dc.SavePNG("images/charts/heatmap.png")
	fmt.Println("Heatmap saved as heatmap.png")
}

func createDashboard() {
	dc := advancegg.NewContext(1200, 800)

	// Background
	dc.SetRGB(0.95, 0.95, 0.95)
	dc.Clear()

	// Title
	dc.SetRGB(0, 0, 0)
	dc.DrawString("Business Dashboard", 50, 50)

	// Create multiple small charts
	createMiniBarChart(dc, 50, 100, 300, 200)
	createMiniLineChart(dc, 400, 100, 300, 200)
	createMiniPieChart(dc, 750, 100, 200, 200)
	createKPICards(dc, 50, 350)
	createMiniHeatmap(dc, 400, 400, 300, 200)
	createGaugeChart(dc, 750, 400, 200, 200)

	dc.SavePNG("images/charts/dashboard.png")
	fmt.Println("Dashboard saved as dashboard.png")
}

func createMiniBarChart(dc *advancegg.Context, x, y, width, height float64) {
	// Mini bar chart implementation
	dc.SetRGB(1, 1, 1)
	dc.DrawRectangle(x, y, width, height)
	dc.Fill()

	dc.SetRGB(0, 0, 0)
	dc.DrawString("Sales by Quarter", x+10, y+20)

	values := []float64{100, 120, 150, 180}
	colors := [][3]float64{{0.8, 0.2, 0.2}, {0.2, 0.8, 0.2}, {0.2, 0.2, 0.8}, {0.8, 0.8, 0.2}}

	barWidth := (width - 40) / float64(len(values))
	maxVal := 200.0

	for i, val := range values {
		barHeight := (val / maxVal) * (height - 60)
		barX := x + 20 + float64(i)*barWidth
		barY := y + height - 20 - barHeight

		dc.SetRGB(colors[i][0], colors[i][1], colors[i][2])
		dc.DrawRectangle(barX, barY, barWidth*0.8, barHeight)
		dc.Fill()
	}
}

func createMiniLineChart(dc *advancegg.Context, x, y, width, height float64) {
	// Mini line chart implementation
	dc.SetRGB(1, 1, 1)
	dc.DrawRectangle(x, y, width, height)
	dc.Fill()

	dc.SetRGB(0, 0, 0)
	dc.DrawString("Revenue Trend", x+10, y+20)

	// Simple trend line
	dc.SetRGB(0.2, 0.6, 0.8)
	dc.SetLineWidth(3)

	points := 20
	for i := 0; i < points-1; i++ {
		x1 := x + 20 + float64(i)*(width-40)/float64(points-1)
		y1 := y + height/2 + 30*math.Sin(float64(i)*0.5)
		x2 := x + 20 + float64(i+1)*(width-40)/float64(points-1)
		y2 := y + height/2 + 30*math.Sin(float64(i+1)*0.5)

		dc.DrawLine(x1, y1, x2, y2)
		dc.Stroke()
	}
}

func createMiniPieChart(dc *advancegg.Context, x, y, width, height float64) {
	// Mini pie chart implementation
	dc.SetRGB(1, 1, 1)
	dc.DrawRectangle(x, y, width, height)
	dc.Fill()

	dc.SetRGB(0, 0, 0)
	dc.DrawString("Market Share", x+10, y+20)

	centerX := x + width/2
	centerY := y + height/2 + 10
	radius := 60.0

	values := []float64{40, 30, 20, 10}
	colors := [][3]float64{{0.8, 0.2, 0.2}, {0.2, 0.8, 0.2}, {0.2, 0.2, 0.8}, {0.8, 0.8, 0.2}}

	currentAngle := 0.0
	total := 100.0

	for i, val := range values {
		sliceAngle := (val / total) * 2 * math.Pi

		dc.SetRGB(colors[i][0], colors[i][1], colors[i][2])
		dc.MoveTo(centerX, centerY)
		dc.DrawArc(centerX, centerY, radius, currentAngle, currentAngle+sliceAngle)
		dc.ClosePath()
		dc.Fill()

		currentAngle += sliceAngle
	}
}

func createKPICards(dc *advancegg.Context, x, y float64) {
	kpis := []struct {
		title string
		value string
		color [3]float64
	}{
		{"Revenue", "$1.2M", [3]float64{0.2, 0.8, 0.2}},
		{"Users", "45.2K", [3]float64{0.2, 0.2, 0.8}},
		{"Growth", "+12.5%", [3]float64{0.8, 0.6, 0.2}},
	}

	cardWidth := 100.0
	cardHeight := 80.0

	for i, kpi := range kpis {
		cardX := x + float64(i)*120

		// Card background with shadow
		dc.SetShadowRGBA(2, 2, 4, 0, 0, 0, 0.2)
		dc.SetRGB(1, 1, 1)
		dc.DrawRoundedRectangleWithShadow(cardX, y, cardWidth, cardHeight, 8)

		// Color accent
		dc.ClearShadow()
		dc.SetRGB(kpi.color[0], kpi.color[1], kpi.color[2])
		dc.DrawRectangle(cardX, y, cardWidth, 5)
		dc.Fill()

		// Text
		dc.SetRGB(0, 0, 0)
		dc.DrawString(kpi.title, cardX+10, y+25)
		dc.DrawString(kpi.value, cardX+10, y+50)
	}
}

func createMiniHeatmap(dc *advancegg.Context, x, y, width, height float64) {
	dc.SetRGB(1, 1, 1)
	dc.DrawRectangle(x, y, width, height)
	dc.Fill()

	dc.SetRGB(0, 0, 0)
	dc.DrawString("Activity Heatmap", x+10, y+20)

	gridSize := 10
	cellSize := (width - 40) / float64(gridSize)

	for i := 0; i < gridSize; i++ {
		for j := 0; j < gridSize; j++ {
			intensity := rand.Float64()
			dc.SetRGB(intensity, intensity*0.5, 1-intensity)
			dc.DrawRectangle(x+20+float64(i)*cellSize, y+40+float64(j)*cellSize,
				cellSize-1, cellSize-1)
			dc.Fill()
		}
	}
}

func createGaugeChart(dc *advancegg.Context, x, y, width, height float64) {
	dc.SetRGB(1, 1, 1)
	dc.DrawRectangle(x, y, width, height)
	dc.Fill()

	dc.SetRGB(0, 0, 0)
	dc.DrawString("Performance", x+10, y+20)

	centerX := x + width/2
	centerY := y + height/2 + 20
	radius := 60.0

	// Draw gauge background
	dc.SetRGB(0.9, 0.9, 0.9)
	dc.SetLineWidth(10)
	dc.DrawArc(centerX, centerY, radius, math.Pi, 2*math.Pi)
	dc.Stroke()

	// Draw gauge value (75%)
	value := 0.75
	dc.SetRGB(0.2, 0.8, 0.2)
	dc.DrawArc(centerX, centerY, radius, math.Pi, math.Pi+value*math.Pi)
	dc.Stroke()

	// Draw needle
	needleAngle := math.Pi + value*math.Pi
	needleX := centerX + math.Cos(needleAngle)*radius*0.8
	needleY := centerY + math.Sin(needleAngle)*radius*0.8

	dc.SetRGB(1, 0, 0)
	dc.SetLineWidth(3)
	dc.DrawLine(centerX, centerY, needleX, needleY)
	dc.Stroke()

	// Value text
	dc.SetRGB(0, 0, 0)
	dc.DrawString("75%", centerX-10, centerY+30)
}
