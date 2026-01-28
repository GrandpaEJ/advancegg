// Package advancegg - CLI utilities and command-line interface
package advancegg

import (
	"flag"
	"fmt"
	"image/color"
	"path/filepath"
	"strings"

	"github.com/GrandpaEJ/advancegg/internal/core"
)

// CLI represents the command-line interface
type CLI struct {
	verbose bool
	output  string
}

// NewCLI creates a new CLI instance
func NewCLI() *CLI {
	return &CLI{}
}

// ParseFlags parses command-line flags
func (cli *CLI) ParseFlags() {
	flag.BoolVar(&cli.verbose, "v", false, "verbose output")
	flag.BoolVar(&cli.verbose, "verbose", false, "verbose output")
	flag.StringVar(&cli.output, "o", "", "output file")
	flag.StringVar(&cli.output, "output", "", "output file")
	flag.Parse()
}

// Version returns the library version
func Version() string {
	return "1.6.0"
}

// Info prints library information
func (cli *CLI) Info() {
	fmt.Printf("AdvanceGG v%s\n", Version())
	fmt.Println("A powerful 2D graphics library for Go")
	fmt.Println("https://github.com/GrandpaEJ/advancegg")
}

// ConvertImage converts an image file to PNG using AdvanceGG
func (cli *CLI) ConvertImage(inputPath, outputPath string) error {
	if cli.verbose {
		fmt.Printf("Converting %s to %s\n", inputPath, outputPath)
	}

	// Load the input image
	img, err := core.LoadImage(inputPath)
	if err != nil {
		return fmt.Errorf("failed to load image: %v", err)
	}

	// Create a context from the image
	dc := core.NewContextForImage(img)

	// Save as PNG
	err = dc.SavePNG(outputPath)
	if err != nil {
		return fmt.Errorf("failed to save PNG: %v", err)
	}

	if cli.verbose {
		fmt.Printf("Successfully converted %s to %s\n", inputPath, outputPath)
	}

	return nil
}

// CreateSample creates a sample image demonstrating AdvanceGG capabilities
func (cli *CLI) CreateSample(outputPath string) error {
	if cli.verbose {
		fmt.Printf("Creating sample image: %s\n", outputPath)
	}

	dc := core.NewContext(800, 600)

	// White background
	dc.SetRGB(1, 1, 1)
	dc.Clear()

	// Draw some shapes to demonstrate capabilities
	// Blue circle
	dc.SetRGB(0, 0, 1)
	dc.DrawCircle(200, 150, 80)
	dc.Fill()

	// Red rectangle
	dc.SetRGB(1, 0, 0)
	dc.DrawRectangle(350, 100, 150, 100)
	dc.Fill()

	// Green triangle (using path)
	dc.SetRGB(0, 1, 0)
	dc.MoveTo(600, 100)
	dc.LineTo(550, 200)
	dc.LineTo(650, 200)
	dc.ClosePath()
	dc.Fill()

	// Draw some text
	dc.SetRGB(0, 0, 0)
	dc.DrawString("AdvanceGG Sample", 50, 300)
	dc.DrawString("Circles, Rectangles, and Paths", 50, 350)

	// Draw a gradient rectangle
	gradient := core.NewLinearGradient(50, 400, 250, 400)
	gradient.AddColorStop(0, color.RGBA{R: 255, G: 0, B: 0, A: 255})
	gradient.AddColorStop(1, color.RGBA{R: 0, G: 0, B: 255, A: 255})
	dc.SetFillStyle(gradient)
	dc.DrawRectangle(50, 400, 200, 50)
	dc.Fill()

	// Save the image
	err := dc.SavePNG(outputPath)
	if err != nil {
		return fmt.Errorf("failed to save sample image: %v", err)
	}

	if cli.verbose {
		fmt.Printf("Sample image created: %s\n", outputPath)
	}

	return nil
}

// Run executes the CLI with the given arguments
func (cli *CLI) Run(args []string) error {
	if len(args) == 0 {
		cli.Info()
		fmt.Println("\nUsage:")
		fmt.Println("  advancegg sample [output.png]    - Create a sample image")
		fmt.Println("  advancegg convert input output   - Convert image to PNG")
		fmt.Println("  advancegg version                - Show version")
		fmt.Println("  advancegg help                   - Show this help")
		return nil
	}

	command := args[0]
	switch command {
	case "version", "-v", "--version":
		fmt.Printf("AdvanceGG v%s\n", Version())
		return nil

	case "help", "-h", "--help":
		cli.Info()
		fmt.Println("\nCommands:")
		fmt.Println("  sample [output.png]    - Create a sample image")
		fmt.Println("  convert input output   - Convert image to PNG")
		fmt.Println("  version               - Show version")
		fmt.Println("  help                  - Show this help")
		return nil

	case "sample":
		outputPath := "sample.png"
		if len(args) > 1 {
			outputPath = args[1]
		}
		return cli.CreateSample(outputPath)

	case "convert":
		if len(args) < 3 {
			return fmt.Errorf("convert command requires input and output paths")
		}
		return cli.ConvertImage(args[1], args[2])

	default:
		return fmt.Errorf("unknown command: %s", command)
	}
}

// GetOutputPath returns the output path, generating one if not specified
func (cli *CLI) GetOutputPath(inputPath string) string {
	if cli.output != "" {
		return cli.output
	}

	// Generate output path based on input
	ext := filepath.Ext(inputPath)
	base := strings.TrimSuffix(inputPath, ext)
	return base + "_converted.png"
}
