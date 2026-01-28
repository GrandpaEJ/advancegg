package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os/exec"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/GrandpaEJ/advancegg"
)

type BenchmarkResult struct {
	Name    string
	NsPerOp float64
}

func main() {
	fmt.Println("Running benchmarks...")
	cmd := exec.Command("go", "test", "-bench=.", "./benchmarks")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}

	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(stdout)
	re := regexp.MustCompile(`^(Benchmark\w+(?:/[\w-]+)?)(?:-\d+)?\s+\d+\s+(\d+)\s+ns/op`)

	var results []BenchmarkResult

	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(line) // Stream output to console
		matches := re.FindStringSubmatch(line)
		if len(matches) == 3 {
			name := matches[1]
			ns, _ := strconv.ParseFloat(matches[2], 64)

			// Clean up name
			name = strings.TrimPrefix(name, "Benchmark")

			results = append(results, BenchmarkResult{
				Name:    name,
				NsPerOp: ns,
			})
		}
	}

	if err := cmd.Wait(); err != nil {
		log.Printf("Benchmark command finished with error: %v", err)
	}

	if len(results) == 0 {
		log.Fatal("No benchmark results found")
	}

	// Sort results by ns/op (descending for horizontal bar chart top-to-bottom)
	sort.Slice(results, func(i, j int) bool {
		return results[i].NsPerOp > results[j].NsPerOp
	})

	createChart(results)
}

func createChart(results []BenchmarkResult) {
	width := 1200
	height := len(results)*40 + 100 // Dynamic height based on number of results
	dc := advancegg.NewContext(width, height)

	// Background
	dc.SetRGB(1, 1, 1) // White
	dc.Clear()

	// Title
	dc.SetRGB(0, 0, 0)
	// dc.SetFontFace(basicfont.Face7x13) // Use default font for now
	dc.DrawString("Benchmark Results (ns/op) - Lower is Better", 50, 30)

	// Layout
	margin := 50.0
	labelWidth := 250.0
	chartWidth := float64(width) - labelWidth - margin*2
	chartHeight := float64(height) - 100
	chartX := margin + labelWidth
	chartY := 50.0

	// Find max value
	maxVal := 0.0
	for _, r := range results {
		if r.NsPerOp > maxVal {
			maxVal = r.NsPerOp
		}
	}

	// Draw axes
	dc.SetRGB(0.2, 0.2, 0.2)
	dc.SetLineWidth(1)
	dc.DrawLine(chartX, chartY, chartX, chartY+chartHeight)                        // Y-axis
	dc.DrawLine(chartX, chartY+chartHeight, chartX+chartWidth, chartY+chartHeight) // X-axis
	dc.Stroke()

	// Draw bars
	barHeight := 30.0
	barSpacing := 40.0

	for i, r := range results {
		y := chartY + float64(i)*barSpacing

		// Bar width
		barW := (r.NsPerOp / maxVal) * chartWidth

		// Unique color per bar
		hue := float64(i) / float64(len(results))
		rVal, gVal, bVal := hsvToRGB(hue, 0.7, 0.9)
		dc.SetRGB(rVal, gVal, bVal)

		dc.DrawRectangle(chartX, y, barW, barHeight)
		dc.Fill()

		// Draw label
		dc.SetRGB(0, 0, 0)
		dc.DrawString(r.Name, margin, y+barHeight/2+5)

		// Draw value
		valStr := fmt.Sprintf("%.0f ns", r.NsPerOp)
		dc.DrawString(valStr, chartX+barW+10, y+barHeight/2+5)
	}

	outputFile := "benchmark_results.png"
	if err := dc.SavePNG(outputFile); err != nil {
		log.Fatalf("Failed to save image: %v", err)
	}
	fmt.Printf("Benchmark chart saved to %s\n", outputFile)
}

// hsvToRGB converts HSV color space to RGB
func hsvToRGB(h, s, v float64) (r, g, b float64) {
	i := math.Floor(h * 6)
	f := h*6 - i
	p := v * (1 - s)
	q := v * (1 - f*s)
	t := v * (1 - (1-f)*s)

	switch int(i) % 6 {
	case 0:
		r, g, b = v, t, p
	case 1:
		r, g, b = q, v, p
	case 2:
		r, g, b = p, v, t
	case 3:
		r, g, b = p, q, v
	case 4:
		r, g, b = t, p, v
	case 5:
		r, g, b = v, p, q
	}
	return
}
