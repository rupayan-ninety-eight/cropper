package main

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

func RunCropper(folderPath string, splitWidthStr string) {
	splitWidth, err := strconv.Atoi(splitWidthStr)
	if err != nil || splitWidth <= 0 {
		log.Fatalf("Invalid split width: %s\n", splitWidthStr)
	}

	outputDir := filepath.Join(folderPath, "output")
	err = os.MkdirAll(outputDir, os.ModePerm)
	if err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	err = filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || strings.Contains(path, "/output/") {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(info.Name()))
		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
			return nil
		}

		fmt.Printf("Processing: %s\n", path)
		err = processImage(path, splitWidth, outputDir)
		if err != nil {
			fmt.Printf("  Skipped (%v)\n", err)
		}
		return nil
	})

	if err != nil {
		log.Fatalf("Error walking folder: %v\n", err)
	}
}

func processImage(imagePath string, splitWidth int, outputDir string) error {
	file, err := os.Open(imagePath)
	if err != nil {
		return fmt.Errorf("failed to open image: %v", err)
	}
	defer file.Close()

	img, _, err := image.DecodeConfig(file)
	if err != nil {
		return fmt.Errorf("failed to decode image config: %v", err)
	}

	width := img.Width
	height := img.Height

	if splitWidth*2+2 > width {
		return fmt.Errorf("image width too small (%d)", width)
	}

	baseName := strings.TrimSuffix(filepath.Base(imagePath), filepath.Ext(imagePath))
	leftOutput := filepath.Join(outputDir, baseName+"_right.png")
	rightOutput := filepath.Join(outputDir, baseName+"_left.png")

	// Crop left
	leftCmd := exec.Command("magick", imagePath, "-crop",
		fmt.Sprintf("%dx%d+0+0", splitWidth, height),
		"+repage", leftOutput)

	// Crop right
	rightCmd := exec.Command("magick", imagePath, "-crop",
		fmt.Sprintf("%dx%d+%d+0", splitWidth, height, width-splitWidth),
		"+repage", rightOutput)

	if err := leftCmd.Run(); err != nil {
		return fmt.Errorf("magick crop left failed: %v", err)
	}
	if err := rightCmd.Run(); err != nil {
		return fmt.Errorf("magick crop right failed: %v", err)
	}

	return nil
}
