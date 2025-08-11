package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
)

// Font setup script for AdvanceGG
// Copies essential fonts from system directories to assets/fonts/

func main() {
	fmt.Println("Setting up fonts for AdvanceGG...")

	// Create assets/fonts directory
	assetsDir := "assets/fonts"
	err := os.MkdirAll(assetsDir, 0755)
	if err != nil {
		fmt.Printf("Error creating assets directory: %v\n", err)
		return
	}

	// Define fonts to copy based on OS
	var fontSources []FontSource
	
	switch runtime.GOOS {
	case "linux":
		fontSources = []FontSource{
			{"/usr/share/fonts/truetype/noto/NotoSans-Regular.ttf", "NotoSans-Regular.ttf"},
			{"/usr/share/fonts/truetype/noto/NotoSans-Bold.ttf", "NotoSans-Bold.ttf"},
			{"/usr/share/fonts/truetype/noto/NotoSerif-Regular.ttf", "NotoSerif-Regular.ttf"},
			{"/usr/share/fonts/truetype/noto/NotoColorEmoji.ttf", "NotoColorEmoji.ttf"},
			{"/usr/share/fonts/truetype/noto/NotoSansArabic-Regular.ttf", "NotoSansArabic-Regular.ttf"},
			{"/usr/share/fonts/truetype/noto/NotoSansHebrew-Regular.ttf", "NotoSansHebrew-Regular.ttf"},
			{"/usr/share/fonts/truetype/noto/NotoSansDevanagari-Regular.ttf", "NotoSansDevanagari-Regular.ttf"},
			{"/usr/share/fonts/truetype/noto/NotoSansThai-Regular.ttf", "NotoSansThai-Regular.ttf"},
			{"/usr/share/fonts/truetype/liberation/LiberationSans-Regular.ttf", "LiberationSans-Regular.ttf"},
		}
	case "darwin": // macOS
		fontSources = []FontSource{
			{"/System/Library/Fonts/Helvetica.ttc", "Helvetica.ttc"},
			{"/System/Library/Fonts/Times.ttc", "Times.ttc"},
			{"/System/Library/Fonts/Apple Color Emoji.ttc", "AppleColorEmoji.ttc"},
		}
	case "windows":
		fontSources = []FontSource{
			{"C:/Windows/Fonts/arial.ttf", "arial.ttf"},
			{"C:/Windows/Fonts/times.ttf", "times.ttf"},
			{"C:/Windows/Fonts/seguiemj.ttf", "seguiemj.ttf"},
		}
	default:
		fmt.Printf("Unsupported OS: %s\n", runtime.GOOS)
		return
	}

	// Copy fonts
	copiedCount := 0
	for _, font := range fontSources {
		destPath := filepath.Join(assetsDir, font.DestName)
		
		// Skip if already exists
		if _, err := os.Stat(destPath); err == nil {
			fmt.Printf("✓ %s already exists\n", font.DestName)
			copiedCount++
			continue
		}
		
		// Check if source exists
		if _, err := os.Stat(font.SourcePath); os.IsNotExist(err) {
			fmt.Printf("⚠ %s not found, skipping\n", font.SourcePath)
			continue
		}
		
		// Copy font
		err := copyFile(font.SourcePath, destPath)
		if err != nil {
			fmt.Printf("✗ Error copying %s: %v\n", font.DestName, err)
			continue
		}
		
		fmt.Printf("✓ Copied %s\n", font.DestName)
		copiedCount++
	}

	fmt.Printf("\nFont setup completed! Copied %d fonts to %s\n", copiedCount, assetsDir)
	
	// Create a font manifest
	createFontManifest(assetsDir)
	
	fmt.Println("\nYou can now run font tests:")
	fmt.Println("  go run examples/unicode-emoji.go")
	fmt.Println("  go run examples/font-loading-test.go")
}

type FontSource struct {
	SourcePath string
	DestName   string
}

func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return err
	}

	return destFile.Sync()
}

func createFontManifest(assetsDir string) {
	manifestPath := filepath.Join(assetsDir, "fonts.json")
	
	manifest := `{
  "fonts": {
    "sans-serif": {
      "regular": "NotoSans-Regular.ttf",
      "bold": "NotoSans-Bold.ttf",
      "fallback": ["LiberationSans-Regular.ttf", "arial.ttf", "Helvetica.ttc"]
    },
    "serif": {
      "regular": "NotoSerif-Regular.ttf",
      "fallback": ["times.ttf", "Times.ttc"]
    },
    "emoji": {
      "color": "NotoColorEmoji.ttf",
      "fallback": ["AppleColorEmoji.ttc", "seguiemj.ttf"]
    },
    "scripts": {
      "arabic": "NotoSansArabic-Regular.ttf",
      "hebrew": "NotoSansHebrew-Regular.ttf",
      "devanagari": "NotoSansDevanagari-Regular.ttf",
      "thai": "NotoSansThai-Regular.ttf"
    }
  },
  "version": "1.0.0",
  "description": "Font manifest for AdvanceGG testing and examples"
}`

	err := os.WriteFile(manifestPath, []byte(manifest), 0644)
	if err != nil {
		fmt.Printf("Warning: Could not create font manifest: %v\n", err)
		return
	}
	
	fmt.Printf("✓ Created font manifest: %s\n", manifestPath)
}
