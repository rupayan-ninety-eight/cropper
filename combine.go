package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

func RunCombiner(folderPath string, direction string) {
	direction = strings.ToLower(direction)
	if direction != "left" && direction != "right" {
		log.Fatalf("Invalid direction: %s (expected 'left' or 'right')", direction)
	}

	files, err := getImageFiles(folderPath)
	if err != nil {
		log.Fatalf("Failed to read image files: %v", err)
	}

	if len(files)%2 != 0 {
		log.Printf("Warning: odd number of images, last one will be skipped\n")
	}

	outputDir := filepath.Join(folderPath, "output")
	err = os.MkdirAll(outputDir, os.ModePerm)
	if err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	for i := 0; i+1 < len(files); i += 2 {
		var left, right string
		if direction == "left" {
			left = files[i]
			right = files[i+1]
		} else {
			left = files[i+1]
			right = files[i]
		}

		outputName := filepath.Base(files[i])
		outputPath := filepath.Join(outputDir, outputName)
		err := combineWithMagick(left, right, outputPath)
		if err != nil {
			fmt.Printf("❌ Failed to combine %s and %s: %v\n", left, right, err)
		} else {
			fmt.Printf("✅ Combined %s + %s → %s\n", left, right, outputPath)
		}
	}
}

func getImageFiles(folder string) ([]string, error) {
	var images []string

	entries, err := os.ReadDir(folder)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		ext := strings.ToLower(filepath.Ext(entry.Name()))
		if ext == ".jpg" || ext == ".jpeg" || ext == ".png" {
			images = append(images, filepath.Join(folder, entry.Name()))
		}
	}

	sort.Strings(images)
	return images, nil
}

func combineWithMagick(leftPath, rightPath, outputPath string) error {
	cmd := exec.Command("magick", leftPath, rightPath, "+append", outputPath)
	return cmd.Run()
}
