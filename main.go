package main

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func main() {
	if len(os.Args) != 3 {
		log.Fatalf("Usage: %s <folder_path> <split_width>\n", os.Args[0])
	}

	folderPath := os.Args[1]
	splitWidth, err := strconv.Atoi(os.Args[2])
	if err != nil || splitWidth <= 0 {
		log.Fatalf("Invalid split width: %s\n", os.Args[2])
	}

	err = filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(info.Name()))
		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
			return nil
		}

		fmt.Printf("Processing: %s\n", path)
		err = processImage(path, splitWidth)
		if err != nil {
			fmt.Printf("  Skipped (%v)\n", err)
		}
		return nil
	})

	if err != nil {
		log.Fatalf("Error walking folder: %v\n", err)
	}
}

func processImage(imagePath string, splitWidth int) error {
	file, err := os.Open(imagePath)
	if err != nil {
		return fmt.Errorf("failed to open image: %v", err)
	}
	defer file.Close()

	img, format, err := image.Decode(file)
	if err != nil {
		return fmt.Errorf("failed to decode image: %v", err)
	}

	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	if splitWidth*2+2 > width {
		return fmt.Errorf("image width too small (%d)", width)
	}

	// Crop left and right
	leftRect := image.Rect(0, 0, splitWidth, height)
	rightRect := image.Rect(width-splitWidth, 0, width, height)

	leftImg := cropImage(img, leftRect)
	rightImg := cropImage(img, rightRect)

	base := strings.TrimSuffix(imagePath, filepath.Ext(imagePath))
	saveImage(leftImg, base+"_right"+filepath.Ext(imagePath), format)
	saveImage(rightImg, base+"_left"+filepath.Ext(imagePath), format)

	return nil
}

func cropImage(src image.Image, rect image.Rectangle) image.Image {
	dst := image.NewRGBA(image.Rect(0, 0, rect.Dx(), rect.Dy()))
	for y := 0; y < rect.Dy(); y++ {
		for x := 0; x < rect.Dx(); x++ {
			dst.Set(x, y, src.At(rect.Min.X+x, rect.Min.Y+y))
		}
	}
	return dst
}

func saveImage(img image.Image, filename string, format string) {
	outFile, err := os.Create(filename)
	if err != nil {
		log.Printf("Failed to create file %s: %v", filename, err)
		return
	}
	defer outFile.Close()

	switch strings.ToLower(format) {
	case "png":
		err = png.Encode(outFile, img)
	case "jpeg", "jpg":
		err = jpeg.Encode(outFile, img, nil)
	default:
		log.Printf("Unsupported format: %s", format)
	}

	if err != nil {
		log.Printf("Failed to encode image %s: %v", filename, err)
	}
}
