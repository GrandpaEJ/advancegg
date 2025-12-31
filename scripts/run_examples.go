package main

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	rootDir := "examples"
	var examples []string

	// Walk through examples directory
	err := filepath.Walk(rootDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".go") {
			// Skip generator script and this runner if implicated
			if strings.Contains(path, "blend_modes_gen.go") {
				return nil
			}
			examples = append(examples, path)
		}
		return nil
	})

	if err != nil {
		fmt.Printf("Error walking examples: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Found %d examples to run...\n", len(examples))

	failed := 0
	for _, ex := range examples {
		fmt.Printf("------------------------------------------------\n")
		fmt.Printf("Running %s...\n", ex)

		cmd := exec.Command("go", "run", ex)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err := cmd.Run(); err != nil {
			fmt.Printf("FAIL: %s (%v)\n", ex, err)
			failed++
		} else {
			fmt.Printf("PASS: %s\n", ex)
		}
	}

	fmt.Printf("------------------------------------------------\n")
	if failed == 0 {
		fmt.Println("All examples passed!")
	} else {
		fmt.Printf("%d examples failed.\n", failed)
		os.Exit(1)
	}
}
